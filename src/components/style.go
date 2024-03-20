package components

import (
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

func SideBarBoardStyle(height int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#A4A2A2")).
		Width(sideBarWidth).
		Height(height)
}

func FilePanelBoardStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#A4A2A2")).
		Width(width).
		Height(height)
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
