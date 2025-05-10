package internal

import (
	"path/filepath"
	"testing"

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

	setupDirectories(t,  curTestDir, dir1, dir2)
	setupFiles(t, file1)



	// We want to test these. Todo : complete important tests
	// 1. Being able to open prompt
	// 1a. Open in shell mode, 1b. Open in prompt mode
	
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

	t.Run("Basic Prompt Opening", func(t *testing.T) {
		m := defaultModelConfig(false, false, []string{curTestDir})
		firstUse = false
		assert.False(t, m.promptModal.IsOpen())
		_, err := TeaUpdate(&m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))
		require.NoError(t, err, "Opening the prompt should not produce an error")
		assert.True(t, m.promptModal.IsOpen())
		assert.True(t, m.promptModal.IsShellMode())
	})

	t.Run("Split Panel", func(t *testing.T) {
		m := defaultModelConfig(false, false, []string{"/"})
		firstUse = false
		assert.Len(t, m.fileModel.filePanels, 1)
		_, _ = TeaUpdate(&m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		_, _ = TeaUpdate(&m, utils.TeaRuneKeyMsg(prompt.SplitCommand))
		_, err := TeaUpdate(&m, tea.KeyMsg{Type: tea.KeyEnter})
		require.NoError(t, err, "Opening the prompt should not produce an error")
		assert.Len(t, m.fileModel.filePanels, 2)
	})
}
