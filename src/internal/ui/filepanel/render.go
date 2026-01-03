package filepanel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) makeColumns(columnThreshold int) []columnDefinition {
	extraColumns := []columnDefinition{
		{Name: "Size", Generator: m.renderFileSize, Size: FileSizeColumnWidth},
	}
	maxColumns := len(extraColumns)
	if columnThreshold < maxColumns {
		maxColumns = columnThreshold
	}

	columns := []columnDefinition{
		{Name: "Name", Generator: m.renderFileName, Size: m.GetContentWidth()},
	}
	// "-1" guards in a cases of rounding numbers.
	extraColumnsThreshold := int(float64(m.GetContentWidth())*FileNameRatio - 1)
	remainWidth := m.GetContentWidth()

	for _, col := range extraColumns[0:maxColumns] {
		widthExtraColumn := lipgloss.Width(ColumnDelimiter) + col.Size
		if remainWidth-widthExtraColumn > extraColumnsThreshold {
			delimiterCol := columnDefinition{
				Name:      "",
				Generator: m.renderDelimiter,
				Size:      lipgloss.Width(ColumnDelimiter),
			}
			columns = append(columns, delimiterCol, col)
			columns[0].Size -= widthExtraColumn
			remainWidth -= widthExtraColumn
		} else {
			break
		}
	}
	return columns
}

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
	columns := m.makeColumns(common.Config.FilePanelExtraColumns)
	if common.Config.FilePanelExtraColumns > 0 {
		m.renderColumnHeaders(r, columns)
	}
	m.renderFileEntries(r, columns)
	return r.Render()
}

func (m *Model) renderTopBar(r *rendering.Renderer) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(m.Location, m.GetContentWidth()-common.InnerPadding, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (m *Model) renderColumnHeaders(r *rendering.Renderer, columns []columnDefinition) {
	var builder strings.Builder
	for _, column := range columns {
		builder.WriteString(column.RenderHeader())
	}
	r.AddLines(builder.String())
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

/*
- The renderer for the mandatory first column in the file panel, with a name, a cursor, and a select option.
*/
func (m *Model) renderFileName(columnWidth int) fileElementRender {
	return func(index int, item Element) string {
		isSelected := m.CheckSelected(m.Element[index].Location)
		cursor := " "
		if index == m.Cursor && !m.SearchBar.Focused() {
			cursor = icon.Cursor
		}

		selectBox := m.renderSelectBox(isSelected)

		// Calculate the actual prefix width for proper alignment
		prefixWidth := lipgloss.Width(cursor+" ") + lipgloss.Width(selectBox)

		isLink := m.Element[index].Info.Mode()&os.ModeSymlink != 0
		renderedName := common.PrettierFileName(
			m.Element[index].Name,
			columnWidth-prefixWidth,
			m.Element[index].Directory,
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
	return func(index int, item Element) string {
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
	return func(index int, item Element) string {
		return common.PrettierFixedWidthItem(
			common.FormatFileSize(item.Info.Size()),
			columnWidth,
			false,
			common.FilePanelBGColor,
			lipgloss.Right,
		)
	}
}

func (m *Model) renderFileEntries(r *rendering.Renderer, columns []columnDefinition) {
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
		for _, column := range columns {
			colData := column.GetRenderer()(itemIndex, m.Element[itemIndex])
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
