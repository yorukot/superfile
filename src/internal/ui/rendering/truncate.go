package rendering

import (
	"log/slog"

	"github.com/charmbracelet/x/exp/term/ansi"
)

type TruncateStyle int

// These truncate styles must preserve ansi escape codes. If something doesn't preserves
// it shouldn't be here
const (
	PlainTruncateRight = iota
	TailsTruncateRight
)

func TruncateBasedOnStyle(line string, maxWidth int, truncateStyle TruncateStyle) string {
	// Todo : Do we need ansi.Width check ?

	switch truncateStyle {
	case PlainTruncateRight:
		return ansi.Truncate(line, maxWidth, "")
	case TailsTruncateRight:
		return ansi.Truncate(line, maxWidth, "...")
	default:
		slog.Error("Invalid truncate style", "style", truncateStyle)
		return ""
	}
}
