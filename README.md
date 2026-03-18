# lark_cli

Command line tools for Feishu Project (Lark Project) OpenAPI.

## Build

```bash
go build -o lark .
```

## Test

```bash
go test ./...
```

## Usage

### 1) Configure static settings (recommended)

Create `~/.lark/config.json`:

```json
{
  "base_url": "https://project.feishu.cn",
  "plugin_id": "<your_plugin_id>",
  "plugin_secret": "<your_plugin_secret>",
  "session_path": "/Users/<you>/.lark/session.json"
}
```

You can still use environment variables, and they have higher priority than file values.

```bash
export LARK_BASE_URL="https://project.feishu.cn"
export LARK_PLUGIN_ID="<your_plugin_id>"
export LARK_PLUGIN_SECRET="<your_plugin_secret>"
export LARK_SESSION_PATH="$HOME/.lark/session.json"
```

### 2) Login (save user_key only)

```bash
lark login -w user_key <user_key>
```

### 3) Check auth status

```bash
# Read local session/user_key status
lark auth status

# Force fetch plugin_access_token from /open_api/authen/plugin_token
lark auth status --refresh-plugin-token
```

### 4) Logout

```bash
lark logout
```

## Auth flow in this CLI

- `login` only saves `user_key` to local session file.
- `plugin_access_token` is fetched when you run `auth status --refresh-plugin-token` (or when a future API command requests headers).
- OpenAPI calls should use:
  - `X-Plugin-Token: <plugin_access_token>`
  - `X-User-Key: <user_key>`

## Configuration priority

Configuration is loaded with this precedence:

1. Environment variables (`LARK_*`) - highest priority
2. `~/.lark/config.json`
3. Built-in defaults

Defaults:

- `base_url`: `https://project.feishu.cn`
- `session_path`: `~/.lark/session.json`

## Config file vs session file

- `~/.lark/config.json` stores static configuration (`base_url`, `plugin_id`, `plugin_secret`, `session_path`).
- `~/.lark/session.json` stores runtime login state (`user_key`) and token cache.

## Environment variables

- `LARK_SESSION_PATH`: session JSON file path (default: `~/.lark/session.json`)
- `LARK_BASE_URL`: default `https://project.feishu.cn`
- `LARK_PLUGIN_ID`: required to fetch plugin token
- `LARK_PLUGIN_SECRET`: required to fetch plugin token
