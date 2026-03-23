package openapitest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"lark_cli/internal/openapi"
)

func TestClient_GetWorkItemDetail_UsesQueryEndpoint(t *testing.T) {
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

	_, _ = client.GetWorkItemDetail(context.Background(), "p_demo", "issue", map[string]any{"work_item_ids": []int{101}})

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/open_api/p_demo/work_item/issue/query" {
		t.Fatalf("path = %s, want /open_api/p_demo/work_item/issue/query", gotPath)
	}
}

func TestClient_GetWorkItemDetail_SendsIDsAndFields(t *testing.T) {
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
			"data": []map[string]any{{
				"id":   101,
				"name": "Issue 101",
			}},
		})
	}))
	defer server.Close()

	client := openapi.NewClient(server.URL, server.Client(), &fakeTokenProvider{})

	payload := map[string]any{
		"work_item_ids": []int{101, 102},
		"fields":        []string{"name", "description"},
	}
	resp, err := client.GetWorkItemDetail(context.Background(), "p_demo", "issue", payload)
	if err != nil {
		t.Fatalf("GetWorkItemDetail: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 row, got %d", len(resp.Data))
	}

	gotIDs, ok := gotBody["work_item_ids"].([]any)
	if !ok {
		t.Fatalf("work_item_ids type = %T, want []any", gotBody["work_item_ids"])
	}
	if !reflect.DeepEqual(gotIDs, []any{float64(101), float64(102)}) {
		t.Fatalf("work_item_ids = %#v", gotIDs)
	}
	gotFields, ok := gotBody["fields"].([]any)
	if !ok {
		t.Fatalf("fields type = %T, want []any", gotBody["fields"])
	}
	if !reflect.DeepEqual(gotFields, []any{"name", "description"}) {
		t.Fatalf("fields = %#v", gotFields)
	}
}
