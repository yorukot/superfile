# Bubble Tea v2 Migration Context

## Issue being worked

- Repository issue: `yorukot/superfile#1367`
- Title: `Migrate to bubbletea v2`
- Additional direction from the issue body:
  - Bubble Tea v2 is considered a large improvement.
  - Beyond the migration guide, the issue specifically calls out:
    - terminal version and name
    - terminfo and termcap capabilities
  - Those features should be used to reduce duplicate code in kitty image preview support.

## Decisions already locked for this branch

- Scope choice from user:
  - minimal Bubble Tea v2 migration plus kitty cleanup
  - not a broader v2 feature adoption pass
- Kitty detection fallback choice from user:
  - query-only
  - no env-var fallback
- Practical interpretation of query-only:
  - use Bubble Tea v2 terminal version/name as the primary signal
  - use Bubble Tea v2 terminal capability queries for terminal name (`TN` / `name`) as the only fallback identity signal
  - do not use:
    - `TERM`
    - `TERM_PROGRAM`
    - `rasterm.IsKittyCapable()`
    - ad hoc terminal-name heuristics outside query responses

## Current repo baseline

- Working branch: `bubbletea-v2`
- Git status at inspection time: clean
- HEAD during inspection:
  - `dd73a40 docs: add download badges to README (#1404)`
- `bubbletea-v2` currently points at the same commit as `main` and `origin/main`.
- There were no local file changes before starting the documentation work for this request.

## Test baseline

- User-reported baseline on this branch: `go test ./...` is green.
- Important because this removes the need to budget for pre-existing failures during the migration.
- Local sandboxed test execution was not a reliable baseline because network/cache access is restricted in this environment, so the user-provided green test run is the relevant project baseline.

User-provided package results:

- `github.com/yorukot/superfile/src/internal` -> `ok`
- `github.com/yorukot/superfile/src/internal/common` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/clipboard` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/filepanel` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/metadata` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/preview` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/processbar` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/prompt` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/rendering` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/sidebar` -> `ok`
- `github.com/yorukot/superfile/src/internal/ui/zoxide` -> `ok`
- `github.com/yorukot/superfile/src/pkg/utils` -> `ok`

Packages without tests are mostly UI composition or glue packages, which is relevant because behavior changes there will need indirect verification via higher-level tests.

## Current dependency state

`go.mod` currently contains direct v1 Charm dependencies:

- `github.com/charmbracelet/bubbles v0.21.0`
- `github.com/charmbracelet/bubbletea v1.3.10`
- `github.com/charmbracelet/lipgloss v1.1.0`

Source:

- `go.mod:10-12`

Interesting detail:

- `go.mod` already has `charm.land/bubbletea/v2 v2.0.2` as an indirect dependency.
- That means v2 is already in the module graph transitively, but the codebase still imports Bubble Tea v1 directly.
- This matters because the migration is not starting from a completely v2-free graph, but it is still a source-level migration.

Source:

- `go.mod:31`

## Where Bubble Tea / Bubbles / Lip Gloss are currently used

### Direct Bubble Tea entry points

- `src/cmd/main.go`
- `src/internal/model.go`
- `src/internal/handle_file_operations.go`
- `src/internal/handle_panel_movement.go`
- `src/internal/model_msg.go`
- `src/internal/key_function.go`
- `src/internal/function.go`
- `src/internal/ui/prompt/model.go`
- `src/internal/ui/zoxide/model.go`
- `src/internal/ui/zoxide/utils.go`
- `src/internal/ui/filemodel/update.go`
- `src/internal/ui/filemodel/dimensions.go`
- `src/internal/ui/sidebar/sidebar.go`
- `src/internal/ui/helpmenu/model_state.go`
- many tests under `src/internal/...`

Source scan:

- `rg -n 'github.com/charmbracelet/bubbletea|github.com/charmbracelet/bubbles|github.com/charmbracelet/lipgloss' src`

### Bubbles usage

- `textinput` is widely used:
  - root model typing modal type in `src/internal/type.go`
  - prompt modal
  - zoxide modal
  - sidebar search
  - file panel search/rename fields
- `progress` is used in common style helpers and process bar rendering.

Source anchors:

- `src/internal/type.go:18`
- `src/internal/common/style_function.go:7-9`
- `src/internal/ui/processbar/process.go:7`
- `src/internal/ui/prompt/type.go:3`
- `src/internal/ui/sidebar/type.go:3`
- `src/internal/ui/filepanel/types.go:7`
- `src/internal/ui/helpmenu/type.go:3`
- `src/internal/ui/zoxide/type.go:4`

### Lip Gloss usage

- Lip Gloss is pervasive across rendering, common styles, preview rendering, and top-level layout.
- The migration should expect broad import churn even where behavior is unchanged.

Representative anchors:

- `src/cmd/main.go:18`
- `src/internal/model.go:19`
- `src/internal/model_render.go:11`
- `src/internal/common/style_function.go:9`
- `src/internal/common/style.go`
- `src/internal/ui/rendering/*`

## Current Bubble Tea root wiring

### Program creation

Current program construction:

- `src/cmd/main.go:153-155`

Current behavior:

- `tea.NewProgram(internal.InitialModel(...), tea.WithAltScreen(), tea.WithMouseCellMotion())`

Implication:

- The code is using Bubble Tea v1 program options for alt screen and mouse mode.
- Bubble Tea v2 expects view-driven configuration instead of those program options.

### Root model init

Current root model `Init()`:

- `src/internal/model.go:52-57`

Behavior:

- returns `tea.Batch(...)`
- sets window title via `tea.SetWindowTitle("superfile")`
- starts `textinput.Blink`
- starts process bar listening command

Implication:

- `tea.SetWindowTitle` currently lives in `Init()`.
- In Bubble Tea v2, window title needs to move to the returned view configuration.

### Root model update

Current root `Update()`:

- `src/internal/model.go:62-112`

Current message handling:

- `tea.WindowSizeMsg`
- `tea.MouseMsg`
- `tea.KeyMsg`
- custom async update messages for zoxide, preview, metadata/process-related work

Current post-update behavior:

- component state updates
- layout updates
- file preview command scheduling
- metadata scheduling

Implication:

- The root message switch is a primary migration hotspot.
- Preview rerender scheduling already exists and can be reused for the kitty query-state transition.

### Root model view

Current `View()`:

- `src/internal/model.go:487-516`

Current signature:

- `func (m *model) View() string`

Current behavior:

- returns `"Loading..."` before first load
- returns warning renders for terminal-too-small cases
- otherwise returns the main composed string render

Implication:

- This must become `View() tea.View`.
- The string rendering logic can remain almost entirely intact if it is wrapped by `tea.NewView(...)`.

## Current key and mouse assumptions

### Legacy key fields are used directly

`src/internal/model.go:275-278` logs:

- `msg.Type`
- `msg.Runes`
- `msg.Paste`
- `msg.Alt`

This is legacy Bubble Tea key-shape code and will need to move to Bubble Tea v2 key fields.

### Test helper currently constructs v1 key messages

`src/pkg/utils/tea_utils.go:5-9`:

- returns `tea.KeyMsg`
- sets `Type: tea.KeyRunes`
- sets `Runes: []rune(msg)`

This helper is heavily reused throughout tests, so changing it is one of the cleanest leverage points in the migration.

### Tests also construct key messages directly

There are many direct literals such as:

- `tea.KeyMsg{Type: tea.KeyEnter}`
- `tea.KeyMsg{Type: tea.KeyDown}`
- `tea.KeyMsg{Type: tea.KeyUp}`
- `tea.KeyMsg{Type: tea.KeyBackspace}`
- `tea.KeyMsg{Type: tea.KeyEsc}`
- `tea.KeyMsg{Type: tea.KeyCtrlC}`
- `tea.KeyMsg{Type: tea.KeyEscape}`
- `tea.KeyMsg{Type: tea.KeyCtrlD}`

Representative locations:

- `src/internal/model_file_operations_test.go`
- `src/internal/model_layout_test.go`
- `src/internal/model_navigation_test.go`
- `src/internal/model_prompt_test.go`
- `src/internal/model_test.go`
- `src/internal/ui/prompt/model_test.go`
- `src/internal/ui/zoxide/utils_test.go`
- `src/internal/ui/zoxide/model_test.go`

Implication:

- The migration is not just import churn.
- Test fixtures and helper semantics are part of the real migration surface.

### Test program harness assumes v1 model/view shape

`src/internal/test_utils_teaprog.go:25-30` creates:

- `tea.NewProgram(m, tea.WithInput(nil), tea.WithOutput(IgnorerWriter{}))`

`src/internal/test_utils_teaprog.go:69-82` also assumes:

- `Update()` returns `tea.Model`
- root model can be cast back to `*model`
- batched commands are still executed as `tea.Batch(...)`

Implication:

- The harness itself may remain conceptually similar, but all direct key message construction and any view-shape assumptions need revalidation after the v2 migration.

## Preview subsystem and kitty-specific state

### Preview model creation

Preview model is created eagerly:

- `src/internal/ui/filemodel/utils.go:13-19`
- `src/internal/ui/preview/model.go:27-47`

Current behavior:

- `filemodel.New(...)` creates `preview.New()`
- `preview.New()` creates:
  - `filepreview.NewImagePreviewer()`
  - thumbnail generator
  - bat command lookup

Existing TODO in code:

- preview initialization is doing IO too early
- it causes unnecessary terminal cell-size detection logs in tests
- it probably should be lazy or moved behind an init path

This is relevant because Bubble Tea terminal query state is also async, so ownership boundaries matter here.

### Image previewer state

`src/pkg/file_preview/image_preview.go:40-61`

Current `ImagePreviewer` fields:

- `cache *cache.Cache[string]`
- `terminalCap *TerminalCapabilities`

Current constructor behavior:

- creates cache
- creates terminal capability object
- immediately calls `InitTerminalCapabilities()`

### Current renderer selection

`src/pkg/file_preview/image_preview.go:64-111`

Current logic:

- validate dimensions
- compute cache key
- if `p.IsKittyCapable()`:
  - try kitty renderer
  - cache kitty result
  - on failure, log and fall back
- otherwise:
  - use ANSI renderer
  - cache ANSI result

Important detail:

- the cache key already includes the renderer enum (`RendererANSI` vs `RendererKitty`).
- That means a later forced rerender after terminal query resolution can coexist cleanly with the existing cache design.

### Current kitty detection

`src/pkg/file_preview/kitty.go:16-45`

Current `isKittyCapable()` behavior:

1. call `rasterm.IsKittyCapable()`
2. if false, read:
   - `TERM_PROGRAM`
   - `TERM`
3. compare case-insensitively against this hard-coded allowlist:
   - `ghostty`
   - `WezTerm`
   - `iTerm2`
   - `xterm-kitty`
   - `kitty`
   - `Konsole`
   - `WarpTerminal`

This is the key duplication and inconsistency target called out by the issue.

### Current kitty clearing behavior

There are two clear paths:

1. package-global clear decision
   - `src/pkg/file_preview/kitty.go:47-54`
2. instance-based clear decision
   - `src/pkg/file_preview/kitty.go:56-63`

The package-global path is currently used by terminal warning renders:

- `src/internal/model_render.go:45`
- `src/internal/model_render.go:72`

The instance path is used by preview rendering:

- `src/internal/ui/preview/render.go:154`
- `src/internal/ui/preview/render.go:165`

Implication:

- After moving to query-owned capability state, the package-global path becomes a poor fit because it has no access to the live query result.
- Warning renders should therefore stop calling the package-global function and instead use the preview instance or a model-owned terminal capability object.

### Current image rendering flow

`src/internal/ui/preview/render.go:62-95`

Current `renderImagePreview(...)` behavior:

- if preview closed: render panel-closed text + clear cmd
- if image preview disabled: render disabled text + clear cmd
- call `m.imagePreviewer.ImagePreview(...)`
- if returned bytes start with kitty escape prefix (`\x1b_G`):
  - render raw kitty output
- otherwise:
  - center ANSI output with Lip Gloss alignment

Important implication:

- The preview layer already treats kitty output and ANSI output differently.
- It does not need a structural redesign for the migration.
- It mostly needs a better capability source and possibly one forced rerender when query state arrives.

### Current kitty geometry dependence

`src/pkg/file_preview/kitty.go:97-162`

Current behavior:

- kitty rendering computes destination rows/cols from:
  - original image dimensions
  - preview area dimensions
  - terminal cell size
- placement row/column depends on:
  - `sideAreaWidth`
  - `common.Config.EnableFilePreviewBorder`

Existing TODOs in file:

- `pkg/file_preview` should not depend on `internal/common`
- kitty preview should ideally not know about global modal width / layout state

Important scope call:

- These TODOs are real, but they are not the core goal of issue `#1367`.
- The migration should avoid pulling them in unless a v2 change forces it.

## Terminal cell-size detection is separate from terminal identity

Current cell-size code lives in:

- `src/pkg/file_preview/utils.go`
- `src/pkg/file_preview/utils_unix.go`
- `src/pkg/file_preview/utils_windows.go`

Current behavior:

- detects terminal cell size lazily or during background init
- uses ioctl on Unix
- uses defaults on Windows
- falls back to default `10x20` or Windows defaults

Important distinction:

- The issue comment mentions terminal version/name and terminfo/termcap capabilities specifically to reduce duplicate kitty support code.
- That is about terminal identity and feature detection.
- It is not the same as the current cell-size mechanism.

Decision recorded:

- keep cell-size detection out of scope for this migration
- do not combine that refactor with the Bubble Tea v2 migration

## Current docs are inconsistent with current code

### Docs currently say env vars drive terminal detection

`website/src/content/docs/getting-started/image-preview.md:17-20` and `:53-60`

Current doc claims:

- superfile automatically detects the terminal using `$TERM` and `$TERM_PROGRAM`
- those variables help determine whether advanced rendering might be possible

### Docs currently publish this support matrix

`website/src/content/docs/getting-started/image-preview.md:21-33`

Notable rows:

- kitty -> supported
- WezTerm -> supported
- Ghostty -> supported
- iTerm2 -> not supported
- Konsole -> not supported

### Code currently disagrees

`src/pkg/file_preview/kitty.go:26-34`

Code treats these as kitty-capable positives:

- `iTerm2`
- `Konsole`
- `WarpTerminal`

This mismatch matters for issue scope because:

- the issue explicitly asks to use Bubble Tea v2 terminal capability features to reduce duplicate code
- this is a concrete case where duplicated/manual detection has already drifted from published docs

Migration implication:

- code and docs should both derive from one normalized query-driven allowlist
- this issue is a good place to remove the drift rather than preserve it

## Current support decision that should be encoded in code

The migration should encode a single allowlist for kitty-positive terminals, based on query results only:

- `kitty`
- `xterm-kitty`
- `ghostty`
- `xterm-ghostty`
- `wezterm`

Reasons:

- `kitty` and `xterm-kitty` are the obvious kitty-native identities
- Ghostty publishes `xterm-ghostty` as a terminal identity
- WezTerm is a known kitty-protocol-capable terminal, but its default `TERM` is `xterm-256color`, so generic xterm-style names must not be treated as capability signals

Practical consequence:

- if WezTerm reports a terminal version/name that identifies itself as WezTerm, treat it as kitty-capable
- if only a generic `xterm-256color` style identifier is available, do not infer WezTerm

## Proposed ownership model after migration

### Preferred direction

Bubble Tea should own terminal-query state.

That means:

- main/root model requests terminal information
- main/root model receives terminal query messages
- main/root model normalizes the terminal identity
- preview/image previewer consumes normalized state

### Why this direction is better than leaving probing inside `pkg/file_preview`

- Bubble Tea v2 is the thing providing terminal version/name and capability query APIs.
- The root model already owns async message handling.
- The preview subsystem already receives async rerenders through model messages.
- Query-only detection is easier to enforce if there is one source of truth in the root model.
- It removes duplicate heuristics currently split across:
  - `rasterm.IsKittyCapable()`
  - env vars
  - docs

### What should not happen

- do not keep `pkg/file_preview` probing `TERM` or `TERM_PROGRAM`
- do not keep `rasterm.IsKittyCapable()` as a hidden fallback
- do not add a second terminal-detection abstraction that competes with Bubble Tea message handling

## Expected async behavior after migration

This is an important behavioral detail for implementation:

1. app starts
2. preview model may exist before terminal query results are returned
3. preview may initially render ANSI
4. terminal query result arrives later
5. root model updates terminal identity state
6. currently focused preview is forced to rerender
7. if terminal is kitty-capable, preview upgrades to kitty output

Why this matters:

- query-only means capability can be unknown at initial render time
- if the code does not explicitly rerender after the query result arrives, image preview may stay stuck on ANSI for the whole session

Good news:

- `filemodel.GetFilePreviewCmd(forcePreviewRender bool)` already has a force-rerender path
- stale preview update protection already exists in `src/internal/ui/filemodel/update.go`

Relevant anchors:

- `src/internal/ui/filemodel/update.go:84-119`
- `src/internal/ui/filemodel/update.go:56-81`

## Files that look like primary migration hotspots

### Program / root model

- `src/cmd/main.go`
- `src/internal/model.go`
- `src/internal/model_render.go`

### Shared key/test helpers

- `src/pkg/utils/tea_utils.go`
- `src/internal/test_utils.go`
- `src/internal/test_utils_teaprog.go`

### Prompt / zoxide / textinput-heavy components

- `src/internal/ui/prompt/model.go`
- `src/internal/ui/zoxide/model.go`
- `src/internal/ui/zoxide/utils.go`

### File preview and kitty support

- `src/internal/ui/preview/model.go`
- `src/internal/ui/preview/render.go`
- `src/internal/ui/filemodel/utils.go`
- `src/internal/ui/filemodel/update.go`
- `src/pkg/file_preview/image_preview.go`
- `src/pkg/file_preview/kitty.go`
- `src/pkg/file_preview/utils.go`

### Docs

- `website/src/content/docs/getting-started/image-preview.md`

### Dependency manifests

- `go.mod`
- `go.sum`
- `gomod2nix.toml`

## Areas that are likely to produce migration bugs

### 1. Key input behavior

Risk:

- many code paths call `msg.String()` today and mix that with direct checks on key type/runes
- migration bugs here can silently break hotkeys, text input, or tests

Highest-risk paths:

- root input handling
- prompt modal
- zoxide modal
- tests that emulate special keys directly

### 2. Mouse wheel behavior

Risk:

- current code handles a single `tea.MouseMsg` and matches string values such as `"wheel up"` / `"wheel down"`
- Bubble Tea v2 splits mouse events more explicitly

Highest-risk path:

- `src/internal/model.go:114-120`

### 3. Preview stuck on ANSI

Risk:

- query-only detection means capability may not be known on first render
- without forced rerender on query arrival, users could lose kitty preview accidentally

### 4. Over-clearing or under-clearing kitty images

Risk:

- current warning renders use package-global clear logic
- after moving to query-owned capability state, stale image clearing can become inconsistent if the code still tries to decide outside the live model

### 5. False positives from generic terminal names

Risk:

- generic xterm-like values must not imply kitty support
- especially important for WezTerm because `xterm-256color` is not a safe capability signal

## Explicitly out of scope

These items came up during inspection but should not be pulled into issue `#1367` unless required for the migration:

- redesigning preview layout math
- refactoring `pkg/file_preview` to remove `internal/common` dependency
- changing terminal cell-size detection strategy
- broader Bubble Tea v2 feature adoption beyond compatibility and kitty detection cleanup
- fixing all existing TODOs around lazy preview subsystem initialization

## Acceptance criteria that should be treated as hard gates

- `go test ./...` remains green on the branch baseline
- all direct Bubble Tea imports are migrated to v2 paths
- all direct Bubbles and Lip Gloss imports are migrated to v2 paths where required
- root `View()` returns a Bubble Tea v2 view object rather than a raw string
- alt screen, mouse mode, and window title are configured via the v2 view model
- kitty capability detection is driven only by terminal query results
- docs match the implemented detection strategy and support matrix

## References gathered during planning

These were the main external references consulted while forming the plan:

- Bubble Tea v2 upgrade guide
  - `https://raw.githubusercontent.com/charmbracelet/bubbletea/v2.0.0/UPGRADE_GUIDE_V2.md`
- Bubble Tea v2 package docs
  - `https://pkg.go.dev/charm.land/bubbletea/v2`
- WezTerm TERM behavior
  - `https://wezterm.org/config/lua/config/term.html`

These references matter mostly for:

- the v1 -> v2 API migration shape
- the move to view-driven alt-screen/window-title configuration
- the availability of terminal query APIs in Bubble Tea v2
- avoiding false positives from generic TERM values

## Short implementation interpretation

If another engineer picks this up, the intended implementation shape is:

- do the source-level Bubble Tea/Bubbles/Lip Gloss v2 migration
- keep superfile behavior stable
- move terminal identity detection into Bubble Tea message handling
- normalize one query-driven kitty allowlist
- inject that result into preview rendering
- rerender preview once capability becomes known
- update docs to match

That is the narrowest change set that fully addresses issue `#1367` without mixing in unrelated cleanup.
