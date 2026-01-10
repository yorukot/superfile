package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Size calculation constants
const (
	KilobyteSize      = 1000 // SI decimal unit
	KibibyteSize      = 1024 // Binary unit
	TabWidth          = 4    // Standard tab expansion width
	DefaultBufferSize = 1024 // Default buffer size for string operations
	NonBreakingSpace  = 0xa0 // Unicode non-breaking space
	EscapeChar        = 0x1b // ANSI escape character
	ASCIIMax          = 0x7f // Maximum ASCII character value
)

// TODO: This has a bug. Remove its usage. Remove all custom truncation
// And audit and evaluate any problem
// The logic truncates to maxChars - len(tails) first, then checks if truncation occurred. This means:
// - "Hello" with maxChars=5 gets truncated to 2 chars (5-3=2), producing "He..."
// - "Hello" with maxChars=6 gets truncated to 3 chars (6-3=3), producing "Hel..."
// Both cases are wrong - "Hello" fits within 5 and 6 characters, so it shouldn't be truncated at all.
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

	//nolint:mnd // standard halving for center truncation
	halfEllipsisLength := (maxChars - 3) / 2
	// TODO : Use ansi.Substring to correctly handle ANSI escape codes
	truncatedText := text[:halfEllipsisLength] + tails + text[utf8.RuneCountInString(text)-halfEllipsisLength:]

	return truncatedText
}

func FilePanelItemRenderWithIcon(
	name string,
	width int,
	isDir bool,
	isLink bool,
	isSelected bool,
	bgColor lipgloss.Color,
) string {
	style := GetElementIcon(name, isDir, isLink, Config.Nerdfont)
	iconData := style.Icon + " "
	filenameWidth := width - ansi.StringWidth(iconData)
	if filenameWidth <= 0 {
		// This should never happen, unless there is extremely low size or programming bug
		slog.Debug("Too low width for rendering file name", "width", width, "filenameWidth", filenameWidth)
		return ""
	}
	return StringColorRender(lipgloss.Color(style.Color), bgColor).
		Background(bgColor).Render(iconData) +
		FilePanelItemRender(name, filenameWidth, isSelected, bgColor, lipgloss.Left)
}
func FilePanelItemRender(data string,
	width int,
	isSelected bool,
	bgColor lipgloss.Color,
	alignment lipgloss.Position,
) string {
	outputData := ansi.Truncate(data, width, "...")
	style := FilePanelStyle
	if isSelected {
		style = FilePanelItemSelectedStyle
	}
	return style.Background(bgColor).Width(width).Align(alignment).Render(outputData)
}

func ClipboardPrettierName(name string, width int, isDir bool, isLink bool, isSelected bool) string {
	style := GetElementIcon(filepath.Base(name), isDir, isLink, Config.Nerdfont)
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

	// TODO : Remove duplication here
	if Config.FileSizeUseSI {
		unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(KilobyteSize)))
		adjustedSize := float64(size) / math.Pow(KilobyteSize, float64(unitIndex))
		return fmt.Sprintf("%.2f %s", adjustedSize, unitsDec[unitIndex])
	}
	unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(KibibyteSize)))
	adjustedSize := float64(size) / math.Pow(KibibyteSize, float64(unitIndex))
	return fmt.Sprintf("%.2f %s", adjustedSize, unitsBin[unitIndex])
}

func GetHelpMenuHotkeyString(hotkeys []string) string {
	hotkey := ""

	for i, key := range hotkeys {
		if key == "" {
			continue
		}

		if i != 0 {
			hotkey += " | "
		}
		if key == " " {
			key = "space"
		}
		hotkey += key
	}

	return hotkey
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
	buffer := make([]byte, DefaultBufferSize)
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
// Allow Any rune that is above ASCII control characters range 0x7f
// for valid unicodes like nerdfont \uf410 \U000f0868
// Also allow \x0b that is for escape sequences
// This function should better not be broken into multiple functions
func MakePrintableWithEscCheck(line string, allowEsc bool) string { //nolint: gocognit // See above
	var sb strings.Builder
	for _, r := range line {
		if r == utf8.RuneError {
			continue
		}
		// It needs to be handled separately since considered a space,
		// It is multi-byte in UTF-8, But it has zero display width
		if r == NonBreakingSpace {
			sb.WriteRune(r)
			continue
		}
		// It needs to be handled separately since considered a space,
		// Since we are using ansi.StringWidth() for truncation, and \t is
		// considered zero width
		if r == '\t' {
			sb.WriteString("    ")
			continue
		}
		if r == EscapeChar {
			if allowEsc {
				sb.WriteRune(r)
			}
			continue
		}
		if r > ASCIIMax {
			if unicode.IsSpace(r) && utf8.RuneLen(r) > 1 {
				// See https://github.com/charmbracelet/x/issues/466
				// Space chacters spanning more than one bytes are not handled well by
				// ansi.Wrap(), and so lipgloss.Render() has issues
				r = ' '
			}
			sb.WriteRune(r)
			continue
		}
		if unicode.IsGraphic(r) || r == rune('\n') {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func MakePrintable(line string) string {
	// We assume default behaviour of allowing ESC is not  a problem
	// If you disallow ESC, then you would see ansi codes afer \x1b and it will look ugly
	// But thats only for files with that kind of data, and its rare.
	// And yazi does it too.
	// We will keep it false only if it can cause a rendering problem
	return MakePrintableWithEscCheck(line, true)
}
