package filepreview

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"strconv"

	"github.com/muesli/termenv"
	_ "golang.org/x/image/webp"
)

type ImageRenderer int

const (
	RendererANSI ImageRenderer = iota
	RendererKitty
)

type colorCache struct {
	rgbaToTermenv map[color.RGBA]termenv.RGBColor
}

func newColorCache() *colorCache {
	return &colorCache{
		rgbaToTermenv: make(map[color.RGBA]termenv.RGBColor),
	}
}

func (c *colorCache) getTermenvColor(col color.Color, fallbackColor string) termenv.RGBColor {
	rgba, ok := color.RGBAModel.Convert(col).(color.RGBA)
	if !ok || rgba.A == 0 {
		return termenv.RGBColor(fallbackColor)
	}

	if termenvColor, exists := c.rgbaToTermenv[rgba]; exists {
		return termenvColor
	}

	termenvColor := termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B))
	c.rgbaToTermenv[rgba] = termenvColor
	return termenvColor
}

// ConvertImageToANSI converts an image to ANSI escape codes with proper aspect ratio
func ConvertImageToANSI(img image.Image, defaultBGColor color.Color) string {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	output := ""
	cache := newColorCache()
	defaultBGHex := colorToHex(defaultBGColor)

	for y := 0; y < height; y += 2 {
		for x := range width {
			upperColor := cache.getTermenvColor(img.At(x, y), defaultBGHex)
			lowerColor := cache.getTermenvColor(defaultBGColor, "")

			if y+1 < height {
				lowerColor = cache.getTermenvColor(img.At(x, y+1), defaultBGHex)
			}

			// Using the "▄" character which fills the lower half
			cell := termenv.String("▄").Foreground(lowerColor).Background(upperColor)
			output += cell.String()
		}
		output += "\n"
	}

	return output
}

func ImagePreview(path string, maxWidth int, maxHeight int, defaultBGColor string) (string, error) {
	if isKittyCapable() {
		return ImagePreviewWithRenderer(path, maxWidth, maxHeight, defaultBGColor, RendererKitty)
	}
	return ImagePreviewWithRenderer(path, maxWidth, maxHeight, defaultBGColor, RendererANSI)
}

func ImagePreviewWithRenderer(path string, maxWidth int, maxHeight int, defaultBGColor string, renderer ImageRenderer) (string, error) {
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
		result, err := renderWithKitty(img, path, maxWidth, maxHeight, originalWidth, originalHeight)
		if err != nil {
			// If kitty fails, fall back to ANSI renderer
			slog.Error("Kitty renderer failed, falling back to ANSI", "error", err)
			bgColor, bgErr := hexToColor(defaultBGColor)
			if bgErr != nil {
				return "", fmt.Errorf("kitty failed and invalid background color: %w", bgErr)
			}
			// For ANSI fallback, fit the image to preview dimensions maintaining aspect ratio
			ansiImg := resizeForANSI(img, maxWidth, maxHeight)
			return ConvertImageToANSI(ansiImg, bgColor), nil
		}
		return result, nil

	case RendererANSI:
		fallthrough
	default:
		// Convert image to ANSI
		bgColor, err := hexToColor(defaultBGColor)
		if err != nil {
			return "", fmt.Errorf("invalid background color: %w", err)
		}

		// For ANSI rendering, resize image appropriately
		fittedImg := resizeForANSI(img, maxWidth, maxHeight)
		return ConvertImageToANSI(fittedImg, bgColor), nil
	}
}

func hexToColor(hex string) (color.RGBA, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return color.RGBA{}, errors.New("invalid hex color format")
	}
	values, err := strconv.ParseUint(hex[1:], 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}, nil
}

func colorToHex(color color.Color) string {
	r, g, b, _ := color.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}
