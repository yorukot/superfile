package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func filePanelBorderStyle(height int, width int, focusType filePanelFocusType, borderBottom string) lipgloss.Style {
	border := generateBorder()
	border.Left = ""
	border.Right = ""
	for i := 0; i < height; i++ {
		if i == 1 {
			border.Left += "┣"
			border.Right += "┫"
		} else {
			border.Left += "┃"
			border.Right += "┃"
		}
	}
	border.Bottom = borderBottom
	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(filePanelFocusColor(focusType)).
		BorderBackground(filePanelBGColor).
		Width(width).
		Height(height).Background(filePanelBGColor)
}

func sideBarBorderStyle(height int, focus focusPanelType) lipgloss.Style {
	border := generateBorder()
	sidebarBorderStateColor := sidebarBorderColor
	if focus == sidebarFocus {
		sidebarBorderStateColor = sidebarBorderActiveColor
	}

	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(sidebarBorderStateColor).
		BorderBackground(sidebarBGColor).
		Width(sidebarWidth).
		Height(height).
		Background(sidebarBGColor).
		Foreground(sidebarFGColor)
}

func procsssBarBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	border := generateBorder()
	border.Top = "━┫ Processes ┣" + strings.Repeat("━", width)
	border.Bottom = borderBottom

	processBorderStateColor := footerBorderColor
	if focusType == processBarFocus {
		processBorderStateColor = footerBorderActiveColor
	}

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(processBorderStateColor).
		BorderBackground(footerBGColor).
		Width(width).
		Height(height).
		Background(footerBGColor).
		Foreground(footerFGColor)
}

func metadataBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	border := generateBorder()
	border.Top = "━┫ Metadata ┣" + strings.Repeat("━", width)
	border.Bottom = borderBottom

	metadataBorderStateColor := footerBorderColor
	if focusType == metadataFocus {
		metadataBorderStateColor = footerBorderActiveColor
	}

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(metadataBorderStateColor).
		BorderBackground(footerBGColor).
		Width(width).
		Height(height).
		Background(footerBGColor).
		Foreground(footerFGColor)
}

func clipboardBoarder(height int, width int, borderBottom string) lipgloss.Style {
	border := generateBorder()
	border.Top = "━┫ Clipboard ┣" + strings.Repeat("━", width)
	border.Bottom = borderBottom

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(footerBorderColor).
		BorderBackground(footerBGColor).
		Width(width).
		Height(height).
		Background(footerBGColor).
		Foreground(footerFGColor)
}

func modalBorderStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.ThickBorder()).
		BorderForeground(modalBorderActiveColor).
		BorderBackground(modalBGColor).
		Background(modalBGColor).
		Foreground(modalFGColor)
}

func helpMenuModalBorderStyle(height int, width int, borderBottom string) lipgloss.Style {
	border := generateBorder()
	border.Bottom = borderBottom

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(modalBorderActiveColor).
		BorderBackground(modalBGColor).
		Width(width).
		Height(height).
		Background(modalBGColor).
		Foreground(modalFGColor)
}
func fullScreenStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Background(fullScreenBGColor).
		Foreground(fullScreenFGColor)
}

func filePanelDividerStyle(focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(filePanelFocusColor(focusType)).
		Background(filePanelBGColor)
}

func filePanelFocusColor(focusType filePanelFocusType) lipgloss.Color {
	if focusType == noneFocus {
		return filePanelBorderColor
	} else {
		return filePanelBorderActiveColor
	}
}

func stringColorRender(fgColor lipgloss.Color, bgColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(fgColor).
		Background(bgColor)
}

func generateBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}
}
