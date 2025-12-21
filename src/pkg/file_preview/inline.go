package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/charmbracelet/lipgloss"
)

// isInlineCapable checks if the terminal supports inline image protocol (iTerm2, WezTerm, etc.)
func isInlineCapable() bool {
	isCapable := rasterm.IsItermCapable()

	// Additional detection for terminals that might not be detected by rasterm
	if !isCapable {
		termProgram := os.Getenv("TERM_PROGRAM")
		term := os.Getenv("TERM")

		// List of known terminal identifiers that support inline image protocol
		knownTerminals := []string{
			"iTerm2",
			"iTerm.app",
			"WezTerm",
			"Hyper",
			"Terminus",
			"Tabby",
		}

		for _, knownTerm := range knownTerminals {
			if strings.EqualFold(termProgram, knownTerm) || strings.EqualFold(term, knownTerm) {
				isCapable = true
				break
			}
		}

		// Additional check for iTerm2 specific environment variables
		if !isCapable && (os.Getenv("ITERM_SESSION_ID") != "" || os.Getenv("ITERM_PROFILE") != "") {
			isCapable = true
		}
	}
	slog.Debug("[TEMP] ", "inlineCapabale", isCapable)

	return isCapable
}

// ClearInlineImage clears all inline image protocol images from the terminal
func (p *ImagePreviewer) ClearInlineImage() string {
	return ""
}

// renderWithInlineUsingTermCap renders an image using inline image protocol
func (p *ImagePreviewer) renderWithInlineUsingTermCap(
	img image.Image,
	path string,
	originalWidth, originalHeight, maxWidth, maxHeight int,
	sideAreaWidth int,
	filePreviewStyle lipgloss.Style,
) (string, error) {
	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	slog.Debug("inline renderer starting", "path", path, "maxWidth", maxWidth, "maxHeight", maxHeight)

	// Calculate display dimensions in character cells
	imgRatio := float64(originalWidth) / float64(originalHeight)
	termRatio := float64(maxWidth) / float64(maxHeight)

	var displayWidthCells, displayHeightCells int

	if imgRatio > termRatio {
		// Image is wider, constrain by width
		displayWidthCells = maxWidth
		displayHeightCells = int(float64(maxWidth) / imgRatio)
	} else {
		// Image is taller, constrain by height
		displayHeightCells = maxHeight
		displayWidthCells = int(float64(maxHeight) * imgRatio)
	}

	// Ensure minimum dimensions
	if displayWidthCells < 1 {
		displayWidthCells = 1
	}
	if displayHeightCells < 1 {
		displayHeightCells = 1
	}

	slog.Debug("inline display dimensions", "widthCells", displayWidthCells, "heightCells", displayHeightCells)

	var buf bytes.Buffer
	line := strings.Repeat("-", maxWidth)
	block := strings.Repeat(line+"\n", maxHeight-1) + line
	//block := strings.Repeat("------------------\n", 10)
	buf.WriteString(filePreviewStyle.Render(block))
	buf.WriteString("\x1b[s")
	buf.WriteString("\x1b[1;" + strconv.Itoa(sideAreaWidth) + "H")
	opts := rasterm.ItermImgOpts{
		Width:             strconv.FormatInt(int64(displayWidthCells), 10),
		Height:            strconv.FormatInt(int64(displayHeightCells), 10),
		IgnoreAspectRatio: true,
		DisplayInline:     true,
	}
	// Use rasterm to write the image using iTerm2/WezTerm protocol
	if err := rasterm.ItermWriteImageWithOptions(&buf, img, opts); err != nil {
		return "", fmt.Errorf("failed to write image using rasterm: %w", err)
	}
	buf.WriteString("\x1b[u")
	slog.Debug("[TEMP]", "res", buf.String())
	return buf.String(), nil
}

// IsInlineCapable checks if the terminal supports inline image protocol
func (p *ImagePreviewer) IsInlineCapable() bool {
	return isInlineCapable()
}
