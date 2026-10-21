package internal

import (
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
)

// updateMouseRegions rebuilds the mouse hit-test registry from the layout the
// next frame will use.
//
// It runs as part of rendering rather than on mouse input on purpose: a mouse
// event has to be resolved against the frame the user is looking at, and by the
// time the event arrives the layout may already have changed. Publishing the
// regions next to the render that produces them keeps the two in step.
func (m *model) updateMouseRegions() {
	reg := &m.mouseRegions
	reg.Reset()

	m.sidebarModel.MouseRegions(reg, m.focusPanel == sidebarFocus, 0, 0)

	// The sidebar occupies the leftmost columns when it is enabled, so the file
	// panels start right after its border.
	panelAreaX := common.Config.SidebarWidth
	if panelAreaX != 0 {
		panelAreaX += common.BorderPadding
	}
	m.fileModel.MouseRegions(reg, panelAreaX, 0)
}

// mouseTargetAt returns the widget drawn at the given terminal coordinates.
func (m *model) mouseTargetAt(x, y int) mouse.Target {
	return m.mouseRegions.TargetAt(x, y)
}
