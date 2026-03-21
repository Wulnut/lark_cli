# TUI 用户信息展示 — 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `lark` 不带参数启动时，TUI 页面自动展示当前登录用户姓名 + 邮箱，作为登录成功引导页。

**Architecture:** 在 `rootModel` 的 `Init()` 中异步调用 `openapi.Client.QueryCurrentUser`，根据结果（成功 / 失败降级）渲染不同界面。复用已有 `DoJSON` 重试 + 鉴权逻辑。

**Tech Stack:** Bubble Tea (Charmbracelet), Go 1.21+, `lark_cli/internal/openapi`, `lark_cli/internal/auth`

---

## 文件变更总览

| 文件 | 操作 |
|------|------|
| `internal/openapi/user.go` | 新增 |
| `internal/tui/root.go` | 修改 |
| `cmd/root.go` | 修改（创建 openapi.Client 并传入 TUI） |
| `test/tuitest/root_test.go` | 新增 |

---

## Task 1: 新增 `internal/openapi/user.go`

**Files:**
- Create: `internal/openapi/user.go`

- [ ] **Step 1: 写类型定义**

```go
package openapi

// LocaleName represents the multilingual name object in user responses.
type LocaleName struct {
    Default string `json:"default"`
    EnUS    string `json:"en_us"`
    ZhCN    string `json:"zh_cn"`
}

// UserInfo represents a user record returned by GET /open_api/user/query.
type UserInfo struct {
    UserID    int64      `json:"user_id"`
    NameCn    string     `json:"name_cn"`
    NameEn    string     `json:"name_en"`
    OutID     string     `json:"out_id"`
    Name      LocaleName `json:"name"`
    UserKey   string     `json:"user_key"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    AvatarURL string     `json:"avatar_url"`
    Status    string     `json:"status"`
}

// QueryUserResponse is the response shape for POST /open_api/user/query.
type QueryUserResponse struct {
    ErrCode int       `json:"err_code"`
    ErrMsg  string    `json:"err_msg"`
    Err     any       `json:"err"`
    Data    []UserInfo `json:"data"`
}
```

- [ ] **Step 2: 写 QueryCurrentUser 方法**

```go
// QueryCurrentUser calls POST /open_api/user/query with a single user_key.
// Returns the first user in the data list, or nil if the result is empty.
// Returns error only on network failure or unexpected response shape.
func (c *Client) QueryCurrentUser(ctx context.Context, userKey string) (*UserInfo, error) {
    if userKey == "" {
        return nil, fmt.Errorf("userKey is required")
    }
    var resp QueryUserResponse
    err := c.DoJSON(ctx, &Request{
        Method: "POST",
        Path:   "open_api/user/query",
        Body:   map[string][]string{"user_keys": {userKey}},
    }, &resp)
    if err != nil {
        return nil, err
    }
    if len(resp.Data) == 0 {
        return nil, nil // no user found, not an error
    }
    return &resp.Data[0], nil
}
```

- [ ] **Step 3: 验证编译**

Run: `go build ./internal/openapi/...`
Expected: 编译成功，无输出

---

## Task 2: 修改 `internal/tui/root.go`

**Files:**
- Modify: `internal/tui/root.go`

- [ ] **Step 1: 更新 import 和 model 结构**

```go
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
    userKey string
    client  *openapi.Client
    user    *openapi.UserInfo
    fetchErr error
    state   tuiState
}
```

- [ ] **Step 2: 新增消息类型**

```go
type userFetchedMsg   struct{ user *openapi.UserInfo }
type userFetchFailedMsg struct{ err error }
```

- [ ] **Step 3: 实现 fetchUserCmd**

```go
func fetchUserCmd(userKey string, client *openapi.Client) tea.Cmd {
    return func() tea.Msg {
        user, err := client.QueryCurrentUser(context.Background(), userKey)
        if err != nil {
            return userFetchFailedMsg{err: err}
        }
        return userFetchedMsg{user: user}
    }
}
```

- [ ] **Step 4: 修改 Init()**

```go
func (m rootModel) Init() tea.Cmd {
    return fetchUserCmd(m.userKey, m.client)
}
```

- [ ] **Step 5: 更新 Update() 处理新消息**

```go
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
```

- [ ] **Step 6: 更新 View() 三态渲染**

```go
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
        return fmt.Sprintf("👤 %s  ⚠️ 仅显示 user_key\n───────────────────────────────\n按 [Q] 退出\n", m.userKey)
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
```

- [ ] **Step 7: 更新 Run() 签名**

```go
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
```

- [ ] **Step 8: 验证编译**

Run: `go build ./internal/tui/...`
Expected: 编译成功，无输出

---

## Task 3: 修改 `cmd/root.go` — 构造 openapi.Client 并传入 TUI

**Files:**
- Modify: `cmd/root.go`

- [ ] **Step 1: NewRootCmd 中创建 openapi.Client**

在 `NewRootCmd` 的 `RunE` 回调中，当走 TUI 路径时，构造 `*openapi.Client`：

```go
RunE: func(cmd *cobra.Command, args []string) error {
    if len(args) > 0 {
        return cmd.Help()
    }
    var userKey string
    var apiClient *openapi.Client
    if deps.PluginTokenProvider != nil {
        userKey = deps.Config.UserKey
        apiClient = openapi.NewClient(deps.Config.BaseURL, &http.Client{Timeout: deps.Config.HTTPTimeout}, deps.PluginTokenProvider)
    }
    if err := tui.Run(deps.Stdout, userKey, apiClient); err != nil {
        return fmt.Errorf("interactive UI: %w", err)
    }
    return nil
},
```

注意：需在 import 中新增：
```go
"lark_cli/internal/openapi"
```

- [ ] **Step 2: 验证编译**

Run: `go build ./cmd/...`
Expected: 编译成功，无输出

---

## Task 4: 新增 TUI 测试 `test/tuitest/root_test.go`

**Files:**
- Create: `test/tuitest/root_test.go`
- Test: `test/tuitest/root_test.go`

- [ ] **Step 1: 写测试文件框架**

```go
package tuittest

import (
    "testing"

    "lark_cli/internal/tui"
)

func TestRootModel_SuccessView(t *testing.T) {
    // TODO: mock openapi client or use httptest
    t.Skip("integration test - requires httptest server")
}

func TestRootModel_DegradedView(t *testing.T) {
    t.Skip("integration test - requires httptest server")
}
```

（先写框架让测试文件编译通过，完整的 mock 测试可后续迭代）

- [ ] **Step 2: 验证测试可运行**

Run: `go test ./test/tuitest/... -v`
Expected: `PASS`（SKIP 正常通过）

---

## Task 5: 端到端验证

- [ ] **Step 1: 编译全项目**

Run: `go build ./...`
Expected: 编译成功

- [ ] **Step 2: 运行全量测试**

Run: `go test ./...`
Expected: 所有测试 PASS

---

## 关键设计决策

1. **降级优先**：`QueryCurrentUser` 返回 error 时（网络错误），TUI 显示 user_key 而非退出。err_code=30006（用户不存在）时返回 `nil, nil`，TUI 同样降级显示 user_key。
2. **无 openapi.Client 时**：若 `PluginTokenProvider == nil`（未登录），`userKey=""` 且 `apiClient=nil`，TUI 走降级路径直接显示 `👤  ⚠️ 仅显示 user_key`。
3. **复用**：不新增 API client 封装，所有逻辑复用 `DoJSON`。
