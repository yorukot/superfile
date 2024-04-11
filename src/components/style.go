package components

import (
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

var (
	minimumHeight       = 35
	minimumWidth        = 96
	bottomBarHeight     = 14
	modalWidth          = 60
	modalHeight         = 7
	terminalTooSmall    lipgloss.Style
	terminalMinimumSize lipgloss.Style

	borderStyle         lipgloss.Style
	cursorStyle         lipgloss.Style
	backgroundWindow    lipgloss.Color
)

var (
	sideBarWidth        = 20
	sideBarTitle        lipgloss.Style
	sideBarItem         lipgloss.Style
	sideBarSelected     lipgloss.Style
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
	backgroundWindow = lipgloss.Color(theme.BackgroundWindow)

	terminalTooSmall = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TerminalTooSmallError)).Background(backgroundWindow)
	terminalMinimumSize = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TerminalSizeCorrect)).Background(backgroundWindow)

	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border)).Background(backgroundWindow)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Cursor)).Bold(true).Background(backgroundWindow)

	sideBarTitle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SideBarTitle)).Bold(true).Background(backgroundWindow)
	sideBarItem = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SideBarItem)).Background(backgroundWindow)

	filePanelTopFolderIcon = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelTopDirectoryIcon)).Bold(true)
	filePanelTopPath = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelTopPath)).Bold(true).Background(backgroundWindow)
	filePanelItem = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelItem)).Background(backgroundWindow)
	filePanelItemSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FilePanelItemSelected)).Background(backgroundWindow)

	sideBarSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SideBarSelected)).Background(backgroundWindow)

	modalCancel = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ModalForeground)).Background(lipgloss.Color(theme.ModalCancel))
	modalConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ModalForeground)).Background(lipgloss.Color(theme.ModalConfirm))
}

func FullScreenStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).Width(width).Align(lipgloss.Center, lipgloss.Center).Background(backgroundWindow)
}

func FocusedModalStyle(height int, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(theme.FilePanelFocus)).BorderBackground(backgroundWindow).Background(backgroundWindow)
}

func SideBarBoardStyle(height int, focus focusPanelType) lipgloss.Style {
	if focus == sideBarFocus {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color(theme.SideBarFocus)).
			BorderBackground(backgroundWindow).
			Width(sideBarWidth).
			Height(height).Bold(true).Background(backgroundWindow)
	} else {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder()).
			BorderBackground(backgroundWindow).
			Width(sideBarWidth).
			Height(height).Bold(true).Background(backgroundWindow)
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
		BorderBackground(backgroundWindow).
		Width(width).
		Height(height).Background(backgroundWindow)
}

func ProcsssBarBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━┫Processes┣" + repeatString("━", width),
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
			BorderForeground(lipgloss.Color(theme.BottomBarFocus)).
			BorderBackground(backgroundWindow).
			Width(width).
			Height(height).Bold(true).Background(backgroundWindow)
	} else {
		return lipgloss.NewStyle().
			Border(filePanelBottomBoard, true, true, true, true).
			BorderForeground(lipgloss.Color(theme.Border)).
			BorderBackground(backgroundWindow).
			Width(width).
			Height(height).Bold(true).Background(backgroundWindow)
	}
}

func MetaDataBoarder(height int, width int, borderBottom string, focusType focusPanelType) lipgloss.Style {
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━┫Metadata┣" + repeatString("━", width),
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
			BorderForeground(lipgloss.Color(theme.BottomBarFocus)).
			BorderBackground(backgroundWindow).
			Width(width).
			Height(height).Bold(true).Background(backgroundWindow)
	} else {
		return lipgloss.NewStyle().
			Border(filePanelBottomBoard, true, true, true, true).
			BorderForeground(lipgloss.Color(theme.Border)).
			BorderBackground(backgroundWindow).
			Width(width).
			Height(height).Bold(true).Background(backgroundWindow)
	}
}

func ClipboardBoarder(height int, width int, borderBottom string) lipgloss.Style {
	filePanelBottomBoard := lipgloss.Border{
		Top:         "━┫Clipboard┣" + repeatString("━", width),
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
		BorderBackground(backgroundWindow).
		Width(width).
		Height(height).Bold(true).Background(backgroundWindow)

}

func FilePanelDividerStyle(focusType filePanelFocusType) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(FilePanelFocusColor(focusType))).Bold(true).Background(backgroundWindow)
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
		return StringColorRender(style.color).Background(backgroundWindow).Render(style.icon) + lipgloss.NewStyle().Background(backgroundWindow).Render("  ") + filePanelItemSelected.Render(TruncateText(name, width))
	} else {
		return StringColorRender(style.color).Background(backgroundWindow).Render(style.icon) + lipgloss.NewStyle().Background(backgroundWindow).Render("  ") + filePanelItem.Render(TruncateText(name, width))
	}
}

func ClipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return StringColorRender(style.color).Background(backgroundWindow).Render(style.icon) + lipgloss.NewStyle().Background(backgroundWindow).Render("  ") + filePanelItemSelected.Render(TruncateTextBeginning(name, width))
	} else {
		return StringColorRender(style.color).Background(backgroundWindow).Render(style.icon) + lipgloss.NewStyle().Background(backgroundWindow).Render("  ") + filePanelItem.Render(TruncateTextBeginning(name, width))
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
	return repeatString("━", width-len(countString)) + "┫" + countString + "┣"
}

func StringColorRender(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func BottomWidth(fullWidth int) int {
	return fullWidth/3 - 2
}
