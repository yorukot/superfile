# MND LINTER REMEDIATION PLAN

## CURRENT STATUS (Updated)

### COMPLETED ✅
**Golines formatting fixes:**
- Fixed all 4 golines issues in src/cmd/main.go, src/internal/function.go, src/internal/handle_file_operations.go, src/internal/model.go
- Extracted lipgloss styles to variables for cleaner fmt.Printf calls
- Build verified after each change

**Phase 1: String/Size Constants:**
- Added KilobyteSize=1000, KibibyteSize=1024, TabWidth=4 constants
- Applied in FormatFileSize and CheckAndTruncateLineLengths functions
- Committed separately

**Phase 2: Internal UI (partial):**
- Added DefaultFilePanelWidth, ExtractedFileMode/DirMode, CenterDivisor constants  
- Fixed centering calculations in model.go (14 /2 operations)
- Updated panel width calculations in handle_panel_navigation.go and model_render.go
- Replaced hardcoded padding values with constants
- Committed separately

**UI Constants created:**
- Added shared constants: src/internal/common/ui_consts.go (HelpKeyColumnWidth, DefaultCLIContextTimeout, PanelPadding, BorderPadding, InnerPadding, FooterGroupCols, FilePanelMax, MinWidthForRename, ResponsiveWidthThreshold, HeightBreakA–D, ReRenderChunkDivisor, FilePanelWidthUnit, DefaultPreviewTimeout)
- Applied constants to core files (cmd/main.go, model.go, handle_panel_navigation.go, model_render.go, function.go, handle_modal.go, handle_panel_movement.go, preview/model.go)
- Build validated successfully

### IN PROGRESS 🔄
**Remaining Internal UI Issues (32 total):**
- src/internal/common/string_function.go: 5 (buffer size check)
- src/internal/common/style_function.go: 2
- src/internal/type_utils.go: 3
- src/internal/ui/: 15 issues across metadata, rendering, sidebar, zoxide
- src/internal/utils/ui_utils.go: 2
- src/internal/model_render.go: 1

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
- mnd: 49 remaining across internal/ and pkg/file_preview/

## Build Command
CGO_ENABLED=0 go build -o bin/spf ./src/cmd

## Lint Command
golangci-lint run --enable=mnd
