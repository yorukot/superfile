package filepanel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (panel *FilePanel) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+common.BorderPadding, filePanelWidth+common.BorderPadding, focussed)

	panel.RenderTopBar(r, filePanelWidth)
	panel.RenderSearchBar(r)
	panel.RenderFooter(r, uint(len(panel.Selected)))
	panel.RenderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}

func (panel *FilePanel) RenderTopBar(r *rendering.Renderer, filePanelWidth int) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(panel.Location, filePanelWidth-common.InnerPadding, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (panel *FilePanel) RenderSearchBar(r *rendering.Renderer) {
	r.AddLines(" " + panel.SearchBar.View())
}

// TODO : Unit test this
func (panel *FilePanel) RenderFooter(r *rendering.Renderer, selectedCount uint) {
	sortLabel, sortIcon := panel.GetSortInfo()
	modeLabel, modeIcon := panel.GetPanelModeInfo(selectedCount)
	cursorStr := panel.GetCursorString()

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

func (panel *FilePanel) RenderFileEntries(r *rendering.Renderer, mainPanelHeight, filePanelWidth int) {
	if len(panel.Element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}

	end := min(panel.RenderIndex+panelElementHeight(mainPanelHeight), len(panel.Element))

	for i := panel.RenderIndex; i < end; i++ {
		// TODO : Fix this, this is O(n^2) complexity. Considered a file panel with 200 files, and 100 selected
		// We will be doing a search in 100 item slice for all 200 files.
		isSelected := slices.Contains(panel.Selected, panel.Element[i].Location)

		if panel.Renaming && i == panel.Cursor {
			r.AddLines(panel.Rename.View())
			continue
		}

		cursor := " "
		if i == panel.Cursor && !panel.SearchBar.Focused() {
			cursor = icon.Cursor
		}

		selectBox := panel.RenderSelectBox(isSelected)

		// Calculate the actual prefix width for proper alignment
		prefixWidth := lipgloss.Width(cursor+" ") + lipgloss.Width(selectBox)

		isLink := panel.Element[i].Info.Mode()&os.ModeSymlink != 0
		renderedName := common.PrettierName(
			panel.Element[i].Name,
			filePanelWidth-prefixWidth,
			panel.Element[i].Directory,
			isLink,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + selectBox + renderedName)
	}
}

func (panel *FilePanel) GetSortInfo() (string, string) {
	opts := panel.SortOptions.Data
	selected := opts.Options[opts.Selected]
	label := selected
	if selected == string(sortingDateModified) {
		label = "Date"
	}

	iconStr := icon.SortAsc

	if opts.Reversed {
		iconStr = icon.SortDesc
	}
	return label, iconStr
}

func (panel *FilePanel) GetPanelModeInfo(selectedCount uint) (string, string) {
	switch panel.PanelMode {
	case BrowserMode:
		return "Browser", icon.Browser
	case SelectMode:
		return "Select" + icon.Space + fmt.Sprintf("(%d)", selectedCount), icon.Select
	default:
		return "", ""
	}
}

func (panel *FilePanel) GetCursorString() string {
	cursor := panel.Cursor
	if len(panel.Element) > 0 {
		cursor++ // Convert to 1-based
	}
	return fmt.Sprintf("%d/%d", cursor, len(panel.Element))
}

func (panel *FilePanel) RenderSelectBox(isSelected bool) string {
	if !common.Config.ShowSelectIcons || !common.Config.Nerdfont || panel.PanelMode != SelectMode {
		return ""
	}

	if panel.IsFocused {
		if isSelected {
			return common.CheckboxCheckedFocused
		}
		return common.CheckboxEmptyFocused
	}
	if isSelected {
		return common.CheckboxChecked
	}
	return common.CheckboxEmpty
}

// Checks whether the focus panel directory changed and forces a re-render.
func (panel *FilePanel) NeedsReRender() bool {
	if len(panel.Element) > 0 {
		return filepath.Dir(panel.Element[0].Location) != panel.Location
	}
	return true
}
