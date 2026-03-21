package session

import "time"

const CurrentVersion = 1

type Session struct {
	Version int `json:"version"`

	LoginType string `json:"login_type"`

	// ConfigFingerprint is used to detect config changes and invalidate cached token.
	ConfigFingerprint string `json:"config_fingerprint,omitempty"`

	PluginAccessToken          string    `json:"plugin_access_token,omitempty"`
	PluginAccessTokenExpiresAt time.Time `json:"plugin_access_token_expires_at,omitzero"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsEmpty returns true if the session contains no plugin token.
func (s *Session) IsEmpty() bool {
	return s == nil || s.PluginAccessToken == ""
}

// IsValid returns true if the session has a valid (non-expired) plugin token
// that matches the given config fingerprint.
func (s *Session) IsValid(nowUnix int64, fingerprint string) bool {
	if s == nil {
		return false
	}
	if s.PluginAccessToken == "" {
		return false
	}
	if s.ConfigFingerprint != fingerprint {
		return false
	}
	// Check if token is expired (now >= expires_at)
	if nowUnix >= s.PluginAccessTokenExpiresAt.Unix() {
		return false
	}
	return true
}
