package querytest

import (
	"testing"

	"lark_cli/internal/query"
)

func TestBuildSearchGroup_FromFlags(t *testing.T) {
	out, err := query.BuildSearchGroup(query.BuildInput{
		CurrentUserKey: "ou_me",
		Me:             true,
		Statuses:       []string{"doing", "todo"},
		CreatedFrom:    "2026-01-01",
		CreatedTo:      "2026-01-31",
		Fields:         []string{"priority=P0", "custom_tag=backend"},
	})
	if err != nil {
		t.Fatalf("BuildSearchGroup: %v", err)
	}

	if out.SearchGroup["conjunction"] != "AND" {
		t.Fatalf("expected conjunction AND, got %v", out.SearchGroup["conjunction"])
	}

	params, ok := out.SearchGroup["search_params"].([]map[string]any)
	if !ok {
		t.Fatalf("search_params type mismatch: %T", out.SearchGroup["search_params"])
	}
	if len(params) == 0 {
		t.Fatalf("expected search params")
	}
}

func TestBuildSearchGroup_MergeRawAndFlags(t *testing.T) {
	out, err := query.BuildSearchGroup(query.BuildInput{
		Persons: []string{"ou_a"},
		RawSearchGroupJSON: `{
			"conjunction":"AND",
			"search_params":[{"param_key":"work_item_status","operator":"HAS ANY OF","value":["doing"]}]
		}`,
	})
	if err != nil {
		t.Fatalf("BuildSearchGroup: %v", err)
	}

	groups, ok := out.SearchGroup["search_groups"].([]map[string]any)
	if !ok || len(groups) == 0 {
		t.Fatalf("expected merged nested search_groups, got %T", out.SearchGroup["search_groups"])
	}
}

func TestBuildSearchGroup_RawOnly(t *testing.T) {
	raw := `{"conjunction":"AND","search_params":[{"param_key":"people","operator":"HAS ANY OF","value":["ou_x"]}]}`
	out, err := query.BuildSearchGroup(query.BuildInput{RawSearchGroupJSON: raw, RawOnly: true})
	if err != nil {
		t.Fatalf("BuildSearchGroup: %v", err)
	}
	if out.SearchGroup["conjunction"] != "AND" {
		t.Fatalf("unexpected output: %v", out.SearchGroup)
	}
}

func TestBuildSearchGroup_InvalidFieldFormat(t *testing.T) {
	_, err := query.BuildSearchGroup(query.BuildInput{Fields: []string{"invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid field format")
	}
}
