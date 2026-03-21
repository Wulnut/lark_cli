package tui

import (
	"context"
	"fmt"
	"io"

	"lark_cli/internal/openapi"
	tea "github.com/charmbracelet/bubbletea"
)

type tuiState int

const (
	stateLoading tuiState = iota
	stateSuccess
	stateDegraded // API failed, show user_key only
)

type rootModel struct {
	userKey  string
	client   *openapi.Client
	user     *openapi.UserInfo
	fetchErr error
	state    tuiState
}

type userFetchedMsg   struct{ user *openapi.UserInfo }
type userFetchFailedMsg struct{ err error }

// fetchUserCmd initiates the user info fetch.
// If client is nil, immediately returns a degraded state.
func fetchUserCmd(userKey string, client *openapi.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil || userKey == "" {
			return userFetchFailedMsg{err: fmt.Errorf("not logged in")}
		}
		user, err := client.QueryCurrentUser(context.Background(), userKey)
		if err != nil {
			return userFetchFailedMsg{err: err}
		}
		if user == nil {
			return userFetchFailedMsg{err: fmt.Errorf("user not found")}
		}
		return userFetchedMsg{user: user}
	}
}

func (m rootModel) Init() tea.Cmd {
	return fetchUserCmd(m.userKey, m.client)
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case userFetchedMsg:
		m.user = msg.user
		m.state = stateSuccess
		return m, nil
	case userFetchFailedMsg:
		m.fetchErr = msg.err
		m.state = stateDegraded
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m rootModel) View() string {
	switch m.state {
	case stateLoading:
		return "🔄 正在获取用户信息...\n"
	case stateSuccess:
		name := displayName(m.user)
		email := m.user.Email
		status := statusEmoji(m.user.Status)
		return fmt.Sprintf("👤 %s <%s>  %s 已登录\n───────────────────────────────\n按 [Q] 退出\n", name, email, status)
	case stateDegraded:
		if m.userKey != "" {
			return fmt.Sprintf("👤 %s  ⚠️ 仅显示 user_key\n───────────────────────────────\n按 [Q] 退出\n", m.userKey)
		}
		return "⚠️ 未登录，请先运行 lark login\n───────────────────────────────\n按 [Q] 退出\n"
	}
	return ""
}

func displayName(u *openapi.UserInfo) string {
	if u == nil {
		return ""
	}
	if u.Name.ZhCN != "" {
		return u.Name.ZhCN
	}
	if u.Name.Default != "" {
		return u.Name.Default
	}
	if u.NameCn != "" {
		return u.NameCn
	}
	return u.UserKey
}

func statusEmoji(status string) string {
	switch status {
	case "activated":
		return "✅"
	default:
		return "⚠️"
	}
}

// Run starts the Bubble Tea program. If out is non-nil it is used as program output (e.g. deps.Stdout).
func Run(out io.Writer, userKey string, client *openapi.Client) error {
	var opts []tea.ProgramOption
	if out != nil {
		opts = append(opts, tea.WithOutput(out))
	}
	model := rootModel{
		userKey: userKey,
		client:  client,
		state:   stateLoading,
	}
	p := tea.NewProgram(model, opts...)
	_, err := p.Run()
	return err
}
