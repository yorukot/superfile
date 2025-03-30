package common

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/config/icon"
)

// Generate border style for file panel
func filePanelBorderStyle(height int, width int, filePanelFocussed bool, borderBottom string) lipgloss.Style {
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
		BorderForeground(filePanelFocusColor(filePanelFocussed)).
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
func sideBarBorderStyle(height int, sidebarFocussed bool) lipgloss.Style {
	border := generateBorder()
	sidebarBorderStateColor := sidebarBorderColor
	if sidebarFocussed {
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
func procsssBarBorder(height int, width int, borderBottom string, processBarFocussed bool) lipgloss.Style {
	border := generateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Processes " + Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
	border.Bottom = borderBottom

	processBorderStateColor := footerBorderColor
	if processBarFocussed {
		processBorderStateColor = footerBorderActiveColor
	}

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(processBorderStateColor).
		BorderBackground(FooterBGColor).
		Width(width).
		Height(height).
		Background(FooterBGColor).
		Foreground(footerFGColor)
}

// Generate border style for metadata and can custom bottom border
func metadataBorder(height int, width int, borderBottom string, metadataFocussed bool) lipgloss.Style {
	border := generateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Metadata " + Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
	border.Bottom = borderBottom

	metadataBorderStateColor := footerBorderColor
	if metadataFocussed {
		metadataBorderStateColor = footerBorderActiveColor
	}

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(metadataBorderStateColor).
		BorderBackground(FooterBGColor).
		Width(width).
		Height(height).
		Background(FooterBGColor).
		Foreground(footerFGColor)
}

// Generate border style for clipboard and can custom bottom border
func clipboardBorder(height int, width int, borderBottom string) lipgloss.Style {
	border := generateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Clipboard " + Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
	border.Bottom = borderBottom

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(footerBorderColor).
		BorderBackground(FooterBGColor).
		Width(width).
		Height(height).
		Background(FooterBGColor).
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

// Generate sort options modal border style
func sortOptionsModalBorderStyle(height int, width int, borderBottom string) lipgloss.Style {
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
func filePanelDividerStyle(filePanelFocussed bool) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(filePanelFocusColor(filePanelFocussed)).
		Background(filePanelBGColor)
}

// Return border color based on file panel status
func filePanelFocusColor(filePanelFocussed bool) lipgloss.Color {
	if filePanelFocussed {
		return filePanelBorderActiveColor
	}
	return filePanelBorderColor
}

// Return only fg and bg color style
func StringColorRender(fgColor lipgloss.Color, bgColor lipgloss.Color) lipgloss.Style {
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
func loadHotkeysError(value string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("■ ERROR: ") + "Hotkeys file \"" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF")).Render(value) + "\" invalidation"
}

// Generate search bar for file panel
func GenerateSearchBar() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = FooterCursorStyle
	ti.Cursor.TextStyle = FooterStyle
	ti.TextStyle = FilePanelStyle
	ti.Prompt = filePanelTopDirectoryIconStyle.Render(icon.Search + icon.Space)
	ti.Cursor.Blink = true
	ti.PlaceholderStyle = FilePanelStyle
	ti.Placeholder = "(" + Hotkeys.SearchBar[0] + ") Type something"
	ti.Blur()
	ti.CharLimit = 156
	return ti
}

func GeneratePromptTextInput() textinput.Model {
	t := textinput.New()
	t.Prompt = ""
	t.CharLimit = 156
	t.SetValue("")
	t.Cursor.Style = ModalCursorStyle
	t.Cursor.TextStyle = ModalStyle
	t.TextStyle = ModalStyle
	t.PlaceholderStyle = ModalStyle

	return t
}

func GenerateNewFileTextInput() textinput.Model {
	t := textinput.New()
	t.Cursor.Style = ModalCursorStyle
	t.Cursor.TextStyle = ModalStyle
	t.TextStyle = ModalStyle
	t.Cursor.Blink = true
	t.Placeholder = "Add \"" + string(filepath.Separator) + "\" transcend folders"
	t.PlaceholderStyle = ModalStyle
	t.Focus()
	t.CharLimit = 156
	t.Width = ModalWidth - 10
	return t
}

func GenerateRenameTextInput(width int, cursorPos int, defaultValue string) textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = FilePanelCursorStyle
	ti.Cursor.TextStyle = FilePanelStyle
	ti.Prompt = FilePanelCursorStyle.Render(icon.Cursor + " ")
	ti.TextStyle = ModalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "New name"
	ti.PlaceholderStyle = ModalStyle
	ti.SetValue(defaultValue)
	ti.SetCursor(cursorPos)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = width

	return ti
}

func GeneratePinnedRenameTextInput(cursorPos int, defaultValue string) textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = FilePanelCursorStyle
	ti.Cursor.TextStyle = FilePanelStyle
	ti.Prompt = FilePanelCursorStyle.Render(icon.Cursor + " ")
	ti.TextStyle = ModalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "New name"
	ti.PlaceholderStyle = ModalStyle
	ti.SetValue(defaultValue)
	ti.SetCursor(cursorPos)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = Config.SidebarWidth - 4
	return ti
}

func GenerateDefaultProgress() progress.Model {
	prog := progress.New(GenerateGradientColor())
	prog.PercentageStyle = FooterStyle
	return prog
}
