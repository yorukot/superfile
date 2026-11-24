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

const (
	rendererAuto  = "auto"
	rendererKitty = "kitty"
	rendererANSI  = "ansi"
)

func normalizeImagePreviewRenderer(renderer string) string {
	switch strings.ToLower(strings.TrimSpace(renderer)) {
	case "", rendererAuto:
		return rendererAuto
	case rendererKitty:
		return rendererKitty
	case rendererANSI:
		return rendererANSI
	default:
		return rendererAuto
	}
}

// TODO: Replace this allowlist with a real Kitty graphics capability check.
// tmux masks the underlying terminal through TERM/TERM_PROGRAM.
var knownKittyTerminals = []string{ //nolint:gochecknoglobals // Static terminal allowlist, effectively const.
	"ghostty",
	// "wezterm",
	"iterm2",
	"xterm-kitty",
	"kitty",
	"konsole",
	"warpterminal",
}

// kittyEnvHints are environment variables set by kitty-capable terminals.
// They are inherited into tmux panes and are checked with or without tmux.
var kittyEnvHints = []string{ //nolint:gochecknoglobals // Static env hint allowlist, effectively const.
	"KITTY_WINDOW_ID",
	"GHOSTTY_RESOURCES_DIR",
	"GHOSTTY_BIN_DIR",
	// NOTE: Out of my scope.
	// Wezterm: freezes (didn't work, only with ascii),
	// iTerm2: sometimes is slow, iTerm2 sometimes freezes and slow (works but poor performance on MacOS)
	// konsole terminal: image position is no right also break the tui (even with tmux) (works but bad position)
	//
	// "WEZTERM_PANE",
	// "WEZTERM_EXECUTABLE",
	// "ITERM_SESSION_ID",
	// "KONSOLE_VERSION",
	//"WARP_IS_LOCAL_SHELL_SESSION",
}

// reports whether the terminal environment indicates support for the Kitty graphics protocol.
func detectKittySupport() bool {
	inTmuxSession := inTmux()
	if inTmuxSession && !ensureTmuxPassthrough() {
		return false
	}

	if isKnownTerminalAllowed() {
		return true
	}

	// Terminal-specific env vars (e.g. KITTY_WINDOW_ID) are checked with or without tmux
	if hasKittyEnvHint() {
		return true
	}

	// Some tmux setups preserve original TERM_PROGRAM in environment.
	if inTmuxSession {
		for _, envKey := range []string{"ORIG_TERM_PROGRAM", "LC_TERMINAL"} {
			val := os.Getenv(envKey)
			for _, knownTerm := range knownKittyTerminals {
				if strings.EqualFold(val, knownTerm) {
					return true
				}
			}
		}
	}

	return false
}

func isKnownTerminalAllowed() bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	term := os.Getenv("TERM")
	for _, knownTerm := range knownKittyTerminals {
		if strings.EqualFold(termProgram, knownTerm) ||
			strings.Contains(strings.ToLower(term), strings.ToLower(knownTerm)) {
			return true
		}
	}
	return false
}

// hasKittyEnvHint reports whether any terminal-specific environment variable of a kitty-capable terminal is present.
// Inside tmux these are inherited from the outer terminal, which is how the outer terminal's graphics capability can be inferred.
func hasKittyEnvHint() bool {
	for _, key := range kittyEnvHints {
		if os.Getenv(key) != "" {
			return true
		}
	}
	return false
}

// GetKittyClearRaw returns the raw APC command to clear all Kitty images.
// This must be sent via tea.Raw(), not embedded in view content.
func (p *ImagePreviewer) GetKittyClearRaw() string {
	if !p.IsKittyCapable() {
		return ""
	}
	raw, err := tmuxPassthroughWrap(ansi.KittyGraphics(nil, "a=d"))
	if err != nil {
		slog.Error("Failed to prepare kitty clear command for tmux passthrough", "error", err)
		return ""
	}
	return raw
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
		Quiet:            kittyQuietAll,
		Chunk:            true,
	}); err != nil {
		return nil, fmt.Errorf("failed to encode kitty graphics: %w", err)
	}

	// Build Unicode placeholder cells for the view
	placeholders := buildKittyPlaceholders(imgID, dstCols, dstRows)

	rawTransmit, err := tmuxPassthroughWrap(transmitBuf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to prepare kitty graphics for tmux passthrough: %w", err)
	}

	return &KittyImageResult{
		Placeholders: placeholders,
		RawTransmit:  rawTransmit,
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
	return p.isKittyCapable()
}

func (p *ImagePreviewer) isKittyCapable() bool {
	switch p.rendererPref {
	case rendererANSI:
		return false
	case rendererKitty:
		return true
	default:
		return detectKittySupport()
	}
}
