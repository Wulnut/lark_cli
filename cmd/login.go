package cmd

import (
	"fmt"
	"strings"
	"time"
	"lark_cli/internal/session"
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

			ctx := cmd.Context()
			now := time.Now()

			sess := &session.Session{
				Version:   session.CurrentVersion,
				LoginType: "user_key",
				UserKey:   userKey,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := deps.Store.Save(ctx, sess); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			fmt.Fprintln(deps.Stdout, "Login successful.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&way, "way", "w", "", "Login method (required: user_key)")
	_ = cmd.MarkFlagRequired("way")

	return cmd
}
