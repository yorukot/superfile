package filepanel

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// The renderer for the mandatory first column in the file panel, with a name, a cursor, and a select option.
func (m *Model) renderFileName(indexElement int, columnWidth int) string {
	elem := m.GetElementAtIdx(indexElement)
	isSelected := m.CheckSelected(elem.Location)
	cursor := emptyCursor
	if indexElement == m.GetCursor() && !m.SearchBar.Focused() {
		cursor = icon.Cursor
	}

	selectBox := m.renderSelectBox(isSelected)

	// Calculate the actual prefix width for proper alignment
	prefixWidth := ansi.StringWidth(cursor+" ") + ansi.StringWidth(selectBox)

	isLink := elem.Info.Mode()&os.ModeSymlink != 0
	renderedName := common.FilePanelItemRenderWithIcon(
		elem.Name,
		columnWidth-prefixWidth,
		elem.Directory,
		isLink,
		isSelected,
		common.FilePanelBGColor,
	)
	return common.FilePanelCursorStyle.Render(cursor+" ") + selectBox + renderedName
}

// The renderer of delimiter spaces. It has a strict fixed size that depends only on the delimiter string.
func (m *Model) renderDelimiter(indexElement int, columnWidth int) string {
	isSelected := m.CheckSelected(m.GetElementAtIdx(indexElement).Location)
	return common.FilePanelItemRender(
		ColumnDelimiter,
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Left,
	)
}

func (m *Model) renderFileSize(indexElement int, columnWidth int) string {
	elem := m.GetElementAtIdx(indexElement)
	isSelected := m.CheckSelected(elem.Location)
	sizeValue := common.FormatFileSize(elem.Info.Size())
	if elem.Info.IsDir() {
		sizeValue = ""
	}
	return common.FilePanelItemRender(
		sizeValue,
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

// TODO: make time template configurable
func (m *Model) renderModifyTime(indexElement int, columnWidth int) string {
	elem := m.GetElementAtIdx(indexElement)
	isSelected := m.CheckSelected(elem.Location)
	modifyTime := elem.Info.ModTime().Format("2006-01-02 15:04")
	return common.FilePanelItemRender(
		modifyTime,
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

func (m *Model) renderPermissions(indexElement int, columnWidth int) string {
	elem := m.GetElementAtIdx(indexElement)
	isSelected := m.CheckSelected(elem.Location)
	return common.FilePanelItemRender(
		elem.Info.Mode().Perm().String(),
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

func (cd *columnDefinition) Render(index int) string {
	return cd.columnRender(index, cd.Size)
}

func (cd *columnDefinition) RenderHeader() string {
	return common.FilePanelItemRender(
		cd.Name,
		cd.Size,
		false,
		common.FilePanelBGColor,
		cd.HeaderAlign,
	)
}

func (m *Model) makeColumns(columnThreshold int, fileNameRatio int) []columnDefinition {
	// TODO: make column set configurable
	// Note: May use a predefined slice for efficiency. This content is static
	extraColumns := []columnDefinition{
		{
			Name:         "Size",
			columnRender: m.renderFileSize,
			Size:         FileSizeColumnWidth,
			HeaderAlign:  lipgloss.Center,
		},
		{
			Name:         "Modify time",
			columnRender: m.renderModifyTime,
			Size:         ModifyTimeSizeColumnWidth,
			HeaderAlign:  lipgloss.Center,
		},
		{
			Name:         "Permission",
			columnRender: m.renderPermissions,
			Size:         PermissionsColumnWidth,
			HeaderAlign:  lipgloss.Center,
		},
	}
	maxColumns := min(columnThreshold, len(extraColumns))
	columns := []columnDefinition{
		{
			Name:         strings.Repeat(" ", ansi.StringWidth(emptyCursor+" ")) + "Name",
			columnRender: m.renderFileName,
			Size:         m.GetContentWidth(),
			HeaderAlign:  lipgloss.Left,
		},
	}

	minWidthForNameColumn := int(float64(m.GetContentWidth() * fileNameRatio / common.FileNameRatioMax))
	// Worst case (5 * 100 / 100) could evaluate to 5.0001
	// Hence, we need this check. Our constraints on Width and ratio guarantee it to be > 0 though
	minWidthForNameColumn = min(minWidthForNameColumn, m.GetContentWidth())

	for _, col := range extraColumns[0:maxColumns] {
		widthExtraColumn := ansi.StringWidth(ColumnDelimiter) + col.Size

		// This condition checks that can we borrow some width from first column for additional columns?
		if columns[0].Size-widthExtraColumn > minWidthForNameColumn {
			delimiterCol := columnDefinition{
				Name:         "",
				columnRender: m.renderDelimiter,
				Size:         ansi.StringWidth(ColumnDelimiter),
			}
			columns = append(columns, delimiterCol, col)
			columns[0].Size -= widthExtraColumn
		} else {
			break
		}
	}
	return columns
}
