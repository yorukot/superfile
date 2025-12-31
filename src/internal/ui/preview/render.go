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
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/ansichroma"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/internal/utils"
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

func (m *Model) renderImagePreview(r *rendering.Renderer, itemPath string, previewWidth,
	previewHeight int, sideAreaWidth int, clearCmd string,
) string {
	if !m.open {
		return r.AddLines("\n --- Preview panel is closed ---").Render() + clearCmd
	}

	if !common.Config.ShowImagePreview {
		return r.AddLines("\n --- Image preview is disabled ---").Render() + clearCmd
	}

	// Use the new auto-detection function to choose the best renderer
	imageRender, err := m.imagePreviewer.ImagePreview(itemPath, previewWidth, previewHeight,
		common.Theme.FilePanelBG, sideAreaWidth)
	if errors.Is(err, image.ErrFormat) {
		return r.AddLines("\n --- "+icon.Error+" Unsupported image formats ---").Render() + clearCmd
	}

	if err != nil {
		slog.Error("Error convert image to ansi", "error", err)
		return r.AddLines("\n --- "+icon.Error+" Error convert image to ansi ---").Render() + clearCmd
	}

	// Check if this looks like Kitty protocol output (starts with escape sequences)
	// For Kitty protocol, avoid using lipgloss alignment to prevent layout drift
	if strings.HasPrefix(imageRender, "\x1b_G") {
		r.AddLines(imageRender)
		return r.Render()
	}

	// For ANSI output, we can safely use vertical alignment
	return r.AddStyleModifier(func(s lipgloss.Style) lipgloss.Style {
		return s.AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center)
	}).AddLines(imageRender).Render() + clearCmd
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
				return r.AddLines("",
					" --- "+icon.Error+" 'bat' is not installed or not found. ---",
					" --- Cannot render file preview. ---").Render()
			}
			fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewHeight, background, m.batCmd)
		} else {
			fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name,
				common.Theme.CodeSyntaxHighlightTheme, background)
		}
		if err != nil {
			slog.Error("Error render code highlight", "error", err)
			return r.AddLines("", common.FilePreviewError).Render()
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
	clearCmd := m.imagePreviewer.ClearKittyImages()
	if width == 0 && height == 0 {
		return clearCmd
	}
	return ui.FilePreviewPanelRenderer(height, width).
		AddLines(text).
		Render() + clearCmd
}

func (m *Model) RenderWithPath(itemPath string, previewWidth int, previewHeight int, fullModelWidth int) string {
	r := ui.FilePreviewPanelRenderer(previewHeight, previewWidth)
	clearCmd := m.imagePreviewer.ClearKittyImages()

	fileInfo, infoErr := os.Stat(itemPath)
	if infoErr != nil {
		return renderFileInfoError(r, infoErr) + clearCmd
	}
	slog.Debug("Attempting to render preview", "itemPath", itemPath,
		"mode", fileInfo.Mode().String(), "isRegular", fileInfo.Mode().IsRegular())

	// For non regular files which are not directories Dont try to read them
	// See Issue #876
	if !fileInfo.Mode().IsRegular() && (fileInfo.Mode()&fs.ModeDir) == 0 {
		return renderUnsupportedFileMode(r) + clearCmd
	}

	ext := filepath.Ext(itemPath)
	if slices.Contains(common.UnsupportedPreviewFormats, ext) {
		return renderUnsupportedFormat(r) + clearCmd
	}

	if fileInfo.IsDir() {
		return renderDirectoryPreview(r, itemPath, previewHeight) + clearCmd
	}

	if m.thumbnailGenerator != nil && m.thumbnailGenerator.SupportsExt(ext) {
		thumbnailPath, err := m.thumbnailGenerator.GetThumbnailOrGenerate(itemPath)
		if err != nil {
			slog.Error("Error generating thumbnail", "error", err)
			return renderThumbnailGenerationError(r) + clearCmd
		}
		// Notes : If renderImagePreview fails, and return some error message
		// render, then we dont apply clearCmd. This might cause issues.
		// same for below usage of renderImagePreview
		return m.renderImagePreview(r, thumbnailPath, previewWidth, previewHeight, fullModelWidth-previewWidth+1, clearCmd)
	}

	if isImageFile(itemPath) {
		return m.renderImagePreview(r, itemPath, previewWidth, previewHeight, fullModelWidth-previewWidth+1, clearCmd)
	}

	return m.renderTextPreview(r, itemPath, previewWidth, previewHeight) + clearCmd
}
