package filepreview

import (
	"fmt"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log/slog"
	"os"
	"time"

	_ "golang.org/x/image/webp" // Register WebP decoder

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/cache"
)

type ImageRenderer int

const (
	RendererANSI ImageRenderer = iota
	RendererKitty
)

func (f ImageRenderer) String() string {
	switch f {
	case RendererANSI:
		return "ANSI"
	case RendererKitty:
		return "Kitty"
	default:
		return common.InvalidTypeString
	}
}

func getPreviewObjKey(path string, dim string, renderer ImageRenderer) string {
	return fmt.Sprintf("%s:%s:%s", path, dim, renderer)
}

// ImagePreviewer encapsulates image preview functionality with caching
type ImagePreviewer struct {
	cache       *cache.Cache[string]
	terminalCap *TerminalCapabilities
}

// NewImagePreviewer creates a new ImagePreviewer with default cache settings
func NewImagePreviewer() *ImagePreviewer {
	return NewImagePreviewerWithConfig(defaultImagePreviewCacheSize, defaultCacheExpiration)
}

// NewImagePreviewerWithConfig creates a new ImagePreviewer with custom cache configuration
func NewImagePreviewerWithConfig(maxEntries int, expiration time.Duration) *ImagePreviewer {
	previewer := &ImagePreviewer{
		cache:       cache.New[string](maxEntries, expiration),
		terminalCap: NewTerminalCapabilities(),
	}

	// Initialize terminal capabilities
	previewer.terminalCap.InitTerminalCapabilities()

	return previewer
}

// ImagePreview generates a preview of an image file
func (p *ImagePreviewer) ImagePreview(path string, maxWidth int, maxHeight int,
	defaultBGColor string, sideAreaWidth int) (string, error) {
	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return "", fmt.Errorf("dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	// Create dimensions string for cache key
	dimensions := fmt.Sprintf("%d,%d,%s,%d", maxWidth, maxHeight, defaultBGColor, sideAreaWidth)

	// Try Kitty first as it's more modern
	if p.IsKittyCapable() {
		// Check cache for Kitty renderer
		if preview, exists := p.cache.Get(getPreviewObjKey(path, dimensions, RendererKitty)); exists {
			return preview, nil
		}

		preview, err := p.ImagePreviewWithRenderer(
			path,
			maxWidth,
			maxHeight,
			defaultBGColor,
			RendererKitty,
			sideAreaWidth,
		)
		if err == nil {
			// Cache the successful result
			p.cache.Set(getPreviewObjKey(path, dimensions, RendererKitty), preview)
			return preview, nil
		}

		// Fall through to ANSI if Kitty fails
		slog.Error("Kitty renderer failed, falling back to ANSI", "error", err)
	}

	// Check cache for ANSI renderer
	if preview, found := p.cache.Get(getPreviewObjKey(path, dimensions, RendererANSI)); found {
		return preview, nil
	}

	// Fall back to ANSI
	preview, err := p.ImagePreviewWithRenderer(path, maxWidth, maxHeight, defaultBGColor, RendererANSI, sideAreaWidth)
	if err == nil {
		// Cache the successful result
		p.cache.Set(getPreviewObjKey(path, dimensions, RendererANSI), preview)
	}
	return preview, err
}

// ImagePreviewWithRenderer generates an image preview using the specified renderer
func (p *ImagePreviewer) ImagePreviewWithRenderer(path string, maxWidth int, maxHeight int,
	defaultBGColor string, renderer ImageRenderer, sideAreaWidth int) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	const maxFileSize = 100 * 1024 * 1024 // 100MB limit
	if info.Size() > maxFileSize {
		return "", fmt.Errorf("image file too large: %d bytes", info.Size())
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Use the new image preparation pipeline
	img, originalWidth, originalHeight, err := prepareImageForPreview(data)
	if err != nil {
		return "", err
	}

	switch renderer {
	case RendererKitty:
		result, err := p.renderWithKittyUsingTermCap(img, originalWidth,
			originalHeight, maxWidth, maxHeight, sideAreaWidth)
		if err != nil {
			// If kitty fails, fall back to ANSI renderer
			slog.Error("Kitty renderer failed, falling back to ANSI", "error", err)
			return p.ANSIRenderer(img, defaultBGColor, maxWidth, maxHeight)
		}
		return result, nil

	case RendererANSI:
		return p.ANSIRenderer(img, defaultBGColor, maxWidth, maxHeight)
	default:
		return "", fmt.Errorf("invalid renderer : %v", renderer)
	}
}
