package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/kitty"
)

// isKittyCapable checks if the terminal supports Kitty graphics protocol
func isKittyCapable() bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	term := os.Getenv("TERM")

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
			return true
		}
	}

	return false
}

// ClearKittyImages returns a Kitty graphics delete-all command.
func ClearKittyImages() string {
	if !isKittyCapable() {
		return ""
	}
	return ansi.KittyGraphics(nil, "a=d")
}

// ClearKittyImages returns a Kitty graphics delete-all command.
func (p *ImagePreviewer) ClearKittyImages() string {
	if !p.IsKittyCapable() {
		return ""
	}
	return ansi.KittyGraphics(nil, "a=d")
}

// generatePlacementID generates a unique placement ID based on file path
func generatePlacementID(path string) int {
	if len(path) == 0 {
		return kittyHashSeed
	}

	hash := 0
	for _, c := range path {
		hash = hash*kittyHashPrime + int(c)
	}
	return (hash & kittyMaxID) + kittyNonZeroOffset
}

// KittyImageResult holds both the placeholder string for the cell buffer
// and the raw transmission data to send directly to the terminal.
type KittyImageResult struct {
	// Placeholders is the Unicode placeholder string for embedding in the view.
	// It contains kitty.Placeholder characters with diacritics.
	Placeholders string
	// RawTransmit is the Kitty graphics APC data to send via tea.Raw().
	// It transmits the image data to the terminal out-of-band.
	RawTransmit string
}

// renderWithKittyUsingTermCap renders an image using Kitty graphics protocol
// with Unicode virtual placeholders (compatible with cell-based renderers).
func (p *ImagePreviewer) renderWithKittyUsingTermCap(img image.Image, path string,
	originalWidth, originalHeight, maxWidth, maxHeight int, _ int,
) (*KittyImageResult, error) {
	if maxWidth <= 0 || maxHeight <= 0 {
		return nil, fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	cellSize := p.terminalCap.GetTerminalCellSize()
	pixelsPerColumn := cellSize.PixelsPerColumn
	pixelsPerRow := cellSize.PixelsPerRow

	slog.Debug("pixelsPerColumn", "pixelsPerColumn", pixelsPerColumn, "pixelsPerRow", pixelsPerRow)

	imgRatio := float64(originalWidth) / float64(originalHeight)
	termRatio := float64(maxWidth*pixelsPerColumn) / float64(maxHeight*pixelsPerRow)

	var dstCols, dstRows int
	if imgRatio > termRatio {
		dstCols = maxWidth
		dstRows = int(float64(dstCols*pixelsPerColumn) / imgRatio / float64(pixelsPerRow))
	} else {
		dstRows = maxHeight
		dstCols = int(float64(dstRows*pixelsPerRow) * imgRatio / float64(pixelsPerColumn))
	}
	if dstCols <= 0 {
		dstCols = 1
	}
	if dstRows <= 0 {
		dstRows = 1
	}

	imgID := generatePlacementID(path)
	imgArea := img.Bounds()

	// Encode image data for transmission via tea.Raw()
	var transmitBuf bytes.Buffer

	// Delete previous image with this ID first
	transmitBuf.WriteString(ansi.KittyGraphics(nil, fmt.Sprintf("a=d,d=i,i=%d", imgID)))

	if err := kitty.EncodeGraphics(&transmitBuf, img, &kitty.Options{
		ID:               imgID,
		Action:           kitty.TransmitAndPut,
		Transmission:     kitty.Direct,
		Format:           kitty.RGBA,
		ImageWidth:       imgArea.Dx(),
		ImageHeight:      imgArea.Dy(),
		Columns:          dstCols,
		Rows:             dstRows,
		VirtualPlacement: true,
		Quite:            kittyQuietAll,
		Chunk:            true,
	}); err != nil {
		return nil, fmt.Errorf("failed to encode kitty graphics: %w", err)
	}

	// Build Unicode placeholder cells for the view
	placeholders := buildKittyPlaceholders(imgID, dstCols, dstRows)

	return &KittyImageResult{
		Placeholders: placeholders,
		RawTransmit:  transmitBuf.String(),
	}, nil
}

// buildKittyPlaceholders builds a string of Kitty Unicode placeholder characters
// that the terminal replaces with the transmitted image.
func buildKittyPlaceholders(imgID int, cols, rows int) string {
	// Encode image ID as foreground color for the placeholder cells.
	// The terminal uses this color to identify which image to display.
	r, g, b := byte((imgID>>rgbShift16)&rgbMask), byte((imgID>>rgbShift8)&rgbMask), byte(imgID&rgbMask)

	var fgSeq string
	if r == 0 && g == 0 {
		// Use 256-color mode for small IDs
		fgSeq = fmt.Sprintf("\x1b[38;5;%dm", b)
	} else {
		fgSeq = fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}
	resetSeq := "\x1b[39m"

	var buf strings.Builder
	for y := range rows {
		if y > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(fgSeq)
		// First cell per row gets placeholder + row diacritic + col(0) diacritic
		buf.WriteRune(kitty.Placeholder)
		buf.WriteRune(kitty.Diacritic(y))
		buf.WriteRune(kitty.Diacritic(0))
		// Subsequent cells just get the placeholder
		for x := 1; x < cols; x++ {
			buf.WriteRune(kitty.Placeholder)
		}
		buf.WriteString(resetSeq)
	}
	return buf.String()
}

// IsKittyCapable checks if the terminal supports Kitty graphics protocol
func (p *ImagePreviewer) IsKittyCapable() bool {
	return isKittyCapable()
}
