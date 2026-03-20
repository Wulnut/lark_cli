package tui

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"
)

type rootModel struct{}

func (rootModel) Init() tea.Cmd { return nil }

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (rootModel) View() string {
	return "Lark CLI — interactive mode (q or Esc to quit)\n"
}

// Run starts the Bubble Tea program. If out is non-nil it is used as program output (e.g. deps.Stdout).
func Run(out io.Writer) error {
	var opts []tea.ProgramOption
	if out != nil {
		opts = append(opts, tea.WithOutput(out))
	}
	p := tea.NewProgram(rootModel{}, opts...)
	_, err := p.Run()
	return err
}
