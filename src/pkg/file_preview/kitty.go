package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/BourgeoisBear/rasterm"
	"github.com/disintegration/imaging"
)

// Terminal cell to pixel conversion constants
// These approximate the pixel dimensions of terminal cells
const (
	DefaultPixelsPerColumn = 10 // approximate pixels per terminal column
	DefaultPixelsPerRow    = 20 // approximate pixels per terminal row
)

// calculateTerminalCells calculates how many terminal cells are needed for given pixel dimensions
func calculateTerminalCells(pixelWidth, pixelHeight, maxCols, maxRows int) (cols, rows int) {
	// Convert pixel dimensions to terminal cell dimensions
	cols = (pixelWidth + DefaultPixelsPerColumn - 1) / DefaultPixelsPerColumn // round up
	rows = (pixelHeight + DefaultPixelsPerRow - 1) / DefaultPixelsPerRow      // round up

	// Ensure we don't exceed the preview area
	if cols > maxCols {
		cols = maxCols
	}
	if rows > maxRows {
		rows = maxRows
	}

	// Ensure minimum bounds
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	return cols, rows
}

// prepareImageForKittyDisplay resizes image to fit within terminal preview area
func prepareImageForKittyDisplay(img image.Image, maxCols, maxRows int) image.Image {
	// Calculate the maximum pixel dimensions we can use for the preview
	maxPixelWidth := maxCols * DefaultPixelsPerColumn
	maxPixelHeight := maxRows * DefaultPixelsPerRow

	// Resize image to fit within these pixel bounds while maintaining aspect ratio
	return imaging.Fit(img, maxPixelWidth, maxPixelHeight, imaging.Lanczos)
}

// isKittyCapable checks if the terminal supports Kitty graphics protocol
func isKittyCapable() bool {
	isCapable := rasterm.IsKittyCapable()

	// Additional detection for terminals that might not be detected by rasterm
	if !isCapable {
		termProgram := os.Getenv("TERM_PROGRAM")
		term := os.Getenv("TERM")

		if termProgram == "ghostty" ||
			term == "xterm-ghostty" ||
			term == "ghostty" ||
			termProgram == "Ghostty" {
			isCapable = true
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

// renderWithKitty renders an image using Kitty graphics protocol
func renderWithKitty(img image.Image, path string, maxWidth, maxHeight, originalWidth, originalHeight int) (string, error) {
	var buf bytes.Buffer

	// Add clearing commands
	buf.WriteString(generateKittyClearCommands())

	// Prepare image for Kitty rendering
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Prepare image for display and calculate terminal cells needed
	previewImg := prepareImageForKittyDisplay(img, maxWidth, maxHeight)
	cols, rows := calculateTerminalCells(previewImg.Bounds().Dx(), previewImg.Bounds().Dy(), maxWidth, maxHeight)

	slog.Info("Kitty rendering",
		"originalSize", fmt.Sprintf("%dx%d", originalWidth, originalHeight),
		"processedSize", fmt.Sprintf("%dx%d", imgWidth, imgHeight),
		"previewSize", fmt.Sprintf("%dx%d", previewImg.Bounds().Dx(), previewImg.Bounds().Dy()),
		"maxPreview", fmt.Sprintf("%dx%d", maxWidth, maxHeight),
		"terminalCells", fmt.Sprintf("%dx%d", cols, rows),
	)

	// Generate placement ID and options
	placementID := generatePlacementID(path)
	opts := rasterm.KittyImgOpts{
		DstCols:     uint32(cols),
		DstRows:     uint32(rows),
		PlacementId: placementID,
	}

	// Write image using Kitty protocol
	if err := rasterm.KittyWriteImage(&buf, previewImg, opts); err != nil {
		return "", err
	}

	// Add explicit position reset at the end
	buf.WriteString("\x1b[H")

	return buf.String(), nil
}
