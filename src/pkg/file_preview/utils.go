//go:build !windows
// Sys do not support terminal detection on windows

package filepreview

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

// Terminal cell to pixel conversion constants
// These approximate the pixel dimensions of terminal cells
const (
	DefaultPixelsPerColumn = 10 // approximate pixels per terminal column
	DefaultPixelsPerRow    = 20 // approximate pixels per terminal row
)

// TerminalCellSize represents the pixel dimensions of a terminal cell
type TerminalCellSize struct {
	PixelsPerColumn int
	PixelsPerRow    int
}

// TerminalCapabilities encapsulates terminal capability detection
type TerminalCapabilities struct {
	cellSize     TerminalCellSize
	cellSizeInit sync.Once
}

// NewTerminalCapabilities creates a new TerminalCapabilities instance
func NewTerminalCapabilities() *TerminalCapabilities {
	return &TerminalCapabilities{
		cellSize: TerminalCellSize{
			PixelsPerColumn: DefaultPixelsPerColumn,
			PixelsPerRow:    DefaultPixelsPerRow,
		},
	}
}

// InitTerminalCapabilities initializes all terminal capabilities detection
// including cell size and Kitty Graphics Protocol support
// This should be called early in the application startup
func (tc *TerminalCapabilities) InitTerminalCapabilities() {
	// Use a goroutine to avoid blocking the application startup
	go func() {
		// Initialize cell size detection
		tc.cellSizeInit.Do(func() {
			tc.cellSize = DetectTerminalCellSize()
			slog.Info("Terminal cell size detection",
				"pixels_per_column", tc.cellSize.PixelsPerColumn,
				"pixels_per_row", tc.cellSize.PixelsPerRow)
		})
	}()
}

// GetTerminalCellSize returns the current terminal cell size
// If detection hasn't been initialized, it performs detection first
func (tc *TerminalCapabilities) GetTerminalCellSize() TerminalCellSize {
	tc.cellSizeInit.Do(func() {
		tc.cellSize = DetectTerminalCellSize()
		slog.Info("Terminal cell size detection (lazy init)",
			"pixels_per_column", tc.cellSize.PixelsPerColumn,
			"pixels_per_row", tc.cellSize.PixelsPerRow)
	})

	return tc.cellSize
}

func DetectTerminalCellSize() TerminalCellSize {
	fd := int(os.Stdin.Fd())

	// Check if we're in a terminal
	if !isTerminal(fd) {
		slog.Warn("Stdin is not a terminal; using default cell size")
		return getDefaultCellSize()
	}

	// Save cursor position
	if _, err := os.Stdout.WriteString("\x1b[s"); err != nil {
		slog.Error("Error saving terminal state", "error", err)
		return getDefaultCellSize()
	}
	defer func() {
		if _, err := os.Stdout.WriteString("\x1b[u"); err != nil {
			slog.Error("Error restoring terminal state", "error", err)
		}
	}()

	// Send CSI 16t query to request terminal pixel cell size
	if _, err := os.Stdout.WriteString("\x1b[16t"); err != nil {
		slog.Error("Error requesting terminal cell size", "error", err)
		return getDefaultCellSize()
	}

	// Read response with timeout (blocking read in goroutine, limited by select)
	responseChan := make(chan string, 1)

	go func() {
		buf := make([]byte, 64)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			slog.Error("Error reading terminal response", "error", err)
			responseChan <- ""
			return
		}
		responseChan <- string(buf[:n])
	}()

	var response string

	select {
	case response = <-responseChan:
		slog.Debug("Received terminal response", "raw_response", fmt.Sprintf("%q", response))
	case <-time.After(100 * time.Millisecond):
		slog.Warn("Timeout waiting for terminal response")
		return getDefaultCellSize()
	}

	if width, height, ok := parseCSI16tResponse(response); ok {
		slog.Debug("Successfully parsed terminal cell size",
			"width", width,
			"height", height)
		return TerminalCellSize{
			PixelsPerColumn: width,
			PixelsPerRow:    height,
		}
	}

	slog.Warn("Failed to parse terminal response, using defaults")
	return getDefaultCellSize()
}

// isTerminal checks whether a file descriptor refers to a terminal
func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	return err == nil
}

// getDefaultCellSize returns default fallback terminal cell size
func getDefaultCellSize() TerminalCellSize {
	return TerminalCellSize{
		PixelsPerColumn: DefaultPixelsPerColumn,
		PixelsPerRow:    DefaultPixelsPerRow,
	}
}

// parseCSI16tResponse parses the CSI 16t response from the terminal
// The format is: ESC[6;height;widtht
func parseCSI16tResponse(response string) (int, int, bool) {
	if !strings.Contains(response, "\x1b[6;") {
		return 0, 0, false
	}

	parts := strings.Split(strings.TrimPrefix(response, "\x1b[6;"), ";")
	if len(parts) < 2 {
		return 0, 0, false
	}

	// Extract height and width
	heightStr := parts[0]
	widthParts := strings.Split(parts[1], "t")
	if len(widthParts) < 1 {
		return 0, 0, false
	}
	widthStr := widthParts[0]

	// Parse as integers
	h, err1 := strconv.Atoi(heightStr)
	w, err2 := strconv.Atoi(widthStr)

	if err1 != nil || err2 != nil || h <= 0 || w <= 0 {
		return 0, 0, false
	}

	return w, h, true
}

// InitTerminalCapabilities initializes terminal capabilities for the ImagePreviewer
func (p *ImagePreviewer) InitTerminalCapabilities() {
	p.terminalCap.InitTerminalCapabilities()
}
