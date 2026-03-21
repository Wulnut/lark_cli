package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"lark_cli/internal/config"
)

// BuildConfigFingerprint generates a SHA256 fingerprint of the config.
// This fingerprint is used to detect config changes and invalidate cached tokens.
func BuildConfigFingerprint(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}

	raw := strings.Join([]string{
		cfg.UserKey,
		cfg.PluginID,
		cfg.PluginSecret,
		cfg.BaseURL,
	}, "|")

	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
