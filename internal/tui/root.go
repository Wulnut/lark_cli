package tui

import (
	"context"
	"fmt"
	"io"
	"strings"

	"lark_cli/internal/openapi"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	tea "github.com/charmbracelet/bubbletea"
)

// tuiState defines which view is currently shown.
type tuiState int

type entryMode int

const (
	stateProjectList tuiState = iota
	stateProjectDetail
	stateWorkItemTypeList
	stateWorkItemSearchBuilder
	stateWorkItemSearchResults
)

const (
	entryMain entryMode = iota
	entryProjectList
	entryProjectDetail
	entryWorkItemTypeStandalone
)

// rootModel is the single tea model that handles all views.
type rootModel struct {
	userKey   string
	client    *openapi.Client
	user      *openapi.UserInfo
	userErr   error
	state     tuiState
	entryMode entryMode
	cursor    int
	projects  []openapi.ProjectDetail
	listErr   error
	selected  *openapi.ProjectDetail
	width     int

	// Admin name resolution
	adminNames        map[string]string
	adminEmails       map[string]string
	adminNamesForProj string
	userCache         map[string]*openapi.UserInfo
	adminLoadingForProj string

	// Work item type state
	projectKey               string
	workItemTypes            []openapi.WorkItemType
	workItemTypesErr         error
	workItemTypeFilterMode   bool
	workItemTypeFilterQuery  string
	workItemTypeFilteredIdx  []int

	// Work item search builder state
	searchBuilderProjectKey  string
	searchBuilderTypeKey     string
	searchBuilderMe          bool
	searchBuilderStatuses    string
	searchBuilderCreatedFrom string
	searchBuilderCreatedTo   string
	searchBuilderRawJSON     string
	searchBuilderFocus       int
	searchBuilderResults     []map[string]any
}

// Async message types.
type userFetchedMsg           struct{ user *openapi.UserInfo }
type userFetchFailedMsg       struct{ err error }
type projectListLoadedMsg      struct{ projects []openapi.ProjectDetail }
type projectListFailedMsg     struct{ err error }
type adminNamesLoadedMsg struct {
	names      map[string]string
	emails     map[string]string
	users      map[string]*openapi.UserInfo
	projectKey string
}
type workItemTypesLoadedMsg struct{ types []openapi.WorkItemType }
type workItemTypesFailedMsg struct{ err error }

// Init fetches data based on entry mode.
func (m rootModel) Init() tea.Cmd {
	switch m.entryMode {
	case entryWorkItemTypeStandalone:
		return tea.Batch(
			fetchUserCmd(m.userKey, m.client),
			fetchWorkItemTypesCmd(m.userKey, m.client, m.projectKey),
		)
	default:
		return tea.Batch(
			fetchUserCmd(m.userKey, m.client),
			fetchProjectListCmd(m.userKey, m.client),
		)
	}
}

// Update handles all messages.
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case userFetchedMsg:
		m.user = msg.user
		return m, nil

	case userFetchFailedMsg:
		m.userErr = msg.err
		return m, nil

	case projectListLoadedMsg:
		m.projects = msg.projects
		m.listErr = nil
		if len(m.projects) == 0 {
			m.cursor = 0
			m.selected = nil
			m.state = stateProjectList
			return m, nil
		}
		if m.cursor < 0 || m.cursor >= len(m.projects) {
			m.cursor = 0
		}
		if m.entryMode == entryProjectDetail {
			m.selected = &m.projects[m.cursor]
			m.adminLoadingForProj = m.selected.ProjectKey
			m.state = stateProjectDetail
			return m, m.loadAdminNames()
		}
		m.state = stateProjectList
		return m, nil

	case projectListFailedMsg:
		m.listErr = msg.err
		m.state = stateProjectList
		return m, nil

	case adminNamesLoadedMsg:
		if m.selected == nil || msg.projectKey != m.selected.ProjectKey {
			return m, nil
		}
		if m.userCache == nil {
			m.userCache = map[string]*openapi.UserInfo{}
		}
		for k, u := range msg.users {
			if u != nil {
				m.userCache[k] = u
			}
		}
		m.adminNames = msg.names
		m.adminEmails = msg.emails
		m.adminNamesForProj = msg.projectKey
		m.adminLoadingForProj = ""
		m.state = stateProjectDetail
		if m.cursor < 0 || m.cursor >= len(m.projects) {
			m.cursor = 0
		}
		if len(m.projects) > 0 {
			m.selected = &m.projects[m.cursor]
		}
		return m, nil

	case workItemTypesLoadedMsg:
		m.workItemTypes = msg.types
		m.workItemTypesErr = nil
		m.state = stateWorkItemTypeList
		m.cursor = 0
		m.workItemTypeFilterMode = false
		m.workItemTypeFilterQuery = ""
		m.rebuildWorkItemTypeFilterIndex()
		m.clampWorkItemTypeCursor()
		return m, nil

	case workItemTypesFailedMsg:
		m.workItemTypesErr = msg.err
		m.state = stateWorkItemTypeList
		m.cursor = 0
		m.workItemTypeFilterMode = false
		m.workItemTypeFilterQuery = ""
		m.workItemTypeFilteredIdx = nil
		return m, nil

	case tea.KeyMsg:
		if m.state == stateWorkItemSearchBuilder {
			switch msg.String() {
			case "tab":
				m.searchBuilderFocus = (m.searchBuilderFocus + 1) % 7
				return m, nil
			case "backspace":
				m.backspaceSearchBuilderField()
				return m, nil
			case "enter":
				m.state = stateWorkItemSearchResults
				if m.cursor < 0 {
					m.cursor = 0
				}
				return m, nil
			}
			if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && msg.Runes[0] == '\t' {
				m.searchBuilderFocus = (m.searchBuilderFocus + 1) % 7
				return m, nil
			}
			if msg.Type == tea.KeyRunes {
				m.appendSearchBuilderInput(string(msg.Runes))
				return m, nil
			}
		}
		if m.state == stateWorkItemSearchResults {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "down", "j":
				if m.cursor < len(m.searchBuilderResults)-1 {
					m.cursor++
				}
				return m, nil
			}
		}

		if m.state == stateWorkItemTypeList {
			if m.workItemTypeFilterMode {
				switch msg.String() {
				case "esc", "enter":
					m.workItemTypeFilterMode = false
					m.clampWorkItemTypeCursor()
					return m, nil
				case "backspace":
					if len(m.workItemTypeFilterQuery) > 0 {
						m.workItemTypeFilterQuery = m.workItemTypeFilterQuery[:len(m.workItemTypeFilterQuery)-1]
						m.rebuildWorkItemTypeFilterIndex()
						m.clampWorkItemTypeCursor()
					}
					return m, nil
				case "up", "k":
					if m.cursor > 0 {
						m.cursor--
						m.clampWorkItemTypeCursor()
					}
					return m, nil
				case "down", "j":
					if m.cursor < m.workItemTypeListLen()-1 {
						m.cursor++
						m.clampWorkItemTypeCursor()
					}
					return m, nil
				}
				if msg.Type == tea.KeyRunes {
					m.workItemTypeFilterQuery += string(msg.Runes)
					m.rebuildWorkItemTypeFilterIndex()
					m.clampWorkItemTypeCursor()
					return m, nil
				}
			}
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.state == stateProjectDetail {
				m.state = stateProjectList
				m.selected = nil
				m.adminNames = nil
				m.adminEmails = nil
				m.adminNamesForProj = ""
				m.adminLoadingForProj = ""
				return m, nil
			}
			if m.state == stateWorkItemSearchBuilder {
				m.state = stateProjectDetail
				return m, nil
			}
			if m.state == stateWorkItemSearchResults {
				m.state = stateWorkItemSearchBuilder
				return m, nil
			}
			if m.state == stateWorkItemTypeList {
				if m.entryMode == entryWorkItemTypeStandalone {
					return m, tea.Quit
				}
				m.state = stateProjectDetail
				m.workItemTypes = nil
				m.workItemTypesErr = nil
				m.workItemTypeFilterMode = false
				m.workItemTypeFilterQuery = ""
				m.workItemTypeFilteredIdx = nil
				return m, nil
			}
			return m, tea.Quit

		case "3":
			if m.state == stateProjectDetail && m.selected != nil {
				m.projectKey = m.selected.ProjectKey
				return m, fetchWorkItemTypesCmd(m.userKey, m.client, m.projectKey)
			}
		case "s":
			if m.state == stateProjectDetail && m.selected != nil {
				m.searchBuilderProjectKey = m.selected.ProjectKey
				if strings.TrimSpace(m.searchBuilderTypeKey) == "" {
					m.searchBuilderTypeKey = "issue"
				}
				m.searchBuilderFocus = 0
				m.state = stateWorkItemSearchBuilder
				return m, nil
			}

		case "/":
			if m.state == stateWorkItemTypeList {
				m.workItemTypeFilterMode = true
				m.rebuildWorkItemTypeFilterIndex()
				m.clampWorkItemTypeCursor()
				return m, nil
			}

		case "up", "k":
			if m.state == stateProjectList && m.cursor > 0 {
				m.cursor--
			}
			if m.state == stateProjectDetail && m.cursor > 0 && m.cursor < len(m.projects) {
				m.cursor--
				m.selected = &m.projects[m.cursor]
				m.adminLoadingForProj = m.selected.ProjectKey
				return m, m.loadAdminNames()
			}
			if m.state == stateWorkItemTypeList && m.cursor > 0 {
				m.cursor--
				m.clampWorkItemTypeCursor()
			}

		case "down", "j":
			if m.state == stateProjectList && m.cursor < len(m.projects)-1 {
				m.cursor++
			}
			if m.state == stateProjectDetail && m.cursor >= 0 && m.cursor < len(m.projects)-1 {
				m.cursor++
				m.selected = &m.projects[m.cursor]
				m.adminLoadingForProj = m.selected.ProjectKey
				return m, m.loadAdminNames()
			}
			if m.state == stateWorkItemTypeList && m.cursor < m.workItemTypeListLen()-1 {
				m.cursor++
				m.clampWorkItemTypeCursor()
			}

		case "enter":
			if m.state == stateProjectList && len(m.projects) > 0 {
				if m.cursor < 0 || m.cursor >= len(m.projects) {
					m.cursor = 0
				}
				m.selected = &m.projects[m.cursor]
				m.adminLoadingForProj = m.selected.ProjectKey
				return m, m.loadAdminNames()
			}
		}
	}
	return m, nil
}

// loadAdminNames fetches admin user display names.
func (m rootModel) loadAdminNames() tea.Cmd {
	if m.selected == nil {
		return nil
	}
	projectKey := m.selected.ProjectKey
	keys := m.selected.Administrators
	names := make(map[string]string, len(keys))
	emails := make(map[string]string, len(keys))
	cachedUsers := map[string]*openapi.UserInfo{}
	missing := make([]string, 0, len(keys))

	for _, k := range keys {
		if m.userCache != nil {
			if u, ok := m.userCache[k]; ok && u != nil {
				names[k] = displayName(u)
				emails[k] = u.Email
				cachedUsers[k] = u
				continue
			}
		}
		names[k] = "loading..."
		emails[k] = ""
		missing = append(missing, k)
	}

	if len(keys) == 0 {
		return func() tea.Msg {
			return adminNamesLoadedMsg{names: names, emails: emails, users: cachedUsers, projectKey: projectKey}
		}
	}
	if len(missing) == 0 {
		return func() tea.Msg {
			return adminNamesLoadedMsg{names: names, emails: emails, users: cachedUsers, projectKey: projectKey}
		}
	}
	if m.client == nil {
		return func() tea.Msg {
			return adminNamesLoadedMsg{names: names, emails: emails, users: cachedUsers, projectKey: projectKey}
		}
	}

	return func() tea.Msg {
		users, err := m.client.QueryUsers(context.Background(), missing)
		if err == nil {
			for _, k := range missing {
				if u, ok := users[k]; ok && u != nil {
					names[k] = displayName(u)
					emails[k] = u.Email
					cachedUsers[k] = u
				}
			}
		}
		return adminNamesLoadedMsg{names: names, emails: emails, users: cachedUsers, projectKey: projectKey}
	}
}

// View renders the current state.
func (m rootModel) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	switch m.state {
	case stateProjectList:
		if m.listErr != nil {
			return renderListError(m, w)
		}
		return renderProjectListView(m, w)
	case stateProjectDetail:
		return renderProjectDetailView(m, w)
	case stateWorkItemTypeList:
		if m.workItemTypesErr != nil {
			return renderWorkItemTypeError(m, w)
		}
		return renderWorkItemTypeListView(m, w)
	case stateWorkItemSearchBuilder:
		return renderWorkItemSearchBuilderView(m, w)
	case stateWorkItemSearchResults:
		return renderWorkItemSearchResultsView(m, w)
	}
	return ""
}

func (m *rootModel) rebuildWorkItemTypeFilterIndex() {
	m.workItemTypeFilteredIdx = m.workItemTypeFilteredIdx[:0]
	query := strings.ToLower(strings.TrimSpace(m.workItemTypeFilterQuery))
	for i, t := range m.workItemTypes {
		if query == "" {
			m.workItemTypeFilteredIdx = append(m.workItemTypeFilteredIdx, i)
			continue
		}
		if strings.Contains(strings.ToLower(t.TypeKey), query) || strings.Contains(strings.ToLower(t.Name), query) || strings.Contains(strings.ToLower(t.APIName), query) {
			m.workItemTypeFilteredIdx = append(m.workItemTypeFilteredIdx, i)
		}
	}
}

func (m rootModel) workItemTypeListLen() int {
	if len(m.workItemTypeFilteredIdx) > 0 || strings.TrimSpace(m.workItemTypeFilterQuery) != "" {
		return len(m.workItemTypeFilteredIdx)
	}
	return len(m.workItemTypes)
}

func (m rootModel) workItemTypeAt(filteredIndex int) (openapi.WorkItemType, bool) {
	if len(m.workItemTypeFilteredIdx) > 0 || strings.TrimSpace(m.workItemTypeFilterQuery) != "" {
		if filteredIndex < 0 || filteredIndex >= len(m.workItemTypeFilteredIdx) {
			return openapi.WorkItemType{}, false
		}
		sourceIdx := m.workItemTypeFilteredIdx[filteredIndex]
		if sourceIdx < 0 || sourceIdx >= len(m.workItemTypes) {
			return openapi.WorkItemType{}, false
		}
		return m.workItemTypes[sourceIdx], true
	}
	if filteredIndex < 0 || filteredIndex >= len(m.workItemTypes) {
		return openapi.WorkItemType{}, false
	}
	return m.workItemTypes[filteredIndex], true
}

func (m *rootModel) clampWorkItemTypeCursor() {
	total := m.workItemTypeListLen()
	if total <= 0 {
		m.cursor = 0
		return
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= total {
		m.cursor = total - 1
	}
}

// --- Project List ---

func fw(text string, w int) string {
	vw := visualWidth(text)
	if vw >= w {
		return truncateToWidth(text, w)
	}
	return text + spaces(w - vw)
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

func visualWidth(s string) int {
	w := 0
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		w += runewidth.RuneWidth(r)
	}
	return w
}

func truncateToWidth(s string, w int) string {
	if w <= 0 {
		return ""
	}
	result := ""
	width := 0
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			result += string(r)
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			result += string(r)
			continue
		}
		if width+runewidth.RuneWidth(r) > w {
			break
		}
		result += string(r)
		width += runewidth.RuneWidth(r)
	}
	return result + "…"
}

func renderProjectListView(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)

	// Header bar: user info
	var userInfo string
	if m.user != nil {
		userInfo = fmt.Sprintf("  User: %s  |  Email: %s",
			displayName(m.user), m.user.Email)
	} else if m.userErr != nil {
		userInfo = fmt.Sprintf("  User error: %v", m.userErr)
	}
	headerText := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(" Lark CLI ")
	infoText := lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true).
		Render(userInfo)
	headerLine := headerText + " " + infoText

	// Column layout: # | Name | Simple Name
	// Calculate widths
	numCol := 5    // "  1.  "
	snCol := 16
	gap := 2
	nameCol := totalWidth - numCol - snCol - gap*3 - 2 // -2 for side padding
	if nameCol < 10 {
		nameCol = 10
		snCol = totalWidth - numCol - nameCol - gap*3 - 2
	}

	// Table header
	thStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
	hdrNum := thStyle.Render(fw(" # ", numCol))
	hdrName := thStyle.Render(fw(" Space Name ", nameCol))
	hdrSN := thStyle.Render(fw(" Simple Name ", snCol))
	hdrRow := fmt.Sprintf(" %s%s%s%s%s ", hdrNum, spaces(gap), hdrName, spaces(gap), hdrSN)

	// Separator
	sep := spaces(1) + "┼" + fw("", numCol) + "┼" + fw("", nameCol) + "┼" + fw("", snCol)

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(subtleColor).Italic(true)
	var footer string
	if m.state == stateProjectList {
		if len(m.projects) == 0 {
			footer = footerStyle.Render("  No project spaces found.")
		} else {
			footer = footerStyle.Render(" [↑/↓] select  [Enter] view detail  [Q] quit")
		}
	} else {
		footer = footerStyle.Render(" [↑/↓] other projects  [Q] back")
	}

	// Build rows
	rows := ""
	evStyle := lipgloss.NewStyle()
	for i, p := range m.projects {
		name := p.Name
		if name == "" {
			name = "(unnamed)"
		}
		sn := p.SimpleName
		if sn == "" {
			sn = "-"
		}

		numCell := fw(fmt.Sprintf(" %d. ", i+1), numCol)
		nameCell := fw(" "+name, nameCol)
		snCell := fw(" "+sn, snCol)
		rowStr := fmt.Sprintf(" %s%s%s%s%s ",
			numCell, spaces(gap), nameCell, spaces(gap), snCell)

		if i == m.cursor {
			selStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(primaryColor)
			rows += selStyle.Render(rowStr) + "\n"
		} else {
			rows += evStyle.Render(rowStr) + "\n"
		}
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		border.Render(headerLine),
		hdrRow,
		sep,
		rows,
		footer)
}

func renderListError(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)
	headerText := lipgloss.NewStyle().
		Foreground(primaryColor).Bold(true).Render(" Lark CLI ")
	body := lipgloss.NewStyle().
		Foreground(dangerColor).Render(fmt.Sprintf("  Error: %v", m.listErr))
	footer := lipgloss.NewStyle().
		Foreground(subtleColor).Italic(true).Render(" [Q] quit")
	return fmt.Sprintf("%s\n%s\n%s\n", border.Render(headerText), body, footer)
}

// --- Project Detail ---

func renderProjectDetailView(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)

	p := m.selected
	if p == nil {
		header := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(" Project Detail ")
		body := lipgloss.NewStyle().Foreground(subtleColor).Render("  No project selected.")
		footer := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).Render(" [Q] back")
		return fmt.Sprintf("%s\n%s\n%s\n", border.Render(header), body, footer)
	}

	header := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(" Project Detail ")
	footer := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).
		Render(fmt.Sprintf(" [↑/↓] other projects  [3] work item types  [Q] back  (%d/%d)",
			m.cursor+1, len(m.projects)))

	labelW := 16
	valueW := totalWidth - labelW - 3
	if valueW < 20 {
		valueW = 20
	}

	lStyle := lipgloss.NewStyle().Foreground(subtleColor).Bold(true).Width(labelW)
	vStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Width(valueW)

	name := p.Name
	if name == "" {
		name = "(unnamed)"
	}

	content := fmt.Sprintf("%s %s\n%s %s\n%s %s\n",
		lStyle.Render("Name:"),
		vStyle.Render(truncateToWidth(name, valueW-2)),
		lStyle.Render("Project Key:"),
		vStyle.Render(p.ProjectKey),
		lStyle.Render("Simple Name:"),
		vStyle.Render(p.SimpleName),
	)

	if len(p.Administrators) > 0 {
		content += fmt.Sprintf("\n%s\n", lStyle.Render("Administrators:"))
		adminValueW := totalWidth - 4
		if adminValueW < 20 {
			adminValueW = 20
		}
		for _, admin := range p.Administrators {
			display := "loading..."
			email := ""
			if m.adminNamesForProj == p.ProjectKey && m.adminNames != nil {
				if n, ok := m.adminNames[admin]; ok && strings.TrimSpace(n) != "" {
					display = n
				}
			}
			if m.adminNamesForProj == p.ProjectKey && m.adminEmails != nil {
				email = m.adminEmails[admin]
			}
			line := "- " + display
			if strings.TrimSpace(email) != "" {
				line += " <" + email + ">"
			}
			content += fmt.Sprintf("%s\n", vStyle.Render(truncateToWidth(line, adminValueW)))
		}
	}

	return fmt.Sprintf("%s\n%s\n%s\n", border.Render(header), content, footer)
}

// --- Helpers ---

// displayName returns the best available name for a user.
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

// fetchUserCmd fetches current user info.
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

// fetchProjectListCmd fetches the project list with details.
func fetchProjectListCmd(userKey string, client *openapi.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil || userKey == "" {
			return projectListFailedMsg{err: fmt.Errorf("not logged in")}
		}
		projects, err := client.ListProjectsWithDetails(context.Background(), userKey, "")
		if err != nil {
			return projectListFailedMsg{err: err}
		}
		return projectListLoadedMsg{projects: projects}
	}
}

// fetchWorkItemTypesCmd fetches work item types for a project.
func fetchWorkItemTypesCmd(userKey string, client *openapi.Client, projectKey string) tea.Cmd {
	return func() tea.Msg {
		if client == nil || userKey == "" || projectKey == "" {
			return workItemTypesFailedMsg{err: fmt.Errorf("not logged in")}
		}
		types, err := client.ListWorkItemTypes(context.Background(), projectKey)
		if err != nil {
			return workItemTypesFailedMsg{err: err}
		}
		return workItemTypesLoadedMsg{types: types}
	}
}

// --- Work Item Type List ---

func renderWorkItemTypeError(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)
	header := lipgloss.NewStyle().
		Foreground(primaryColor).Bold(true).Render(" Work Item Types ")
	body := lipgloss.NewStyle().
		Foreground(dangerColor).Render(fmt.Sprintf("  Error: %v", m.workItemTypesErr))
	footer := lipgloss.NewStyle().
		Foreground(subtleColor).Italic(true).Render(" [Q] back ")
	return fmt.Sprintf("%s\n%s\n%s\n", border.Render(header), body, footer)
}

func workItemTypeVisibleRange(total, cursor int) (int, int) {
	if total <= 0 {
		return 0, 0
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= total {
		cursor = total - 1
	}
	pageSize := 12
	start := cursor - pageSize/2
	if start < 0 {
		start = 0
	}
	end := start + pageSize
	if end > total {
		end = total
		start = end - pageSize
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

func (m *rootModel) appendSearchBuilderInput(s string) {
	switch m.searchBuilderFocus {
	case 0:
		m.searchBuilderProjectKey += s
	case 1:
		m.searchBuilderTypeKey += s
	case 2:
		if strings.TrimSpace(s) != "" {
			m.searchBuilderMe = !m.searchBuilderMe
		}
	case 3:
		m.searchBuilderStatuses += s
	case 4:
		m.searchBuilderCreatedFrom += s
	case 5:
		m.searchBuilderCreatedTo += s
	case 6:
		m.searchBuilderRawJSON += s
	}
}

func (m *rootModel) backspaceSearchBuilderField() {
	trimLast := func(v string) string {
		if len(v) == 0 {
			return v
		}
		return v[:len(v)-1]
	}
	switch m.searchBuilderFocus {
	case 0:
		m.searchBuilderProjectKey = trimLast(m.searchBuilderProjectKey)
	case 1:
		m.searchBuilderTypeKey = trimLast(m.searchBuilderTypeKey)
	case 2:
		m.searchBuilderMe = false
	case 3:
		m.searchBuilderStatuses = trimLast(m.searchBuilderStatuses)
	case 4:
		m.searchBuilderCreatedFrom = trimLast(m.searchBuilderCreatedFrom)
	case 5:
		m.searchBuilderCreatedTo = trimLast(m.searchBuilderCreatedTo)
	case 6:
		m.searchBuilderRawJSON = trimLast(m.searchBuilderRawJSON)
	}
}

func renderWorkItemSearchBuilderView(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)

	title := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(" Search Builder ")
	tagline := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).Render(" minimal, but not simple ")
	header := title + tagline

	fieldLabel := lipgloss.NewStyle().Foreground(subtleColor).Bold(true)
	fieldValue := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	focusedLabel := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	focusedValue := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("236"))
	marker := func(idx int, label string, value string) string {
		if m.searchBuilderFocus == idx {
			return fmt.Sprintf("› %s %s", focusedLabel.Render(label), focusedValue.Render(value))
		}
		return fmt.Sprintf("  %s %s", fieldLabel.Render(label), fieldValue.Render(value))
	}

	meVal := "No"
	if m.searchBuilderMe {
		meVal = "Yes"
	}
	body := ""
	body += marker(0, "Project:", m.searchBuilderProjectKey) + "\n"
	body += marker(1, "Type:", m.searchBuilderTypeKey) + "\n"
	body += marker(2, "Me:", meVal) + "\n"
	body += marker(3, "Statuses:", m.searchBuilderStatuses) + "\n"
	body += marker(4, "Created From:", m.searchBuilderCreatedFrom) + "\n"
	body += marker(5, "Created To:", m.searchBuilderCreatedTo) + "\n"
	body += marker(6, "Raw JSON:", truncateToWidth(m.searchBuilderRawJSON, totalWidth-20)) + "\n"

	focusNames := []string{"project", "type", "me", "status", "created_from", "created_to", "raw"}
	focus := "-"
	if m.searchBuilderFocus >= 0 && m.searchBuilderFocus < len(focusNames) {
		focus = focusNames[m.searchBuilderFocus]
	}
	body += fmt.Sprintf("\n%s %s\n", fieldLabel.Render("Focus:"), fieldValue.Render(focus))

	footer := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).
		Render(" [Tab] next field  [Enter] preview results  [Q] back")

	return fmt.Sprintf("%s\n%s\n%s\n", border.Render(header), body, footer)
}

func renderWorkItemSearchResultsView(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)

	header := lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(" Search Results ")
	rows := ""
	if len(m.searchBuilderResults) == 0 {
		rows = lipgloss.NewStyle().Foreground(subtleColor).Render("  No results.")
	} else {
		for i, row := range m.searchBuilderResults {
			id := row["id"]
			name := row["name"]
			if name == nil || strings.TrimSpace(fmt.Sprint(name)) == "" {
				name = "(unnamed)"
			}
			prefix := "  "
			style := lipgloss.NewStyle()
			if i == m.cursor {
				prefix = "› "
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(primaryColor)
			} else if i%2 == 1 {
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
			} else {
				style = lipgloss.NewStyle().Foreground(subtleColor)
			}
			line := fmt.Sprintf("%s[%v] %v", prefix, id, name)
			rows += style.Render(line) + "\n"
		}
	}

	footer := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).
		Render(" [↑/↓] move  [Q] back to builder")

	return fmt.Sprintf("%s\n%s\n%s\n", border.Render(header), rows, footer)
}

func renderWorkItemTypeListView(m rootModel, totalWidth int) string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(totalWidth)

	// Header
	headerText := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(" Work Item Types ")
	infoText := lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true).
		Render(fmt.Sprintf("  %s ", m.projectKey))
	headerLine := headerText + infoText

	// Column widths
	numCol := 5
	typeKeyCol := 14
	nameCol := 16
	apiCol := 14
	gap := 2
	enabledCol := totalWidth - numCol - typeKeyCol - nameCol - apiCol - gap*4 - 2
	if enabledCol < 8 {
		enabledCol = 8
	}

	// Table header
	thStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	hdrNum := thStyle.Render(fw(" # ", numCol))
	hdrTK := thStyle.Render(fw(" Type Key ", typeKeyCol))
	hdrName := thStyle.Render(fw(" Name ", nameCol))
	hdrAPI := thStyle.Render(fw(" API Name ", apiCol))
	hdrEnabled := thStyle.Render(fw(" Enabled ", enabledCol))
	hdrRow := fmt.Sprintf(" %s%s%s%s%s%s%s%s%s ", hdrNum, spaces(gap), hdrTK, spaces(gap), hdrName, spaces(gap), hdrAPI, spaces(gap), hdrEnabled)

	// Footer
	footerStyle := lipgloss.NewStyle().Foreground(subtleColor).Italic(true)
	total := m.workItemTypeListLen()
	start, end := workItemTypeVisibleRange(total, m.cursor)
	var footer string
	if total == 0 {
		if strings.TrimSpace(m.workItemTypeFilterQuery) != "" {
			footer = footerStyle.Render("  No matched work item types.")
		} else {
			footer = footerStyle.Render("  No work item types found.")
		}
	} else {
		base := fmt.Sprintf(" [↑/↓] select  [/] filter  [Q] back to project detail  (%d-%d/%d)", start+1, end, total)
		if m.workItemTypeFilterMode {
			base += fmt.Sprintf("  filter: %s  [Enter/Esc] exit filter", m.workItemTypeFilterQuery)
		} else if strings.TrimSpace(m.workItemTypeFilterQuery) != "" {
			base += fmt.Sprintf("  filter=%s  [/] edit", m.workItemTypeFilterQuery)
		}
		footer = footerStyle.Render(base)
	}

	// Build rows
	rows := ""
	evStyle := lipgloss.NewStyle()
	for visibleIdx := start; visibleIdx < end; visibleIdx++ {
		t, ok := m.workItemTypeAt(visibleIdx)
		if !ok {
			continue
		}
		enabled := "No"
		if t.IsDisable == 0 {
			enabled = "Yes"
		}
		numCell := fw(fmt.Sprintf(" %d. ", visibleIdx+1), numCol)
		tkCell := fw(" "+t.TypeKey, typeKeyCol)
		nameCell := fw(" "+t.Name, nameCol)
		apiCell := fw(" "+t.APIName, apiCol)
		enabledCell := fw(" "+enabled, enabledCol)
		rowStr := fmt.Sprintf(" %s%s%s%s%s%s%s%s%s ",
			numCell, spaces(gap), tkCell, spaces(gap), nameCell, spaces(gap), apiCell, spaces(gap), enabledCell)

		if visibleIdx == m.cursor {
			selStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(primaryColor)
			rows += selStyle.Render(rowStr) + "\n"
		} else {
			rows += evStyle.Render(rowStr) + "\n"
		}
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n",
		border.Render(headerLine),
		hdrRow,
		rows,
		footer)
}

// Run starts the Bubble Tea program.
func Run(out io.Writer, userKey string, client *openapi.Client) error {
	var opts []tea.ProgramOption
	if out != nil {
		opts = append(opts, tea.WithOutput(out))
	}
	model := rootModel{
		userKey:   userKey,
		client:    client,
		state:     stateProjectList,
		entryMode: entryMain,
		width:     80,
	}
	p := tea.NewProgram(model, opts...)
	_, err := p.Run()
	return err
}
