package tui

import (
	"io"

	"lark_cli/internal/openapi"
	tea "github.com/charmbracelet/bubbletea"
)

// RunWorkItemTypeList starts a standalone Bubble Tea program for work item type list.
func RunWorkItemTypeList(out io.Writer, userKey string, client *openapi.Client, projectKey string) error {
	var opts []tea.ProgramOption
	if out != nil {
		opts = append(opts, tea.WithOutput(out))
	}
	model := rootModel{
		userKey:    userKey,
		client:     client,
		state:      stateWorkItemTypeList,
		entryMode:  entryWorkItemTypeStandalone,
		projectKey: projectKey,
		width:      80,
	}
	p := tea.NewProgram(model, opts...)
	_, err := p.Run()
	return err
}
