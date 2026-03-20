/*
 * @Author: liangyz liangyz@seirobitcs.net
 * @Date: 2026-03-19 18:26:09
 * @LastEditors: liangyz liangyz@seirobitcs.net
 * @LastEditTime: 2026-03-20 11:25:06
 * @FilePath: \lark_cli\cmd\root.go
 * @Description: root.go
 */
package cmd

import (
	"fmt"
	"lark_cli/internal/auth"
	"lark_cli/internal/config"
	"lark_cli/internal/session"
	"lark_cli/internal/tui"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns the root command with subcommands registered.
func NewRootCmd(deps Deps) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lark",
		Short: "Lark CLI - Feishu Project command line tool",
		Long: `Lark CLI is a tool for interacting with the Feishu Project OpenAPI.
It provides a command line interface to manage and interact with various Feishu Project resources.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return cmd.Help()
			}
			if err := tui.Run(deps.Stdout); err != nil {
				return fmt.Errorf("interactive UI: %w", err)
			}
			return nil
		},
	}

	rootCmd.AddCommand(NewLoginCmd(deps))
	rootCmd.AddCommand(NewLogoutCmd(deps))
	rootCmd.AddCommand(NewAuthCmd(deps))

	return rootCmd
}

// Execute loads config, builds dependencies, and runs the root command.
func Execute() {
	cfg, err := config.Load()
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
