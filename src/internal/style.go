package internal

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

var (
	minimumHeight = 24
	minimumWidth  = 80
	footerHeight  = 14
	modalWidth    = 60
	modalHeight   = 7
	sidebarWidth  = 20
)

var (
	terminalTooSmall    lipgloss.Style
	terminalCorrectSize lipgloss.Style
)

var (
	mainStyle      lipgloss.Style
	filePanelStyle lipgloss.Style
	sidebarStyle   lipgloss.Style
	footerStyle    lipgloss.Style
	modalStyle     lipgloss.Style
)

var (
	sidebarDividerStyle  lipgloss.Style
	sidebarTitleStyle    lipgloss.Style
	sidebarSelectedStyle lipgloss.Style
)

var (
	filePanelCursorStyle lipgloss.Style
	footerCursorStyle    lipgloss.Style
	modalCursorStyle     lipgloss.Style
)

var (
	filePanelTopDirectoryIconStyle lipgloss.Style
	filePanelTopPathStyle          lipgloss.Style
	filePanelItemSelectedStyle     lipgloss.Style
)

var (
	processErrorStyle       lipgloss.Style
	processInOperationStyle lipgloss.Style
	processCancelStyle      lipgloss.Style
	processSuccessfulStyle  lipgloss.Style
)

var (
	modalCancel  lipgloss.Style
	modalConfirm lipgloss.Style
)

var (
	helpMenuHotkeyStyle lipgloss.Style
	helpMenuTitleStyle  lipgloss.Style
)

func LoadThemeConfig() {
	filePanelBorderColor = lipgloss.Color(theme.FilePanelBorder)
	sidebarBorderColor = lipgloss.Color(theme.SidebarBorder)
	footerBorderColor = lipgloss.Color(theme.FooterBorder)

	filePanelBorderActiveColor = lipgloss.Color(theme.FilePanelBorderActive)
	sidebarBorderActiveColor = lipgloss.Color(theme.SidebarBorderActive)
	footerBorderActiveColor = lipgloss.Color(theme.FooterBorderActive)
	modalBorderActiveColor = lipgloss.Color(theme.ModalBorderActive)

	fullScreenBGColor = lipgloss.Color(theme.FullScreenBG)
	filePanelBGColor = lipgloss.Color(theme.FilePanelBG)
	sidebarBGColor = lipgloss.Color(theme.SidebarBG)
	footerBGColor = lipgloss.Color(theme.FooterBG)
	modalBGColor = lipgloss.Color(theme.ModalBG)

	fullScreenFGColor = lipgloss.Color(theme.FullScreenFG)
	filePanelFGColor = lipgloss.Color(theme.FilePanelFG)
	sidebarFGColor = lipgloss.Color(theme.SidebarFG)
	footerFGColor = lipgloss.Color(theme.FooterFG)
	modalFGColor = lipgloss.Color(theme.ModalFG)

	cursorColor = lipgloss.Color(theme.Cursor)
	correctColor = lipgloss.Color(theme.Correct)
	errorColor = lipgloss.Color(theme.Error)
	hintColor = lipgloss.Color(theme.Hint)
	cancelColor = lipgloss.Color(theme.Cancel)

	filePanelTopDirectoryIconColor = lipgloss.Color(theme.FilePanelTopDirectoryIcon)
	filePanelTopPathColor = lipgloss.Color(theme.FilePanelTopPath)
	filePanelItemSelectedFGColor = lipgloss.Color(theme.FilePanelItemSelectedFG)
	filePanelItemSelectedBGColor = lipgloss.Color(theme.FilePanelItemSelectedBG)

	sidebarTitleColor = lipgloss.Color(theme.SidebarTitle)
	sidebarItemSelectedFGColor = lipgloss.Color(theme.SidebarItemSelectedFG)
	sidebarItemSelectedBGColor = lipgloss.Color(theme.SidebarItemSelectedBG)
	sidebarDividerColor = lipgloss.Color(theme.SidebarDivider)

	modalCancelFGColor = lipgloss.Color(theme.ModalCancelFG)
	modalCancelBGColor = lipgloss.Color(theme.ModalCancelBG)
	modalConfirmFGColor = lipgloss.Color(theme.ModalConfirmFG)
	modalConfirmBGColor = lipgloss.Color(theme.ModalConfirmBG)

	helpMenuHotkeyColor = lipgloss.Color(theme.HelpMenuHotkey)
	helpMenuTitleColor = lipgloss.Color(theme.HelpMenuTitle)

	// All Panel Main Color
	// (full screen and default color)
	mainStyle = lipgloss.NewStyle().Foreground(fullScreenFGColor).Background(fullScreenBGColor)
	filePanelStyle = lipgloss.NewStyle().Foreground(filePanelFGColor).Background(filePanelBGColor)
	sidebarStyle = lipgloss.NewStyle().Foreground(sidebarFGColor).Background(sidebarBGColor)
	footerStyle = lipgloss.NewStyle().Foreground(footerFGColor).Background(footerBGColor)
	modalStyle = lipgloss.NewStyle().Foreground(modalFGColor).Background(modalBGColor)

	// Terminal Size Error
	terminalTooSmall = lipgloss.NewStyle().Foreground(errorColor).Background(fullScreenBGColor)
	terminalCorrectSize = lipgloss.NewStyle().Foreground(cursorColor).Background(fullScreenBGColor)

	// Cursor
	filePanelCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(filePanelBGColor)
	footerCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(footerBGColor)
	modalCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(modalBGColor)

	// File Panel Special Style
	filePanelTopDirectoryIconStyle = lipgloss.NewStyle().Foreground(filePanelTopDirectoryIconColor).Background(filePanelBGColor)
	filePanelTopPathStyle = lipgloss.NewStyle().Foreground(filePanelTopPathColor).Background(filePanelBGColor)
	filePanelItemSelectedStyle = lipgloss.NewStyle().Foreground(filePanelItemSelectedFGColor).Background(filePanelItemSelectedBGColor)

	// Sidebar Special Style
	sidebarDividerStyle = lipgloss.NewStyle().Foreground(sidebarDividerColor).Background(sidebarBGColor)
	sidebarTitleStyle = lipgloss.NewStyle().Foreground(sidebarTitleColor).Background(sidebarBGColor)
	sidebarSelectedStyle = lipgloss.NewStyle().Foreground(sidebarItemSelectedFGColor).Background(sidebarItemSelectedBGColor)

	// Footer Special Style
	processErrorStyle = lipgloss.NewStyle().Foreground(errorColor).Background(footerBGColor)
	processInOperationStyle = lipgloss.NewStyle().Foreground(hintColor).Background(footerBGColor)
	processCancelStyle = lipgloss.NewStyle().Foreground(cancelColor).Background(footerBGColor)
	processSuccessfulStyle = lipgloss.NewStyle().Foreground(correctColor).Background(footerBGColor)

	// Modal Special Style
	modalCancel = lipgloss.NewStyle().Foreground(modalCancelFGColor).Background(modalCancelBGColor)
	modalConfirm = lipgloss.NewStyle().Foreground(modalConfirmFGColor).Background(modalConfirmBGColor)

	// Help Menu Style
	helpMenuHotkeyStyle = lipgloss.NewStyle().Foreground(helpMenuHotkeyColor).Background(modalBGColor)
	helpMenuTitleStyle = lipgloss.NewStyle().Foreground(helpMenuTitleColor).Background(modalBGColor)
}

func generateGradientColor() progress.Option {
	return progress.WithScaledGradient(theme.GradientColor[0], theme.GradientColor[1])
}

func generateFooterBorder(countString string, width int) string {
	return strings.Repeat("━", width-len(countString)) + "┫" + countString + "┣"
}

func footerWidth(fullWidth int) int {
	return fullWidth/3 - 2
}
