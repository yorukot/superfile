package internal

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

//nolint:gocognit,maintidx // Integration test with multiple subtests for comprehensive zoxide functionality
func TestZoxide(t *testing.T) {
	zoxideDataDir := t.TempDir()
	zClient, err := zoxidelib.New(zoxidelib.WithDataDir(zoxideDataDir))
	if err != nil {
		if runtime.GOOS != utils.OsLinux {
			t.Skipf("Skipping zoxide tests in non-Linux because zoxide client cannot be initialized")
		} else {
			t.Fatalf("zoxide initialization failed")
		}
	}

	originalZoxideSupport := common.Config.ZoxideSupport
	defer func() {
		common.Config.ZoxideSupport = originalZoxideSupport
	}()

	curTestDir := filepath.Join(testDir, "TestZoxide")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	dir3 := filepath.Join(curTestDir, "dir3")
	utils.SetupDirectories(t, curTestDir, dir1, dir2, dir3)

	t.Run("Zoxide tracking and navigation", func(t *testing.T) {
		common.Config.ZoxideSupport = true
		m := defaultTestModelWithZClient(zClient, dir1)
		p := NewTestTeaProgWithEventLoop(t, m)

		err := p.getModel().updateCurrentFilePanelDir(dir2)
		require.NoError(t, err, "Failed to navigate to dir2")
		assert.Equal(t, dir2, p.getModel().getFocusedFilePanel().location, "Should be in dir2 after navigation")

		err = p.getModel().updateCurrentFilePanelDir(dir3)
		require.NoError(t, err, "Failed to navigate to dir3")
		assert.Equal(t, dir3, p.getModel().getFocusedFilePanel().location, "Should be in dir3 after navigation")

		p.SendKey(common.Hotkeys.OpenZoxide[0])
		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open when pressing 'z' key")

		// Type "dir2" to search for it
		for _, char := range "dir2" {
			p.SendKey(string(char))
		}

		// Wait for async query results to arrive
		assert.Eventually(t, func() bool {
			results := p.getModel().zoxideModal.GetResults()
			if len(results) == 0 {
				return false
			}
			for _, result := range results {
				if result.Path == dir2 {
					return true
				}
			}
			return false
		}, DefaultTestTimeout, DefaultTestTick, "dir2 should be found by zoxide UI search")

		results := p.getModel().zoxideModal.GetResults()
		assert.GreaterOrEqual(t, len(results), 1, "Should have at least 1 directory found by zoxide UI search")

		resultPaths := make([]string, len(results))
		for i, result := range results {
			resultPaths[i] = result.Path
		}
		assert.Contains(t, resultPaths, dir2, "dir2 should be found by zoxide UI search")

		// Find dir2 in results and navigate to it
		dir2Index := -1
		for i, result := range results {
			if result.Path == dir2 {
				dir2Index = i
				break
			}
		}
		require.NotEqual(t, -1, dir2Index, "dir2 should be in results")

		// Navigate to dir2's position in the list
		for range dir2Index {
			p.SendKey(common.Hotkeys.ListDown[0])
		}

		// Press enter to navigate to dir2
		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			return !p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should close after navigation")
		assert.Equal(
			t,
			dir2,
			p.getModel().getFocusedFilePanel().location,
			"Should navigate back to dir2 after zoxide selection",
		)
	})

	t.Run("Zoxide disabled shows no results", func(t *testing.T) {
		common.Config.ZoxideSupport = false
		m := defaultTestModelWithZClient(zClient, dir1)

		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenZoxide[0]))
		assert.True(t, m.zoxideModal.IsOpen(), "Zoxide modal should open even when ZoxideSupport is disabled")

		results := m.zoxideModal.GetResults()
		assert.Empty(t, results, "Zoxide modal should show no results when ZoxideSupport is disabled")
	})

	t.Run("Zoxide modal size on window resize", func(t *testing.T) {
		common.Config.ZoxideSupport = true
		m := defaultTestModelWithZClient(zClient, dir1)
		p := NewTestTeaProgWithEventLoop(t, m)

		p.SendKey(common.Hotkeys.OpenZoxide[0])
		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open")

		initialWidth := p.getModel().zoxideModal.GetWidth()
		initialMaxHeight := p.getModel().zoxideModal.GetMaxHeight()

		newWidth := 4 * common.MinimumWidth
		newHeight := 4 * common.MinimumHeight
		p.Send(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})

		assert.Eventually(t, func() bool {
			updatedWidth := p.getModel().zoxideModal.GetWidth()
			updatedMaxHeight := p.getModel().zoxideModal.GetMaxHeight()
			return updatedWidth != initialWidth && updatedMaxHeight != initialMaxHeight
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal dimensions should update on window resize")

		updatedWidth := p.getModel().zoxideModal.GetWidth()
		updatedMaxHeight := p.getModel().zoxideModal.GetMaxHeight()
		assert.Greater(t, updatedWidth, initialWidth, "Width should increase with larger window")
		assert.Greater(t, updatedMaxHeight, initialMaxHeight, "MaxHeight should increase with larger window")
	})

	t.Run("Zoxide 'z' key suppression on open", func(t *testing.T) {
		common.Config.ZoxideSupport = true
		m := defaultTestModelWithZClient(zClient, dir1)
		p := NewTestTeaProgWithEventLoop(t, m)

		p.SendKey(common.Hotkeys.OpenZoxide[0])
		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open")

		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.GetTextInputValue() == ""
		}, DefaultTestTimeout, DefaultTestTick, "The 'z' key should not be added to textInput")

		p.SendKey("a")
		p.SendKey("b")
		p.SendKey("c")

		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.GetTextInputValue() == "abc"
		}, DefaultTestTimeout, DefaultTestTick, "Subsequent keys should be added to textInput")
	})

	t.Run("Multi-space directory name navigation", func(t *testing.T) {
		common.Config.ZoxideSupport = true
		multiSpaceDir := filepath.Join(curTestDir, "test  dir")
		utils.SetupDirectories(t, multiSpaceDir)
		defer os.RemoveAll(multiSpaceDir)

		m := defaultTestModelWithZClient(zClient, dir1)
		p := NewTestTeaProgWithEventLoop(t, m)

		err := p.getModel().updateCurrentFilePanelDir(multiSpaceDir)
		require.NoError(t, err, "Failed to navigate to multi-space directory")

		err = p.getModel().updateCurrentFilePanelDir(dir1)
		require.NoError(t, err, "Failed to navigate back to dir1")

		p.SendKey(common.Hotkeys.OpenZoxide[0])
		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open")

		for _, char := range "test  dir" {
			p.SendKey(string(char))
		}

		assert.Eventually(t, func() bool {
			results := p.getModel().zoxideModal.GetResults()
			for _, result := range results {
				if result.Path == multiSpaceDir {
					return true
				}
			}
			return false
		}, DefaultTestTimeout, DefaultTestTick, "Multi-space directory should be found by zoxide")

		results := p.getModel().zoxideModal.GetResults()
		multiSpaceDirIndex := -1
		for i, result := range results {
			if result.Path == multiSpaceDir {
				multiSpaceDirIndex = i
				break
			}
		}
		require.NotEqual(t, -1, multiSpaceDirIndex, "Multi-space directory should be in results")

		for range multiSpaceDirIndex {
			p.SendKey(common.Hotkeys.ListDown[0])
		}

		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			return p.getModel().getFocusedFilePanel().location == multiSpaceDir
		}, DefaultTestTimeout, DefaultTestTick, "Should navigate to multi-space directory")
	})

	t.Run("Zoxide escape key closes modal", func(t *testing.T) {
		common.Config.ZoxideSupport = true
		m := defaultTestModelWithZClient(zClient, dir1)
		p := NewTestTeaProgWithEventLoop(t, m)

		p.SendKey(common.Hotkeys.OpenZoxide[0])
		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open")

		p.SendKey(common.Hotkeys.CancelTyping[0])
		assert.Eventually(t, func() bool {
			return !p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should close on escape key")
	})

	t.Run("Zoxide with invalid/non-existent directory", func(t *testing.T) {
		common.Config.ZoxideSupport = true
		m := defaultTestModelWithZClient(zClient, dir1)
		p := NewTestTeaProgWithEventLoop(t, m)

		currentLocation := p.getModel().getFocusedFilePanel().location

		p.SendKey(common.Hotkeys.OpenZoxide[0])
		assert.Eventually(t, func() bool {
			return p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open")

		for _, char := range "uniquexyzabc123" {
			p.SendKey(string(char))
		}

		assert.Eventually(t, func() bool {
			results := p.getModel().zoxideModal.GetResults()
			return len(results) == 0
		}, DefaultTestTimeout, DefaultTestTick, "Should have no results for unique search")

		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			return !p.getModel().zoxideModal.IsOpen()
		}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should close")

		assert.Equal(t, currentLocation, p.getModel().getFocusedFilePanel().location,
			"Should stay in current location when confirming with no results")
	})
}
