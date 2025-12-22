package filepanel

import "github.com/yorukot/superfile/src/internal/common"

// UpdateDimensions sets the panel dimensions with validation
func (m *Model) UpdateDimensions(width, height int) {
	m.SetWidth(width)
	m.SetHeight(height)
}

// GetWidth returns the total panel width
func (m *Model) SetWidth(width int) {
	if width < FilePanelMinWidth {
		width = FilePanelMinWidth
	}
	m.width = width
	m.SearchBar.Width = m.width - common.InnerPadding
}

func (m *Model) SetHeight(height int) {
	if height < FilePanelMinHeight {
		height = FilePanelMinHeight
	}
	m.height = height
	// Adjust scroll if needed
	m.scrollToCursor(m.Cursor)
}

// GetWidth returns the total panel width
func (m *Model) GetWidth() int {
	return m.width
}

// GetHeight returns the total panel height
func (m *Model) GetHeight() int {
	return m.height
}

// GetMainPanelHeight returns content height (total height minus borders)
func (m *Model) GetMainPanelHeight() int {
	return m.height - common.BorderPadding
}

// GetContentWidth returns content width (total width minus borders)
func (m *Model) GetContentWidth() int {
	return m.width - common.BorderPadding
}

// PanelElementHeight calculates the number of visible elements in content area
func (m *Model) PanelElementHeight() int {
	return m.GetMainPanelHeight() - FilePanelContentPadding
}
