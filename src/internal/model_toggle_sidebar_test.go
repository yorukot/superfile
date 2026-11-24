package internal

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Redirect the persisted toggle state, so tests never touch the real data dir.
func setupToggleSidebarTest(t *testing.T, sidebarWidth int) string {
	t.Helper()
	oldPath, oldWidth := variable.ToggleSidebar, common.Config.SidebarWidth
	//nolint:reassign // Needed to tests
	variable.ToggleSidebar = filepath.Join(t.TempDir(), "toggleSidebar")
	common.Config.SidebarWidth = sidebarWidth
	t.Cleanup(func() {
		//nolint:reassign // Needed to tests
		variable.ToggleSidebar = oldPath
		common.Config.SidebarWidth = oldWidth
	})
	return variable.ToggleSidebar
}

// A directory with files, as the layout validation expects a valid cursor
func newToggleSidebarTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	utils.SetupFiles(t, filepath.Join(dir, "file1.txt"), filepath.Join(dir, "file2.txt"))
	return dir
}

func TestToggleSidebar(t *testing.T) {
	toggleKey := utils.TeaRuneKeyMsg(common.Hotkeys.ToggleSidebar[0])

	t.Run("Hide and show", func(t *testing.T) {
		stateFile := setupToggleSidebarTest(t, 20)
		m := defaultTestModel(newToggleSidebarTestDir(t))
		fullSidebarWidth := m.sidebarModel.GetWidth()
		require.Equal(t, 20+common.BorderPadding, fullSidebarWidth)
		require.Equal(t, m.fullWidth, fullSidebarWidth+m.fileModel.Width)

		TeaUpdate(m, toggleKey)
		assert.True(t, m.sidebarModel.Disabled())
		assert.Zero(t, m.sidebarModel.GetWidth())
		assert.Equal(t, m.fullWidth, m.fileModel.Width)
		assert.Zero(t, m.fileModel.SidebarWidth)
		assert.False(t, utils.ReadBoolFile(stateFile, true), "hidden state must be persisted")
		assertLayoutValidity(t, m)

		TeaUpdate(m, toggleKey)
		assert.False(t, m.sidebarModel.Disabled())
		assert.Equal(t, fullSidebarWidth, m.sidebarModel.GetWidth())
		assert.Equal(t, m.fullWidth, fullSidebarWidth+m.fileModel.Width)
		assert.Equal(t, fullSidebarWidth, m.fileModel.SidebarWidth)
		assert.True(t, utils.ReadBoolFile(stateFile, false), "shown state must be persisted")
		assertLayoutValidity(t, m)
	})

	t.Run("Hidden state restored at startup", func(t *testing.T) {
		setupToggleSidebarTest(t, 20)
		dirs := []string{newToggleSidebarTestDir(t)}
		m := setModelParamsForTest(defaultModelConfig(false, false, false, false, dirs, nil), true)
		assert.True(t, m.sidebarModel.Disabled())
		assert.Equal(t, m.fullWidth, m.fileModel.Width)
		assertLayoutValidity(t, m)
	})

	t.Run("Focus leaves a hidden sidebar", func(t *testing.T) {
		setupToggleSidebarTest(t, 20)
		m := defaultTestModel(newToggleSidebarTestDir(t))
		focusKey := utils.TeaRuneKeyMsg(common.Hotkeys.FocusOnSidebar[0])

		TeaUpdate(m, focusKey)
		require.Equal(t, sidebarFocus, m.focusPanel)

		TeaUpdate(m, toggleKey)
		assert.Equal(t, nonePanelFocus, m.focusPanel)
		assert.True(t, m.getFocusedFilePanel().IsFocused)

		// A hidden sidebar must not take focus
		TeaUpdate(m, focusKey)
		assert.Equal(t, nonePanelFocus, m.focusPanel)
	})

	t.Run("sidebar_width = 0 always wins", func(t *testing.T) {
		stateFile := setupToggleSidebarTest(t, 0)
		m := defaultTestModel(newToggleSidebarTestDir(t))
		require.True(t, m.sidebarModel.Disabled())

		TeaUpdate(m, toggleKey)
		TeaUpdate(m, toggleKey)
		assert.True(t, m.sidebarModel.Disabled(), "showing must not override the config")
		assert.Zero(t, m.sidebarModel.GetWidth())
		assert.Equal(t, m.fullWidth, m.fileModel.Width)
		assert.NoFileExists(t, stateFile, "a disabled sidebar must not persist toggle state")
		assertLayoutValidity(t, m)
	})
}
