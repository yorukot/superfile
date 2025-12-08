MND LINTER REMEDIATION PLAN (HIGHLY DETAILED)

Objective
- Enable and satisfy the mnd (magic number) linter across the codebase.
- Replace repeated numeric literals with named constants.
- Add targeted //nolint:mnd with concise justification where constants reduce clarity.
- Preserve behavior 100%. No feature changes.

Prerequisites
- Ensure golangci-lint is available (already present: v2.5.0).
- Work on a feature branch (e.g., mnd-remediation).
- Build baseline: CGO_ENABLED=0 go build -o bin/spf ./src/cmd
- Lint baseline: golangci-lint run --enable=mnd

Step 1 — Enable mnd in existing linter config (uncomment only)
1. File to edit: .golangci.yaml (do NOT add .golangci.yml).
2. In the `linters.enable` list, locate the commented mnd entry:
   `# - mnd # detects magic numbers` and UNCOMMENT it to: `- mnd # detects magic numbers`.
   - Location: around lines 100–115 in .golangci.yaml, under `linters: enable:`.
3. Do not change the existing `mnd:` settings block (already present at ~line 356). Keep it as-is.
4. Exclusions: follow the repo’s current convention only if needed.
   - If tests begin flagging noisy mnd hits, append `mnd` to the existing per-path exclusion rule for tests:
     `.golangci.yaml -> exclusions.rules -> - path: '_test\.go' -> linters: [ … add mnd here … ]`.
   - Keep order/format consistent with current style. Do not add new exclusion sections.

Step 2 — Add shared UI/layout constants
1. Create file: src/internal/common/ui_consts.go
2. Add the following constants with comments:
   - HelpKeyColumnWidth = 55 // width of help key column in CLI help
   - DefaultCLIContextTimeout = 5 * time.Second // default CLI context timeout
   - PanelPadding = 3 // rows reserved around file list for borders/header/footer
   - BorderPadding = 2 // rows/cols for outer border frame
   - InnerPadding = 4 // cols for inner content padding (truncate widths)
   - FooterGroupCols = 3 // columns per group in footer layout math
   - FilePanelMax = 10 // max number of file panels supported
   - MinWidthForRename = 18 // minimal width for rename input to render
   - ResponsiveWidthThreshold = 95 // width breakpoint for layout behavior
   - HeightBreakA = 30; HeightBreakB = 35; HeightBreakC = 40; HeightBreakD = 45 // stacked height breakpoints
   - ReRenderChunkDivisor = 100 // divisor for re-render throttling
3. Import time at file top. Ensure package is common.

Step 3 — CLI fixes (src/cmd/main.go)
1. Replace 55 in fmt.Printf("%-*s %s\n", 55, ...) with common.HelpKeyColumnWidth in lines printing colored help entries (approx lines 52, 54, 56, 58, 60).
   - Add above the first print: // use shared help column width (mnd)
2. Replace 5*time.Second in context.WithTimeout(..., 5*time.Second) with common.DefaultCLIContextTimeout (approx line 270).
   - Add comment: // shared CLI timeout (mnd)

Step 4 — Core model/layout (src/internal/model.go)
1. Replace literal 10 with common.FilePanelMax for file panel max check (approx line 237).
2. Replace height threshold literals:
   - if height < 30 → if height < common.HeightBreakA
   - else if height < 35 → else if height < common.HeightBreakB
   - else if height < 40 → else if height < common.HeightBreakC
   - else if height < 45 → else if height < common.HeightBreakD
   - if m.fullHeight > 35 → > common.HeightBreakB
   - if m.fullWidth > 95 → > common.ResponsiveWidthThreshold
   - Add comment near block: // responsive layout breakpoints (mnd)
3. Replace +2 with +common.BorderPadding for:
   - m.fileModel.filePreview.SetHeight(m.mainPanelHeight + 2)
   - Bottom/footer widths where +2 is used (search for SetHeight/SetDimensions with +2).
4. Replace /3, /2 usages for modal sizing:
   - m.promptModal.SetMaxHeight(m.fullHeight / 3)
   - m.promptModal.SetWidth(m.fullWidth / 2)
   - m.zoxideModal.SetMaxHeight(m.fullHeight / 2)
   - m.zoxideModal.SetWidth(m.fullWidth / 2)
   Options:
   - Prefer constants ModalThird=3, ModalHalf=2 in common; OR
   - Keep divisions and add //nolint:mnd with comment “modal uses thirds/halves”.
5. Replace -4 and -3 style paddings using constants when adjusting widths/heights in render/layout helpers:
   - Use common.InnerPadding for -4
   - Use common.PanelPadding for -3
6. Replace m.fileModel.width < 18 with < common.MinWidthForRename (approx line 566). Add comment.
7. Replace reRenderTime := int(float64(len(...))/100) with /common.ReRenderChunkDivisor (approx line 731). Add comment.

Step 5 — Panel navigation & operations
- src/internal/handle_panel_navigation.go: replace +2 with +common.BorderPadding when setting preview height (approx line 111). Add comment.
- src/internal/handle_file_operations.go: replace width-4 with width-common.InnerPadding when creating rename input (approx line 100). Add comment.

Step 6 — Rendering (src/internal/model_render.go)
1. Replace +2 with +common.BorderPadding at:
   - FilePanelRenderer(mainPanelHeight+2, filePanelWidth+2, …) (approx line 65)
   - ClipboardRenderer(m.footerHeight+2, bottomWidth+2) (approx line 217)
2. Replace filePanelWidth-4 with filePanelWidth-common.InnerPadding (approx line 77)
3. Replace bottom width calc utils.FooterWidth(m.fullWidth + m.fullWidth%3 + 2):
   - Use %common.FooterGroupCols
   - Replace +2 with +common.BorderPadding (approx line 213)
4. Replace -3 when truncating display width within metadata/preview draws with -common.PanelPadding (approx line 236). Add comment.
5. Replace ModalWidth-4 with ModalWidth-common.InnerPadding (approx line 297).
6. Replace panel.sortOptions.width-2 with -common.BorderPadding (approx line 457).

Step 7 — Directory listing & sorting (src/internal/function.go)
1. func panelElementHeight(mainPanelHeight int) int { return mainPanelHeight - 3 }
   - Replace 3 with common.PanelPadding (approx line 244). Add comment.
2. In suffixRegexp.FindStringSubmatch(name); len(match) == 3 (approx line 294):
   - Keep 3 and add inline: //nolint:mnd — 3 = full match + 2 capture groups (regex)

Step 8 — Preview subsystem
- src/internal/ui/preview/model.go: replace 500*time.Millisecond with a named constant.
  Options:
  - Add in common: DefaultPreviewTimeout = 500*time.Millisecond
  - Or local const in preview package with comment.
  Add comment: // preview operation timeout (mnd)

Step 9 — Image preview utils (src/pkg/file_preview/image_preview.go)
1. Replace 100 with DefaultThumbnailWidth (const in this file). Comment: // default thumb width (mnd)
2. Replace 5*time.Minute with DefaultCacheExpiration. Comment.
3. Replace /2 on ticker with either const DividerTwo or //nolint:mnd “half expiration interval”.
4. Replace 16, 8, 0xFF, 255 with named consts: Shift16, Shift8, MaskFF, OpaqueAlpha; add comment block explaining RGB math.
5. Replace 8 in fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), …) with Shift8.

Step 10 — Image resize (src/pkg/file_preview/image_resize.go)
1. Switch cases 2..8:
   - Prefer //nolint:mnd on each case with comment: // 2=low … 8=ultra quality levels
   - Or introduce Quality2..Quality8 if reused elsewhere.
2. imaging.Fit(img, maxWidth, maxHeight*2, …): replace 2 with HeightScaleFactor const or //nolint:mnd with comment “fit scales height x2”.

Step 11 — Kitty utils (src/pkg/file_preview/kitty.go)
1. Replace 42 with KittyDefaultSeed
2. Replace 31 with HashPrime
3. Replace 0xFFFF with MaskFFFF
4. Replace +1000 with NonZeroOffset
5. Add one-line comments next to const block.

Step 12 — Preview terminal metrics (src/pkg/file_preview/utils.go)
1. Replace PixelsPerColumn: 8, PixelsPerRow: 16 with named consts (PixelsPerColumnDefault, PixelsPerRowDefault) and comment on typical terminal cell sizes.

Step 13 — External path detection (optional cleanup)
- src/internal/function.go:isExternalDiskPath
  - Create consts: TimeMachinePrefix, VolumesPrefix, MediaPrefixes (slice)
  - Use them in HasPrefix checks. Rationale: clarity; not necessarily for mnd.

Step 14 — Commenting & //nolint practices
- Only add //nolint:mnd where constants reduce clarity or are inherently part of API/math:
  - Regex submatch count (len(match)==3): add concise reason
  - Switch cases for fixed, human-defined quality levels (2..8): add mapping comment
  - Simple halves/thirds if not centralized: add “half/third sizing” comments
- For every //nolint:mnd, add a short, explicit justification on the same line.

Validation Checklist
1. golangci-lint run --enable=mnd — should show a decreasing count; iterate until 0 or only justified //nolint sites remain.
2. Build: CGO_ENABLED=0 go build -o bin/spf ./src/cmd
3. Tests: run focused suites to ensure no behavior change
   - go test ./src/internal -run '^TestInitialFilePathPositionsCursor|TestInitialFilePathPositionsCursorWindow$'
   - go test ./src/internal -run '^TestReturnDirElement$'
4. Manual smoke:
   - Launch spf in a directory with many files; verify layout unchanged.
   - Open with a file path; ensure cursor targets file and remains visible.

Commit Strategy
- Commit 1: Add ui_consts.go and .golangci.yml mnd enablement.
- Commit 2: Apply constants to src/cmd and internal model/layout.
- Commit 3: Rendering replacements.
- Commit 4: Preview and image utils (with //nolint where needed).
- Commit 5: Kitty utils and optional path-prefix constants.
- Keep each commit small and scoped; include brief messages referencing mnd.
