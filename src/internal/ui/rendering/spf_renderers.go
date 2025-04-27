package rendering

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/internal/common"
)

// Todo : rendering package should not be aware of sidebar
func SidebarRenderer(totalHeight int, totalWidth int, sidebarFocussed bool) Renderer {
	cfg := DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.SidebarFGColor
	cfg.ContentBGColor = common.SidebarBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.SidebarBGColor
	cfg.BorderFGColor = common.SidebarBorderColor
	if sidebarFocussed {
		cfg.BorderFGColor = common.SidebarBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()

	return NewRenderer(cfg)
}

func FilePanelRenderer(totalHeight int, totalWidth int, filePanelFocussed bool) Renderer {
	cfg := DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.FilePanelFGColor
	cfg.ContentBGColor = common.FilePanelBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.FilePanelBGColor
	cfg.BorderFGColor = common.FilePanelBorderColor
	if filePanelFocussed {
		cfg.BorderFGColor = common.FilePanelBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()

	return NewRenderer(cfg)
}

func PromptRenderer(totalHeight int, totalWidth int) Renderer {
	cfg := DefaultRendererConfig(totalHeight, totalWidth)
	cfg.TruncateHeight = true
	cfg.ContentFGColor = common.ModalFGColor
	cfg.ContentBGColor = common.ModalBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.ModalBGColor
	cfg.BorderFGColor = common.ModalBorderActiveColor

	cfg.Border = DefaultLipglossBorder()

	return NewRenderer(cfg)
}

// Todo : Move to diff package
func DefaultFooterRenderer(totalHeight int, totalWidth int, focussed bool) Renderer {
	cfg := DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.FooterFGColor
	cfg.ContentBGColor = common.FooterBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.FooterBGColor
	cfg.BorderFGColor = common.FooterBorderColor
	if focussed {
		cfg.BorderFGColor = common.FooterBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()

	return NewRenderer(cfg)
}

func ProcessBarRenderer(totalHeight int, totalWidth int, processBarFocussed bool) Renderer {
	r := DefaultFooterRenderer(totalHeight, totalWidth, processBarFocussed)
	r.SetBorderTitle("Process")
	return r
}

func MetadataRenderer(totalHeight int, totalWidth int, metadataFocussed bool) Renderer {
	r := DefaultFooterRenderer(totalHeight, totalWidth, metadataFocussed)
	// Todo : Move hardcoded string to constant
	r.SetBorderTitle("Metadata")
	return r
}

func ClipboardRenderer(totalHeight int, totalWidth int) Renderer {
	r := DefaultFooterRenderer(totalHeight, totalWidth, false)
	// Todo : Move hardcoded string to constant
	r.SetBorderTitle("Clipboard")
	return r
}

// Todo : Maybe this doesn't belongs in here ? Put it in style functions ?
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
