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

// QueryUsers calls POST /open_api/user/query with a list of user_keys.
// Returns a map of user_key -> UserInfo for all found users.
func (c *Client) QueryUsers(ctx context.Context, userKeys []string) (map[string]*UserInfo, error) {
	if len(userKeys) == 0 {
		return nil, nil
	}
	var resp QueryUserResponse
	err := c.DoJSON(ctx, &Request{
		Method: "POST",
		Path:   "open_api/user/query",
		Body:   map[string][]string{"user_keys": userKeys},
	}, &resp)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*UserInfo, len(resp.Data))
	for i := range resp.Data {
		result[resp.Data[i].UserKey] = &resp.Data[i]
	}
	return result, nil
}

// QueryCurrentUser calls POST /open_api/user/query with a single user_key.
// Returns the first user in the data list, or nil if the result is empty.
func (c *Client) QueryCurrentUser(ctx context.Context, userKey string) (*UserInfo, error) {
	if userKey == "" {
		return nil, fmt.Errorf("userKey is required")
	}
	users, err := c.QueryUsers(ctx, []string{userKey})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	if u, ok := users[userKey]; ok && u != nil {
		return u, nil
	}
	for _, u := range users {
		if u != nil {
			return u, nil
		}
	}
	return nil, nil
}
