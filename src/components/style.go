package components

import (
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

var (
	downBarHeight = 13
)

var (
	sideBarWidth      = 20
	projectTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC241D"))
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C287"))
	pinnedLineStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A4A2A2"))
	pinnedTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC241D"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8751A"))
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8EC07C"))
)

var (
	fileIconStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8EC07C"))
	fileLocation  = lipgloss.NewStyle().Foreground(lipgloss.Color("#458588"))
)

func SideBarBoardStyle(height int, focus bool) lipgloss.Style {
	if focus {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#D79921")).
			Width(sideBarWidth).
			Height(height)
	} else {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A4A2A2")).
			Width(sideBarWidth).
			Height(height)
	}

}

func FilePanelBoardStyle(height int, width int, focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(FilePanelFocusColor(focusType))).
		Width(width).
		Height(height)
}

func FilePanelDividerStyle(focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(FilePanelFocusColor(focusType)))
}

func TruncateText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	return string(runes[:maxChars-3]) + "..."
}

func TruncateTextBeginning(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	// 计算应该保留的字符数量
	charsToKeep := maxChars - 3 // 减去省略号的长度
	// 截取字符串，从开头到指定的字符数
	truncatedRunes := append([]rune("..."), runes[len(runes)-charsToKeep:]...)
	return string(truncatedRunes)
}
