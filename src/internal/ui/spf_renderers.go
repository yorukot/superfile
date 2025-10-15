package ui

import (
	"log/slog"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func SidebarRenderer(totalHeight int, totalWidth int, sidebarFocussed bool) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.SidebarFGColor
	cfg.ContentBGColor = common.SidebarBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.SidebarBGColor
	cfg.BorderFGColor = common.SidebarBorderColor
	if sidebarFocussed {
		cfg.BorderFGColor = common.SidebarBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
	return r
}

func FilePanelRenderer(totalHeight int, totalWidth int, filePanelFocussed bool) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.FilePanelFGColor
	cfg.ContentBGColor = common.FilePanelBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.FilePanelBGColor
	cfg.BorderFGColor = common.FilePanelBorderColor
	if filePanelFocussed {
		cfg.BorderFGColor = common.FilePanelBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
	return r
}

func FilePreviewPanelRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)
	cfg.ContentFGColor = common.FilePanelFGColor
	cfg.ContentBGColor = common.FilePanelBGColor
	cfg.BorderRequired = false

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
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

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
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

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
	return r
}

func DefaultFooterRenderer(totalHeight int, totalWidth int, focussed bool) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.FooterFGColor
	cfg.ContentBGColor = common.FooterBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.FooterBGColor
	cfg.BorderFGColor = common.FooterBorderColor
	if focussed {
		cfg.BorderFGColor = common.FooterBorderActiveColor
	}
	cfg.Border = DefaultLipglossBorder()

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
	return r
}

func ProcessBarRenderer(totalHeight int, totalWidth int, processBarFocussed bool) *rendering.Renderer {
	r := DefaultFooterRenderer(totalHeight, totalWidth, processBarFocussed)
	r.SetBorderTitle("Processes")
	return r
}

func MetadataRenderer(totalHeight int, totalWidth int, metadataFocussed bool) *rendering.Renderer {
	r := DefaultFooterRenderer(totalHeight, totalWidth, metadataFocussed)
	r.SetBorderTitle("Metadata")
	return r
}

func ClipboardRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	r := DefaultFooterRenderer(totalHeight, totalWidth, false)
	r.SetBorderTitle("Clipboard")
	return r
}

func BulkRenameRenderer(totalHeight int, totalWidth int) *rendering.Renderer {
	cfg := rendering.DefaultRendererConfig(totalHeight, totalWidth)

	cfg.ContentFGColor = common.ModalFGColor
	cfg.ContentBGColor = common.ModalBGColor

	cfg.BorderRequired = true
	cfg.BorderBGColor = common.ModalBGColor
	cfg.BorderFGColor = common.ModalBorderActiveColor
	cfg.Border = DefaultLipglossBorder()

	r, err := rendering.NewRenderer(cfg)
	if err != nil {
		slog.Error("Error in creating renderer. Falling back to default renderer", "error", err)
		r = &rendering.Renderer{}
	}
	return r
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
