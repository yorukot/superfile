package internal

import (
	"fmt"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/common/utils"
)

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}

// ================ Model related utils =======================

// Non fatal Validations. This indicates bug / programming errors, not user configuration mistake
func (m *model) validateLayout() error {
	if 0 < m.footerHeight && m.footerHeight < common.MinFooterHeight {
		return fmt.Errorf("footerHeight %v is too small", m.footerHeight)
	}
	if !m.toggleFooter && m.footerHeight != 0 {
		return fmt.Errorf("footer closed and footerHeight %v is non zero", m.footerHeight)
	}
	// PanelHeight + 2 lines (main border) + actual footer height
	if m.fullHeight != (m.mainPanelHeight+2)+utils.FullFooterHeight(m.footerHeight, m.toggleFooter) {
		return fmt.Errorf("Invalid model layout, fullHeight : %v, mainPanelHeight : %v, footerHeight : %v\n",
			m.fullHeight, m.mainPanelHeight, m.footerHeight)
	}
	// Todo : Add check for width as well
	return nil
}

// ================ Sidebar related utils =====================
// Hopefully compiler inlines it
func (d directory) isDivider() bool {
	return d == pinnedDividerDir || d == diskDividerDir
}
func (d directory) requiredHeight() int {
	if d.isDivider() {
		return 3
	}
	return 1
}

// ================ String method for easy logging =====================

func (f focusPanelType) String() string {
	switch f {
	case nonePanelFocus:
		return "nonePanelFocus"
	case processBarFocus:
		return "processBarFocus"
	case sidebarFocus:
		return "sidebarFocus"
	case metadataFocus:
		return "metadataFocus"
	default:
		return invalidTypeString
	}
}

func (f filePanelFocusType) String() string {
	switch f {
	case noneFocus:
		return "noneFocus"
	case secondFocus:
		return "secondFocus"
	case focus:
		return "focus"
	default:
		return invalidTypeString
	}

}

func (p panelMode) String() string {
	switch p {
	case selectMode:
		return "selectMode"
	case browserMode:
		return "browserMode"
	default:
		return invalidTypeString
	}
}
