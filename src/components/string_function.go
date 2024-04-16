package components

import (
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

func truncateText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	return text[:maxChars-3] + "..."
}

func truncateTextBeginning(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	charsToKeep := maxChars - 3
	truncatedRunes := append([]rune("..."), runes[len(runes)-charsToKeep:]...)
	return string(truncatedRunes)
}

func truncateMiddleText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	halfEllipsisLength := (maxChars - 3) / 2

	truncatedText := text[:halfEllipsisLength] + "..." + text[utf8.RuneCountInString(text)-halfEllipsisLength:]

	return truncatedText
}

func prettierName(name string, width int, isDir bool, isSelected bool, bgColor lipgloss.Color) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.color), bgColor).
		Background(bgColor).
		Render(style.icon + "") + 
		filePanelItemSelectedStyle.
		Render(truncateText(name, width))
	} else {
		return stringColorRender(lipgloss.Color(style.color), bgColor).
		Background(bgColor).
		Render(style.icon + " ") + 
		filePanelStyle.Render(truncateText(name, width))
	}
}

func clipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.color), footerBGColor).
		Background(footerBGColor).
		Render(style.icon + " ") + 
		filePanelItemSelectedStyle.Render(truncateTextBeginning(name, width))
	} else {
		return stringColorRender(lipgloss.Color(style.color), footerBGColor).
		Background(footerBGColor).
		Render(style.icon + " ") + 
		filePanelStyle.Render(truncateTextBeginning(name, width))
	}
}