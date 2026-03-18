# GitHub CI Design (Go)

**Goal:** Add a minimal GitHub Actions CI gate that runs on PRs and pushes to main, verifying the Go code builds and tests pass.

## Scope

In scope:
- Run `go test ./...` and `go test -race ./...`
- Run `go vet ./...`
- Optionally enforce formatting via `gofmt -l .` (fail if output is non-empty)

Out of scope:
- Release automation (GoReleaser)
- Multi-OS test matrix
- Dependency vulnerability scanning

## Triggers

- `pull_request` (all PRs)
- `push` on `main`

## Implementation

Create:
- `.github/workflows/ci.yml`

Workflow:
- Use `actions/checkout@v4`
- Use `actions/setup-go@v5` with `go-version-file: go.mod` and caching enabled
- Steps:
  1. `go test ./... -count=1`
  2. `go test -race ./... -count=1`
  3. `go vet ./...`
  4. `test -z "$(gofmt -l .)"`

## Security / Permissions

- `permissions: contents: read`

## Success Criteria

- A PR with failing tests is blocked by CI
- A PR with data races (if detectable) is blocked by CI
- A PR with `go vet` issues is blocked by CI
