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

	return isCapable
}




// renderWithInlineUsingTermCap renders an image using inline image protocol
func (p *ImagePreviewer) renderWithInlineUsingTermCap(img image.Image, path string,
	originalWidth, originalHeight, maxWidth, maxHeight int, sideAreaWidth int) (string, error) {
	
	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	var buf bytes.Buffer

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

	// Use rasterm to write the image using iTerm2/WezTerm protocol
	if err := rasterm.ItermWriteImage(&buf, img); err != nil {
		return "", fmt.Errorf("failed to write image using rasterm: %w", err)
	}

	// Position cursor after the image output (following kitty.go pattern)
	buf.WriteString("\x1b[2;" + strconv.Itoa(sideAreaWidth) + "H")

	// Add a newline to ensure proper display
	buf.WriteString("\n")

	return buf.String(), nil
}

// IsInlineCapable checks if the terminal supports inline image protocol
func (p *ImagePreviewer) IsInlineCapable() bool {
	return isInlineCapable()
}