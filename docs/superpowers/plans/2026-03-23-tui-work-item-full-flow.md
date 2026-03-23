# TUI Work Item Full Flow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete TUI work-item flow that supports `lark` and `lark project work-item -t`, including work item list, work item detail, search results, and search result detail.

**Architecture:** Keep a single Bubble Tea `rootModel`, but move page-specific rendering and update helpers into focused files so `internal/tui/root.go` only owns shared state, routing, and top-level transitions. Add two new OpenAPI clients for work item list/detail, then drive new TUI states through TDD, preserving back-navigation, pagination, loading/error states, and detail source tracking.

**Tech Stack:** Go, Cobra, Bubble Tea, Lipgloss, existing `internal/openapi` client, existing `test/` root test layout.

---

## File Structure

### Existing files to modify

- `cmd/project/work_item.go`
  - Add the TUI shortcut entry command under `work-item`.
- `cmd/root.go`
  - Keep root `lark` entry wired to the same TUI engine.
- `internal/openapi/work_item_search.go`
  - Reuse existing search response shape if needed for pagination helpers.
- `internal/tui/root.go`
  - Reduce to shared state, routing, helper types, run entrypoints.
- `internal/tui/project_detail.go`
  - Add transitions from project detail to type list and search builder if needed.
- `internal/tui/work_item_type.go`
  - Add enter-to-list behavior.
- `internal/tui/style.go`
  - Add shared styles for loading/error/detail cards and selected rows.
- `internal/tui/root_state_test.go`
  - Add state transition tests for new TUI behaviors.
- `test/openapitest/work_item_search_test.go`
  - Extend only if shared pagination parsing is extracted there.

### New files to create

- `internal/openapi/work_item_list.go`
  - Client for `POST /open_api/:project_key/work_item/filter`.
- `internal/openapi/work_item_detail.go`
  - Client for `POST /open_api/:project_key/work_item/:work_item_type_key/query`.
- `internal/tui/work_item_list.go`
  - Work item list rendering and list-specific helpers.
- `internal/tui/work_item_detail.go`
  - Work item detail rendering, summary/raw toggle, source-aware footer.
- `internal/tui/work_item_search.go`
  - Search builder/results/result-detail helpers split out of `root.go`.
- `test/openapitest/work_item_list_test.go`
  - API contract tests for list endpoint.
- `test/openapitest/work_item_detail_test.go`
  - API contract tests for detail endpoint.
- `test/cmdtest/project_work_item_tui_cmd_test.go`
  - Command-tree tests for `lark project work-item -t` entry.

### Files explicitly not in scope for this plan

- CLI UX redesign for `work-item-type list`
- Position-argument compatibility work
- Full CLI query ergonomics cleanup

---

## Task 1: Add failing OpenAPI tests for work item list

**Files:**
- Create: `test/openapitest/work_item_list_test.go`
- Modify: none
- Test: `test/openapitest/work_item_list_test.go`

- [ ] **Step 1: Write the failing test for the request path and method**

```go
func TestClient_ListWorkItems_UsesFilterEndpoint(t *testing.T) {
    // Expect POST /open_api/p_demo/work_item/filter
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./test/openapitest -run TestClient_ListWorkItems_UsesFilterEndpoint -v`
Expected: FAIL because `ListWorkItems` does not exist yet.

- [ ] **Step 3: Write the failing test for request payload**

```go
func TestClient_ListWorkItems_SendsTypeKeysAndPagination(t *testing.T) {
    // Expect work_item_type_keys, page_num, page_size in request body
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `go test ./test/openapitest -run 'TestClient_ListWorkItems_(UsesFilterEndpoint|SendsTypeKeysAndPagination)' -v`
Expected: FAIL because the client is not implemented.

- [ ] **Step 5: Write the failing test for response parsing**

```go
func TestClient_ListWorkItems_ParsesDataAndPagination(t *testing.T) {
    // Assert parsed items and pagination fields are available
}
```

- [ ] **Step 6: Run test to verify it fails**

Run: `go test ./test/openapitest -run 'TestClient_ListWorkItems_' -v`
Expected: FAIL because types/method are missing.

---

## Task 2: Implement minimal OpenAPI work item list client

**Files:**
- Create: `internal/openapi/work_item_list.go`
- Test: `test/openapitest/work_item_list_test.go`

- [ ] **Step 1: Add minimal response and item types**

```go
type WorkItemListResponse struct {
    ErrCode    int              `json:"err_code"`
    ErrMsg     string           `json:"err_msg"`
    Err        any              `json:"err"`
    Data       []map[string]any `json:"data"`
    Pagination map[string]any   `json:"pagination"`
}
```

- [ ] **Step 2: Add the minimal method**

```go
func (c *Client) ListWorkItems(ctx context.Context, projectKey string, payload map[string]any) (*WorkItemListResponse, error) {
    // validate projectKey and payload
    // POST open_api/{projectKey}/work_item/filter
}
```

- [ ] **Step 3: Run focused tests to verify they pass**

Run: `go test ./test/openapitest -run 'TestClient_ListWorkItems_' -v`
Expected: PASS.

- [ ] **Step 4: Run package tests for regression**

Run: `go test ./test/openapitest -v`
Expected: PASS.

---

## Task 3: Add failing OpenAPI tests for work item detail

**Files:**
- Create: `test/openapitest/work_item_detail_test.go`
- Test: `test/openapitest/work_item_detail_test.go`

- [ ] **Step 1: Write the failing test for the detail endpoint path**

```go
func TestClient_GetWorkItemDetail_UsesQueryEndpoint(t *testing.T) {
    // Expect POST /open_api/p_demo/work_item/issue/query
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./test/openapitest -run TestClient_GetWorkItemDetail_UsesQueryEndpoint -v`
Expected: FAIL because `GetWorkItemDetail` does not exist yet.

- [ ] **Step 3: Write the failing test for request payload**

```go
func TestClient_GetWorkItemDetail_SendsIDsAndFields(t *testing.T) {
    // Expect work_item_ids and fields in request body
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `go test ./test/openapitest -run 'TestClient_GetWorkItemDetail_' -v`
Expected: FAIL because the client is not implemented.

---

## Task 4: Implement minimal OpenAPI work item detail client

**Files:**
- Create: `internal/openapi/work_item_detail.go`
- Test: `test/openapitest/work_item_detail_test.go`

- [ ] **Step 1: Add minimal response type**

```go
type WorkItemDetailResponse struct {
    ErrCode int              `json:"err_code"`
    ErrMsg  string           `json:"err_msg"`
    Err     any              `json:"err"`
    Data    []map[string]any `json:"data"`
}
```

- [ ] **Step 2: Add the minimal method**

```go
func (c *Client) GetWorkItemDetail(ctx context.Context, projectKey, workItemTypeKey string, payload map[string]any) (*WorkItemDetailResponse, error) {
    // validate required values
    // POST open_api/{projectKey}/work_item/{workItemTypeKey}/query
}
```

- [ ] **Step 3: Run focused tests to verify they pass**

Run: `go test ./test/openapitest -run 'TestClient_GetWorkItemDetail_' -v`
Expected: PASS.

- [ ] **Step 4: Run package tests for regression**

Run: `go test ./test/openapitest -v`
Expected: PASS.

---

## Task 5: Add failing command-tree test for TUI shortcut entry

**Files:**
- Create: `test/cmdtest/project_work_item_tui_cmd_test.go`
- Test: `test/cmdtest/project_work_item_tui_cmd_test.go`

- [ ] **Step 1: Write the failing test for the command flag**

```go
func TestProjectWorkItemCommand_HasTUIFlag(t *testing.T) {
    // Assert project work-item command has -t/--tui
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./test/cmdtest -run TestProjectWorkItemCommand_HasTUIFlag -v`
Expected: FAIL because the flag is not registered yet.

---

## Task 6: Implement TUI shortcut command wiring and Entry B fallback behavior

**Files:**
- Modify: `cmd/project/work_item.go`
- Test: `test/cmdtest/project_work_item_tui_cmd_test.go`

- [ ] **Step 1: Add `--tui` flag to `work-item` command**

```go
var tuiMode bool
cmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Launch interactive TUI")
```

- [ ] **Step 2: Add minimal RunE behavior for TUI entry**

- If `--tui` is set, invoke a new TUI runner for work-item flow.
- Keep non-TUI behavior unchanged.

- [ ] **Step 3: Add failing command tests for Entry B bootstrap matrix**

```go
func TestProjectWorkItemTUI_DirectProjectKeyEntryOpensTypeList(t *testing.T) {}
func TestProjectWorkItemTUI_DirectProjectAndTypeEntryOpensWorkItemList(t *testing.T) {}
func TestProjectWorkItemTUI_MissingProjectKeyFallsBackToInteractiveBootstrap(t *testing.T) {}
func TestProjectWorkItemTUI_InvalidProjectKeyFallsBackToProjectList(t *testing.T) {}
func TestProjectWorkItemTUI_InvalidWorkItemTypeKeyFallsBackToTypeList(t *testing.T) {}
```

- [ ] **Step 4: Run focused command tests to verify they fail**

Run: `go test ./test/cmdtest -run 'TestProjectWorkItem(TUI_|Command_HasTUIFlag)' -v`
Expected: FAIL until bootstrap behavior is implemented.

- [ ] **Step 5: Implement the bootstrap targets in the TUI entry wiring**

- provided `project_key` => direct to type list
- provided `project_key` + `work_item_type_key` => direct to work item list
- missing `project_key` => interactive project list
- invalid `project_key` => project error page with deterministic `q`/`r` target = project list
- invalid `work_item_type_key` => type error page with deterministic `q`/`r` target = type list

- [ ] **Step 6: Run focused command tests to verify they pass**

Run: `go test ./test/cmdtest -run 'TestProjectWorkItem(TUI_|Command_HasTUIFlag)' -v`
Expected: PASS.

- [ ] **Step 7: Run command package regression tests**

Run: `go test ./test/cmdtest -v`
Expected: PASS.

---

## Task 7: Add failing TUI tests for type-list to work-item-list flow

**Files:**
- Modify: `internal/tui/root_state_test.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Write the failing state transition test**

```go
func TestRootModel_WorkItemTypeList_EnterOpensWorkItemList(t *testing.T) {
    // stateWorkItemTypeList + selected type + enter => stateWorkItemList
}
```

- [ ] **Step 2: Write the failing view test for work item list**

```go
func TestRootModel_WorkItemListView_ShowsRowsAndHints(t *testing.T) {
    // item rows and footer hints visible
}
```

- [ ] **Step 3: Run focused tests to verify they fail**

Run: `go test ./internal/tui -run 'TestRootModel_WorkItem(TypeList_EnterOpensWorkItemList|ListView_ShowsRowsAndHints)' -v`
Expected: FAIL because state/view do not exist yet.

---

## Task 8: Implement minimal work item list state and view

**Files:**
- Modify: `internal/tui/root.go`
- Create: `internal/tui/work_item_list.go`
- Modify: `internal/tui/work_item_type.go`
- Modify: `internal/tui/style.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Add new state and rootModel fields**

Add at minimum:

```go
stateWorkItemList
workItems []map[string]any
workItemListErr error
workItemPageNum int
workItemPageSize int
workItemTotal int
workItemHasNextPage bool
workItemFilterMode bool
workItemFilterQuery string
workItemFilteredIdx []int
```

- [ ] **Step 2: Add minimal routing in `View()` and `Update()`**

- `enter` from work item type list opens work item list state.
- `q` from work item list returns to type list.

- [ ] **Step 3: Add `renderWorkItemListView` in `internal/tui/work_item_list.go`**

Render minimal rows:

```go
[id] title
```

Plus status/assignee when present, and footer hints for select/filter/back.

- [ ] **Step 4: Run focused tests to verify they pass**

Run: `go test ./internal/tui -run 'TestRootModel_WorkItem(TypeList_EnterOpensWorkItemList|ListView_ShowsRowsAndHints)' -v`
Expected: PASS.

---

## Task 9: Add failing TUI tests for work item detail flow

**Files:**
- Modify: `internal/tui/root_state_test.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Write the failing test for list -> detail transition**

```go
func TestRootModel_WorkItemList_EnterOpensDetail(t *testing.T) {
    // stateWorkItemList + enter => stateWorkItemDetail
}
```

- [ ] **Step 2: Write the failing test for detail summary/raw toggle**

```go
func TestRootModel_WorkItemDetail_TabTogglesRawView(t *testing.T) {
    // tab toggles detail mode
}
```

- [ ] **Step 3: Write the failing test for back-navigation source restoration**

```go
func TestRootModel_WorkItemDetail_QReturnsToList(t *testing.T) {
    // q returns to list preserving cursor
}
```

- [ ] **Step 4: Run focused tests to verify they fail**

Run: `go test ./internal/tui -run 'TestRootModel_WorkItem(Detail_|List_EnterOpensDetail)' -v`
Expected: FAIL because detail state is not implemented.

---

## Task 10: Implement minimal work item detail state and view

**Files:**
- Modify: `internal/tui/root.go`
- Create: `internal/tui/work_item_detail.go`
- Modify: `internal/tui/work_item_list.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Add detail state fields**

Add at minimum:

```go
stateWorkItemDetail
selectedWorkItem map[string]any
workItemDetail map[string]any
workItemDetailErr error
workItemDetailSource string
workItemDetailRaw bool
```

- [ ] **Step 2: Implement list -> detail navigation and `q` return**

- `enter` on list opens detail.
- `q` from detail returns to list.
- Re-entry resets to summary mode.

- [ ] **Step 3: Implement minimal summary/raw rendering**

Summary should show:
- id
- name/title
- status if present
- assignee/owner if present
- created_at / updated_at if present
- description if present

Raw view can start as pretty-printed JSON text.

- [ ] **Step 4: Run focused tests to verify they pass**

Run: `go test ./internal/tui -run 'TestRootModel_WorkItem(Detail_|List_EnterOpensDetail)' -v`
Expected: PASS.

---

## Task 11: Add failing TUI tests for search builder semantics and keyboard precedence

**Files:**
- Modify: `internal/tui/root_state_test.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Write failing tests for raw-mode precedence**

```go
func TestRootModel_SearchBuilder_RawJSONDisablesStructuredExecution(t *testing.T) {}
func TestRootModel_SearchBuilderView_ShowsRawModeIndicator(t *testing.T) {}
```

- [ ] **Step 2: Write failing tests for builder validation errors**

```go
func TestRootModel_SearchBuilder_InvalidDateBlocksExecution(t *testing.T) {}
func TestRootModel_SearchBuilder_InvalidRawJSONBlocksExecution(t *testing.T) {}
```

- [ ] **Step 3: Write failing tests for keyboard precedence**

```go
func TestRootModel_SearchBuilder_QReturnsInsteadOfEditing(t *testing.T) {}
func TestRootModel_SearchBuilder_JKAreTextInputInEditableField(t *testing.T) {}
func TestRootModel_FilterMode_JKNavigateWithoutChangingQuery(t *testing.T) {}
```

- [ ] **Step 4: Run focused tests to verify they fail**

Run: `go test ./internal/tui -run 'TestRootModel_(SearchBuilder_|FilterMode_)' -v`
Expected: FAIL until raw precedence, validation, and keyboard rules are implemented.

---

## Task 12: Implement search builder semantics and keyboard precedence

**Files:**

- Create: `internal/tui/work_item_search.go`
- Modify: `internal/tui/root.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Implement raw-json precedence behavior**

- non-empty raw json => raw mode
- structured fields ignored for execution in raw mode
- explicit raw-mode indicator visible in builder view

- [ ] **Step 2: Implement builder-side validation before execution**

- invalid date blocks execution and shows diagnostic error
- invalid raw json blocks execution and shows diagnostic error

- [ ] **Step 3: Implement documented keyboard precedence**

- `q` returns, not text input
- `j/k` are text in builder editable fields
- existing filter mode keeps `j/k` navigation semantics

- [ ] **Step 4: Run focused tests to verify they pass**

Run: `go test ./internal/tui -run 'TestRootModel_(SearchBuilder_|FilterMode_)' -v`
Expected: PASS.

---

## Task 13: Add failing TUI tests for search results -> detail flow

**Files:**
- Modify: `internal/tui/root_state_test.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Write the failing test for results -> detail transition**

```go
func TestRootModel_SearchResults_EnterOpensResultDetail(t *testing.T) {
    // enter from results opens detail state/source=searchResults
}
```

- [ ] **Step 2: Write the failing test for detail -> results back behavior**

```go
func TestRootModel_SearchResultDetail_QReturnsToResults(t *testing.T) {
    // q returns to search results preserving cursor
}
```

- [ ] **Step 3: Run focused tests to verify they fail**

Run: `go test ./internal/tui -run 'TestRootModel_Search(Result|Results)_' -v`
Expected: FAIL because this flow is not implemented yet.

---

## Task 14: Implement minimal search result detail reuse

**Files:**
- Modify: `internal/tui/root.go`
- Create: `internal/tui/work_item_search.go`
- Modify: `internal/tui/work_item_detail.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Add `stateWorkItemSearchResultDetail` and source-aware transitions**

- `enter` from search results opens detail.
- `q` from detail returns to results.

- [ ] **Step 2: Reuse detail rendering instead of duplicating it**

Use one rendering path and branch footer/back target by `workItemDetailSource`.

- [ ] **Step 3: Run focused tests to verify they pass**

Run: `go test ./internal/tui -run 'TestRootModel_Search(Result|Results)_' -v`
Expected: PASS.

---

## Task 15: Add failing tests for loading/error/retry behavior

**Files:**
- Modify: `internal/tui/root_state_test.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Write failing tests for loading and error views**

```go
func TestRootModel_WorkItemListView_ShowsLoadingState(t *testing.T) {}
func TestRootModel_WorkItemDetailView_ShowsErrorAndRetryHint(t *testing.T) {}
```

- [ ] **Step 2: Write failing test for retry context**

```go
func TestRootModel_ErrorState_RRetriesLastRequest(t *testing.T) {}
```

- [ ] **Step 3: Run focused tests to verify they fail**

Run: `go test ./internal/tui -run 'TestRootModel_(WorkItem.*(Loading|Error)|ErrorState_RRetriesLastRequest)' -v`
Expected: FAIL because loading/error/retry are not fully wired.

---

## Task 16: Implement minimal loading/error/retry contract

**Files:**
- Modify: `internal/tui/root.go`
- Modify: `internal/tui/work_item_list.go`
- Modify: `internal/tui/work_item_detail.go`
- Modify: `internal/tui/work_item_search.go`
- Modify: `internal/tui/style.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Add last-request context fields to model**

```go
lastRequestKind string
lastRequestProjectKey string
lastRequestWorkItemTypeKey string
lastRequestWorkItemID any
lastRequestPayload map[string]any
lastRequestPageNum int
lastRequestPageSize int
```

- [ ] **Step 2: Add loading and error views with footer hints**

At minimum show:
- what is loading
- what failed
- `[R] retry`
- `[Q] back`

- [ ] **Step 3: Add `r` handling that replays the last failed request only**

- [ ] **Step 4: Run focused tests to verify they pass**

Run: `go test ./internal/tui -run 'TestRootModel_(WorkItem.*(Loading|Error)|ErrorState_RRetriesLastRequest)' -v`
Expected: PASS.

---

## Task 17: Add failing tests for pagination and restoration rules

**Files:**
- Modify: `internal/tui/root_state_test.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Write failing test for results/list page metadata rendering**

```go
func TestRootModel_WorkItemListView_ShowsPaginationHint(t *testing.T) {}
```

- [ ] **Step 2: Write failing test for restoration after back-navigation**

```go
func TestRootModel_WorkItemDetail_BackRestoresCursorAndPage(t *testing.T) {}
func TestRootModel_SearchBuilder_BackRestoresInputsAndFocus(t *testing.T) {}
```

- [ ] **Step 3: Run focused tests to verify they fail**

Run: `go test ./internal/tui -run 'TestRootModel_(WorkItemListView_ShowsPaginationHint|WorkItemDetail_BackRestoresCursorAndPage|SearchBuilder_BackRestoresInputsAndFocus)' -v`
Expected: FAIL until restoration contract is implemented.

---

## Task 18: Implement pagination metadata and restoration behavior

**Files:**
- Modify: `internal/tui/root.go`
- Modify: `internal/tui/work_item_list.go`
- Modify: `internal/tui/work_item_search.go`
- Test: `internal/tui/root_state_test.go`

- [ ] **Step 1: Add footer rendering for page state**

Show at minimum:
- current page
- total or loaded count
- no-more-data hint when applicable

- [ ] **Step 2: Ensure back-navigation restores stored cursor/page/filter/focus state**

- list -> detail -> list restores cursor/page/filter
- builder -> results -> builder restores inputs/focus/errors
- results -> detail -> results restores cursor/page

- [ ] **Step 3: Run focused tests to verify they pass**

Run: `go test ./internal/tui -run 'TestRootModel_(WorkItemListView_ShowsPaginationHint|WorkItemDetail_BackRestoresCursorAndPage|SearchBuilder_BackRestoresInputsAndFocus)' -v`
Expected: PASS.

---

## Task 19: Run targeted regression suites

**Files:**
- Test only

- [ ] **Step 1: Run TUI package tests**

Run: `go test ./internal/tui -v`
Expected: PASS.

- [ ] **Step 2: Run external TUI regression tests**

Run: `go test ./test/tuitest -v`
Expected: PASS.

- [ ] **Step 3: Run command tests**

Run: `go test ./test/cmdtest -v`
Expected: PASS.

- [ ] **Step 4: Run OpenAPI tests**

Run: `go test ./test/openapitest -v`
Expected: PASS.

---

## Task 20: Run full verification

**Files:**
- Test only

- [ ] **Step 1: Run full Go test suite**

Run: `go test ./...`
Expected: PASS.

- [ ] **Step 2: Manual verification for root TUI path**

Run: `./lark`
Expected:
- Can reach project list
- Can navigate to project detail
- Can open work item type list
- Can open work item list
- Can open detail

- [ ] **Step 3: Manual verification for search path**

Run: `./lark`
Expected:
- Can open search builder
- Can move to results
- Can open result detail
- `q` returns correctly

- [ ] **Step 4: Manual verification for shortcut path**

Run: `./lark project work-item -t`
Expected:
- Enters the same TUI work-item flow
- Missing `project_key` falls back to interactive bootstrap instead of immediate hard failure

---

## Notes for the implementing agent

- Use `@superpowers:test-driven-development` literally for each behavior change: write test, watch it fail, implement minimal code, run it green.
- Keep `internal/tui/root.go` small by moving page-specific rendering and update helpers into new files immediately when adding a new page.
- Do not fix the unrelated CLI ergonomics in this plan.
- Prefer `map[string]any` for detail/list payloads at first; only extract helpers for stable summary fields actually needed by the TUI.
- Keep tests in the dedicated `test/` root or existing TUI test file per project preference.

## Next Phase Backlog (not included in this plan)

The user explicitly identified these as the next-stage scope after the current TUI full-flow work:

- CLI 参数体系修复
- 跨空间工作项聚合浏览
- search builder 的复杂表达式编辑器
- 结果页多维排序、批量操作、内联编辑
- 过度抽象成多 tea.Model 架构（如后续复杂度继续上升再评估）
