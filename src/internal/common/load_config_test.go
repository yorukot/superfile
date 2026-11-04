package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
)

func TestReloadThemeFile(t *testing.T) {
	themeDir := t.TempDir()
	themeFile := filepath.Join(themeDir, "reload-test.toml")

	origThemeFolder := variable.ThemeFolder
	origTheme := Theme
	origConfigTheme := Config.Theme
	t.Cleanup(func() {
		//nolint:reassign // Needed to restore the theme folder after the test
		variable.ThemeFolder = origThemeFolder
		Theme = origTheme
		Config.Theme = origConfigTheme
	})

	//nolint:reassign // Needed to point theme loading at a temporary folder
	variable.ThemeFolder = themeDir
	Config.Theme = "reload-test"
	Theme = ThemeType{FilePanelFG: "#111111", GradientColor: []string{"#000000", "#ffffff"}}

	writeTheme := func(content string) {
		t.Helper()
		require.NoError(t, os.WriteFile(themeFile, []byte(content), 0o600))
	}

	t.Run("valid theme replaces current theme", func(t *testing.T) {
		writeTheme("file_panel_fg = \"#abcdef\"\ngradient_color = [\"#123456\", \"#654321\"]\n")
		require.NoError(t, ReloadThemeFile())
		assert.Equal(t, "#abcdef", Theme.FilePanelFG)
		assert.Equal(t, []string{"#123456", "#654321"}, Theme.GradientColor)
	})

	t.Run("invalid toml keeps current theme", func(t *testing.T) {
		writeTheme("this is = not [valid\n")
		require.Error(t, ReloadThemeFile())
		assert.Equal(t, "#abcdef", Theme.FilePanelFG)
	})

	t.Run("wrong gradient length keeps current theme", func(t *testing.T) {
		writeTheme("file_panel_fg = \"#999999\"\ngradient_color = [\"#123456\"]\n")
		require.Error(t, ReloadThemeFile())
		assert.Equal(t, "#abcdef", Theme.FilePanelFG)
		assert.Equal(t, []string{"#123456", "#654321"}, Theme.GradientColor)
	})

	t.Run("missing file keeps current theme", func(t *testing.T) {
		require.NoError(t, os.Remove(themeFile))
		require.Error(t, ReloadThemeFile())
		assert.Equal(t, "#abcdef", Theme.FilePanelFG)
	})
}
