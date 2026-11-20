package filepreview

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/cache"
)

// Set SPF_BENCH_IMAGE to a real photo to benchmark it. Without it, the
// benchmarks use a synthetic 6000x4000 JPEG.
const (
	benchPaneCols = 60
	benchPaneRows = 30
	benchBGColor  = "#000000"
)

func benchImagePath(b *testing.B) string {
	b.Helper()
	if p := os.Getenv("SPF_BENCH_IMAGE"); p != "" {
		return p
	}
	const w, h = 6000, 4000
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(x ^ y), alphaOpaque})
		}
	}
	path := filepath.Join(b.TempDir(), "bench.jpg")
	f, err := os.Create(path)
	require.NoError(b, err)
	defer f.Close()
	require.NoError(b, jpeg.Encode(f, img, &jpeg.Options{Quality: 90}))
	return path
}

func newTestPreviewer() *ImagePreviewer {
	tc := NewTerminalCapabilities()
	tc.cellSizeInit.Do(func() {})
	return &ImagePreviewer{
		cache:       cache.New[string](1, time.Minute),
		terminalCap: tc,
	}
}

func BenchmarkImagePreview(b *testing.B) {
	path := benchImagePath(b)
	p := newTestPreviewer()

	b.Run("Kitty", func(b *testing.B) {
		var raw string
		for b.Loop() {
			var err error
			_, raw, err = p.ImagePreviewWithRenderer(path, benchPaneCols, benchPaneRows,
				benchBGColor, RendererKitty, 0)
			require.NoError(b, err)
		}
		b.ReportMetric(float64(len(raw)), "payload-bytes")
	})

	b.Run("ANSI", func(b *testing.B) {
		for b.Loop() {
			_, _, err := p.ImagePreviewWithRenderer(path, benchPaneCols, benchPaneRows,
				benchBGColor, RendererANSI, 0)
			require.NoError(b, err)
		}
	})
}

func BenchmarkImagePreviewStages(b *testing.B) {
	path := benchImagePath(b)
	data, err := os.ReadFile(path)
	require.NoError(b, err)

	b.Run("Decode", func(b *testing.B) {
		for b.Loop() {
			_, _, err := image.Decode(bytes.NewReader(data))
			require.NoError(b, err)
		}
	})

	b.Run("Prepare", func(b *testing.B) {
		for b.Loop() {
			_, err := prepareImageForPreview(data,
				benchPaneCols*DefaultPixelsPerColumn, benchPaneRows*DefaultPixelsPerRow)
			require.NoError(b, err)
		}
	})
}
