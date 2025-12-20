package internal

import (
	"fmt"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

const invalidTypeString = "InvalidType"

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
	if m.fullHeight != (m.mainPanelHeight+common.BorderPadding)+utils.FullFooterHeight(m.footerHeight, m.toggleFooter) {
		return fmt.Errorf("invalid model layout, fullHeight : %v, mainPanelHeight : %v, footerHeight : %v",
			m.fullHeight, m.mainPanelHeight, m.footerHeight)
	}
	// TODO : Add check for width as well
	return nil
}

// ================ filepanel

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

func FilePanelSlice(paths []string) []FilePanel {
	res := make([]FilePanel, len(paths))
	for i := range paths {
		// Making the first panel as the focussed
		isFocus := i == 0
		res[i] = DefaultFilePanel(paths[i], isFocus)
	}
	return res
}
