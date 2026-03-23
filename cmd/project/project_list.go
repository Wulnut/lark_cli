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

func NewProjectListCmd(deps cmddeps.Deps) *cobra.Command {
	var order string
	var tuiMode bool

	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List Feishu Project spaces",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if deps.PluginTokenProvider == nil {
				return fmt.Errorf("not logged in: run 'lark login' first")
			}

			client := openapi.NewClient(deps.Config.BaseURL, &http.Client{Timeout: deps.Config.HTTPTimeout}, deps.PluginTokenProvider)

			if tuiMode {
				return tui.RunProjectList(os.Stdout, deps.Config.UserKey, client)
			}

			return runList(ctx, deps.Stdout, client, deps.Config.UserKey, order)
		},
	}

	cmd.Flags().StringVar(&order, "order", "", "Sort order: last_visited, +last_visited, -last_visited")
	cmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Launch interactive TUI")

	return cmd
}

func runList(ctx context.Context, out io.Writer, client *openapi.Client, userKey, order string) error {
	projects, err := client.ListProjectsWithDetails(ctx, userKey, order)
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Fprintln(out, "No project spaces found.")
		return nil
	}

	for i, p := range projects {
		num := fmt.Sprintf("%d.", i+1)
		name := p.Name
		if name == "" {
			name = "(unnamed)"
		}
		fmt.Fprintf(out, "%s %s (%s)\n", num, name, p.SimpleName)
	}
	return nil
}
