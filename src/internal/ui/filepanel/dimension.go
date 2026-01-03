package filepanel

import (
	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) UpdateDimensions(width, height int) {
	m.SetWidth(width)
	m.SetHeight(height)
}

func (m *Model) SetWidth(width int) {
	if width < MinWidth {
		width = MinWidth
	}
	m.width = width
	m.SearchBar.Width = m.width - common.InnerPadding
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

// PanelElementHeight calculates the number of visible elements in content area
func (m *Model) PanelElementHeight() int {
	headerHeight := 0
	if common.Config.FilePanelExtraColumns > 0 {
		headerHeight = ColumnHeaderHeight
	}

	return m.GetMainPanelHeight() - contentPadding - headerHeight - headerHeight
}
