package internal

import (
	"fmt"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
)

var SampleDataBytes = []byte("This is sample") //nolint: gochecknoglobals // Effectively const

func defaultTestModel(dirs ...string) model {
	m := defaultModelConfig(false, false, false, dirs)
	_, _ = TeaUpdate(&m, tea.WindowSizeMsg{Width: 2 * common.MinimumWidth, Height: 2 * common.MinimumHeight})
	return m
}

func setupDirectories(t *testing.T, dirs ...string) {
	t.Helper()
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
	}
}

func setupFilesWithData(t *testing.T, data []byte, files ...string) {
	t.Helper()
	for _, file := range files {
		err := os.WriteFile(file, data, 0644)
		require.NoError(t, err)
	}
}

func setupFiles(t *testing.T, files ...string) {
	setupFilesWithData(t, SampleDataBytes, files...)
}

// TeaUpdate : Utility to send update to model , majorly used in tests
// Not using pointer receiver as this is more like a utility, than
// a member function of model
// Todo : Consider wrapping TeaUpdate with a helper that both forwards the return
// values and does a require.NoError(t, err)
func TeaUpdate(m *model, msg tea.Msg) (tea.Cmd, error) {
	resModel, cmd := m.Update(msg)

	mObj, ok := resModel.(model)
	if !ok {
		return cmd, fmt.Errorf("unexpected model type: %T", resModel)
	}
	*m = mObj
	return cmd, nil
}

func TeaUpdateWithErrCheck(t *testing.T, m *model, msg tea.Msg) tea.Cmd {
	cmd, err := TeaUpdate(m, msg)
	require.NoError(t, err)
	return cmd
}

// Is the command tea.quit, or a batch that contains tea.quit
func IsTeaQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	switch msg := msg.(type) {
	case tea.QuitMsg:
		return true
	case tea.BatchMsg:
		for _, curCmd := range msg {
			if IsTeaQuit(curCmd) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
