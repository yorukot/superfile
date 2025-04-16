package rendering

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/internal/common"
)

type BorderConfig struct {
	// Optional title at the top of the border
	title string

	// Optional info items at the bottom of the border
	infoItems []string

	lipglossBorder lipgloss.Border
}

func (b *BorderConfig) SetTitle(title string) {
	b.title = title
}

func (b *BorderConfig) SetInfoItems(infoItems []string) {
	b.infoItems = infoItems
}

func NewBorderConfig() BorderConfig {
	return BorderConfig{
		lipglossBorder: DefaultLipglossBorder(),
	}
}

func DefaultLipglossBorder() lipgloss.Border {
	return lipgloss.Border{
		Top:         common.Config.BorderTop,
		Bottom:      common.Config.BorderBottom,
		Left:        common.Config.BorderLeft,
		Right:       common.Config.BorderRight,
		TopLeft:     common.Config.BorderTopLeft,
		TopRight:    common.Config.BorderTopRight,
		BottomLeft:  common.Config.BorderBottomLeft,
		BottomRight: common.Config.BorderBottomRight,
		MiddleLeft:  common.Config.BorderMiddleLeft,
		MiddleRight: common.Config.BorderMiddleRight,
	}
}
