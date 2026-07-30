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
	slog.Debug("toggleContentPanelFocus called", "currentFocus", m.focusPanel)
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

	// For code files with embed_nvim=true, start embedded nvim
	if contentpanel.IsCodeFile(path, common.Config.ContentPanelCodeExtensions) && common.Config.ContentPanelEmbedNvim {
		err := m.fileModel.FilePreview.EnterNvimMode(path)
		if err != nil {
			slog.Error("Failed to start embedded nvim", "error", err)
		}
		return nil
	}

	// For code files without embed, launch nvim full-screen
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
	// In nvim mode, save via ":w" escape sequence
	if m.fileModel.FilePreview.IsNvimRunning() {
		m.fileModel.FilePreview.WriteNvimKey([]byte("\x1b:w\r"))
		return nil
	}
	err := m.fileModel.FilePreview.SaveEdit()
	if err != nil {
		slog.Error("Error saving file in content panel", "error", err)
	}
	return nil
}

func (m *model) contentPanelExitEdit() {
	// In nvim mode, quit via ":q" (esc to normal mode first)
	if m.fileModel.FilePreview.IsNvimRunning() {
		m.fileModel.FilePreview.StopNvim()
		return
	}
	m.fileModel.FilePreview.ExitEditMode()
}

// contentPanelHandleEditKey forwards a key message to the content panel's
// textarea when in edit mode, or to nvim PTY when embedded nvim is running.
func (m *model) contentPanelHandleEditKey(msg tea.KeyPressMsg) tea.Cmd {
	// If embedded nvim is running, send key to PTY
	if m.fileModel.FilePreview.IsNvimRunning() {
		m.fileModel.FilePreview.WriteNvimKey(nvimKeyBytes(msg))
		return nil
	}
	return m.fileModel.FilePreview.HandleEditKey(msg)
}

// nvimKeyBytes converts a Bubble Tea key message to bytes suitable for
// sending to nvim's PTY stdin.
func nvimKeyBytes(msg tea.KeyPressMsg) []byte {
	s := msg.String()
	switch {
	case s == "enter":
		return []byte("\r")
	case s == "esc":
		return []byte("\x1b")
	case s == "tab":
		return []byte("\t")
	case s == "backspace":
		return []byte("\x7f")
	case s == "ctrl+c":
		return []byte("\x03")
	case s == "ctrl+d":
		return []byte("\x04")
	case s == "ctrl+z":
		return []byte("\x1a")
	case s == "ctrl+u":
		return []byte("\x15")
	case s == "ctrl+w":
		return []byte("\x17")
	case len(s) == 1:
		return []byte(s)
	default:
		// For arrow keys and other special keys, send the escape sequence
		return nvimSpecialKey(s)
	}
}

func nvimSpecialKey(key string) []byte {
	switch key {
	case "up":
		return []byte("\x1b[A")
	case "down":
		return []byte("\x1b[B")
	case "right":
		return []byte("\x1b[C")
	case "left":
		return []byte("\x1b[D")
	case "home":
		return []byte("\x1b[H")
	case "end":
		return []byte("\x1b[F")
	case "pgup":
		return []byte("\x1b[5~")
	case "pgdown":
		return []byte("\x1b[6~")
	case "delete":
		return []byte("\x1b[3~")
	case "insert":
		return []byte("\x1b[2~")
	case "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12":
		// Not supporting function keys for now
		return nil
	default:
		return nil
	}
}
