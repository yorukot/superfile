package internal

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Executes cmd, and returns the DirSizeMsg from it, if any
func findDirSizeMsg(cmd tea.Cmd) (DirSizeMsg, bool) {
	if cmd == nil {
		return DirSizeMsg{}, false
	}
	switch msg := ExecuteTeaCmdWithTimeout(cmd, DefaultTestTimeout).(type) {
	case DirSizeMsg:
		return msg, true
	case tea.BatchMsg:
		for _, curCmd := range msg {
			if res, ok := findDirSizeMsg(curCmd); ok {
				return res, true
			}
		}
	}
	return DirSizeMsg{}, false
}

func TestCalculateDirSize(t *testing.T) {
	spaceKey := tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}
	require.Contains(t, common.Hotkeys.CalculateDirSize, spaceKey.String())

	t.Run("Size is shown in the file panel", func(t *testing.T) {
		curTestDir := t.TempDir()
		dir1 := filepath.Join(curTestDir, "dir1")
		utils.SetupDirectories(t, filepath.Join(dir1, "sub"))
		utils.SetupFilesWithData(t, []byte("12345"), filepath.Join(dir1, "f1"), filepath.Join(dir1, "sub", "f2"))

		// Size column is only rendered with extra columns enabled
		origExtraColumns := common.Config.FilePanelExtraColumns
		common.Config.FilePanelExtraColumns = 1
		t.Cleanup(func() { common.Config.FilePanelExtraColumns = origExtraColumns })

		m := defaultTestModel(curTestDir)
		TeaUpdate(m, tea.WindowSizeMsg{Width: 200, Height: 40})
		setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), dir1)
		assert.Contains(t, m.getFocusedFilePanel().Render(true), "2 items")

		msg, ok := findDirSizeMsg(TeaUpdate(m, spaceKey))
		require.True(t, ok, "Expected a DirSizeMsg")
		assert.Contains(t, m.getFocusedFilePanel().Render(true), "???")

		TeaUpdate(m, msg)
		assert.Contains(t, m.getFocusedFilePanel().Render(true), common.FormatFileSize(10))
		assert.NotContains(t, m.getFocusedFilePanel().Render(true), "2 items")
		assert.False(t, m.notifyModel.IsOpen())
	})

	t.Run("Error is shown when directory cannot be read", func(t *testing.T) {
		curTestDir := t.TempDir()
		dir1 := filepath.Join(curTestDir, "dir1")
		utils.SetupDirectories(t, dir1)

		m := defaultTestModel(curTestDir)
		setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), dir1)
		require.NoError(t, os.Remove(dir1))

		msg, ok := findDirSizeMsg(TeaUpdate(m, spaceKey))
		require.True(t, ok, "Expected a DirSizeMsg")
		TeaUpdate(m, msg)
		assert.True(t, m.notifyModel.IsOpen())
	})
}
