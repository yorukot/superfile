package preview

import (
	"errors"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/pkg/cliamp"
)

func TestIsAudioFile(t *testing.T) {
	assert.True(t, IsAudioFile("/music/song.MP3"))
	assert.True(t, IsAudioFile("track.flac"))
	assert.True(t, IsAudioFile("a/b/c.ogg"))
	assert.False(t, IsAudioFile("video.mp4"))
	assert.False(t, IsAudioFile("notes.txt"))
	assert.False(t, IsAudioFile("noextension"))
}

func TestLevelRune(t *testing.T) {
	const height = 4
	// Empty band produces an empty column.
	for row := 0; row < height; row++ {
		assert.Equal(t, ' ', levelRune(0, row, height))
	}
	// Full band fills every cell.
	for row := 0; row < height; row++ {
		assert.Equal(t, '█', levelRune(1, row, height))
	}
	// Half band fills the bottom two of four rows.
	assert.Equal(t, '█', levelRune(0.5, 3, height))
	assert.Equal(t, '█', levelRune(0.5, 2, height))
	assert.Equal(t, ' ', levelRune(0.5, 1, height))
	assert.Equal(t, ' ', levelRune(0.5, 0, height))
	// Out of range values are clamped.
	assert.Equal(t, '█', levelRune(2, 0, height))
	assert.Equal(t, ' ', levelRune(-1, 0, height))
}

func TestRenderSpectrumRowsDimensions(t *testing.T) {
	bands := []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 1.0}
	rows := renderSpectrumRows(bands, 40, 8)
	assert.Len(t, rows, 8)
	for _, row := range rows {
		assert.Equal(t, 40, lipgloss.Width(row))
	}

	// No bands still yields full-size blank rows.
	rows = renderSpectrumRows(nil, 12, 3)
	assert.Len(t, rows, 3)
	for _, row := range rows {
		assert.Equal(t, 12, lipgloss.Width(row))
		assert.Equal(t, "            ", ansi.Strip(row))
	}
}

func TestBuildCliampLinesDaemonDown(t *testing.T) {
	lines := buildCliampLines("song.mp3", cliamp.Snapshot{}, errors.New("down"), nil, 40, 12)
	joined := ""
	for _, l := range lines {
		joined += ansi.Strip(l) + "\n"
	}
	assert.Contains(t, joined, "cliamp")
	assert.Contains(t, joined, "daemon not running")
	assert.Contains(t, joined, "Enter to start & play")
	assert.Len(t, lines, 12)
}

func TestFormatClock(t *testing.T) {
	assert.Equal(t, "0:00", formatClock(0))
	assert.Equal(t, "0:05", formatClock(5))
	assert.Equal(t, "1:05", formatClock(65))
	assert.Equal(t, "0:00", formatClock(-3))
}

func TestTruncateVisible(t *testing.T) {
	assert.Equal(t, "abc", truncateVisible("abc", 5))
	assert.Equal(t, "ab", truncateVisible("abcdef", 2))
	assert.Equal(t, "", truncateVisible("abc", 0))
	assert.Equal(t, "日本", truncateVisible("日本語", 4))
}
