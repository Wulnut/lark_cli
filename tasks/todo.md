# Todo

- [x] Implement layered config loading in `internal/config/config.go` (defaults + `~/.lark/config.json` + env overrides)
- [x] Switch root command config entrypoint to `config.Load()` in `cmd/root.go`
- [x] Expand config tests in `test/configtest/config_test.go`
- [x] Update README for `config.json` usage and precedence
- [x] Run verification (`go test`, `go test -race`, `go vet`, `gofmt`)

## Review

- Implemented `config.Load()` with precedence: env > config.json > defaults.
- `cmd/root.go` now uses `config.Load()`.
- Added tests for defaults/file/env precedence/missing file/invalid JSON/config-based credential validation.
- Updated README with `~/.lark/config.json` usage and config/session responsibility split.

1. 没有创建config.json
2. plugin id和secret 没有配置的地方
3. 如果创建空的config.json会出现问题