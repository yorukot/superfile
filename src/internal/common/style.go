package common

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

var (
	minimumHeight = 24
	minimumWidth  = 60

	// Todo : These are model object properties, not global properties
	// We are modifying them in the code many time. They need to be part of model struct.
	minFooterHeight = 6
	modalWidth      = 60
	modalHeight     = 7
)

var (
	bottomMiddleBorderSplit string
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
	modalCancel     lipgloss.Style
	modalConfirm    lipgloss.Style
	modalTitleStyle lipgloss.Style
)

var (
	helpMenuHotkeyStyle lipgloss.Style
	helpMenuTitleStyle  lipgloss.Style
)

var (
	filePanelBorderColor lipgloss.Color
	sidebarBorderColor   lipgloss.Color
	footerBorderColor    lipgloss.Color

	filePanelBorderActiveColor lipgloss.Color
	sidebarBorderActiveColor   lipgloss.Color
	footerBorderActiveColor    lipgloss.Color
	modalBorderActiveColor     lipgloss.Color

	fullScreenBGColor lipgloss.Color
	filePanelBGColor  lipgloss.Color
	sidebarBGColor    lipgloss.Color
	footerBGColor     lipgloss.Color
	modalBGColor      lipgloss.Color

	fullScreenFGColor lipgloss.Color
	filePanelFGColor  lipgloss.Color
	sidebarFGColor    lipgloss.Color
	footerFGColor     lipgloss.Color
	modalFGColor      lipgloss.Color

	cursorColor  lipgloss.Color
	correctColor lipgloss.Color
	errorColor   lipgloss.Color
	hintColor    lipgloss.Color
	cancelColor  lipgloss.Color

	filePanelTopDirectoryIconColor lipgloss.Color
	filePanelTopPathColor          lipgloss.Color
	filePanelItemSelectedFGColor   lipgloss.Color
	filePanelItemSelectedBGColor   lipgloss.Color

	sidebarTitleColor          lipgloss.Color
	sidebarItemSelectedFGColor lipgloss.Color
	sidebarItemSelectedBGColor lipgloss.Color
	sidebarDividerColor        lipgloss.Color

	modalCancelFGColor  lipgloss.Color
	modalCancelBGColor  lipgloss.Color
	modalConfirmFGColor lipgloss.Color
	modalConfirmBGColor lipgloss.Color

	helpMenuHotkeyColor lipgloss.Color
	helpMenuTitleColor  lipgloss.Color
)

func LoadThemeConfig() {
	bottomMiddleBorderSplit = Config.BorderMiddleLeft + Config.BorderBottom + Config.BorderMiddleRight

	filePanelBorderColor = lipgloss.Color(Theme.FilePanelBorder)
	sidebarBorderColor = lipgloss.Color(Theme.SidebarBorder)
	footerBorderColor = lipgloss.Color(Theme.FooterBorder)

	filePanelBorderActiveColor = lipgloss.Color(Theme.FilePanelBorderActive)
	sidebarBorderActiveColor = lipgloss.Color(Theme.SidebarBorderActive)
	footerBorderActiveColor = lipgloss.Color(Theme.FooterBorderActive)
	modalBorderActiveColor = lipgloss.Color(Theme.ModalBorderActive)

	fullScreenBGColor = lipgloss.Color(Theme.FullScreenBG)
	filePanelBGColor = lipgloss.Color(Theme.FilePanelBG)
	sidebarBGColor = lipgloss.Color(Theme.SidebarBG)
	footerBGColor = lipgloss.Color(Theme.FooterBG)
	modalBGColor = lipgloss.Color(Theme.ModalBG)

	fullScreenFGColor = lipgloss.Color(Theme.FullScreenFG)
	filePanelFGColor = lipgloss.Color(Theme.FilePanelFG)
	sidebarFGColor = lipgloss.Color(Theme.SidebarFG)
	footerFGColor = lipgloss.Color(Theme.FooterFG)
	modalFGColor = lipgloss.Color(Theme.ModalFG)

	cursorColor = lipgloss.Color(Theme.Cursor)
	correctColor = lipgloss.Color(Theme.Correct)
	errorColor = lipgloss.Color(Theme.Error)
	hintColor = lipgloss.Color(Theme.Hint)
	cancelColor = lipgloss.Color(Theme.Cancel)

	filePanelTopDirectoryIconColor = lipgloss.Color(Theme.FilePanelTopDirectoryIcon)
	filePanelTopPathColor = lipgloss.Color(Theme.FilePanelTopPath)
	filePanelItemSelectedFGColor = lipgloss.Color(Theme.FilePanelItemSelectedFG)
	filePanelItemSelectedBGColor = lipgloss.Color(Theme.FilePanelItemSelectedBG)

	sidebarTitleColor = lipgloss.Color(Theme.SidebarTitle)
	sidebarItemSelectedFGColor = lipgloss.Color(Theme.SidebarItemSelectedFG)
	sidebarItemSelectedBGColor = lipgloss.Color(Theme.SidebarItemSelectedBG)
	sidebarDividerColor = lipgloss.Color(Theme.SidebarDivider)

	modalCancelFGColor = lipgloss.Color(Theme.ModalCancelFG)
	modalCancelBGColor = lipgloss.Color(Theme.ModalCancelBG)
	modalConfirmFGColor = lipgloss.Color(Theme.ModalConfirmFG)
	modalConfirmBGColor = lipgloss.Color(Theme.ModalConfirmBG)

	helpMenuHotkeyColor = lipgloss.Color(Theme.HelpMenuHotkey)
	helpMenuTitleColor = lipgloss.Color(Theme.HelpMenuTitle)

	if Config.TransparentBackground {
		transparentAllBackgroundColor()
	}

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
	modalTitleStyle = lipgloss.NewStyle().Foreground(hintColor).Background(modalBGColor)

	// Help Menu Style
	helpMenuHotkeyStyle = lipgloss.NewStyle().Foreground(helpMenuHotkeyColor).Background(modalBGColor)
	helpMenuTitleStyle = lipgloss.NewStyle().Foreground(helpMenuTitleColor).Background(modalBGColor)
}

func generateGradientColor() progress.Option {
	return progress.WithScaledGradient(Theme.GradientColor[0], Theme.GradientColor[1])
}

func generateFooterBorder(countString string, width int) string {
	repeatCount := width - len(countString)
	if repeatCount < 0 {
		repeatCount = 0
	}
	return strings.Repeat(Config.BorderBottom, repeatCount) + Config.BorderMiddleRight + countString + Config.BorderMiddleLeft
}

func footerWidth(fullWidth int) int {
	return fullWidth/3 - 2
}

var transparentBackgroundColor string

func transparentAllBackgroundColor() {

	if sidebarBGColor == sidebarItemSelectedBGColor {
		sidebarItemSelectedBGColor = lipgloss.Color(transparentBackgroundColor)
	}

	if filePanelBGColor == filePanelItemSelectedBGColor {
		filePanelItemSelectedBGColor = lipgloss.Color(transparentBackgroundColor)
	}

	fullScreenBGColor = lipgloss.Color(transparentBackgroundColor)
	filePanelBGColor = lipgloss.Color(transparentBackgroundColor)
	sidebarBGColor = lipgloss.Color(transparentBackgroundColor)
	footerBGColor = lipgloss.Color(transparentBackgroundColor)
	modalBGColor = lipgloss.Color(transparentBackgroundColor)
}
