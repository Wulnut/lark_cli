package project

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"lark_cli/internal/cmddeps"
	"lark_cli/internal/openapi"
	"lark_cli/internal/tui"

	"github.com/spf13/cobra"
)

func NewWorkItemTypeCmd(deps cmddeps.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "work-item-type",
		Short: "Manage work item types in a Feishu Project space",
	}
	cmd.AddCommand(NewWorkItemTypeListCmd(deps))
	return cmd
}

func NewWorkItemTypeListCmd(deps cmddeps.Deps) *cobra.Command {
	var projectKey string
	var tuiMode bool

	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List work item types in a Feishu Project space",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if deps.PluginTokenProvider == nil {
				return fmt.Errorf("not logged in: run 'lark login' first")
			}

			if projectKey == "" {
				return fmt.Errorf("--project-key is required")
			}

			client := openapi.NewClient(deps.Config.BaseURL, &http.Client{Timeout: deps.Config.HTTPTimeout}, deps.PluginTokenProvider)

			if tuiMode {
				return tui.RunWorkItemTypeList(os.Stdout, deps.Config.UserKey, client, projectKey)
			}

			return runWorkItemTypeList(ctx, deps.Stdout, client, projectKey)
		},
	}

	cmd.Flags().StringVarP(&projectKey, "project-key", "k", "", "Project key (required)")
	cmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Launch interactive TUI")

	return cmd
}

func runWorkItemTypeList(ctx context.Context, out io.Writer, client *openapi.Client, projectKey string) error {
	types, err := client.ListWorkItemTypes(ctx, projectKey)
	if err != nil {
		return fmt.Errorf("failed to list work item types: %w", err)
	}

	if len(types) == 0 {
		fmt.Fprintln(out, "No work item types found.")
		return nil
	}

	fmt.Fprintf(out, "Space: %s\n", projectKey)
	fmt.Fprintf(out, "%-12s %-12s %-12s %s\n", "Type Key", "Name", "API Name", "Enabled")
	fmt.Fprintf(out, "%-12s %-12s %-12s %s\n", "---------", "------------", "---------", "-------")

	for _, t := range types {
		enabled := "No"
		if t.IsDisable == 0 {
			enabled = "Yes"
		}
		fmt.Fprintf(out, "%-12s %-12s %-12s %s\n", t.TypeKey, t.Name, t.APIName, enabled)
	}
	return nil
}
