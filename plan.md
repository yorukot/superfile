# FilePanel Export Refactoring Plan

## Phase 1: Disable GitHub Workflows
1. Remove or rename `.github/workflows/` directory to `.github/workflows.disabled/`
2. This will prevent workflows from running on push to save GitHub CPU

## Phase 2: Export Core Types (src/internal/type.go)
Note - Update the usage and ensure compilation after each change.
Keep it in many commits - ensure each commit compiles.

1. Run `go build -o bin/spf` after each major phase
2. Fix any compilation errors
3. Run tests: `go test ./...`
4. Linter `golangci-lint --fix`


### 2.1 Export filePanel struct → FilePanel
- Line 156: `type filePanel struct` → `type FilePanel struct`
- Export all fields:
  - `cursor` → `Cursor`
  - `render` → `Render` (or `RenderIndex` for clarity)
  - `location` → `Location`
  - `sortOptions` → `SortOptions`
  - `panelMode` → `PanelMode`
  - `selected` → `Selected`
  - `element` → `Element`
  - `directoryRecords` → `DirectoryRecords`
  - `rename` → `Rename`
  - `renaming` → `Renaming`
  - `searchBar` → `SearchBar`
  - `lastTimeGetElement` → `LastTimeGetElement`

### 2.2 Export element struct → Element
- Line 195: `type element struct` → `type Element struct`
- Export all fields:
  - `name` → `Name`
  - `location` → `Location`
  - `directory` → `Directory`
  - `metaData` → `MetaData`
  - `info` → `Info`

### 2.3 Export directoryRecord struct → DirectoryRecord
- `type directoryRecord struct` → `type DirectoryRecord struct`
- Fields:
  - `directoryCursor` → `DirectoryCursor`
  - `directoryRender` → `DirectoryRender`

### 2.4 Export sortOptionsModel struct → SortOptionsModel
- `type sortOptionsModel struct` → `type SortOptionsModel struct`
- Fields:
  - `width` → `Width`
  - `height` → `Height`
  - `open` → `Open`
  - `cursor` → `Cursor`
  - `data` → `Data`

### 2.5 Export sortOptionsModelData struct → SortOptionsModelData
- `type sortOptionsModelData struct` → `type SortOptionsModelData struct`
- Fields:
  - `options` → `Options`
  - `selected` → `Selected`
  - `reversed` → `Reversed`

### 2.6 Export panelMode type → PanelMode
- Need to find where panelMode is defined (likely as int or uint)
- Export constants: `browserMode` → `BrowserMode`, `selectMode` → `SelectMode`

### Phase 2.7: Export FilePanel Methods (src/internal/file_panel.go)

All methods should be exported (capitalize first letter):
- `getSelectedItem()` → `GetSelectedItem()`
- `resetSelected()` → `ResetSelected()`
- `getSelectedItemPtr()` → `GetSelectedItemPtr()`
- `changeFilePanelMode()` → `ChangeFilePanelMode()`
- `updateCurrentFilePanelDir()` → `UpdateCurrentFilePanelDir()`
- `parentDirectory()` → `ParentDirectory()`
- `handleResize()` → `HandleResize()`
- Any other panel methods found

## Phase 2.8: Update fileModel struct (src/internal/type.go)
- Update `filePanels []filePanel` → `filePanels []FilePanel`

- This prepares the codebase for moving FilePanel to a separate package later
- All changes maintain backward compatibility within the internal package
- After this, FilePanel and related types can be moved to `src/pkg/filepanel/` or similar