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

	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/ansichroma"

	"github.com/yorukot/superfile/src/config/icon"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

type Model struct {
	open               bool
	width              int
	height             int
	location           string
	content            string
	imagePreviewer     *filepreview.ImagePreviewer
	batCmd             string
	thumbnailGenerator *filepreview.ThumbnailGenerator
}

func New() Model {
	generator, err := filepreview.NewThumbnailGenerator()
	if err != nil {
		slog.Error("Could not NewThumbnailGenerator object", "error", err)
	}

	// Debug logging for bat command detection
	slog.Debug("[BAT DEBUG] Initializing preview model")
	batCommandPath := checkBatCmd()
	slog.Debug("[BAT DEBUG] bat command check result", "batCmd", batCommandPath, "isEmpty", batCommandPath == "")

	return Model{
		open:               common.Config.DefaultOpenFilePreview,
		imagePreviewer:     filepreview.NewImagePreviewer(),
		thumbnailGenerator: generator,
		// TODO:  This is an IO operation, move to async ?
		batCmd: batCommandPath,
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
		Render() + m.imagePreviewer.ClearKittyImages()
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

func (m *Model) CleanUp() {
	if m.thumbnailGenerator != nil {
		err := m.thumbnailGenerator.CleanUp()
		if err != nil {
			slog.Error("Error While cleaning up TempDirectory", "Error:", err)
		}
	}
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

func (m *Model) renderImagePreview(box lipgloss.Style, itemPath string, previewWidth,
	previewHeight int, sideAreaWidth int,
) string {
	if !m.open {
		return box.Render("\n --- Preview panel is closed ---")
	}

	if !common.Config.ShowImagePreview {
		return box.Render("\n --- Image preview is disabled ---")
	}

	// Use the new auto-detection function to choose the best renderer
	imageRender, err := m.imagePreviewer.ImagePreview(itemPath, previewWidth, previewHeight,
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
	if strings.HasPrefix(imageRender, "\x1b_G") {
		rendered := common.FilePreviewBox(previewHeight, previewWidth).Render(imageRender)
		return rendered
	}

	// For ANSI output, we can safely use vertical alignment
	return box.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(imageRender)
}

func (m *Model) renderTextPreview(r *rendering.Renderer, box lipgloss.Style, itemPath string,
	previewWidth, previewHeight int,
) string {
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
	slog.Debug("[BAT DEBUG] Code preview check", "CodePreviewer", common.Config.CodePreviewer, "format", format.Config().Name)
	if common.Config.CodePreviewer == "bat" {
		slog.Debug("[BAT DEBUG] Using bat previewer", "batCmd", m.batCmd, "isEmpty", m.batCmd == "")
		if m.batCmd == "" {
			slog.Error("[BAT DEBUG] bat command is empty - bat not found", "batCmd", m.batCmd)
			return box.Render("\n --- " + icon.Error +
				" 'bat' is not installed or not found. ---\n --- Cannot render file preview. ---")
		}
		slog.Debug("[BAT DEBUG] Calling getBatSyntaxHighlightedContent", "itemPath", itemPath, "previewHeight", previewHeight, "batCmd", m.batCmd)
		fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewHeight, background, m.batCmd)
		if err != nil {
			slog.Error("[BAT DEBUG] bat execution failed", "error", err)
		} else {
			slog.Debug("[BAT DEBUG] bat execution successful", "contentLength", len(fileContent))
		}
	} else {
		slog.Debug("[BAT DEBUG] Using chroma previewer (not bat)")
		fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name,
			common.Theme.CodeSyntaxHighlightTheme, background)
	}
		if err != nil {
			slog.Error("[BAT DEBUG] Error rendering code highlight", "error", err, "previewer", common.Config.CodePreviewer)
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
		return renderUnsupportedFormat(box) + clearCmd
	}

	if fileInfo.IsDir() {
		return renderDirectoryPreview(r, itemPath, previewHeight) + clearCmd
	}

	if isVideoFile(itemPath) {
		if m.thumbnailGenerator == nil {
			return renderUnsupportedFormat(box) + clearCmd
		}
		thumbnailPath, err := m.thumbnailGenerator.GetThumbnailOrGenerate(itemPath)
		if err != nil {
			slog.Error("Error generating thumbnail", "error", err)
			return renderUnsupportedFormat(box) + clearCmd
		}
		return m.renderImagePreview(box, thumbnailPath, previewWidth, previewHeight, fullModelWidth-previewWidth+1)
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
	slog.Debug("[BAT DEBUG] getBatSyntaxHighlightedContent called", "itemPath", itemPath, "previewLine", previewLine, "batCmd", batCmd)
	
	// Check if file exists and is readable
	if fileInfo, err := os.Stat(itemPath); err != nil {
		slog.Error("[BAT DEBUG] File stat error", "error", err, "path", itemPath)
		return "", fmt.Errorf("file stat error: %w", err)
	} else {
		slog.Debug("[BAT DEBUG] File info", "path", itemPath, "size", fileInfo.Size(), "mode", fileInfo.Mode().String())
	}
	
	// --plain: use the plain style without line numbers and decorations
	// --force-colorization: force colorization for non-interactive shell output
	// --line-range <:m>: only read from line 1 to line "m"
	batArgs := []string{itemPath, "--plain", "--force-colorization", "--line-range", fmt.Sprintf(":%d", previewLine-1)}
	slog.Debug("[BAT DEBUG] Executing bat command", "cmd", batCmd, "args", batArgs)

	// set timeout for the external command execution to 500ms max
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultPreviewTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, batCmd, batArgs...)
	slog.Debug("[BAT DEBUG] Command created", "path", cmd.Path, "args", cmd.Args)

	fileContentBytes, err := cmd.Output()
	slog.Debug("[BAT DEBUG] Command executed", "error", err, "outputLen", len(fileContentBytes))
	
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			slog.Error("[BAT DEBUG] Command failed with exit error", "error", err, "stderr", string(exitError.Stderr), "stdout", string(fileContentBytes))
		} else {
			slog.Error("[BAT DEBUG] Command failed", "error", err, "type", fmt.Sprintf("%T", err))
		}
	}
	
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
	slog.Debug("[BAT DEBUG] Checking for bat executable in PATH")
	
	// Check PATH env var
	pathEnv := os.Getenv("PATH")
	slog.Debug("[BAT DEBUG] PATH environment", "PATH", pathEnv)
	
	if _, err := exec.LookPath("bat"); err == nil {
		slog.Debug("[BAT DEBUG] Found 'bat' executable")
		// Get full path
		if fullPath, err := exec.LookPath("bat"); err == nil {
			slog.Debug("[BAT DEBUG] bat full path", "path", fullPath)
			// Test executing bat to ensure it works
			if output, testErr := exec.Command(fullPath, "--version").CombinedOutput(); testErr == nil {
				slog.Debug("[BAT DEBUG] bat version test successful", "output", string(output))
			} else {
				slog.Error("[BAT DEBUG] bat version test failed", "error", testErr, "output", string(output))
			}
		}
		return "bat"
	} else {
		slog.Debug("[BAT DEBUG] 'bat' not found in PATH", "error", err)
	}
	
	// on ubuntu bat executable is called batcat
	if _, err := exec.LookPath("batcat"); err == nil {
		slog.Debug("[BAT DEBUG] Found 'batcat' executable")
		// Get full path
		if fullPath, err := exec.LookPath("batcat"); err == nil {
			slog.Debug("[BAT DEBUG] batcat full path", "path", fullPath)
			// Test executing batcat
			if output, testErr := exec.Command(fullPath, "--version").CombinedOutput(); testErr == nil {
				slog.Debug("[BAT DEBUG] batcat version test successful", "output", string(output))
			} else {
				slog.Error("[BAT DEBUG] batcat version test failed", "error", testErr, "output", string(output))
			}
		}
		return "batcat"
	} else {
		slog.Debug("[BAT DEBUG] 'batcat' not found in PATH", "error", err)
	}
	
	slog.Debug("[BAT DEBUG] No bat executable found")
	return ""
}

func isImageFile(filename string) bool {
	return common.ImageExtensions[strings.ToLower(filepath.Ext(filename))]
}

func isVideoFile(filename string) bool {
	return common.VideoExtensions[strings.ToLower(filepath.Ext(filename))]
}
