package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/kitty"
)

// isKittyCapable checks if the terminal supports Kitty graphics protocol
func isKittyCapable() bool {
	termProgram := os.Getenv("TERM_PROGRAM")
	term := os.Getenv("TERM")

	// This allowlist is only a provisional guess, used until the terminal
	// answers KittyGraphicsQuery(). It cannot see through tmux, which masks
	// both variables with its own values.
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

// Kitty graphics support, as reported by the terminal itself.
const (
	kittyCapUnknown int32 = iota
	kittyCapSupported
	kittyCapUnsupported
)

// kittyCapability holds the terminal's answer to KittyGraphicsQuery(). It stays
// at kittyCapUnknown until a reply arrives, during which isKittyCapable()'s
// allowlist serves as a provisional guess.
var kittyCapability atomic.Int32

// KittyGraphicsQuery returns the escape sequences that ask the terminal whether
// it supports the kitty graphics protocol. Send it once at startup via
// tea.Raw(); the replies arrive as messages (see MarkKittySupported and
// MarkKittyUnsupportedIfUnknown).
//
// It transmits a 1x1 RGB image with a=q, which asks the terminal to report
// support without displaying anything, followed by a primary device attributes
// request. Every terminal answers DA1, so a DA1 reply arriving without a
// graphics reply means the protocol is unsupported.
//
// Both are wrapped in a single tmux passthrough so that the host terminal, not
// tmux, answers both — otherwise tmux would answer DA1 itself and the fence
// could beat the host's graphics reply back.
func KittyGraphicsQuery() string {
	query := ansi.KittyGraphics([]byte("AAAA"), "i=31", "s=1", "v=1", "a=q", "t=d", "f=24")
	return tmuxPassthrough(query + ansi.RequestPrimaryDeviceAttributes)
}

// MarkKittySupported records a graphics reply from the terminal. It reports
// whether this changed the known capability, i.e. whether previews rendered
// before now need redrawing.
func MarkKittySupported() bool {
	return kittyCapability.Swap(kittyCapSupported) != kittyCapSupported
}

// MarkKittyUnsupportedIfUnknown records the DA1 fence arriving with no graphics
// reply before it. A positive reply always wins, so this only downgrades a
// still-unknown capability. It reports whether it changed anything.
func MarkKittyUnsupportedIfUnknown() bool {
	return kittyCapability.CompareAndSwap(kittyCapUnknown, kittyCapUnsupported)
}

// tmuxPassthrough wraps each escape sequence in s in tmux's DCS passthrough,
// so tmux forwards it to the host terminal instead of discarding it as an
// unrecognised sequence. Requires "set -g allow-passthrough on" in tmux.
//
// Each APC is wrapped individually rather than wrapping the whole buffer in one
// DCS: image data is chunked into many sequences and a single multi-megabyte DCS
// risks hitting tmux's buffer limits.
func tmuxPassthrough(s string) string {
	if os.Getenv("TMUX") == "" || s == "" {
		return s
	}

	var b strings.Builder
	for len(s) > 0 {
		seq := s
		// Sequences are terminated by ST (ESC \); keep the terminator with its
		// sequence. A trailing fragment without ST is passed through as-is.
		if i := strings.Index(s, "\x1b\\"); i >= 0 {
			seq, s = s[:i+2], s[i+2:]
		} else {
			s = ""
		}
		b.WriteString("\x1bPtmux;")
		// tmux strips one level of escaping, so every ESC must be doubled.
		b.WriteString(strings.ReplaceAll(seq, "\x1b", "\x1b\x1b"))
		b.WriteString("\x1b\\")
	}
	return b.String()
}

// GetKittyClearRaw returns the raw APC command to clear all Kitty images.
// This must be sent via tea.Raw(), not embedded in view content.
func (p *ImagePreviewer) GetKittyClearRaw() string {
	if !p.IsKittyCapable() {
		return ""
	}
	return tmuxPassthrough(ansi.KittyGraphics(nil, "a=d"))
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
		RawTransmit:  tmuxPassthrough(transmitBuf.String()),
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

// IsKittyCapable reports whether the terminal supports the Kitty graphics
// protocol, preferring the terminal's own answer to KittyGraphicsQuery() and
// falling back to the $TERM/$TERM_PROGRAM allowlist until that answer arrives.
func (p *ImagePreviewer) IsKittyCapable() bool {
	switch kittyCapability.Load() {
	case kittyCapSupported:
		return true
	case kittyCapUnsupported:
		return false
	default:
		return isKittyCapable()
	}
}
