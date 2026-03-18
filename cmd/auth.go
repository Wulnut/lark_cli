package cmd

import (
	"github.com/spf13/cobra"
)

func NewAuthCmd(deps Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication management",
	}

	cmd.AddCommand(NewAuthStatusCmd(deps))

	return cmd
}
