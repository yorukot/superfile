package internal

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}

// String method for easy logging
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
		return "Invalid"
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
		return "Invalid"
	}
}

func (p panelMode) String() string {
	switch p {
	case selectMode:
		return "selectMode"
	case browserMode:
		return "browserMode"
	default:
		return "Invalid"
	}
}
