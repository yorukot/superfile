# Progress

## Completed

- Bubble Tea v2 synthetic key-message migration is done for the current repo state.
- Replaced all `tea.KeyMsg{...}` composite literals with Bubble Tea v2 `tea.KeyPressMsg{Code: ...}` forms.
- Updated direct test-driving call sites that were still constructing legacy key messages:
  - `p.Send(...)`
  - `p.SendDirectly(...)`
  - `TeaUpdate(...)`
  - direct `HandleUpdate(...)` calls in prompt and zoxide tests
- Updated `src/pkg/utils/tea_utils.go`:
  - `TeaRuneKeyMsg` now returns a `tea.KeyPressMsg`
  - single-rune input uses `Code` plus `Text`
  - multi-rune input uses `Code: tea.KeyExtended` plus `Text`
- Updated ctrl-key test cases to v2 form:
  - `ctrl+c` -> `tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}`
  - `ctrl+d` -> `tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl}`
- Updated migration notes to use `KeyPressMsg` terminology in `plan.md`, `migration_context.md`, and `info.md`.
- Bubbles v2 Width-field migration is done for `textinput` and `progress` usages:
  - replaced direct `textinput.Width = ...` writes with `SetWidth(...)`
  - replaced direct `progress.Width = ...` writes with `SetWidth(...)`
  - replaced `textinput.Width` reads with `Width()`
  - updated file panel, sidebar, help menu, prompt, zoxide, validation, and process bar call sites

## Verified

- `rg -n "tea\\.KeyMsg\\{" .` returns no matches.
- `rg -n "SearchBar\\.Width|Rename\\.Width|rename\\.Width|textInput\\.Width|Progress\\.Width|searchBar\\.Width" src` only returns `Width()` getter usage in validation, not direct field access.
- `src/pkg/utils` builds and tests with the updated `TeaRuneKeyMsg` helper.
- Targeted package test runs for the width migration still stop in unrelated `lipgloss.Color` compile errors under `src/internal/common`.

## Current Known Blocker

- Broader test runs are still blocked by unrelated `lipgloss.Color` compile errors in `src/internal/common`, which are outside this key-message-only change set.
