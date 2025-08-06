package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/notify"
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

		p := NewTestTeaProgWithEventLoop(t, defaultTestModel(dir1))

		require.Equal(t, "file1.txt",
			p.getModel().getFocusedFilePanel().element[0].name)
		p.SendKeyDirectly(common.Hotkeys.CopyItems[0])
		assert.False(t, p.getModel().copyItems.cut)
		assert.Equal(t, file1, p.getModel().copyItems.items[0])

		p.getModel().updateCurrentFilePanelDir("../dir2")
		p.SendKey(common.Hotkeys.PasteItems[0])

		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1.txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)

		assert.False(t, p.getModel().copyItems.cut)
		assert.Equal(t, file1, p.getModel().copyItems.items[0])

		p.SendKey(common.Hotkeys.PasteItems[0])
		assert.Eventually(t, func() bool {
			_, err := os.Lstat(filepath.Join(dir2, "file1(1).txt"))
			return err == nil
		}, time.Second, 10*time.Millisecond)
		assert.FileExists(t, filepath.Join(dir2, "file1(1).txt"))
		//TODO: Also verify if there are only 2 items in process bar
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

		TeaUpdateWithErrCheck(m, nil)
		TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.FilePanelItemCreate[0]))

		assert.Equal(t, "", m.typingModal.errorMesssage)

		m.typingModal.textInput.SetValue(tt.fileName)

		TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.ConfirmTyping[0]))

		if tt.expectedError {
			assert.NotEqual(t, "", m.typingModal.errorMesssage, "expected an error for input: %q", tt.fileName)
		} else {
			assert.Empty(t, m.typingModal.errorMesssage, "expected an error for input: %q", tt.fileName)
			assert.FileExists(t, filepath.Join(testChildDir, tt.fileName), "expected file to be created: %q", tt.fileName)
		}
	}
}

func TestFileRename(t *testing.T) {
	curTestDir := t.TempDir()
	file1 := filepath.Join(curTestDir, "file1.txt")
	file2 := filepath.Join(curTestDir, "file2.txt")
	file3 := filepath.Join(curTestDir, "file3.txt")

	setupFilesWithData(t, []byte("f1"), file1)
	setupFilesWithData(t, []byte("f2"), file2)
	setupFilesWithData(t, []byte("f3"), file3)

	file1New := filepath.Join(curTestDir, "file1_new.txt")

	t.Run("Basic rename", func(t *testing.T) {
		m := defaultTestModel(curTestDir)
		p := NewTestTeaProgWithEventLoop(t, m)
		idx := findItemIndexInPanelByLocation(m.getFocusedFilePanel(), file1)
		require.NotEqual(t, -1, idx, "%s should be found in panel", file1)
		m.getFocusedFilePanel().cursor = idx

		p.SendKey(common.Hotkeys.FilePanelItemRename[0])
		p.SendKey("_new")
		p.Send(tea.KeyMsg{Type: tea.KeyEnter})

		assert.Eventually(t, func() bool {
			_, err1 := os.Stat(file1)
			_, err1New := os.Stat(file1New)
			return err1New == nil && os.IsNotExist(err1)
		}, time.Second, 10*time.Millisecond, "File never got renamed")
	})

	t.Run("Rename confirmation for same name", func(t *testing.T) {
		actualTest := func(doRename bool) {
			m := defaultTestModel(curTestDir)
			p := NewTestTeaProgWithEventLoop(t, m)
			idx := findItemIndexInPanelByLocation(m.getFocusedFilePanel(), file3)
			require.NotEqual(t, -1, idx, "%s should be found in panel", file3)

			m.getFocusedFilePanel().cursor = idx

			p.SendKeyDirectly(common.Hotkeys.FilePanelItemRename[0])
			m.getFocusedFilePanel().rename.SetValue("file2.txt")
			p.Send(tea.KeyMsg{Type: tea.KeyEnter})

			// This will result in async
			assert.Eventually(t, func() bool {
				return m.notifyModel.IsOpen()
			}, time.Second, 10*time.Millisecond, "Notify modal never opened, filepanel items : %v",
				m.getFocusedFilePanel().element)

			assert.Equal(t, notify.New(true,
				common.SameRenameWarnTitle,
				common.SameRenameWarnContent,
				notify.RenameAction), m.notifyModel, "Notify model should be as expected")

			if doRename {
				p.Send(tea.KeyMsg{Type: tea.KeyEnter})
			} else {
				p.SendKey(common.Hotkeys.CancelTyping[0])
			}

			assert.Eventually(t, func() bool {
				_, err2 := os.Stat(file2)
				_, err3 := os.Stat(file3)
				f2Data, err := os.ReadFile(file2)
				require.NoError(t, err)
				if doRename {
					// f3 should be gone. f2 should have content of f3
					return os.IsNotExist(err3) && err2 == nil &&
						string(f2Data) == "f3"
				}
				return err2 == nil && err3 == nil
			}, time.Second, 10*time.Millisecond,
				"Rename should be done/not done appropriately, file : %v",
				m.getFocusedFilePanel().element)
		}

		actualTest(false)
		actualTest(true)
	})
}
