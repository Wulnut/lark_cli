package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PluginTokenClient struct {
	httpClient   *http.Client
	baseURL      string
	pluginID     string
	pluginSecret string
}

type pluginTokenRequest struct {
	PluginID     string `json:"plugin_id"`
	PluginSecret string `json:"plugin_secret"`
	Type         int    `json:"type"`
}

type pluginTokenResponse struct {
	Data struct {
		ExpireTime int    `json:"expire_time"`
		Token      string `json:"token"`
	} `json:"data"`
	Error struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
}

func NewPluginTokenClient(httpClient *http.Client, baseURL, pluginID, pluginSecret string) *PluginTokenClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &PluginTokenClient{
		httpClient:   httpClient,
		baseURL:      baseURL,
		pluginID:     pluginID,
		pluginSecret: pluginSecret,
	}
}

func (c *PluginTokenClient) Fetch(ctx context.Context) (string, time.Duration, error) {
	reqBody := pluginTokenRequest{
		PluginID:     c.pluginID,
		PluginSecret: c.pluginSecret,
		Type:         0,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + "/open_api/authen/plugin_token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBytes))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp pluginTokenResponse
	if err := json.Unmarshal(respBytes, &tokenResp); err != nil {
		return "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenResp.Error.Code != 0 {
		return "", 0, fmt.Errorf("API error %d: %s", tokenResp.Error.Code, tokenResp.Error.Msg)
	}

	expiresIn := time.Duration(tokenResp.Data.ExpireTime) * time.Second
	return tokenResp.Data.Token, expiresIn, nil
}
