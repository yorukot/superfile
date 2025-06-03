package filepreview

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"log/slog"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/muesli/termenv"
	"github.com/rwcarlsen/goexif/exif"
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
	var output strings.Builder
	cache := newColorCache()
	defaultBGHex := colorToHex(defaultBGColor)

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			bg := cache.getTermenvColor(img.At(x, y), defaultBGHex)
			cell := termenv.String(" ").Background(bg)
			output.WriteString(cell.String())
		}
		output.WriteString("\n")
	}

	return output.String()
}

func ImagePreview(path string, maxWidth int, maxHeight int, defaultBGColor string) (string, error) {

	return ImagePreviewWithRenderer(path, maxWidth, maxHeight, defaultBGColor, RendererANSI)
}

func ImagePreviewWithRenderer(path string, maxWidth int, maxHeight int, defaultBGColor string, renderer ImageRenderer) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	const maxFileSize = 100 * 1024 * 1024
	if info.Size() > maxFileSize {
		return "", fmt.Errorf("image file too large: %d bytes", info.Size())
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	imgReader := bytes.NewReader(data)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		return "", err
	}

	exifReader := bytes.NewReader(data)
	img = adjustImageOrientation(exifReader, img)

	resizedImg := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)

	switch renderer {
	case RendererANSI:
		fallthrough
	default:
		bgColor, err := hexToColor(defaultBGColor)
		if err != nil {
			return "", fmt.Errorf("invalid background color: %w", err)
		}
		return ConvertImageToANSI(resizedImg, bgColor), nil
	}
}

func adjustImageOrientation(r *bytes.Reader, img image.Image) image.Image {
	exifData, err := exif.Decode(r)
	if err != nil {
		slog.Error("exif error", "error", err)
		return img
	}
	tag, err := exifData.Get(exif.Orientation)
	if err != nil {
		slog.Error("exif orientation error", "error", err)
		return img
	}
	orientation, err := tag.Int(0)
	if err != nil {
		slog.Error("exif orientation value error", "error", err)
		return img
	}
	return adjustOrientation(img, orientation)
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
		slog.Error("Invalid orientation value", "error", orientation)
		return img
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
