package rendering

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
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
	contentFGColor lipgloss.TerminalColor
	contentBGColor lipgloss.TerminalColor

	// Should this go in borderConfig ?
	borderFGColor lipgloss.TerminalColor
	borderBGColor lipgloss.TerminalColor

	// Final rendered string should have exactly this many lines, including borders
	totalHeight int
	// Every line should have at most this many characters, including borders
	totalWidth     int
	borderRequired bool

	// Dont add any colors.
	// Todo : Is it needed, Is using .Foreground(lipgloss.NoColor{}) equivalent to not Using .Foreground()
	noColor bool
}

// Add lines as much as the remaining capacity allows
func (r *Renderer) AddLines(lines ...string) {
	r.content.AddLines(lines...)
}

// Lines until now will belong to current section, and
// Any new lines will belong to a new section
func (r *Renderer) AddSection() {

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
		s = s.Border(r.border.GetBorder())

		if !r.noColor {
			s = s.BorderForeground(r.borderFGColor).
				BorderBackground(r.borderBGColor)
		}
	}
	if !r.noColor {
		s = s.Background(r.contentBGColor).
			Foreground(r.contentFGColor)
	}
	return s
}

type RendererConfig struct{
	TotalHeight int 
	TotalWidth int 
	
	DefTruncateStyle TruncateStyle

	BorderRequired bool

	ContentFGColor lipgloss.TerminalColor
	ContentBGColor lipgloss.TerminalColor

	BorderFGColor lipgloss.TerminalColor
	BorderBGColor lipgloss.TerminalColor

	Border lipgloss.Border
}

func DefaultRendererConfig(totalHeight int, totalWidth int) RendererConfig {
	return RendererConfig{
		TotalHeight: totalHeight,
		TotalWidth:  totalWidth,
		BorderRequired: false,
		DefTruncateStyle: PlainTruncateRight,
		ContentFGColor: lipgloss.NoColor{},
		ContentBGColor: lipgloss.NoColor{},
		BorderFGColor: lipgloss.NoColor{},
		BorderBGColor: lipgloss.NoColor{},
	}
}

func NewRenderer(cfg RendererConfig) Renderer {
	
	// Validations of config
	cfg, err := ValidateAndFix(cfg)
	if err != nil {
		// Config cannot be fixed. Too bad
		panic(fmt.Sprintf("Invalid renderer config : %v", err))
	}

	contentLines := cfg.TotalHeight
	if cfg.BorderRequired {
		contentLines -= 2
	}
	contentWidth := cfg.TotalWidth
	if cfg.BorderRequired {
		contentWidth -= 2
	}

	return Renderer{
		content: NewContentRenderer(contentLines, contentWidth, cfg.DefTruncateStyle),
		border:  NewBorderConfig(cfg.TotalHeight, cfg.TotalWidth, cfg.Border),
		
		contentFGColor: cfg.ContentFGColor,
		contentBGColor: cfg.ContentBGColor,
		borderFGColor: cfg.BorderFGColor,
		borderBGColor: cfg.BorderBGColor,

		totalHeight: cfg.TotalHeight,
		totalWidth: cfg.TotalWidth,

		borderRequired: cfg.BorderRequired,
	}
}

// Log any fix that is needed
// Todo : What is better ? This or pass by pointer ?
// Does passing by pointer means object has to be moved to heap ?
func ValidateAndFix(cfg RendererConfig) (RendererConfig, error) {
	// Todo : Validations
	// 1 - Width and Height should be >=2 if border is required
	// 2 - Border should have single runewidth strings 


	return cfg, nil
}

