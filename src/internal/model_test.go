package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/utils"
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
	if err := os.Mkdir(testDir, 0755); err != nil {
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

		assert.Equal(t, dir1, m.getFocusedFilePanel().location)
	})
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
		cmd := TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.Quit[0]))
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
		cmd := TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.Quit[0]))
		assert.Equal(t, quitConfirmationInitiated, m.modelQuitState)
		assert.False(t, IsTeaQuit(cmd))

		// Now we would be asked for confirmation.
		// Cancel the quit
		cmd = TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.CancelTyping[0]))
		assert.Equal(t, notQuitting, m.modelQuitState)
		assert.False(t, IsTeaQuit(cmd))

		// Again trigger quit
		cmd = TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.Quit[0]))
		assert.Equal(t, quitConfirmationInitiated, m.modelQuitState)
		assert.False(t, IsTeaQuit(cmd))

		// Confirm this time
		cmd = TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.Confirm[0]))
		assert.Equal(t, quitDone, m.modelQuitState)
		assert.True(t, IsTeaQuit(cmd))
	})

	t.Run("Cd on quit test that LastDir is written on", func(t *testing.T) {
		lastDirFile := filepath.Join(variable.SuperFileStateDir, "lastdir")
		require.NoError(t, os.MkdirAll(filepath.Dir(lastDirFile), 0o755))
		m := defaultTestModel(testDir)

		assert.Equal(t, notQuitting, m.modelQuitState)

		cmd := TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.CdQuit[0]))

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
				err := os.WriteFile(tt.chooserFile, []byte{}, 0644)
				require.NoError(t, err)
			}
			variable.SetChooserFile(tt.chooserFile)
			cmd := TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(tt.hotkey))

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

		err := m.updateCurrentFilePanelDir(dir2)
		require.NoError(t, err, "Failed to navigate to dir2")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location, "Should be in dir2 after navigation")

		err = m.updateCurrentFilePanelDir(dir3)
		require.NoError(t, err, "Failed to navigate to dir3")
		assert.Equal(t, dir3, m.getFocusedFilePanel().location, "Should be in dir3 after navigation")

		TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenZoxide[0]))
		assert.True(t, m.zoxideModal.IsOpen(), "Zoxide modal should open when pressing 'z' key")

		// Type "dir2" to search for it
		for _, char := range "dir2" {
			TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(string(char)))
		}

		results := m.zoxideModal.GetResults()
		assert.GreaterOrEqual(t, len(results), 1, "Should have at least 1 directory found by zoxide UI search")

		resultPaths := make([]string, len(results))
		for i, result := range results {
			resultPaths[i] = result.Path
		}
		assert.Contains(t, resultPaths, dir2, "dir2 should be found by zoxide UI search")

		// Press enter to navigate to dir2
		TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.ConfirmTyping[0]))
		assert.False(t, m.zoxideModal.IsOpen(), "Zoxide modal should close after navigation")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location, "Should navigate back to dir2 after zoxide selection")
	})

	t.Run("Zoxide disabled shows no results", func(t *testing.T) {
		common.Config.ZoxideSupport = false
		m := defaultTestModelWithZClient(zClient, dir1)

		TeaUpdateWithErrCheck(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenZoxide[0]))
		assert.True(t, m.zoxideModal.IsOpen(), "Zoxide modal should open even when ZoxideSupport is disabled")

		results := m.zoxideModal.GetResults()
		assert.Empty(t, results, "Zoxide modal should show no results when ZoxideSupport is disabled")
	})
}
