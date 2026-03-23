package cmdtest

import (
	"testing"

	"lark_cli/cmd/project"
	"lark_cli/internal/cmddeps"
)

func TestProjectWorkItem_CommandHierarchy(t *testing.T) {
	root := project.NewProjectCmd(cmddeps.Deps{})

	workItemCmd, _, err := root.Find([]string{"work-item"})
	if err != nil {
		t.Fatalf("find work-item command: %v", err)
	}
	if workItemCmd == nil || workItemCmd.Use != "work-item" {
		t.Fatalf("expected parent command 'work-item', got %v", workItemCmd)
	}

	searchCmd, _, err := root.Find([]string{"work-item", "search"})
	if err != nil {
		t.Fatalf("find work-item search command: %v", err)
	}
	if searchCmd == nil || searchCmd.Use != "search" {
		t.Fatalf("expected child command 'search', got %v", searchCmd)
	}
}
