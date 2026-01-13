package filepreview

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/charmbracelet/lipgloss"
)

// isInlineCapable checks if the terminal supports inline image protocol (iTerm2, WezTerm, etc.)
func (p *ImagePreviewer) IsInlineCapable() bool {
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

func itermInlineSeq(name string, data []byte, width, height string) string {
	nameEnc := base64.StdEncoding.EncodeToString([]byte(name))
	dataB64 := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf(
		"\033]1337;File=name=%s;inline=1;width=%s;height=%s;preserveAspectRatio=1:%s\a",
		nameEnc, width, height, dataB64,
	)
}

// renderWithInlineUsingTermCap renders an image using inline image protocol
func (p *ImagePreviewer) renderWithInlineUsingTermCap(
	img image.Image,
	path string,
	originalWidth, originalHeight, maxWidth, maxHeight int,
	sideAreaWidth int,
	filePreviewStyle lipgloss.Style,
) (string, error) {
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	slog.Debug("inline renderer starting", "path", path, "maxWidth", maxWidth, "maxHeight", maxHeight)

	w, h := originalWidth, originalHeight
	if (w <= 0 || h <= 0) && img != nil {
		b := img.Bounds()
		w, h = b.Dx(), b.Dy()
	}
	if w <= 0 || h <= 0 {
		return "", fmt.Errorf("invalid original dimensions (w=%d, h=%d)", w, h)
	}

	imgRatio := float64(w) / float64(h)

	cellSize := p.terminalCap.GetTerminalCellSize()
	ppc := cellSize.PixelsPerColumn
	ppr := cellSize.PixelsPerRow

	usePixelRatio := ppc > 0 && ppr > 0

	displayW, displayH := 1, 1

	if usePixelRatio {
		termRatio := float64(maxWidth*ppc) / float64(maxHeight*ppr)

		if imgRatio > termRatio {
			displayW = maxWidth
			widthPx := float64(displayW * ppc)
			heightPx := widthPx / imgRatio
			displayH = int(heightPx / float64(ppr))
		} else {
			displayH = maxHeight
			heightPx := float64(displayH * ppr)
			widthPx := heightPx * imgRatio
			displayW = int(widthPx / float64(ppc))
		}
	} else {
		termRatio := float64(maxWidth) / float64(maxHeight)

		if imgRatio > termRatio {
			displayW = maxWidth
			displayH = int(float64(maxWidth) / imgRatio)
		} else {
			displayH = maxHeight
			displayW = int(float64(maxHeight) * imgRatio)
		}
	}

	// clamp
	if displayW < 1 {
		displayW = 1
	}
	if displayH < 1 {
		displayH = 1
	}
	if displayW > maxWidth {
		displayW = maxWidth
	}
	if displayH > maxHeight {
		displayH = maxHeight
	}

	var buf bytes.Buffer
	line := strings.Repeat(" ", maxWidth)
	block := strings.Repeat(line+"\n", maxHeight-1) + line
	buf.WriteString(filePreviewStyle.Render(block))

	var data []byte
	var name string

	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read image file: %w", err)
	}
	data = b
	name = filepath.Base(path)

	seq := itermInlineSeq(
		name,
		data,
		strconv.Itoa(displayW),
		strconv.Itoa(displayH),
	)

	buf.WriteString("\x1b[s")
	buf.WriteString(fmt.Sprintf("\x1b[1;%dH", sideAreaWidth))
	buf.WriteString(seq)
	buf.WriteString("\x1b[u")

	return buf.String(), nil
}
