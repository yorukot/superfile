package filepanel

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) Render(focussed bool) string {
	r := ui.FilePanelRenderer(m.height, m.width, focussed)

	m.renderTopBar(r)
	m.renderSearchBar(r)
	m.renderFooter(r, uint(len(m.Selected)))
	m.renderFileEntries(r)

	return r.Render()
}

func (m *Model) renderTopBar(r *rendering.Renderer) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(m.Location, m.GetContentWidth()-common.InnerPadding, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (m *Model) renderSearchBar(r *rendering.Renderer) {
	r.AddLines(" " + m.SearchBar.View())
}

// TODO : Unit test this
func (m *Model) renderFooter(r *rendering.Renderer, selectedCount uint) {
	sortLabel, sortIcon := m.getSortInfo()
	modeLabel, modeIcon := m.getPanelModeInfo(selectedCount)
	cursorStr := m.getCursorString()

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

func (m *Model) renderFileEntries(r *rendering.Renderer) {
	if len(m.Element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}

	end := min(m.RenderIndex+m.PanelElementHeight(), len(m.Element))

	for i := m.RenderIndex; i < end; i++ {
		// TODO : Fix this, this is O(n^2) complexity. Considered a file panel with 200 files, and 100 selected
		// We will be doing a search in 100 item slice for all 200 files.
		isSelected := slices.Contains(m.Selected, m.Element[i].Location)

		if m.Renaming && i == m.Cursor {
			r.AddLines(m.Rename.View())
			continue
		}

		cursor := " "
		if i == m.Cursor && !m.SearchBar.Focused() {
			cursor = icon.Cursor
		}

		selectBox := m.renderSelectBox(isSelected)

		// Calculate the actual prefix width for proper alignment
		prefixWidth := lipgloss.Width(cursor+" ") + lipgloss.Width(selectBox)

		isLink := m.Element[i].Info.Mode()&os.ModeSymlink != 0
		renderedName := common.PrettierName(
			m.Element[i].Name,
			m.GetContentWidth()-prefixWidth,
			m.Element[i].Directory,
			isLink,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + selectBox + renderedName)
	}
}

func (m *Model) getSortInfo() (string, string) {
	opts := m.SortOptions.Data
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

func (m *Model) getPanelModeInfo(selectedCount uint) (string, string) {
	switch m.PanelMode {
	case BrowserMode:
		return "Browser", icon.Browser
	case SelectMode:
		return "Select" + icon.Space + fmt.Sprintf("(%d)", selectedCount), icon.Select
	default:
		return "", ""
	}
}

func (m *Model) getCursorString() string {
	cursor := m.Cursor
	if len(m.Element) > 0 {
		cursor++ // Convert to 1-based
	}
	return fmt.Sprintf("%d/%d", cursor, len(m.Element))
}

func (m *Model) renderSelectBox(isSelected bool) string {
	if !common.Config.ShowSelectIcons || !common.Config.Nerdfont || m.PanelMode != SelectMode {
		return ""
	}

	if m.IsFocused {
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
func (m *Model) NeedsReRender() bool {
	if len(m.Element) > 0 {
		return filepath.Dir(m.Element[0].Location) != m.Location
	}
	return true
}
