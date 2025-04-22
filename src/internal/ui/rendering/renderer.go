package rendering

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/superfile/src/internal/common"
)

type Renderer struct {
	content ContentRenderer

	// Border should not be coupled with content at all. A separate struct is better
	// If you want a borderless content, you will have all these extra variables
	// Renderer's interaction with Border isn't related to content , so usage as a separate
	// struct is okay ?
	// No... Border itself is not an independent renderer. You need width and height
	// If you want to save memory, use pointers
	border BorderConfig

	// Should this go in contentRenderer - No . ContentRenderer is not for storing style configs
	contentFGColor lipgloss.Color
	contentBGColor lipgloss.Color

	// Should this go in borderConfig ?
	borderFGColor lipgloss.Color
	borderBGColor lipgloss.Color

	// Final rendered string should have exactly this many lines, including borders
	totalHeight int
	// Every line should have at most this many characters, including borders
	totalWidth     int
	borderRequired bool
}

// Add lines as much as the remaining capacity allows
func (r *Renderer) AddLines(lines ...string) {
	r.content.AddLines(lines...)
}

// Truncate would always preserve ansi codes.
func (r *Renderer) AddLineWithCustomTruncate(line string, truncateStyle TruncateStyle) {
	r.content.AddLineWithCustomTruncate(line, truncateStyle)
}

func (r *Renderer) SetBorderTitle(title string) {
	r.border.SetTitle(title)
}

func (r *Renderer) SetBorderInfoItems(infoItems []string) {
	r.border.SetInfoItems(infoItems)
}

func (r *Renderer) Render() string {
	res := r.Style().Render(r.content.Render())

	maxW := 0
	for line := range strings.Lines(res) {
		maxW = max(maxW, ansi.StringWidth(line))
	}
	lineCnt := strings.Count(res, "\n") + 1
	slog.Debug("Rendered output", "line count", lineCnt, "totalHeight", r.totalHeight,
		"totalWidth", r.totalWidth, "maxW", maxW)

	return res
}

func (r *Renderer) Style() lipgloss.Style {
	s := lipgloss.NewStyle().
		Width(r.content.maxLineWidth).
		Height(r.content.maxLines)

	if r.borderRequired {
		s = s.Border(r.border.GetBorder()).
			BorderForeground(r.borderFGColor).
			BorderBackground(r.borderBGColor)
	}
	s = s.Background(r.contentBGColor).
		Foreground(r.contentFGColor)
	return s
}

func NewRenderer(totalHeight int, totalWidth int, borderRequired bool, truncateStyle TruncateStyle,
	contentFGColor lipgloss.Color, contentBGColor lipgloss.Color,
	borderFGColor lipgloss.Color, borderBGColor lipgloss.Color) Renderer {
	contentLines := totalHeight
	if borderRequired {
		contentLines -= 2
	}
	contentWidth := totalWidth
	if borderRequired {
		contentWidth -= 2
	}
	return Renderer{
		content: NewContentRenderer(contentLines, contentWidth, truncateStyle),
		border:  NewBorderConfig(totalWidth, totalHeight),

		contentFGColor: contentFGColor,
		contentBGColor: contentBGColor,
		borderFGColor:  borderFGColor,
		borderBGColor:  borderBGColor,

		totalWidth:     totalWidth,
		totalHeight:    totalHeight,
		borderRequired: borderRequired,
	}
}

// Todo : rendering package should not be aware of sidebar
func SidebarRenderer(totalHeight int, totalWidth int, sidebarFocussed bool) Renderer {
	borderFG := common.SidebarBorderColor
	if sidebarFocussed {
		borderFG = common.SidebarBorderActiveColor
	}
	return NewRenderer(totalHeight, totalWidth, true, PlainTruncateRight,
		common.SidebarFGColor, common.SidebarBGColor, borderFG, common.SidebarBGColor)
}

// Todo : Move to diff package
func ProcessBarRenderer(totalHeight int, totalWidth int, processBarFocussed bool) Renderer {
	borderFGColor := common.FooterBorderColor
	if processBarFocussed {
		borderFGColor = common.FooterBorderActiveColor
	}

	r := NewRenderer(totalHeight, totalWidth, true, PlainTruncateRight,
		common.FooterFGColor, common.FooterBGColor, borderFGColor, common.FooterBGColor)
	r.SetBorderTitle("Process")

	return r
}
