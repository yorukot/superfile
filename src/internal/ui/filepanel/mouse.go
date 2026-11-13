package filepanel

import (
	"github.com/yorukot/superfile/src/internal/mouse"
)

// RenderedItemRange returns the half-open range of element indexes that the
// panel currently draws, that is the elements between renderIndex and the
// bottom border. Rendering and mouse hit-testing both read the range from here
// so a click always addresses a row that is actually on screen.
func (m *Model) RenderedItemRange() (int, int) {
	if m.Empty() {
		return 0, 0
	}
	return m.renderIndex, min(m.renderIndex+m.PanelElementHeight(), m.ElemCount())
}

// MouseRegions registers one region per file entry the panel currently draws.
// originX and originY are the terminal coordinates of the panel's top left
// corner, its border included, and panelIndex identifies the panel among its
// siblings.
//
// Only the rows are published. A click anywhere on a row addresses that entry,
// which keeps hit-testing independent of the configured extra columns and of
// the nerdfont-dependent select icons.
func (m *Model) MouseRegions(reg *mouse.Registry, panelIndex, originX, originY int) {
	start, end := m.RenderedItemRange()
	if start == end {
		return
	}

	top := originY + m.ItemAreaTop()
	for itemIndex := start; itemIndex < end; itemIndex++ {
		reg.Add(mouse.FilePanelItemTarget(panelIndex, itemIndex),
			originX, top, originX+m.GetWidth(), top+elementRowHeight)
		top += elementRowHeight
	}
}
