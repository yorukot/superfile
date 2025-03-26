package filepreview

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Import to enable GIF support
	_ "image/jpeg" // Import to enable JPEG support
	_ "image/png"  // Import to enable PNG support
	"log/slog"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
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
		for x := 0; x < width; x++ {
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

// Return image preview ansi string
func ImagePreview(path string, maxWidth, maxHeight int, defaultBGColor string) (string, error) {
	// Load image file
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Seek back to the beginning of the file before reading EXIF data
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Try to adjust image orientation based on EXIF data
	img = adjustImageOrientation(file, img)

	// Resize image to fit terminal
	resizedImg := resize.Thumbnail(uint(maxWidth), uint(maxHeight), img, resize.Lanczos3)

	// Convert image to ANSI
	bgColor, err := hexToColor(defaultBGColor)
	if err != nil {
		return "", fmt.Errorf("invalid background color: %v", err)
	}
	ansiImage := ConvertImageToANSI(resizedImg, bgColor)

	return ansiImage, nil
}

func adjustImageOrientation(file *os.File, img image.Image) image.Image {
	exifData, err := exif.Decode(file)
	if err != nil {
		slog.Error("exif error", "error", err)
		return img
	}

	orientation, err := exifData.Get(exif.Orientation)
	if err != nil {
		slog.Error("exif orientation error", "error", err)
		return img
	}

	orientationValue, err := orientation.Int(0)
	if err != nil {
		slog.Error("exif orientation value error", "error", err)
		return img
	}

	return adjustOrientation(img, orientationValue)
}

func adjustOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		return img
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		slog.Error("invalid orientation value", "orientation", orientation)
		return img // Return the original image for invalid orientation values
	}
}

func hexToColor(hex string) (color.RGBA, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return color.RGBA{}, fmt.Errorf("invalid hex color format")
	}
	values, err := strconv.ParseUint(string(hex[1:]), 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}, nil
}

func colorToHex(color color.Color) (fullbackHex string) {
	r, g, b, _ := color.RGBA()

	fullbackHex = fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
	return fullbackHex

}
