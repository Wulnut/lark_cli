package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lark_cli/internal/config"

	"github.com/spf13/cobra"
)

func NewLoginCmd(deps Deps) *cobra.Command {
	var way string

	cmd := &cobra.Command{
		Use:          "login",
		Short:        "Log in to Feishu Project",
		Long:         "Save user credentials for subsequent API calls.\nCurrently supports: lark login -w user_key <user_key>",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if way != "user_key" {
				return fmt.Errorf("unsupported login method %q: only 'user_key' is supported", way)
			}

			userKey := args[0]
			if strings.TrimSpace(userKey) == "" {
				return fmt.Errorf("invalid user_key: must not be empty")
			}

			// Load existing config or create new one
			cfg := deps.Config
			cfg.UserKey = userKey

			// Save to config.json
			if err := saveConfig(&cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Fprintln(deps.Stdout, "Login successful.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&way, "way", "w", "", "Login method (required: user_key)")
	_ = cmd.MarkFlagRequired("way")

	return cmd
}

// saveConfig saves the config to ~/.lark/config.json.
func saveConfig(cfg *config.Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dir := filepath.Join(home, ".lark")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path := filepath.Join(dir, "config.json")

	// Start with existing file config to avoid overwriting previously saved fields.
	existing := fileConfig{}
	if b, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(b, &existing); err != nil {
			return fmt.Errorf("failed to parse existing config file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing config file: %w", err)
	}

	// Merge: preserve existing persisted fields and update user_key.
	fc := fileConfig{
		BaseURL:      existing.BaseURL,
		PluginID:     existing.PluginID,
		PluginSecret: existing.PluginSecret,
		UserKey:      cfg.UserKey,
		SessionPath:  existing.SessionPath,
	}

	if fc.BaseURL == "" {
		fc.BaseURL = cfg.BaseURL
	}
	if fc.PluginID == "" {
		fc.PluginID = cfg.PluginID
	}
	if fc.PluginSecret == "" {
		fc.PluginSecret = cfg.PluginSecret
	}
	if fc.SessionPath == "" {
		fc.SessionPath = cfg.SessionPath
	}

	data, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// fileConfig matches the config.json structure.
type fileConfig struct {
	BaseURL      string `json:"base_url"`
	PluginID     string `json:"plugin_id"`
	PluginSecret string `json:"plugin_secret"`
	UserKey      string `json:"user_key"`
	SessionPath  string `json:"session_path"`
}
