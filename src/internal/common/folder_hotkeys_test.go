package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestFolderShortcutConfigPreservesOverrides(t *testing.T) {
	defaults, err := os.ReadFile("../../superfile_config/hotkeys.toml")
	require.NoError(t, err)
	path := filepath.Join(t.TempDir(), "hotkeys.toml")
	data := []byte("# keep my settings\nquit = ['ctrl+q']\ngo_to_home = ['alt+b']\ngo_to_downloads = []\n")
	require.NoError(t, os.WriteFile(path, data, utils.ConfigFilePerm))
	var hotkeys HotkeysType
	require.NoError(t, utils.LoadTomlFile(path, string(defaults), &hotkeys, false, true))
	assert.Equal(t, []string{"ctrl+q"}, hotkeys.Quit)
	assert.Equal(t, []string{"alt+b"}, hotkeys.GoToHome)
	assert.Empty(t, hotkeys.GoToDownloads)
	assert.Equal(t, []string{"alt+9"}, hotkeys.GoToPin9)
	after, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, data, after)
}

func TestFolderShortcutDefaultsInBothPresets(t *testing.T) {
	for _, name := range []string{"hotkeys.toml", "vimHotkeys.toml"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("../../superfile_config", name))
			require.NoError(t, err)
			var h HotkeysType
			require.NoError(t, toml.Unmarshal(data, &h))
			shortcuts := h.FolderShortcuts()
			require.Len(t, shortcuts, 19)
			for _, shortcut := range shortcuts {
				assert.NotEmpty(t, shortcut.Keys, shortcut.Name)
				assert.NotEmpty(t, shortcut.AltHint(), shortcut.Name)
			}
		})
	}
}

func TestFolderShortcutConflictsAndHints(t *testing.T) {
	h := HotkeysType{
		OpenHelpMenu:  []string{"alt+d"},
		GoToHome:      []string{"ctrl+g", "alt+b", "alt+d"},
		GoToDownloads: []string{"alt+d", "alt+b", "alt+o"},
		GoToDesktop:   []string{"ctrl+e"},
	}
	shortcuts := h.FolderShortcuts()
	assert.Equal(t, []string{"ctrl+g", "alt+b"}, shortcuts[0].Keys)
	assert.Equal(t, "b", shortcuts[0].AltHint())
	assert.Equal(t, []string{"alt+o"}, shortcuts[1].Keys)
	assert.Empty(t, shortcuts[2].AltHint())
	assert.Equal(t, []string{"ctrl+g", "alt+b", "alt+d"}, h.GoToHome, "filter must not mutate settings")
}
