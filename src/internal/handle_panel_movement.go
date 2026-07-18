package internal

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/pkg/utils"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

const remoteNavigationTimeout = 10 * time.Second

// Back to parent directory
func (m *model) parentDirectory() tea.Cmd {
	panel := m.getFocusedFilePanel()
	if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
		return m.remoteNavigationCmd(panel.CurrentLocation().Path.Dir())
	}
	err := panel.ParentDirectory()
	if err != nil {
		slog.Error("Error while changing to parent directory", "error", err)
	}
	return nil
}

// Enter directory or open file with default application
// TODO: Unit test this
func (m *model) enterPanel() tea.Cmd {
	panel := m.getFocusedFilePanel()

	if panel.Empty() {
		return nil
	}
	selectedItem := panel.GetFocusedItem()
	if selectedItem.Directory {
		targetPath := selectedItem.Location
		if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
			return m.remoteNavigationCmd(selectedItem.Path)
		}

		if selectedItem.Info.Mode()&os.ModeSymlink != 0 {
			var symlinkErr error
			targetPath, symlinkErr = filepath.EvalSymlinks(targetPath)
			if symlinkErr != nil {
				return nil
			}

			// targetPath shouldn't be a link now, so Stat and Lstat should be same
			if targetInfo, lstatErr := os.Lstat(targetPath); lstatErr != nil || !targetInfo.IsDir() {
				return nil
			}
		}
		// TODO : Propagate error out from this this function. Return here, instead of logging
		err := m.updateCurrentFilePanelDir(targetPath)
		if err != nil {
			slog.Error("Error while changing to directory", "error", err, "target", targetPath)
		}
		return nil
	}
	if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
		return m.unsupportedRemoteOperationCmd(panel.CurrentLocation(), filesystem.OperationOpenWith)
	}

	if variable.ChooserFile != "" {
		chooserErr := m.chooserFileWriteAndQuit(panel.GetFocusedItem().Location)
		if chooserErr == nil {
			return nil
		}
		// Continue with preview if file is not writable
		slog.Error("Error while writing to chooser file, continuing with file open", "error", chooserErr)
	}
	m.executeOpenCommand()
	return nil
}

func (m *model) remoteNavigationCmd(target filesystem.Path) tea.Cmd {
	panelIndex := m.fileModel.FocusedPanelIndex
	source := m.getFocusedFilePanel().CurrentLocation()
	sessionState, err := m.fileModel.PaneSession(panelIndex)
	reqID := m.nextIoReqCnt()
	if err != nil || sessionState.Browser == nil {
		if err == nil {
			err = filesystem.NewDisconnectedError(
				source.Provider,
				filesystem.OperationNavigate,
				target,
				"session is unavailable",
			)
		}
		return func() tea.Msg {
			return NewRemoteNavigationMsg(panelIndex, source, target, 0, err, reqID)
		}
	}
	browser := sessionState.Browser
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), remoteNavigationTimeout)
		defer cancel()
		_, listErr := browser.List(ctx, target)
		return NewRemoteNavigationMsg(panelIndex, source, target, sessionState.Generation, listErr, reqID)
	}
}

func (m *model) executeOpenCommand() {
	panel := m.getFocusedFilePanel()

	filePath := panel.GetFocusedItem().Location

	openCommand := "xdg-open"
	switch runtime.GOOS {
	case utils.OsDarwin:
		openCommand = "open"
	case utils.OsWindows:
		dllpath := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
		dllfile := "url.dll,FileProtocolHandler"

		//nolint:gosec // Uses Windows system handler to open the selected file.
		cmd := exec.Command(dllpath, dllfile, filePath)
		err := cmd.Start()
		if err != nil {
			slog.Error("Error while open file with", "error", err)
		}

		return
	}

	// For now open_with works only for mac and linux
	// TODO: Make it in parity with windows.
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	if extEditor, ok := common.Config.OpenWith[ext]; ok {
		openCommand = extEditor
	}

	cmd := exec.Command(openCommand, filePath)
	utils.DetachFromTerminal(cmd)
	err := cmd.Start()
	if err != nil {
		// TODO: This kind of errors should go to user facing pop ups
		slog.Error("Error while open file with", "error", err)
	}
}

// Switch to the directory where the sidebar cursor is located
func (m *model) sidebarSelectDirectory() {
	// We can't do this when we have only divider directories
	// m.sidebarModel.directories[m.sidebarModel.cursor].location would point to a divider dir.
	if m.sidebarModel.NoActualDir() {
		return
	}
	// TODO(Refactor): Move this to a function m.ResetFocus()
	m.focusPanel = nonePanelFocus
	panel := m.getFocusedFilePanel()

	location := filepanel.NewLocalLocation(m.sidebarModel.GetCurrentDirectoryLocation())
	err := m.fileModel.SetPaneLocation(m.fileModel.FocusedPanelIndex, location)
	if err == nil {
		panel.UpdateElementsIfNeeded(true, m.fileModel.DisplayDotFiles)
	}
	if err != nil {
		slog.Error("Error switching to sidebar directory", "error", err)
	}
	panel.IsFocused = true
}

// Toggle dotfile display or not
func (m *model) toggleDotFileController() tea.Cmd {
	cmd := m.fileModel.ToggleDotFile()
	err := utils.WriteBoolFile(variable.ToggleDotFile, m.fileModel.DisplayDotFiles)
	if err != nil {
		slog.Error("Error while updating toggleDotFile data", "error", err)
	}
	return cmd
}

// Toggle dotfile display or not
func (m *model) toggleFooterController() tea.Cmd {
	m.toggleFooter = !m.toggleFooter
	err := utils.WriteBoolFile(variable.ToggleFooter, m.toggleFooter)
	if err != nil {
		slog.Error("Error while updating toggleFooter data", "error", err)
	}
	m.setHeightValues()
	return m.updateComponentDimensions()
}

// Focus on search bar
func (m *model) searchBarFocus() {
	panel := m.getFocusedFilePanel()
	if panel.SearchBar.Focused() {
		panel.SearchBar.Blur()
	} else {
		panel.SearchBar.Focus()
		m.firstTextInput = true
	}

	// config search bar width
	panel.SearchBar.SetWidth(m.fileModel.SinglePanelWidth - common.InnerPadding)
}

func (m *model) sidebarSearchBarFocus() {
	if m.sidebarModel.SearchBarFocused() {
		// Ideally Code should never reach here. Once sidebar is focused, we should
		// not cause sidebarSearchBarFocus() event by pressing search key
		slog.Error("sidebarSearchBarFocus() called on Focused sidebar")
		m.sidebarModel.SearchBarBlur()
		return
	}
	m.sidebarModel.SearchBarFocus()
	m.firstTextInput = true
}
