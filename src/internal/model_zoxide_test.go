package internal

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	tea "charm.land/bubbletea/v2"
	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
)

// This runs a zoxide query and doesn't waits for it.
// Can cause race with zoxide add. See https://github.com/ajeetdsouza/zoxide/issues/1219
func setupProgAndOpenZoxide(t *testing.T, zClient *zoxidelib.Client, dir string) *TeaProg {
	p := setupProgWithZoxide(t, zClient, dir)
	openZoxide(t, p)
	return p
}

func setupProgWithZoxide(t *testing.T, zClient *zoxidelib.Client, dir string) *TeaProg {
	t.Helper()
	common.Config.ZoxideSupport = true
	m := defaultTestModelWithZClient(zClient, dir)
	return NewTestTeaProgWithEventLoop(t, m)
}

func openZoxide(t *testing.T, p *TeaProg) {
	p.SendKey(common.Hotkeys.OpenZoxide[0])
	assert.Eventually(t, func() bool {
		return p.getModel().zoxideModal.IsOpen()
	}, DefaultTestTimeout, DefaultTestTick, "Zoxide modal should open")
}

func updateCurrentFilePanelDirOfTestModel(t *testing.T, p *TeaProg, dir string) {
	err := p.getModel().updateCurrentFilePanelDir(dir)
	require.NoError(t, err, "Failed to navigate to %s", dir)
	assert.Equal(t, dir, p.getModel().getFocusedFilePanel().Location, "Should be in %s after navigation", dir)
}

func TestZoxide(t *testing.T) {
	// Cannot use t.TempDir() here. Bubble Tea intentionally leaks in-flight tea.Cmd
	// goroutines on shutdown (bubbletea/v2 tea.go:697-700). A zoxide query subprocess
	// may still be writing to the data dir after p.Close(), and t.TempDir() treats
	// cleanup failure as a test failure ("directory not empty").
	//
	// Permanent fix: pass a parent context.Context through the program lifecycle
	// (model → zoxide modal → go-zoxide Client) so that p.Close() cancels all
	// in-flight subprocesses. Currently go-zoxide's execCmd creates its own
	// context.Background(), so this requires an API change in go-zoxide.
	zoxideDataDir, err := os.MkdirTemp("", "zoxide-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(zoxideDataDir)

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
	multiSpaceDir := filepath.Join(curTestDir, "test  dir")
	utils.SetupDirectories(t, curTestDir, dir1, dir2, dir3, multiSpaceDir)

	t.Run("Zoxide tracking and navigation", func(t *testing.T) {
		p := setupProgWithZoxide(t, zClient, dir1)
		updateCurrentFilePanelDirOfTestModel(t, p, dir2)
		updateCurrentFilePanelDirOfTestModel(t, p, dir3)

		openZoxide(t, p)

		p.SendKey("dir2")
		assert.Eventually(t, func() bool {
			results := p.getModel().zoxideModal.GetResults()
			return len(results) == 1 && results[0].Path == dir2
		}, DefaultTestTimeout, DefaultTestTick, "dir2 should be found by zoxide UI search")

		// Press enter to navigate to dir2
		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		// Wait for both modal to close AND location to change to avoid race condition
		assert.Eventually(t, func() bool {
			return !p.getModel().zoxideModal.IsOpen() &&
				p.getModel().getFocusedFilePanel().Location == dir2
		}, DefaultTestTimeout, DefaultTestTick,
			"Zoxide modal should close and navigate to %s (current location: %s)",
			dir2, p.getModel().getFocusedFilePanel().Location)
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
		p := setupProgAndOpenZoxide(t, zClient, dir1)

		initialWidth := p.getModel().zoxideModal.GetWidth()
		initialMaxHeight := p.getModel().zoxideModal.GetMaxHeight()

		p.SendDirectly(tea.WindowSizeMsg{Width: 2 * DefaultTestModelWidth, Height: 2 * DefaultTestModelHeight})

		updatedWidth := p.getModel().zoxideModal.GetWidth()
		updatedMaxHeight := p.getModel().zoxideModal.GetMaxHeight()
		assert.Greater(t, updatedWidth, initialWidth, "Width should increase with larger window")
		assert.Greater(t, updatedMaxHeight, initialMaxHeight, "MaxHeight should increase with larger window")
	})

	t.Run("Zoxide 'z' key suppression on open", func(t *testing.T) {
		p := setupProgAndOpenZoxide(t, zClient, dir1)
		assert.Empty(t, p.getModel().zoxideModal.GetTextInputValue(),
			"The 'z' key should not be added to textInput")
		p.SendKeyDirectly("abc")
		assert.Equal(t, "abc", p.getModel().zoxideModal.GetTextInputValue())
	})

	t.Run("Multi-space directory name navigation", func(t *testing.T) {
		p := setupProgWithZoxide(t, zClient, dir1)

		updateCurrentFilePanelDirOfTestModel(t, p, multiSpaceDir)
		updateCurrentFilePanelDirOfTestModel(t, p, dir1)

		openZoxide(t, p)

		p.SendKey(filepath.Base(multiSpaceDir))
		assert.Eventually(t, func() bool {
			results := p.getModel().zoxideModal.GetResults()
			for _, result := range results {
				if result.Path == multiSpaceDir {
					return true
				}
			}
			return false
		}, DefaultTestTimeout, DefaultTestTick, "Multi-space directory should be found by zoxide")

		// Reset textinput via Close-Open
		p.SendKey(common.Hotkeys.Quit[0])
		p.SendKey(common.Hotkeys.OpenZoxide[0])

		p.SendKey("di r 1")
		assert.Eventually(t, func() bool {
			results := p.getModel().zoxideModal.GetResults()
			for _, result := range results {
				if result.Path == dir1 {
					return true
				}
			}
			return false
		}, DefaultTestTimeout, DefaultTestTick, "dir1 should be found by zoxide")
	})

	t.Run("Zoxide escape key closes modal", func(t *testing.T) {
		p := setupProgAndOpenZoxide(t, zClient, dir1)
		p.SendKeyDirectly(common.Hotkeys.CancelTyping[0])
		assert.False(t, p.getModel().zoxideModal.IsOpen(),
			"Zoxide modal should close on escape key")
	})
}
