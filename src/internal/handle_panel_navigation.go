package internal

import (
	"log/slog"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/contentpanel"
)

// Pinned directory
func (m *model) pinnedDirectory() {
	panel := m.getFocusedFilePanel()
	err := m.sidebarModel.TogglePinnedDirectory(panel.Location)
	if err != nil {
		slog.Error("Error while toggling pinned directory", "error", err)
	}
}

// Focus on sidebar
func (m *model) focusOnSideBar() {
	if common.Config.SidebarWidth == 0 {
		return
	}
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.getFocusedFilePanel().IsFocused = true
	} else {
		m.focusPanel = sidebarFocus
		m.getFocusedFilePanel().IsFocused = false
	}
}

// Focus on processbar
func (m *model) focusOnProcessBar() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.getFocusedFilePanel().IsFocused = true
	} else {
		m.focusPanel = processBarFocus
		m.getFocusedFilePanel().IsFocused = false
	}
}

// focus on metadata
func (m *model) focusOnMetadata() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.getFocusedFilePanel().IsFocused = true
	} else {
		m.focusPanel = metadataFocus
		m.getFocusedFilePanel().IsFocused = false
	}
}

// toggleContentPanelFocus cycles focus into/out of the content (preview) panel.
func (m *model) toggleContentPanelFocus() {
	if m.focusPanel == contentPanelFocus {
		m.focusPanel = nonePanelFocus
		m.getFocusedFilePanel().IsFocused = true
		m.fileModel.FilePreview.SetFocused(false)
	} else {
		m.focusPanel = contentPanelFocus
		m.getFocusedFilePanel().IsFocused = false
		m.fileModel.FilePreview.SetFocused(true)
	}
}

func (m *model) contentPanelEnterEdit() tea.Cmd {
	panel := m.getFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		return nil
	}
	path := panel.GetFocusedItem().Location

	// For code files, use nvim full-screen (if not embed mode)
	if contentpanel.IsCodeFile(path, common.Config.ContentPanelCodeExtensions) && !common.Config.ContentPanelEmbedNvim {
		return m.launchNvimForContentPanel(path)
	}

	ok, reason := m.fileModel.FilePreview.EnterEditMode(path)
	if !ok {
		slog.Debug("Cannot enter edit mode", "path", path, "reason", reason)
		return nil
	}
	return nil
}

// launchNvimForContentPanel opens nvim full-screen via tea.ExecProcess, then
// refreshes the content panel with the saved file.
func (m *model) launchNvimForContentPanel(path string) tea.Cmd {
	nvimCmd := common.Config.ContentPanelNvimCmd
	if nvimCmd == "" {
		nvimCmd = "nvim"
	}
	parts := strings.Fields(nvimCmd)
	args := append(parts[1:], path)
	c := exec.Command(parts[0], args...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return contentPanelNvimFinishedMsg{path: path, err: err}
	})
}

func (m *model) contentPanelSaveEdit() tea.Cmd {
	err := m.fileModel.FilePreview.SaveEdit()
	if err != nil {
		slog.Error("Error saving file in content panel", "error", err)
	}
	return nil
}

func (m *model) contentPanelExitEdit() {
	m.fileModel.FilePreview.ExitEditMode()
}

// contentPanelHandleEditKey forwards a key message to the content panel's
// textarea when in edit mode.
func (m *model) contentPanelHandleEditKey(msg tea.KeyPressMsg) tea.Cmd {
	return m.fileModel.FilePreview.HandleEditKey(msg)
}
