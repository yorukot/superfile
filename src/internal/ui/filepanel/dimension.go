package filepanel

import (
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) UpdateDimensions(width, height int) {
	m.SetWidth(width)
	m.SetHeight(height)
}

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
		widthExtraColumn := ansi.StringWidth(ColumnDelimiter) + col.Size
		if remainWidth-widthExtraColumn > extraColumnsThreshold {
			delimiterCol := columnDefinition{
				Name:      "",
				Generator: m.renderDelimiter,
				Size:      ansi.StringWidth(ColumnDelimiter),
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

func (m *Model) SetWidth(width int) {
	if width < MinWidth {
		width = MinWidth
	}
	m.width = width
	m.SearchBar.Width = m.width - common.InnerPadding
	m.columns = m.makeColumns(common.Config.FilePanelExtraColumns)
}

func (m *Model) SetHeight(height int) {
	if height < MinHeight {
		height = MinHeight
	}
	m.height = height
	// Adjust scroll if needed
	m.scrollToCursor(m.Cursor)
}

func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) GetHeight() int {
	return m.height
}

func (m *Model) GetMainPanelHeight() int {
	return m.height - common.BorderPadding
}

func (m *Model) GetContentWidth() int {
	return m.width - common.BorderPadding
}

func (m *Model) NeedRenderHeaders() bool {
	return common.Config.FilePanelExtraColumns > 0 && len(m.columns) > 1
}

// PanelElementHeight calculates the number of visible elements in content area
func (m *Model) PanelElementHeight() int {
	headerHeight := 0
	if m.NeedRenderHeaders() {
		headerHeight = ColumnHeaderHeight
	}
	return m.GetMainPanelHeight() - contentPadding - headerHeight
}
