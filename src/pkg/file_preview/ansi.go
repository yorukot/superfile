package filepreview

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/muesli/termenv"
)

// ConvertImageToANSI converts an image to ANSI escape codes with proper aspect ratio
func ConvertImageToANSI(img image.Image, defaultBGColor color.Color) string {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// TODO: Use renderer here to prevent newline management,and overflows
	var output strings.Builder
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
			output.WriteString(cell.String())
		}
		// Only add newline if this is not the last row
		if y+2 < height {
			output.WriteByte('\n')
		}
	}

	return output.String()
}

// Convert image to ansi
func (p *ImagePreviewer) ANSIRenderer(img image.Image, defaultBGColor string,
	maxWidth int, maxHeight int) (string, error) {
	bgColor, err := hexToColor(defaultBGColor)
	if err != nil {
		return "", fmt.Errorf("invalid background color: %w", err)
	}

	// For ANSI rendering, resize image appropriately
	fittedImg := resizeForANSI(img, maxWidth, maxHeight)
	return ConvertImageToANSI(fittedImg, bgColor), nil
}

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
