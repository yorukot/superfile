package preview

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/yorukot/superfile/src/internal/common"
)

func getBatSyntaxHighlightedContent(
	itemPath string,
	startLine int,
	previewLine int,
	background string,
	batCmd string,
) (string, bool, error) {
	// --plain: use the plain style without line numbers and decorations
	// --force-colorization: force colorization for non-interactive shell output
	// --line-range m:n: only read from line m to line n (1-indexed)
	firstLine := startLine + 1
	lastLine := startLine + previewLine
	batArgs := []string{
		itemPath,
		"--plain",
		"--force-colorization",
		"--line-range",
		fmt.Sprintf("%d:%d", firstLine, lastLine),
	}

	// set timeout for the external command execution to 500ms max
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultPreviewTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, batCmd, batArgs...)

	fileContentBytes, err := cmd.Output()
	if err != nil {
		slog.Error("Error render code highlight", "error", err)
		return "", false, err
	}

	fileContent := string(fileContentBytes)
	if !common.Config.TransparentBackground {
		fileContent = setBatBackground(fileContent, background)
	}
	lineCount := strings.Count(fileContent, "\n")
	if fileContent != "" && !strings.HasSuffix(fileContent, "\n") {
		lineCount++
	}
	hasMore := previewLine > 0 && lineCount >= previewLine
	return fileContent, hasMore, nil
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
	return common.ImageExtensions[strings.ToLower(filepath.Ext(filename))]
}
