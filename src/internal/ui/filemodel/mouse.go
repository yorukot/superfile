package filemodel

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/mouse"
)

// PanelOriginX returns the column, relative to the file model area, at which
// the panel with the given index starts. It returns zero for an out of range
// index so that callers probing an unknown index fall back to the first column
// instead of producing a negative offset.
func (m *Model) PanelOriginX(index int) int {
	if index < 0 || index >= len(m.panelOriginsX) {
		slog.Error("Unexpected panel index while resolving panel origin", "index", index,
			"panel count", m.PanelCount())
		return 0
	}
	return m.panelOriginsX[index]
}

// MouseRegions registers the entries drawn by every file panel. originX and
// originY are the terminal coordinates of the file model area's top left
// corner, its panels' borders included.
func (m *Model) MouseRegions(reg *mouse.Registry, originX, originY int) {
	if m.PanelCount() == 0 {
		return
	}

	for i := range m.FilePanels {
		m.FilePanels[i].MouseRegions(reg, i, originX+m.PanelOriginX(i), originY)
	}
}
