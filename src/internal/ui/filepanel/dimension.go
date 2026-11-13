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
	m.SearchBar.SetWidth(m.width - common.InnerPadding)
	m.columns = m.makeColumns(common.Config.FilePanelExtraColumns, common.Config.FilePanelNamePercent)
}

func (m *Model) SetHeight(height int) {
	if height < MinHeight {
		height = MinHeight
	}
	m.height = height
	// Adjust scroll if needed
	m.scrollToCursor(m.GetCursor())
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

// NeedRenderHeaders returns whether column headers take a line above the entries.
func (m *Model) NeedRenderHeaders() bool {
	return common.Config.FilePanelExtraColumns > 0 && len(m.columns) > 1
}

// ItemAreaTop returns the terminal row, counted from the top of the panel and
// including its top border, on which the first file entry is drawn.
//
// Rendering and mouse hit-testing both derive their row offsets from here, so
// the two cannot disagree about where the entries begin.
func (m *Model) ItemAreaTop() int {
	top := common.BorderTopHeight + contentPadding
	if m.NeedRenderHeaders() {
		top += ColumnHeaderHeight
	}
	return top
}

// PanelElementHeight calculates the number of visible elements in content area
func (m *Model) PanelElementHeight() int {
	// Everything below the item area is taken up by the bottom border.
	return m.GetHeight() - m.ItemAreaTop() - common.BorderBottomHeight
}
