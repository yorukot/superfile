package rendering

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
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

func NewRenderer(cfg RendererConfig) (*Renderer, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}

	contentHeight := cfg.TotalHeight
	if cfg.BorderRequired {
		contentHeight -= 2
	}
	contentWidth := cfg.TotalWidth
	if cfg.BorderRequired {
		contentWidth -= 2
	}

	return &Renderer{

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
	}, nil
}

func validate(cfg RendererConfig) error {
	if cfg.TotalHeight < 0 || cfg.TotalWidth < 0 {
		return fmt.Errorf("dimensions must be non-negative (h=%d, w=%d)", cfg.TotalHeight, cfg.TotalWidth)
	}
	if cfg.BorderRequired {
		if cfg.TotalWidth < 2 || cfg.TotalHeight < 2 {
			return errors.New("need at least 2 width and height for borders")
		}
	}
	return nil
}
