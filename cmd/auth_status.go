package cmd

import (
	"fmt"
	"time"

	"lark_cli/internal/auth"

	"github.com/spf13/cobra"
)

func NewAuthStatusCmd(deps Deps) *cobra.Command {
	var refreshToken bool

	cmd := &cobra.Command{
		Use:          "status",
		Short:        "Show current authentication status",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			exists, err := deps.Store.Exists(ctx)
			if err != nil {
				return fmt.Errorf("failed to check session: %w", err)
			}
			if !exists {
				fmt.Fprintln(deps.Stderr, "Not logged in.")
				return fmt.Errorf("not logged in")
			}

			sess, err := deps.Store.Load(ctx)
			if err != nil {
				return fmt.Errorf("failed to load session: %w", err)
			}

			fmt.Fprintf(deps.Stdout, "Session path: %s\n", deps.Config.SessionPath)
			fmt.Fprintf(deps.Stdout, "User key:     %s\n", auth.Mask(sess.UserKey))
			fmt.Fprintf(deps.Stdout, "Login type:   %s\n", sess.LoginType)

			if refreshToken && deps.PluginTokenProvider != nil {
				token, err := deps.PluginTokenProvider.ForceRefresh(ctx)
				if err != nil {
					fmt.Fprintf(deps.Stderr, "Failed to refresh plugin token: %v\n", err)
				} else {
					fmt.Fprintf(deps.Stdout, "Plugin token: %s (refreshed)\n", auth.Mask(token))
					// Reload session to show updated expiry
					sess, _ = deps.Store.Load(ctx)
				}
			}

			if sess.PluginAccessToken != "" {
				fmt.Fprintf(deps.Stdout, "Plugin token: %s\n", auth.Mask(sess.PluginAccessToken))
				fmt.Fprintf(deps.Stdout, "Expires at:   %s\n", sess.PluginAccessTokenExpiresAt.Format(time.RFC3339))
				if sess.PluginAccessTokenExpiresAt.Before(time.Now()) {
					fmt.Fprintln(deps.Stdout, "Status:       expired")
				} else {
					fmt.Fprintln(deps.Stdout, "Status:       valid")
				}
			} else if !refreshToken {
				fmt.Fprintln(deps.Stdout, "Plugin token: none")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&refreshToken, "refresh-plugin-token", false, "Force refresh the plugin access token")

	return cmd
}
