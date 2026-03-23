/*
 * @Author: wulnut carepdime@gmail.com
 * @Date: 2026-03-21 22:58:46
 * @LastEditors: wulnut carepdime@gmail.com
 * @LastEditTime: 2026-03-21 23:01:19
 * @FilePath: /lark_cli/test/openapitest/user_manual_test.go
 * @Description: go build -tags=manual ./...
 */

package openapitest

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"lark_cli/internal/auth"
	"lark_cli/internal/config"
	"lark_cli/internal/openapi"
	"lark_cli/internal/session"
)

// configStore matches cmd.configStore — implements auth.ConfigStore.
type configStore struct{}

func (c *configStore) Load(ctx context.Context) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// TestManualQueryCurrentUser calls QueryCurrentUser against the real API using
// ~/.lark/config.json (and env overrides). Requires a valid session for plugin token.
//
//	go test -tags=manual -v -run TestManualQueryCurrentUser ./test/openapitest
func TestManualQueryCurrentUser(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	if err := cfg.ValidateForOpenAPI(); err != nil {
		t.Skipf("OpenAPI config incomplete (set ~/.lark/config.json or LARK_* env): %v", err)
	}

	store := session.NewFileStore(cfg.SessionPath)
	httpClient := &http.Client{Timeout: cfg.HTTPTimeout}
	pluginClient := auth.NewPluginTokenClient(httpClient, cfg.BaseURL, cfg.PluginID, cfg.PluginSecret)
	tokenProvider := auth.NewPluginTokenProvider(&configStore{}, store, pluginClient, cfg.RefreshLeeway)
	client := openapi.NewClient(cfg.BaseURL, httpClient, tokenProvider)

	user, err := client.QueryCurrentUser(ctx, cfg.UserKey)
	if err != nil {
		t.Fatalf("QueryCurrentUser: %v", err)
	}
	if user == nil {
		t.Log("QueryCurrentUser: success, data empty (no user for this user_key)")
		return
	}

	b, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		t.Fatalf("marshal UserInfo: %v", err)
	}
	t.Logf("QueryCurrentUser response:\n%s", string(b))
}
