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

func NewProjectDetailCmd(deps cmddeps.Deps) *cobra.Command {
	var projectKeys []string
	var simpleNames []string
	var tuiMode bool

	cmd := &cobra.Command{
		Use:          "detail",
		Short:        "Get Feishu Project space details",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if deps.PluginTokenProvider == nil {
				return fmt.Errorf("not logged in: run 'lark login' first")
			}

			client := openapi.NewClient(deps.Config.BaseURL, &http.Client{Timeout: deps.Config.HTTPTimeout}, deps.PluginTokenProvider)

			if tuiMode {
				return tui.RunProjectDetail(os.Stdout, deps.Config.UserKey, client)
			}

			if len(projectKeys) == 0 && len(simpleNames) == 0 {
				return tui.RunProjectDetail(os.Stdout, deps.Config.UserKey, client)
			}

			return runDetail(ctx, deps.Stdout, client, deps.Config.UserKey, projectKeys, simpleNames)
		},
	}

	cmd.Flags().StringSliceVarP(&projectKeys, "key", "k", nil, "Project key (can be specified multiple times, max 100)")
	cmd.Flags().StringSliceVarP(&simpleNames, "simple-name", "s", nil, "Simple name (can be specified multiple times)")
	cmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Launch interactive TUI")

	return cmd
}

func runDetail(ctx context.Context, out io.Writer, client *openapi.Client, userKey string, projectKeys, simpleNames []string) error {
	details, err := client.GetProjectDetails(ctx, userKey, projectKeys, simpleNames)
	if err != nil {
		return fmt.Errorf("failed to get project details: %w", err)
	}

	if len(details) == 0 {
		fmt.Fprintln(out, "No project details found.")
		return nil
	}

	for _, d := range details {
		name := d.Name
		if name == "" {
			name = "(unnamed)"
		}
		fmt.Fprintf(out, "Name:         %s\n", name)
		fmt.Fprintf(out, "Project Key:   %s\n", d.ProjectKey)
		fmt.Fprintf(out, "Simple Name:   %s\n", d.SimpleName)
		if len(d.Administrators) > 0 {
			fmt.Fprintln(out, "Administrators:")
			for _, admin := range d.Administrators {
				fmt.Fprintf(out, "  - %s\n", admin)
			}
		}
		fmt.Fprintln(out)
	}
	return nil
}
