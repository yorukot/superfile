package common

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

var (
	BottomMiddleBorderSplit string
)
var (
	TerminalTooSmall    lipgloss.Style
	TerminalCorrectSize lipgloss.Style
)

var (
	mainStyle      lipgloss.Style
	FilePanelStyle lipgloss.Style
	SidebarStyle   lipgloss.Style
	FooterStyle    lipgloss.Style
	ModalStyle     lipgloss.Style
)

var (
	SidebarDividerStyle  lipgloss.Style
	SidebarTitleStyle    lipgloss.Style
	SidebarSelectedStyle lipgloss.Style
)

var (
	FilePanelCursorStyle lipgloss.Style
	FooterCursorStyle    lipgloss.Style
	ModalCursorStyle     lipgloss.Style
)

var (
	FilePanelTopDirectoryIconStyle lipgloss.Style
	FilePanelTopPathStyle          lipgloss.Style
	FilePanelItemSelectedStyle     lipgloss.Style
)

var (
	ProcessErrorStyle       lipgloss.Style
	ProcessInOperationStyle lipgloss.Style
	ProcessCancelStyle      lipgloss.Style
	ProcessSuccessfulStyle  lipgloss.Style
)

var (
	ModalCancel     lipgloss.Style
	ModalConfirm    lipgloss.Style
	ModalTitleStyle lipgloss.Style
)

var (
	HelpMenuHotkeyStyle lipgloss.Style
	HelpMenuTitleStyle  lipgloss.Style
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
	FooterBGColor     lipgloss.Color
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
	BottomMiddleBorderSplit = Config.BorderMiddleLeft + Config.BorderBottom + Config.BorderMiddleRight

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
	FooterBGColor = lipgloss.Color(Theme.FooterBG)
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
	FilePanelStyle = lipgloss.NewStyle().Foreground(filePanelFGColor).Background(filePanelBGColor)
	SidebarStyle = lipgloss.NewStyle().Foreground(sidebarFGColor).Background(sidebarBGColor)
	FooterStyle = lipgloss.NewStyle().Foreground(footerFGColor).Background(FooterBGColor)
	ModalStyle = lipgloss.NewStyle().Foreground(modalFGColor).Background(modalBGColor)

	// Terminal Size Error
	TerminalTooSmall = lipgloss.NewStyle().Foreground(errorColor).Background(fullScreenBGColor)
	TerminalCorrectSize = lipgloss.NewStyle().Foreground(cursorColor).Background(fullScreenBGColor)

	// Cursor
	FilePanelCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(filePanelBGColor)
	FooterCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(FooterBGColor)
	ModalCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(modalBGColor)

	// File Panel Special Style
	FilePanelTopDirectoryIconStyle = lipgloss.NewStyle().Foreground(filePanelTopDirectoryIconColor).Background(filePanelBGColor)
	FilePanelTopPathStyle = lipgloss.NewStyle().Foreground(filePanelTopPathColor).Background(filePanelBGColor)
	FilePanelItemSelectedStyle = lipgloss.NewStyle().Foreground(filePanelItemSelectedFGColor).Background(filePanelItemSelectedBGColor)

	// Sidebar Special Style
	SidebarDividerStyle = lipgloss.NewStyle().Foreground(sidebarDividerColor).Background(sidebarBGColor)
	SidebarTitleStyle = lipgloss.NewStyle().Foreground(sidebarTitleColor).Background(sidebarBGColor)
	SidebarSelectedStyle = lipgloss.NewStyle().Foreground(sidebarItemSelectedFGColor).Background(sidebarItemSelectedBGColor)

	// Footer Special Style
	ProcessErrorStyle = lipgloss.NewStyle().Foreground(errorColor).Background(FooterBGColor)
	ProcessInOperationStyle = lipgloss.NewStyle().Foreground(hintColor).Background(FooterBGColor)
	ProcessCancelStyle = lipgloss.NewStyle().Foreground(cancelColor).Background(FooterBGColor)
	ProcessSuccessfulStyle = lipgloss.NewStyle().Foreground(correctColor).Background(FooterBGColor)

	// Modal Special Style
	ModalCancel = lipgloss.NewStyle().Foreground(modalCancelFGColor).Background(modalCancelBGColor)
	ModalConfirm = lipgloss.NewStyle().Foreground(modalConfirmFGColor).Background(modalConfirmBGColor)
	ModalTitleStyle = lipgloss.NewStyle().Foreground(hintColor).Background(modalBGColor)

	// Help Menu Style
	HelpMenuHotkeyStyle = lipgloss.NewStyle().Foreground(helpMenuHotkeyColor).Background(modalBGColor)
	HelpMenuTitleStyle = lipgloss.NewStyle().Foreground(helpMenuTitleColor).Background(modalBGColor)
}

func GenerateGradientColor() progress.Option {
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
	FooterBGColor = lipgloss.Color(transparentBackgroundColor)
	modalBGColor = lipgloss.Color(transparentBackgroundColor)
}
