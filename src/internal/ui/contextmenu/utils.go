package contextmenu

import (
	"charm.land/lipgloss/v2"
)

const (
	// borderRows is the number of rows taken by the top and bottom borders.
	borderRows = 2
	// cursorColumn is the width of the column holding the cursor indicator.
	cursorColumn = 1
	// cursorPad is the space between the cursor indicator and the label.
	cursorPad = 1
	// shortcutGap is the space between the label and the shortcut column.
	shortcutGap = 1
	// sideBorders is the number of columns taken by the left and right borders.
	sideBorders = 2
)

// contentWidth returns the width of the rendered menu for the given entries.
//
// Every row is padded to the widest label and the widest shortcut so that the
// columns line up, which is also what lets hit-test regions and the drawn rows
// agree on where each entry starts and ends.
func contentWidth(items []Item) int {
	labelWidth := 0
	shortcutWidth := 0
	for _, item := range items {
		labelWidth = max(labelWidth, lipgloss.Width(item.Label))
		shortcutWidth = max(shortcutWidth, lipgloss.Width(item.Shortcut))
	}

	inner := cursorColumn + cursorPad + labelWidth
	if shortcutWidth > 0 {
		inner += shortcutGap + shortcutWidth
	}
	return inner + sideBorders
}

// itemRowY returns the terminal row of an entry relative to the top of the menu.
// The first entry sits directly below the top border.
func (m *Model) itemRowY(index int) int {
	return m.originY + 1 + index
}
