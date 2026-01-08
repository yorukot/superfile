package filepanel

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// The renderer for the mandatory first column in the file panel, with a name, a cursor, and a select option.
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

// The renderer of delimiter spaces. It has a strict fixed size that depends only on the delimiter string.
func (m *Model) renderDelimiter(indexElement int, columnWidth int) string {
	isSelected := m.CheckSelected(m.Element[indexElement].Location)
	return common.PrettierFixedWidthItem(
		ColumnDelimiter,
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Left,
	)
}

func (m *Model) renderFileSize(indexElement int, columnWidth int) string {
	isSelected := m.CheckSelected(m.Element[indexElement].Location)
	sizeValue := common.FormatFileSize(m.Element[indexElement].Info.Size())
	if m.Element[indexElement].Info.IsDir() {
		sizeValue = ""
	}
	return common.PrettierFixedWidthItem(
		sizeValue,
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

// TODO: make time template configurable
func (m *Model) renderModifyTime(indexElement int, columnWidth int) string {
	isSelected := m.CheckSelected(m.Element[indexElement].Location)
	modifyTime := m.Element[indexElement].Info.ModTime().Format("2006-01-02 15:04")
	return common.PrettierFixedWidthItem(
		modifyTime,
		columnWidth,
		isSelected,
		common.FilePanelBGColor,
		lipgloss.Right,
	)
}

func (m *Model) renderPermissions(indexElement int, columnWidth int) string {
	isSelected := m.CheckSelected(m.Element[indexElement].Location)
	return common.PrettierFixedWidthItem(
		m.Element[indexElement].Info.Mode().Perm().String(),
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
	return common.PrettierFixedWidthItem(
		cd.Name,
		cd.Size,
		false,
		common.FilePanelBGColor,
		lipgloss.Center,
	)
}

func (m *Model) makeColumns(columnThreshold int, filePanelNamePercent int) []columnDefinition {
	// TODO: make column set configurable
	// Note: May use a predefined slice for efficiency. This content is static
	extraColumns := []columnDefinition{
		{Name: "Size", columnRender: m.renderFileSize, Size: FileSizeColumnWidth},
		{Name: "Modify time", columnRender: m.renderModifyTime, Size: ModifyTimeSizeColumnWidth},
		{Name: "Permission", columnRender: m.renderPermissions, Size: PermissionsColumnWidth},
	}
	maxColumns := min(columnThreshold, len(extraColumns))

	columns := []columnDefinition{
		{Name: "Name", columnRender: m.renderFileName, Size: m.GetContentWidth()},
	}

	fileNameRatio := FileNameRatioDefault
	if filePanelNamePercent > FileNameRatioDefault &&
		filePanelNamePercent <= FileNameRatioMax {
		fileNameRatio = filePanelNamePercent
	}

	// "-1" guards in a cases of rounding numbers.
	extraColumnsThreshold := int(float64(m.GetContentWidth()*fileNameRatio/FileNameRatioMax) - 1)
	if extraColumnsThreshold <= 0 {
		extraColumnsThreshold = m.GetContentWidth()
	}

	for _, col := range extraColumns[0:maxColumns] {
		widthExtraColumn := ansi.StringWidth(ColumnDelimiter) + col.Size

		// This condition checks that can we borrow some width from first column for additional columns?
		if columns[0].Size-widthExtraColumn > extraColumnsThreshold {
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
