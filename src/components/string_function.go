package components

import (
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

func truncateText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	return text[:maxChars-3] + "..."
}

func truncateTextBeginning(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}
	runes := []rune(text)
	charsToKeep := maxChars - 3
	truncatedRunes := append([]rune("..."), runes[len(runes)-charsToKeep:]...)
	return string(truncatedRunes)
}

func truncateMiddleText(text string, maxChars int) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	halfEllipsisLength := (maxChars - 3) / 2

	truncatedText := text[:halfEllipsisLength] + "..." + text[utf8.RuneCountInString(text)-halfEllipsisLength:]

	return truncatedText
}

func prettierName(name string, width int, isDir bool, isSelected bool, bgColor lipgloss.Color) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.color), bgColor).
		Background(bgColor).
		Render(style.icon + " ") + 
		filePanelItemSelectedStyle.
		Render(truncateText(name, width))
	} else {
		return stringColorRender(lipgloss.Color(style.color), bgColor).
		Background(bgColor).
		Render(style.icon + " ") + 
		filePanelStyle.Render(truncateText(name, width))
	}
}

func clipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.color), footerBGColor).
		Background(footerBGColor).
		Render(style.icon + " ") + 
		filePanelItemSelectedStyle.Render(truncateTextBeginning(name, width))
	} else {
		return stringColorRender(lipgloss.Color(style.color), footerBGColor).
		Background(footerBGColor).
		Render(style.icon + " ") + 
		filePanelStyle.Render(truncateTextBeginning(name, width))
	}
}

// func placeOverlay(x, y int,background, placeModal string) string {
// 	lines := strings.Split(placeModal, "\n")
// 	lines = lines
// 	re := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	
// 	// 示例字符串
// 	str := "[38;2;134;134;134;48;2;30;30;46m┏A我[0m"
	
// 	// 使用 FindAllStringIndex 找出所有匹配的位置
// 	indexes := re.FindAllStringIndex(str, -1)
// 	outPutLog(str)
// 	// 檢查是否找到匹配
// 	if indexes != nil {
// 		for _, loc := range indexes {
// 			loc = mapCoords(str, loc)
// 			outPutLog(fmt.Sprintf("匹配的開始位置: %d, 結束位置: %d", loc[0], loc[1]))
// 		}
// 	} else {
// 		outPutLog("沒有找到匹配")
// 	}

// 	return ""
// }

// func mapCoords(s string, byteCoords []int) (graphemeCoords []int) {
//     graphemeCoords = make([]int, 2)
//     gr := uniseg.NewGraphemes(s)
//     graphemeIndex := -1
//     for gr.Next() {
//         graphemeIndex++
//         a, b := gr.Positions()
//         if a == byteCoords[0] {
//             graphemeCoords[0] = graphemeIndex
//         }
//         if b == byteCoords[1] {
//             graphemeCoords[1] = graphemeIndex + 1
//             break
//         }
//     }
//     return
// }