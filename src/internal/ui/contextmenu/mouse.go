package contextmenu

import (
	"github.com/yorukot/superfile/src/internal/mouse"
)

// MouseRegions publishes one region per entry, using the same origin and widths
// that Render draws with.
//
// These are registered after the panels have published theirs, so the
// registry's last-added-wins lookup hands clicks inside the menu to the menu
// rather than to whatever it is covering.
func (m *Model) MouseRegions(reg *mouse.Registry) {
	if !m.open {
		return
	}

	for i := range m.items {
		reg.AddRow(mouse.ContextMenuItemTarget(i),
			m.originX, m.itemRowY(i), m.width, 1)
	}
}
