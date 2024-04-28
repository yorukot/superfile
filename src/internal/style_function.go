package internal

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Generate border style for file panel
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

// Generate border style for sidebar
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

// Generate border style for process and can custom bottom border
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

// Generate border style for metadata and can custom bottom border
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

// Generate border style for clipboard and can custom bottom border
func clipboardBoarder(height int, width int, borderBottom string) lipgloss.Style {
	border := generateBorder()
	border.Top = "━┫ Clipboard w" + strings.Repeat("━", width)
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

// Generate modal (pop up widnwos) border style
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

// Generate first use modal style (This modal pop up when user first use superfile)
func firstUseModal(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Left, lipgloss.Center).
		Border(lipgloss.ThickBorder()).
		BorderForeground(modalBorderActiveColor).
		BorderBackground(modalBGColor).
		Background(modalBGColor).
		Foreground(modalFGColor)
}

// Generate help menu modal border style
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

// Generate full screen style for terminal size too small etc
func fullScreenStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Background(fullScreenBGColor).
		Foreground(fullScreenFGColor)
}

// Generate file panel divider style
func filePanelDividerStyle(focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(filePanelFocusColor(focusType)).
		Background(filePanelBGColor)
}

// Return border color based on file panel status
func filePanelFocusColor(focusType filePanelFocusType) lipgloss.Color {
	if focusType == noneFocus {
		return filePanelBorderColor
	} else {
		return filePanelBorderActiveColor
	}
}

// Return only fg and bg color style
func stringColorRender(fgColor lipgloss.Color, bgColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(fgColor).
		Background(bgColor)
}

// Geerate border style
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