package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

/*
The purpose of this test file is to have the
(1) common global data for tests
(2) common setup for tests, and cleanup
(3) Basic model fuctionality tests
    - Initialization
	- Resize
	- Update
	- Quitting
*/

// Helps to have centralized cleanup
var testDir string //nolint: gochecknoglobals // One-time initialized, and then read-only global test variable

func cleanupTestDir() {
	err := os.RemoveAll(testDir)
	if err != nil {
		fmt.Printf("error while cleaning up test directory, err : %v", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	err := common.PopulateGlobalConfigs()
	if err != nil {
		fmt.Printf("error while populating config, err : %v", err)
		os.Exit(1)
	}

	// A cleanup before is required in case the previous test run had a panic, and then
	// deferred cleanup never executed

	// Create testDir
	testDir = filepath.Join(os.TempDir(), "spf_testdir")
	cleanupTestDir()
	if err := os.Mkdir(testDir, 0o755); err != nil {
		fmt.Printf("error while creating test directory, err : %v", err)
		os.Exit(1)
	}
	defer cleanupTestDir()

	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	m.Run()
	// Maybe catch panic
}

func TestBasic(t *testing.T) {
	curTestDir := filepath.Join(testDir, "TestBasic")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")

	t.Run("Basic Checks", func(t *testing.T) {
		utils.SetupDirectories(t, curTestDir, dir1, dir2)
		utils.SetupFiles(t, file1)
		t.Cleanup(func() {
			os.RemoveAll(curTestDir)
		})

		m := defaultTestModel(dir1)

		// Validate the most of the data stored in model object
		// Inspect model struct to see what more can be validated.
		// 1 - File panel location, cursor, render index, etc.
		// 2 - Directory Items are listed
		// 3 - sidebar items pinned items are listed
		// 4 - process panel is empty
		// 5 - clipboard is empty
		// 6 - model's dimenstion

		assert.Equal(t, dir1, m.getFocusedFilePanel().Location)
	})
}

func TestInitialFilePathPositionsCursorWindow(t *testing.T) {
	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")

	utils.SetupDirectories(t, curTestDir, dir1)

	var file7 string
	var file2 string
	for i := range 10 {
		f := filepath.Join(dir1, fmt.Sprintf("file%d.txt", i))
		utils.SetupFiles(t, f)
		if i == 7 {
			file7 = f
		}
		if i == 2 {
			file2 = f
		}
	}

	m := defaultTestModel(dir1, file2, file7)
	// View port of 5
	TeaUpdate(m, tea.WindowSizeMsg{Width: common.MinimumWidth, Height: 10})
	// Uncomment below to understand the distribution
	// t.Logf("Heights : %d [%d - [%d] %d]\n", m.fullHeight, m.footerHeight, m.mainPanelHeight,
	//	panelElementHeight(m.mainPanelHeight))
	require.Len(t, m.fileModel.FilePanels, 3)
	assert.Equal(t, dir1, m.fileModel.FilePanels[0].Location)
	assert.Equal(t, file2, m.fileModel.FilePanels[1].GetFocusedItem().Location)
	assert.Equal(t, 2, m.fileModel.FilePanels[1].GetCursor())
	assert.Equal(t, 0, m.fileModel.FilePanels[1].GetRenderIndex())
	assert.Equal(t, file7, m.fileModel.FilePanels[2].GetFocusedItem().Location)
	assert.Equal(t, 7, m.fileModel.FilePanels[2].GetCursor())
	assert.Equal(t, 3, m.fileModel.FilePanels[2].GetRenderIndex())
}

func TestQuit(t *testing.T) {
	// Test
	// 1 - Normal quit
	// 2 - Normal quit with running process causing a warn modal
	//     2a - Cancelling quit
	//     2b - Proceeding with the quit
	// 3 - Cd on quit test that LastDir is written on

	t.Run("Normal Quit", func(t *testing.T) {
		m := defaultTestModel(testDir)
		assert.Equal(t, notQuitting, m.modelQuitState)
		cmd := TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.Quit[0]))
		assert.Equal(t, quitDone, m.modelQuitState)
		assert.True(t, IsTeaQuit(cmd))
	})
	t.Run("Quit with running process", func(t *testing.T) {
		m := defaultTestModel(testDir)
		m.processBarModel.AddOrUpdateProcess(processbar.Process{
			State: processbar.InOperation,
			Done:  0,
			Total: 100,
			ID:    "1",
		})

		assert.Equal(t, notQuitting, m.modelQuitState)
		cmd := TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.Quit[0]))
		assert.Equal(t, quitConfirmationInitiated, m.modelQuitState)
		assert.False(t, IsTeaQuit(cmd))

		// Now we would be asked for confirmation.
		// Cancel the quit
		cmd = TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.CancelTyping[0]))
		assert.Equal(t, notQuitting, m.modelQuitState)
		assert.False(t, IsTeaQuit(cmd))

		// Again trigger quit
		cmd = TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.Quit[0]))
		assert.Equal(t, quitConfirmationInitiated, m.modelQuitState)
		assert.False(t, IsTeaQuit(cmd))

		// Confirm this time
		cmd = TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.Confirm[0]))
		assert.Equal(t, quitDone, m.modelQuitState)
		assert.True(t, IsTeaQuit(cmd))
	})

	t.Run("Cd on quit test that LastDir is written on", func(t *testing.T) {
		lastDirFile := filepath.Join(variable.SuperFileStateDir, "lastdir")
		require.NoError(t, os.MkdirAll(filepath.Dir(lastDirFile), 0o755))
		m := defaultTestModel(testDir)

		assert.Equal(t, notQuitting, m.modelQuitState)

		cmd := TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.CdQuit[0]))

		assert.Equal(t, quitDone, m.modelQuitState)
		assert.True(t, IsTeaQuit(cmd))

		data, err := os.ReadFile(lastDirFile)
		require.NoError(t, err)
		assert.Equal(t, "cd '"+testDir+"'", string(data), "LastDir file should contain the tempDir path")

		err = os.Remove(lastDirFile)
		require.NoError(t, err)
	})
}

func TestChooserFile(t *testing.T) {
	// 1 - No quit - blank chooser file
	// 2 - Quit with valid chooser file
	//     2a - file preview
	//     2b - directory preview
	// 3 - No quit - invalid chooser file
	curTestDir := filepath.Join(testDir, "TestChooserFile")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")
	testChooserFile := filepath.Join(dir2, "chooser_file.txt")
	utils.SetupDirectories(t, curTestDir, dir1, dir2)
	utils.SetupFiles(t, file1)

	testdata := []struct {
		name            string
		chooserFile     string
		hotkey          string
		expectedQuit    bool
		expectedContent string
	}{
		{
			name:            "Open with default app with valid chooser file",
			chooserFile:     testChooserFile,
			hotkey:          common.Hotkeys.Confirm[0],
			expectedQuit:    true,
			expectedContent: file1,
		},
		{
			name:            "Open with file editor with valid chooser file",
			chooserFile:     testChooserFile,
			hotkey:          common.Hotkeys.OpenFileWithEditor[0],
			expectedQuit:    true,
			expectedContent: file1,
		},
		{
			name:            "Open with directory editor valid chooser file",
			hotkey:          common.Hotkeys.OpenCurrentDirectoryWithEditor[0],
			chooserFile:     testChooserFile,
			expectedQuit:    true,
			expectedContent: dir1,
		},
		{
			name:            "Open with file editor with Blank chooser file",
			chooserFile:     "",
			hotkey:          common.Hotkeys.OpenFileWithEditor[0],
			expectedQuit:    false,
			expectedContent: "",
		},
		{
			name:            "Open with file editor with Invalid chooser file",
			chooserFile:     filepath.Join(curTestDir, "non_existent_dir", "file.txt"),
			hotkey:          common.Hotkeys.OpenFileWithEditor[0],
			expectedQuit:    false,
			expectedContent: "",
		},
	}

	// Must be sequential as we are using global variable chooserfile
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := defaultTestModel(dir1)
			if tt.expectedQuit {
				err := os.WriteFile(tt.chooserFile, []byte{}, 0o644)
				require.NoError(t, err)
			}
			variable.SetChooserFile(tt.chooserFile)
			cmd := TeaUpdate(m, utils.TeaRuneKeyMsg(tt.hotkey))

			if tt.expectedQuit {
				assert.Equal(t, quitDone, m.modelQuitState)
				assert.True(t, IsTeaQuit(cmd))
				assert.FileExists(t, tt.chooserFile)
				data, err := os.ReadFile(tt.chooserFile)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedContent, string(data))
			} else {
				assert.Equal(t, notQuitting, m.modelQuitState)
				assert.False(t, IsTeaQuit(cmd))
			}
		})
	}
}

func eventuallyEnsurePreviewContent(t *testing.T, m *model, content string, msgAndArgs ...any) {
	contains := false
	assert.Eventually(t, func() bool {
		contains = strings.Contains(m.fileModel.FilePreview.GetContent(), content)
		return contains
	}, DefaultTestTimeout, DefaultTestTick, msgAndArgs...)
	if !contains {
		pContent := ansi.Strip(m.fileModel.FilePreview.GetContent())
		pContent = pContent[:min(len(pContent), 20)]
		t.Logf("%s was not found in '%s'", content, pContent)
	}
}

func TestAsyncPreviewPanelSync(t *testing.T) {
	curTestDir := t.TempDir()

	originalPreviewWidth := common.Config.FilePreviewWidth
	common.Config.FilePreviewWidth = 0
	t.Cleanup(func() {
		common.Config.FilePreviewWidth = originalPreviewWidth
	})

	file1, content1 := filepath.Join(curTestDir, "file1.txt"), "File 1 content"
	file2, content2 := filepath.Join(curTestDir, "file2.txt"), "File 2 content"
	utils.SetupFilesWithData(t, []byte(content1), file1)
	utils.SetupFilesWithData(t, []byte(content2), file2)

	m := defaultTestModelWithFilePreview(curTestDir)
	p := NewTestTeaProgWithEventLoop(t, m)

	// We need to send message via event loop to ensure that preview load command
	// is actually processed, also we want a size bigger than default
	// to allow more number of panels
	p.Send(tea.WindowSizeMsg{Width: 4 * DefaultTestModelWidth, Height: 4 * DefaultTestModelHeight})

	eventuallyEnsurePreviewContent(t, m, content1, "file1 content should load initially")
	pW := m.fileModel.FilePreview.GetContentWidth()

	// Create two panels
	splitPanelAsync(p)
	splitPanelAsync(p)
	eventuallyEnsurePreviewContent(t, m, content1, "file1 content should reload after new panel")

	assert.NotEqual(t, pW, m.fileModel.FilePreview.GetContentWidth(),
		"width should change on new panel creation")

	p.Send(tea.KeyMsg{Type: tea.KeyDown})
	t.Logf("Current element : %s", m.getFocusedFilePanel().GetFocusedItem().Location)
	eventuallyEnsurePreviewContent(t, m, content2, "content should update to file2")

	p.SendKey(common.Hotkeys.CloseFilePanel[0])
	eventuallyEnsurePreviewContent(t, m, content1, "content should update to file1 after closing panel")

	// Upscale
	p.Send(tea.WindowSizeMsg{Width: 8 * DefaultTestModelWidth,
		Height: 8 * DefaultTestModelHeight})
	eventuallyEnsurePreviewContent(t, m, content1, "content should update to file1 after resize")

	// Downscale
	p.Send(tea.WindowSizeMsg{Width: 6 * DefaultTestModelWidth,
		Height: 6 * DefaultTestModelHeight})
	eventuallyEnsurePreviewContent(t, m, content1, "content should update to file1 after resize")
}
