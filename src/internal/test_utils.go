package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

var SampleDataBytes = []byte("This is sample") //nolint: gochecknoglobals // Effectively const

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

// -------------------- Model setup utils

func defaultTestModel(dirs ...string) model {
	m := defaultModelConfig(false, false, false, dirs)
	cmd := m.Init()
	msg := cmd
	slog.Debug("[Test] defaultTestModel()", "msg", msg)
	_, _ = TeaUpdate(&m, tea.WindowSizeMsg{Width: 2 * common.MinimumWidth, Height: 2 * common.MinimumHeight})
	return m
}

// Helper function to setup panel mode and selection
func setupPanelModeAndSelection(t *testing.T, m *model, useSelectMode bool, itemName string, selectedItems []string) {
	t.Helper()
	panel := m.getFocusedFilePanel()

	if useSelectMode {
		// Switch to select mode and set selected items
		m.changeFilePanelMode()
		require.Equal(t, selectMode, panel.panelMode)
		panel.selected = selectedItems
	} else {
		// Find the item in browser mode
		itemIndex := findItemIndexInPanel(panel, itemName)
		require.NotEqual(t, -1, itemIndex, "%s should be found in panel", itemName)
		panel.cursor = itemIndex
	}
}

// --------------------  Bubletea utilities

type TeaProgram struct {
	m          model
	msgs       chan tea.Msg
	sentToChan  int
	sentMsgCnt int
	mutex      sync.Mutex
	mutex2      sync.Mutex
}

func DefaultTeaProgram(startdirs ...string) *TeaProgram {
	return &TeaProgram{
		defaultModelConfig(false, false, false, startdirs),
		// TODO: This might be a hacky way. Figure out if we can
		// get it working with single buffer channel
		make(chan tea.Msg, 100),
		0,
		0,
		sync.Mutex{},
		sync.Mutex{},
	}
}

func (p *TeaProgram) SendMessage(msg tea.Msg) {
	p.mutex2.Lock()
	defer p.mutex2.Unlock()
	p.sentToChan++
	slog.Debug("[Test] Writing msg to channel", "type", reflect.TypeOf(msg), "cnt", p.sentToChan)
	p.msgs <- msg
	slog.Debug("[Test] Done Writing msg to channel", "type", reflect.TypeOf(msg), "cnt", p.sentToChan)

}

// Only one instance should be running at a time
func (p *TeaProgram) SendMessageBlocking(msg tea.Msg) {
	p.SendMessage(msg)
	p.SendAllMsg()
}

// Only one instance should be running at a time
func (p *TeaProgram) sendMessageToModel(msg tea.Msg) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Todo : Add a way to fail tests on err
	curId := p.sentMsgCnt
	p.sentMsgCnt++
	slog.Debug("[Test] Sending message", "id", curId, "type", reflect.TypeOf(msg))
	slog.Debug("model clipboard info before send", "id", curId, "items", p.m.copyItems.items, "cut", p.m.copyItems.cut)

	cmd, err := TeaUpdate(&p.m, msg)
	slog.Debug("model clipboard info after send", "id", curId, "items", p.m.copyItems.items, "cut", p.m.copyItems.cut)

	if err == nil {
		p.RunCmd(cmd)
	}
	slog.Debug("model clipboard info after run command", "id", curId, "items", p.m.copyItems.items, "cut", p.m.copyItems.cut)

}

// Block till all current messages are dealt with
func (p *TeaProgram) SendAllMsg() {
	for {
		select {
		case msg := <-p.msgs:
			p.sendMessageToModel(msg)
		default:
			slog.Debug("[Test] SendAllMsg() : No messages in channel, exiting")
			return
		}
	}
}

func (p *TeaProgram) Run() {
	slog.Debug("[Test] TeaProgram : Started")

	p.RunCmd(p.m.Init())
	p.SendAllMsg()
}

func (p *TeaProgram) RunCmdBlocking(cmd tea.Cmd) {
	msg := cmd()
	switch msg := msg.(type) {
	case tea.BatchMsg:
		for _, cur_cmd := range msg {
			go p.RunCmdBlocking(cur_cmd)
		}
	default:
		p.SendMessage(msg)
	}
}

func (p *TeaProgram) RunCmd(cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	go p.RunCmdBlocking(cmd)
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

// Helper function to perform copy or cut operation
func performCopyOrCutOperation(t *testing.T, p *TeaProgram, isCut bool) {
	t.Helper()
	if isCut {
		p.SendMessageBlocking(utils.TeaRuneKeyMsg(common.Hotkeys.CutItems[0]))
	} else {
		p.SendMessageBlocking(utils.TeaRuneKeyMsg(common.Hotkeys.CopyItems[0]))
	}
}

// -------------- Validation Utilities

// Helper function to verify clipboard state after copy/cut
func verifyClipboardState(t *testing.T, m *model, isCut bool, useSelectMode bool, selectedItemsCount int) {
	t.Helper()
	assert.Equal(t, isCut, m.copyItems.cut, "Clipboard cut state should match operation")
	if useSelectMode {
		assert.Len(t, m.copyItems.items, selectedItemsCount, "Clipboard should contain all selected items")
	} else {
		assert.Len(t, m.copyItems.items, 1, "Clipboard should contain one item")
	}
}

// Helper function to verify file or directory exists
func verifyPathExists(t *testing.T, path, message string) {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err, message)
	if info.IsDir() {
		assert.DirExists(t, path, message)
	} else {
		assert.FileExists(t, path, message)
	}
}

// Helper function to verify file or directory doesn't exist after cut
func verifyPathNotExistsEventually(t *testing.T, path, message string) {
	t.Helper()
	assert.Eventually(t, func() bool {
		_, err := os.Stat(path)
		return os.IsNotExist(err)
	}, time.Second, 10*time.Millisecond, message)
}

// Helper function to verify expected destination files exist
func verifyDestinationFiles(t *testing.T, targetDir string, expectedDestFiles []string) {
	t.Helper()
	for _, expectedFile := range expectedDestFiles {
		destPath := filepath.Join(targetDir, expectedFile)
		assert.Eventually(t, func() bool {
			_, err := os.Stat(destPath)
			return err == nil
		}, time.Second, 10*time.Millisecond, "%s should exist in destination", expectedFile)
	}
}

// Helper function to verify prevented paste results
func verifyPreventedPasteResults(t *testing.T, m *model, originalPath string) {
	t.Helper()
	if originalPath != "" {
		verifyPathExists(t, originalPath, "Original file should still exist when paste is prevented")
	}
	// Clipboard should not be cleared when paste is prevented
	assert.NotEmpty(t, m.copyItems.items, "Clipboard should not be cleared when paste is prevented")
}

// Helper function to verify successful paste results
func verifySuccessfulPasteResults(t *testing.T, targetDir string, expectedDestFiles []string, originalPath string, shouldOriginalExist bool) {
	t.Helper()
	// Verify expected files were created in destination
	verifyDestinationFiles(t, targetDir, expectedDestFiles)

	// Verify original file existence based on operation type
	if originalPath != "" {
		if shouldOriginalExist {
			verifyPathExists(t, originalPath, "Original file should exist after copy operation")
		} else {
			verifyPathNotExistsEventually(t, originalPath, "Original file should not exist after cut operation")
		}
	}
}

// -------------- Other utilities

// Helper function to find item index in panel by name
func findItemIndexInPanel(panel *filePanel, itemName string) int {
	for i, elem := range panel.element {
		if elem.name == itemName {
			return i
		}
	}
	return -1
}

// Helper function to navigate to target directory if different from start
func navigateToTargetDir(t *testing.T, p *TeaProgram, startDir, targetDir string) {
	t.Helper()
	if targetDir != startDir {
		slog.Debug("model clipboard info 1a", "items", p.m.copyItems.items, "cut", p.m.copyItems.cut)
		err := p.m.updateCurrentFilePanelDir(targetDir)
		require.NoError(t, err)
		slog.Debug("model clipboard info 1b", "items", p.m.copyItems.items, "cut", p.m.copyItems.cut)

		p.SendMessageBlocking(nil)
		slog.Debug("model clipboard info 1c", "items", p.m.copyItems.items, "cut", p.m.copyItems.cut)
	}
}

// Helper function to get original path for existence check
func getOriginalPath(useSelectMode bool, itemName, startDir string) string {
	if !useSelectMode && itemName != "" {
		return filepath.Join(startDir, itemName)
	}
	return ""
}
