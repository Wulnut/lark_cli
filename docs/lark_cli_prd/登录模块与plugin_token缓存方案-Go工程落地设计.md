---
title: 登录模块与 plugin_token 缓存方案 - Go 工程落地设计
source:
author:
published:
created: 2026-03-19
updated: 2026-03-19
description: lark_cli plugin_token 缓存方案的Go工程落地设计，包含模块目录、接口清单、token_provider实现、请求中间层及错误策略等Go skeleton代码
tags: [lark_cli, PRD, Go, skeleton, auth, token, cache, plugin_token, OpenAPI]
category: lark_cli
status: draft
related_docs:
  - "[[lark_cli PRD/登录模块与plugin_token缓存方案]]"
  - "[[飞书项目/飞书项目OpenAPI获取访问凭证]]"
  - "[[飞书项目/飞书项目OpenAPI鉴权流程]]"
---

# 登录模块与 plugin_token 缓存方案 - Go 工程落地设计

本文是 [[lark_cli PRD/登录模块与plugin_token缓存方案]] 的工程落地补充稿，目标是将前文中的设计结论拆解为：

1. 接口清单
2. 文件级实现顺序
3. 更接近真实代码的 Go skeleton

---

# 一、开发任务拆解

## 1. 模块与目录建议

```text
internal/
  config/
    config.go
  session/
    session.go
  auth/
    fingerprint.go
    token_provider.go
    auth_api.go
  openapi/
    client.go
    request.go
    response.go
    middleware.go
    errors/
      types.go
      policies.go
      normalize.go
```

职责如下：

- `config`：读取与保存 `.lark/config.json`
- `session`：读取与保存 `.lark/session.json`
- `auth/fingerprint`：生成配置指纹
- `auth/token_provider`：获取、缓存、刷新 `plugin_token`
- `auth/auth_api`：调用 `/authen/plugin_token`
- `openapi/client`：统一发起 API 请求
- `openapi/request`：定义请求结构
- `openapi/response`：定义通用响应结构
- `openapi/middleware`：统一注入认证头与重试策略
- `openapi/errors`：统一错误分类、策略表、错误规范化

---

## 2. 接口清单

## 2.1 config 模块

### `Config`

```go
type Config struct {
    UserKey      string `json:"user_key"`
    PluginID     string `json:"plugin_id"`
    PluginSecret string `json:"plugin_secret"`
    BaseURL      string `json:"base_url"`
}
```

### `Store`

```go
type Store interface {
    Load(ctx context.Context) (*Config, error)
    Save(ctx context.Context, cfg *Config) error
}
```

### 关键方法

```go
func DefaultConfig() *Config
func (c *Config) ValidateForOpenAPI() error
```

职责：
- 提供默认配置
- 校验调用 OpenAPI 所需的关键字段是否完整

---

## 2.2 session 模块

### `Session`

```go
type Session struct {
    PluginToken       string `json:"plugin_token"`
    ExpireAt          int64  `json:"expire_at"`
    ObtainedAt        int64  `json:"obtained_at"`
    ConfigFingerprint string `json:"config_fingerprint"`
}
```

### `Store`

```go
type Store interface {
    Load(ctx context.Context) (*Session, error)
    Save(ctx context.Context, sess *Session) error
    Clear(ctx context.Context) error
}
```

### 关键方法

```go
func (s *Session) IsEmpty() bool
func (s *Session) IsValid(now int64, fingerprint string) bool
```

职责：
- 表示 token 缓存状态
- 判断缓存是否可复用

---

## 2.3 auth/fingerprint 模块

### 关键方法

```go
func BuildConfigFingerprint(cfg *config.Config) string
```

职责：
- 根据 `user_key + plugin_id + plugin_secret + base_url` 生成配置指纹
- 用于判定 `session.json` 中缓存是否仍与当前配置匹配

---

## 2.4 auth/auth_api 模块

### `PluginTokenResponse`

```go
type PluginTokenResponse struct {
    Token      string
    ExpireTime int64
}
```

### `AuthAPI`

```go
type AuthAPI interface {
    FetchPluginToken(ctx context.Context, pluginID, pluginSecret string) (*PluginTokenResponse, error)
}
```

职责：
- 封装调用 `/authen/plugin_token` 的 HTTP 细节
- 返回统一 token 响应结构

---

## 2.5 auth/token_provider 模块

### `AuthContext`

```go
type AuthContext struct {
    UserKey     string
    PluginToken string
}
```

### `TokenProvider`

```go
type TokenProvider interface {
    GetAuthContext(ctx context.Context) (*AuthContext, error)
    ForceRefresh(ctx context.Context) (*AuthContext, error)
}
```

职责：
- 对外暴露“获取当前可用认证上下文”的统一入口
- 自动处理本地缓存、过期、刷新

---

## 2.6 openapi/errors 模块

### `Category`

```go
type Category string
```

推荐枚举：

```go
const (
    CategoryAuth          Category = "auth"
    CategoryPermission    Category = "permission"
    CategoryParam         Category = "param"
    CategoryNotFound      Category = "not_found"
    CategoryStateConflict Category = "state_conflict"
    CategoryRateLimit     Category = "rate_limit"
    CategoryServer        Category = "server"
    CategoryUnknown       Category = "unknown"
)
```

### `Policy`

```go
type Policy struct {
    Code         int
    Message      string
    Category     Category
    Retryable    bool
    RefreshToken bool
    MaxRetry     int
    UserHint     string
    DevHint      string
}
```

### `OpenAPIError`

```go
type OpenAPIError struct {
    HTTPStatus int
    Code       int
    Message    string
    Category   Category
    UserHint   string
    DevHint    string
    RawBody    []byte
}
```

### 关键函数

```go
func LookupPolicy(code int) Policy
func NormalizeError(httpStatus int, body []byte, code int, msg string) error
```

职责：
- 根据错误码查策略
- 将响应规范化为统一错误对象

---

## 2.7 openapi/client 模块

### `Request`

```go
type Request struct {
    Method  string
    Path    string
    Body    any
    Headers map[string]string
}
```

### `Client`

```go
type Client struct {
    BaseURL       string
    HTTPClient    *http.Client
    TokenProvider auth.TokenProvider
}
```

### 关键方法

```go
func (c *Client) DoJSON(ctx context.Context, req *Request, out any) error
func (c *Client) buildRequest(ctx context.Context, req *Request, authCtx *auth.AuthContext) (*http.Request, error)
func (c *Client) parseAndNormalize(httpStatus int, body []byte) (int, string, error)
```

职责：
- 统一构造请求
- 自动注入认证头
- 统一执行刷新与重试策略

---

## 3. 文件级实现顺序

建议按如下顺序逐步落地，确保每一层都能单独验证。

## 第 1 阶段：配置与会话落盘

### 目标文件
- `internal/config/config.go`
- `internal/session/session.go`

### 目标
- 能读取/保存 `.lark/config.json`
- 能读取/保存 `.lark/session.json`
- 明确 config 与 session 的边界

### 验收点
- 手动编辑 `config.json` 后程序能正常加载
- 删除 `session.json` 不会造成致命错误

---

## 第 2 阶段：配置指纹与 token 获取接口

### 目标文件
- `internal/auth/fingerprint.go`
- `internal/auth/auth_api.go`

### 目标
- 生成配置指纹
- 打通 `/authen/plugin_token` 获取逻辑

### 验收点
- 能通过 `plugin_id + plugin_secret` 拿到 `plugin_token`
- 配置变更时指纹发生变化

---

## 第 3 阶段：TokenProvider

### 目标文件
- `internal/auth/token_provider.go`

### 目标
- 实现 token 缓存读取
- 识别过期
- 自动刷新
- 写回 `session.json`

### 验收点
- 首次请求自动获取 token
- token 有效时优先复用缓存
- token 即将过期时自动刷新
- 配置变更后自动废弃旧 token

---

## 第 4 阶段：错误码体系

### 目标文件
- `internal/openapi/errors/types.go`
- `internal/openapi/errors/policies.go`
- `internal/openapi/errors/normalize.go`

### 目标
- 建立统一错误分类
- 建立错误码策略表
- 统一包装 OpenAPIError

### 验收点
- 能识别 `10021/10022/10211/10301/10429/50006`
- 能根据错误码给出策略信息

---

## 第 5 阶段：统一请求客户端

### 目标文件
- `internal/openapi/request.go`
- `internal/openapi/response.go`
- `internal/openapi/client.go`
- `internal/openapi/middleware.go`

### 目标
- 自动注入 `X-PLUGIN-TOKEN`
- 自动注入 `X-USER-KEY`
- 对 token 相关错误执行刷新 + 重试一次
- 对限流 / 服务端错误执行短重试

### 验收点
- 业务命令无需自己处理 token
- 业务命令无需自己处理错误码分类

---

## 第 6 阶段：命令层接入

### 目标
- 现有 `user query` 等命令统一切到 `openapi.Client`
- CLI 输出支持用户提示与 debug 提示两层文案

### 验收点
- 业务命令代码显著简化
- 错误输出风格统一

---

# 二、Go skeleton（接近真实代码的完整框架）

下面按文件给出一版可直接作为起点的代码框架。

---

## 文件 1：`internal/config/config.go`

```go
package config

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
)

type Config struct {
    UserKey      string `json:"user_key"`
    PluginID     string `json:"plugin_id"`
    PluginSecret string `json:"plugin_secret"`
    BaseURL      string `json:"base_url"`
}

type Store interface {
    Load(ctx context.Context) (*Config, error)
    Save(ctx context.Context, cfg *Config) error
}

type FileStore struct {
    Path string
}

func DefaultConfig() *Config {
    return &Config{
        BaseURL: "https://project.feishu.cn/open_api",
    }
}

func (c *Config) ValidateForOpenAPI() error {
    if c == nil {
        return errors.New("config is nil")
    }
    if c.UserKey == "" {
        return fmt.Errorf("missing user_key")
    }
    if c.PluginID == "" {
        return fmt.Errorf("missing plugin_id")
    }
    if c.PluginSecret == "" {
        return fmt.Errorf("missing plugin_secret")
    }
    if c.BaseURL == "" {
        return fmt.Errorf("missing base_url")
    }
    return nil
}

func (s *FileStore) Load(ctx context.Context) (*Config, error) {
    _ = ctx

    cfg := DefaultConfig()
    data, err := os.ReadFile(s.Path)
    if err != nil {
        if os.IsNotExist(err) {
            return cfg, nil
        }
        return nil, err
    }

    if err := json.Unmarshal(data, cfg); err != nil {
        return nil, err
    }

    if cfg.BaseURL == "" {
        cfg.BaseURL = DefaultConfig().BaseURL
    }

    return cfg, nil
}

func (s *FileStore) Save(ctx context.Context, cfg *Config) error {
    _ = ctx

    if cfg == nil {
        return errors.New("config is nil")
    }

    if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
        return err
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return atomicWrite(s.Path, data, 0o600)
}

func atomicWrite(path string, data []byte, perm os.FileMode) error {
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, data, perm); err != nil {
        return err
    }
    return os.Rename(tmp, path)
}
```

---

## 文件 2：`internal/session/session.go`

```go
package session

import (
    "context"
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
)

type Session struct {
    PluginToken       string `json:"plugin_token"`
    ExpireAt          int64  `json:"expire_at"`
    ObtainedAt        int64  `json:"obtained_at"`
    ConfigFingerprint string `json:"config_fingerprint"`
}

type Store interface {
    Load(ctx context.Context) (*Session, error)
    Save(ctx context.Context, sess *Session) error
    Clear(ctx context.Context) error
}

type FileStore struct {
    Path string
}

func (s *Session) IsEmpty() bool {
    return s == nil || s.PluginToken == ""
}

func (s *Session) IsValid(now int64, fingerprint string) bool {
    if s == nil {
        return false
    }
    if s.PluginToken == "" {
        return false
    }
    if s.ConfigFingerprint != fingerprint {
        return false
    }
    if now >= s.ExpireAt {
        return false
    }
    return true
}

func (fs *FileStore) Load(ctx context.Context) (*Session, error) {
    _ = ctx

    data, err := os.ReadFile(fs.Path)
    if err != nil {
        if os.IsNotExist(err) {
            return &Session{}, nil
        }
        return nil, err
    }

    sess := &Session{}
    if err := json.Unmarshal(data, sess); err != nil {
        return nil, err
    }
    return sess, nil
}

func (fs *FileStore) Save(ctx context.Context, sess *Session) error {
    _ = ctx

    if sess == nil {
        return errors.New("session is nil")
    }

    if err := os.MkdirAll(filepath.Dir(fs.Path), 0o755); err != nil {
        return err
    }

    data, err := json.MarshalIndent(sess, "", "  ")
    if err != nil {
        return err
    }

    return atomicWrite(fs.Path, data, 0o600)
}

func (fs *FileStore) Clear(ctx context.Context) error {
    _ = ctx
    if err := os.Remove(fs.Path); err != nil && !os.IsNotExist(err) {
        return err
    }
    return nil
}

func atomicWrite(path string, data []byte, perm os.FileMode) error {
    tmp := path + ".tmp"
    if err := os.WriteFile(tmp, data, perm); err != nil {
        return err
    }
    return os.Rename(tmp, path)
}
```

---

## 文件 3：`internal/auth/fingerprint.go`

```go
package auth

import (
    "crypto/sha256"
    "encoding/hex"
    "strings"

    "your/module/internal/config"
)

func BuildConfigFingerprint(cfg *config.Config) string {
    if cfg == nil {
        return ""
    }

    raw := strings.Join([]string{
        cfg.UserKey,
        cfg.PluginID,
        cfg.PluginSecret,
        cfg.BaseURL,
    }, "|")

    sum := sha256.Sum256([]byte(raw))
    return hex.EncodeToString(sum[:])
}
```

---

## 文件 4：`internal/auth/auth_api.go`

```go
package auth

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
)

type PluginTokenResponse struct {
    Token      string
    ExpireTime int64
}

type AuthAPI interface {
    FetchPluginToken(ctx context.Context, pluginID, pluginSecret string) (*PluginTokenResponse, error)
}

type HTTPAuthAPI struct {
    BaseURL    string
    HTTPClient *http.Client
}

type fetchPluginTokenRequest struct {
    PluginID     string `json:"plugin_id"`
    PluginSecret string `json:"plugin_secret"`
}

type fetchPluginTokenResponse struct {
    Error struct {
        Code int    `json:"code"`
        Msg  string `json:"msg"`
    } `json:"error"`
    Data struct {
        Token      string `json:"token"`
        ExpireTime int64  `json:"expire_time"`
    } `json:"data"`
}

func (a *HTTPAuthAPI) FetchPluginToken(ctx context.Context, pluginID, pluginSecret string) (*PluginTokenResponse, error) {
    body := fetchPluginTokenRequest{
        PluginID:     pluginID,
        PluginSecret: pluginSecret,
    }

    payload, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }

    url := strings.TrimRight(a.BaseURL, "/") + "/authen/plugin_token"
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := a.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var parsed fetchPluginTokenResponse
    if err := json.Unmarshal(respBody, &parsed); err != nil {
        return nil, err
    }

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("token api http status: %d", resp.StatusCode)
    }
    if parsed.Error.Code != 0 {
        return nil, fmt.Errorf("token api error: code=%d msg=%s", parsed.Error.Code, parsed.Error.Msg)
    }

    return &PluginTokenResponse{
        Token:      parsed.Data.Token,
        ExpireTime: parsed.Data.ExpireTime,
    }, nil
}
```

---

## 文件 5：`internal/auth/token_provider.go`

```go
package auth

import (
    "context"
    "sync"
    "time"

    "your/module/internal/config"
    "your/module/internal/session"
)

type ConfigStore interface {
    Load(ctx context.Context) (*config.Config, error)
}

type SessionStore interface {
    Load(ctx context.Context) (*session.Session, error)
    Save(ctx context.Context, sess *session.Session) error
    Clear(ctx context.Context) error
}

type AuthContext struct {
    UserKey     string
    PluginToken string
}

type TokenProvider interface {
    GetAuthContext(ctx context.Context) (*AuthContext, error)
    ForceRefresh(ctx context.Context) (*AuthContext, error)
}

type CachedTokenProvider struct {
    ConfigStore  ConfigStore
    SessionStore SessionStore
    AuthAPI      AuthAPI

    SafetyWindow time.Duration
    mu           sync.Mutex
}

func (p *CachedTokenProvider) GetAuthContext(ctx context.Context) (*AuthContext, error) {
    cfg, err := p.ConfigStore.Load(ctx)
    if err != nil {
        return nil, err
    }
    if err := cfg.ValidateForOpenAPI(); err != nil {
        return nil, err
    }

    fingerprint := BuildConfigFingerprint(cfg)
    now := time.Now().Unix()

    sess, err := p.SessionStore.Load(ctx)
    if err != nil {
        return nil, err
    }

    if sess != nil && sess.IsValid(now, fingerprint) {
        return &AuthContext{
            UserKey:     cfg.UserKey,
            PluginToken: sess.PluginToken,
        }, nil
    }

    return p.refreshLocked(ctx, cfg, fingerprint)
}

func (p *CachedTokenProvider) ForceRefresh(ctx context.Context) (*AuthContext, error) {
    cfg, err := p.ConfigStore.Load(ctx)
    if err != nil {
        return nil, err
    }
    if err := cfg.ValidateForOpenAPI(); err != nil {
        return nil, err
    }

    fingerprint := BuildConfigFingerprint(cfg)
    return p.refreshLocked(ctx, cfg, fingerprint)
}

func (p *CachedTokenProvider) refreshLocked(ctx context.Context, cfg *config.Config, fingerprint string) (*AuthContext, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    now := time.Now().Unix()
    sess, err := p.SessionStore.Load(ctx)
    if err == nil && sess != nil && sess.IsValid(now, fingerprint) {
        return &AuthContext{
            UserKey:     cfg.UserKey,
            PluginToken: sess.PluginToken,
        }, nil
    }

    tokenResp, err := p.AuthAPI.FetchPluginToken(ctx, cfg.PluginID, cfg.PluginSecret)
    if err != nil {
        return nil, err
    }

    obtainedAt := time.Now().Unix()
    expireAt := obtainedAt + tokenResp.ExpireTime - int64(p.SafetyWindow.Seconds())
    if expireAt <= obtainedAt {
        expireAt = obtainedAt + tokenResp.ExpireTime
    }

    newSession := &session.Session{
        PluginToken:       tokenResp.Token,
        ObtainedAt:        obtainedAt,
        ExpireAt:          expireAt,
        ConfigFingerprint: fingerprint,
    }

    if err := p.SessionStore.Save(ctx, newSession); err != nil {
        return nil, err
    }

    return &AuthContext{
        UserKey:     cfg.UserKey,
        PluginToken: tokenResp.Token,
    }, nil
}
```

---

## 文件 6：`internal/openapi/errors/types.go`

```go
package errors

import "fmt"

type Category string

const (
    CategoryAuth          Category = "auth"
    CategoryPermission    Category = "permission"
    CategoryParam         Category = "param"
    CategoryNotFound      Category = "not_found"
    CategoryStateConflict Category = "state_conflict"
    CategoryRateLimit     Category = "rate_limit"
    CategoryServer        Category = "server"
    CategoryUnknown       Category = "unknown"
)

type Policy struct {
    Code         int
    Message      string
    Category     Category
    Retryable    bool
    RefreshToken bool
    MaxRetry     int
    UserHint     string
    DevHint      string
}

type OpenAPIError struct {
    HTTPStatus int
    Code       int
    Message    string
    Category   Category
    UserHint   string
    DevHint    string
    RawBody    []byte
}

func (e *OpenAPIError) Error() string {
    if e.Code != 0 {
        return fmt.Sprintf("openapi error: code=%d msg=%s", e.Code, e.Message)
    }
    return fmt.Sprintf("openapi error: http_status=%d", e.HTTPStatus)
}
```

---

## 文件 7：`internal/openapi/errors/policies.go`

```go
package errors

var Policies = map[int]Policy{
    10021: {
        Code:         10021,
        Message:      "Token Not Exist",
        Category:     CategoryAuth,
        Retryable:    true,
        RefreshToken: true,
        MaxRetry:     1,
        UserHint:     "未检测到 plugin_token，正在尝试重新获取。",
        DevHint:      "请检查 session 缓存或 X-PLUGIN-TOKEN 注入逻辑。",
    },
    10022: {
        Code:         10022,
        Message:      "Check Token Failed",
        Category:     CategoryAuth,
        Retryable:    true,
        RefreshToken: true,
        MaxRetry:     1,
        UserHint:     "plugin_token 校验失败，可能已过期，正在尝试重新获取。",
        DevHint:      "请检查 token 是否已过期或配置变更后仍使用旧 token。",
    },
    10211: {
        Code:         10211,
        Message:      "Token Info Is Invalid",
        Category:     CategoryAuth,
        Retryable:    true,
        RefreshToken: true,
        MaxRetry:     1,
        UserHint:     "plugin_token 信息不合法，正在尝试重新获取。",
        DevHint:      "也请检查是否缺少 X-USER-KEY。",
    },
    20039: {
        Code:         20039,
        Message:      "Plugin Token Must Have User Key, But X-USER-KEY Is Not Set In Request Header",
        Category:     CategoryAuth,
        Retryable:    false,
        RefreshToken: false,
        MaxRetry:     0,
        UserHint:     "请求缺少 X-USER-KEY，请检查认证头注入逻辑。",
        DevHint:      "这是请求层问题，不是 token 过期问题。",
    },
    10301: {
        Code:         10301,
        Message:      "Check Token Perm Failed",
        Category:     CategoryPermission,
        Retryable:    false,
        RefreshToken: false,
        MaxRetry:     0,
        UserHint:     "权限校验失败，请检查接口权限、版本发布、插件安装状态和空间权限。",
        DevHint:      "不要把 10301 视为 token 过期。",
    },
    10429: {
        Code:         10429,
        Message:      "API Request Frequency Limit",
        Category:     CategoryRateLimit,
        Retryable:    true,
        RefreshToken: false,
        MaxRetry:     3,
        UserHint:     "请求过于频繁，正在稍后重试。",
        DevHint:      "建议使用指数退避。",
    },
    50006: {
        Code:         50006,
        Message:      "openapi system err, please try again later",
        Category:     CategoryServer,
        Retryable:    true,
        RefreshToken: false,
        MaxRetry:     2,
        UserHint:     "服务暂时异常，正在稍后重试。",
        DevHint:      "服务端瞬时错误，可短暂重试。",
    },
}

func LookupPolicy(code int) Policy {
    if p, ok := Policies[code]; ok {
        return p
    }
    return Policy{
        Code:         code,
        Category:     CategoryUnknown,
        Retryable:    false,
        RefreshToken: false,
        MaxRetry:     0,
        UserHint:     "请求失败，请检查错误信息或开启 debug 查看详情。",
        DevHint:      "unknown openapi error code",
    }
}
```

---

## 文件 8：`internal/openapi/errors/normalize.go`

```go
package errors

func NormalizeError(httpStatus int, body []byte, code int, msg string) error {
    if code == 0 && httpStatus >= 200 && httpStatus < 300 {
        return nil
    }

    policy := LookupPolicy(code)

    return &OpenAPIError{
        HTTPStatus: httpStatus,
        Code:       code,
        Message:    msg,
        Category:   policy.Category,
        UserHint:   policy.UserHint,
        DevHint:    policy.DevHint,
        RawBody:    body,
    }
}
```

---

## 文件 9：`internal/openapi/request.go`

```go
package openapi

type Request struct {
    Method  string
    Path    string
    Body    any
    Headers map[string]string
}
```

---

## 文件 10：`internal/openapi/response.go`

```go
package openapi

import "encoding/json"

type CommonResponse struct {
    ErrCode int             `json:"err_code"`
    ErrMsg  string          `json:"err_msg"`
    Data    json.RawMessage `json:"data"`
}
```

---

## 文件 11：`internal/openapi/middleware.go`

```go
package openapi

import (
    "time"

    openapierrors "your/module/internal/openapi/errors"
)

func backoffDuration(category openapierrors.Category, attempt int) time.Duration {
    switch category {
    case openapierrors.CategoryRateLimit:
        switch attempt {
        case 0:
            return 200 * time.Millisecond
        case 1:
            return 500 * time.Millisecond
        default:
            return 1 * time.Second
        }
    case openapierrors.CategoryServer:
        switch attempt {
        case 0:
            return 300 * time.Millisecond
        default:
            return 800 * time.Millisecond
        }
    default:
        return 0
    }
}
```

---

## 文件 12：`internal/openapi/client.go`

```go
package openapi

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"

    "your/module/internal/auth"
    openapierrors "your/module/internal/openapi/errors"
)

type Client struct {
    BaseURL       string
    HTTPClient    *http.Client
    TokenProvider auth.TokenProvider
}

func (c *Client) DoJSON(ctx context.Context, req *Request, out any) error {
    var lastErr error

    for attempt := 0; attempt < 4; attempt++ {
        authCtx, err := c.TokenProvider.GetAuthContext(ctx)
        if err != nil {
            return err
        }

        httpReq, err := c.buildRequest(ctx, req, authCtx)
        if err != nil {
            return err
        }

        httpResp, body, err := c.do(httpReq)
        if err != nil {
            return err
        }

        code, msg, normErr := c.parseAndNormalize(httpResp.StatusCode, body)
        if normErr == nil {
            if out != nil {
                if err := json.Unmarshal(body, out); err != nil {
                    return err
                }
            }
            return nil
        }

        lastErr = normErr
        policy := openapierrors.LookupPolicy(code)

        if policy.RefreshToken && attempt < policy.MaxRetry {
            if _, err := c.TokenProvider.ForceRefresh(ctx); err != nil {
                return err
            }
            continue
        }

        if policy.Retryable && attempt < policy.MaxRetry {
            time.Sleep(backoffDuration(policy.Category, attempt))
            continue
        }

        _ = msg
        return normErr
    }

    return lastErr
}

func (c *Client) buildRequest(ctx context.Context, req *Request, authCtx *auth.AuthContext) (*http.Request, error) {
    var bodyReader io.Reader
    if req.Body != nil {
        payload, err := json.Marshal(req.Body)
        if err != nil {
            return nil, err
        }
        bodyReader = bytes.NewReader(payload)
    }

    fullURL := strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(req.Path, "/")
    httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
    if err != nil {
        return nil, err
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("X-PLUGIN-TOKEN", authCtx.PluginToken)
    httpReq.Header.Set("X-USER-KEY", authCtx.UserKey)

    for k, v := range req.Headers {
        httpReq.Header.Set(k, v)
    }

    return httpReq, nil
}

func (c *Client) do(req *http.Request) (*http.Response, []byte, error) {
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, nil, err
    }

    return resp, body, nil
}

func (c *Client) parseAndNormalize(httpStatus int, body []byte) (int, string, error) {
    var parsed CommonResponse
    if err := json.Unmarshal(body, &parsed); err == nil {
        if parsed.ErrCode != 0 {
            return parsed.ErrCode, parsed.ErrMsg, openapierrors.NormalizeError(httpStatus, body, parsed.ErrCode, parsed.ErrMsg)
        }
        if httpStatus >= 200 && httpStatus < 300 {
            return 0, "", nil
        }
    }

    if httpStatus < 200 || httpStatus >= 300 {
        msg := fmt.Sprintf("http status %d", httpStatus)
        return 0, msg, openapierrors.NormalizeError(httpStatus, body, 0, msg)
    }

    return 0, "", nil
}
```

---

# 三、首轮实现的最小验收清单

建议首轮只先验收以下能力：

## 配置与会话
- [ ] 能读取 `.lark/config.json`
- [ ] 能读取 `.lark/session.json`
- [ ] 删除 `session.json` 后能自动恢复

## token 流程
- [ ] 首次请求自动获取 `plugin_token`
- [ ] token 有效时复用缓存
- [ ] token 即将过期时自动刷新
- [ ] 配置变更后旧 token 自动失效

## 错误处理
- [ ] `10021` 自动刷新并重试
- [ ] `10022` 自动刷新并重试
- [ ] `10211` 自动刷新并重试
- [ ] `10301` 不刷新，只提示权限问题
- [ ] `10429` 执行退避重试
- [ ] `50006` 执行短重试

## 请求层
- [ ] 统一注入 `X-PLUGIN-TOKEN`
- [ ] 统一注入 `X-USER-KEY`
- [ ] 业务命令不再自己处理 token 逻辑

---

# 四、下一步建议

完成本文档后，建议继续补以下两个方向：

1. 增加命令层接入示例：
   - `lark login`
   - `lark plugin set`
   - `lark auth status`
   - `lark doctor`
2. 增加测试清单：
   - config/session 单测
   - token provider 单测
   - error policy 单测
   - request middleware 单测

---

# 五、命令层接入设计

## 5.1 设计原则

1. **零配置优先**：用户首次使用时，引导式交互完成全部配置
2. **渐进式增强**：每个命令都有 `--help`，关键操作支持 `--dry-run` 和 `--debug`
3. **诊断友好**：`lark doctor` 可快速定位问题
4. **错误友好**：错误信息包含用户提示 + 开发者提示两层文案

---

## 5.2 命令清单

### 命令总览

| 命令 | 用途 | 交互模式 |
|------|------|----------|
| `lark login` | 首次登录引导 | 全交互式 |
| `lark plugin set` | 配置插件凭证 | 半交互式（可指定参数） |
| `lark auth status` | 查看当前认证状态 | 非交互式 |
| `lark doctor` | 诊断配置与环境 | 非交互式 + 多级检查 |

---

## 5.3 `lark login` - 首次登录引导

### 使用场景
- 首次使用 CLI
- 用户想重置全部配置

### 交互流程

```
$ lark login

  ██████╗ ██████╗ ███╗   ██╗███╗   ██╗███████╗ ██████╗████████╗███████╗██████╗
  ██╔════╝██╔═══██╗████╗  ██║████╗  ██║██╔════╝██╔════╝╚══██╔══╝██╔════╝██╔══██╗
  ██║     ██║   ██║██╔██╗ ██║██╔██╗ ██║█████╗  ██║        ██║   █████╗  ██║  ██║
  ██║     ██║   ██║██║╚██╗██║██║╚██╗██║██╔══╝  ██║        ██║   ██╔══╝  ██║  ██║
  ╚██████╗╚██████╔╝██║ ╚████║██║ ╚████║███████╗╚██████╗   ██║   ███████╗██████╔╝
   ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝ ╚═════╝   ╚═╝   ╚══════╝╚═════╝

  飞书 Open API CLI 工具
  首次使用需要配置插件凭证

  开始配置...

  ? 请输入 User Key: › _
```

### 详细交互步骤

```
Step 1: 引导用户输入 User Key
  ? 请输入 User Key: › <输入user_key>
  ✓ User Key 已记录

Step 2: 引导用户输入 Plugin ID
  ? 请输入 Plugin ID (Plugin_ID): › <输入plugin_id>
  ✓ Plugin ID 已记录

Step 3: 引导用户输入 Plugin Secret
  ? 请输入 Plugin Secret (Plugin_Secret): › <输入secret>
  ✓ Plugin Secret 已记录（已脱敏显示）

Step 4: 确认 Base URL（可选，有默认值）
  ? 请确认 Base URL [https://project.feishu.cn/open_api]: ›
  ✓ 使用默认值

Step 5: 尝试验证凭证
  正在验证凭证...
  ✓ 凭证验证成功

Step 6: 保存配置
  正在保存配置到 ~/.lark/config.json...
  ✓ 配置保存成功

Step 7: 获取初始 token
  正在获取 plugin_token...
  ✓ 已获取 plugin_token
  ✓ 已保存 session 到 ~/.lark/session.json

完成！您现在可以使用 lark CLI 访问飞书 Open API。

运行 lark auth status 查看当前认证状态
运行 lark doctor 检查配置是否正常
```

### 参数支持

```bash
# 全程静默模式（用于脚本）
lark login --non-interactive \
  --user-key <key> \
  --plugin-id <id> \
  --plugin-secret <secret>

# 指定自定义 Base URL
lark login --base-url https://custom.feishu.cn/open_api

# 输出详细调试信息
lark login --debug

# 跳过凭证验证（不推荐）
lark login --skip-verify
```

### 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 登录成功 |
| 1 | 登录失败（凭证错误、网络问题等） |
| 2 | 参数错误 |

---

## 5.4 `lark plugin set` - 配置插件凭证

### 使用场景
- 已登录用户想更换插件凭证
-只想更新部分配置（如换了新插件）

### 交互流程

```
$ lark plugin set

  配置当前插件凭证

  ? 请选择要配置的项目: (Use arrow keys)
    ❯ User Key
      Plugin ID
      Plugin Secret
      Base URL
      全部重新配置
      取消
```

### 详细交互步骤

#### 场景 1：单独修改某一项

```
$ lark plugin set

? 请选择要配置的项目: User Key
? 请输入新的 User Key: <输入新key>
✓ User Key 已更新

? 是否立即验证新凭证? (Y/n): Y
正在验证凭证...
✓ 凭证验证成功
✓ Token 已刷新
```

#### 场景 2：全部重新配置

```
$ lark plugin set --all

? 请输入 User Key: <输入>
? 请输入 Plugin ID: <输入>
? 请输入 Plugin Secret: <输入>
? 请确认 Base URL [https://project.feishu.cn/open_api]: <确认或修改>

正在验证凭证...
✓ 凭证验证成功
✓ 配置已更新
✓ Token 已刷新
```

### 参数支持

```bash
# 单独设置某个字段
lark plugin set --user-key <key>
lark plugin set --plugin-id <id>
lark plugin set --plugin-secret <secret>
lark plugin set --base-url <url>

# 组合设置（半交互式）
lark plugin set --plugin-id <id> --plugin-secret <secret>

# 非交互式（用于脚本）
lark plugin set --non-interactive \
  --user-key <key> \
  --plugin-id <id> \
  --plugin-secret <secret>

# 验证后保存（默认）
lark plugin set --verify

# 跳过验证（不推荐）
lark plugin set --no-verify
```

### 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 配置成功 |
| 1 | 配置失败 |
| 2 | 参数错误 |
| 3 | 验证失败（凭证无效） |

---

## 5.5 `lark auth status` - 查看认证状态

### 使用场景
- 用户想确认当前登录状态
- 调试时确认 token 是否有效

### 输出示例

```
$ lark auth status

  认证状态

  配置信息:
    User Key:     xxxxxx...（已脱敏）
    Plugin ID:    cli_xxxxxxxxx
    Base URL:     https://project.feishu.cn/open_api

  Token 状态:
    状态:         ✓ 有效
    获取时间:     2026-03-19 10:30:00
    过期时间:     2026-03-19 11:30:00
    剩余有效期:   58 分钟
    缓存来源:     ~/.lark/session.json

  配置指纹:       a1b2c3d4e5f6...（已脱敏）
```

### Token 过期时的输出

```
$ lark auth status

  认证状态

  配置信息:
    User Key:     xxxxxx...（已脱敏）
    Plugin ID:    cli_xxxxxxxxx
    Base URL:     https://project.feishu.cn/open_api

  Token 状态:
    状态:         ⚠ 已过期
    获取时间:     2026-03-19 09:30:00
    过期时间:     2026-03-19 10:30:00
    已过期:       30 分钟

  建议: 运行 lark auth refresh 刷新 Token
```

### Token 与配置不匹配时

```
$ lark auth status

  认证状态

  ⚠ 检测到配置已变更，缓存的 Token 已失效

  配置指纹:       a1b2c3d4e5f6...（当前）
  Token 指纹:     x1y2z3a4b5c6...（缓存）

  建议: 运行 lark auth refresh 重新获取 Token
```

### 参数支持

```bash
# 输出 JSON 格式（便于脚本解析）
lark auth status --json

# 仅显示状态码（0=正常，1=异常）
lark auth status --quiet

# 显示完整信息（不过敏）
lark auth status --show-full
```

### JSON 输出格式

```json
{
  "config": {
    "user_key": "xxxxxx***",
    "plugin_id": "cli_xxx",
    "base_url": "https://..."
  },
  "token": {
    "status": "valid",  // valid | expired | invalid | missing
    "obtained_at": "2026-03-19T10:30:00Z",
    "expire_at": "2026-03-19T11:30:00Z",
    "remaining_minutes": 58,
    "fingerprint": "a1b2c3..."
  },
  "match": true  // config fingerprint == token fingerprint
}
```

---

## 5.6 `lark auth refresh` - 强制刷新 Token

### 使用场景
- Token 过期但想手动刷新
- 配置变更后想立即获取新 Token

### 使用方式

```bash
# 刷新 Token
lark auth refresh

# 刷新并显示详情
lark auth refresh --verbose

# 强制刷新（即使现有 Token 有效）
lark auth refresh --force
```

### 输出示例

```
$ lark auth refresh

  正在刷新 Token...
  ✓ 已获取新 Token
  ✓ 已保存到 ~/.lark/session.json
  剩余有效期: 60 分钟
```

---

## 5.7 `lark doctor` - 诊断工具

### 使用场景
- 用户遇到问题时的第一排查步骤
- 提交 Issue 前的自助诊断

### 设计理念

`lark doctor` 执行多级检查，每项检查有独立的通过/警告/失败状态，最终给出汇总报告。

### 检查项清单

| 序号 | 检查项 | 说明 | 状态码 |
|------|--------|------|--------|
| 1 | 配置文件存在性 | 检查 `.lark/config.json` 是否存在 | pass/warn/fail |
| 2 | Session 存在性 | 检查 `.lark/session.json` 是否存在 | pass/warn |
| 3 | 配置完整性 | 检查 User Key、Plugin ID、Plugin Secret 是否齐全 | pass/fail |
| 4 | 配置格式 | JSON 格式是否合法 | pass/fail |
| 5 | Token 有效性 | 检查缓存的 token 是否未过期 | pass/warn/fail |
| 6 | 配置指纹匹配 | 检查 config 和 session 的指纹是否一致 | pass/fail |
| 7 | 网络连通性 | 测试是否能访问飞书 API | pass/fail |
| 8 | 凭证有效性 | 实际调用 API 验证凭证 | pass/fail |

### 输出示例

```
$ lark doctor

  ██╗  ██╗██╗   ██╗███╗  ██╗████████╗███████╗██████╗ ███╗   ███╗██╗███╗  ██╗██╗   ██╗███████╗
  ██║  ██║██║   ██║████╗ ██║╚══██╔══╝██╔════╝██╔══██╗████╗ ████║██║████╗ ██║██║   ██║██╔════╝
  ███████║██║   ██║██╔██╗██║   ██║   █████╗  ██████╔╝██╔████╔██║██║██╔██╗██║██║   ██║███████╗
  ██╔══██║██║   ██║██║╚████║   ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║██║██║╚████║██║   ██║╚════██║
  ██║  ██║╚██████╔╝██║ ╚███║   ██║   ███████╗██║  ██║██║ ╚═╝ ██║██║██║ ╚███║╚██████╔╝███████║
  ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝╚═╝  ╚══╝ ╚═════╝ ╚══════╝

  配置诊断工具 v1.0

  [1/8] 检查配置文件...
        ~/.lark/config.json          ✓ 存在

  [2/8] 检查 Session 文件...
        ~/.lark/session.json         ✓ 存在

  [3/8] 检查配置完整性...
        User Key                     ✓ 已配置
        Plugin ID                    ✓ 已配置
        Plugin Secret                ✓ 已配置
        Base URL                     ✓ 已配置 (https://project.feishu.cn/open_api)

  [4/8] 检查配置格式...
                                ✓ JSON 格式正确

  [5/8] 检查 Token 有效性...
                                ✓ Token 未过期 (剩余 45 分钟)

  [6/8] 检查配置指纹匹配...
                                ✓ 配置与 Token 匹配

  [7/8] 检查网络连通性...
        project.feishu.cn            ✓ 可访问 (23ms)

  [8/8] 验证凭证有效性...
                                ✓ 凭证有效，Token 获取成功

  ─────────────────────────────────────────────────

  诊断结果: ✓ 全部检查通过

  您的 CLI 配置正常，可以正常使用。
  如遇问题，可运行 lark auth status 查看详细状态。
```

### 检查失败时的输出

```
$ lark doctor

  ... (前几步检查同上) ...

  [5/8] 检查 Token 有效性...
                                ⚠ Token 已过期 (过期 15 分钟前)

  [6/8] 检查配置指纹匹配...
                                ✓ 配置与 Token 匹配

  ... (后续检查) ...

  ─────────────────────────────────────────────────

  诊断结果: ⚠ 1 个警告

  检测到以下问题:

    [5/8] Token 已过期
      建议: 运行 lark auth refresh 刷新 Token

  ─────────────────────────────────────────────────

  运行 lark auth refresh 解决上述问题。
  如问题持续，请运行 lark doctor --verbose 获取详细信息。
```

### 严重错误时的输出

```
$ lark doctor

  ... (前几步检查) ...

  [3/8] 检查配置完整性...
        User Key                     ✗ 未配置
        Plugin ID                    ✗ 未配置
        Plugin Secret                ✗ 未配置

  ─────────────────────────────────────────────────

  诊断结果: ✗ 配置缺失

  检测到以下错误:

    [3/8] 缺少必需配置
      请运行 lark login 进行首次配置

  ─────────────────────────────────────────────────

  运行 lark login 开始配置。
```

### 参数支持

```bash
# 静默模式（只返回退出码）
lark doctor --quiet

# 详细输出
lark doctor --verbose

# 输出 JSON 格式
lark doctor --json

# 只执行特定检查
lark doctor --check token      # 只检查 token
lark doctor --check network    # 只检查网络
lark doctor --check config     # 只检查配置

# 不执行实际 API 调用（跳过验证）
lark doctor --no-verify
```

### 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 全部检查通过 |
| 1 | 有警告（不影响使用） |
| 2 | 有错误（无法正常使用） |
| 3 | 参数错误 |

---

## 5.8 命令行框架建议

### 推荐库

| 库 | 特点 | 推荐场景 |
|----|------|----------|
| cobra | 最流行，功能完整 | 生产级 CLI |
| urfave/cli | 轻量，API 简洁 | 中小型 CLI |
| kingpin | 声明式，类型安全 | Go 1.17+ 项目 |

**推荐**：使用 **cobra** + **survey** 组合
- `cobra`：命令注册、参数解析、帮助文档
- `survey`：交互式输入

### 目录结构建议

```
cmd/
  root.go           # Root command
  login.go          # lark login
  plugin/
    set.go          # lark plugin set
  auth/
    status.go       # lark auth status
    refresh.go      # lark auth refresh
  doctor.go         # lark doctor
```

### Root Command 骨架

```go
// cmd/root.go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    debug bool
)

var rootCmd = &cobra.Command{
    Use:   "lark",
    Short: "飞书 Open API CLI 工具",
    Long: `lark CLI - 飞书 Open API 命令行工具

常用命令:
  lark login         首次登录引导
  lark plugin set    配置插件凭证
  lark auth status   查看认证状态
  lark doctor        诊断配置问题

更多帮助请使用 --help`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        if debug {
            fmt.Println("[DEBUG] Running command:", cmd.Name())
        }
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "输出调试信息")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "输出详细信息")
}
```

### 错误输出规范

```go
// 所有错误统一使用以下格式

// 用户可见错误（友好提示）
fmt.Fprintf(os.Stderr, "✗ %s\n", userHint)
if debug || verbose {
    fmt.Fprintf(os.Stderr, "  [DEV] %s\n", devHint)
}

// 致命错误（程序无法继续）
fmt.Fprintf(os.Stderr, "✗ %s: %v\n", userHint, err)
os.Exit(1)
```

---

## 5.9 命令间依赖关系

```
                    ┌─────────────┐
                    │ lark login  │
                    └──────┬──────┘
                           │ 创建 config.json
                           ▼
              ┌────────────────────────┐
              │   lark plugin set     │◄────── 更新配置
              └───────────┬────────────┘
                          │ 生成/更新 session.json
                          ▼
              ┌─────────────────────────┐
              │   lark auth refresh     │◄────── 手动刷新
              └───────────┬─────────────┘
                          │
         ┌────────────────┼────────────────┐
         ▼                ▼                ▼
  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
  │ lark auth   │  │ lark doctor │  │ 业务命令    │
  │ status      │  │             │  │ (user query │
  └─────────────┘  └─────────────┘  │  等)        │
                                    └─────────────┘
```

---

## 5.10 实现优先级

| 优先级 | 命令 | 理由 |
|--------|------|------|
| P0 | `lark doctor` | 诊断友好，用户遇到问题时先跑这个 |
| P0 | `lark auth status` | 最基础的检查命令 |
| P1 | `lark auth refresh` | Token 过期的补救措施 |
| P1 | `lark plugin set` | 配置更新的核心命令 |
| P2 | `lark login` | 首次引导，可与 plugin set 复用逻辑 |

---

## 5.11 与 OpenAPI Client 的集成

每个命令在执行时都通过 `openapi.Client` 访问 API：

```go
// cmd/auth/status.go
func runAuthStatus(cmd *cobra.Command, args []string) error {
    // 1. 创建 TokenProvider
    cfgStore := &config.FileStore{Path: configPath}
    sessStore := &session.FileStore{Path: sessionPath}
    authAPI := &auth.HTTPAuthAPI{BaseURL: baseURL, HTTPClient: http.DefaultClient}
    tokenProvider := &auth.CachedTokenProvider{
        ConfigStore:  cfgStore,
        SessionStore: sessStore,
        AuthAPI:      authAPI,
    }

    // 2. 获取认证上下文
    authCtx, err := tokenProvider.GetAuthContext(ctx)
    if err != nil {
        return fmt.Errorf("获取认证上下文失败: %w", err)
    }

    // 3. 显示状态
    printAuthStatus(authCtx, cfg, sess)

    return nil
}
```

---

# 六、测试清单与骨架

## 6.1 测试分层与工具选型

| 层级 | 测试内容 | 工具选型 | 覆盖目标 |
|------|----------|----------|----------|
| **单元测试** | config/session 读写、fingerprint、错误策略 | `testing` + `testify/assert` | 每个 public 方法 100% 覆盖 |
| **集成测试** | token 获取/刷新/缓存全链路 | `net/http/httptest` mock server | 各模块协作正确性 |
| **端到端测试** | CLI 命令执行 | `cobra` + `os/exec` | 用户实际使用场景 |

### 测试依赖注入策略

为便于单元测试，所有依赖均通过 interface 注入：

```go
// 测试友好的设计模式
type FileStore struct {
    Path    string
    readFile  func(path string) ([]byte, error)
    writeFile func(path string, data []byte, perm os.FileMode) error
}
```

---

## 6.2 config 模块测试

### 测试文件：`internal/config/config_test.go`

#### 6.2.1 `TestConfig_ValidateForOpenAPI` - 配置校验

| 编号 | 用例名称 | 输入 | 期望输出 |
|------|----------|------|----------|
| TC-CONFIG-001 | 全部字段有效 | 完整 Config | `nil` |
| TC-CONFIG-002 | Config 为 nil | `nil` | error: "config is nil" |
| TC-CONFIG-003 | UserKey 为空 | 缺 UserKey | error: "missing user_key" |
| TC-CONFIG-004 | PluginID 为空 | 缺 PluginID | error: "missing plugin_id" |
| TC-CONFIG-005 | PluginSecret 为空 | 缺 PluginSecret | error: "missing plugin_secret" |
| TC-CONFIG-006 | BaseURL 为空 | 缺 BaseURL | error: "missing base_url" |
| TC-CONFIG-007 | 多个字段缺失 | 缺 UserKey + PluginID | error 包含 "user_key" |

```go
func TestConfig_ValidateForOpenAPI(t *testing.T) {
    tests := []struct {
        name    string
        cfg     *Config
        wantErr bool
        errMsg  string
    }{
        {
            name:    "TC-CONFIG-001: valid config",
            cfg:     &Config{UserKey: "k", PluginID: "i", PluginSecret: "s", BaseURL: "u"},
            wantErr: false,
        },
        {
            name:    "TC-CONFIG-002: config is nil",
            cfg:     nil,
            wantErr: true,
            errMsg:  "config is nil",
        },
        {
            name:    "TC-CONFIG-003: missing user_key",
            cfg:     &Config{PluginID: "i", PluginSecret: "s", BaseURL: "u"},
            wantErr: true,
            errMsg:  "missing user_key",
        },
        // ... TC-CONFIG-004 ~ TC-CONFIG-007
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.cfg.ValidateForOpenAPI()
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### 6.2.2 `TestFileStore_Load` - 配置加载

| 编号 | 用例名称 | 场景 | 期望输出 |
|------|----------|------|----------|
| TC-CONFIG-010 | 文件存在正常加载 | 完整 JSON | Config.UserKey == "test" |
| TC-CONFIG-011 | BaseURL 默认值填充 | 缺 BaseURL | BaseURL == 默认值 |
| TC-CONFIG-012 | 文件不存在返回默认 | 文件不存在 | 返回 DefaultConfig() |
| TC-CONFIG-013 | JSON 格式错误 | 非法 JSON | error |
| TC-CONFIG-014 | 部分字段覆盖 | 缺 PluginID | PluginID == "" |

```go
func TestFileStore_Load(t *testing.T) {
    defaultBaseURL := "https://project.feishu.cn/open_api"

    tests := []struct {
        name          string
        fileContent   string
        wantUserKey   string
        wantBaseURL   string
        wantErr       bool
    }{
        {
            name:        "TC-CONFIG-010: file exists and valid",
            fileContent: `{"user_key":"test_key","plugin_id":"id","plugin_secret":"sec","base_url":"https://custom"}`,
            wantUserKey: "test_key",
            wantBaseURL: "https://custom",
            wantErr:     false,
        },
        {
            name:        "TC-CONFIG-011: missing base_url uses default",
            fileContent: `{"user_key":"test_key","plugin_id":"id","plugin_secret":"sec"}`,
            wantUserKey: "test_key",
            wantBaseURL: defaultBaseURL,
            wantErr:     false,
        },
        {
            name:    "TC-CONFIG-013: invalid json",
            fileContent: `{invalid`,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 创建临时文件
            tmp, err := os.CreateTemp("", "config-*.json")
            assert.NoError(t, err)
            defer os.Remove(tmp.Name())

            if tt.fileContent != "NO_FILE" {
                _, err = tmp.WriteString(tt.fileContent)
                assert.NoError(t, err)
                tmp.Close()
            }

            store := &FileStore{Path: tmp.Name()}
            cfg, err := store.Load(context.Background())

            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.wantUserKey, cfg.UserKey)
            assert.Equal(t, tt.wantBaseURL, cfg.BaseURL)
        })
    }
}
```

#### 6.2.3 `TestFileStore_Save` - 配置保存

| 编号 | 用例名称 | 场景 | 期望输出 |
|------|----------|------|----------|
| TC-CONFIG-020 | 正常保存 | 有效 Config | 文件存在且可重新加载 |
| TC-CONFIG-021 | 父目录不存在自动创建 | 嵌套路径 | 目录和文件都创建成功 |
| TC-CONFIG-022 | Config 为 nil | `nil` | error |
| TC-CONFIG-023 | 原子写入一致性 | 并发写入 | 文件不损坏 |

```go
func TestFileStore_Save(t *testing.T) {
    t.Run("TC-CONFIG-020: normal save and reload", func(t *testing.T) {
        tmp, _ := os.CreateTemp("", "config-save-*.json")
        defer os.Remove(tmp.Name())
        tmp.Close()

        cfg := &Config{
            UserKey:      "key",
            PluginID:     "id",
            PluginSecret: "secret",
            BaseURL:      "url",
        }

        store := &FileStore{Path: tmp.Name()}
        err := store.Save(context.Background(), cfg)
        assert.NoError(t, err)

        // 验证可以重新加载
        loaded, err := store.Load(context.Background())
        assert.NoError(t, err)
        assert.Equal(t, cfg.UserKey, loaded.UserKey)
    })

    t.Run("TC-CONFIG-022: nil config", func(t *testing.T) {
        store := &FileStore{Path: "/tmp/test.json"}
        err := store.Save(context.Background(), nil)
        assert.Error(t, err)
    })
}
```

---

## 6.3 session 模块测试

### 测试文件：`internal/session/session_test.go`

#### 6.3.1 `TestSession_IsEmpty` - 空判断

| 编号 | 用例名称 | 输入 | 期望输出 |
|------|----------|------|----------|
| TC-SESSION-001 | Session 为 nil | `nil` | `true` |
| TC-SESSION-002 | PluginToken 为空 | `&Session{}` | `true` |
| TC-SESSION-003 | PluginToken 非空 | Token 有值 | `false` |

```go
func TestSession_IsEmpty(t *testing.T) {
    tests := []struct {
        name  string
        sess  *Session
        want  bool
    }{
        {"TC-SESSION-001: nil session", nil, true},
        {"TC-SESSION-002: empty token", &Session{}, true},
        {"TC-SESSION-003: has token", &Session{PluginToken: "token"}, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.sess.IsEmpty()
            assert.Equal(t, tt.want, got)
        })
    }
}
```

#### 6.3.2 `TestSession_IsValid` - 有效性判断

| 编号 | 用例名称 | 场景 | 期望输出 |
|------|----------|------|----------|
| TC-SESSION-010 | 全部有效 | 未过期+指纹匹配 | `true` |
| TC-SESSION-011 | Session 为 nil | `nil` | `false` |
| TC-SESSION-012 | Token 为空 | 空 Token | `false` |
| TC-SESSION-013 | 已过期 | ExpireAt <= now | `false` |
| TC-SESSION-014 | 指纹不匹配 | Config 变更 | `false` |
| TC-SESSION-015 | 刚好过期 | ExpireAt == now | `false` |

```go
func TestSession_IsValid(t *testing.T) {
    validFingerprint := "fp123"
    now := int64(1000000)

    tests := []struct {
        name        string
        sess        *Session
        now         int64
        fingerprint string
        want        bool
    }{
        {
            name:        "TC-SESSION-010: valid session",
            sess:        &Session{PluginToken: "tok", ExpireAt: now + 3600, ConfigFingerprint: validFingerprint},
            now:         now,
            fingerprint: validFingerprint,
            want:        true,
        },
        {
            name:        "TC-SESSION-011: nil session",
            sess:        nil,
            now:         now,
            fingerprint: validFingerprint,
            want:        false,
        },
        {
            name:        "TC-SESSION-012: empty token",
            sess:        &Session{PluginToken: "", ExpireAt: now + 3600, ConfigFingerprint: validFingerprint},
            now:         now,
            fingerprint: validFingerprint,
            want:        false,
        },
        {
            name:        "TC-SESSION-013: expired",
            sess:        &Session{PluginToken: "tok", ExpireAt: now - 1, ConfigFingerprint: validFingerprint},
            now:         now,
            fingerprint: validFingerprint,
            want:        false,
        },
        {
            name:        "TC-SESSION-014: fingerprint mismatch",
            sess:        &Session{PluginToken: "tok", ExpireAt: now + 3600, ConfigFingerprint: "old_fp"},
            now:         now,
            fingerprint: validFingerprint,
            want:        false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.sess.IsValid(tt.now, tt.fingerprint)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

#### 6.3.3 `TestFileStore_Load` / `Save` / `Clear`

| 编号 | 用例名称 | 场景 |
|------|----------|------|
| TC-SESSION-020 | 加载有效 session | JSON → Session 结构 |
| TC-SESSION-021 | 文件不存在返回空 Session | 新建 `&Session{}` |
| TC-SESSION-022 | 保存后重新加载 | 数据一致性 |
| TC-SESSION-023 | Clear 删除文件 | 文件不存在 |

---

## 6.4 auth/fingerprint 模块测试

### 测试文件：`internal/auth/fingerprint_test.go`

#### 6.4.1 `TestBuildConfigFingerprint`

| 编号 | 用例名称 | 输入 | 期望输出 |
|------|----------|------|----------|
| TC-FP-001 | nil Config | `nil` | `""` |
| TC-FP-002 | 相同配置生成相同指纹 | 两份相同 Config | 指纹相同 |
| TC-FP-003 | 不同配置生成不同指纹 | Config 任一字段不同 | 指纹不同 |
| TC-FP-004 | 指纹长度固定 | 任意 Config | SHA256 长度 (64 字符) |
| TC-FP-005 | 顺序无关 | 字段值相同顺序不同 | 指纹相同 |

```go
func TestBuildConfigFingerprint(t *testing.T) {
    cfg := &config.Config{
        UserKey:      "key1",
        PluginID:     "id1",
        PluginSecret: "secret1",
        BaseURL:      "url1",
    }

    t.Run("TC-FP-001: nil config returns empty", func(t *testing.T) {
        fp := BuildConfigFingerprint(nil)
        assert.Equal(t, "", fp)
    })

    t.Run("TC-FP-002: same config produces same fingerprint", func(t *testing.T) {
        fp1 := BuildConfigFingerprint(cfg)
        fp2 := BuildConfigFingerprint(cfg)
        assert.Equal(t, fp1, fp2)
    })

    t.Run("TC-FP-003: different config produces different fingerprint", func(t *testing.T) {
        fp1 := BuildConfigFingerprint(cfg)
        fp2 := BuildConfigFingerprint(&config.Config{
            UserKey:      "key2", // changed
            PluginID:     "id1",
            PluginSecret: "secret1",
            BaseURL:      "url1",
        })
        assert.NotEqual(t, fp1, fp2)
    })

    t.Run("TC-FP-004: fingerprint length is 64 (SHA256)", func(t *testing.T) {
        fp := BuildConfigFingerprint(cfg)
        assert.Len(t, fp, 64)
    })
}
```

---

## 6.5 auth/auth_api 模块测试

### 测试文件：`internal/auth/auth_api_test.go`

#### 6.5.1 Mock Server 设置

```go
// httptest.NewServer mock
func setupMockAuthServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 验证请求路径和内容
        assert.Equal(t, "/authen/plugin_token", r.URL.Path)
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

        handler(w, r)
    }))
}
```

#### 6.5.2 `TestHTTPAuthAPI_FetchPluginToken`

| 编号 | 用例名称 | Mock 响应 | 期望输出 |
|------|----------|-----------|----------|
| TC-AUTH-001 | 成功获取 Token | code=0, 返回 token | Token="tok", Exp=3600 |
| TC-AUTH-002 | HTTP 500 | status=500 | error |
| TC-AUTH-003 | 业务错误码 | code=99999 | error |
| TC-AUTH-004 | 无响应体 | 空 body | error |
| TC-AUTH-005 | 非法 JSON | 非 JSON | error |

```go
func TestHTTPAuthAPI_FetchPluginToken(t *testing.T) {
    t.Run("TC-AUTH-001: successful token fetch", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"error":{"code":0,"msg":""},"data":{"token":"test_token","expire_time":3600}}`))
        }))
        defer server.Close()

        api := &HTTPAuthAPI{BaseURL: server.URL, HTTPClient: server.Client()}
        resp, err := api.FetchPluginToken(context.Background(), "plugin_id", "plugin_secret")

        assert.NoError(t, err)
        assert.Equal(t, "test_token", resp.Token)
        assert.Equal(t, int64(3600), resp.ExpireTime)
    })

    t.Run("TC-AUTH-003: business error", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"error":{"code":99999,"msg":"invalid plugin"},"data":{}}`))
        }))
        defer server.Close()

        api := &HTTPAuthAPI{BaseURL: server.URL, HTTPClient: server.Client()}
        _, err := api.FetchPluginToken(context.Background(), "plugin_id", "plugin_secret")

        assert.Error(t, err)
        assert.Contains(t, err.Error(), "99999")
    })
}
```

---

## 6.6 auth/token_provider 模块测试

### 测试文件：`internal/auth/token_provider_test.go`

#### 6.6.1 Mock Store 实现

```go
type mockConfigStore struct {
    cfg *config.Config
    err error
}

func (m *mockConfigStore) Load(ctx context.Context) (*config.Config, error) {
    return m.cfg, m.err
}

type mockSessionStore struct {
    sess *session.Session
    err  error
    saveCalls int
}

func (m *mockSessionStore) Load(ctx context.Context) (*session.Session, error) {
    return m.sess, m.err
}

func (m *mockSessionStore) Save(ctx context.Context, sess *session.Session) error {
    m.saveCalls++
    m.sess = sess
    return m.err
}

func (m *mockSessionStore) Clear(ctx context.Context) error {
    m.sess = nil
    return nil
}

type mockAuthAPI struct {
    resp *PluginTokenResponse
    err  error
}

func (m *mockAuthAPI) FetchPluginToken(ctx context.Context, pluginID, pluginSecret string) (*PluginTokenResponse, error) {
    return m.resp, m.err
}
```

#### 6.6.2 `TestCachedTokenProvider_GetAuthContext`

| 编号 | 用例名称 | 场景 | 期望输出 |
|------|----------|------|----------|
| TC-TP-001 | 缓存有效直接返回 | Session 未过期+指纹匹配 | 返回缓存 Token |
| TC-TP-002 | 缓存过期自动刷新 | Session 已过期 | 调用 AuthAPI，保存新 Session |
| TC-TP-003 | 缓存指纹不匹配 | Config 变更 | 调用 AuthAPI，获取新 Token |
| TC-TP-004 | 缓存不存在 | Session 为空 | 调用 AuthAPI |
| TC-TP-005 | 缓存不存在但 API 失败 | API 返回 error | error |
| TC-TP-006 | 并发安全 | 多次并发调用 | 只调用一次 AuthAPI |

```go
func TestCachedTokenProvider_GetAuthContext(t *testing.T) {
    validCfg := &config.Config{
        UserKey:      "key", PluginID: "id", PluginSecret: "sec", BaseURL: "url",
    }
    validFP := "valid_fp"
    now := time.Now().Unix()

    t.Run("TC-TP-001: valid cache returns cached token", func(t *testing.T) {
        cfgStore := &mockConfigStore{cfg: validCfg}
        sessStore := &mockSessionStore{
            sess: &session.Session{
                PluginToken:       "cached_token",
                ExpireAt:          now + 3600,
                ConfigFingerprint: validFP,
            },
        }
        authAPI := &mockAuthAPI{}

        provider := NewCachedTokenProvider(cfgStore, sessStore, authAPI, 5*time.Minute)
        ctx, err := provider.GetAuthContext(context.Background())

        assert.NoError(t, err)
        assert.Equal(t, "cached_token", ctx.PluginToken)
        assert.Equal(t, 0, authAPI.callCount) // 不应调用 API
    })

    t.Run("TC-TP-002: expired cache triggers refresh", func(t *testing.T) {
        cfgStore := &mockConfigStore{cfg: validCfg}
        sessStore := &mockSessionStore{
            sess: &session.Session{
                PluginToken:       "old_token",
                ExpireAt:          now - 1, // 已过期
                ConfigFingerprint: validFP,
            },
        }
        authAPI := &mockAuthAPI{
            resp: &PluginTokenResponse{Token: "new_token", ExpireTime: 3600},
        }

        provider := NewCachedTokenProvider(cfgStore, sessStore, authAPI, 5*time.Minute)
        ctx, err := provider.GetAuthContext(context.Background())

        assert.NoError(t, err)
        assert.Equal(t, "new_token", ctx.PluginToken)
        assert.Equal(t, 1, authAPI.callCount)
        assert.Equal(t, 1, sessStore.saveCalls) // 保存新 Session
    })

    t.Run("TC-TP-006: concurrent calls only refresh once", func(t *testing.T) {
        cfgStore := &mockConfigStore{cfg: validCfg}
        sessStore := &mockSessionStore{
            sess: &session.Session{PluginToken: "", ExpireAt: 0}, // 空缓存
        }
        authAPI := &mockAuthAPI{
            resp: &PluginTokenResponse{Token: "token", ExpireTime: 3600},
        }

        provider := NewCachedTokenProvider(cfgStore, sessStore, authAPI, 5*time.Minute)

        var wg sync.WaitGroup
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                provider.GetAuthContext(context.Background())
            }()
        }
        wg.Wait()

        assert.Equal(t, 1, authAPI.callCount) // 只调用一次
    })
}
```

---

## 6.7 openapi/errors 模块测试

### 测试文件：`internal/openapi/errors/errors_test.go`

#### 6.7.1 `TestLookupPolicy`

| 编号 | 用例名称 | 输入 | 期望输出 |
|------|----------|------|----------|
| TC-ERR-001 | 已知错误码 10021 | 10021 | CategoryAuth, RefreshToken=true |
| TC-ERR-002 | 已知错误码 10301 | 10301 | CategoryPermission, RefreshToken=false |
| TC-ERR-003 | 已知错误码 10429 | 10429 | CategoryRateLimit, Retryable=true |
| TC-ERR-004 | 未知错误码 | 99999 | CategoryUnknown, Retryable=false |

```go
func TestLookupPolicy(t *testing.T) {
    tests := []struct {
        name           string
        code           int
        wantCategory   Category
        wantRefresh    bool
        wantRetryable  bool
    }{
        {"TC-ERR-001: token not exist", 10021, CategoryAuth, true, true},
        {"TC-ERR-002: permission denied", 10301, CategoryPermission, false, false},
        {"TC-ERR-003: rate limit", 10429, CategoryRateLimit, false, true},
        {"TC-ERR-004: unknown code", 99999, CategoryUnknown, false, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := LookupPolicy(tt.code)
            assert.Equal(t, tt.wantCategory, p.Category)
            assert.Equal(t, tt.wantRefresh, p.RefreshToken)
            assert.Equal(t, tt.wantRetryable, p.Retryable)
        })
    }
}
```

#### 6.7.2 `TestNormalizeError`

| 编号 | 用例名称 | 输入 | 期望输出 |
|------|----------|------|----------|
| TC-ERR-010 | 成功响应 | code=0, status=200 | `nil` |
| TC-ERR-011 | 业务错误 | code=10022, status=200 | OpenAPIError, UserHint 非空 |
| TC-ERR-012 | HTTP 错误 | status=500, code=0 | OpenAPIError, HTTPStatus=500 |
| TC-ERR-013 | 未知错误码 | code=99999 | Category=Unknown |

```go
func TestNormalizeError(t *testing.T) {
    t.Run("TC-ERR-010: success returns nil", func(t *testing.T) {
        err := NormalizeError(200, []byte(`{"err_code":0}`), 0, "")
        assert.Nil(t, err)
    })

    t.Run("TC-ERR-011: token expired normalizes correctly", func(t *testing.T) {
        err := NormalizeError(200, []byte(`{"err_code":10022}`), 10022, "Check Token Failed")
        assert.NotNil(t, err)

        apiErr, ok := err.(*OpenAPIError)
        assert.True(t, ok)
        assert.Equal(t, CategoryAuth, apiErr.Category)
        assert.Equal(t, 10022, apiErr.Code)
        assert.NotEmpty(t, apiErr.UserHint)
    })
}
```

---

## 6.8 openapi/client 模块测试

### 测试文件：`internal/openapi/client_test.go`

#### 6.8.1 `TestClient_DoJSON` - 完整重试流程

| 编号 | 用例名称 | Mock 行为 | 期望调用次数 |
|------|----------|-----------|--------------|
| TC-CLIENT-001 | 首次成功 | 返回 success | 1 次 |
| TC-CLIENT-002 | Token 过期重试 | 1st: 10022, 2nd: success | 2 次 API + 1 次 token refresh |
| TC-CLIENT-003 | 限流重试 | 1st: 10429, 2nd: success | 2 次 |
| TC-CLIENT-004 | 权限错误不重试 | 1st: 10301 | 1 次 |
| TC-CLIENT-005 | 全部重试失败 | 连续 3 次 10429 | 3 次后 error |

```go
func TestClient_DoJSON(t *testing.T) {
    t.Run("TC-CLIENT-001: first call succeeds", func(t *testing.T) {
        server := mockAuthServer(t, func(count *int) http.HandlerFunc {
            return func(w http.ResponseWriter, r *http.Request) {
                *count++
                writeJSON(w, http.StatusOK, `{"err_code":0}`)
            }
        }(new(int)))

        client := newTestClient(server.URL)
        err := client.DoJSON(context.Background(), &Request{Path: "/test"}, nil)
        assert.NoError(t, err)
    })

    t.Run("TC-CLIENT-002: token expired triggers retry", func(t *testing.T) {
        callCount := 0
        tokenRefreshCount := 0

        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            callCount++
            // 第一次返回 token 过期，第二次成功
            if callCount == 1 {
                writeJSON(w, http.StatusOK, `{"err_code":10022}`)
            } else {
                writeJSON(w, http.StatusOK, `{"err_code":0}`)
            }
        }))
        defer server.Close()

        client := newTestClient(server.URL, WithTokenRefreshHook(func() {
            tokenRefreshCount++
        }))

        err := client.DoJSON(context.Background(), &Request{Path: "/test"}, nil)
        assert.NoError(t, err)
        assert.Equal(t, 2, callCount)
        assert.Equal(t, 1, tokenRefreshCount)
    })
}
```

---

## 6.9 CLI 集成测试

### 测试文件：`cmd/cmd_test.go`

#### 6.9.1 `lark auth status` 测试

```go
func TestAuthStatusCommand(t *testing.T) {
    t.Run("shows valid status", func(t *testing.T) {
        // 准备临时配置
        tmpDir := t.TempDir()
        writeConfig(t, tmpDir, &config.Config{
            UserKey: "key", PluginID: "id", PluginSecret: "sec", BaseURL: "url",
        })
        writeSession(t, tmpDir, &session.Session{
            PluginToken: "token", ExpireAt: time.Now().Unix() + 3600,
        })

        cmd := lark.New()
        cmd.SetArgs([]string{"auth", "status"})
        cmd.SetOut(os.Stdout)

        // 需要 mock TokenProvider 返回成功
        // ...
    })
}
```

---

## 6.10 测试覆盖率目标

| 模块 | 行覆盖率目标 | 关键路径 |
|------|-------------|----------|
| config | 95%+ | Load, Save, ValidateForOpenAPI |
| session | 95%+ | Load, Save, Clear, IsEmpty, IsValid |
| auth/fingerprint | 100% | BuildConfigFingerprint |
| auth/auth_api | 90%+ | FetchPluginToken (含错误分支) |
| auth/token_provider | 95%+ | GetAuthContext, ForceRefresh, 并发场景 |
| openapi/errors | 100% | LookupPolicy, NormalizeError |
| openapi/client | 90%+ | DoJSON 重试逻辑 |

---

## 6.11 测试执行指南

### 本地运行

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 运行特定模块
go test ./internal/config/...
go test ./internal/auth/...

# 运行集成测试（需要 mock server）
go test -tags=integration ./...
```

### CI/CD

```yaml
# .github/workflows/test.yml
- name: Run Tests
  run: |
    go test -race -coverprofile=coverage.out ./...
    go vet ./...
    golint ./...
```

---

## 6.12 Mock 工具建议

| 场景 | 推荐工具 |
|------|----------|
| HTTP Server Mock | `httptest` (标准库) |
| 接口 Mock | `gomock` 或手动 mock |
| 文件系统 Mock | `testhelper` 或 `afero` |
| 时间控制 | `stretchr/testify` mock |

### 手动 Mock 示例（推荐，简单场景）

```go
// 无需代码生成，手动实现接口
type mockTokenProvider struct {
    authCtx *auth.AuthContext
    err     error
}

func (m *mockTokenProvider) GetAuthContext(ctx context.Context) (*auth.AuthContext, error) {
    return m.authCtx, m.err
}

func (m *mockTokenProvider) ForceRefresh(ctx context.Context) (*auth.AuthContext, error) {
    return m.authCtx, m.err
}
```

---

## 一句话总结

这份文档的目标不是直接给出最终生产代码，而是提供一套**可立即开始编码的 Go 工程骨架与实现顺序**，让 [[lark_cli PRD/登录模块与plugin_token缓存方案]] 能从 PRD 顺利过渡到具体实现。
