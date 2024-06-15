package internal

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/config/icon"
)

// Generate border style for file panel
func filePanelBorderStyle(height int, width int, focusType filePanelFocusType, borderBottom string) lipgloss.Style {
	border := generateBorder()
	border.Left = ""
	border.Right = ""
	for i := 0; i < height; i++ {
		if i == 1 {
			border.Left += Config.BorderMiddleLeft
			border.Right += Config.BorderMiddleRight
		} else {
			border.Left += Config.BorderLeft
			border.Right += Config.BorderRight
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

// Generate filePreview Box
func filePreviewBox(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
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
		Width(Config.SidebarWidth).
		Height(height).
		Background(sidebarBGColor).
		Foreground(sidebarFGColor)
}

// Generate border style for process and can custom bottom border
func procsssBarBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	border := generateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Processes " + Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
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
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Metadata " + Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
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
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Clipboard " + Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
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
	border := generateBorder()
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Border(border).
		BorderForeground(modalBorderActiveColor).
		BorderBackground(modalBGColor).
		Background(modalBGColor).
		Foreground(modalFGColor)
}

// Generate first use modal style (This modal pop up when user first use superfile)
func firstUseModal(height int, width int) lipgloss.Style {
	border := generateBorder()
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Left, lipgloss.Center).
		Border(border).
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

// Generate border style
func generateBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         Config.BorderTop,
		Bottom:      Config.BorderBottom,
		Left:        Config.BorderLeft,
		Right:       Config.BorderRight,
		TopLeft:     Config.BorderTopLeft,
		TopRight:    Config.BorderTopRight,
		BottomLeft:  Config.BorderBottomLeft,
		BottomRight: Config.BorderBottomRight,
	}
}

// Generate config error style
func loadConfigError(value string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("■ ERROR: ") + "Config file \"" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF")).Render(value) + "\" invalidation"
}

// Generate config error style
func lodaHotkeysError(value string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("■ ERROR: ") + "Hotkeys file \"" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF")).Render(value) + "\" invalidation"
}

// Generate search bar for file panel
func generateSearchBar() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = footerCursorStyle
	ti.Cursor.TextStyle = footerStyle
	ti.TextStyle = filePanelStyle
	ti.Prompt = filePanelTopDirectoryIconStyle.Render(icon.Search + icon.Space)
	ti.Cursor.Blink = true
	ti.PlaceholderStyle = filePanelStyle
	ti.Placeholder = "(" + hotkeys.SearchBar[0] + ") Type something"
	ti.Blur()
	ti.CharLimit = 156
	return ti
}

// Generate command line in the bottom
func generateCommandLineInputBox() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = footerCursorStyle
	ti.Cursor.TextStyle = footerStyle
	ti.TextStyle = filePanelStyle
	ti.Prompt = filePanelTopDirectoryIconStyle.Render(icon.Cursor + icon.Space)
	ti.Cursor.Blink = true
	ti.PlaceholderStyle = filePanelStyle
	ti.Blur()
	return ti
}