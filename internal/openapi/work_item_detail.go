package openapi

import (
	"context"
	"fmt"
)

// WorkItemDetailResponse is the response shape for POST /open_api/{project_key}/work_item/{work_item_type_key}/query.
type WorkItemDetailResponse struct {
	ErrCode int              `json:"err_code"`
	ErrMsg  string           `json:"err_msg"`
	Err     any              `json:"err"`
	Data    []map[string]any `json:"data"`
}

// GetWorkItemDetail calls the work-item detail query endpoint for one project and work-item type.
func (c *Client) GetWorkItemDetail(ctx context.Context, projectKey, workItemTypeKey string, payload map[string]any) (*WorkItemDetailResponse, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("project_key is required")
	}
	if workItemTypeKey == "" {
		return nil, fmt.Errorf("work_item_type_key is required")
	}
	if payload == nil {
		return nil, fmt.Errorf("payload is required")
	}

	path := fmt.Sprintf("open_api/%s/work_item/%s/query", projectKey, workItemTypeKey)
	var resp WorkItemDetailResponse
	err := c.DoJSON(ctx, &Request{
		Method: "POST",
		Path:   path,
		Body:   payload,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
