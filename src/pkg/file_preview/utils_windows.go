//go:build windows

// Sys do not support terminal detection on windows

package filepreview

import (
	"log/slog"
	"sync"
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
	// On Windows, we just use default values
	tc.cellSizeInit.Do(func() {
		tc.cellSize = getDefaultCellSize()
		slog.Info("Terminal cell size detection (Windows - using defaults)",
			"pixels_per_column", tc.cellSize.PixelsPerColumn,
			"pixels_per_row", tc.cellSize.PixelsPerRow)
	})
}

// GetTerminalCellSize returns the current terminal cell size
// If detection hasn't been initialized, it performs detection first
func (tc *TerminalCapabilities) GetTerminalCellSize() TerminalCellSize {
	tc.cellSizeInit.Do(func() {
		tc.cellSize = getDefaultCellSize()
		slog.Info("Terminal cell size detection (Windows - lazy init, using defaults)",
			"pixels_per_column", tc.cellSize.PixelsPerColumn,
			"pixels_per_row", tc.cellSize.PixelsPerRow)
	})

	return tc.cellSize
}

func DetectTerminalCellSize() TerminalCellSize {
	// On Windows, always return default values
	slog.Info("Terminal cell size detection not supported on Windows, using defaults")
	return getDefaultCellSize()
}

// getDefaultCellSize returns default fallback terminal cell size
func getDefaultCellSize() TerminalCellSize {
	return TerminalCellSize{
		PixelsPerColumn: DefaultPixelsPerColumn,
		PixelsPerRow:    DefaultPixelsPerRow,
	}
}

// InitTerminalCapabilities initializes terminal capabilities for the ImagePreviewer
func (p *ImagePreviewer) InitTerminalCapabilities() {
	p.terminalCap.InitTerminalCapabilities()
}
