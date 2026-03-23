package cmdtest

import (
	"context"
	"strings"
	"testing"

	"lark_cli/cmd/project"
	"lark_cli/internal/auth"
	"lark_cli/internal/cmddeps"
)

type stubPluginTokenProvider struct{}

func (stubPluginTokenProvider) GetAuthContext(context.Context) (*auth.AuthContext, error) {
	return &auth.AuthContext{}, nil
}

func (stubPluginTokenProvider) ForceRefresh(context.Context) (*auth.AuthContext, error) {
	return &auth.AuthContext{}, nil
}

func TestProjectWorkItem_HasTUIFlags(t *testing.T) {
	root := project.NewProjectCmd(cmddeps.Deps{})

	workItemCmd, _, err := root.Find([]string{"work-item"})
	if err != nil {
		t.Fatalf("find work-item command: %v", err)
	}

	if workItemCmd.Flags().Lookup("tui") == nil {
		t.Fatalf("expected --tui flag on work-item command")
	}
	if workItemCmd.Flags().ShorthandLookup("t") == nil {
		t.Fatalf("expected -t shorthand flag on work-item command")
	}
}

func TestProjectWorkItem_TUIRequiresProjectKey(t *testing.T) {
	root := project.NewProjectCmd(cmddeps.Deps{PluginTokenProvider: stubPluginTokenProvider{}})
	root.SetArgs([]string{"work-item", "--tui"})
	root.SilenceUsage = true

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when --tui is used without --project-key")
	}
	if !strings.Contains(err.Error(), "--project-key is required") {
		t.Fatalf("expected missing project key error, got: %v", err)
	}
}
