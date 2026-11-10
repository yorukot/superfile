package filepreview

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi/kitty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeOrientedJPEG writes a 600x400 JPEG, gray with a red top-left quarter,
// tagged with the given EXIF orientation.
func writeOrientedJPEG(t *testing.T, orientation uint16) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 600, 400))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Gray{Y: 128}), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(0, 0, 300, 200), image.NewUniform(color.RGBA{R: 255, A: 255}),
		image.Point{}, draw.Src)

	var encoded bytes.Buffer
	require.NoError(t, jpeg.Encode(&encoded, img, &jpeg.Options{Quality: 95}))

	// A big-endian TIFF header with one IFD entry: Orientation, SHORT, count 1.
	tiff := []byte{'M', 'M', 0, 42, 0, 0, 0, 8, 0, 1, 0x01, 0x12, 0, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint16(tiff[18:], orientation)
	payload := append([]byte("Exif\x00\x00"), tiff...)
	app1 := []byte{0xFF, 0xE1, 0, 0}
	binary.BigEndian.PutUint16(app1[2:], uint16(len(payload)+2))

	jpg := encoded.Bytes()
	out := append(append(append([]byte{}, jpg[:2]...), append(app1, payload...)...), jpg[2:]...)
	path := filepath.Join(t.TempDir(), "oriented.jpg")
	require.NoError(t, os.WriteFile(path, out, 0o600))
	return path
}

func isRed(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r>>8 > 200 && g>>8 < 60 && b>>8 < 60
}

func TestPrepareImageForPreviewOrientation(t *testing.T) {
	tests := []struct {
		name        string
		orientation uint16
		wantSize    image.Point
		redCorner   func(image.Rectangle) image.Point
	}{
		{"upright", 1, image.Pt(150, 100), func(r image.Rectangle) image.Point { return r.Min }},
		{"rotate 90 clockwise", 6, image.Pt(100, 150),
			func(r image.Rectangle) image.Point { return image.Pt(r.Max.X-1, r.Min.Y) }},
		{"rotate 90 counterclockwise", 8, image.Pt(100, 150),
			func(r image.Rectangle) image.Point { return image.Pt(r.Min.X, r.Max.Y-1) }},
		{"rotate 180", 3, image.Pt(150, 100),
			func(r image.Rectangle) image.Point { return image.Pt(r.Max.X-1, r.Max.Y-1) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(writeOrientedJPEG(t, tt.orientation))
			require.NoError(t, err)

			img, err := prepareImageForPreview(data, 150, 150)
			require.NoError(t, err)

			assert.Equal(t, tt.wantSize, img.Bounds().Size())
			assert.True(t, isRed(img.At(tt.redCorner(img.Bounds()).X, tt.redCorner(img.Bounds()).Y)),
				"red quarter should land in the displayed corner")
		})
	}
}

func TestKittyPlacementFollowsOrientation(t *testing.T) {
	p := newTestPreviewer()
	path := writeOrientedJPEG(t, 6)

	render, raw, err := p.ImagePreviewWithRenderer(path, 60, 30, "#000000", RendererKitty, 0)
	require.NoError(t, err)
	require.NotEmpty(t, raw)

	// 60x30 cells of 10x20 pixels make a 600x600 box. The displayed image is
	// 400x600, so it fills the 30 rows and takes 40 columns.
	rows := strings.Split(render, "\n")
	assert.Len(t, rows, 30)
	assert.Equal(t, 40, strings.Count(rows[0], string(kitty.Placeholder)))
}
