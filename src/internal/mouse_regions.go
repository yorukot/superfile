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

	m.fileModel.MouseRegions(reg, m.filePanelAreaX(), 0)

	// The context menu is registered last so that it claims the cells it covers,
	// which is what makes a click on a menu entry reach the menu rather than the
	// widget it is drawn on top of.
	m.contextMenu.MouseRegions(reg)
}

// filePanelAreaX returns the terminal column at which the file panel area
// starts, the sidebar's border included.
//
// Every coordinate that is resolved against a panel has to be measured from here
// rather than from the left edge of the terminal, so this is the one place that
// knows about the sidebar's width. Regions and drops both read it, which keeps
// them from disagreeing about which panel a column belongs to.
func (m *model) filePanelAreaX() int {
	if common.Config.SidebarWidth == 0 {
		return 0
	}
	return common.Config.SidebarWidth + common.BorderPadding
}

// mouseTargetAt returns the widget drawn at the given terminal coordinates.
func (m *model) mouseTargetAt(x, y int) mouse.Target {
	return m.mouseRegions.TargetAt(x, y)
}
