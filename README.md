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

### Login

```bash
# Save user_key
lark login -w user_key <user_key>
```

### Auth status

```bash
lark auth status

# Force refresh plugin_access_token (requires env vars)
lark auth status --refresh-plugin-token
```

### Logout

```bash
lark logout
```

## Environment variables

- `LARK_SESSION_PATH`: session JSON file path (default: OS user config dir + `lark/session.json`)
- `LARK_BASE_URL`: default `https://project.feishu.cn`
- `LARK_PLUGIN_ID`: required to fetch plugin token
- `LARK_PLUGIN_SECRET`: required to fetch plugin token
