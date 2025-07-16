package filepreview

import (
	"log/slog"
	"runtime"
	"sync"
)

// Terminal cell to pixel conversion constants
// These approximate the pixel dimensions of terminal cells
const (
	DefaultPixelsPerColumn = 10 // approximate pixels per terminal column
	DefaultPixelsPerRow    = 20 // approximate pixels per terminal row
)

// TerminalCellSize represents the pixel dimensions of terminal cells
type TerminalCellSize struct {
	PixelsPerColumn int
	PixelsPerRow    int
}

// TerminalCapabilities encapsulates terminal capability detection
type TerminalCapabilities struct {
	cellSize       TerminalCellSize
	cellSizeInit   sync.Once
	detectionMutex sync.Mutex
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
			tc.cellSize = tc.detectTerminalCellSize()
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
		tc.cellSize = tc.detectTerminalCellSize()
		slog.Info("Terminal cell size detection (lazy init)",
			"pixels_per_column", tc.cellSize.PixelsPerColumn,
			"pixels_per_row", tc.cellSize.PixelsPerRow)
	})

	return tc.cellSize
}

// winsize struct for ioctl TIOCGWINSZ
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// detectTerminalCellSize detects the terminal cell size using ioctl system calls
// This method is non-blocking and doesn't interfere with stdin
func (tc *TerminalCapabilities) detectTerminalCellSize() TerminalCellSize {
	tc.detectionMutex.Lock()
	defer tc.detectionMutex.Unlock()

	// Try platform-specific detection
	if runtime.GOOS == "windows" {
		if cellSize, ok := getTerminalCellSizeWindows(); ok {
			slog.Info("Successfully detected terminal cell size on Windows",
				"pixels_per_column", cellSize.PixelsPerColumn,
				"pixels_per_row", cellSize.PixelsPerRow)
			return cellSize
		}
	} else {
		// Unix-like systems (Linux, macOS, etc.)
		if cellSize, ok := getTerminalCellSizeViaIoctl(); ok {
			slog.Info("Successfully detected terminal cell size via ioctl",
				"pixels_per_column", cellSize.PixelsPerColumn,
				"pixels_per_row", cellSize.PixelsPerRow)
			return cellSize
		}
	}

	// Fallback to default values
	slog.Info("Using default terminal cell size", "os", runtime.GOOS)
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

// Windows-specific terminal detection functions
// getTerminalCellSizeWindows uses Windows Console API to detect terminal cell size
func getTerminalCellSizeWindows() (TerminalCellSize, bool) {
	if runtime.GOOS != "windows" {
		return TerminalCellSize{}, false
	}

	// For Windows, just return reasonable defaults
	// Windows terminal detection is complex and varies greatly between
	// different terminal emulators (Windows Terminal, ConEmu, etc.)
	slog.Info("Using Windows default terminal cell size")
	// TODO: Implement actual Windows Console API calls when running on Windows
	return getWindowsDefaultCellSize(), true
}

// getWindowsDefaultCellSize returns reasonable defaults for Windows
func getWindowsDefaultCellSize() TerminalCellSize {
	return TerminalCellSize{
		PixelsPerColumn: 8,  // Windows Terminal/CMD typical width
		PixelsPerRow:    16, // Windows Terminal/CMD typical height
	}
}
