package common

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
)

func ModalBorderStyle(height int, width int) lipgloss.Style {
	return modalBorderStyleWithAlign(height, width, lipgloss.Center)
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

func LoadConfigError(value string, msg string) string {
	return UserConfigInvalidationErrorString(value, "Config", msg)
}

func LoadHotkeysError(value string, msg string) string {
	return UserConfigInvalidationErrorString(value, "Hotkey", msg)
}

func LoadThemeError(value string, msg string) string {
	return UserConfigInvalidationErrorString(value, "Theme", msg)
}

func UserConfigInvalidationErrorString(value string, configType string, msg string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("■ ERROR: ") +
		configType + " value for \"" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF")).Render(value) +
		"\" is invalid : " + msg
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
