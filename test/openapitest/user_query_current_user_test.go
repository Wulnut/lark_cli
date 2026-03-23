package openapitest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lark_cli/internal/openapi"
)

func TestQueryCurrentUser_FallbackWhenExactKeyMiss(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"err_code": 0,
			"err_msg":  "",
			"err":      map[string]any{},
			"data": []map[string]any{
				{
					"user_id":   1,
					"name_cn":   "测试用户",
					"name_en":   "Test User",
					"out_id":    "",
					"name":      map[string]string{"default": "测试用户", "zh_cn": "测试用户", "en_us": "Test User"},
					"user_key":  "different_key",
					"username":  "tester",
					"email":     "test@example.com",
					"avatar_url": "",
					"status":    "activated",
				},
			},
		})
	}))
	defer server.Close()

	client := openapi.NewClient(server.URL, server.Client(), &fakeTokenProvider{})

	user, err := client.QueryCurrentUser(context.Background(), "requested_key")
	if err != nil {
		t.Fatalf("QueryCurrentUser error: %v", err)
	}
	if user == nil {
		t.Fatalf("expected fallback user, got nil")
	}
	if user.UserKey != "different_key" {
		t.Fatalf("unexpected user key: %s", user.UserKey)
	}
}
