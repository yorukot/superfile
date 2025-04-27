package rendering

import (
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/charmbracelet/x/exp/term/ansi"

	"github.com/charmbracelet/lipgloss"
)

type BorderConfig struct {
	// ANSI encoded strings are not allowed in border title and info items, for now.
	// The style is overridden with border's style.
	title string

	// Optional info items at the bottom of the border
	infoItems []string

	// Section dividers - A slice of values within [0,height-2]
	// Signifying usage of MiddleLeft and MiddleRight borders in Left and Right borders for
	// Section divider line.
	dividerIdx []int

	// Including corners. Both should be >= 2
	width  int
	height int

	titleLeftMargin int
}

func (b *BorderConfig) SetTitle(title string) {
	b.title = ansi.Strip(title)
}

func (b *BorderConfig) SetInfoItems(infoItems ...string) {
	for i := range infoItems {
		infoItems[i] = ansi.Strip(infoItems[i])
	}
	b.infoItems = infoItems
}

func (b *BorderConfig) AreInfoItemsTruncated() bool {
	reqWidth := 0
	for _, item := range b.infoItems {
		// border.MiddleLeft <content> border.MiddleRight border.Bottom
		reqWidth += 3 + runewidth.StringWidth(item)
	}

	return reqWidth > b.width-2
}

func (b *BorderConfig) AddDivider(idx int) {
	b.dividerIdx = append(b.dividerIdx, idx)
}

// border.Top with something that takes up more than 1 runewidth will not work, so
// we only allow 1 runewidth for now, in the config. multiple things like
// border corner characters must be single rune, or else it would break rendering.
// This is all filled in one function to prevent passing around too many values
// in helper functions
func (b *BorderConfig) GetBorder(borderStrings lipgloss.Border) lipgloss.Border {
	res := borderStrings

	// excluding corners
	actualWidth := b.width - 2
	actualHeight := b.height - 2

	// Min 5 width is needed for title so that at least one character can be
	// rendered
	if b.title != "" && actualWidth >= 5 {
		// We need to plain truncate the title if needed.
		// topWidth - 1( for BorderMiddleLeft) - 1 (for BorderMiddleRight) - 2 (padding)
		titleAvailWidth := actualWidth - 4

		// Basic Right truncation
		truncatedTitle := runewidth.Truncate(b.title, titleAvailWidth, "")
		remainingWidth := actualWidth - 4 - runewidth.StringWidth(truncatedTitle)

		margin := ""
		if remainingWidth > b.titleLeftMargin {
			margin = strings.Repeat(borderStrings.Top, b.titleLeftMargin)
			remainingWidth -= b.titleLeftMargin
		}

		// Title alignment is by default Left for now
		res.Top = margin + borderStrings.MiddleRight + " " + truncatedTitle + " " + borderStrings.MiddleLeft +
			strings.Repeat(borderStrings.Top, remainingWidth)
	}

	cnt := len(b.infoItems)
	// Minimum 4 character for each info item so that at least first character is rendered
	if cnt > 0 && actualWidth >= cnt*4 {
		// Max available width for each item's actual content
		// border.MiddleLeft <content> border.MiddleRight border.Bottom
		availWidth := actualWidth/cnt - 3
		infoText := ""
		for _, item := range b.infoItems {
			item = runewidth.Truncate(item, availWidth, "")
			infoText += borderStrings.MiddleRight + item + borderStrings.MiddleLeft + borderStrings.Bottom
		}

		// Fill the rest with border char.
		remainingWidth := actualWidth - runewidth.StringWidth(infoText)

		res.Bottom = strings.Repeat(borderStrings.Bottom, remainingWidth) + infoText
	}

	if len(b.dividerIdx) > 0 {
		// Update res.Left and res.Right
		leftBorder := strings.Builder{}
		rightBorder := strings.Builder{}
		di := 0
		for i := range actualHeight {
			if di < len(b.dividerIdx) && b.dividerIdx[di] == i {
				di++
				leftBorder.WriteString(borderStrings.MiddleLeft)
				rightBorder.WriteString(borderStrings.MiddleRight)
			} else {
				leftBorder.WriteString(borderStrings.Left)
				rightBorder.WriteString(borderStrings.Right)
			}
		}

		res.Left = leftBorder.String()
		res.Right = rightBorder.String()
	}

	return res
}

func NewBorderConfig(height int, width int) BorderConfig {
	return BorderConfig{
		height:          height,
		width:           width,
		titleLeftMargin: 1,
	}
}
