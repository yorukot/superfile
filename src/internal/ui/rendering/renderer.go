package rendering

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

type StyleModifier func(lipgloss.Style) lipgloss.Style

// For now we are not allowing to add/update/remove lines to previous sections
// We may allow that later.
// Also we could have functions about getting sections count, line count, adding updating a
// specific line in a specific section, and adjusting section sizes. But not needed now.
// NOTE: Renderer's zero value isn't safe to use, always use NewRenderer()
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

	// Use this to add additional style modifications
	// This is applied before any style update that are defined by other configurations,
	// like border, height, width. Hence if conflicting styles are used, they can get
	// overridden
	styleModifiers []StyleModifier

	// Maybe better rename these to maxHeight
	// Final rendered string should have exactly this many lines, including borders
	// But if truncateHeight is true, it maybe be <= totalHeight
	totalHeight int
	// Every line should have at most this many characters, including borders
	totalWidth int

	contentHeight int
	contentWidth  int

	// Note: Must pass non empty borderStrings if borderRequired is set as true
	// TODO: Have ansi.StringWidth checks in `ValidateConfig`
	// If you silently pass empty border, rendering will be unexpectd and,
	// it might take some time to RCA.
	borderRequired bool
	borderStrings  lipgloss.Border
	// for logging
	name string
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

	Border       lipgloss.Border
	RendererName string
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
		//nolint: gosec // Not for security purpose, only for logging
		RendererName: "R-" + strconv.Itoa(rand.IntN(rendererNameMax)),
	}
}

func NewRenderer(cfg RendererConfig) (*Renderer, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}
	return createRendererWithValidatedConfig(cfg), nil
}

func NewRendererWithAutoFixConfig(cfg RendererConfig) *Renderer {
	validateAndAutoFix(&cfg)
	return createRendererWithValidatedConfig(cfg)
}

func createRendererWithValidatedConfig(cfg RendererConfig) *Renderer {
	contentHeight := cfg.TotalHeight
	if cfg.BorderRequired {
		contentHeight -= 2
	}
	contentWidth := cfg.TotalWidth
	if cfg.BorderRequired {
		contentWidth -= 2
	}

	return &Renderer{

		contentSections: []ContentRenderer{
			NewContentRenderer(contentHeight, contentWidth, cfg.DefTruncateStyle, cfg.RendererName),
		},
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
		name:           cfg.RendererName,
	}
}

// There is code duplication with `validate` but, I can't think of any clean design pattern to fix that.
// Note: Having a function validate(cfg,autoFix) error and ensure err is not nil via panic is not clean.
func validateAndAutoFix(cfg *RendererConfig) {
	if cfg.TotalHeight < 0 || cfg.TotalWidth < 0 {
		slog.Debug("AutoFixConfig: clamping negative dimensions", "h", cfg.TotalHeight, "w", cfg.TotalWidth)
		cfg.TotalHeight = max(0, cfg.TotalHeight)
		cfg.TotalWidth = max(0, cfg.TotalWidth)
	}
	if cfg.BorderRequired {
		if cfg.TotalWidth < MinWidthForBorder || cfg.TotalHeight < MinHeightForBorder {
			slog.Debug("AutoFixConfig: disabling border due to insufficient dimensions",
				"h", cfg.TotalHeight, "w", cfg.TotalWidth)
			cfg.BorderRequired = false
		}
	}
}

func validate(cfg RendererConfig) error {
	if cfg.TotalHeight < 0 || cfg.TotalWidth < 0 {
		return fmt.Errorf("dimensions must be non-negative (h=%d, w=%d)", cfg.TotalHeight, cfg.TotalWidth)
	}
	if cfg.BorderRequired {
		if cfg.TotalWidth < MinWidthForBorder || cfg.TotalHeight < MinHeightForBorder {
			return errors.New("need at least 2 width and height for borders")
		}
	}
	return nil
}
