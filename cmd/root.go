package cmd

import (
	"net/http"
	"os"

	"lark_cli/internal/auth"
	"lark_cli/internal/config"
	"lark_cli/internal/session"

	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns the root command with subcommands registered.
func NewRootCmd(deps Deps) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lark",
		Short: "Lark CLI - Feishu Project command line tool",
		Long: `Lark CLI is a tool for interacting with the Feishu Project OpenAPI.
It provides a command line interface to manage and interact with various Feishu Project resources.`,
	}

	rootCmd.AddCommand(NewLoginCmd(deps))
	rootCmd.AddCommand(NewLogoutCmd(deps))
	rootCmd.AddCommand(NewAuthCmd(deps))

	return rootCmd
}

// Execute loads config, builds dependencies, and runs the root command.
func Execute() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		os.Stderr.WriteString("Error loading config: " + err.Error() + "\n")
		os.Exit(1)
	}

	store := session.NewFileStore(cfg.SessionPath)

	var tokenProvider auth.PluginTokenProvider
	var headerProvider auth.HeaderProvider

	// Only build token/header providers if plugin credentials are available.
	if cfg.ValidateForPluginToken() == nil {
		httpClient := &http.Client{Timeout: cfg.HTTPTimeout}
		client := auth.NewPluginTokenClient(httpClient, cfg.BaseURL, cfg.PluginID, cfg.PluginSecret)
		tokenProvider = auth.NewPluginTokenProvider(store, client, cfg.RefreshLeeway)
		headerProvider = auth.NewHeaderProvider(store, tokenProvider)
	}

	deps := Deps{
		Config:              cfg,
		Store:               store,
		PluginTokenProvider: tokenProvider,
		HeaderProvider:      headerProvider,
		Stdout:              os.Stdout,
		Stderr:              os.Stderr,
	}

	rootCmd := NewRootCmd(deps)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
