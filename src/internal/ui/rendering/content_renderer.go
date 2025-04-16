package rendering

import (
	"log/slog"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
)

// This is not zero safe. Please only use New function
type ContentRenderer struct {
	lines    []string
	maxLines int
	// Every line should have at most this many characters
	maxLineWidth    int
	sanitizeContent bool

	// Todo - Do we need alignStyle ? Not yet
	truncateStyle TruncateStyle
}

func (r *ContentRenderer) AddLines(lines ...string) {
	for _, line := range lines {
		r.AddLineWithCustomTruncate(line, r.truncateStyle)
	}
}

// AddLineWithCustomTruncate adds lines to the renderer, truncating each line according to the specified style.
// It does not trims whitespace, and its possible to add multiple empty lines using this.
func (r *ContentRenderer) AddLineWithCustomTruncate(lineStr string, truncateStyle TruncateStyle) {
	// If string is multiline, add individual lines separately

	// Dont use strings.Lines() we need to allow adding empty strings "" as line.
	for line := range strings.SplitSeq(lineStr, "\n") {
		// Todo : AnsiTree calculation should not be needed in non-debug mode.
		slog.Debug("Adding line", "cur size", len(r.lines), "line", AnsiTree(line))
		if len(r.lines) >= r.maxLines {
			slog.Error("Max lines reached", "maxLines", r.maxLines)
			return
		}
		// Todo : Move this tails to const
		// Todo : What if there is a "\t" will truncate take care of it ?
		line = TruncateBasedOnStyle(line, r.maxLineWidth, truncateStyle)
		if r.sanitizeContent {
			line = common.MakePrintableWithEscCheck(line, true)
		}
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
