package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/sixel"
	"github.com/disintegration/imaging"
)

// Sixel images are drawn in pixels at the cursor position, outside of the
// cell buffer that Bubble Tea renders. To make them work with it:
//
//   - The view reserves the image's area with marker cells: an invisible
//     character that differs from a plain space. Whenever anything else is
//     shown there (text, another image, an overlay), those cells change, so the
//     renderer rewrites them, and the terminal erases the image pixels under
//     rewritten cells. No explicit erase sequence is needed.
//   - The image itself is sent separately via tea.Raw, a short delay after the
//     frame containing the marker cells, because Bubble Tea flushes raw output
//     before each frame. See SixelPlacement.

const (
	// SixelMarker is the character for the cells under a Sixel image.
	SixelMarker = "⠀" // BRAILLE PATTERN BLANK
	// SixelMarkerAlt is used for every other image, so that the cells of one
	// image never look identical to the cells of the next one, and every cell
	// is rewritten (erasing the previous image) before the next is drawn.
	SixelMarkerAlt = " " // NO-BREAK SPACE

	// Windows Terminal maps Sixel images onto a virtual 10x20 pixel cell,
	// whatever the font size.
	windowsTerminalPixelsPerColumn = 10
	windowsTerminalPixelsPerRow    = 20

	// Sixel DCS parameters: P1 is the pixel aspect ratio (0 = default 2:1, which
	// raster attributes override), P2 = 1 leaves pixels not set by the image
	// unchanged, so transparent areas show the cell background.
	sixelAspectDefault  = 0
	sixelTransparentBG  = 1
	sixelGridSizeUnused = 0
)

// SixelImage is an encoded Sixel image laid out inside a preview area.
type SixelImage struct {
	// DCS is the complete Sixel device control string.
	DCS string
	// OffsetX and OffsetY are the image's top-left cell inside the content area.
	OffsetX int
	OffsetY int
	// Cols and Rows are the image's size in cells.
	Cols int
	Rows int
}

// Placement returns the raw sequence that draws the image with its top-left
// cell at the given 0-based screen position of the content area's top-left
// cell. The cursor is saved and restored so the renderer's idea of the cursor
// position stays correct.
func (s *SixelImage) Placement(originX, originY int) string {
	return ansi.SaveCursor +
		ansi.CursorPosition(originX+s.OffsetX+1, originY+s.OffsetY+1) +
		s.DCS +
		ansi.RestoreCursor
}

// isSixelCapable reports whether the terminal is known to support Sixel.
// Only Windows Terminal is detected for now; other terminals can use the
// image_renderer config option.
func isSixelCapable() bool {
	return isWindowsTerminal()
}

func isWindowsTerminal() bool {
	return os.Getenv("WT_SESSION") != ""
}

// sixelCellSize returns the pixel size of a cell for Sixel images.
func (p *ImagePreviewer) sixelCellSize() TerminalCellSize {
	if isWindowsTerminal() {
		return TerminalCellSize{
			PixelsPerColumn: windowsTerminalPixelsPerColumn,
			PixelsPerRow:    windowsTerminalPixelsPerRow,
		}
	}
	return p.terminalCap.GetTerminalCellSize()
}

// renderWithSixel encodes img as Sixel, sized to fit in maxWidth x maxHeight
// cells. It returns the cell block to show in the view (spaces, with marker
// cells where the image goes) and the encoded image.
func (p *ImagePreviewer) renderWithSixel(img image.Image, maxWidth, maxHeight int) (string, *SixelImage, error) {
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", nil, fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	cellSize := p.sixelCellSize()
	bounds := img.Bounds()
	cols, rows := fitImageInCells(bounds.Dx(), bounds.Dy(), maxWidth, maxHeight, cellSize)

	scaled := imaging.Fit(img, cols*cellSize.PixelsPerColumn, rows*cellSize.PixelsPerRow, imaging.Lanczos)

	var payload bytes.Buffer
	if err := (&sixel.Encoder{}).Encode(&payload, scaled); err != nil {
		return "", nil, fmt.Errorf("failed to encode sixel: %w", err)
	}

	result := &SixelImage{
		DCS:     ansi.SixelGraphics(sixelAspectDefault, sixelTransparentBG, sixelGridSizeUnused, payload.Bytes()),
		OffsetX: (maxWidth - cols) / centerDivisor,
		OffsetY: (maxHeight - rows) / centerDivisor,
		Cols:    cols,
		Rows:    rows,
	}
	return buildSixelBlock(result, maxWidth, maxHeight), result, nil
}

// buildSixelBlock builds a maxWidth x maxHeight block of spaces with marker
// cells where the image goes.
func buildSixelBlock(s *SixelImage, maxWidth, maxHeight int) string {
	blankLine := strings.Repeat(" ", maxWidth)
	imageLine := strings.Repeat(" ", s.OffsetX) + strings.Repeat(SixelMarker, s.Cols) +
		strings.Repeat(" ", maxWidth-s.OffsetX-s.Cols)

	lines := make([]string, maxHeight)
	for i := range lines {
		if i >= s.OffsetY && i < s.OffsetY+s.Rows {
			lines[i] = imageLine
		} else {
			lines[i] = blankLine
		}
	}
	return strings.Join(lines, "\n")
}

// fitImageInCells returns the largest cell size, within maxWidth x maxHeight,
// that keeps the image's aspect ratio.
func fitImageInCells(imgWidth, imgHeight, maxWidth, maxHeight int, cellSize TerminalCellSize) (int, int) {
	imgRatio := float64(imgWidth) / float64(imgHeight)
	areaRatio := float64(maxWidth*cellSize.PixelsPerColumn) / float64(maxHeight*cellSize.PixelsPerRow)

	var cols, rows int
	if imgRatio > areaRatio {
		cols = maxWidth
		rows = int(float64(cols*cellSize.PixelsPerColumn) / imgRatio / float64(cellSize.PixelsPerRow))
	} else {
		rows = maxHeight
		cols = int(float64(rows*cellSize.PixelsPerRow) * imgRatio / float64(cellSize.PixelsPerColumn))
	}
	return max(1, min(cols, maxWidth)), max(1, min(rows, maxHeight))
}
