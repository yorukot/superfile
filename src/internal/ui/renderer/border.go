package renderer

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
	fgColor      lipgloss.Color
	bgColor      lipgloss.Color
}

func NewBorderConfig(title string, infoItems []string, fgColor lipgloss.Color, bgColor lipgloss.Color) BorderConfig {
	return BorderConfig{
		title:          title,
		infoItems:      infoItems,
		lipglossBorder: DefaultLipglossBorder(),
		fgColor:        fgColor,
		bgColor:        bgColor,
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
