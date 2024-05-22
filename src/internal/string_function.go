package internal

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	charmansi "github.com/charmbracelet/x/exp/term/ansi"
)

func truncateText(text string, maxChars int) string {

	truncatedText := charmansi.Truncate(text, maxChars - 3, "")
	if text != truncatedText {
		return truncatedText + "..."
	}
	
	return text
}

func truncateTextBeginning(text string, maxChars int) string {
	if charmansi.StringWidth(text) <= maxChars {
		return text
	}

	runes := []rune(text)
	var truncatedRunes []rune

	truncatedRunes = runes

	truncatedWidth := charmansi.StringWidth(string(truncatedRunes))

	for truncatedWidth > maxChars {
		truncatedRunes = truncatedRunes[1:]
		truncatedWidth = charmansi.StringWidth(string(truncatedRunes))
	}

	if len(truncatedRunes) > 3 {
		truncatedRunes = append([]rune("..."), truncatedRunes[3:]...)
	}

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
			Render(style.icon+" ") +
			filePanelItemSelectedStyle.
				Render(truncateText(name, width))
	} else {
		return stringColorRender(lipgloss.Color(style.color), bgColor).
			Background(bgColor).
			Render(style.icon+" ") +
			filePanelStyle.Render(truncateText(name, width))
	}
}

func prettierDirectoryPreviewName(name string, isDir bool, bgColor lipgloss.Color) string {
	style := getElementIcon(name, isDir)
	return stringColorRender(lipgloss.Color(style.color), bgColor).
		Background(bgColor).
		Render(style.icon+" ") +
		filePanelStyle.Render(name)
}

func clipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.color), footerBGColor).
			Background(footerBGColor).
			Render(style.icon+" ") +
			filePanelItemSelectedStyle.Render(truncateTextBeginning(name, width))
	} else {
		return stringColorRender(lipgloss.Color(style.color), footerBGColor).
			Background(footerBGColor).
			Render(style.icon+" ") +
			filePanelStyle.Render(truncateTextBeginning(name, width))
	}
}

func fileNameWithoutExtension(fileName string) string {
	for {
		pos := strings.LastIndexByte(fileName, '.')
		if pos <= 0 {
			break
		}
		fileName = fileName[:pos]
	}
	return fileName
}

func formatFileSize(size int64) string {
	if size == 0 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(1024)))
	adjustedSize := float64(size) / math.Pow(1024, float64(unitIndex))

	return fmt.Sprintf("%.2f %s", adjustedSize, units[unitIndex])
}

// Truncate line lengths and keep ANSI
func checkAndTruncateLineLengths(text string, maxLength int) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		truncatedLine := charmansi.Truncate(line, maxLength, "")
		result.WriteString(truncatedLine + "\n")
	}

	finalResult := strings.TrimRight(result.String(), "\n")

	return finalResult
}

// Check file is text file or not
func isTextFile(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	_, err = reader.Read(buffer)
	if err != nil {
		return false, err
	}

	for _, b := range buffer {
		if b == 0 {
			return false, nil
		}
		if !unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			return false, nil
		}
	}

	return true, nil
}
