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

// isKittyCapable checks if the terminal supports Kitty graphics protocol
func isKittyCapable() bool {
	isCapable := rasterm.IsKittyCapable()

	// Additional detection for terminals that might not be detected by rasterm
	if !isCapable {
		termProgram := os.Getenv("TERM_PROGRAM")
		term := os.Getenv("TERM")

		// List of known terminal identifiers that support Kitty protocol
		knownTerminals := []string{
			"ghostty",
			"WezTerm",
			"iTerm2",
			"xterm-kitty",
			"kitty",
			"Konsole",
			"WarpTerminal",
		}

		for _, knownTerm := range knownTerminals {
			if strings.EqualFold(termProgram, knownTerm) || strings.EqualFold(term, knownTerm) {
				isCapable = true
				break
			}
		}
	}

	return isCapable
}

// ClearKittyImages clears all Kitty protocol images from the terminal
func ClearKittyImages() string {
	if !isKittyCapable() {
		return "" // No need to clear if terminal doesn't support Kitty protocol
	}

	return generateKittyClearCommands()
}

// ClearKittyImages clears all Kitty protocol images from the terminal
func (p *ImagePreviewer) ClearKittyImages() string {
	if !p.IsKittyCapable() {
		return "" // No need to clear if terminal doesn't support Kitty protocol
	}

	return generateKittyClearCommands()
}

// generateKittyClearCommands generates the clearing commands for Kitty protocol
func generateKittyClearCommands() string {
	var buf bytes.Buffer

	// Clear all images first
	clearAllCmd := "\x1b_Ga=d\x1b\\"
	buf.WriteString(clearAllCmd)

	// Clear all placements
	clearPlacementsCmd := "\x1b_Ga=d,p=1\x1b\\"
	buf.WriteString(clearPlacementsCmd)

	// Add a small delay command to ensure clearing is processed
	buf.WriteString("\x1b[0m")

	return buf.String()
}

// generatePlacementID generates a unique placement ID based on file path
func generatePlacementID(path string) uint32 {
	if len(path) == 0 {
		return 42 // Default fallback
	}

	hash := 0
	for _, c := range path {
		hash = hash*31 + int(c)
	}
	return uint32(hash&0xFFFF) + 1000 // Ensure it's not 0 and avoid low numbers
}

// renderWithKittyUsingTermCap renders an image using Kitty graphics protocol with terminal capabilities
func (p *ImagePreviewer) renderWithKittyUsingTermCap(img image.Image, path string, originalWidth, originalHeight, maxWidth, maxHeight int, sideAreaWidth int) (string, error) {
	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	var buf bytes.Buffer

	// Add clearing commands
	buf.WriteString(generateKittyClearCommands())

	opts := rasterm.KittyImgOpts{
		PlacementId: generatePlacementID(path),
	}

	// Get terminal cell size from ImagePreviewer's terminal capabilities
	cellSize := p.terminalCap.GetTerminalCellSize()
	pixelsPerColumn := cellSize.PixelsPerColumn
	pixelsPerRow := cellSize.PixelsPerRow

	slog.Debug("pixelsPerColumn", "pixelsPerColumn", pixelsPerColumn, "pixelsPerRow", pixelsPerRow)

	imgRatio := float64(originalWidth) / float64(originalHeight)
	termRatio := float64(maxWidth*pixelsPerColumn) / float64(maxHeight*pixelsPerRow)

	slog.Debug("imgRatio", "imgRatio", imgRatio, "termRatio", termRatio)

	if imgRatio > termRatio {
		dstCols := maxWidth
		dstRows := int(float64(dstCols*pixelsPerColumn) / imgRatio / float64(pixelsPerRow))
		opts.DstCols = uint32(dstCols)
		opts.DstRows = uint32(dstRows)
	} else {
		dstRows := maxHeight
		dstCols := int(float64(dstRows*pixelsPerRow) * imgRatio / float64(pixelsPerColumn))
		opts.DstRows = uint32(dstRows)
		opts.DstCols = uint32(dstCols)
	}

	// Write image using Kitty protocol
	if err := rasterm.KittyWriteImage(&buf, img, opts); err != nil {
		return "", err
	}

	buf.WriteString("\x1b[1;" + strconv.Itoa(sideAreaWidth) + "H")

	return buf.String(), nil
}

// IsKittyCapable checks if the terminal supports Kitty graphics protocol
func (p *ImagePreviewer) IsKittyCapable() bool {
	return isKittyCapable()
}
