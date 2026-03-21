# TUI 用户信息展示设计

## Context

用户以 `lark` 不带参数启动 CLI 时会进入交互 TUI 页面。首个功能是展示当前登录用户的基本信息，作为"登录成功引导页"——让用户知道自己已登录并显示真实身份（姓名 + 邮箱），而非仅显示冷冰冰的 user_key。

## 目标

- 用户运行 `lark` 后，立即通过 API 获取并展示姓名 + 邮箱
- 界面简洁，作为登录成功确认页
- API 失败时降级展示，只显示 user_key 不阻塞

## 涉及文件

| 文件 | 操作 |
|------|------|
| `internal/openapi/user.go` | 新增 |
| `internal/tui/root.go` | 修改 |
| `test/tuitest/user_model_test.go` | 新增（可选） |

## API 规范

**接口：** `POST /open_api/user/query`
**鉴权：** `X-Plugin-Token` + `X-User-Key`（由 OpenAPI Client 自动注入）

**请求体：**
```json
{
  "user_keys": ["<user_key from config>"]
}
```

**响应字段（使用）：**
- `name.zh_cn` 或 `name.default` → 显示姓名
- `email` → 显示邮箱
- `status` → 映射为 ✅ 已激活 / ⚠️ 未激活

**特殊错误码处理：**
- `err_code == 30006`（用户不存在）：降级显示 user_key
- 其他 API 错误：降级显示 user_key + 提示

## TUI View 渲染

### 加载中
```
🔄 正在获取用户信息...
```

### 成功
```
👤 张三 <user@company.com>  ✅ 已登录
──────────────────────────────────────
按 [Q] 退出
```

### 降级（API 失败）
```
👤 <user_key>  ⚠️ 仅显示 user_key
──────────────────────────────────────
按 [Q] 退出
```

## 组件设计

### internal/openapi/user.go（新增）

```go
package openapi

// QueryUserRequest / QueryUserResponse / UserInfo
// LocaleName structs (同 API 文档字段)

// QueryCurrentUser calls POST /open_api/user/query with a single user_key.
// Returns the first user in data list or nil if empty.
// Returns error only on network / unexpected failure.
func (c *Client) QueryCurrentUser(ctx context.Context, userKey string) (*UserInfo, error)
```

### internal/tui/root.go（修改）

- `Init()` 中新增 `userStatusCmd`：异步调用 `client.QueryCurrentUser`
- 新增 `userFetchedMsg(user *openapi.UserInfo)` 消息类型
- 新增 `userFetchFailedMsg(err error)` 消息类型
- `View()` 根据 model 状态渲染上述三种界面之一
- 移除/简化占位文字，替换为真实用户信息展示

### 降级逻辑

```go
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case userFetchedMsg:
        m.user = msg.user
        m.state = stateSuccess
        return m, nil
    case userFetchFailedMsg:
        m.fetchError = msg.err
        m.state = stateDegraded
        return m, nil
    // ...
    }
}
```

## 测试策略

- Mock `openapi.Client`（或 fake token provider + httptest）验证 View 输出包含预期字符串
- 测试成功/失败/降级三种状态渲染
