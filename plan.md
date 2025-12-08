# MND LINTER REMEDIATION PLAN

## CURRENT STATUS (Updated)

### COMPLETED ✅
**Golines formatting fixes:**
- Fixed all 4 golines issues in src/cmd/main.go, src/internal/function.go, src/internal/handle_file_operations.go, src/internal/model.go
- Extracted lipgloss styles to variables for cleaner fmt.Printf calls
- Build verified after each change

**UI Constants created:**
- Added shared constants: src/internal/common/ui_consts.go (HelpKeyColumnWidth, DefaultCLIContextTimeout, PanelPadding, BorderPadding, InnerPadding, FooterGroupCols, FilePanelMax, MinWidthForRename, ResponsiveWidthThreshold, HeightBreakA–D, ReRenderChunkDivisor, FilePanelWidthUnit, DefaultPreviewTimeout)
- Applied constants to core files (cmd/main.go, model.go, handle_panel_navigation.go, model_render.go, function.go, handle_modal.go, handle_panel_movement.go, preview/model.go)
- Build validated successfully
- 1 commit made with initial changes

### IN PROGRESS 🔄
**Phase 1: String/Size Constants (5 issues)**
- src/internal/common/string_function.go:
  - Lines 119-120: Add KilobyteSize = 1000
  - Lines 123-124: Add KibibyteSize = 1024
  - Line 135: Add TabWidth = 4

**Phase 2: Internal UI (28 issues)**
- src/internal/ui/metadata/const.go: Document keyDataModified: 2
- src/internal/default_config.go: Replace width: 10
- src/internal/file_operations_extract.go: 2 extraction constants
- src/internal/handle_panel_navigation.go: 3 width calculations
- src/internal/model.go: 16 centering /2 operations → CenterDivisor = 2
- src/internal/model_render.go: 8 render dimensions
- src/internal/type_utils.go: 3 instances of +2 → BorderPadding

### TODO 📝
**Phase 3: Image Preview (17 issues)**
- src/pkg/file_preview/image_preview.go: 4 issues (DefaultThumbnailWidth, bit shifts, masks)
- src/pkg/file_preview/image_resize.go: 8 issues (quality levels 2-8)
- src/pkg/file_preview/kitty.go: 3 issues (hash seed, prime, max ID)
- src/pkg/file_preview/utils.go: 2 issues (pixels per column/row)

**Phase 4: Enable Linter**
- Uncomment `- mnd` in .golangci.yaml

## REMAINING MND ISSUES
- golines: 0 ✅ (all fixed)
- mnd: 50 remaining
  - src/internal/common/string_function.go: 5
  - src/internal/ui/metadata/const.go: 1
  - src/internal/default_config.go: 1
  - src/internal/file_operations_extract.go: 2
  - src/internal/handle_panel_navigation.go: 3
  - src/internal/model.go: 16
  - src/internal/model_render.go: 8
  - src/internal/type_utils.go: 3
  - src/pkg/file_preview/: 17 total

## Build Command
CGO_ENABLED=0 go build -o bin/spf ./src/cmd

## Lint Command
golangci-lint run --enable=mnd
