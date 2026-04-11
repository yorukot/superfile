# Bubble Tea v2 + Query-Only Kitty Detection Plan

## Summary

- Upgrade the Charm stack used by superfile from Bubble Tea v1 to Bubble Tea v2.
- Keep the migration intentionally narrow: preserve current product behavior except where Bubble Tea v2 requires API changes, and use the migration to clean up kitty preview capability detection.
- Move kitty capability detection to Bubble Tea v2 terminal queries.
- Use query-only detection:
  - Primary signal: Bubble Tea terminal version/name query result.
  - Fallback signal: Bubble Tea capability query for terminal name (`TN` / `name`).
  - No fallback to `TERM`, `TERM_PROGRAM`, or `rasterm` heuristics.
- Keep terminal cell-size detection out of scope for this issue.

## Locked Decisions

- Scope: minimal Bubble Tea v2 migration plus kitty capability cleanup.
- Detection policy: query-only.
- Supported kitty-positive identities in this migration:
  - `kitty`
  - `xterm-kitty`
  - `ghostty`
  - `xterm-ghostty`
  - `wezterm`
- Do not special-case `iTerm2`, `Konsole`, or `WarpTerminal` in this issue.
- Do not use `xterm-256color` as a signal for kitty support.

## Implementation Changes

### 1. Dependency and import migration

- Update direct dependencies in `go.mod`, `go.sum`, and `gomod2nix.toml`:
  - `github.com/charmbracelet/bubbletea` -> `charm.land/bubbletea/v2`
  - `github.com/charmbracelet/lipgloss` -> `charm.land/lipgloss/v2`
  - `github.com/charmbracelet/bubbles/...` -> `charm.land/bubbles/v2/...`
- Migrate every import site in `src/` and tests.
- Keep the rest of the dependency graph unchanged unless the v2 migration forces a matching version bump.

### 2. Bubble Tea root model migration

- Update the root program setup in `src/cmd/main.go`:
  - stop configuring alt screen and mouse mode via `tea.NewProgram(...)` options
  - keep program construction otherwise unchanged
- Update the root model in `src/internal/model.go`:
  - change `View() string` to `View() tea.View`
  - wrap the existing final render string with `tea.NewView(...)`
  - set view properties on the returned view:
    - alt screen enabled
    - mouse mode enabled
    - window title set to `superfile`
- Remove `tea.SetWindowTitle("superfile")` from `Init()`.

### 3. Key and mouse event migration

- Replace legacy `tea.KeyMsg` usage with Bubble Tea v2 key press handling.
- Replace top-level message switching in `src/internal/model.go` from:
  - `tea.KeyMsg`
  - `tea.MouseMsg`
  to:
  - `tea.KeyPressMsg`
  - Bubble Tea v2 mouse message types, especially wheel handling
- Update all key helpers and tests that currently construct messages through:
  - `src/pkg/utils/tea_utils.go`
  - `src/internal/test_utils.go`
  - `src/internal/test_utils_teaprog.go`
  - direct `tea.KeyMsg{Type: ...}` literals in tests
- Replace code that reads legacy key fields such as:
  - `msg.Type`
  - `msg.Runes`
  - `msg.Alt`
  with Bubble Tea v2 key fields such as:
  - `Code`
  - `Text`
  - `Mod`

### 4. Terminal-query ownership and kitty detection

- Introduce model-owned terminal query state in the Bubble Tea layer instead of probing terminal capability directly from `pkg/file_preview`.
- Request terminal information during startup:
  - terminal version/name query
  - terminal capability query for terminal name (`TN`)
- Handle the resulting Bubble Tea messages in the root model and normalize them into one internal terminal identity model.
- Treat terminal version/name as the primary source.
- Use `TN` capability only as a fallback when terminal version/name is unavailable.
- Remove current kitty-capability probing logic from `src/pkg/file_preview/kitty.go`:
  - `rasterm.IsKittyCapable()`
  - `TERM_PROGRAM`
  - `TERM`
  - hard-coded env-var allowlist matching

### 5. Preview integration

- Thread terminal identity/capability state into the preview subsystem instead of letting `ImagePreviewer` inspect the terminal by itself.
- Update preview construction so the preview model and image previewer receive terminal capability state from the main model.
- Keep the current preview cache model:
  - separate renderer cache keys for ANSI and kitty
  - no global cache invalidation
- When terminal query state changes from unknown or ANSI-only to kitty-capable:
  - force one preview re-render for the currently selected item
  - keep behavior otherwise unchanged

### 6. Kitty clearing behavior

- Remove package-global kitty clearing decisions from warning renders and route clearing through the live preview instance.
- Update warning render paths so they only emit kitty clear commands when the running session is actually known to be kitty-capable.
- Keep current clear-command behavior in preview rendering, but make its capability decision depend on injected query state.

### 7. Documentation sync

- Update `website/src/content/docs/getting-started/image-preview.md` so it matches the migrated implementation.
- Remove the current explanation that superfile detects terminals via `TERM` and `TERM_PROGRAM`.
- Replace it with query-based detection language.
- Fix the support matrix to match the single source of truth used in code.

## Test Plan

- Hard acceptance gate: `go test ./...` stays green.
- Add unit tests for terminal identity normalization:
  - terminal version/name positive cases
  - `TN` fallback positive cases
  - unknown and no-response cases
- Add tests for preview behavior when terminal query results arrive after the preview subsystem is already open:
  - ANSI first
  - upgrade to kitty after query result
  - stale renders not applied
- Update existing model/test helpers to use Bubble Tea v2 message construction.
- Keep current layout and preview tests passing.

## Files Most Likely To Change

- `go.mod`
- `go.sum`
- `gomod2nix.toml`
- `src/cmd/main.go`
- `src/internal/model.go`
- `src/internal/model_render.go`
- `src/pkg/utils/tea_utils.go`
- `src/internal/test_utils.go`
- `src/internal/test_utils_teaprog.go`
- `src/internal/ui/preview/model.go`
- `src/internal/ui/preview/render.go`
- `src/internal/ui/filemodel/utils.go`
- `src/internal/ui/filemodel/update.go`
- `src/pkg/file_preview/image_preview.go`
- `src/pkg/file_preview/kitty.go`
- `website/src/content/docs/getting-started/image-preview.md`

## Explicit Non-Goals

- Do not redesign preview layout.
- Do not refactor terminal cell-size detection.
- Do not remove `pkg/file_preview`'s existing dependency on `internal/common` as part of this issue.
- Do not adopt broader Bubble Tea v2 features that are unrelated to compatibility or kitty detection.

## Acceptance Criteria

- The codebase builds and tests cleanly on the existing branch baseline.
- All direct Bubble Tea, Bubbles, and Lip Gloss imports are on the v2 paths.
- The root view is migrated to Bubble Tea v2 view semantics.
- Key and mouse handling use Bubble Tea v2 message types.
- Kitty capability detection is driven only by terminal queries.
- Docs describe the migrated behavior accurately.
