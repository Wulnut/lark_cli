---
title: 登录模块与 plugin_token 缓存方案
source:
author:
published:
created: 2026-03-19
updated: 2026-03-19
description: lark_cli plugin_token 缓存方案设计，定义config.json与session.json分离、token按需获取与自动刷新、配置变更导致缓存失效、请求头统一注入及错误码集中处理等策略
tags: [lark_cli, PRD, auth, token, cache, plugin_token, OpenAPI]
category: lark_cli
status: draft
related_docs:
  - "[[飞书项目/飞书项目OpenAPI获取访问凭证]]"
  - "[[飞书项目/飞书项目OpenAPI鉴权流程]]"
  - "[[lark_cli PRD/登录模块设计]]"
  - "[[lark_cli PRD/登录模块与plugin_token缓存方案-Go工程落地设计]]"
---

# 登录模块与 plugin_token 缓存方案

## 背景

当前 `lark_cli` 需要支持调用飞书项目 OpenAPI。现有使用方式中，用户通过如下命令登录：

```bash
lark login --user-key <user_key>
```

当前实现只会把部分会话信息保存到 `.lark/session.json` 中，但在实际调用 OpenAPI 时，除了 `user_key` 之外，还需要以下插件身份配置：

- `plugin_id`
- `plugin_secret`

通过 `plugin_id + plugin_secret` 可以获取 `plugin_token`，后续业务请求需在 Header 中携带：

- `X-PLUGIN-TOKEN`
- `X-USER-KEY`

获取 `plugin_token` 的请求示例如下：

```bash
curl --location 'https://project.feishu.cn//open_api/authen/plugin_token' \
--header 'Content-Type: application/json' \
--data '{
  "plugin_id": "<your_plugin_id>",
  "plugin_secret": "<your_plugin_secret>"
}'
```

返回示例：

```json
{
  "error": {
    "display_msg": {},
    "code": 0,
    "msg": "success"
  },
  "data": {
    "token": "p-example-plugin-token",
    "expire_time": 7200
  }
}
```

业务接口调用示例：

```bash
curl --location 'https://project.feishu.cn//open_api/user/query' \
--header 'X-PLUGIN-TOKEN: p-example-plugin-token' \
--header 'X-USER-KEY: <your_user_key>' \
--header 'Content-Type: application/json' \
--data-raw '{
  "user_keys": [""],
  "out_ids": [""],
  "emails": ["user@example.com"],
  "tenant_key": ""
}'
```

因此，CLI 需要系统性解决以下问题：

1. 如何保存长期配置
2. 如何缓存和刷新 `plugin_token`
3. 如何统一注入请求头
4. 如何集中处理 OpenAPI 错误码

---

## 当前问题

### 1. `login` 命令职责混杂

如果把 `user_key`、`plugin_id`、`plugin_secret` 全部放到 `login` 命令中，会导致：

- `login` 语义过重
- 登录与插件配置耦合
- 命令扩展性差
- 后续支持多 profile / 多环境时不易演进

### 2. `plugin_token` 不是静态配置

`plugin_token` 由接口动态签发，并且有有效期（例如 7200 秒）。如果在登录阶段就获取并长期保存，会面临：

- 长时间不使用导致 token 过期
- 配置修改后仍沿用旧 token
- 用户误以为 token 是稳定配置
- 业务请求阶段仍需处理过期刷新

### 3. 当前缺少集中式错误处理机制

OpenAPI 错误码较多，且错误类型不同：

- 凭证错误
- 权限错误
- 参数错误
- 资源不存在
- 限流错误
- 服务端错误

如果由各命令分别处理，会导致：

- 错误逻辑分散
- 提示信息不统一
- 重试策略不一致
- 维护成本高

---

## 设计目标

本方案的目标如下：

1. 区分长期配置与短期会话缓存
2. 让 `login` 命令职责更清晰
3. 实现 `plugin_token` 的按需获取与自动刷新
4. 统一管理请求头注入逻辑
5. 建立集中式错误码处理机制
6. 为未来多 profile / 多环境扩展预留空间

---

## 核心设计原则

### 1. 配置与会话分离

- `config.json` 保存长期配置
- `session.json` 保存运行期缓存

### 2. 登录与插件配置分离

- `login` 只负责用户身份配置
- `plugin set` 只负责插件配置

### 3. token 按需获取

不在 `login` 阶段强制获取 `plugin_token`，而是在首次真正发起 OpenAPI 请求时自动获取。

### 4. 请求统一走中间层

所有 OpenAPI 请求必须统一经过请求中间层，由中间层负责：

- 获取有效 token
- 注入 `X-PLUGIN-TOKEN`
- 注入 `X-USER-KEY`
- 处理特定错误码的刷新与重试

---

## 推荐命令设计

### 1. 用户登录

```bash
lark login --user-key <user_key>
```

职责：

- 保存当前用户身份
- 建立 CLI 默认上下文

说明：

- 不强制请求 `plugin_token`
- 不强制校验 OpenAPI 是否可用

### 2. 插件配置

```bash
lark plugin set --id <plugin_id> --secret <plugin_secret>
```

职责：

- 保存 `plugin_id`
- 保存 `plugin_secret`

说明：

- 插件配置不等于登录动作
- 单独拆出后命令语义更清晰

### 3. 状态检查（推荐新增）

```bash
lark auth status
```

职责：

- 检查当前 `user_key` 是否存在
- 检查 `plugin_id` / `plugin_secret` 是否完整
- 检查本地是否存在 token 缓存
- 检查 token 是否过期

### 4. 手动刷新（可选）

```bash
lark auth refresh
```

职责：

- 主动重新获取 `plugin_token`
- 更新本地会话缓存

### 5. 环境诊断（推荐新增）

```bash
lark doctor
```

职责：

- 检查配置完整性
- 检查 token 获取是否正常
- 检查请求头是否正确注入
- 输出人类可读的诊断信息

---

## 本地文件设计

### 1. `.lark/config.json`

用于保存长期配置，允许手动编辑。

示例：

```json
{
  "user_key": "<your_user_key>",
  "plugin_id": "<your_plugin_id>",
  "plugin_secret": "******",
  "base_url": "https://project.feishu.cn/open_api"
}
```

建议字段：

- `user_key`
- `plugin_id`
- `plugin_secret`
- `base_url`

职责：

- 保存用户显式输入的配置
- 作为长期配置来源
- 不保存短期 token

### 2. `.lark/session.json`

用于保存短期会话缓存，不建议手动编辑。

示例：

```json
{
  "plugin_token": "p-example-plugin-token",
  "expire_at": 1760000000,
  "obtained_at": 1759992800,
  "config_fingerprint": "hash(user_key+plugin_id+plugin_secret+base_url)"
}
```

建议字段：

- `plugin_token`
- `expire_at`
- `obtained_at`
- `config_fingerprint`

职责：

- 保存短期 token
- 记录 token 失效时间
- 检测 token 是否仍匹配当前配置

---

## token 缓存方案

### 1. 获取时机

推荐采用 **按需获取**：

- 用户执行 `login` 或 `plugin set` 时，只保存配置
- 第一次真正调用业务 API 时，自动获取 `plugin_token`
- 后续优先使用本地缓存

原因：

1. 登录动作不依赖网络
2. 避免获取后长时间不用导致立即进入过期倒计时
3. 更符合 token 作为短期凭证的语义
4. 可减少命令副作用

### 2. 过期策略

服务端返回的 `expire_time` 例如为 `7200` 秒。

本地建议增加安全窗口（例如 300 秒）：

```text
expire_at = now + expire_time - safety_window
```

这样可避免：

- 本地时钟偏差
- 网络延迟
- token 在请求过程中失效

### 3. 刷新策略

请求前执行以下逻辑：

1. 如果本地无 token，则获取
2. 如果 token 已过期或即将过期，则刷新
3. 如果 token 存在且有效，则直接使用

### 4. 配置变更导致缓存失效

若用户修改下列任一字段：

- `user_key`
- `plugin_id`
- `plugin_secret`
- `base_url`

则旧 token 应立即失效。

推荐方式：

- 在 `session.json` 中记录 `config_fingerprint`
- 每次请求前重新计算当前配置指纹
- 如果指纹不一致，则丢弃旧 token 并重新获取

---

## 请求流程设计

### 请求前

1. 读取 `config.json`
2. 校验 `user_key`、`plugin_id`、`plugin_secret`
3. 读取 `session.json`
4. 根据缓存和过期时间获取有效 `plugin_token`
5. 统一构建请求头：
   - `X-PLUGIN-TOKEN`
   - `X-USER-KEY`
   - `Content-Type: application/json`

### 请求后

如果返回以下 token 相关错误：

- `10021 Token Not Exist`
- `10022 Check Token Failed`
- `10211 Token Info Is Invalid`

则执行：

1. 强制刷新 token
2. 自动重试一次请求
3. 若仍失败，则直接报错

注意：

- 自动重试最多一次
- 避免无限循环

---

## 错误码集中处理方案

错误码不应只作为文档表格保存，而应做成“错误码 -> 分类 -> 处理策略”的程序化映射。

### 推荐分类

- `auth`：认证与凭证错误
- `permission`：权限错误
- `param`：参数错误
- `not_found`：资源不存在
- `state_conflict`：状态流转冲突
- `rate_limit`：限流错误
- `server`：服务端错误
- `unknown`：未知错误

### 推荐错误策略字段

每个错误码至少需要以下信息：

- 错误码 `Code`
- 默认文案 `Message`
- 错误分类 `Category`
- 是否可重试 `Retryable`
- 是否需要刷新 token `RefreshToken`
- 最大重试次数 `MaxRetry`
- 面向用户的提示文案 `UserHint`
- 面向开发调试的提示 `DevHint`

### 和 token 缓存最相关的关键错误码

#### 认证类

| err_code | err_msg                            | 建议处理                       |
| -------- | ---------------------------------- | -------------------------- |
| `10021`  | Token Not Exist                    | 刷新 token 并重试一次             |
| `10022`  | Check Token Failed                 | 刷新 token 并重试一次             |
| `10211`  | Token Info Is Invalid              | 刷新 token 并重试一次             |
| `20039`  | Plugin Token Must Have User Key... | 不刷新 token，修复请求头注入逻辑        |
| `20042`  | X-User-Key Is Wrong...             | 不刷新 token，检查 `user_key` 配置 |
| `30006`  | User Not Found                     | 不刷新 token，检查用户身份或可见性       |

#### 权限类

| err_code | err_msg                            | 建议处理                              |
| -------- | ---------------------------------- | --------------------------------- |
| `10001`  | No Permission                      | 提示当前操作无权限                         |
| `10301`  | Check Token Perm Failed            | 不刷新 token，优先检查接口权限、版本发布、插件安装、空间权限 |
| `10302`  | Check User Error, User Is Resigned | 提示用户已离职                           |
| `10404`  | No Project Permission              | 提示用户无空间访问权限                       |

#### 限流与服务端类

| err_code     | err_msg                             | 建议处理         |
| ------------ | ----------------------------------- | ------------ |
| `10429`      | API Request Frequency Limit         | 指数退避重试       |
| `10430`      | API Request Idempotent Limit        | 检查幂等串，不自动重试  |
| `50006`      | openapi system err / RPC Call Error | 短暂重试 1~2 次   |
| `1000051942` | commercial usage exceeded           | 提示额度超限，不自动重试 |

---

## 推荐的错误策略矩阵

| 分类                  | 代表错误码                     | 自动刷新 token | 自动重试 | 最大次数 | 处理建议                   |
| ------------------- | ------------------------- | ---------- | ---- | ---- | ---------------------- |
| auth/token missing  | `10021`                   | 是          | 是    | 1    | 刷新 token 后重试一次         |
| auth/token invalid  | `10022`, `10211`          | 是          | 是    | 1    | 强制刷新 token 后重试一次       |
| auth/header missing | `20039`, `20042`          | 否          | 否    | 0    | 修请求头注入逻辑或检查 `user_key` |
| permission          | `10001`, `10301`, `10404` | 否          | 否    | 0    | 提示权限、发布、安装、空间权限问题      |
| param               | `20005`, `20006`, `9999`  | 否          | 否    | 0    | 直接报参数错误                |
| not_found           | `13001`, `30005`, `30009` | 否          | 否    | 0    | 提示资源不存在                |
| rate_limit          | `10429`                   | 否          | 是    | 2~3  | 指数退避重试                 |
| idempotent_conflict | `10430`                   | 否          | 否    | 0    | 检查 `X-IDEM-UUID`       |
| server              | `50006`                   | 否          | 是    | 1~2  | 服务端瞬时异常，短重试            |
| business_limit      | `1000051942`              | 否          | 否    | 0    | 提示额度超限                 |

---

## 安全建议

### 1. `plugin_secret` 属于敏感配置

不建议长期通过命令行明文传入，例如：

```bash
lark plugin set --secret xxx
```

因为可能出现在 shell history 中。

更推荐：

- 交互式输入 secret（不回显）
- 通过环境变量注入
- 后续版本接入系统密钥链

### 2. `session.json` 应允许自动重建

即使 `session.json` 被删除，只要 `config.json` 完整，CLI 仍应能自动获取新 token 并恢复工作。

---

## 并发刷新建议

如果 CLI 后续会并发请求多个接口，需要避免多个请求同时刷新 token。

建议在 token 刷新逻辑中增加：

- 进程内互斥锁，或
- `singleflight` 机制

目的：

- 避免同时重复调用 token 接口
- 避免多个协程同时覆盖 `session.json`
- 避免随机出现旧 token / 新 token 混用问题

---

## 非目标

当前版本建议明确以下内容不在本期范围内：

1. 不实现完整多 profile 切换能力，只预留结构
2. 不接入系统级安全存储（如 Keychain / Credential Manager）
3. 不对全部 OpenAPI 参数做本地预校验
4. 不做无限重试与自动恢复
5. 不在首期覆盖全部业务错误码的人类语义化解释，只先覆盖通用错误与高频错误

---

## 推荐落地优先级

### P0

1. 拆分 `config.json` 与 `session.json`
2. 将 token 获取改为按需获取
3. 统一请求头注入
4. 接入核心认证错误码：
   - `10021`
   - `10022`
   - `10211`

### P1

5. 接入权限错误码：
   - `10001`
   - `10301`
   - `10404`
6. 接入限流与服务端重试：
   - `10429`
   - `50006`
7. 增加 `config_fingerprint`

### P2

8. 增加 `lark auth status`
9. 增加 `lark doctor`
10. 加入并发刷新保护

---

## 最终结论

1. `plugin_token` 是短期会话凭证，不是长期静态配置
2. `login` 命令只负责用户身份配置，不负责完整插件配置
3. `plugin_id/plugin_secret` 应独立配置
4. `plugin_token` 推荐采用“按需获取 + 本地缓存 + 提前刷新 + 特定错误自动刷新一次”的策略
5. 所有 OpenAPI 请求必须统一通过请求中间层自动注入：
   - `X-PLUGIN-TOKEN`
   - `X-USER-KEY`
6. 错误码处理应集中建模为“错误码 -> 分类 -> 策略”的映射体系
7. 优先支持以下自动化处理：
   - `10021/10022/10211`：刷新 token 并重试一次
   - `10301`：提示权限、发布、安装问题，不自动刷新
   - `10429`：退避重试
   - `50006`：短重试
8. 通过 `config_fingerprint` 保证配置变更后旧 token 自动失效

---

## 相关待落地工程项

- 输出一版 Go 伪代码与接口设计稿：
  - config
  - session
  - token provider
  - error policy
  - request middleware
- 将错误码表拆分为程序可消费的数据结构
- 为 CLI 输出增加用户提示与 debug 提示两层文案
