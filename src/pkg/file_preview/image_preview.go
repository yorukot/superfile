package filepreview

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Import to enable GIF support
	_ "image/jpeg" // Import to enable JPEG support
	_ "image/png"  // Import to enable PNG support
	"os"
	"strconv"

	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
)

// ConvertImageToANSI converts an image to ANSI escape codes with proper aspect ratio
func ConvertImageToANSI(img image.Image, defaultBGColor color.Color) string {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	output := ""

	r, g, b, _ := defaultBGColor.RGBA()

	BGColor := fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			upperColor := colorToTermenv(img.At(x, y), BGColor)
			lowerColor := termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8)))

			if y + 1 < height {
				lowerColor = colorToTermenv(img.At(x, y + 1), BGColor)
			}

			// Using the "▄" character which fills the lower half
			cell := termenv.String("▄").Foreground(lowerColor).Background(upperColor)
			output += cell.String()
		}
		output += "\n"
	}

	return output
}

// colorToTermenv converts a color.Color to a termenv.RGBColor
func colorToTermenv(c color.Color, BGColor string) termenv.RGBColor {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return termenv.RGBColor(BGColor)
	}

	return termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8)))
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

	// Resize image to fit terminal
	resizedImg := resize.Thumbnail(uint(maxWidth), uint(maxHeight), img, resize.Lanczos3)

	// Convert image to ANSI
	ansiImage := ConvertImageToANSI(resizedImg, hexToColor(defaultBGColor))

	return ansiImage, nil
}

func hexToColor(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}