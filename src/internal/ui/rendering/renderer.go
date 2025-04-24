package rendering

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
)

// For now we are not allowing to add/update/remove lines to previous sections
// We may allow that later.
// Also we could have functions about getting sections count, line count, adding updating a
// specific line in a specific section, and adjusting section sizes. But not needed now.
type Renderer struct {

	// Current sectionization will not allow to predefine section
	// but only allow adding them via AddSection(). Hence trucateWill be applicable to
	// last section only.
	contentSections []ContentRenderer

	// Empty for last section . len(sectionDividers) should be equal to len(contentSections) - 1
	sectionDividers []string
	curSectionIdx   int
	// Including Dividers - Count of actual lines that were added. It maybe <= totalHeight - 2
	actualContentHeight int
	defTruncateStyle    TruncateStyle

	// Whether to reduce rendered height to fit number of lines
	truncateHeight bool

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

	// Maybe better rename these to maxHeight
	// Final rendered string should have exactly this many lines, including borders
	// But if truncateHeight is true, it maybe be <= totalHeight
	totalHeight int
	// Every line should have at most this many characters, including borders
	totalWidth int

	contentHeight int
	contentWidth  int

	borderRequired bool

	borderStrings lipgloss.Border

	// Dont add any colors.
	// Todo : Is it needed, Is using .Foreground(lipgloss.NoColor{}) equivalent to not Using .Foreground()
	noColor bool
}

// Add lines as much as the remaining capacity allows
func (r *Renderer) AddLines(lines ...string) {
	r.contentSections[r.curSectionIdx].AddLines(lines...)
}

// Lines until now will belong to current section, and
// Any new lines will belong to a new section
func (r *Renderer) AddSection() {
	// r.actualContentHeight before this point only includes sections
	// before r.curSectionIdx
	r.actualContentHeight += r.contentSections[r.curSectionIdx].CntLines()

	// Silently Fail if cannot add
	if r.contentHeight <= r.actualContentHeight {
		slog.Error("Cannot add any more sections", "actualHeight", r.actualContentHeight,
			"contentHeight", r.contentHeight)
		return
	}

	// Add divider
	r.border.AddDivider(r.actualContentHeight)
	r.sectionDividers = append(r.sectionDividers, strings.Repeat(r.borderStrings.Top, r.contentWidth))
	r.actualContentHeight++

	remainingHeight := r.contentHeight - r.actualContentHeight
	r.contentSections = append(r.contentSections,
		NewContentRenderer(remainingHeight, r.contentWidth, r.defTruncateStyle))
	// Adjust index
	r.curSectionIdx++
}

// Truncate would always preserve ansi codes.
func (r *Renderer) AddLineWithCustomTruncate(line string, truncateStyle TruncateStyle) {
	r.contentSections[r.curSectionIdx].AddLineWithCustomTruncate(line, truncateStyle)
}

func (r *Renderer) SetBorderTitle(title string) {
	r.border.SetTitle(title)
}

func (r *Renderer) SetBorderInfoItems(infoItems []string) {
	r.border.SetInfoItems(infoItems)
}

// Should not do any updates on 'r'
func (r *Renderer) Render() string {
	// Todo : Do a validate before performing render
	// Check that 	contentSections []ContentRenderer
	// len(sectionDividers) should be equal to len(contentSections) - 1
	// curSectionIdx is okay
	// actualContentHeight is <= contentHeight
	content := strings.Builder{}
	for i := range r.contentSections {
		// After every iteration, current cursor will be on next newline
		curContent := r.contentSections[i].Render()
		content.WriteString(curContent)
		// == "" check cant differentiate between no data, vs empty line
		if r.contentSections[i].CntLines() > 0 {
			content.WriteString("\n")
		}

		if i < len(r.contentSections)-1 {
			// True for all except last section
			content.WriteString(r.sectionDividers[i])
			content.WriteString("\n")
		}
	}
	contentStr := strings.TrimSuffix(content.String(), "\n")

	slog.Debug("contentStr is", "contentStr", contentStr)

	res := r.Style().Render(contentStr)

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
	contentHeight := r.contentHeight
	if r.truncateHeight {
		contentHeight = r.actualContentHeight
	}
	s := lipgloss.NewStyle().
		Width(r.contentWidth).
		Height(contentHeight)

	if r.borderRequired {
		s = s.Border(r.border.GetBorder(r.borderStrings))

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

type RendererConfig struct {
	TotalHeight int
	TotalWidth  int

	DefTruncateStyle TruncateStyle
	TruncateHeight   bool
	BorderRequired   bool

	ContentFGColor lipgloss.TerminalColor
	ContentBGColor lipgloss.TerminalColor

	BorderFGColor lipgloss.TerminalColor
	BorderBGColor lipgloss.TerminalColor

	Border lipgloss.Border
}

func DefaultRendererConfig(totalHeight int, totalWidth int) RendererConfig {
	return RendererConfig{
		TotalHeight:      totalHeight,
		TotalWidth:       totalWidth,
		TruncateHeight:   false,
		BorderRequired:   false,
		DefTruncateStyle: PlainTruncateRight,
		ContentFGColor:   lipgloss.NoColor{},
		ContentBGColor:   lipgloss.NoColor{},
		BorderFGColor:    lipgloss.NoColor{},
		BorderBGColor:    lipgloss.NoColor{},
	}
}

func NewRenderer(cfg RendererConfig) Renderer {
	// Validations of config
	cfg, err := ValidateAndFix(cfg)
	if err != nil {
		// Config cannot be fixed. Too bad
		panic(fmt.Sprintf("Invalid renderer config : %v", err))
	}

	contentHeight := cfg.TotalHeight
	if cfg.BorderRequired {
		contentHeight -= 2
	}
	contentWidth := cfg.TotalWidth
	if cfg.BorderRequired {
		contentWidth -= 2
	}

	return Renderer{

		contentSections:     []ContentRenderer{NewContentRenderer(contentHeight, contentWidth, cfg.DefTruncateStyle)},
		sectionDividers:     nil,
		curSectionIdx:       0,
		actualContentHeight: 0,
		defTruncateStyle:    cfg.DefTruncateStyle,
		truncateHeight:      cfg.TruncateHeight,

		border: NewBorderConfig(cfg.TotalHeight, cfg.TotalWidth),

		contentFGColor: cfg.ContentFGColor,
		contentBGColor: cfg.ContentBGColor,
		borderFGColor:  cfg.BorderFGColor,
		borderBGColor:  cfg.BorderBGColor,

		totalHeight:   cfg.TotalHeight,
		totalWidth:    cfg.TotalWidth,
		contentHeight: contentHeight,
		contentWidth:  contentWidth,

		borderRequired: cfg.BorderRequired,
		borderStrings:  cfg.Border,
	}
}

// Log any fix that is needed
// Todo : What is better ? This or pass by pointer ?
// Does passing by pointer means object has to be moved to heap ?
func ValidateAndFix(cfg RendererConfig) (RendererConfig, error) {
	// Todo : Validations
	// 1 - Width and Height should be >=2 if border is required
	// 2 - Border should have single runewidth strings
	if cfg.BorderRequired {
		if cfg.TotalWidth < 2 || cfg.TotalHeight < 2 {
			cfg.TotalHeight = 0
			cfg.TotalWidth = 0
			cfg.BorderRequired = false
		}
	}

	return cfg, nil
}
