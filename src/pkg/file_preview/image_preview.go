package filepreview

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"log/slog"
	"os"
	"strconv"

	"github.com/BourgeoisBear/rasterm"
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
	// Enhanced kitty capability detection for Ghostty and other terminals
	isKittyCapable := rasterm.IsKittyCapable()

	// Additional detection for terminals that might not be detected by rasterm
	if !isKittyCapable {
		// Check for Ghostty specifically
		termProgram := os.Getenv("TERM_PROGRAM")
		term := os.Getenv("TERM")

		// Log environment for debugging
		slog.Info("Terminal detection",
			"TERM_PROGRAM", termProgram,
			"TERM", term,
			"rasterm_detected", isKittyCapable,
		)

		// Manual detection for additional terminals
		if termProgram == "ghostty" ||
			term == "xterm-ghostty" ||
			term == "ghostty" ||
			termProgram == "Ghostty" {
			isKittyCapable = true
			slog.Info("Manually detected Ghostty terminal, enabling kitty protocol")
		}
	}

	if isKittyCapable {
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
	imgReader := bytes.NewReader(data)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		return "", err
	}

	exifReader := bytes.NewReader(data)
	img = adjustImageOrientation(exifReader, img)

	switch renderer {
	case RendererKitty:
		// Create a buffer to capture the kitty output
		var buf bytes.Buffer

		// More comprehensive clearing approach for Kitty protocol
		// Clear all images first
		clearAllCmd := "\x1b_Ga=d\x1b\\"
		buf.WriteString(clearAllCmd)

		// Clear all placements
		clearPlacementsCmd := "\x1b_Ga=d,p=1\x1b\\"
		buf.WriteString(clearPlacementsCmd)

		// Add a small delay command to ensure clearing is processed
		buf.WriteString("\x1b[0m")

		// Use the full available preview space for Kitty renderer
		cols := maxWidth
		rows := maxHeight

		// Ensure minimum bounds
		if cols < 1 {
			cols = 1
		}
		if rows < 1 {
			rows = 1
		}

		// For Kitty protocol, we can let the terminal handle the scaling
		// by setting the destination size to the full preview area
		slog.Info("Kitty rendering",
			"originalSize", fmt.Sprintf("%dx%d", img.Bounds().Dx(), img.Bounds().Dy()),
			"maxPreview", fmt.Sprintf("%dx%d", maxWidth, maxHeight),
			"terminalCells", fmt.Sprintf("%dx%d", cols, rows),
		)

		// Use a unique placement ID based on path hash to avoid conflicts
		placementId := uint32(42) // Default fallback
		if len(path) > 0 {
			hash := 0
			for _, c := range path {
				hash = hash*31 + int(c)
			}
			placementId = uint32(hash&0xFFFF) + 1000 // Ensure it's not 0 and avoid low numbers
		}

		opts := rasterm.KittyImgOpts{
			DstCols:     uint32(cols),
			DstRows:     uint32(rows),
			PlacementId: placementId,
		}

		if err := rasterm.KittyWriteImage(&buf, img, opts); err != nil {
			// If kitty fails, fall back to ANSI renderer
			slog.Error("Kitty renderer failed, falling back to ANSI", "error", err)
			bgColor, bgErr := hexToColor(defaultBGColor)
			if bgErr != nil {
				return "", fmt.Errorf("kitty failed and invalid background color: %w", bgErr)
			}
			// Resize image for ANSI fallback to fill the preview area
			resizedImg := imaging.Fill(img, maxWidth, maxHeight*2, imaging.Center, imaging.Lanczos)
			return ConvertImageToANSI(resizedImg, bgColor), nil
		}

		// Add explicit position reset at the end
		buf.WriteString("\x1b[H")

		// Return the captured kitty protocol string
		return buf.String(), nil
	case RendererANSI:
		fallthrough
	default:
		// Convert image to ANSI
		bgColor, err := hexToColor(defaultBGColor)
		if err != nil {
			return "", fmt.Errorf("invalid background color: %w", err)
		}

		// For ANSI rendering, resize image to fill the preview area
		// Use maxHeight*2 because each terminal row represents 2 pixel rows in ANSI rendering
		resizedImg := imaging.Fill(img, maxWidth, maxHeight*2, imaging.Center, imaging.Lanczos)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
