package components

import "github.com/charmbracelet/lipgloss"

func FilePanelFocusColor(focusType filePanelFocusType) string {
	if focusType == noneFocus {
		return "#A4A2A2"
	} else if focusType == secondFocus {
		return "#656565"
	} else {
		return "#D79921"
	}
}

func FilePanelBoard(focusType filePanelFocusType) lipgloss.Border {
	if focusType == noneFocus {
		return lipgloss.RoundedBorder()
	} else {
		return lipgloss.DoubleBorder()
	}
}