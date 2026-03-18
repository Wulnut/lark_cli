package session

import "time"

const CurrentVersion = 1

type Session struct {
	Version int `json:"version"`

	LoginType string `json:"login_type"`
	UserKey   string `json:"user_key"`

	PluginAccessToken          string    `json:"plugin_access_token,omitempty"`
	PluginAccessTokenExpiresAt time.Time `json:"plugin_access_token_expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
