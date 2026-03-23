package openapi

import (
	"context"
	"fmt"
)

// WorkItemListResponse is the response shape for POST /open_api/{project_key}/work_item/filter.
type WorkItemListResponse struct {
	ErrCode    int              `json:"err_code"`
	ErrMsg     string           `json:"err_msg"`
	Err        any              `json:"err"`
	Data       []map[string]any `json:"data"`
	Pagination map[string]any   `json:"pagination"`
}

// ListWorkItems calls POST /open_api/{project_key}/work_item/filter.
func (c *Client) ListWorkItems(ctx context.Context, projectKey string, payload map[string]any) (*WorkItemListResponse, error) {
	if projectKey == "" {
		return nil, fmt.Errorf("projectKey is required")
	}
	if payload == nil {
		return nil, fmt.Errorf("payload is required")
	}

	var resp WorkItemListResponse
	err := c.DoJSON(ctx, &Request{
		Method: "POST",
		Path:   fmt.Sprintf("open_api/%s/work_item/filter", projectKey),
		Body:   payload,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
