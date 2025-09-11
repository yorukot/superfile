package internal

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/internal/utils"
)

// ================ File Panel Core Methods ================

func (panel *FilePanel) GetSelectedItem() Element {
	if panel.Cursor < 0 || len(panel.Element) <= panel.Cursor {
		return Element{}
	}
	return panel.Element[panel.Cursor]
}

func (panel *FilePanel) ResetSelected() {
	panel.Selected = panel.Selected[:0]
}

func (panel *FilePanel) GetSelectedItemPtr() *Element {
	if panel.Cursor < 0 || len(panel.Element) <= panel.Cursor {
		return nil
	}
	return &panel.Element[panel.Cursor]
}

func (panel *FilePanel) ChangeFilePanelMode() {
	switch panel.PanelMode {
	case SelectMode:
		panel.Selected = panel.Selected[:0]
		panel.PanelMode = BrowserMode
	case BrowserMode:
		panel.PanelMode = SelectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", panel.PanelMode)
	}
}

func (panel *FilePanel) UpdateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", panel.Location, "path", path)
	path = utils.ResolveAbsPath(panel.Location, path)

	if path == panel.Location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	panel.DirectoryRecords[panel.Location] = DirectoryRecord{
		directoryCursor: panel.Cursor,
		directoryRender: panel.RenderIndex,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	panel.Location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := panel.DirectoryRecords[panel.Location]
	if hasRecord {
		panel.Cursor = curDirectoryRecord.directoryCursor
		panel.RenderIndex = curDirectoryRecord.directoryRender
	} else {
		panel.Cursor = 0
		panel.RenderIndex = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", panel.Cursor, "render", panel.RenderIndex)

	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	panel.SearchBar.SetValue("")

	return nil
}

func (panel *FilePanel) ParentDirectory() error {
	return panel.UpdateCurrentFilePanelDir("..")
}

// ================ Navigation Methods ================

func (panel *FilePanel) ListUp(mainPanelHeight int) {
	if len(panel.Element) == 0 {
		return
	}
	if panel.Cursor > 0 {
		panel.Cursor--
		if panel.Cursor < panel.RenderIndex {
			panel.RenderIndex--
		}
	} else {
		if len(panel.Element) > panelElementHeight(mainPanelHeight) {
			panel.RenderIndex = len(panel.Element) - panelElementHeight(mainPanelHeight)
			panel.Cursor = len(panel.Element) - 1
		} else {
			panel.Cursor = len(panel.Element) - 1
		}
	}
}

func (panel *FilePanel) ListDown(mainPanelHeight int) {
	if len(panel.Element) == 0 {
		return
	}
	if panel.Cursor < len(panel.Element)-1 {
		panel.Cursor++
		if panel.Cursor > panel.RenderIndex+panelElementHeight(mainPanelHeight)-1 {
			panel.RenderIndex++
		}
	} else {
		panel.RenderIndex = 0
		panel.Cursor = 0
	}
}

func (panel *FilePanel) PgUp(mainPanelHeight int) {
	panlen := len(panel.Element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2

	if panHeight >= panlen {
		panel.Cursor = 0
	} else {
		if panel.Cursor-panHeight <= 0 {
			panel.Cursor = 0
			panel.RenderIndex = 0
		} else {
			panel.Cursor -= panHeight
			panel.RenderIndex = panel.Cursor - panCenter

			if panel.RenderIndex < 0 {
				panel.RenderIndex = 0
			}
		}
	}
}

func (panel *FilePanel) PgDown(mainPanelHeight int) {
	panlen := len(panel.Element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2

	if panHeight >= panlen {
		panel.Cursor = panlen - 1
	} else {
		if panel.Cursor+panHeight >= panlen {
			panel.Cursor = panlen - 1
			panel.RenderIndex = panel.Cursor - panCenter
		} else {
			panel.Cursor += panHeight
			panel.RenderIndex = panel.Cursor - panCenter
		}
	}
}

func (panel *FilePanel) ItemSelectUp(mainPanelHeight int) {
	panel.SingleItemSelect()
	panel.ListUp(mainPanelHeight)
}

func (panel *FilePanel) ItemSelectDown(mainPanelHeight int) {
	panel.SingleItemSelect()
	panel.ListDown(mainPanelHeight)
}

// ================ Selection Methods ================

func (panel *FilePanel) SingleItemSelect() {
	if len(panel.Element) > 0 && panel.Cursor >= 0 && panel.Cursor < len(panel.Element) {
		elementLocation := panel.Element[panel.Cursor].location

		if arrayContains(panel.Selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			panel.Selected = removeElementByValue(panel.Selected, elementLocation)
		} else {
			panel.Selected = append(panel.Selected, elementLocation)
		}
	}
}

// ================ Rendering Methods ================

func (panel *FilePanel) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+2, filePanelWidth+2, focussed)

	panel.renderTopBar(r, filePanelWidth)
	panel.renderSearchBar(r)
	panel.renderFooter(r)
	panel.renderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}

func (panel *FilePanel) renderTopBar(r *rendering.Renderer, filePanelWidth int) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(panel.Location, filePanelWidth-4, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (panel *FilePanel) renderSearchBar(r *rendering.Renderer) {
	r.AddLines(" " + panel.SearchBar.View())
}

func (panel *FilePanel) renderFooter(r *rendering.Renderer) {
	sortLabel, sortIcon := panel.getSortInfo()
	modeLabel, modeIcon := panel.getPanelModeInfo()
	cursorStr := panel.getCursorString()

	if common.Config.Nerdfont {
		sortLabel = sortIcon + icon.Space + sortLabel
		modeLabel = modeIcon + icon.Space + modeLabel
	} else {
		// TODO : Figure out if we can set icon.Space to " " if nerdfont is false
		// That would simplify code
		sortLabel = sortIcon + " " + sortLabel
	}

	if common.Config.ShowPanelFooterInfo {
		r.SetBorderInfoItems(sortLabel, modeLabel, cursorStr)
		if r.AreInfoItemsTruncated() {
			r.SetBorderInfoItems(sortIcon, modeIcon, cursorStr)
		}
	} else {
		r.SetBorderInfoItems(cursorStr)
	}
}

func (panel *FilePanel) renderFileEntries(r *rendering.Renderer, mainPanelHeight, filePanelWidth int) {
	if len(panel.Element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}

	end := min(panel.RenderIndex+panelElementHeight(mainPanelHeight), len(panel.Element))

	for i := panel.RenderIndex; i < end; i++ {
		// TODO : Fix this, this is O(n^2) complexity. Considered a file panel with 200 files, and 100 selected
		// We will be doing a search in 100 item slice for all 200 files.
		isSelected := arrayContains(panel.Selected, panel.Element[i].location)

		if panel.Renaming && i == panel.Cursor {
			r.AddLines(panel.Rename.View())
			continue
		}

		cursor := " "
		if i == panel.Cursor && !panel.SearchBar.Focused() {
			cursor = icon.Cursor
		}

		// Performance TODO: Remove or cache this if not needed at render time
		// This will unnecessarily slow down rendering. There should be a way to avoid this at render
		_, err := os.ReadDir(panel.Element[i].location)
		dirExists := err == nil || panel.Element[i].directory

		renderedName := common.PrettierName(
			panel.Element[i].name,
			filePanelWidth-5,
			dirExists,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + renderedName)
	}
}

func (panel *FilePanel) getSortInfo() (string, string) {
	opts := panel.SortOptions.data
	selected := opts.options[opts.selected]
	label := selected
	if selected == string(sortingDateModified) {
		label = "Date"
	}

	iconStr := icon.SortAsc
	if opts.reversed {
		iconStr = icon.SortDesc
	}

	return label, iconStr
}

func (panel *FilePanel) getPanelModeInfo() (string, string) {
	label := "Browser"
	iconStr := icon.Browser

	if panel.PanelMode == SelectMode {
		label = "Select"
		iconStr = icon.Select
	}

	return label, iconStr
}

func (panel *FilePanel) getCursorString() string {
	if len(panel.Element) == 0 {
		return "0/0"
	}
	return fmt.Sprintf("%d/%d", panel.Cursor+1, len(panel.Element))
}

// ================ Accessor Methods ================

func (panel *FilePanel) GetLocation() string {
	return panel.Location
}

func (panel *FilePanel) GetCursor() int {
	return panel.Cursor
}

func (panel *FilePanel) SetCursor(cursor int) {
	panel.Cursor = cursor
}

func (panel *FilePanel) GetRender() int {
	return panel.RenderIndex
}

func (panel *FilePanel) SetRender(render int) {
	panel.RenderIndex = render
}

func (panel *FilePanel) GetElement() []Element {
	return panel.Element
}

func (panel *FilePanel) SetElement(elements []Element) {
	panel.Element = elements
}

func (panel *FilePanel) GetSelected() []string {
	return panel.Selected
}

func (panel *FilePanel) SetSelected(selected []string) {
	panel.Selected = selected
}

func (panel *FilePanel) GetPanelMode() PanelMode {
	return panel.PanelMode
}

func (panel *FilePanel) SetPanelMode(mode PanelMode) {
	panel.PanelMode = mode
}

func (panel *FilePanel) IsFocused() bool {
	return panel.isFocused
}

func (panel *FilePanel) SetFocused(focused bool) {
	panel.isFocused = focused
}

func (panel *FilePanel) IsRenaming() bool {
	return panel.Renaming
}

func (panel *FilePanel) SetRenaming(renaming bool) {
	panel.Renaming = renaming
}
