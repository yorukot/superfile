# CI Test Failure Investigations

## 1. TestAsyncPreviewPanelSync (FIXED)

### Root Cause
Bubble Tea's `Program.Run()` sends an initial `WindowSizeMsg` via `go p.Send(resizeMsg)` (a goroutine). With no real terminal (`WithInput(nil)` + `IgnorerWriter`), the size defaults to `{0, 0}`. This goroutine races with the test's `p.Send(WindowSizeMsg{480, 192})`. On CI, the `{0, 0}` was scheduled after the test's message, overwriting good dimensions.

With preview width reduced to 9, content was truncated ("File 1 co" instead of "File 1 content"), and subsequent async preview results were rejected due to dimension mismatches.

### Fix Applied
Added `tea.WithWindowSize(DefaultTestModelWidth, DefaultTestModelHeight)` to `NewTeaProg()` in `test_utils_teaprog.go`, so `Run()` sends `{120, 48}` instead of `{0, 0}`.

### Key Evidence (CI logs)
```
WindowSizeMsg → fullWidth=120   (sync from setModelParamsForTest)
WindowSizeMsg → fullWidth=480   (test's p.Send - processed 2nd)
WindowSizeMsg → fullWidth=0     (tea.Run() goroutine - processed 3rd, overwrites!)
Preview id=1 (w=229) → IGNORED (current dimensions now w=9)
Preview id=2 (w=9)   → accepted, content truncated
```

### Relevant Code
- `tea.go:1073`: `go p.Send(resizeMsg)` — goroutine that races with user sends
- `test_utils_teaprog.go:26`: Program creation (fix location)
- `ui/filemodel/update.go:74-80`: Stale dimension check that rejects preview results

---

## 2. TestCursorOutOfBoundsAfterDirectorySwitch (RCA CONFIRMED)

### Symptom
After navigating from dir1 (10 files, cursor at 8) to dir2 (5 files), `ElemCount()` returns 10 instead of 5. Fails on CI (macOS), passes locally 20/20.

### Root Cause: TOCTOU race in `getDirectoryElements`

`getDirectoryElements` reads `m.Location` at two separate points:
- Line 18: `os.ReadDir(m.Location)` — reads filesystem entries
- Line 34: `sortFileElement(..., m.Location)` — builds Element structs with `Location: filepath.Join(location, item.Name())`

When the test thread changes `m.Location` between these two reads, it creates phantom elements:
- `os.ReadDir(dir1)` returns 10 filesystem entries from dir1
- `sortFileElement(..., dir2)` builds Elements with `Location: dir2/a.txt, dir2/b.txt...`
- Result: 10 elements that appear to belong to dir2 but contain dir1's entries

This makes `NeedsReRender()` return `false` (because `filepath.Dir(element[0].Location) == m.Location`), so the corruption is undetectable by subsequent update cycles.

### Race Sequence (confirmed via debug logging)

**Locally (passes):**
1. Event loop: enters `UpdateElementsIfNeeded`, calls `getDirectoryElements(dir1)` → 10 entries
2. Test thread: `updateCurrentFilePanelDir(dir2)` → changes `m.Location` to dir2
3. Event loop: `sortFileElement(..., m.Location=dir2)` → writes 10 phantom elements with dir2 paths
4. Test thread: `TeaUpdate(nil)` → `UpdateElementsIfNeeded` → reloads from dir2 → 5 elements ✓

**CI (fails — timing flips steps 3-4):**
1. Event loop: enters `UpdateElementsIfNeeded`, calls `getDirectoryElements(dir1)` → 10 entries
2. Test thread: changes Location to dir2, `TeaUpdate(nil)` reloads → 5 elements
3. Event loop: writes 10 phantom elements → overwrites the correct 5 ✗
4. Test asserts `ElemCount() == 5` → gets 10

### Why `assert.Eventually(cursor==8)` Doesn't Guarantee Safety

In `model.Update()`:
- Line 80: `handleKeyInput(down)` → sets cursor=8
- Line 104: `updateModelStateAfterMsg()` → `UpdateElementsIfNeeded` → `getDirectoryElements` → `os.ReadDir` (I/O)

`assert.Eventually` reads cursor==8 (set at line 80) while the event loop is still executing line 104's `getDirectoryElements`. The test thread then changes `m.Location` between `os.ReadDir(m.Location)` (line 18) and `sortFileElement(..., m.Location)` (line 34).

### TOCTOU Confirmed Locally

Added a check comparing `m.Location` at ReadDir vs sortFileElement. Hit 2 out of 30 runs:
```
WARN TOCTOU detected in getDirectoryElements
  locAtReadDir=.../dir1
  locAtSort=.../dir2
  entryCount=10
```

### Root Problem
`navigateToTargetDir` mutates the model directly while the event loop is still inside the same `Update()` call that set cursor==8. The event loop's `getDirectoryElements` reads `m.Location` twice (lines 18 and 34), and the test thread changes it between those reads.

### Fix Options
1. **Test fix**: Have `navigateToTargetDir` send directory change through the event loop (`p.Send`) instead of direct model mutation, so all model access is serialized through the event loop
2. **Test fix**: Wait for event loop to be idle before direct mutation (e.g., add a short sleep or barrier after `assert.Eventually`)
3. **Code fix**: Capture `m.Location` once at the start of `getDirectoryElements` and pass it through, preventing the TOCTOU

### Relevant Files
- `src/internal/ui/filepanel/get_elements.go:17-34` — `getDirectoryElements` with TOCTOU (lines 18 and 34 read `m.Location`)
- `src/internal/ui/filepanel/sort.go:113-126` — `sortFileElement` uses the `location` param for Element.Location
- `src/internal/model_navigation_test.go:160-214` — test code
- `src/internal/test_utils.go:209-216` — `navigateToTargetDir` (direct model mutation)
- `src/internal/model.go:104` — `updateModelStateAfterMsg()` called every Update