package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
)

func TruncateText(text string, maxChars int, tails string) string {
	truncatedText := ansi.Truncate(text, maxChars-len(tails), "")
	if text != truncatedText {
		return truncatedText + tails
	}

	return text
}

func TruncateTextBeginning(text string, maxChars int, tails string) string {
	if ansi.StringWidth(text) <= maxChars {
		return text
	}

	truncatedRunes := []rune(text)

	truncatedWidth := ansi.StringWidth(string(truncatedRunes))

	for truncatedWidth > maxChars {
		truncatedRunes = truncatedRunes[1:]
		truncatedWidth = ansi.StringWidth(string(truncatedRunes))
	}

	if len(truncatedRunes) > len(tails) {
		truncatedRunes = append([]rune(tails), truncatedRunes[len(tails):]...)
	}

	return string(truncatedRunes)
}

func TruncateMiddleText(text string, maxChars int, tails string) string {
	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	halfEllipsisLength := (maxChars - 3) / 2
	// Todo : Use ansi.Substring to correctly handle ANSI escape codes
	truncatedText := text[:halfEllipsisLength] + tails + text[utf8.RuneCountInString(text)-halfEllipsisLength:]

	return truncatedText
}

func PrettierName(name string, width int, isDir bool, isSelected bool, bgColor lipgloss.Color) string {
	style := GetElementIcon(name, isDir, Config.Nerdfont)
	if isSelected {
		return StringColorRender(lipgloss.Color(style.Color), bgColor).
			Background(bgColor).
			Render(style.Icon+" ") +
			FilePanelItemSelectedStyle.
				Render(TruncateText(name, width, "..."))
	}
	return StringColorRender(lipgloss.Color(style.Color), bgColor).
		Background(bgColor).
		Render(style.Icon+" ") +
		FilePanelStyle.Render(TruncateText(name, width, "..."))
}

func PrettierDirectoryPreviewName(name string, isDir bool, bgColor lipgloss.Color) string {
	style := GetElementIcon(name, isDir, Config.Nerdfont)
	return StringColorRender(lipgloss.Color(style.Color), bgColor).
		Background(bgColor).
		Render(style.Icon+" ") +
		FilePanelStyle.Render(name)
}

func ClipboardPrettierName(name string, width int, isDir bool, isSelected bool) string {
	style := GetElementIcon(name, isDir, Config.Nerdfont)
	if isSelected {
		return StringColorRender(lipgloss.Color(style.Color), FooterBGColor).
			Background(FooterBGColor).
			Render(style.Icon+" ") +
			FilePanelItemSelectedStyle.Render(TruncateTextBeginning(name, width, "..."))
	}
	return StringColorRender(lipgloss.Color(style.Color), FooterBGColor).
		Background(FooterBGColor).
		Render(style.Icon+" ") +
		FilePanelStyle.Render(TruncateTextBeginning(name, width, "..."))
}

func FileNameWithoutExtension(fileName string) string {
	for {
		pos := strings.LastIndexByte(fileName, '.')
		if pos <= 0 {
			break
		}
		fileName = fileName[:pos]
	}
	return fileName
}

func FormatFileSize(size int64) string {
	if size == 0 {
		return "0B"
	}

	unitsDec := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	unitsBin := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	// Todo : Remove duplication here
	if Config.FileSizeUseSI {
		unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(1000)))
		adjustedSize := float64(size) / math.Pow(1000, float64(unitIndex))
		return fmt.Sprintf("%.2f %s", adjustedSize, unitsDec[unitIndex])
	}
	unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(1024)))
	adjustedSize := float64(size) / math.Pow(1024, float64(unitIndex))
	return fmt.Sprintf("%.2f %s", adjustedSize, unitsBin[unitIndex])
}

// Truncate line lengths and keep ANSI
func CheckAndTruncateLineLengths(text string, maxLength int) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	for _, line := range lines {
		// Replace tabs with spaces
		expandedLine := strings.ReplaceAll(line, "\t", strings.Repeat(" ", 4))
		truncatedLine := ansi.Truncate(expandedLine, maxLength, "")
		result.WriteString(truncatedLine + "\n")
	}

	finalResult := strings.TrimRight(result.String(), "\n")

	return finalResult
}

// Separated this out out for easy testing
func IsBufferPrintable(buffer []byte) bool {
	for _, b := range buffer {
		// This will also handle b==0
		if !unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			return false
		}
	}
	return true
}

// IsExtensionExtractable checks if a string is a valid compressed archive file extension.
func IsExtensionExtractable(ext string) bool {
	// Extensions based on the types that package: `xtractr` `ExtractFile` function handles.
	validExtensions := map[string]struct{}{
		".zip":     {},
		".bz":      {},
		".gz":      {},
		".iso":     {},
		".rar":     {},
		".7z":      {},
		".tar":     {},
		".tar.gz":  {},
		".tar.bz2": {},
	}
	_, exists := validExtensions[strings.ToLower(ext)]
	return exists
}

// Check file is text file or not
func IsTextFile(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	cnt, err := reader.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	return IsBufferPrintable(buffer[:cnt]), nil
}

// Although some characters like `\x0b`(vertical tab) are printable,
// previewing them breaks the layout.
// So, among the "non-graphic" printable characters, we only need \n and \t
// Space and NBSP are already considered graphic by unicode.
func MakePrintable(line string) string {
	var sb strings.Builder
	// This has to be looped byte-wise, looping it rune-wise
	// or by using strings.Map would cause issues with strings like
	// "(NBSP)\xa0"
	for i := range len(line) {
		r := rune(line[i])
		if unicode.IsGraphic(r) || r == rune('\t') || r == rune('\n') {
			sb.WriteByte(line[i])
		}
	}
	return sb.String()
}
