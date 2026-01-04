package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"
	"os"
	"strings"

	"github.com/blacktop/go-termimg"

	"github.com/yorukot/superfile/src/internal/common"
)

// isKittyCapable checks if the terminal supports Kitty graphics protocol
func isKittyCapable() bool {
	isCapable := termimg.KittySupported()

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

	// Reset text formatting to default
	buf.WriteString("\x1b[0m")

	return buf.String()
}

// renderWithKittyUsingTermCap renders an image using Kitty graphics protocol with terminal capabilities
func (p *ImagePreviewer) renderWithKittyUsingTermCap(img image.Image,
	originalWidth, originalHeight, maxWidth, maxHeight int, sideAreaWidth int,
) (string, error) {
	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	var buf bytes.Buffer

	// Add clearing commands
	if err := termimg.ClearAll(); err != nil {
		slog.Error("Failed to clear previous images", "error", err)
	}

	// Get terminal cell size from ImagePreviewer's terminal capabilities
	cellSize := p.terminalCap.GetTerminalCellSize()
	pixelsPerColumn := cellSize.PixelsPerColumn
	pixelsPerRow := cellSize.PixelsPerRow

	slog.Debug("pixelsPerColumn", "pixelsPerColumn", pixelsPerColumn, "pixelsPerRow", pixelsPerRow)

	imgRatio := float64(originalWidth) / float64(originalHeight)
	termRatio := float64(maxWidth*pixelsPerColumn) / float64(maxHeight*pixelsPerRow)

	slog.Debug("imgRatio", "imgRatio", imgRatio, "termRatio", termRatio)

	widget := termimg.NewImageWidget(termimg.New(img))

	if imgRatio > termRatio {
		dstCols := maxWidth
		dstRows := int(float64(dstCols*pixelsPerColumn) / imgRatio / float64(pixelsPerRow))

		widget.SetPosition(dstCols, dstRows).SetSizeWithCorrection(dstCols, dstRows)
	} else {
		dstRows := maxHeight
		dstCols := int(float64(dstRows*pixelsPerRow) * imgRatio / float64(pixelsPerColumn))

		widget.SetPosition(dstCols, dstRows).SetSizeWithCorrection(dstCols, dstRows)
	}

	// Write image using Kitty protocol
	widget.SetProtocol(termimg.Kitty)
	rendered, renderError := widget.Render()
	if renderError != nil {
		return "", renderError
	}
	buf.WriteString(rendered)

	// TODO: using internal/common package in pkg package is against the standards
	// We shouldn't use that here.
	// Other usage of common in `file_preview` should be removed too.
	// common.VideoExtensions should be moved to fixed_variables
	// and internal/common/utils shoud move to pkg/utils so that it can
	// be used by everyone

	// TODO : Ideally we should not need the kitty previewer to be
	// aware of full modal width and make decisions based on global config
	// A better solutions than this is needed for it.
	row := 1
	col := sideAreaWidth + 1
	if common.Config.EnableFilePreviewBorder {
		row++
		col++
	}
	buf.WriteString(fmt.Sprintf("\x1b[%d;%dH", row, col))

	return buf.String(), nil
}

// IsKittyCapable checks if the terminal supports Kitty graphics protocol
func (p *ImagePreviewer) IsKittyCapable() bool {
	return isKittyCapable()
}
