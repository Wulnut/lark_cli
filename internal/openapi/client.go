package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"lark_cli/internal/auth"
	openapierrors "lark_cli/internal/openapi/errors"
)

// Request represents an OpenAPI request.
type Request struct {
	Method  string
	Path    string
	Body    any
	Headers map[string]string
}

// CommonResponse is the standard Feishu Project API response structure.
type CommonResponse struct {
	ErrCode int             `json:"err_code"`
	ErrMsg  string          `json:"err_msg"`
	Data    json.RawMessage `json:"data"`
}

// Client is the unified OpenAPI client with automatic auth and retry handling.
type Client struct {
	BaseURL       string
	HTTPClient    *http.Client
	TokenProvider auth.PluginTokenProvider
	MaxAttempts   int
}

// NewClient creates a new OpenAPI client.
func NewClient(baseURL string, httpClient *http.Client, tokenProvider auth.PluginTokenProvider) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		BaseURL:       strings.TrimRight(baseURL, "/"),
		HTTPClient:    httpClient,
		TokenProvider: tokenProvider,
		MaxAttempts:   4, // Total attempts: 1 initial + up to 3 retries
	}
}

// DoJSON performs an HTTP request with automatic auth injection and retry logic.
func (c *Client) DoJSON(ctx context.Context, req *Request, out any) error {
	var lastErr error

	for attempt := 0; attempt < c.MaxAttempts; attempt++ {
		// 1. Get auth context (user_key + plugin_token)
		authCtx, err := c.TokenProvider.GetAuthContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth context: %w", err)
		}

		// 2. Build the HTTP request
		httpReq, err := c.buildRequest(ctx, req, authCtx)
		if err != nil {
			return fmt.Errorf("failed to build request: %w", err)
		}

		// 3. Execute the request
		httpResp, body, err := c.do(httpReq)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		// 4. Parse and normalize the response
		code, _, apiErr := c.parseAndNormalize(httpResp.StatusCode, body)
		if apiErr == nil {
			// Success
			if out != nil && len(body) > 0 {
				if err := json.Unmarshal(body, out); err != nil {
					return fmt.Errorf("failed to unmarshal response: %w", err)
				}
			}
			return nil
		}

		lastErr = apiErr
		policy := openapierrors.LookupPolicy(code)

		// 5. Handle token refresh errors
		if policy.RefreshToken && attempt < policy.MaxRetry {
			if _, err := c.TokenProvider.ForceRefresh(ctx); err != nil {
				return fmt.Errorf("failed to refresh token: %w", err)
			}
			continue
		}

		// 6. Handle retryable errors with backoff
		if policy.Retryable && attempt < policy.MaxRetry {
			sleepDuration := backoffDuration(policy.Category, attempt)
			if sleepDuration > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(sleepDuration):
				}
			}
			continue
		}

		// 7. Non-retryable error or max retries exceeded
		return apiErr
	}

	return lastErr
}

func (c *Client) buildRequest(ctx context.Context, req *Request, authCtx *auth.AuthContext) (*http.Request, error) {
	var bodyReader io.Reader
	if req.Body != nil {
		payload, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	fullURL := c.BaseURL + "/" + strings.TrimLeft(req.Path, "/")
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set standard headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Plugin-Token", authCtx.PluginToken)
	httpReq.Header.Set("X-User-Key", authCtx.UserKey)

	// Set custom headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	return httpReq, nil
}

func (c *Client) do(req *http.Request) (*http.Response, []byte, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("http request failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	return resp, body, nil
}

func (c *Client) parseAndNormalize(httpStatus int, body []byte) (int, string, error) {
	var parsed CommonResponse
	if err := json.Unmarshal(body, &parsed); err == nil {
		if parsed.ErrCode != 0 {
			return parsed.ErrCode, parsed.ErrMsg,
				openapierrors.NormalizeError(httpStatus, body, parsed.ErrCode, parsed.ErrMsg)
		}
		if httpStatus >= 200 && httpStatus < 300 {
			return 0, "", nil
		}
	}

	// Non-2xx HTTP status without parsed error code
	if httpStatus < 200 || httpStatus >= 300 {
		msg := fmt.Sprintf("http status %d", httpStatus)
		if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
			msg = fmt.Sprintf("http status %d: %s", httpStatus, trimmed)
		}
		return 0, "", openapierrors.NormalizeError(httpStatus, body, 0, msg)
	}

	return 0, "", nil
}

// backoffDuration returns the sleep duration for a given error category and attempt.
func backoffDuration(category openapierrors.Category, attempt int) time.Duration {
	switch category {
	case openapierrors.CategoryRateLimit:
		switch attempt {
		case 0:
			return 200 * time.Millisecond
		case 1:
			return 500 * time.Millisecond
		default:
			return 1 * time.Second
		}
	case openapierrors.CategoryServer:
		switch attempt {
		case 0:
			return 300 * time.Millisecond
		default:
			return 800 * time.Millisecond
		}
	default:
		return 0
	}
}
