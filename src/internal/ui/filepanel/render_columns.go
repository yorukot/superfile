package filepanel

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

/*
- The renderer for the mandatory first column in the file panel, with a name, a cursor, and a select option.
*/
func (m *Model) renderFileName(indexElement int, columnWidth int) string {
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

/*
- The renderer of delimiter spaces. It has a strict fixed size that depends only on the delimiter string.
*/
func (m *Model) renderDelimiter(indexElement int, columnWidth int) string {
	return common.PrettierFixedWidthItem(
		ColumnDelimiter,
		columnWidth,
		false,
		common.FilePanelBGColor,
		lipgloss.Left,
	)
}

/*
- The renderer of a file size column.
*/
func (m *Model) renderFileSize(indexElement int, columnWidth int) string {
	return common.PrettierFixedWidthItem(
		common.FormatFileSize(m.Element[indexElement].Info.Size()),
		columnWidth,
		false,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

/*
- The renderer of a modify time column.
TODO: make time template configurable
*/
func (m *Model) renderModifyTime(indexElement int, columnWidth int) string {
	modifyTime := m.Element[indexElement].Info.ModTime().Format("2006-01-02 15:04")
	return common.PrettierFixedWidthItem(
		modifyTime,
		columnWidth,
		false,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

/*
- The renderer of a permission column.
*/
func (m *Model) renderPermissions(indexElement int, columnWidth int) string {
	return common.PrettierFixedWidthItem(
		m.Element[indexElement].Info.Mode().Perm().String(),
		columnWidth,
		false,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}
