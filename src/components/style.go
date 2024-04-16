package components

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

var (
	minimumHeight       = 35
	minimumWidth        = 96
	footerHeight        = 14
	modalWidth          = 60
	modalHeight         = 7
	terminalTooSmall    lipgloss.Style
	terminalCorrectSize lipgloss.Style

	borderStyle      lipgloss.Style
	cursorStyle      lipgloss.Style
	mainBackgroundColor lipgloss.Color

	textStyle lipgloss.Style
)

var (
	sidebarWidth    = 20
	sidebarTitle    lipgloss.Style
	sidebarItem     lipgloss.Style
	sidebarSelected lipgloss.Style
)

var (
	filePanelTopFolderIcon lipgloss.Style
	filePanelTopPath       lipgloss.Style
	filePanelItem          lipgloss.Style
	filePanelItemSelected  lipgloss.Style
)

var (
	modalCancel  lipgloss.Style
	modalConfirm lipgloss.Style
)

func LoadThemeConfig() {
	filePanelBorderColor = lipgloss.Color(theme.FilePanelBorder)
	sidebarBorderColor = lipgloss.Color(theme.SidebarBorder)
	footerBorderColor = lipgloss.Color(theme.FooterBorder)
	modalBorderColor = lipgloss.Color(theme.ModalBorder)

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
	warnColor = lipgloss.Color(theme.Warn)
	
	filePanelTopDirectoryIconColor = lipgloss.Color(theme.FilePanelTopDirectoryIcon)
	filePanelTopPathColor = lipgloss.Color(theme.FilePanelTopPath)
	filePanelItemSelectedFGColor = lipgloss.Color(theme.FilePanelItemSelectedFG)
	filePanelItemSelectedBGColor = lipgloss.Color(theme.FilePanelItemSelectedBG)
	
	sidebarTitleColor = lipgloss.Color(theme.SidebarTitle)
	sidebarItemSelectedFGColor = lipgloss.Color(theme.SidebarItemSelectedFG)
	sidebarItemSelectedBGColor = lipgloss.Color(theme.SidebarItemSelectedBG)
	
	ModalCancelFGColor = lipgloss.Color(theme.ModalCancelFG)
	ModalCancelBGColor = lipgloss.Color(theme.ModalCancelBG)
	ModalConfirmFGColor = lipgloss.Color(theme.ModalConfirmFG)
	ModalConfirmBGColor = lipgloss.Color(theme.ModalConfirmBG)

	
	mainBackgroundColor = lipgloss.Color(theme.MainBackground)

	terminalTooSmall = lipgloss.NewStyle().Foreground(errorColor).Background(fullScreenBGColor)
	terminalCorrectSize = lipgloss.NewStyle().Foreground(cursorColor).Background(fullScreenBGColor)

	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border)).Background(mainBackgroundColor)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Cursor)).Bold(true).Background(mainBackgroundColor)

	sidebarTitle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SidebarTitle)).Bold(true).Background(mainBackgroundColor)
	sidebarItem = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SidebarItem)).Background(mainBackgroundColor)

	filePanelTopFolderIcon = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelTopDirectoryIcon)).Bold(true)
	filePanelTopPath = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelTopPath)).Bold(true).Background(mainBackgroundColor)
	filePanelItem = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelItem)).Background(mainBackgroundColor)
	filePanelItemSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelItemSelected)).Background(mainBackgroundColor)

	sidebarSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SidebarSelected)).Background(mainBackgroundColor)

	modalCancel = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ModalForeground)).Background(lipgloss.Color(theme.ModalCancel))
	modalConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ModalForeground)).Background(lipgloss.Color(theme.ModalConfirm))

	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ModalForeground)).Background(mainBackgroundColor)
}

func FullScreenStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).Width(width).Align(lipgloss.Center, lipgloss.Center).Background(mainBackgroundColor)
}

func FocusedModalStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(theme.FilePanelFocus)).BorderBackground(mainBackgroundColor).Background(mainBackgroundColor)
}

func SideBarBoardStyle(height int, focus focusPanelType) lipgloss.Style {
	if focus == sidebarFocus {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color(theme.SidebarFocus)).
			BorderBackground(mainBackgroundColor).
			Width(sidebarWidth).
			Height(height).Bold(true).Background(mainBackgroundColor)
	} else {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder()).
			BorderBackground(mainBackgroundColor).
			Width(sidebarWidth).
			Height(height).Bold(true).Background(mainBackgroundColor)
	}
}

func FilePanelBoardStyle(height int, width int, focusType filePanelFocusType, borderBottom string) lipgloss.Style {
	leftBorder := ""
	rightBorder := ""
	for i := 0; i < height; i++ {
		if i == 1 {
			leftBorder += "┣"
			rightBorder += "┫"
		} else {
			leftBorder += "┃"
			rightBorder += "┃"
		}
	}
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━",
		Bottom:      borderBottom,
		Left:        leftBorder,
		Right:       rightBorder,
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}
	return lipgloss.NewStyle().
		Border(filePanelBottomBoard, true, true, true, true).
		BorderForeground(lipgloss.Color(FilePanelFocusColor(focusType))).
		BorderBackground(mainBackgroundColor).
		Width(width).
		Height(height).Background(mainBackgroundColor)
}

func ProcsssBarBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━┫Processes┣" + strings.Repeat("━", width),
		Bottom:      borderBottom,
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}
	if focusType == processBarFocus {
		return lipgloss.NewStyle().
			Border(filePanelBottomBoard, true, true, true, true).
			BorderForeground(lipgloss.Color(theme.FooterFocus)).
			BorderBackground(mainBackgroundColor).
			Width(width).
			Height(height).Bold(true).Background(mainBackgroundColor)
	} else {
		return lipgloss.NewStyle().
			Border(filePanelBottomBoard, true, true, true, true).
			BorderForeground(lipgloss.Color(theme.Border)).
			BorderBackground(mainBackgroundColor).
			Width(width).
			Height(height).Bold(true).Background(mainBackgroundColor)
	}
}

func MetaDataBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━┫Metadata┣" + strings.Repeat("━", width),
		Bottom:      borderBottom,
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}
	if focusType == metaDataFocus {
		return lipgloss.NewStyle().
			Border(filePanelBottomBoard, true, true, true, true).
			BorderForeground(lipgloss.Color(theme.FooterFocus)).
			BorderBackground(mainBackgroundColor).
			Width(width).
			Height(height).Bold(true).Background(mainBackgroundColor)
	} else {
		return lipgloss.NewStyle().
			Border(filePanelBottomBoard, true, true, true, true).
			BorderForeground(lipgloss.Color(theme.Border)).
			BorderBackground(mainBackgroundColor).
			Width(width).
			Height(height).Bold(true).Background(mainBackgroundColor)
	}
}

func ClipboardBoarder(height int, width int, borderBottom string) lipgloss.Style {
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━┫Clipboard┣" + strings.Repeat("━", width),
		Bottom:      borderBottom,
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	return lipgloss.NewStyle().
		Border(filePanelBottomBoard, true, true, true, true).
		BorderForeground(lipgloss.Color(theme.Border)).
		BorderBackground(mainBackgroundColor).
		Width(width).
		Height(height).Bold(true).Background(mainBackgroundColor)

}

func FilePanelDividerStyle(focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(FilePanelFocusColor(focusType))).Bold(true).Background(mainBackgroundColor)
}

func TruncateText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	return text[:maxChars-3] + "..."
}

func TruncateTextBeginning(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	charsToKeep := maxChars - 3
	truncatedRunes := append([]rune("..."), runes[len(runes)-charsToKeep:]...)
	return string(truncatedRunes)
}

func TruncateMiddleText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	halfEllipsisLength := (maxChars - 3) / 2

	truncatedText := text[:halfEllipsisLength] + "..." + text[utf8.RuneCountInString(text)-halfEllipsisLength:]

	return truncatedText
}

func PrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return StringColorRender(style.color).Background(mainBackgroundColor).Render(style.icon) + lipgloss.NewStyle().Background(mainBackgroundColor).Render("  ") + filePanelItemSelected.Render(TruncateText(name, width))
	} else {
		return StringColorRender(style.color).Background(mainBackgroundColor).Render(style.icon) + lipgloss.NewStyle().Background(mainBackgroundColor).Render("  ") + filePanelItem.Render(TruncateText(name, width))
	}
}

func ClipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return StringColorRender(style.color).Background(mainBackgroundColor).Render(style.icon) + lipgloss.NewStyle().Background(mainBackgroundColor).Render("  ") + filePanelItemSelected.Render(TruncateTextBeginning(name, width))
	} else {
		return StringColorRender(style.color).Background(mainBackgroundColor).Render(style.icon) + lipgloss.NewStyle().Background(mainBackgroundColor).Render("  ") + filePanelItem.Render(TruncateTextBeginning(name, width))
	}
}

// CHOOSE STYLE FUNCTION
func FilePanelFocusColor(focusType filePanelFocusType) string {
	if focusType == noneFocus {
		return theme.Border
	} else {
		return theme.FilePanelFocus
	}
}

func FilePanelBoard(focusType filePanelFocusType) lipgloss.Border {
	if focusType == noneFocus {
		return lipgloss.RoundedBorder()
	} else {
		return lipgloss.ThickBorder()
	}
}

func GenerateBottomBorder(countString string, width int) string {
	return strings.Repeat("━", width-len(countString)) + "┫" + countString + "┣"
}

func StringColorRender(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Background(mainBackgroundColor)
}

func BottomWidth(fullWidth int) int {
	return fullWidth/3 - 2
}
