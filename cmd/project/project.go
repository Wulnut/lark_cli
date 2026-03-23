package project

import (
	"lark_cli/internal/cmddeps"

	"github.com/spf13/cobra"
)

func NewProjectCmd(deps cmddeps.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage Feishu Project spaces",
		Long:  "List and view details of Feishu Project spaces you have access to.",
	}
	cmd.AddCommand(NewProjectListCmd(deps))
	cmd.AddCommand(NewProjectDetailCmd(deps))
	cmd.AddCommand(NewWorkItemTypeCmd(deps))
	cmd.AddCommand(NewWorkItemCmd(deps))
	return cmd
}
