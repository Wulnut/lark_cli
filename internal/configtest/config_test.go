package configtest

import (
	"os"
	"testing"

	"lark_cli/internal/config"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	// Clear relevant env vars
	os.Unsetenv("LARK_BASE_URL")
	os.Unsetenv("LARK_PLUGIN_ID")
	os.Unsetenv("LARK_PLUGIN_SECRET")
	os.Unsetenv("LARK_SESSION_PATH")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != config.DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, config.DefaultBaseURL)
	}
	if cfg.SessionPath == "" {
		t.Error("SessionPath should have a default")
	}
}

func TestLoadFromEnv_CustomValues(t *testing.T) {
	t.Setenv("LARK_BASE_URL", "https://custom.example.com")
	t.Setenv("LARK_PLUGIN_ID", "test_id")
	t.Setenv("LARK_PLUGIN_SECRET", "test_secret")
	t.Setenv("LARK_SESSION_PATH", "/tmp/test-session.json")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != "https://custom.example.com" {
		t.Errorf("BaseURL = %q, want custom", cfg.BaseURL)
	}
	if cfg.PluginID != "test_id" {
		t.Errorf("PluginID = %q", cfg.PluginID)
	}
	if cfg.PluginSecret != "test_secret" {
		t.Errorf("PluginSecret = %q", cfg.PluginSecret)
	}
	if cfg.SessionPath != "/tmp/test-session.json" {
		t.Errorf("SessionPath = %q", cfg.SessionPath)
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
