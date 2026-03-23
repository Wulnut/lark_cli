package openapi

import (
	"context"
	"fmt"
)

// SearchWorkItemsResponse is the response shape for
// POST /open_api/{project_key}/work_item/{work_item_type_key}/search/params.
type SearchWorkItemsResponse struct {
	ErrCode    int              `json:"err_code"`
	ErrMsg     string           `json:"err_msg"`
	Err        any              `json:"err"`
	Data       []map[string]any `json:"data"`
	Pagination map[string]any   `json:"pagination"`
}

// SearchWorkItems calls complex work-item search endpoint for one project and work-item type.
func (c *Client) SearchWorkItems(ctx context.Context, projectKey, workItemTypeKey string, payload map[string]any) (*SearchWorkItemsResponse, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("project_key is required")
	}
	if workItemTypeKey == "" {
		return nil, fmt.Errorf("work_item_type_key is required")
	}
	if payload == nil {
		return nil, fmt.Errorf("payload is required")
	}
	if _, ok := payload["search_group"]; !ok {
		return nil, fmt.Errorf("search_group is required")
	}

	path := fmt.Sprintf("open_api/%s/work_item/%s/search/params", projectKey, workItemTypeKey)
	var resp SearchWorkItemsResponse
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
