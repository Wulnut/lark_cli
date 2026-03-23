package tui

import (
	"fmt"
	"strings"
	"testing"

	"lark_cli/internal/openapi"
	tea "github.com/charmbracelet/bubbletea"
)

func TestRootModel_EnterWithOutOfRangeCursor_DoesNotPanic(t *testing.T) {
	m := rootModel{
		state:    stateProjectList,
		cursor:   3,
		projects: []openapi.ProjectDetail{{ProjectKey: "p1", Name: "n1"}},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Update panicked on out-of-range cursor: %v", r)
		}
	}()

	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
}

func TestRootModel_QInStandaloneWorkItemType_Quits(t *testing.T) {
	m := rootModel{state: stateWorkItemTypeList, entryMode: entryWorkItemTypeStandalone}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatalf("expected quit command, got nil")
	}

	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestRootModel_WorkItemTypeView_HidesFarRowsWhenCursorMoves(t *testing.T) {
	types := make([]openapi.WorkItemType, 0, 50)
	for i := 1; i <= 50; i++ {
		types = append(types, openapi.WorkItemType{TypeKey: fmt.Sprintf("type-%02d", i), Name: fmt.Sprintf("Type %02d", i), APIName: fmt.Sprintf("api_%02d", i)})
	}

	m := rootModel{
		state:         stateWorkItemTypeList,
		workItemTypes: types,
		cursor:        20,
		projectKey:    "p1",
		width:         100,
	}

	view := m.View()
	if strings.Contains(view, "type-01") {
		t.Fatalf("expected viewport to hide early rows when cursor moved, but type-01 is visible")
	}
}

func TestRootModel_WorkItemTypeFilter_BySlashInput(t *testing.T) {
	m := rootModel{
		state: stateWorkItemTypeList,
		workItemTypes: []openapi.WorkItemType{
			{TypeKey: "task", Name: "Task", APIName: "task_api"},
			{TypeKey: "bug", Name: "Bug", APIName: "bug_api"},
		},
		width: 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2, _ := m1.(rootModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	view := m2.(rootModel).View()

	if strings.Contains(strings.ToLower(view), "task") {
		t.Fatalf("expected filter query 'b' to hide task row, got view: %s", view)
	}
	if !strings.Contains(strings.ToLower(view), "bug") {
		t.Fatalf("expected filter query 'b' to keep bug row, got view: %s", view)
	}
}

func TestRenderProjectDetail_Admins_ShowLoadingBeforeNamesLoaded(t *testing.T) {
	m := rootModel{
		state: stateProjectDetail,
		selected: &openapi.ProjectDetail{
			ProjectKey:     "p1",
			Name:           "Project A",
			Administrators: []string{"ou_admin1"},
		},
		projects: []openapi.ProjectDetail{{ProjectKey: "p1", Name: "Project A", Administrators: []string{"ou_admin1"}}},
		cursor:   0,
		width:    100,
	}

	view := m.View()
	if strings.Contains(view, "ou_admin1") {
		t.Fatalf("expected loading placeholder before admin names are loaded, but raw user key is visible: %s", view)
	}
	if !strings.Contains(strings.ToLower(view), "loading") {
		t.Fatalf("expected loading placeholder for unresolved admin, got view: %s", view)
	}
}

func TestRootModel_WorkItemTypeFilterMode_ShowsExitHint(t *testing.T) {
	m := rootModel{
		state: stateWorkItemTypeList,
		workItemTypes: []openapi.WorkItemType{
			{TypeKey: "task", Name: "Task", APIName: "task_api"},
		},
		width: 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	view := m1.(rootModel).View()
	if !strings.Contains(strings.ToLower(view), "enter/esc") {
		t.Fatalf("expected filter mode footer to include exit hint, got view: %s", view)
	}
}

func TestRootModel_OpenSearchBuilder_FromProjectDetail(t *testing.T) {
	m := rootModel{
		state: stateProjectDetail,
		selected: &openapi.ProjectDetail{
			ProjectKey: "p1",
			Name:       "Project A",
		},
		width: 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	after := m1.(rootModel)
	if after.state != stateWorkItemSearchBuilder {
		t.Fatalf("expected stateWorkItemSearchBuilder, got %v", after.state)
	}
	if after.searchBuilderProjectKey != "p1" {
		t.Fatalf("expected project key p1 in search builder, got %q", after.searchBuilderProjectKey)
	}
}

func TestRootModel_SearchBuilderView_ShowsMinimalLayoutAndHints(t *testing.T) {
	m := rootModel{
		state:                   stateWorkItemSearchBuilder,
		searchBuilderProjectKey: "p1",
		searchBuilderTypeKey:    "issue",
		searchBuilderMe:         true,
		searchBuilderStatuses:   "doing,todo",
		searchBuilderCreatedFrom: "2026-01-01",
		searchBuilderCreatedTo:   "2026-01-31",
		searchBuilderRawJSON:    "{}",
		width:                   100,
	}

	view := strings.ToLower(m.View())
	if !strings.Contains(view, "search builder") {
		t.Fatalf("expected search builder title, got: %s", view)
	}
	if !strings.Contains(view, "minimal") {
		t.Fatalf("expected minimal aesthetic tagline, got: %s", view)
	}
	if !strings.Contains(view, "[tab]") || !strings.Contains(view, "[enter]") {
		t.Fatalf("expected key hints for tab and enter, got: %s", view)
	}
}

func TestRootModel_SearchBuilder_TabCyclesFocus(t *testing.T) {
	m := rootModel{state: stateWorkItemSearchBuilder, width: 100}
	start := m.searchBuilderFocus
	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'\t'}})
	after := m1.(rootModel)
	if after.searchBuilderFocus == start {
		t.Fatalf("expected focus to move on tab")
	}
}

func TestRootModel_SearchBuilder_EnterShowsResultView(t *testing.T) {
	m := rootModel{
		state:                   stateWorkItemSearchBuilder,
		searchBuilderProjectKey: "p1",
		searchBuilderTypeKey:    "issue",
		searchBuilderResults: []map[string]any{
			{"id": 1, "name": "Issue A"},
		},
		width: 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	after := m1.(rootModel)
	if after.state != stateWorkItemSearchResults {
		t.Fatalf("expected stateWorkItemSearchResults, got %v", after.state)
	}
}

func TestRootModel_SearchResultView_ShowsRowsAndBackHint(t *testing.T) {
	m := rootModel{
		state: stateWorkItemSearchResults,
		searchBuilderResults: []map[string]any{
			{"id": 101, "name": "Issue 101"},
			{"id": 102, "name": "Issue 102"},
		},
		width: 100,
	}

	view := strings.ToLower(m.View())
	if !strings.Contains(view, "issue 101") || !strings.Contains(view, "issue 102") {
		t.Fatalf("expected result rows in view, got: %s", view)
	}
	if !strings.Contains(view, "[q]") || !strings.Contains(view, "builder") {
		t.Fatalf("expected back hint to builder, got: %s", view)
	}
}

func TestRootModel_SearchBuilder_EditStatusesAtFocusedField(t *testing.T) {
	m := rootModel{
		state:              stateWorkItemSearchBuilder,
		searchBuilderFocus: 3,
		width:              100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m2, _ := m1.(rootModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	after := m2.(rootModel)
	if after.searchBuilderStatuses != "do" {
		t.Fatalf("expected statuses to be edited, got %q", after.searchBuilderStatuses)
	}
}

func TestRootModel_SearchBuilder_BackspaceEditsFocusedField(t *testing.T) {
	m := rootModel{
		state:                 stateWorkItemSearchBuilder,
		searchBuilderFocus:    4,
		searchBuilderCreatedFrom: "2026-01-01",
		width:                 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	after := m1.(rootModel)
	if after.searchBuilderCreatedFrom != "2026-01-0" {
		t.Fatalf("expected created-from backspace edit, got %q", after.searchBuilderCreatedFrom)
	}
}

func TestRootModel_SearchBuilderView_HighlightsFocusedField(t *testing.T) {
	m := rootModel{
		state:                 stateWorkItemSearchBuilder,
		searchBuilderProjectKey: "p1",
		searchBuilderTypeKey:  "issue",
		searchBuilderFocus:    3,
		searchBuilderStatuses: "doing",
		width:                 100,
	}

	view := strings.ToLower(m.View())
	if !strings.Contains(view, "› statuses") {
		t.Fatalf("expected focused statuses field marker in view, got: %s", view)
	}
}

func TestRootModel_SearchResults_JKNavigatesSelection(t *testing.T) {
	m := rootModel{
		state: stateWorkItemSearchResults,
		searchBuilderResults: []map[string]any{
			{"id": 101, "name": "Issue 101"},
			{"id": 102, "name": "Issue 102"},
			{"id": 103, "name": "Issue 103"},
		},
		width: 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	afterJ := m1.(rootModel)
	if afterJ.cursor != 1 {
		t.Fatalf("expected cursor 1 after j, got %d", afterJ.cursor)
	}

	m2, _ := afterJ.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	afterK := m2.(rootModel)
	if afterK.cursor != 0 {
		t.Fatalf("expected cursor 0 after k, got %d", afterK.cursor)
	}
}

func TestRootModel_SearchResultView_HighlightsSelectedRow(t *testing.T) {
	m := rootModel{
		state: stateWorkItemSearchResults,
		cursor: 1,
		searchBuilderResults: []map[string]any{
			{"id": 101, "name": "Issue 101"},
			{"id": 102, "name": "Issue 102"},
		},
		width: 100,
	}

	view := strings.ToLower(m.View())
	if !strings.Contains(view, "› [102] issue 102") {
		t.Fatalf("expected selected row marker in results, got: %s", view)
	}
}

func TestRootModel_FilterMode_JKNavigateWithoutChangingQuery(t *testing.T) {
	m := rootModel{
		state: stateWorkItemTypeList,
		workItemTypes: []openapi.WorkItemType{
			{TypeKey: "task", Name: "Task", APIName: "task_api"},
			{TypeKey: "bug", Name: "Bug", APIName: "bug_api"},
		},
		width: 100,
	}

	m1, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2, _ := m1.(rootModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	before := m2.(rootModel)
	if before.cursor != 0 || before.workItemTypeFilterQuery != "b" {
		t.Fatalf("unexpected precondition: cursor=%d query=%q", before.cursor, before.workItemTypeFilterQuery)
	}

	m3, _ := before.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	afterJ := m3.(rootModel)
	if afterJ.workItemTypeFilterQuery != "b" {
		t.Fatalf("expected query to stay 'b' when pressing j in filter mode, got %q", afterJ.workItemTypeFilterQuery)
	}
	if afterJ.cursor != 0 {
		t.Fatalf("expected cursor to stay at 0 for single filtered result, got %d", afterJ.cursor)
	}

	m4, _ := afterJ.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	afterK := m4.(rootModel)
	if afterK.workItemTypeFilterQuery != "b" {
		t.Fatalf("expected query to stay 'b' when pressing k in filter mode, got %q", afterK.workItemTypeFilterQuery)
	}
}
