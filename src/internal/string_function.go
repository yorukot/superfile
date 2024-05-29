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
	"github.com/charmbracelet/x/exp/term/ansi"
)

func truncateText(text string, maxChars int, talis string) string {

	truncatedText := ansi.Truncate(text, maxChars-len(talis), "")
	if text != truncatedText {
		return truncatedText + talis
	}

	return text
}

func truncateTextBeginning(text string, maxChars int, talis string) string {
	if ansi.StringWidth(text) <= maxChars {
		return text
	}

	truncatedRunes := []rune(text)

	truncatedWidth := ansi.StringWidth(string(truncatedRunes))

	for truncatedWidth > maxChars {
		truncatedRunes = truncatedRunes[1:]
		truncatedWidth = ansi.StringWidth(string(truncatedRunes))
	}

	if len(truncatedRunes) > len(talis) {
		truncatedRunes = append([]rune(talis), truncatedRunes[len(talis):]...)
	}

	return string(truncatedRunes)
}

func truncateMiddleText(text string, maxChars int, talis string) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	halfEllipsisLength := (maxChars - 3) / 2

	truncatedText := text[:halfEllipsisLength] + talis + text[utf8.RuneCountInString(text)-halfEllipsisLength:]

	return truncatedText
}

func prettierName(name string, width int, isDir bool, isSelected bool, bgColor lipgloss.Color) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.Color), bgColor).
			Background(bgColor).
			Render(style.Icon+" ") +
			filePanelItemSelectedStyle.
				Render(truncateText(name, width, "..."))
	} else {
		return stringColorRender(lipgloss.Color(style.Color), bgColor).
			Background(bgColor).
			Render(style.Icon+" ") +
			filePanelStyle.Render(truncateText(name, width, "..."))
	}
}

func prettierDirectoryPreviewName(name string, isDir bool, bgColor lipgloss.Color) string {
	style := getElementIcon(name, isDir)
	return stringColorRender(lipgloss.Color(style.Color), bgColor).
		Background(bgColor).
		Render(style.Icon+" ") +
		filePanelStyle.Render(name)
}

func clipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := getElementIcon(name, isDir)
	if isSelected {
		return stringColorRender(lipgloss.Color(style.Color), footerBGColor).
			Background(footerBGColor).
			Render(style.Icon+" ") +
			filePanelItemSelectedStyle.Render(truncateTextBeginning(name, width, "..."))
	} else {
		return stringColorRender(lipgloss.Color(style.Color), footerBGColor).
			Background(footerBGColor).
			Render(style.Icon+" ") +
			filePanelStyle.Render(truncateTextBeginning(name, width, "..."))
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
		truncatedLine := ansi.Truncate(line, maxLength, "")
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
