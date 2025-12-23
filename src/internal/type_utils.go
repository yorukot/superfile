package internal

import (
	"github.com/yorukot/superfile/src/internal/common"
)

// reset the items slice and set the cut value
func (c *CopyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
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
		return common.InvalidTypeString
	}
}
