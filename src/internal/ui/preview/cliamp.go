package preview

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/pkg/cliamp"
)

// CliampTickInterval is how often the cliamp preview refreshes while an audio
// file is highlighted (~12 fps), which keeps the spectrum animating.
const CliampTickInterval = 80 * time.Millisecond

// CliampTickMsg asks the model to re-render the cliamp audio preview.
type CliampTickMsg struct{}

// audioExtensions are the file types superfile hands to cliamp.
var audioExtensions = map[string]bool{
	".mp3":  true,
	".flac": true,
	".wav":  true,
	".ogg":  true,
	".oga":  true,
	".opus": true,
	".m4a":  true,
	".aac":  true,
	".alac": true,
	".wma":  true,
	".aif":  true,
	".aiff": true,
	".m4b":  true,
	".ape":  true,
	".wv":   true,
}

// IsAudioFile reports whether path has an extension cliamp can play.
func IsAudioFile(path string) bool {
	return audioExtensions[strings.ToLower(filepath.Ext(path))]
}

// ShouldAnimateCliamp reports whether the preview panel should keep ticking to
// refresh the spectrum for the given location.
func (m *Model) ShouldAnimateCliamp(location string) bool {
	return m.open && common.Config.CliampPreview && IsAudioFile(location)
}

// renderCliampPreview renders the now-playing panel with a live spectrum. State
// and spectrum are fetched synchronously here (this runs inside the preview
// render goroutine), so no shared state is needed.
func (m *Model) renderCliampPreview(r *rendering.Renderer, itemPath string, width, height int) string {
	client := cliamp.New()
	snapshot, stateErr := client.State()

	var bands []float64
	if stateErr == nil && snapshot.IsPlaying() {
		// Spectrum errors are non-fatal: fall back to a flat baseline.
		bands, _ = client.Spectrum()
	}

	r.AddLines(buildCliampLines(itemPath, snapshot, stateErr, bands, width, height)...)
	return r.Render()
}

func buildCliampLines(_ string, snapshot cliamp.Snapshot, stateErr error, bands []float64, width, height int) []string {
	if width <= 0 || height <= 0 {
		return nil
	}

	bg := common.FilePanelBGColor
	titleStyle := lipgloss.NewStyle().Bold(true).
		Foreground(common.FilePanelBorderActiveColor).Background(bg)
	fgStyle := lipgloss.NewStyle().Foreground(common.FilePanelFGColor).Background(bg)
	dimStyle := lipgloss.NewStyle().Foreground(common.FilePanelBorderColor).Background(bg)

	header := []string{titleStyle.Render(" cliamp")}
	if stateErr != nil {
		header = append(header,
			"",
			fgStyle.Render(" daemon not running"),
			dimStyle.Render(" Enter to start & play"),
		)
	} else {
		status := "idle"
		if snapshot.IsPlaying() {
			status = "playing"
		}
		header = append(header, dimStyle.Render(" "+status))
		if snapshot.Track.Title != "" {
			header = append(header, fgStyle.Render(" "+truncateVisible(snapshot.Track.Title, width-1)))
		}
		if meta := joinNonEmpty(" - ", snapshot.Track.Artist, snapshot.Track.Album); meta != "" {
			header = append(header, dimStyle.Render(" "+truncateVisible(meta, width-1)))
		}
		if snapshot.IsPlaying() && snapshot.Duration > 0 {
			header = append(header, dimStyle.Render(
				" "+formatClock(snapshot.Position)+" / "+formatClock(snapshot.Duration)))
		}
	}

	hint := dimStyle.Render(" Enter: play in cliamp")

	// Degrade gracefully when there is no room for anything but text.
	if height <= len(header) {
		return fitLines(header, height)
	}

	budget := height - len(header)
	hintLines := []string{}
	if budget >= 2 {
		hintLines = append(hintLines, hint)
		budget--
	}

	spectrum := renderSpectrumRows(bands, width, budget)
	for len(spectrum) < budget {
		spectrum = append(spectrum, strings.Repeat(" ", width))
	}

	out := make([]string, 0, height)
	out = append(out, header...)
	out = append(out, spectrum...)
	out = append(out, hintLines...)
	return fitLines(out, height)
}

// renderSpectrumRows draws vertical bars from the given bands (each 0..1) into
// exactly height rows of exactly width cells. Bars grow from the bottom and use
// a classic green/yellow/red gradient. The rows carry their own background so
// they blend with the preview panel.
func renderSpectrumRows(bands []float64, width, height int) []string {
	rows := make([]string, 0, height)
	if width <= 0 || height <= 0 {
		return rows
	}

	n := len(bands)
	bandWidth := 1
	if n > 0 {
		bandWidth = width / n
		if bandWidth < 1 {
			bandWidth = 1
		}
	}

	for row := 0; row < height; row++ {
		var b strings.Builder
		for col := 0; col < width; col++ {
			var value float64
			if n > 0 {
				band := col / bandWidth
				if band > n-1 {
					band = n - 1
				}
				value = bands[band]
			}
			b.WriteRune(levelRune(value, row, height))
		}
		rows = append(rows, lipgloss.NewStyle().
			Foreground(spectrumRowColor(row, height)).
			Background(common.FilePanelBGColor).
			Render(b.String()))
	}
	return rows
}

func levelRune(value float64, row, height int) rune {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	filled := value * float64(height)
	// How much of this cell is filled, measured from the bottom edge.
	fill := filled - float64(height-row-1)
	switch {
	case fill <= 0:
		return ' '
	case fill >= 1:
		return '█'
	}
	idx := int(fill*8 + 0.5)
	if idx < 1 {
		idx = 1
	}
	if idx > 7 {
		idx = 7
	}
	return []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇'}[idx-1]
}

func spectrumRowColor(row, height int) color.Color {
	if height <= 1 {
		return lipgloss.Color("#5fff87")
	}
	fromTop := float64(row+1) / float64(height)
	switch {
	case fromTop <= 0.2:
		return lipgloss.Color("#ff5f5f")
	case fromTop <= 0.5:
		return lipgloss.Color("#ffd75f")
	default:
		return lipgloss.Color("#5fff87")
	}
}

func formatClock(seconds float64) string {
	if seconds < 0 {
		seconds = 0
	}
	total := int(seconds)
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}

func joinNonEmpty(sep string, parts ...string) string {
	kept := parts[:0]
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, sep)
}

func truncateVisible(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	var b strings.Builder
	used := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if used+rw > width {
			break
		}
		b.WriteRune(r)
		used += rw
	}
	return b.String()
}

func fitLines(lines []string, height int) []string {
	if height <= 0 {
		return nil
	}
	if len(lines) > height {
		return lines[:height]
	}
	return lines
}
