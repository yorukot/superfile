package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

// TODO : Add test for model initialized with multiple directories
// TODO : Add test for clipboard different variations, cut paste
// TODO : Add test for tea resizing
// TODO : Add test for quitting

func TestCopy(t *testing.T) {
	curTestDir := filepath.Join(testDir, "TestCopy")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")
	t.Run("Basic Copy", func(t *testing.T) {
		setupDirectories(t, curTestDir, dir1, dir2)
		setupFiles(t, file1)
		t.Cleanup(func() {
			os.RemoveAll(curTestDir)
		})

		m := defaultTestModel(dir1)

		// TODO validate current panel is "dir1"
		// TODO : Move all basic validation to a separate test
		// Everything that doesn't have anything to do with copy paste

		// validate file1
		// TODO : improve the interface we use to interact with filepanel

		// TODO : file1.txt should not be duplicated

		// TODO : Having to send a random keypress to initiate model init.
		// Should not have to do that
		TeaUpdateWithErrCheck(t, &m, nil)

		assert.Equal(t, "file1.txt",
			m.fileModel.filePanels[m.filePanelFocusIndex].element[0].name)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.CopyItems[0]))

		// TODO : validate clipboard
		assert.False(t, m.copyItems.cut)
		assert.Equal(t, file1, m.copyItems.items[0])

		// move to dir2
		m.updateCurrentFilePanelDir("../dir2")
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Actual paste may take time, since its an os operations
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1.txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)

		// TODO : still on clipboard
		assert.False(t, m.copyItems.cut)
		assert.Equal(t, file1, m.copyItems.items[0])

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Actual paste may take time, since its an os operations
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1(1).txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)
		assert.FileExists(t, filepath.Join(dir2, "file1(1).txt"))
		// TODO : Also validate process bar having two processes.
	})
}

type IgnorerWriter struct{}

func (w IgnorerWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestCopy2(t *testing.T) {
	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")
	setupDirectories(t, curTestDir, dir1, dir2)
	setupFiles(t, file1)

	t.Run("Basic Copy", func(t *testing.T) {

		m := defaultTestModel(dir1)
		p := tea.NewProgram(m, tea.WithInput(nil), tea.WithOutput(IgnorerWriter{}))
		go p.Run()
		// TODO validate current panel is "dir1"
		// TODO : Move all basic validation to a separate test
		// Everything that doesn't have anything to do with copy paste

		// validate file1
		// TODO : improve the interface we use to interact with filepanel

		// TODO : file1.txt should not be duplicated

		// TODO : Having to send a random keypress to initiate model init.
		// Should not have to do that
		t.Log("Sending nil")

		p.Send(nil)
		//TeaUpdateWithErrCheck(t, &m, nil)
		t.Log("Sent nil")
		assert.Equal(t, "file1.txt",
			m.fileModel.filePanels[m.filePanelFocusIndex].element[0].name)

		p.Send(utils.TeaRuneKeyMsg(common.Hotkeys.CopyItems[0]))

		// TODO : validate clipboard
		assert.False(t, m.copyItems.cut)
		assert.Equal(t, file1, m.copyItems.items[0])

		// move to dir2
		m.updateCurrentFilePanelDir("../dir2")
		p.Send(utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Actual paste may take time, since its an os operations
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1.txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)

		// TODO : still on clipboard
		assert.False(t, m.copyItems.cut)
		assert.Equal(t, file1, m.copyItems.items[0])

		p.Send(utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Actual paste may take time, since its an os operations
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1(1).txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)
		assert.FileExists(t, filepath.Join(dir2, "file1(1).txt"))
		// TODO : Also validate process bar having two processes.
	})
}

func TestFileCreation(t *testing.T) {
	// TODO Also add directory creation test to this
	curTestDir := filepath.Join(testDir, "TestNaming")
	testParentDir := filepath.Join(curTestDir, "parentDir")
	testChildDir := filepath.Join(testParentDir, "childDir")

	setupDirectories(t, curTestDir, testParentDir, testChildDir)

	t.Cleanup(func() {
		os.RemoveAll(curTestDir)
	})

	testdata := []struct {
		name          string
		fileName      string
		expectedError bool
	}{
		{"valid name", "file.txt", false},
		{"invalid single dot", ".", true},
		{"invalid double dot", "..", true},
		{"invalid trailing slash-dot", fmt.Sprintf("test%c.", filepath.Separator), true},
		{"invalid trailing slash-dot-dot", fmt.Sprintf("test%c..", filepath.Separator), true},
		{"valid name with trailing .", "abc.", false},
	}

	for _, tt := range testdata {
		m := defaultTestModel(testChildDir)

		TeaUpdateWithErrCheck(t, &m, nil)
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.FilePanelItemCreate[0]))

		assert.Equal(t, "", m.typingModal.errorMesssage)

		m.typingModal.textInput.SetValue(tt.fileName)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.ConfirmTyping[0]))

		if tt.expectedError {
			assert.NotEqual(t, "", m.typingModal.errorMesssage, "expected an error for input: %q", tt.fileName)
		} else {
			assert.Empty(t, m.typingModal.errorMesssage, "expected an error for input: %q", tt.fileName)
			assert.FileExists(t, filepath.Join(testChildDir, tt.fileName), "expected file to be created: %q", tt.fileName)
		}
	}
}
