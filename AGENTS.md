# AGENTS.md — superfile

## Build / Lint / Test Commands

```bash
make build        # Build only (skip tests)
make test         # Run unit tests: go test ./...
make lint         # Run golangci-lint
make dev          # Full workflow: tidy → fmt → lint → test → build
make testsuite    # Full workflow + Python integration tests
make clean        # Remove bin/
```

### Running a single test

```bash
go test ./src/internal/ -run TestName
go test ./src/internal/ -run TestName -v          # verbose
go test ./src/internal/ -run TestName/SubTest     # subtest
go test ./src/internal/ -run TestName -count=1    # disable cache
```

### Lint formatting

```bash
golangci-lint fmt          # Format with goimports + golines
go fmt ./...               # Standard Go formatting
```

### Build binary

```bash
CGO_ENABLED=0 go build -o ./bin/spf   # or: ./build.sh
```

## Code Style

### Imports

- Three groups separated by blank lines: stdlib → third-party → internal
- Internal imports use full path: `github.com/yorukot/superfile/src/internal/common`
- Use aliases for disambiguation: `tea "github.com/charmbracelet/bubbletea"`
- `goimports` enforces local-prefixes: `github.com/yorukot/superfile`

### Formatting

- Max line length: **120 chars** (enforced by `golines`)
- Run `golangci-lint fmt` or `go fmt ./...` before committing

### Naming

- `PascalCase` for exported identifiers, `camelCase` for unexported
- Sentinel errors: prefix with `Err` (e.g., `ErrNotFound`)
- Error types: suffix with `Error` (e.g., `TomlLoadError`)
- Test structs: descriptive field names (`expectedZipName`, `shouldClear`)
- Method receivers: short, consistent names per type

### Error Handling

- Return errors, never panic (except in test setup via `os.Exit`)
- Wrap with `fmt.Errorf("context: %w", err)`; create with `errors.New()`
- Always check: `if err != nil { return ... }`
- Use `slog` for logging: `slog.Error("msg", "key", value)`
- Inline error handling (`if err := ...; err != nil`) is acceptable

### Types & Structs

- No embedded `sync.Mutex` or `sync.RWMutex` (enforced by linter)
- Use struct tags for (un)marshaling (enforced by `musttag`)
- No named returns (enforced by `nonamedreturns`)
- Global variables allowed only in config/icon files (use `//nolint:gochecknoglobals` elsewhere with justification)

### Testing

- **Table-driven tests** are the dominant pattern
- Use `testify/assert` and `testify/require` for assertions
- Test files use same package (`package internal`), not `_test`
- Use `t.TempDir()` and `t.Cleanup()` for cleanup
- Async tests: `assert.Eventually()` with `1s` timeout / `10ms` tick
- `TestMain` for shared setup/teardown in `model_test.go`
- No `t.Parallel()` (disabled in CI due to false positives)

### Linter

- 80+ linters enabled via `.golangci.yaml` (very strict)
- Always run `make lint` before committing
- `//nolint:lintname` requires explanation
- PR titles must follow [Conventional Commits](https://www.conventionalcommits.org/): `feat(scope): description`

### Architecture

- Bubble Tea Elm architecture: `Init()`, `Update(msg tea.Msg)`, `View()`
- UI components in `src/internal/ui/` (filepanel, sidebar, preview, etc.)
- Shared utilities in `src/internal/common/` and `src/pkg/`
- Config/theme TOML files embedded via `//go:embed`
