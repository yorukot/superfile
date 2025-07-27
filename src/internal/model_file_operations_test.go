package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		assert.Len(t, p.getModel().processBarModel.process, 2)
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
