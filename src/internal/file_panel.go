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

func (panel *filePanel) GetSelectedItem() element {
	if panel.cursor < 0 || len(panel.element) <= panel.cursor {
		return element{}
	}
	return panel.element[panel.cursor]
}

func (panel *filePanel) ResetSelected() {
	panel.selected = panel.selected[:0]
}

func (panel *filePanel) GetSelectedItemPtr() *element {
	if panel.cursor < 0 || len(panel.element) <= panel.cursor {
		return nil
	}
	return &panel.element[panel.cursor]
}

func (panel *filePanel) ChangeFilePanelMode() {
	switch panel.panelMode {
	case selectMode:
		panel.selected = panel.selected[:0]
		panel.panelMode = browserMode
	case browserMode:
		panel.panelMode = selectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", panel.panelMode)
	}
}

func (panel *filePanel) UpdateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", panel.location, "path", path)
	path = utils.ResolveAbsPath(panel.location, path)

	if path == panel.location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	panel.directoryRecords[panel.location] = directoryRecord{
		directoryCursor: panel.cursor,
		directoryRender: panel.render,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	panel.location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := panel.directoryRecords[panel.location]
	if hasRecord {
		panel.cursor = curDirectoryRecord.directoryCursor
		panel.render = curDirectoryRecord.directoryRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", panel.cursor, "render", panel.render)

	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	panel.searchBar.SetValue("")

	return nil
}

func (panel *filePanel) ParentDirectory() error {
	return panel.UpdateCurrentFilePanelDir("..")
}

// ================ Navigation Methods ================

func (panel *filePanel) ListUp(mainPanelHeight int) {
	if len(panel.element) == 0 {
		return
	}
	if panel.cursor > 0 {
		panel.cursor--
		if panel.cursor < panel.render {
			panel.render--
		}
	} else {
		if len(panel.element) > panelElementHeight(mainPanelHeight) {
			panel.render = len(panel.element) - panelElementHeight(mainPanelHeight)
			panel.cursor = len(panel.element) - 1
		} else {
			panel.cursor = len(panel.element) - 1
		}
	}
}

func (panel *filePanel) ListDown(mainPanelHeight int) {
	if len(panel.element) == 0 {
		return
	}
	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+panelElementHeight(mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
}

func (panel *filePanel) PgUp(mainPanelHeight int) {
	panlen := len(panel.element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2

	if panHeight >= panlen {
		panel.cursor = 0
	} else {
		if panel.cursor-panHeight <= 0 {
			panel.cursor = 0
			panel.render = 0
		} else {
			panel.cursor -= panHeight
			panel.render = panel.cursor - panCenter

			if panel.render < 0 {
				panel.render = 0
			}
		}
	}
}

func (panel *filePanel) PgDown(mainPanelHeight int) {
	panlen := len(panel.element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2

	if panHeight >= panlen {
		panel.cursor = panlen - 1
	} else {
		if panel.cursor+panHeight >= panlen {
			panel.cursor = panlen - 1
			panel.render = panel.cursor - panCenter
		} else {
			panel.cursor += panHeight
			panel.render = panel.cursor - panCenter
		}
	}
}

func (panel *filePanel) ItemSelectUp(mainPanelHeight int) {
	panel.SingleItemSelect()
	panel.ListUp(mainPanelHeight)
}

func (panel *filePanel) ItemSelectDown(mainPanelHeight int) {
	panel.SingleItemSelect()
	panel.ListDown(mainPanelHeight)
}

// ================ Selection Methods ================

func (panel *filePanel) SingleItemSelect() {
	if len(panel.element) > 0 && panel.cursor >= 0 && panel.cursor < len(panel.element) {
		elementLocation := panel.element[panel.cursor].location

		if arrayContains(panel.selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			panel.selected = removeElementByValue(panel.selected, elementLocation)
		} else {
			panel.selected = append(panel.selected, elementLocation)
		}
	}
}

// ================ Rendering Methods ================

func (panel *filePanel) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+2, filePanelWidth+2, focussed)

	panel.renderTopBar(r, filePanelWidth)
	panel.renderSearchBar(r)
	panel.renderFooter(r)
	panel.renderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}

func (panel *filePanel) renderTopBar(r *rendering.Renderer, filePanelWidth int) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(panel.location, filePanelWidth-4, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (panel *filePanel) renderSearchBar(r *rendering.Renderer) {
	r.AddLines(" " + panel.searchBar.View())
}

func (panel *filePanel) renderFooter(r *rendering.Renderer) {
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

func (panel *filePanel) renderFileEntries(r *rendering.Renderer, mainPanelHeight, filePanelWidth int) {
	if len(panel.element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}

	end := min(panel.render+panelElementHeight(mainPanelHeight), len(panel.element))

	for i := panel.render; i < end; i++ {
		// TODO : Fix this, this is O(n^2) complexity. Considered a file panel with 200 files, and 100 selected
		// We will be doing a search in 100 item slice for all 200 files.
		isSelected := arrayContains(panel.selected, panel.element[i].location)

		if panel.renaming && i == panel.cursor {
			r.AddLines(panel.rename.View())
			continue
		}

		cursor := " "
		if i == panel.cursor && !panel.searchBar.Focused() {
			cursor = icon.Cursor
		}

		// Performance TODO: Remove or cache this if not needed at render time
		// This will unnecessarily slow down rendering. There should be a way to avoid this at render
		_, err := os.ReadDir(panel.element[i].location)
		dirExists := err == nil || panel.element[i].directory

		renderedName := common.PrettierName(
			panel.element[i].name,
			filePanelWidth-5,
			dirExists,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + renderedName)
	}
}

func (panel *filePanel) getSortInfo() (string, string) {
	opts := panel.sortOptions.data
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

func (panel *filePanel) getPanelModeInfo() (string, string) {
	label := "Browser"
	iconStr := icon.Browser

	if panel.panelMode == selectMode {
		label = "Select"
		iconStr = icon.Select
	}

	return label, iconStr
}

func (panel *filePanel) getCursorString() string {
	if len(panel.element) == 0 {
		return "0/0"
	}
	return fmt.Sprintf("%d/%d", panel.cursor+1, len(panel.element))
}

// ================ Accessor Methods ================

func (panel *filePanel) GetLocation() string {
	return panel.location
}

func (panel *filePanel) GetCursor() int {
	return panel.cursor
}

func (panel *filePanel) SetCursor(cursor int) {
	panel.cursor = cursor
}

func (panel *filePanel) GetRender() int {
	return panel.render
}

func (panel *filePanel) SetRender(render int) {
	panel.render = render
}

func (panel *filePanel) GetElement() []element {
	return panel.element
}

func (panel *filePanel) SetElement(elements []element) {
	panel.element = elements
}

func (panel *filePanel) GetSelected() []string {
	return panel.selected
}

func (panel *filePanel) SetSelected(selected []string) {
	panel.selected = selected
}

func (panel *filePanel) GetPanelMode() panelMode {
	return panel.panelMode
}

func (panel *filePanel) SetPanelMode(mode panelMode) {
	panel.panelMode = mode
}

func (panel *filePanel) IsFocused() bool {
	return panel.isFocused
}

func (panel *filePanel) SetFocused(focused bool) {
	panel.isFocused = focused
}

func (panel *filePanel) IsRenaming() bool {
	return panel.renaming
}

func (panel *filePanel) SetRenaming(renaming bool) {
	panel.renaming = renaming
}
