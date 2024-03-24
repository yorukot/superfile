package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	downBarHeight = 13
)

var (
	sideBarWidth      = 20
	projectTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC241D")).Bold(true)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C287"))
	pinnedLineStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A4A2A2"))
	pinnedTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC241D")).Bold(true)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8751A")).Bold(true)
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8EC07C")).Bold(true)
)

var (
	fileIconStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8EC07C")).Bold(true)
	fileLocation  = lipgloss.NewStyle().Foreground(lipgloss.Color("#458588")).Bold(true)
)

func StringColorRender(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func SideBarBoardStyle(height int, focus bool) lipgloss.Style {
	if focus {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#D79921")).
			MaxWidth(height).
			Width(sideBarWidth).
			Height(height).Bold(true)
	} else {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A4A2A2")).
			MaxWidth(height).
			Width(sideBarWidth).
			Height(height).Bold(true)
	}

}

func FilePanelBoardStyle(height int, width int, focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(FilePanelBoard(focusType)).
		BorderForeground(lipgloss.Color(FilePanelFocusColor(focusType))).
		Width(width).
		Height(height).Bold(true)
}

func FilePanelDividerStyle(focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(FilePanelFocusColor(focusType))).Bold(true)
}

func TruncateText(text string, maxChars int) string {
	if len(text) <= maxChars {
		return text
	}
	return text[:maxChars-3] + "..."
}

func TruncateTextBeginning(text string, maxChars int) string {
	if len(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	charsToKeep := maxChars - 3
	truncatedRunes := append([]rune("..."), runes[len(runes)-charsToKeep:]...)
	return string(truncatedRunes)
}

func PrettierName(name string, isDir bool) string {
	style := getElementIcon(name, isDir)
	return StringColorRender(style.color).Render(style.icon) + "  " + itemStyle.Render(name)
}