package tui

import (
	"io"

	"lark_cli/internal/openapi"
	tea "github.com/charmbracelet/bubbletea"
)

// RunProjectList starts a standalone Bubble Tea program for project list.
func RunProjectList(out io.Writer, userKey string, client *openapi.Client) error {
	var opts []tea.ProgramOption
	if out != nil {
		opts = append(opts, tea.WithOutput(out))
	}
	model := rootModel{
		userKey:   userKey,
		client:    client,
		state:     stateProjectList,
		entryMode: entryProjectList,
		width:     80,
	}
	p := tea.NewProgram(model, opts...)
	_, err := p.Run()
	return err
}
