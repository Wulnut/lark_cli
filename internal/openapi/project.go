package openapi

import (
	"context"
	"fmt"
)

// ListProjectsResponse is the response shape for POST /open_api/projects.
type ListProjectsResponse struct {
	ErrCode int      `json:"err_code"`
	ErrMsg  string   `json:"err_msg"`
	Err     any      `json:"err"`
	Data    []string `json:"data"` // list of project_key
}

// ProjectDetail represents the detail of a single project space.
type ProjectDetail struct {
	ProjectKey     string   `json:"project_key"`
	Name           string   `json:"name"`
	SimpleName     string   `json:"simple_name"`
	Administrators []string `json:"administrators"`
}

// GetProjectDetailResponse is the response shape for POST /open_api/projects/detail.
type GetProjectDetailResponse struct {
	ErrCode int                       `json:"err_code"`
	ErrMsg  string                    `json:"err_msg"`
	Err     any                       `json:"err"`
	Data    map[string]ProjectDetail  `json:"data"` // key = project_key
}

// ListProjects calls POST /open_api/projects to get the list of project_keys.
// Returns a list of project_key strings, or error on network failure.
func (c *Client) ListProjects(ctx context.Context, userKey, order string) ([]string, error) {
	body := map[string]string{"user_key": userKey}
	if order != "" {
		body["order"] = order
	}
	var resp ListProjectsResponse
	err := c.DoJSON(ctx, &Request{
		Method: "POST",
		Path:   "open_api/projects",
		Body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetProjectDetails calls POST /open_api/projects/detail to get project details.
// projectKeys and simpleNames are optional but cannot both be empty.
// Returns a map of project_key -> ProjectDetail.
func (c *Client) GetProjectDetails(ctx context.Context, userKey string, projectKeys, simpleNames []string) (map[string]ProjectDetail, error) {
	body := map[string]any{"user_key": userKey}
	if len(projectKeys) > 0 {
		body["project_keys"] = projectKeys
	}
	if len(simpleNames) > 0 {
		body["simple_names"] = simpleNames
	}
	var resp GetProjectDetailResponse
	err := c.DoJSON(ctx, &Request{
		Method: "POST",
		Path:   "open_api/projects/detail",
		Body:   body,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListProjectsWithDetails fetches project list and then enriches with details.
// It returns a list of ProjectDetail sorted by the order returned from ListProjects.
func (c *Client) ListProjectsWithDetails(ctx context.Context, userKey, order string) ([]ProjectDetail, error) {
	keys, err := c.ListProjects(ctx, userKey, order)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	if len(keys) == 0 {
		return nil, nil
	}

	details, err := c.GetProjectDetails(ctx, userKey, keys, nil)
	if err != nil {
		return nil, fmt.Errorf("get project details: %w", err)
	}

	// Preserve order from ListProjects
	result := make([]ProjectDetail, 0, len(keys))
	for _, k := range keys {
		if d, ok := details[k]; ok {
			result = append(result, d)
		}
	}
	return result, nil
}
