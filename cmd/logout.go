package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewLogoutCmd(deps Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out and clear local session",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := deps.Store.Delete(ctx); err != nil {
				return fmt.Errorf("failed to delete session: %w", err)
			}

			fmt.Fprintln(deps.Stdout, "Logged out.")
			return nil
		},
	}
}
