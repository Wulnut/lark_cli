package openapitest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lark_cli/internal/openapi"
)

func TestSearchWorkItems_SendsSearchGroupPayload(t *testing.T) {
	var gotPath string
	var gotMethod string
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"err_code": 0,
			"err_msg":  "",
			"data":     []map[string]any{{"id": 1, "name": "Issue A"}},
		})
	}))
	defer server.Close()

	tp := &fakeTokenProvider{}
	client := openapi.NewClient(server.URL, server.Client(), tp)

	payload := map[string]any{
		"search_group": map[string]any{
			"conjunction": "AND",
			"search_params": []map[string]any{{
				"param_key": "people",
				"operator":  "HAS ANY OF",
				"value":     []string{"ou_test"},
			}},
		},
		"page_size": 20,
		"page_num":  1,
	}

	resp, err := client.SearchWorkItems(context.Background(), "p_demo", "issue", payload)
	if err != nil {
		t.Fatalf("SearchWorkItems: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 row, got %d", len(resp.Data))
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/open_api/p_demo/work_item/issue/search/params" {
		t.Fatalf("path = %s", gotPath)
	}
	if _, ok := gotBody["search_group"]; !ok {
		t.Fatalf("search_group missing in body")
	}
}

func TestSearchWorkItems_RequiresSearchGroup(t *testing.T) {
	tp := &fakeTokenProvider{}
	client := openapi.NewClient("https://example.invalid", nil, tp)

	_, err := client.SearchWorkItems(context.Background(), "p_demo", "issue", map[string]any{"page_size": 20})
	if err == nil {
		t.Fatal("expected error when search_group is missing")
	}
}
