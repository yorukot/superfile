package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	spfConfigDir := filepath.Join(filepath.Dir(filepath.Dir(filename)),
		"superfile_config")

	err := common.PopulateGlobalConfigs(
		filepath.Join(spfConfigDir, "config.toml"),
		filepath.Join(spfConfigDir, "hotkeys.toml"),
		filepath.Join(spfConfigDir, "theme", "monokai.toml"))

	if err != nil {
		fmt.Printf("error while populating config, err : %v", err)
		os.Exit(1)
	}

	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	m.Run()
}

// Model is huge. Just one test file ain't enough

func TestModel_Update_Prompt(t *testing.T) {
	// We want to test these. Todo : complete important tests
	// 1. Being able to open prompt
	// 2. Being able to execute shell commands
	// 3. Shell command scenarios like failure (validate failure)
	// 4. Successful Model actions - Split, Cd, Open new panel
	// 5. Split - Failure due to reaching max no. of panels
	// 6. cd - failure due to invalid path
	// 7. open - failure due to reaching max no. of panels
	// 8. open - failure due to invalid path
	// 9. cd and open - handling absolute and relative paths correctly
	// We might want to wrap os command execution in an interface and
	// Use a mock os command executor to have timeouts, and
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
		m := defaultModelConfig(false, false, []string{"/"})
		firstUse = false
		_, err := TeaUpdate(&m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenCommandLine[0]))
		require.NoError(t, err, "Opening the prompt should not produce an error")
		assert.True(t, m.promptModal.IsOpen())
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
