package openapitest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lark_cli/internal/openapi"
)

func TestClient_ListWorkItems_UsesFilterEndpoint(t *testing.T) {
	var gotPath string
	var gotMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"err_code": 0,
			"err_msg":  "",
			"data":     []map[string]any{},
		})
	}))
	defer server.Close()

	client := openapi.NewClient(server.URL, server.Client(), &fakeTokenProvider{})

	_, _ = client.ListWorkItems(context.Background(), "p_demo", map[string]any{
		"work_item_type_keys": []string{"issue"},
		"page_num":            1,
		"page_size":           20,
	})

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/open_api/p_demo/work_item/filter" {
		t.Fatalf("path = %s, want /open_api/p_demo/work_item/filter", gotPath)
	}
}

func TestClient_ListWorkItems_SendsTypeKeysAndPagination(t *testing.T) {
	var gotBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"err_code": 0,
			"err_msg":  "",
			"data":     []map[string]any{},
		})
	}))
	defer server.Close()

	client := openapi.NewClient(server.URL, server.Client(), &fakeTokenProvider{})

	_, _ = client.ListWorkItems(context.Background(), "p_demo", map[string]any{
		"work_item_type_keys": []string{"issue", "story"},
		"page_num":            2,
		"page_size":           50,
	})

	if _, ok := gotBody["work_item_type_keys"]; !ok {
		t.Fatalf("work_item_type_keys missing in body")
	}
	if gotBody["page_num"] != float64(2) {
		t.Fatalf("page_num = %v, want 2", gotBody["page_num"])
	}
	if gotBody["page_size"] != float64(50) {
		t.Fatalf("page_size = %v, want 50", gotBody["page_size"])
	}
}

func TestClient_ListWorkItems_ParsesDataAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"err_code": 0,
			"err_msg":  "",
			"data": []map[string]any{
				{"work_item_id": 101, "name": "Issue A"},
				{"work_item_id": 102, "name": "Issue B"},
			},
			"pagination": map[string]any{
				"page_num":  3,
				"page_size": 20,
				"total":     88,
				"has_more":  true,
			},
		})
	}))
	defer server.Close()

	client := openapi.NewClient(server.URL, server.Client(), &fakeTokenProvider{})

	resp, err := client.ListWorkItems(context.Background(), "p_demo", map[string]any{
		"work_item_type_keys": []string{"issue"},
		"page_num":            3,
		"page_size":           20,
	})
	if err != nil {
		t.Fatalf("ListWorkItems: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 work items, got %d", len(resp.Data))
	}
	if resp.Pagination["page_num"] != float64(3) {
		t.Fatalf("pagination.page_num = %v, want 3", resp.Pagination["page_num"])
	}
	if resp.Pagination["has_more"] != true {
		t.Fatalf("pagination.has_more = %v, want true", resp.Pagination["has_more"])
	}
}
