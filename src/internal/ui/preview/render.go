package preview

import (
	"errors"
	"image"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/yorukot/ansichroma"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func renderDirectoryPreview(r *rendering.Renderer, itemPath string, previewHeight int) string {
	files, err := os.ReadDir(itemPath)
	if err != nil {
		slog.Error("Error render directory preview", "error", err)
		r.AddLines(common.FilePreviewDirectoryUnreadableText)
		return r.Render()
	}

	if len(files) == 0 {
		r.AddLines(common.FilePreviewEmptyText)
		return r.Render()
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})

	for i := 0; i < previewHeight && i < len(files); i++ {
		file := files[i]
		isLink := false
		if info, err := file.Info(); err == nil {
			isLink = info.Mode()&os.ModeSymlink != 0
		}
		style := common.GetElementIcon(file.Name(), file.IsDir(), isLink, common.Config.Nerdfont)
		res := lipgloss.NewStyle().Foreground(lipgloss.Color(style.Color)).Background(common.FilePanelBGColor).
			Render(style.Icon+" ") + common.FilePanelStyle.Render(file.Name())
		r.AddLines(res)
	}
	return r.Render()
}

// renderImagePreview returns (render, rawTransmit). rawTransmit is non-empty
// only for Kitty protocol and must be sent via tea.Raw().
func (m *Model) renderImagePreview(r *rendering.Renderer, itemPath string, previewWidth,
	previewHeight int, sideAreaWidth int, kittyClear string,
) (string, string) {
	if !m.open {
		return r.AddLines(common.FilePreviewPanelClosedText).Render(), kittyClear
	}

	if !common.Config.ShowImagePreview {
		return r.AddLines(common.FilePreviewImagePreviewDisabledText).Render(), kittyClear
	}

	imageRender, rawTransmit, err := m.imagePreviewer.ImagePreview(itemPath, previewWidth, previewHeight,
		common.Theme.FilePanelBG, sideAreaWidth)
	if errors.Is(err, image.ErrFormat) {
		return r.AddLines(common.FilePreviewUnsupportedImageFormatsText).Render(), kittyClear
	}

	if err != nil {
		slog.Error("Error convert image to ansi", "error", err)
		return r.AddLines(common.FilePreviewImageConversionErrorText).Render(), kittyClear
	}

	// For Kitty placeholders or ANSI output, use vertical alignment
	return r.AddStyleModifier(func(s lipgloss.Style) lipgloss.Style {
		return s.AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center)
	}).AddLines(imageRender).Render(), rawTransmit
}

func (m *Model) renderTextPreview(r *rendering.Renderer, itemPath string,
	previewWidth, previewHeight int,
) string {
	format := lexers.Match(filepath.Base(itemPath))
	if format == nil {
		isText, err := common.IsTextFile(itemPath)
		if err != nil {
			slog.Error("Error while checking text file", "error", err)
			return r.AddLines(common.FilePreviewError).Render()
		} else if !isText {
			return r.AddLines(common.FilePreviewUnsupportedFormatText).Render()
		}
	}

	fileContent, err := utils.ReadFileContent(itemPath, previewWidth, previewHeight)
	if err != nil {
		slog.Error("Error open file", "error", err)
		return r.AddLines(common.FilePreviewError).Render()
	}

	if fileContent == "" {
		return r.AddLines(common.FilePreviewEmptyText).Render()
	}

	if format != nil {
		background := ""
		if !common.Config.TransparentBackground {
			background = common.Theme.FilePanelBG
		}
		if common.Config.CodePreviewer == "bat" {
			if m.batCmd == "" {
				return r.AddLines(common.FilePreviewBatNotInstalledText).Render()
			}
			fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewHeight, background, m.batCmd)
		} else {
			fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name,
				common.Theme.CodeSyntaxHighlightTheme, background)
		}
		if err != nil {
			slog.Error("Error render code highlight", "error", err)
			return r.AddLines(common.FilePreviewError).Render()
		}
	}

	r.AddLines(fileContent)
	return r.Render()
}

// Only use this when height and width are synced with filemodel's expectations
func (m *Model) RenderText(text string) string {
	return m.RenderTextWithDimension(text, m.contentHeight, m.contentWidth)
}

func (m *Model) RenderTextWithDimension(text string, height int, width int) string {
	// For zero size, don't need to render anything. Its kinda hack, but
	// its to prevent error logs
	if width == 0 && height == 0 {
		return ""
	}
	return ui.FilePreviewPanelRenderer(height, width).
		AddLines(text).
		Render()
}

// RenderWithPath returns (render, rawTransmit). rawTransmit is non-empty
// for Kitty images (transmit data) or when clearing Kitty images (delete-all).
// It must be sent via tea.Raw().
func (m *Model) RenderWithPath(
	itemPath string,
	previewWidth int,
	previewHeight int,
	fullModelWidth int,
) (string, string) {
	r := ui.FilePreviewPanelRenderer(previewHeight, previewWidth)
	// Raw command to clear any previous Kitty images when showing non-image content
	kittyClear := m.imagePreviewer.GetKittyClearRaw()

	// Adjust dimensions if border is enabled
	contentWidth := previewWidth
	contentHeight := previewHeight
	if common.Config.EnableFilePreviewBorder {
		contentWidth = previewWidth - common.BorderPadding
		contentHeight = previewHeight - common.BorderPadding
	}

	fileInfo, infoErr := os.Stat(itemPath)
	if infoErr != nil {
		slog.Error("Error get file info", "error", infoErr)
		return r.AddLines(common.FilePreviewNoFileInfoText).Render(), kittyClear
	}
	slog.Debug("Attempting to render preview", "itemPath", itemPath,
		"mode", fileInfo.Mode().String(), "isRegular", fileInfo.Mode().IsRegular())

	// For non regular files which are not directories Dont try to read them
	// See Issue #876
	if !fileInfo.Mode().IsRegular() && (fileInfo.Mode()&fs.ModeDir) == 0 {
		return r.AddLines(common.FilePreviewUnsupportedFileMode).Render(), kittyClear
	}

	ext := filepath.Ext(itemPath)
	if slices.Contains(common.UnsupportedPreviewFormats, ext) {
		return r.AddLines(common.FilePreviewUnsupportedFormatText).Render(), kittyClear
	}

	if fileInfo.IsDir() {
		return renderDirectoryPreview(r, itemPath, contentHeight), kittyClear
	}

	if m.thumbnailGenerator != nil && m.thumbnailGenerator.SupportsExt(ext) {
		thumbnailPath, err := m.thumbnailGenerator.GetThumbnailOrGenerate(itemPath)
		if err != nil {
			slog.Error("Error generating thumbnail", "error", err)
			return r.AddLines(common.FilePreviewThumbnailGenerationErrorText).Render(), kittyClear
		}
		return m.renderImagePreview(
			r, thumbnailPath, contentWidth, contentHeight,
			fullModelWidth-previewWidth, kittyClear)
	}

	if isImageFile(itemPath) {
		return m.renderImagePreview(
			r, itemPath, contentWidth, contentHeight,
			fullModelWidth-previewWidth, kittyClear)
	}

	return m.renderTextPreview(r, itemPath, contentWidth, contentHeight), kittyClear
}
