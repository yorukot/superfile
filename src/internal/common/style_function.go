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
func FilePanelBorderStyle(height int, width int, filePanelFocussed bool, borderBottom string) lipgloss.Style {
	border := GenerateBorder()
	border.Left = ""
	border.Right = ""

	for i := range height {
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
		BorderForeground(FilePanelFocusColor(filePanelFocussed)).
		BorderBackground(FilePanelBGColor).
		Width(width).
		Height(height).Background(FilePanelBGColor)
}

// Generate filePreview Box
func FilePreviewStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(FilePanelBGColor).
		Foreground(FilePanelFGColor)
}

// Generate border style for sidebar
func SideBarBorderStyle(height int, sidebarFocussed bool) lipgloss.Style {
	border := GenerateBorder()
	sidebarBorderStateColor := SidebarBorderColor
	if sidebarFocussed {
		sidebarBorderStateColor = SidebarBorderActiveColor
	}

	return lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(sidebarBorderStateColor).
		BorderBackground(SidebarBGColor).
		Width(Config.SidebarWidth).
		Height(height).
		Background(SidebarBGColor).
		Foreground(SidebarFGColor)
}

// Generate border style for process and can custom bottom border
func ProcsssBarBorder(height int, width int, borderBottom string, processBarFocussed bool) lipgloss.Style {
	border := GenerateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Processes " +
		Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
	border.Bottom = borderBottom

	processBorderStateColor := FooterBorderColor
	if processBarFocussed {
		processBorderStateColor = FooterBorderActiveColor
	}

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(processBorderStateColor).
		BorderBackground(FooterBGColor).
		Width(width).
		Height(height).
		Background(FooterBGColor).
		Foreground(FooterFGColor)
}

// Generate border style for metadata and can custom bottom border
func MetadataBorder(height int, width int, borderBottom string, metadataFocussed bool) lipgloss.Style {
	border := GenerateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Metadata " +
		Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
	border.Bottom = borderBottom

	metadataBorderStateColor := FooterBorderColor
	if metadataFocussed {
		metadataBorderStateColor = FooterBorderActiveColor
	}

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(metadataBorderStateColor).
		BorderBackground(FooterBGColor).
		Width(width).
		Height(height).
		Background(FooterBGColor).
		Foreground(FooterFGColor)
}

// Generate border style for clipboard and can custom bottom border
func ClipboardBorder(height int, width int, borderBottom string) lipgloss.Style {
	border := GenerateBorder()
	border.Top = Config.BorderTop + Config.BorderMiddleRight + " Clipboard " +
		Config.BorderMiddleLeft + strings.Repeat(Config.BorderTop, width)
	border.Bottom = borderBottom

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(FooterBorderColor).
		BorderBackground(FooterBGColor).
		Width(width).
		Height(height).
		Background(FooterBGColor).
		Foreground(FooterFGColor)
}

func ModalBorderStyle(height int, width int) lipgloss.Style {
	return modalBorderStyleWithAlign(height, width, lipgloss.Center)
}

func ModalBorderStyleLeft(height int, width int) lipgloss.Style {
	return modalBorderStyleWithAlign(height, width, lipgloss.Left)
}

// Generate modal (pop up widnwos) border style
func modalBorderStyleWithAlign(height int, width int, horizontalAlignment lipgloss.Position) lipgloss.Style {
	border := GenerateBorder()
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(horizontalAlignment, lipgloss.Center).
		Border(border).
		BorderForeground(ModalBorderActiveColor).
		BorderBackground(ModalBGColor).
		Background(ModalBGColor).
		Foreground(ModalFGColor)
}

// Generate first use modal style (This modal pop up when user first use superfile)
func FirstUseModal(height int, width int) lipgloss.Style {
	border := GenerateBorder()
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Left, lipgloss.Center).
		Border(border).
		BorderForeground(ModalBorderActiveColor).
		BorderBackground(ModalBGColor).
		Background(ModalBGColor).
		Foreground(ModalFGColor)
}

// Generate sort options modal border style
func SortOptionsModalBorderStyle(height int, width int, borderBottom string) lipgloss.Style {
	border := GenerateBorder()
	border.Bottom = borderBottom

	return lipgloss.NewStyle().
		Border(border).
		BorderForeground(ModalBorderActiveColor).
		BorderBackground(ModalBGColor).
		Width(width).
		Height(height).
		Background(ModalBGColor).
		Foreground(ModalFGColor)
}

// Generate full screen style for terminal size too small etc
func FullScreenStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Background(FullScreenBGColor).
		Foreground(FullScreenFGColor)
}

// Generate file panel divider style
func FilePanelDividerStyle(filePanelFocussed bool) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(FilePanelFocusColor(filePanelFocussed)).
		Background(FilePanelBGColor)
}

// Return border color based on file panel status
func FilePanelFocusColor(filePanelFocussed bool) lipgloss.Color {
	if filePanelFocussed {
		return FilePanelBorderActiveColor
	}
	return FilePanelBorderColor
}

// Return only fg and bg color style
func StringColorRender(fgColor lipgloss.Color, bgColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(fgColor).
		Background(bgColor)
}

// Generate border style
func GenerateBorder() lipgloss.Border {
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
func LoadConfigError(value string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("■ ERROR: ") +
		"Config file \"" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF")).Render(value) + "\" invalidation"
}

// Generate config error style
func LoadHotkeysError(value string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("■ ERROR: ") +
		"Hotkeys file \"" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF")).Render(value) + "\" invalidation"
}

// TODO : Fix Code duplication in textInput.Model creation
// This eventually caused a bug, where we created new model for sidebar search, and
// Didn't set `Width` in that. Take Width and other parameters as input in one function
// Generate search bar for file panel
func GenerateSearchBar() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = FooterCursorStyle
	ti.Cursor.TextStyle = FooterStyle
	ti.TextStyle = FilePanelStyle
	ti.Prompt = FilePanelTopDirectoryIconStyle.Render(icon.Search + icon.Space)
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
	//nolint:mnd // modal width minus padding
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
	ti.Width = Config.SidebarWidth - PanelPadding
	return ti
}

func GenerateDefaultProgress() progress.Model {
	prog := progress.New(GenerateGradientColor())
	prog.PercentageStyle = FooterStyle
	return prog
}

func GenerateGradientColor() progress.Option {
	return progress.WithScaledGradient(Theme.GradientColor[0], Theme.GradientColor[1])
}

func GenerateFooterBorder(countString string, width int) string {
	repeatCount := width - len(countString)
	if repeatCount < 0 {
		repeatCount = 0
	}
	return strings.Repeat(Config.BorderBottom, repeatCount) + Config.BorderMiddleRight +
		countString + Config.BorderMiddleLeft
}
