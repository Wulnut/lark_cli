package cmdtest

import (
	"testing"

	"lark_cli/cmd/project"
	"lark_cli/internal/cmddeps"
)

func TestProjectWorkItemType_CommandHierarchy(t *testing.T) {
	root := project.NewProjectCmd(cmddeps.Deps{})

	workItemTypeCmd, _, err := root.Find([]string{"work-item-type"})
	if err != nil {
		t.Fatalf("find work-item-type command: %v", err)
	}
	if workItemTypeCmd == nil || workItemTypeCmd.Use != "work-item-type" {
		t.Fatalf("expected parent command 'work-item-type', got %v", workItemTypeCmd)
	}

	listCmd, _, err := root.Find([]string{"work-item-type", "list"})
	if err != nil {
		t.Fatalf("find work-item-type list command: %v", err)
	}
	if listCmd == nil || listCmd.Use != "list" {
		t.Fatalf("expected child command 'list', got %v", listCmd)
	}
}
