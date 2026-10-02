package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
)

func setupProgForFind(t *testing.T, dir string) *TeaProg {
	t.Helper()
	common.Config.FindFileSupport = true
	m := defaultTestModel(dir)
	return NewTestTeaProgWithEventLoop(t, m)
}

func openFind(t *testing.T, p *TeaProg) {
	t.Helper()
	p.SendKey(common.Hotkeys.FindFile[0])
	assert.Eventually(t, func() bool {
		return p.getModel().findModal.IsOpen()
	}, DefaultTestTimeout, DefaultTestTick, "Find modal should open")
}

func TestFind(t *testing.T) {
	if _, err := exec.LookPath("fd"); err != nil {
		t.Skipf("fd not installed")
	}

	originalFindFileSupport := common.Config.FindFileSupport
	defer func() {
		common.Config.FindFileSupport = originalFindFileSupport
	}()

	curTestDir := filepath.Join(testDir, "TestFind")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	nestedDir := filepath.Join(dir2, "nested")
	utils.SetupDirectories(t, curTestDir, dir1, dir2, nestedDir)
	utils.SetupFiles(t,
		filepath.Join(dir1, "file1.txt"),
		filepath.Join(dir2, "topfile.txt"),
		filepath.Join(nestedDir, "target_file.txt"))
	t.Cleanup(func() {
		os.RemoveAll(curTestDir)
	})

	t.Run("Find file navigation", func(t *testing.T) {
		p := setupProgForFind(t, curTestDir)
		openFind(t, p)

		p.SendKey("target")
		targetFile := filepath.Join(nestedDir, "target_file.txt")
		assert.Eventually(t, func() bool {
			results := p.getModel().findModal.GetResults()
			return len(results) == 1 && results[0].Path == targetFile && !results[0].IsDir
		}, DefaultTestTimeout, DefaultTestTick, "target_file.txt should be found by find")

		// Press enter to navigate to the file's directory
		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			panel := p.getModel().getFocusedFilePanel()
			return !p.getModel().findModal.IsOpen() &&
				panel.Location == nestedDir &&
				panel.TargetFile == "" &&
				panel.GetFocusedItem().Name == "target_file.txt"
		}, DefaultTestTimeout, DefaultTestTick,
			"Find modal should close and navigate to %s (current location: %s)",
			nestedDir, p.getModel().getFocusedFilePanel().Location)
	})

	t.Run("Find folder navigation", func(t *testing.T) {
		p := setupProgForFind(t, curTestDir)
		openFind(t, p)

		p.SendKey("nested")
		assert.Eventually(t, func() bool {
			results := p.getModel().findModal.GetResults()
			return len(results) == 1 && results[0].Path == nestedDir && results[0].IsDir
		}, DefaultTestTimeout, DefaultTestTick, "nested dir should be found by find")

		// Press enter to navigate to the directory
		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			return !p.getModel().findModal.IsOpen() &&
				p.getModel().getFocusedFilePanel().Location == nestedDir
		}, DefaultTestTimeout, DefaultTestTick,
			"Find modal should close and navigate to %s (current location: %s)",
			nestedDir, p.getModel().getFocusedFilePanel().Location)
	})

	t.Run("Find disabled shows no results", func(t *testing.T) {
		common.Config.FindFileSupport = false
		m := defaultTestModel(dir1)

		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.FindFile[0]))
		assert.True(t, m.findModal.IsOpen(), "Find modal should open even when FindFileSupport is disabled")

		results := m.findModal.GetResults()
		assert.Empty(t, results, "Find modal should show no results when FindFileSupport is disabled")
	})

	t.Run("Find 'ctrl+f' key suppression on open", func(t *testing.T) {
		p := setupProgForFind(t, curTestDir)
		openFind(t, p)
		assert.Empty(t, p.getModel().findModal.GetTextInputValue(),
			"The 'ctrl+f' key should not be added to textInput")
		p.SendKeyDirectly("abc")
		assert.Equal(t, "abc", p.getModel().findModal.GetTextInputValue())
	})

	t.Run("Search filter cleared on file selection", func(t *testing.T) {
		target2File := filepath.Join(dir2, "target2.txt")
		utils.SetupFiles(t, target2File)
		t.Cleanup(func() {
			os.Remove(target2File)
		})

		p := setupProgForFind(t, dir2)
		p.getModel().getFocusedFilePanel().SearchBar.SetValue("topfile")

		openFind(t, p)

		p.SendKey("target2")
		assert.Eventually(t, func() bool {
			results := p.getModel().findModal.GetResults()
			return len(results) == 1 && results[0].Path == target2File && !results[0].IsDir
		}, DefaultTestTimeout, DefaultTestTick, "target2.txt should be found by find")

		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			panel := p.getModel().getFocusedFilePanel()
			return !p.getModel().findModal.IsOpen() &&
				panel.Location == dir2 &&
				panel.SearchBar.Value() == "" &&
				panel.TargetFile == "" &&
				panel.GetFocusedItem().Name == "target2.txt"
		}, DefaultTestTimeout, DefaultTestTick,
			"Search filter should be cleared and cursor should move to %s (search bar: %s)",
			target2File, p.getModel().getFocusedFilePanel().SearchBar.Value())
	})
}
