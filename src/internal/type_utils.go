package internal

import "fmt"

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}

// ================ Model related utils =======================

// Non fatal Validations. This indicates bug / programming errors, not user configuration mistake
func (m *model) validateLayout() error {
	if 0 < m.footerHeight && m.footerHeight < minFooterHeight {
		return fmt.Errorf("footerHeight %v is too small", m.footerHeight)
	}
	// PanelHeight + 2 lines (main border) + actual footer height
	if m.fullHeight != (m.mainPanelHeight+2)+actualfooterHeight(m.footerHeight, m.commandLine.input.Focused()) {
		return fmt.Errorf("Invalid model layout, fullHeight : %v, mainPanelHeight : %v, footerHeight : %v\n",
			m.fullHeight, m.mainPanelHeight, m.footerHeight)
	}
	// Todo : Add check for width as well
	return nil
}

func actualfooterHeight(footerHeight int, commandLineFocussed bool) int {
	// footerHeight + 2 or 0 lines (footer border)
	// + 1 lines ( commmand line only if footersize is >0)
	footerBorder := 2
	if footerHeight == 0 {
		footerBorder = 0
	}
	commandLineHeight := 0
	if commandLineFocussed && footerHeight != 0 {
		commandLineHeight = 1
	}
	return footerHeight + footerBorder + commandLineHeight
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
