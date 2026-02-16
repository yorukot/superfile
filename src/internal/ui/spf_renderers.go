package ui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func SidebarRenderer(totalHeight int, totalWidth int, sidebarFocused bool) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.SidebarFGColor
	cfg.ContentBGColor = common.SidebarBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.SidebarBGColor
	cfg.BorderFGColor = common.SidebarBorderColor
	if sidebarFocused {
		cfg.BorderFGColor = common.SidebarBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()
	cfg.RendererName += "-sidebar"

	r := rendering.NewRendererWithAutoFixConfig(cfg)

	return r
}

func FilePanelRenderer(totalHeight int, totalWidth int, filePanelFocused bool) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.FilePanelFGColor
	cfg.ContentBGColor = common.FilePanelBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.FilePanelBGColor
	cfg.BorderFGColor = common.FilePanelBorderColor
	if filePanelFocused {
		cfg.BorderFGColor = common.FilePanelBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()
	cfg.RendererName += "-filepanel"

	r := rendering.NewRendererWithAutoFixConfig(cfg)
	return r
}

func FilePreviewPanelRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)
	cfg.ContentFGColor = common.FilePanelFGColor
	cfg.ContentBGColor = common.FilePanelBGColor

	if common.Config.EnableFilePreviewBorder {
		cfg.BorderRequired = true
		cfg.BorderBGColor = common.FilePanelBGColor
		cfg.BorderFGColor = common.FilePanelBorderColor
		cfg.Border = DefaultLipglossBorder()
	}
	cfg.RendererName += "-preview"

	r := rendering.NewRendererWithAutoFixConfig(cfg)
	return r
}

func PromptRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)
	cfg.TruncateHeight = true
	cfg.ContentFGColor = common.ModalFGColor
	cfg.ContentBGColor = common.ModalBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.ModalBGColor
	cfg.BorderFGColor = common.ModalBorderActiveColor

	cfg.Border = DefaultLipglossBorder()

	r := rendering.NewRendererWithAutoFixConfig(cfg)
	return r
}

func ZoxideRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	return PromptRenderer(totalHeight, totalWidth)
}

func HelpMenuRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)
	cfg.ContentFGColor = common.ModalFGColor
	cfg.ContentBGColor = common.ModalBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.ModalBGColor
	cfg.BorderFGColor = common.ModalBorderActiveColor

	cfg.Border = DefaultLipglossBorder()

	r := rendering.NewRendererWithAutoFixConfig(cfg)
	return r
}

func DefaultFooterRenderer(totalHeight int, totalWidth int, focused bool, name string) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.FooterFGColor
	cfg.ContentBGColor = common.FooterBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.FooterBGColor
	cfg.BorderFGColor = common.FooterBorderColor
	if focused {
		cfg.BorderFGColor = common.FooterBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()
	cfg.RendererName = name

	r := rendering.NewRendererWithAutoFixConfig(cfg)

	r.SetBorderTitle(name)
	return r
}

func ProcessBarRenderer(totalHeight int, totalWidth int, processBarFocused bool) *rendering.Renderer {
	return DefaultFooterRenderer(totalHeight, totalWidth, processBarFocused, "Processes")
}

func MetadataRenderer(totalHeight int, totalWidth int, metadataFocused bool) *rendering.Renderer {
	return DefaultFooterRenderer(totalHeight, totalWidth, metadataFocused, "Metadata")
}

func ClipboardRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	return DefaultFooterRenderer(totalHeight, totalWidth, false, "Clipboard")
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
