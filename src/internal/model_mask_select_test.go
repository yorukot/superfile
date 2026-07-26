package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Types the given mask into an open mask modal
func typeMask(m *model, mask string) {
	for _, r := range mask {
		TeaUpdate(m, utils.TeaRuneKeyMsg(string(r)))
	}
}

func TestModel_Update_SelectByMask(t *testing.T) {
	curTestDir := filepath.Join(testDir, "TestSelectByMask")
	dir1 := filepath.Join(curTestDir, "dir1")
	file1 := filepath.Join(dir1, "file1.txt")
	file2 := filepath.Join(dir1, "file2.txt")
	file3 := filepath.Join(dir1, "image.png")

	utils.SetupDirectories(t, curTestDir, dir1)
	utils.SetupFiles(t, file1, file2, file3)
	t.Cleanup(func() {
		os.RemoveAll(curTestDir)
	})

	selectKey := common.Hotkeys.FilePanelSelectItemsByMask[0]
	unselectKey := common.Hotkeys.FilePanelUnselectItemsByMask[0]
	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}

	t.Run("Select by mask in select mode", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ChangePanelMode[0]))
		require.Equal(t, filepanel.SelectMode, m.getFocusedFilePanel().PanelMode)

		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		require.True(t, m.maskModal.open)
		assert.True(t, m.maskModal.selecting)
		// The hotkey that opened the modal must not land in the text input
		assert.Empty(t, m.maskModal.textInput.Value())

		typeMask(m, "*.txt")
		assert.Equal(t, "*.txt", m.maskModal.textInput.Value())

		TeaUpdate(m, enterMsg)
		assert.False(t, m.maskModal.open)
		assert.ElementsMatch(t, []string{file1, file2},
			m.getFocusedFilePanel().GetSelectedLocations())
	})

	t.Run("Select by mask switches browser mode to select mode", func(t *testing.T) {
		m := defaultTestModel(dir1)
		require.Equal(t, filepanel.BrowserMode, m.getFocusedFilePanel().PanelMode)

		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		require.True(t, m.maskModal.open)
		typeMask(m, "*.png")
		TeaUpdate(m, enterMsg)

		assert.Equal(t, filepanel.SelectMode, m.getFocusedFilePanel().PanelMode)
		assert.Equal(t, []string{file3}, m.getFocusedFilePanel().GetSelectedLocations())

		// The selection has to be visible, not just recorded. Checkboxes are
		// rendered only in select mode, a ticked one means the matched item is
		// drawn as selected
		view := m.viewContent()
		assert.Contains(t, view, common.CheckboxCheckedFocused, "matched item should render as selected")
		assert.Contains(t, view, common.CheckboxEmptyFocused, "unmatched items should render as unselected")
		// The panel footer reports the mode along with the count of selections
		assert.Contains(t, view, "Select"+icon.Space+"(1)")
	})

	t.Run("Unselect by mask", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ChangePanelMode[0]))
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.FilePanelSelectAllItem[0]))
		require.Equal(t, uint(3), m.getFocusedFilePanel().SelectedCount())

		TeaUpdate(m, utils.TeaRuneKeyMsg(unselectKey))
		require.True(t, m.maskModal.open)
		assert.False(t, m.maskModal.selecting)

		typeMask(m, "file*")
		TeaUpdate(m, enterMsg)
		assert.False(t, m.maskModal.open)
		assert.Equal(t, []string{file3}, m.getFocusedFilePanel().GetSelectedLocations())
	})

	t.Run("Unselect hotkey does nothing in browser mode", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(unselectKey))
		assert.False(t, m.maskModal.open)
	})

	t.Run("Cancelling the modal keeps the selection unchanged", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ChangePanelMode[0]))
		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		typeMask(m, "*.txt")

		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.CancelTyping[1]))
		assert.False(t, m.maskModal.open)
		assert.Equal(t, uint(0), m.getFocusedFilePanel().SelectedCount())
	})

	t.Run("Invalid mask keeps the modal open with an error", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ChangePanelMode[0]))
		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		typeMask(m, "[a-")
		TeaUpdate(m, enterMsg)

		assert.True(t, m.maskModal.open)
		assert.NotEmpty(t, m.maskModal.errorMsg)
		assert.Equal(t, uint(0), m.getFocusedFilePanel().SelectedCount())

		// Editing the mask clears the error
		TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyBackspace})
		assert.Empty(t, m.maskModal.errorMsg)

		// Fixing the mask applies it
		TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyBackspace})
		TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyBackspace})
		require.Empty(t, m.maskModal.textInput.Value())
		typeMask(m, "*.png")
		TeaUpdate(m, enterMsg)
		assert.False(t, m.maskModal.open)
		assert.Equal(t, []string{file3}, m.getFocusedFilePanel().GetSelectedLocations())
	})

	t.Run("Confirming an empty mask closes the modal", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ChangePanelMode[0]))
		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		require.True(t, m.maskModal.open)

		TeaUpdate(m, enterMsg)
		assert.False(t, m.maskModal.open)
		assert.Empty(t, m.maskModal.errorMsg)
		assert.Equal(t, uint(0), m.getFocusedFilePanel().SelectedCount())
	})

	t.Run("Mask matching nothing keeps the modal open", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ChangePanelMode[0]))
		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		typeMask(m, "*.go")
		TeaUpdate(m, enterMsg)

		assert.True(t, m.maskModal.open)
		assert.Equal(t, common.MaskNoMatchText, m.maskModal.errorMsg)
		assert.Equal(t, uint(0), m.getFocusedFilePanel().SelectedCount())
	})

	t.Run("Model renders with mask modal open", func(t *testing.T) {
		m := defaultTestModel(dir1)
		TeaUpdate(m, utils.TeaRuneKeyMsg(selectKey))
		require.True(t, m.maskModal.open)
		assert.Contains(t, m.viewContent(), strings.TrimSpace(common.MaskSelectTitle))
	})
}
