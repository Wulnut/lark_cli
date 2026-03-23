package project

import (
	"fmt"
	"net/http"
	"os"

	"lark_cli/internal/cmddeps"
	"lark_cli/internal/openapi"
	"lark_cli/internal/tui"

	"github.com/spf13/cobra"
)

func NewWorkItemCmd(deps cmddeps.Deps) *cobra.Command {
	var projectKey string
	var tuiMode bool

	cmd := &cobra.Command{
		Use:          "work-item",
		Short:        "Manage work items in a Feishu Project space",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !tuiMode {
				return cmd.Help()
			}

			if deps.PluginTokenProvider == nil {
				return fmt.Errorf("not logged in: run 'lark login' first")
			}
			if projectKey == "" {
				return fmt.Errorf("--project-key is required")
			}

			client := openapi.NewClient(deps.Config.BaseURL, &http.Client{Timeout: deps.Config.HTTPTimeout}, deps.PluginTokenProvider)
			return tui.RunWorkItemTypeList(os.Stdout, deps.Config.UserKey, client, projectKey)
		},
	}

	cmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Launch interactive TUI")
	cmd.Flags().StringVarP(&projectKey, "project-key", "k", "", "Project key (required for --tui)")

	cmd.AddCommand(NewWorkItemSearchCmd(deps))
	return cmd
}
