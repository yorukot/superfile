package common

import (
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
	MainStyle      lipgloss.Style
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
	ModalErrorStyle lipgloss.Style
)

var (
	HelpMenuHotkeyStyle lipgloss.Style
	HelpMenuTitleStyle  lipgloss.Style
)

var (
	PromptSuccessStyle lipgloss.Style
	PromptFailureStyle lipgloss.Style
)
var TransparentBackgroundColor string

var (
	FilePanelBorderColor lipgloss.Color
	SidebarBorderColor   lipgloss.Color
	FooterBorderColor    lipgloss.Color

	FilePanelBorderActiveColor lipgloss.Color
	SidebarBorderActiveColor   lipgloss.Color
	FooterBorderActiveColor    lipgloss.Color
	ModalBorderActiveColor     lipgloss.Color

	FullScreenBGColor lipgloss.Color
	FilePanelBGColor  lipgloss.Color
	SidebarBGColor    lipgloss.Color
	FooterBGColor     lipgloss.Color
	ModalBGColor      lipgloss.Color

	FullScreenFGColor lipgloss.Color
	FilePanelFGColor  lipgloss.Color
	SidebarFGColor    lipgloss.Color
	FooterFGColor     lipgloss.Color
	ModalFGColor      lipgloss.Color

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

	promptSuccessColor lipgloss.Color
	promptFailureColor lipgloss.Color
)

func LoadThemeConfig() { //nolint: funlen // Variable initialization
	BottomMiddleBorderSplit = Config.BorderMiddleLeft + Config.BorderBottom + Config.BorderMiddleRight

	FilePanelBorderColor = lipgloss.Color(Theme.FilePanelBorder)
	SidebarBorderColor = lipgloss.Color(Theme.SidebarBorder)
	FooterBorderColor = lipgloss.Color(Theme.FooterBorder)

	FilePanelBorderActiveColor = lipgloss.Color(Theme.FilePanelBorderActive)
	SidebarBorderActiveColor = lipgloss.Color(Theme.SidebarBorderActive)
	FooterBorderActiveColor = lipgloss.Color(Theme.FooterBorderActive)
	ModalBorderActiveColor = lipgloss.Color(Theme.ModalBorderActive)

	FullScreenBGColor = lipgloss.Color(Theme.FullScreenBG)
	FilePanelBGColor = lipgloss.Color(Theme.FilePanelBG)
	SidebarBGColor = lipgloss.Color(Theme.SidebarBG)
	FooterBGColor = lipgloss.Color(Theme.FooterBG)
	ModalBGColor = lipgloss.Color(Theme.ModalBG)

	FullScreenFGColor = lipgloss.Color(Theme.FullScreenFG)
	FilePanelFGColor = lipgloss.Color(Theme.FilePanelFG)
	SidebarFGColor = lipgloss.Color(Theme.SidebarFG)
	FooterFGColor = lipgloss.Color(Theme.FooterFG)
	ModalFGColor = lipgloss.Color(Theme.ModalFG)

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

	promptSuccessColor = lipgloss.Color(Theme.Correct)
	promptFailureColor = lipgloss.Color(Theme.Error)

	if Config.TransparentBackground {
		TransparentAllBackgroundColor()
	}

	// All Panel Main Color
	// (full screen and default color)
	MainStyle = lipgloss.NewStyle().Foreground(FullScreenFGColor).Background(FullScreenBGColor)
	FilePanelStyle = lipgloss.NewStyle().Foreground(FilePanelFGColor).Background(FilePanelBGColor)
	SidebarStyle = lipgloss.NewStyle().Foreground(SidebarFGColor).Background(SidebarBGColor)
	FooterStyle = lipgloss.NewStyle().Foreground(FooterFGColor).Background(FooterBGColor)
	ModalStyle = lipgloss.NewStyle().Foreground(ModalFGColor).Background(ModalBGColor)

	// Terminal Size Error
	TerminalTooSmall = lipgloss.NewStyle().Foreground(errorColor).Background(FullScreenBGColor)
	TerminalCorrectSize = lipgloss.NewStyle().Foreground(cursorColor).Background(FullScreenBGColor)

	// Cursor
	FilePanelCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(FilePanelBGColor)
	FooterCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(FooterBGColor)
	ModalCursorStyle = lipgloss.NewStyle().Foreground(cursorColor).Background(ModalBGColor)

	// File Panel Special Style
	FilePanelTopDirectoryIconStyle = lipgloss.NewStyle().Foreground(filePanelTopDirectoryIconColor).
		Background(FilePanelBGColor)
	FilePanelTopPathStyle = lipgloss.NewStyle().Foreground(filePanelTopPathColor).Background(FilePanelBGColor)
	FilePanelItemSelectedStyle = lipgloss.NewStyle().Foreground(filePanelItemSelectedFGColor).
		Background(filePanelItemSelectedBGColor)

	// Sidebar Special Style
	SidebarDividerStyle = lipgloss.NewStyle().Foreground(sidebarDividerColor).Background(SidebarBGColor)
	SidebarTitleStyle = lipgloss.NewStyle().Foreground(sidebarTitleColor).Background(SidebarBGColor)
	SidebarSelectedStyle = lipgloss.NewStyle().Foreground(sidebarItemSelectedFGColor).
		Background(sidebarItemSelectedBGColor)

	// Footer Special Style
	ProcessErrorStyle = lipgloss.NewStyle().Foreground(errorColor).Background(FooterBGColor)
	ProcessInOperationStyle = lipgloss.NewStyle().Foreground(hintColor).Background(FooterBGColor)
	ProcessCancelStyle = lipgloss.NewStyle().Foreground(cancelColor).Background(FooterBGColor)
	ProcessSuccessfulStyle = lipgloss.NewStyle().Foreground(correctColor).Background(FooterBGColor)

	// Modal Special Style
	ModalCancel = lipgloss.NewStyle().Foreground(modalCancelFGColor).Background(modalCancelBGColor)
	ModalConfirm = lipgloss.NewStyle().Foreground(modalConfirmFGColor).Background(modalConfirmBGColor)
	ModalTitleStyle = lipgloss.NewStyle().Foreground(hintColor).Background(ModalBGColor)
	ModalErrorStyle = lipgloss.NewStyle().Foreground(errorColor).Background(ModalBGColor)
	// Help Menu Style
	HelpMenuHotkeyStyle = lipgloss.NewStyle().Foreground(helpMenuHotkeyColor).Background(ModalBGColor)
	HelpMenuTitleStyle = lipgloss.NewStyle().Foreground(helpMenuTitleColor).Background(ModalBGColor)

	// Prompt Style
	PromptSuccessStyle = lipgloss.NewStyle().Foreground(promptSuccessColor).Background(ModalBGColor)
	PromptFailureStyle = lipgloss.NewStyle().Foreground(promptFailureColor).Background(ModalBGColor)
}

func TransparentAllBackgroundColor() {
	if SidebarBGColor == sidebarItemSelectedBGColor {
		sidebarItemSelectedBGColor = lipgloss.Color(TransparentBackgroundColor)
	}

	if FilePanelBGColor == filePanelItemSelectedBGColor {
		filePanelItemSelectedBGColor = lipgloss.Color(TransparentBackgroundColor)
	}

	FullScreenBGColor = lipgloss.Color(TransparentBackgroundColor)
	FilePanelBGColor = lipgloss.Color(TransparentBackgroundColor)
	SidebarBGColor = lipgloss.Color(TransparentBackgroundColor)
	FooterBGColor = lipgloss.Color(TransparentBackgroundColor)
	ModalBGColor = lipgloss.Color(TransparentBackgroundColor)
}
