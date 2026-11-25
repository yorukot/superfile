package filepreview

import (
	"fmt"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log/slog"
	"os"
	"strings"
	"time"

	_ "golang.org/x/image/webp" // Register WebP decoder

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/cache"
)

type ImageRenderer int

const (
	RendererANSI ImageRenderer = iota
	RendererKitty
	RendererSixel
)

// Values of the image_renderer config option ("" or "auto" means detect).
// Validated in common.ValidateConfig.
const (
	ImageRendererConfigKitty = "kitty"
	ImageRendererConfigSixel = "sixel"
	ImageRendererConfigANSI  = "ansi"
)

func (f ImageRenderer) String() string {
	switch f {
	case RendererANSI:
		return "ANSI"
	case RendererKitty:
		return "Kitty"
	case RendererSixel:
		return "Sixel"
	default:
		return common.InvalidTypeString
	}
}

func getPreviewObjKey(path string, dim string, renderer ImageRenderer) string {
	return fmt.Sprintf("%s:%s:%s", path, dim, renderer)
}

// ImagePreviewResult is a rendered image preview.
type ImagePreviewResult struct {
	// Render is the content for the cell-based view: Kitty placeholders,
	// ANSI half-blocks, or the blank block with marker cells for Sixel.
	Render string
	// RawTransmit is Kitty graphics data. It must be sent via tea.Raw().
	RawTransmit string
	// Sixel is the Sixel image, which must be drawn after the frame
	// containing Render. See SixelImage.Placement.
	Sixel *SixelImage
}

// ImagePreviewer encapsulates image preview functionality with caching
type ImagePreviewer struct {
	cache       *cache.Cache[string]
	sixelCache  *cache.Cache[*SixelImage]
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
		sixelCache:  cache.New[*SixelImage](maxEntries, expiration),
		terminalCap: NewTerminalCapabilities(),
	}

	// Initialize terminal capabilities
	previewer.terminalCap.InitTerminalCapabilities()

	return previewer
}

// SelectedRenderer returns the renderer to use: the image_renderer config
// option if set, otherwise the best one the terminal supports
// (Kitty > Sixel > ANSI).
func (p *ImagePreviewer) SelectedRenderer() ImageRenderer {
	return selectRenderer(common.Config.ImageRenderer, isKittyCapable(), isSixelCapable())
}

func selectRenderer(configValue string, kittyCapable, sixelCapable bool) ImageRenderer {
	switch strings.ToLower(configValue) {
	case ImageRendererConfigKitty:
		return RendererKitty
	case ImageRendererConfigSixel:
		return RendererSixel
	case ImageRendererConfigANSI:
		return RendererANSI
	}
	if kittyCapable {
		return RendererKitty
	}
	if sixelCapable {
		return RendererSixel
	}
	return RendererANSI
}

// ImagePreview generates a preview of an image file, using the renderer from
// SelectedRenderer and falling back to ANSI if that fails.
func (p *ImagePreviewer) ImagePreview(path string, maxWidth int, maxHeight int,
	defaultBGColor string, sideAreaWidth int) (ImagePreviewResult, error) {
	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return ImagePreviewResult{}, fmt.Errorf(
			"dimensions must be positive (maxWidth=%d, maxHeight=%d)", maxWidth, maxHeight)
	}

	// Create dimensions string for cache key
	dimensions := fmt.Sprintf("%d,%d,%s,%d", maxWidth, maxHeight, defaultBGColor, sideAreaWidth)

	if renderer := p.SelectedRenderer(); renderer != RendererANSI {
		if result, ok := p.getCached(path, dimensions, renderer); ok {
			return result, nil
		}

		result, err := p.ImagePreviewWithRenderer(path, maxWidth, maxHeight, defaultBGColor, renderer, sideAreaWidth)
		if err == nil && (result.RawTransmit != "" || result.Sixel != nil) {
			p.setCached(path, dimensions, renderer, result)
			return result, nil
		}
		if err == nil {
			// The renderer fell back to ANSI. Return it directly; don't
			// pollute this renderer's cache with ANSI results.
			return result, nil
		}

		// Fall through to ANSI if the renderer fails
		slog.Error("Image renderer failed, falling back to ANSI", "renderer", renderer, "error", err)
	}

	// Check cache for ANSI renderer
	if preview, found := p.cache.Get(getPreviewObjKey(path, dimensions, RendererANSI)); found {
		return ImagePreviewResult{Render: preview}, nil
	}

	// Fall back to ANSI
	result, err := p.ImagePreviewWithRenderer(path, maxWidth, maxHeight, defaultBGColor, RendererANSI, sideAreaWidth)
	if err == nil {
		// Cache the successful result
		p.cache.Set(getPreviewObjKey(path, dimensions, RendererANSI), result.Render)
	}
	return result, err
}

// getCached returns a cached Kitty or Sixel result. Both parts must be
// present, since they are cached and evicted separately.
func (p *ImagePreviewer) getCached(path, dimensions string, renderer ImageRenderer) (ImagePreviewResult, bool) {
	cacheKey := getPreviewObjKey(path, dimensions, renderer)
	render, exists := p.cache.Get(cacheKey)
	if !exists {
		return ImagePreviewResult{}, false
	}
	switch renderer {
	case RendererKitty:
		if raw, ok := p.cache.Get(cacheKey + ":raw"); ok && raw != "" {
			return ImagePreviewResult{Render: render, RawTransmit: raw}, true
		}
	case RendererSixel:
		if img, ok := p.sixelCache.Get(cacheKey); ok && img != nil {
			return ImagePreviewResult{Render: render, Sixel: img}, true
		}
	case RendererANSI:
	}
	return ImagePreviewResult{}, false
}

func (p *ImagePreviewer) setCached(path, dimensions string, renderer ImageRenderer, result ImagePreviewResult) {
	cacheKey := getPreviewObjKey(path, dimensions, renderer)
	p.cache.Set(cacheKey, result.Render)
	switch renderer {
	case RendererKitty:
		p.cache.Set(cacheKey+":raw", result.RawTransmit)
	case RendererSixel:
		p.sixelCache.Set(cacheKey, result.Sixel)
	case RendererANSI:
	}
}

// ImagePreviewWithRenderer generates an image preview using the specified
// renderer. If a Kitty or Sixel renderer fails, it falls back to ANSI.
func (p *ImagePreviewer) ImagePreviewWithRenderer(path string, maxWidth int, maxHeight int,
	defaultBGColor string, renderer ImageRenderer, sideAreaWidth int) (ImagePreviewResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return ImagePreviewResult{}, err
	}
	const maxFileSize = 100 * 1024 * 1024 // 100MB limit
	if info.Size() > maxFileSize {
		return ImagePreviewResult{}, fmt.Errorf("image file too large: %d bytes", info.Size())
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ImagePreviewResult{}, err
	}

	// Use the new image preparation pipeline
	img, originalWidth, originalHeight, err := prepareImageForPreview(data)
	if err != nil {
		return ImagePreviewResult{}, err
	}

	switch renderer {
	case RendererKitty:
		result, err := p.renderWithKittyUsingTermCap(img, path, originalWidth,
			originalHeight, maxWidth, maxHeight, sideAreaWidth)
		if err != nil {
			slog.Error("Kitty renderer failed, falling back to ANSI", "error", err)
			render, ansiErr := p.ANSIRenderer(img, defaultBGColor, maxWidth, maxHeight)
			return ImagePreviewResult{Render: render}, ansiErr
		}
		return ImagePreviewResult{Render: result.Placeholders, RawTransmit: result.RawTransmit}, nil

	case RendererSixel:
		block, sixelImg, err := p.renderWithSixel(img, maxWidth, maxHeight)
		if err != nil {
			slog.Error("Sixel renderer failed, falling back to ANSI", "error", err)
			render, ansiErr := p.ANSIRenderer(img, defaultBGColor, maxWidth, maxHeight)
			return ImagePreviewResult{Render: render}, ansiErr
		}
		return ImagePreviewResult{Render: block, Sixel: sixelImg}, nil

	case RendererANSI:
		render, err := p.ANSIRenderer(img, defaultBGColor, maxWidth, maxHeight)
		return ImagePreviewResult{Render: render}, err
	default:
		return ImagePreviewResult{}, fmt.Errorf("invalid renderer : %v", renderer)
	}
}
