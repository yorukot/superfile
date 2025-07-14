package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
	"github.com/yorukot/superfile/src/internal/utils"
)

func TestModel_Update_Prompt(t *testing.T) {
	curTestDir := filepath.Join(testDir, "TestPrompt")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")

	setupDirectories(t, curTestDir, dir1, dir2)
	setupFiles(t, file1)
	t.Cleanup(func() {
		os.RemoveAll(curTestDir)
	})

	// We want to test these. TODO : complete important tests
	// 1. Being able to open prompt
	// 1a. Open in shell mode, 1b. Open in prompt mode 1c. Switching between then

	// 2. Being able to execute shell commands
	// 3. Shell command failure is handled and prompt stays open
	// 4. Successful Model actions - Split, Cd, Open new panel
	// 4a. Working split
	// 4b. Working cd : cd to abs path, cd to relative path, cd to home
	// 4c. Working open : open to abs path, open to relative path, open to home
	// 5. Split - Failure due to reaching max no. of panels
	// 6. cd - failure due to invalid path
	// 7. open - failure due to reaching max no. of panels
	// 8. open - failure due to invalid path
	// 9. cd and open - handling absolute and relative paths correctly
	// 10. Model closing
	// 10a. Pressing escape or ctrl+c and model closes
	// 10b. Autoclose based on config

	// Dont test shell command substitution here.

	// We might want to wrap os command execution in an interface and
	// ? Use a mock os command executor to have timeouts, and
	// custom command behaviour

	// Other tests cases
	// -- UI
	// 1. Entire model's rendering with promptModel open/closed
	// 2. Rendering not breaking when user pastes/enter special character or too much text
	// 3. Prompt gets resized based on total screen size. And always fits in

	// -- Functionality
	// 1. Shell command Timeout. Testing timeout is a pain. We should use async, and configure low timeout
	// like 1 sec for testing
	// 2. In case we plan to show output, we need to test case of
	// too big Shell command output

	testBasicPromptFunctionality(t, dir1)
	testPanelOperations(t, dir1, dir2, curTestDir)
	testDirectoryHandlingWithQuotes(t, curTestDir, dir1)
	testShellCommandsWithQuotes(t, curTestDir, dir1)
}

// testBasicPromptFunctionality tests opening, closing and basic command execution
func testBasicPromptFunctionality(t *testing.T, dir1 string) {
	t.Run("Basic Prompt Opening", func(t *testing.T) {
		m := defaultTestModel(dir1)
		assert.False(t, m.promptModal.IsOpen())
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))
		assert.True(t, m.promptModal.IsOpen())
		assert.True(t, m.promptModal.IsShellMode())

		// Switching between modes
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		assert.False(t, m.promptModal.IsShellMode(), "Pressing prompt key should switch to prompt mode")
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))
		assert.True(t, m.promptModal.IsShellMode(), "Pressing shell key should switch to shell mode")

		// Closing and opening in prompt mode
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.CancelTyping[0]))
		assert.False(t, m.promptModal.IsOpen())
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		assert.True(t, m.promptModal.IsOpen())
		assert.False(t, m.promptModal.IsShellMode())
	})

	t.Run("Shell command execution", func(t *testing.T) {
		m := defaultTestModel(dir1)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))
		// Prefer cross platform command
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg("mkdir test_dir"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded())
		assert.DirExists(t, filepath.Join(dir1, "test_dir"))

		// Invalid command shouldn't cause issues.
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg("xyz_non_exisiting_command"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded())
		assert.True(t, m.promptModal.IsOpen())
	})

	t.Run("Model closing", func(t *testing.T) {
		m := defaultTestModel(dir1)
		for _, key := range common.Hotkeys.CancelTyping {
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
			assert.True(t, m.promptModal.IsOpen())
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(key))
			assert.False(t, m.promptModal.IsOpen(), "Prompt should get closed")
		}
	})
}

// testPanelOperations tests split, cd, and open panel operations
func testPanelOperations(t *testing.T, dir1, dir2, curTestDir string) {
	t.Run("Split Panel", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		require.True(t, m.promptModal.IsOpen())
		for len(m.fileModel.filePanels) < m.fileModel.maxFilePanel {
			prevCnt := len(m.fileModel.filePanels)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.SplitCommand))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			require.Len(t, m.fileModel.filePanels, prevCnt+1)
			assert.Equal(t, dir1, m.fileModel.filePanels[prevCnt].location)
			assert.True(t, m.promptModal.LastActionSucceeded())
		}

		// Now doing a split should fail
		prevCnt := len(m.fileModel.filePanels)
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.SplitCommand))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded())
		assert.Len(t, m.fileModel.filePanels, prevCnt)
	})

	t.Run("cd Panel", func(t *testing.T) {
		m := defaultTestModel(dir1)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+" "+dir2))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "cd using absolute path should work")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+" .."))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "cd using relative path should work")
		assert.Equal(t, curTestDir, m.getFocusedFilePanel().location)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+" "+filepath.Base(dir2)))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "cd using relative path should work")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+" "+filepath.Join(dir2, "non_existing_dir")))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded(), "cd invalid abs path should not work")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+" non_existing_dir"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded(), "cd invalid relative path should not work")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+" ~"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "cd using tilde should work")
		assert.Equal(t, xdg.Home, m.getFocusedFilePanel().location)
	})

	t.Run("open Panel", func(t *testing.T) {
		m := defaultTestModel(dir1)
		orgCnt := len(m.fileModel.filePanels)
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" "+dir2))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "open using absolute path should work")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location)

		m.closeFilePanel()
		assert.Len(t, m.fileModel.filePanels, orgCnt)
		assert.Equal(t, dir1, m.getFocusedFilePanel().location)

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" ../dir2"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "open using relative path should work")
		assert.Equal(t, dir2, m.getFocusedFilePanel().location)

		m.closeFilePanel()

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" ~"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "open using tilde should work")
		assert.Equal(t, xdg.Home, m.getFocusedFilePanel().location)

		m.closeFilePanel()

		userHomeEnv := "HOME"
		if runtime.GOOS == utils.OsWindows {
			userHomeEnv = "USERPROFILE"
		}
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+fmt.Sprintf(" ${%s}", userHomeEnv)))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "open using variable substitution should work")
		assert.Equal(t, xdg.Home, m.getFocusedFilePanel().location)

		m.closeFilePanel()

		// Note : resolving shell subsitution is flaky in windows.
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" $(echo \"~\")"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.True(t, m.promptModal.LastActionSucceeded(), "open using command substitution should work")
		assert.Equal(t, xdg.Home, m.getFocusedFilePanel().location)

		m.closeFilePanel()

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" non_existing_dir"))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded(), "open using invalid relative path should not work")

		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" "+filepath.Join(dir2, "non_existing_dir")))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded(), "open using invalid abs path should not work")

		for len(m.fileModel.filePanels) < m.fileModel.maxFilePanel {
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" ."))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded())
		}

		// Now doing a open should fail
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+" ."))
		TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
		assert.False(t, m.promptModal.LastActionSucceeded())
	})
}

// testDirectoryHandlingWithQuotes tests handling directories with spaces and quotes
func testDirectoryHandlingWithQuotes(t *testing.T, curTestDir, dir1 string) {
	t.Run("Directory names with spaces and quotes", func(t *testing.T) {
		// Create test directories with spaces and special characters
		dirWithSpaces := filepath.Join(curTestDir, "dir with spaces")
		dirWithQuotes := filepath.Join(curTestDir, "dir'with'quotes")

		// Windows doesn't allow double quotes in directory names
		var dirWithSpecialChars, dirWithMixed string
		var directoriesToCreate []string

		if runtime.GOOS == "windows" {
			// On Windows, use alternative characters that don't conflict with filesystem restrictions
			dirWithSpecialChars = filepath.Join(curTestDir, `dir[with]quotes`)
			dirWithMixed = filepath.Join(curTestDir, `dir with 'mixed' [quotes]`)
			directoriesToCreate = []string{dirWithSpaces, dirWithQuotes, dirWithSpecialChars, dirWithMixed}
		} else {
			// On Unix-like systems, double quotes are allowed in directory names
			dirWithSpecialChars = filepath.Join(curTestDir, `dir"with"quotes`)
			dirWithMixed = filepath.Join(curTestDir, `dir with 'mixed' "quotes"`)
			directoriesToCreate = []string{dirWithSpaces, dirWithQuotes, dirWithSpecialChars, dirWithMixed}
		}

		setupDirectories(t, directoriesToCreate...)

		t.Run("cd with double quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` "`+dirWithSpaces+`"`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with double quotes should work")
			assert.Equal(t, dirWithSpaces, m.getFocusedFilePanel().location)
		})

		t.Run("cd with single quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` '`+dirWithSpaces+`'`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with single quotes should work")
			assert.Equal(t, dirWithSpaces, m.getFocusedFilePanel().location)
		})

		t.Run("cd with single quotes in path using double quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` "`+dirWithQuotes+`"`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with single quotes in path should work")
			assert.Equal(t, dirWithQuotes, m.getFocusedFilePanel().location)
		})

		t.Run("cd with double quotes in path using single quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` '`+dirWithSpecialChars+`'`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with double quotes in path should work")
			assert.Equal(t, dirWithSpecialChars, m.getFocusedFilePanel().location)
		})

		t.Run("cd with escaped spaces", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` `+strings.ReplaceAll(dirWithSpaces, " ", `\ `)))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with escaped spaces should work")
			assert.Equal(t, dirWithSpaces, m.getFocusedFilePanel().location)
		})

		t.Run("open with double quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.OpenCommand+` "`+dirWithSpaces+`"`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "open with double quotes should work")
			assert.Equal(t, dirWithSpaces, m.getFocusedFilePanel().location)

			m.closeFilePanel()
		})

		t.Run("cd with quoted environment variable", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			userHomeEnv := "HOME"
			if runtime.GOOS == utils.OsWindows {
				userHomeEnv = "USERPROFILE"
			}

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` "${`+userHomeEnv+`}"`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with quoted env var should work")
			assert.Equal(t, xdg.Home, m.getFocusedFilePanel().location)
		})

		t.Run("cd with single quoted environment variable", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))

			userHomeEnv := "HOME"
			if runtime.GOOS == utils.OsWindows {
				userHomeEnv = "USERPROFILE"
			}

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(prompt.CdCommand+` '${`+userHomeEnv+`}'`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "cd with single quoted env var works in superfile (unlike bash)")
			assert.Equal(t, xdg.Home, m.getFocusedFilePanel().location)
		})
	})
}

// testShellCommandsWithQuotes tests shell command execution with quoted arguments
func testShellCommandsWithQuotes(t *testing.T, curTestDir, dir1 string) {
	t.Run("Shell command with quotes", func(t *testing.T) {
		dirWithSpaces := filepath.Join(curTestDir, "test dir with spaces")
		setupDirectories(t, dirWithSpaces)

		t.Run("shell command with double quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(`mkdir "`+filepath.Join(dir1, "new dir with spaces")+`"`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "shell command with quotes should work")
			assert.DirExists(t, filepath.Join(dir1, "new dir with spaces"))
		})

		t.Run("shell command with single quotes", func(t *testing.T) {
			m := defaultTestModel(dir1)
			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(`mkdir '`+filepath.Join(dir1, "another dir with spaces")+`'`))
			TeaUpdateWithErrCheck(t, &m, tea.KeyMsg{Type: tea.KeyEnter})
			assert.True(t, m.promptModal.LastActionSucceeded(), "shell command with single quotes should work")
			assert.DirExists(t, filepath.Join(dir1, "another dir with spaces"))
		})
	})
}
