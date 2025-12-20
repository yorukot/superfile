# FilePanel Export Refactoring Plan

## Phase 1: Disable GitHub Workflows
1. Remove or rename `.github/workflows/` directory to `.github/workflows.disabled/`
2. This will prevent workflows from running on push to save GitHub CPU

## Phase 2: Export Core Types (src/internal/type.go)

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

## Phase 3: Export FilePanel Methods (src/internal/file_panel.go)

All methods should be exported (capitalize first letter):
- `getSelectedItem()` → `GetSelectedItem()`
- `resetSelected()` → `ResetSelected()`
- `getSelectedItemPtr()` → `GetSelectedItemPtr()`
- `changeFilePanelMode()` → `ChangeFilePanelMode()`
- `updateCurrentFilePanelDir()` → `UpdateCurrentFilePanelDir()`
- `parentDirectory()` → `ParentDirectory()`
- `handleResize()` → `HandleResize()`
- Any other panel methods found

## Phase 4: Update fileModel struct (src/internal/type.go)
- Update `filePanels []filePanel` → `filePanels []FilePanel`

## Phase 5: Update All References

### 5.1 Update type references
- All `filePanel` → `FilePanel`
- All `element` → `Element`
- All `directoryRecord` → `DirectoryRecord`
- All `sortOptionsModel` → `SortOptionsModel`
- All `panelMode` → `PanelMode`

### 5.2 Update field access (throughout codebase)
Examples:
- `panel.cursor` → `panel.Cursor`
- `panel.location` → `panel.Location`
- `panel.element` → `panel.Element`
- `elem.name` → `elem.Name`
- `elem.directory` → `elem.Directory`
- etc.

### 5.3 Update method calls
- `panel.getSelectedItem()` → `panel.GetSelectedItem()`
- `panel.updateCurrentFilePanelDir()` → `panel.UpdateCurrentFilePanelDir()`
- etc.

## Phase 6: Files to Update
Main files that will need updates:
1. `src/internal/type.go` - Type definitions
2. `src/internal/file_panel.go` - All panel methods
3. `src/internal/model.go` - Uses filePanel extensively
4. `src/internal/model_render.go` - Renders panel fields
5. `src/internal/handle_panel_movement.go` - Panel navigation
6. `src/internal/handle_panel_navigation.go` - Panel operations
7. `src/internal/handle_file_operations.go` - File operations on panels
8. `src/internal/function.go` - Helper functions using element
9. `src/internal/type_utils.go` - Panel initialization
10. Test files:
    - `src/internal/model_test.go`
    - `src/internal/model_navigation_test.go`
    - `src/internal/function_test.go`
    - etc.

## Phase 7: Compilation & Testing
1. Run `go build -o bin/spf ./src/cmd` after each major phase
2. Fix any compilation errors
3. Run tests: `go test ./...`

## Implementation Order
1. Disable workflows (Phase 1)
2. Export types (Phase 2) - Single commit
3. Export methods (Phase 3) - Single commit
4. Update all references (Phase 4-5) - Can be multiple commits
5. Ensure compilation (Phase 6-7)

## Notes
- This prepares the codebase for moving FilePanel to a separate package later
- All changes maintain backward compatibility within the internal package
- After this, FilePanel and related types can be moved to `src/pkg/filepanel/` or similar