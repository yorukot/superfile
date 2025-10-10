package preview

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/ansichroma"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/internal/utils"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

type Model struct {
	open           bool
	width          int
	height         int
	location       string
	content        string
	imagePreviewer *filepreview.ImagePreviewer
	batCmd         string
}

func New() Model {
	return Model{
		open:           common.Config.DefaultOpenFilePreview,
		imagePreviewer: filepreview.NewImagePreviewer(),
		// TODO:  This is an IO operation, move to async ?
		batCmd: checkBatCmd(),
	}
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) Close() {
	m.open = false
}

// Simple rendered string with given text
func (m *Model) RenderText(text string) string {
	return ui.FilePreviewPanelRenderer(m.height, m.width).
		AddLines(text).
		Render() + m.imagePreviewer.ClearAllImages()
}

func (m *Model) SetContentWithRenderText(text string) {
	m.content = m.RenderText(text)
}

func (m *Model) GetContent() string {
	return m.content
}

func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) GetHeight() int {
	return m.height
}

func (m *Model) GetLocation() string {
	return m.location
}

func (m *Model) SetOpen(open bool) {
	m.open = open
}

func (m *Model) SetLocation(location string) {
	m.location = location
}

func (m *Model) SetContent(content string) {
	m.content = content
}

func (m *Model) ToggleOpen() {
	m.open = !m.open
}

func renderFileInfoError(r *rendering.Renderer, err error) string {
	slog.Error("Error get file info", "error", err)
	return r.Render()
}

func renderUnsupportedFormat(box lipgloss.Style) string {
	return box.Render(common.FilePreviewUnsupportedFormatText)
}

func renderUnsupportedFileMode(r *rendering.Renderer) string {
	r.AddLines(common.FilePreviewUnsupportedFileMode)
	return r.Render()
}

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
		style := common.GetElementIcon(file.Name(), file.IsDir(), common.Config.Nerdfont)
		res := lipgloss.NewStyle().Foreground(lipgloss.Color(style.Color)).Background(common.FilePanelBGColor).
			Render(style.Icon+" ") + common.FilePanelStyle.Render(file.Name())
		r.AddLines(res)
	}
	return r.Render()
}

func (m *Model) renderImagePreview(box lipgloss.Style, itemPath string, previewWidth,
	previewHeight int, sideAreaWidth int) string {
	if !m.open {
		return box.Render("\n --- Preview panel is closed ---")
	}

	if !common.Config.ShowImagePreview {
		return box.Render("\n --- Image preview is disabled ---")
	}

	// Use the new auto-detection function to choose the best renderer
	imageRender, err, renderedType := m.imagePreviewer.ImagePreview(itemPath, previewWidth, previewHeight,
		common.Theme.FilePanelBG, sideAreaWidth)
	if errors.Is(err, image.ErrFormat) {
		return box.Render("\n --- " + icon.Error + " Unsupported image formats ---")
	}

	if err != nil {
		slog.Error("Error convert image to ansi", "error", err)
		return box.Render("\n --- " + icon.Error + " Error convert image to ansi ---")
	}

	// Check if this looks like Kitty protocol output (starts with escape sequences)
	// For Kitty protocol, avoid using lipgloss alignment to prevent layout drift
	if renderedType != filepreview.RendererANSI {
		rendered := common.FilePreviewBox(previewHeight, previewWidth).Render(imageRender)
		return rendered
	}

	// For ANSI output, we can safely use vertical alignment
	return box.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(imageRender)
}

func (m *Model) renderTextPreview(r *rendering.Renderer, box lipgloss.Style, itemPath string,
	previewWidth, previewHeight int) string {
	format := lexers.Match(filepath.Base(itemPath))
	if format == nil {
		isText, err := common.IsTextFile(itemPath)
		if err != nil {
			slog.Error("Error while checking text file", "error", err)
			return box.Render(common.FilePreviewError)
		} else if !isText {
			return box.Render(common.FilePreviewUnsupportedFormatText)
		}
	}

	fileContent, err := utils.ReadFileContent(itemPath, previewWidth, previewHeight)
	if err != nil {
		slog.Error("Error open file", "error", err)
		return box.Render(common.FilePreviewError)
	}

	if fileContent == "" {
		return box.Render(common.FilePreviewEmptyText)
	}

	if format != nil {
		background := ""
		if !common.Config.TransparentBackground {
			background = common.Theme.FilePanelBG
		}
		if common.Config.CodePreviewer == "bat" {
			if m.batCmd == "" {
				return box.Render("\n --- " + icon.Error +
					" 'bat' is not installed or not found. ---\n --- Cannot render file preview. ---")
			}
			fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewHeight, background, m.batCmd)
		} else {
			fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name,
				common.Theme.CodeSyntaxHighlightTheme, background)
		}
		if err != nil {
			slog.Error("Error render code highlight", "error", err)
			return box.Render("\n" + common.FilePreviewError)
		}
	}

	r.AddLines(fileContent)
	return r.Render()
}

func (m *Model) RenderWithPath(itemPath string, fullModelWidth int) string {
	previewHeight := m.height
	previewWidth := m.width

	box := common.FilePreviewBox(previewHeight, previewWidth)
	r := ui.FilePreviewPanelRenderer(previewHeight, previewWidth)
	clearCmd := m.imagePreviewer.ClearAllImages()

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
		return renderUnsupportedFormat(box) + clearCmd
	}

	if fileInfo.IsDir() {
		return renderDirectoryPreview(r, itemPath, previewHeight) + clearCmd
	}

	if isImageFile(itemPath) {
		return m.renderImagePreview(box, itemPath, previewWidth, previewHeight, fullModelWidth-previewWidth+1)
	}

	return m.renderTextPreview(r, box, itemPath, previewWidth, previewHeight) + clearCmd
}

func getBatSyntaxHighlightedContent(
	itemPath string,
	previewLine int,
	background string,
	batCmd string,
) (string, error) {
	// --plain: use the plain style without line numbers and decorations
	// --force-colorization: force colorization for non-interactive shell output
	// --line-range <:m>: only read from line 1 to line "m"
	batArgs := []string{itemPath, "--plain", "--force-colorization", "--line-range", fmt.Sprintf(":%d", previewLine-1)}

	// set timeout for the external command execution to 500ms max
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, batCmd, batArgs...)

	fileContentBytes, err := cmd.Output()
	if err != nil {
		slog.Error("Error render code highlight", "error", err)
		return "", err
	}

	fileContent := string(fileContentBytes)
	if !common.Config.TransparentBackground {
		fileContent = setBatBackground(fileContent, background)
	}
	return fileContent, nil
}

func setBatBackground(input string, background string) string {
	tokens := strings.Split(input, "\x1b[0m")
	backgroundStyle := lipgloss.NewStyle().Background(lipgloss.Color(background))
	for idx, token := range tokens {
		tokens[idx] = backgroundStyle.Render(token)
	}
	return strings.Join(tokens, "\x1b[0m")
}

// Check if bat is an executable in PATH and whether to use bat or batcat as command
func checkBatCmd() string {
	if _, err := exec.LookPath("bat"); err == nil {
		return "bat"
	}
	// on ubuntu bat executable is called batcat
	if _, err := exec.LookPath("batcat"); err == nil {
		return "batcat"
	}
	return ""
}

func isImageFile(filename string) bool {
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".svg":  true,
		".webp": true,
		".ico":  true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return imageExtensions[ext]
}
