package filepreview

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func newTestImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	return img
}

func TestSelectRenderer(t *testing.T) {
	tests := []struct {
		name         string
		config       string
		kittyCapable bool
		sixelCapable bool
		want         ImageRenderer
	}{
		{"auto prefers kitty", "", true, true, RendererKitty},
		{"auto uses sixel without kitty", "auto", false, true, RendererSixel},
		{"auto falls back to ansi", "", false, false, RendererANSI},
		{"forced sixel", "sixel", true, false, RendererSixel},
		{"forced kitty", "kitty", false, true, RendererKitty},
		{"forced ansi", "ansi", true, true, RendererANSI},
		{"config is case insensitive", "SIXEL", false, false, RendererSixel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, selectRenderer(tt.config, tt.kittyCapable, tt.sixelCapable))
		})
	}
}

func TestIsSixelCapable(t *testing.T) {
	t.Setenv("WT_SESSION", "")
	assert.False(t, isSixelCapable())
	t.Setenv("WT_SESSION", "0f6a2a1e-0000-0000-0000-000000000000")
	assert.True(t, isSixelCapable())
}

func TestFitImageInCells(t *testing.T) {
	cell := TerminalCellSize{PixelsPerColumn: 10, PixelsPerRow: 20}
	tests := []struct {
		name                   string
		imgW, imgH, maxW, maxH int
		wantCols, wantRows     int
	}{
		{"wide image limited by width", 800, 400, 60, 30, 60, 15},
		{"tall image limited by height", 300, 700, 60, 30, 25, 30},
		{"tiny area never returns zero", 1000, 1, 5, 5, 5, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols, rows := fitImageInCells(tt.imgW, tt.imgH, tt.maxW, tt.maxH, cell)
			assert.Equal(t, tt.wantCols, cols, "cols")
			assert.Equal(t, tt.wantRows, rows, "rows")
		})
	}
}

func TestRenderWithSixel(t *testing.T) {
	t.Setenv("WT_SESSION", "test") // use Windows Terminal's 10x20 cell
	p := NewImagePreviewer()

	block, img, err := p.renderWithSixel(newTestImage(800, 400), 60, 30)
	require.NoError(t, err)

	assert.Equal(t, 60, img.Cols)
	assert.Equal(t, 15, img.Rows)
	assert.Equal(t, 0, img.OffsetX)
	assert.Equal(t, 7, img.OffsetY)
	assert.True(t, strings.HasPrefix(img.DCS, "\x1bP0;1q"), "DCS should start with the Sixel introducer")
	assert.True(t, strings.HasSuffix(img.DCS, "\x1b\\"), "DCS should end with ST")

	lines := strings.Split(block, "\n")
	require.Len(t, lines, 30)
	for i, line := range lines {
		assert.Equal(t, 60, ansi.StringWidth(line), "line %d width", i)
		markers := strings.Count(line, SixelMarker)
		if i >= img.OffsetY && i < img.OffsetY+img.Rows {
			assert.Equal(t, img.Cols, markers, "line %d should be marker cells", i)
		} else {
			assert.Zero(t, markers, "line %d should be blank", i)
		}
	}

	_, _, err = p.renderWithSixel(newTestImage(10, 10), 0, 10)
	assert.Error(t, err)
}

func TestSixelPlacement(t *testing.T) {
	img := &SixelImage{DCS: "<dcs>", OffsetX: 3, OffsetY: 7}
	// Content origin (70, 1) + offset (3, 7) = 0-based (73, 8) = 1-based (74, 9)
	assert.Equal(t, ansi.SaveCursor+"\x1b[9;74H<dcs>"+ansi.RestoreCursor, img.Placement(70, 1))
}

func TestImagePreviewSixelUsesCache(t *testing.T) {
	t.Setenv("WT_SESSION", "test")
	oldRenderer := common.Config.ImageRenderer
	common.Config.ImageRenderer = "sixel"
	t.Cleanup(func() { common.Config.ImageRenderer = oldRenderer })

	path := filepath.Join(t.TempDir(), "image.png")
	f, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, png.Encode(f, newTestImage(200, 100)))
	require.NoError(t, f.Close())

	p := NewImagePreviewer()
	first, err := p.ImagePreview(path, 40, 20, "#000000", 0)
	require.NoError(t, err)
	require.NotNil(t, first.Sixel)
	assert.Empty(t, first.RawTransmit)
	assert.Contains(t, first.Render, SixelMarker)

	second, err := p.ImagePreview(path, 40, 20, "#000000", 0)
	require.NoError(t, err)
	assert.Same(t, first.Sixel, second.Sixel, "second call should come from the cache")
	assert.Equal(t, first.Render, second.Render)

	assert.Empty(t, p.GetKittyClearRaw(), "no Kitty clear when Sixel is selected")
}
