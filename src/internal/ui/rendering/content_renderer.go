package rendering

import (
	"log/slog"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
)

type ContentRenderer struct {
	lines []string

	// Allow at max this many lines. If there are lesser lines
	maxLines int
	// Every line should have at most this many characters
	maxLineWidth    int
	sanitizeContent bool

	// We can add alignStyle if needed
	truncateStyle TruncateStyle
}

func (r *ContentRenderer) CntLines() int {
	return len(r.lines)
}

func (r *ContentRenderer) AddLines(lines ...string) {
	for _, line := range lines {
		r.AddLineWithCustomTruncate(line, r.truncateStyle)
	}
}

func (r *ContentRenderer) ClearLines() {
	r.lines = r.lines[:0]
}

// Maybe better return an error ?
// AddLineWithCustomTruncate adds lines to the renderer, truncating each line according to the specified style.
// It does not trims whitespace, and its possible to add multiple empty lines using this.
func (r *ContentRenderer) AddLineWithCustomTruncate(lineStr string, truncateStyle TruncateStyle) {
	// If string is multiline, add individual lines separately
	// We dont use strings.Lines() we need to allow adding empty strings "" as line.
	for line := range strings.SplitSeq(lineStr, "\n") {
		if len(r.lines) >= r.maxLines {
			slog.Error("Max lines reached", "maxLines", r.maxLines)
			return
		}
		// Sanitazation should be done before truncate. Sanitization can increase width
		// For ex: Converting problematic unicode nbsp to spaces.
		if r.sanitizeContent {
			line = common.MakePrintableWithEscCheck(line, true)
		}
		// Some characters like "\t" are considered 1 width
		line = TruncateBasedOnStyle(line, r.maxLineWidth, truncateStyle)

		r.lines = append(r.lines, line)
	}
}

func (r *ContentRenderer) Render() string {
	return strings.Join(r.lines, "\n")
}

func NewContentRenderer(maxLines int, maxLineWidth int, truncateStyle TruncateStyle) ContentRenderer {
	return ContentRenderer{
		lines:           make([]string, 0),
		maxLines:        maxLines,
		maxLineWidth:    maxLineWidth,
		truncateStyle:   truncateStyle,
		sanitizeContent: true,
	}
}
