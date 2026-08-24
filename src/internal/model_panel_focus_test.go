package internal

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// #1497: opening a new panel while the sidebar (or process bar / metadata) was
// focused left both focused at once. Focus should end up on the new panel only.
func TestNewPanelClearsNonFilePanelFocus(t *testing.T) {
	curTestDir := t.TempDir()
	subDir := filepath.Join(curTestDir, "sub")
	utils.SetupDirectories(t, subDir)

	// Every way of opening a panel should land focus on it.
	createPaths := []struct {
		name   string
		create func(m *model) error
	}{
		{
			name: "split panel hotkey",
			create: func(m *model) error {
				_, err := m.splitPanel()
				return err
			},
		},
		{
			name: "create new file panel at home",
			create: func(m *model) error {
				_, err := m.createNewFilePanel(variable.HomeDir)
				return err
			},
		},
		{
			name: "open panel relative to current",
			create: func(m *model) error {
				_, err := m.createNewFilePanelRelativeToCurrent("sub")
				return err
			},
		},
	}

	for _, cp := range createPaths {
		t.Run("sidebar focused / "+cp.name, func(t *testing.T) {
			m := defaultTestModel(curTestDir)
			m.focusOnSideBar()
			require.Equal(t, sidebarFocus, m.focusPanel, "precondition: sidebar focused")
			require.False(t, m.getFocusedFilePanel().IsFocused, "precondition: file panel not focused")

			require.NoError(t, cp.create(m))

			assert.Equal(t, nonePanelFocus, m.focusPanel,
				"sidebar must lose focus after a new file panel is created")
			assert.True(t, m.getFocusedFilePanel().IsFocused,
				"the newly created file panel must be focused")
		})
	}

	// Same deal for the process bar (only focusable when the footer is on).
	t.Run("process bar focused / split panel", func(t *testing.T) {
		m := defaultTestModelWithFooterAndFilePreview(curTestDir)
		m.focusOnProcessBar()
		require.Equal(t, processBarFocus, m.focusPanel, "precondition: process bar focused")

		_, err := m.splitPanel()
		require.NoError(t, err)

		assert.Equal(t, nonePanelFocus, m.focusPanel,
			"process bar must lose focus after a new file panel is created")
		assert.True(t, m.getFocusedFilePanel().IsFocused,
			"the newly created file panel must be focused")
	})
}
