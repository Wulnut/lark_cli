package openapitest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lark_cli/internal/auth"
	"lark_cli/internal/openapi"
)

type fakeTokenProvider struct {
	forceRefreshCalls int
}

func (f *fakeTokenProvider) GetAuthContext(ctx context.Context) (*auth.AuthContext, error) {
	return &auth.AuthContext{
		UserKey:     "ou_test",
		PluginToken: "p_test_token",
	}, nil
}

func (f *fakeTokenProvider) ForceRefresh(ctx context.Context) (*auth.AuthContext, error) {
	f.forceRefreshCalls++
	return &auth.AuthContext{
		UserKey:     "ou_test",
		PluginToken: "p_refreshed",
	}, nil
}

func TestDoJSON_Non2xxWithoutErrCode_IncludesBodyDetailInErrorString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("upstream unavailable"))
	}))
	defer server.Close()

	tp := &fakeTokenProvider{}
	client := openapi.NewClient(server.URL, server.Client(), tp)

	err := client.DoJSON(context.Background(), &openapi.Request{Method: http.MethodGet, Path: "/v1/test"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "http_status=503") {
		t.Fatalf("error should include http status, got: %v", err)
	}
	if !strings.Contains(err.Error(), "upstream unavailable") {
		t.Fatalf("error should include response detail, got: %v", err)
	}
}

func TestDoJSON_AttemptLimitUsesMaxAttemptsSemantics(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"err_code": 10022,
			"err_msg":  "Check Token Failed",
			"data":     map[string]any{},
		})
	}))
	defer server.Close()

	tp := &fakeTokenProvider{}
	client := openapi.NewClient(server.URL, server.Client(), tp)
	client.MaxAttempts = 1

	err := client.DoJSON(context.Background(), &openapi.Request{Method: http.MethodGet, Path: "/v1/test"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if callCount != 1 {
		t.Fatalf("expected 1 attempt due to MaxAttempts cap, got %d", callCount)
	}
	if tp.forceRefreshCalls != 1 {
		t.Fatalf("expected force refresh once on auth error, got %d", tp.forceRefreshCalls)
	}
}
