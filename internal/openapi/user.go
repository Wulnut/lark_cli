package openapi

import (
	"context"
	"fmt"
)

// LocaleName represents the multilingual name object in user responses.
type LocaleName struct {
	Default string `json:"default"`
	EnUS    string `json:"en_us"`
	ZhCN    string `json:"zh_cn"`
}

// UserInfo represents a user record returned by POST /open_api/user/query.
type UserInfo struct {
	UserID    int64      `json:"user_id"`
	NameCn    string     `json:"name_cn"`
	NameEn    string     `json:"name_en"`
	OutID     string     `json:"out_id"`
	Name      LocaleName `json:"name"`
	UserKey   string     `json:"user_key"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	AvatarURL string     `json:"avatar_url"`
	Status    string     `json:"status"`
}

// QueryUserResponse is the response shape for POST /open_api/user/query.
type QueryUserResponse struct {
	ErrCode int        `json:"err_code"`
	ErrMsg  string     `json:"err_msg"`
	Err     any        `json:"err"`
	Data    []UserInfo `json:"data"`
}

// QueryCurrentUser calls POST /open_api/user/query with a single user_key.
// Returns the first user in the data list, or nil if the result is empty.
// Returns error only on network failure or unexpected response shape.
func (c *Client) QueryCurrentUser(ctx context.Context, userKey string) (*UserInfo, error) {
	if userKey == "" {
		return nil, fmt.Errorf("userKey is required")
	}
	var resp QueryUserResponse
	err := c.DoJSON(ctx, &Request{
		Method: "POST",
		Path:   "open_api/user/query",
		Body:   map[string][]string{"user_keys": {userKey}},
	}, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, nil // no user found, not an error
	}
	return &resp.Data[0], nil
}
