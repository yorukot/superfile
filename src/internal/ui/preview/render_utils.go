package preview

import (
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/yorukot/superfile/src/internal/common"
)

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
