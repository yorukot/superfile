package contextmenu

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// Render draws the menu.
//
// The layout it produces is exactly what MouseRegions describes: one row per
// entry, starting one row below the top border, each as wide as the whole menu.
func (m *Model) Render() string {
	if !m.open {
		return ""
	}

	labelWidth := 0
	shortcutWidth := 0
	for _, item := range m.items {
		labelWidth = max(labelWidth, lipgloss.Width(item.Label))
		shortcutWidth = max(shortcutWidth, lipgloss.Width(item.Shortcut))
	}

	var content strings.Builder
	for i, item := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor)
		}

		content.WriteString(cursor)
		content.WriteString(" ")
		content.WriteString(lipgloss.NewStyle().
			Width(labelWidth).
			Render(common.ModalStyle.Render(item.Label)))
		if shortcutWidth > 0 {
			content.WriteString(" ")
			content.WriteString(lipgloss.NewStyle().
				Width(shortcutWidth).
				Align(lipgloss.Right).
				Render(common.HelpMenuHotkeyStyle.Render(item.Shortcut)))
		}
		content.WriteString("\n")
	}

	// GenerateFooterBorder returns its argument's width plus the two junction
	// characters it wraps the count in, while lipgloss only leaves the box width
	// minus its corners for the bottom border. Sizing for both is what keeps the
	// last part of the count from being cut off.
	bottomBorder := common.GenerateFooterBorder(
		fmt.Sprintf("%s/%s", strconv.Itoa(m.cursor+1), strconv.Itoa(len(m.items))),
		m.width-sideBorders-common.BorderPadding)

	return common.ContextMenuBorderStyle(m.height, m.width, bottomBorder).
		Render(strings.TrimSuffix(content.String(), "\n"))
}
