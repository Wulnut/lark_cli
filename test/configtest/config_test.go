package configtest

import (
	"os"
	"path/filepath"
	"testing"

	"lark_cli/internal/config"
)

func setupTempHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("LARK_BASE_URL", "")
	t.Setenv("LARK_PLUGIN_ID", "")
	t.Setenv("LARK_PLUGIN_SECRET", "")
	t.Setenv("LARK_SESSION_PATH", "")
	return home
}

func writeConfigFile(t *testing.T, home string, body string) {
	t.Helper()
	dir := filepath.Join(home, ".lark")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll error: %v", err)
	}
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
}

func TestLoad_DefaultsWithoutEnvAndConfigFile(t *testing.T) {
	home := setupTempHome(t)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BaseURL != config.DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, config.DefaultBaseURL)
	}
	wantSessionPath := filepath.Join(home, ".lark", "session.json")
	if cfg.SessionPath != wantSessionPath {
		t.Errorf("SessionPath = %q, want %q", cfg.SessionPath, wantSessionPath)
	}
}

func TestLoad_UsesConfigFileValues(t *testing.T) {
	home := setupTempHome(t)
	writeConfigFile(t, home, `{
  "base_url": "https://file.example.com",
  "plugin_id": "file_id",
  "plugin_secret": "file_secret",
  "session_path": "/tmp/file-session.json"
}`)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BaseURL != "https://file.example.com" {
		t.Errorf("BaseURL = %q", cfg.BaseURL)
	}
	if cfg.PluginID != "file_id" {
		t.Errorf("PluginID = %q", cfg.PluginID)
	}
	if cfg.PluginSecret != "file_secret" {
		t.Errorf("PluginSecret = %q", cfg.PluginSecret)
	}
	if cfg.SessionPath != "/tmp/file-session.json" {
		t.Errorf("SessionPath = %q", cfg.SessionPath)
	}
}

func TestLoad_EnvOverridesConfigFile(t *testing.T) {
	home := setupTempHome(t)
	writeConfigFile(t, home, `{
  "base_url": "https://file.example.com",
  "plugin_id": "file_id",
  "plugin_secret": "file_secret",
  "session_path": "/tmp/file-session.json"
}`)
	t.Setenv("LARK_BASE_URL", "https://env.example.com")
	t.Setenv("LARK_PLUGIN_ID", "env_id")
	t.Setenv("LARK_PLUGIN_SECRET", "env_secret")
	t.Setenv("LARK_SESSION_PATH", "/tmp/env-session.json")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BaseURL != "https://env.example.com" {
		t.Errorf("BaseURL = %q", cfg.BaseURL)
	}
	if cfg.PluginID != "env_id" {
		t.Errorf("PluginID = %q", cfg.PluginID)
	}
	if cfg.PluginSecret != "env_secret" {
		t.Errorf("PluginSecret = %q", cfg.PluginSecret)
	}
	if cfg.SessionPath != "/tmp/env-session.json" {
		t.Errorf("SessionPath = %q", cfg.SessionPath)
	}
}

func TestLoad_MissingConfigFileDoesNotError(t *testing.T) {
	setupTempHome(t)

	if _, err := config.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoad_InvalidConfigFileJSONReturnsError(t *testing.T) {
	home := setupTempHome(t)
	writeConfigFile(t, home, `{"plugin_id":`)

	if _, err := config.Load(); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_ValidateForPluginTokenPassesWithConfigFileCredentials(t *testing.T) {
	home := setupTempHome(t)
	writeConfigFile(t, home, `{
  "base_url": "https://file.example.com",
  "plugin_id": "file_id",
  "plugin_secret": "file_secret"
}`)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := cfg.ValidateForPluginToken(); err != nil {
		t.Fatalf("expected validate success, got: %v", err)
	}
}

func TestValidateForPluginToken_MissingRequired(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
	}{
		{"missing BaseURL", config.Config{PluginID: "x", PluginSecret: "y"}},
		{"missing PluginID", config.Config{BaseURL: "x", PluginSecret: "y"}},
		{"missing PluginSecret", config.Config{BaseURL: "x", PluginID: "y"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.ValidateForPluginToken(); err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestValidateForPluginToken_AllPresent(t *testing.T) {
	c := config.Config{BaseURL: "https://x.com", PluginID: "id", PluginSecret: "secret"}
	if err := c.ValidateForPluginToken(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
