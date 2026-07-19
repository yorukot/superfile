package preview

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
	lastLine := startLine + previewLine + 1
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
	lines := strings.Split(strings.TrimSuffix(fileContent, "\n"), "\n")
	if fileContent == "" {
		return "", false, nil
	}
	hasMore := previewLine > 0 && len(lines) > previewLine
	if hasMore {
		fileContent = strings.Join(lines[:previewLine], "\n") + "\n"
	}
	return fileContent, hasMore, nil
}

// countFileLines counts newline-delimited lines in itemPath. Scanning stops when
// common.DefaultPreviewTimeout elapses; the second return value is false when
// the count may be incomplete.
func countFileLines(itemPath string) (int, bool, error) {
	return countFileLinesBefore(itemPath, time.Now().Add(common.DefaultPreviewTimeout))
}

func countFileLinesBefore(itemPath string, deadline time.Time) (int, bool, error) {
	file, err := os.Open(itemPath)
	if err != nil {
		return 0, false, err
	}
	defer file.Close()

	reader := transform.NewReader(file, unicode.BOMOverride(unicode.UTF8.NewDecoder()))
	scanner := bufio.NewScanner(reader)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if time.Now().After(deadline) {
			return lineCount, false, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, false, err
	}
	return lineCount, true, nil
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
