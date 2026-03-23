package query

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// BuildInput is the high-level query input from CLI flags and optional raw JSON.
type BuildInput struct {
	CurrentUserKey   string
	Me               bool
	Persons          []string
	Statuses         []string
	CreatedFrom      string
	CreatedTo        string
	UpdatedFrom      string
	UpdatedTo        string
	Fields           []string // key=value
	RawSearchGroupJSON string
	RawOnly          bool
}

// BuildOutput contains the final search_group and optional warnings.
type BuildOutput struct {
	SearchGroup map[string]any
	Warnings    []string
}

func BuildSearchGroup(in BuildInput) (*BuildOutput, error) {
	rawGroup, err := parseRawSearchGroup(in.RawSearchGroupJSON)
	if err != nil {
		return nil, err
	}
	structured, err := buildStructuredSearchGroup(in)
	if err != nil {
		return nil, err
	}

	if in.RawOnly {
		if rawGroup == nil {
			return nil, fmt.Errorf("--raw-only requires --search-group-json")
		}
		return &BuildOutput{SearchGroup: rawGroup}, nil
	}

	if rawGroup != nil && structured != nil {
		return &BuildOutput{SearchGroup: map[string]any{
			"conjunction":  "AND",
			"search_params": []map[string]any{},
			"search_groups": []map[string]any{rawGroup, structured},
		}, Warnings: []string{"raw search_group merged with compiled flags by AND"}}, nil
	}
	if rawGroup != nil {
		return &BuildOutput{SearchGroup: rawGroup}, nil
	}
	if structured != nil {
		return &BuildOutput{SearchGroup: structured}, nil
	}

	return nil, fmt.Errorf("no query conditions provided")
}

func parseRawSearchGroup(raw string) (map[string]any, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var g map[string]any
	if err := json.Unmarshal([]byte(raw), &g); err != nil {
		return nil, fmt.Errorf("invalid search-group-json: %w", err)
	}
	if len(g) == 0 {
		return nil, fmt.Errorf("search_group json cannot be empty")
	}
	return g, nil
}

func buildStructuredSearchGroup(in BuildInput) (map[string]any, error) {
	params := make([]map[string]any, 0)

	persons := make([]string, 0, len(in.Persons)+1)
	seen := map[string]struct{}{}
	if in.Me {
		if strings.TrimSpace(in.CurrentUserKey) == "" {
			return nil, fmt.Errorf("--me requires logged-in user key")
		}
		seen[in.CurrentUserKey] = struct{}{}
		persons = append(persons, in.CurrentUserKey)
	}
	for _, p := range in.Persons {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		persons = append(persons, p)
	}
	if len(persons) > 0 {
		params = append(params, map[string]any{
			"param_key": "people",
			"operator":  "HAS ANY OF",
			"value":     persons,
		})
	}

	if len(in.Statuses) > 0 {
		vals := make([]string, 0, len(in.Statuses))
		for _, s := range in.Statuses {
			s = strings.TrimSpace(s)
			if s != "" {
				vals = append(vals, s)
			}
		}
		if len(vals) > 0 {
			params = append(params, map[string]any{
				"param_key": "work_item_status",
				"operator":  "HAS ANY OF",
				"value":     vals,
			})
		}
	}

	var createdFromMs, createdToMs int64
	var hasCreatedFrom, hasCreatedTo bool
	if in.CreatedFrom != "" {
		ms, err := parseDateTimeToMillis(in.CreatedFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid --created-from: %w", err)
		}
		hasCreatedFrom = true
		createdFromMs = ms
		params = append(params, map[string]any{"param_key": "created_at", "operator": ">=", "value": ms})
	}
	if in.CreatedTo != "" {
		ms, err := parseDateTimeToMillis(in.CreatedTo)
		if err != nil {
			return nil, fmt.Errorf("invalid --created-to: %w", err)
		}
		hasCreatedTo = true
		createdToMs = ms
		params = append(params, map[string]any{"param_key": "created_at", "operator": "<=", "value": ms})
	}
	if hasCreatedFrom && hasCreatedTo && createdFromMs > createdToMs {
		return nil, fmt.Errorf("--created-from cannot be later than --created-to")
	}

	var updatedFromMs, updatedToMs int64
	var hasUpdatedFrom, hasUpdatedTo bool
	if in.UpdatedFrom != "" {
		ms, err := parseDateTimeToMillis(in.UpdatedFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid --updated-from: %w", err)
		}
		hasUpdatedFrom = true
		updatedFromMs = ms
		params = append(params, map[string]any{"param_key": "updated_at", "operator": ">=", "value": ms})
	}
	if in.UpdatedTo != "" {
		ms, err := parseDateTimeToMillis(in.UpdatedTo)
		if err != nil {
			return nil, fmt.Errorf("invalid --updated-to: %w", err)
		}
		hasUpdatedTo = true
		updatedToMs = ms
		params = append(params, map[string]any{"param_key": "updated_at", "operator": "<=", "value": ms})
	}
	if hasUpdatedFrom && hasUpdatedTo && updatedFromMs > updatedToMs {
		return nil, fmt.Errorf("--updated-from cannot be later than --updated-to")
	}

	for _, f := range in.Fields {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --field format %q: expected key=value", f)
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" || v == "" {
			return nil, fmt.Errorf("invalid --field format %q: key and value must be non-empty", f)
		}
		params = append(params, map[string]any{
			"param_key": k,
			"operator":  "=",
			"value":     v,
		})
	}

	if len(params) == 0 {
		return nil, nil
	}
	if len(params) > 50 {
		return nil, fmt.Errorf("too many search conditions: %d (max 50)", len(params))
	}
	return map[string]any{
		"conjunction":  "AND",
		"search_params": params,
		"search_groups": []map[string]any{},
	}, nil
}

func parseDateTimeToMillis(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty datetime")
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UnixMilli(), nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.UnixMilli(), nil
	}
	return 0, fmt.Errorf("unsupported format %q, use RFC3339 or YYYY-MM-DD", s)
}
