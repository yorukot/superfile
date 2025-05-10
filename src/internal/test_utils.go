package internal

import (
	"fmt"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
)

var SampleDataBytes = []byte("This is sample") //nolint: gochecknoglobals // Effectively const

func defaultTestModel(dirs ...string) model {
	return defaultModelConfig(false, false, false, dirs)
}

func setupDirectories(t *testing.T, dirs ...string) {
	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		require.NoError(t, err)
	}
}

func setupFiles(t *testing.T, files ...string) {
	for _, file := range files {
		err := os.WriteFile(file, SampleDataBytes, 0755)
		require.NoError(t, err)
	}
}

// TeaUpdate : Utility to send update to model , majorly used in tests
// Not using pointer receiver as this is more like a utility, than
// a member function of model
func TeaUpdate(m *model, msg tea.Msg) (tea.Cmd, error) {
	resModel, cmd := m.Update(msg)

	mObj, ok := resModel.(model)
	if !ok {
		return cmd, fmt.Errorf("unexpected model type: %T", resModel)
	}
	*m = mObj
	return cmd, nil
}
