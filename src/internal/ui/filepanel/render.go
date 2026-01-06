package filepanel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

/*
- TODO: Write File Panel Specific unit test
  - Individual panel resizes
  - Footer content of filepanel changes due to resizing
  - i Only mode icons remains on smaller
  - ii Other things that change too
  - Other panels like clipboard and metadata's content changes too on resize
*/
func (m *Model) Render(focused bool) string {
	r := ui.FilePanelRenderer(m.height, m.width, focused)

	m.renderTopBar(r)
	m.renderSearchBar(r)
	m.renderFooter(r, m.SelectedCount())
	if m.NeedRenderHeaders() {
		m.renderColumnHeaders(r)
	}
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

func (m *Model) renderColumnHeaders(r *rendering.Renderer) {
	var builder strings.Builder
	for _, column := range m.columns {
		builder.WriteString(column.RenderHeader())
	}
	r.AddLines(builder.String())
}

/*
- The renderer for the mandatory first column in the file panel, with a name, a cursor, and a select option.
*/
func (m *Model) renderFileName(columnWidth int) fileElementRender {
	return func(indexElement int) string {
		isSelected := m.CheckSelected(m.Element[indexElement].Location)
		cursor := " "
		if indexElement == m.Cursor && !m.SearchBar.Focused() {
			cursor = icon.Cursor
		}

		selectBox := m.renderSelectBox(isSelected)

		// Calculate the actual prefix width for proper alignment
		prefixWidth := ansi.StringWidth(cursor+" ") + ansi.StringWidth(selectBox)

		isLink := m.Element[indexElement].Info.Mode()&os.ModeSymlink != 0
		renderedName := common.PrettierFilePanelItemName(
			m.Element[indexElement].Name,
			columnWidth-prefixWidth,
			m.Element[indexElement].Directory,
			isLink,
			isSelected,
			common.FilePanelBGColor,
		)
		return common.FilePanelCursorStyle.Render(cursor+" ") + selectBox + renderedName
	}
}

/*
- The renderer of delimiter spaces. It has a strict fixed size that depends only on the delimiter string.
*/
func (m *Model) renderDelimiter(columnWidth int) fileElementRender {
	return func(_ int) string {
		return common.PrettierFixedWidthItem(
			ColumnDelimiter,
			columnWidth,
			false,
			common.FilePanelBGColor,
			lipgloss.Left,
		)
	}
}

/*
- The renderer of a file size column.
*/
func (m *Model) renderFileSize(columnWidth int) fileElementRender {
	return func(indexElement int) string {
		return common.PrettierFixedWidthItem(
			common.FormatFileSize(m.Element[indexElement].Info.Size()),
			columnWidth,
			false,
			common.FilePanelBGColor,
			lipgloss.Right,
		)
	}
}

/*
- The renderer of a modify time column.
TODO: make time template configurable
*/
func (m *Model) renderModifyTime(columnWidth int) fileElementRender {
	return func(indexElement int) string {
		modifyTime := m.Element[indexElement].Info.ModTime().Format("2006-01-02 15:04")
		return common.PrettierFixedWidthItem(
			modifyTime,
			columnWidth,
			false,
			common.FilePanelBGColor,
			lipgloss.Right,
		)
	}
}

/*
- The renderer of a permission column.
*/
func (m *Model) renderPermissions(columnWidth int) fileElementRender {
	return func(indexElement int) string {
		return common.PrettierFixedWidthItem(
			m.Element[indexElement].Info.Mode().Perm().String(),
			columnWidth,
			false,
			common.FilePanelBGColor,
			lipgloss.Right,
		)
	}
}

func (m *Model) renderFileEntries(r *rendering.Renderer) {
	if len(m.Element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}
	end := min(m.RenderIndex+m.PanelElementHeight(), len(m.Element))

	for itemIndex := m.RenderIndex; itemIndex < end; itemIndex++ {
		if m.Renaming && itemIndex == m.Cursor {
			r.AddLines(m.Rename.View())
			continue
		}
		var builder strings.Builder
		for _, column := range m.columns {
			colData := column.GetRenderer()(itemIndex)
			builder.WriteString(colData)
		}
		r.AddLines(builder.String())
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
