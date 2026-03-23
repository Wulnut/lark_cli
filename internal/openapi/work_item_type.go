package openapi

import (
	"context"
	"fmt"
)

// WorkItemType represents a work item type in a project space.
type WorkItemType struct {
	TypeKey              string `json:"type_key"`
	Name                 string `json:"name"`
	IsDisable           int    `json:"is_disable"`
	APIName             string `json:"api_name"`
	EnableModelResourceLib bool `json:"enable_model_resource_lib"`
}

// ListWorkItemTypesResponse is the response shape for GET /open_api/{project_key}/work_item/all-types.
type ListWorkItemTypesResponse struct {
	ErrCode int            `json:"err_code"`
	ErrMsg  string         `json:"err_msg"`
	Err     any            `json:"err"`
	Data    []WorkItemType `json:"data"`
}

// ListWorkItemTypes calls GET /open_api/{project_key}/work_item/all-types.
// Returns all work item types for the given project.
func (c *Client) ListWorkItemTypes(ctx context.Context, projectKey string) ([]WorkItemType, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("project_key is required")
	}
	path := fmt.Sprintf("open_api/%s/work_item/all-types", projectKey)
	var resp ListWorkItemTypesResponse
	err := c.DoJSON(ctx, &Request{
		Method: "GET",
		Path:   path,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
