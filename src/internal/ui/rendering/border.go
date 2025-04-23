package rendering

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/mattn/go-runewidth"
)

type BorderConfig struct {

	// ANSI encoded strings are not allowed in border title and info items, for now.
	// The style is overidden with border's style.
	// Todo : Allow it. Would need a ansiTruncateLeft function for this.
	// That can trucate strings towards Left, while preserving ansi escape sequences
	// Optional title at the top of the border
	title string

	// Optional info items at the bottom of the border
	infoItems []string

	borderStrings lipgloss.Border

	// Including corners. Both should be >= 2
	width  int
	height int

	titleLeftMargin int
}

func (b *BorderConfig) SetTitle(title string) {
	b.title = ansi.Strip(title)
}

// Todo - take varidiac args
func (b *BorderConfig) SetInfoItems(infoItems []string) {
	for i := range infoItems {
		infoItems[i] = ansi.Strip(infoItems[i])
	}
	b.infoItems = infoItems
}

// Todo - unit test with border.Top with something that takes up more than 1 runewidth
// Sadly that might now work, so maybe only allow 1 runewidth for now, in the config ?
// multiple things like corner characters must be single rune, or else it would break things.
// Todo - Write thorough unit tests that have bigger title which needs to be truncated.
func (b *BorderConfig) GetBorder() lipgloss.Border {
	res := b.borderStrings

	// width excluding corners
	actualWidth := b.width - 2

	// Min 5 width is needed for title so that at least one character can be
	// rendered
	if b.title != "" && actualWidth >= 5 {
		// We need to plain truncate the title if needed.
		// topWidth - 1( for BorderMiddleLeft) - 1 (for BorderMiddleRight) - 2 (padding)
		titleAvailWidth := actualWidth - 4

		// This is okay, because we are not yet allowing ansi escaped text
		// Basic Left truncate without preserving ansi
		truncatedTitle := runewidth.Truncate(b.title, titleAvailWidth, "")
		remainingWidth := actualWidth - 4 - runewidth.StringWidth(truncatedTitle)

		margin := ""
		if remainingWidth > b.titleLeftMargin {
			margin = strings.Repeat(b.borderStrings.Top, b.titleLeftMargin)
			remainingWidth -= b.titleLeftMargin
		}

		// Title alignment is by default Left for now
		res.Top = margin + b.borderStrings.MiddleRight + " " + truncatedTitle + " " + b.borderStrings.MiddleLeft +
			strings.Repeat(b.borderStrings.Top, remainingWidth)
	}

	cnt := len(b.infoItems)
	// Minimum 3 character for each info item
	// We can make it 4 if we want a padding of 1 border.Bottom character
	// after each item - Todo - Do it.
	if cnt > 0 && actualWidth >= cnt*3 {
		// Todo : Do this. What if maxCnt > cnt ?
		// maxCnt := actualWidth / 4
		// infoItems := b.infoItems[:maxCnt]

		// Right aligned // Individually Truncated

		// Max available width for each item's actual content
		availWidth := actualWidth/cnt - 2
		infoText := ""
		for _, item := range b.infoItems {
			item = runewidth.Truncate(item, availWidth, "")
			infoText += b.borderStrings.MiddleRight + item + b.borderStrings.MiddleLeft
		}

		// Fill the rest with border char.
		remainingWidth := actualWidth - runewidth.StringWidth(infoText)

		res.Bottom = strings.Repeat(b.borderStrings.Bottom, remainingWidth) + infoText

		slog.Debug("Border rendering", "bottom len", len(res.Bottom),
			"actualWidth", actualWidth, "infoText Len", len(infoText),
			"bottom", res.Bottom, "bottom bytes", fmt.Sprintf("%v", []byte(res.Bottom)))
	}
	return res
}

func NewBorderConfig(height int, width int, borderStrings lipgloss.Border) BorderConfig {
	return BorderConfig{
		borderStrings:   borderStrings,
		height:          height,
		width:           width,
		titleLeftMargin: 1,
	}
}


