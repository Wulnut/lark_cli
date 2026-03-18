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

### 1) Configure environment variables

```bash
export LARK_BASE_URL="https://project.feishu.cn"
export LARK_PLUGIN_ID="<your_plugin_id>"
export LARK_PLUGIN_SECRET="<your_plugin_secret>"
# Optional, default is ~/.lark/session.json
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

## Environment variables

- `LARK_SESSION_PATH`: session JSON file path (default: `~/.lark/session.json`)
- `LARK_BASE_URL`: default `https://project.feishu.cn`
- `LARK_PLUGIN_ID`: required to fetch plugin token
- `LARK_PLUGIN_SECRET`: required to fetch plugin token
