package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

var errInvalidOpenWithChoice = errors.New("invalid Open With application choice")

func (m *model) openWithMenu() {
	panel := m.getFocusedFilePanel()
	if panel.Empty() {
		return
	}
	selectedItem := panel.GetFocusedItem()
	if selectedItem.Directory {
		return
	}

	applications, err := common.ResolveOpenWith(common.Config.OpenWith, selectedItem.Location)
	if err != nil {
		slog.Error("Could not resolve Open With applications", "error", err, "path", selectedItem.Location)
		return
	}
	m.openWithModal.Open(selectedItem.Location, applications)
}

func (m *model) confirmOpenWith() tea.Cmd {
	application, ok := m.openWithModal.Selected()
	if !ok {
		return nil
	}
	filePath := m.openWithModal.Path()
	m.openWithModal.Close()

	cmd, err := newOpenWithCommand(application, filePath)
	if err != nil {
		slog.Error("Could not create Open With command", "error", err, "path", filePath)
		return nil
	}
	if application.Terminal {
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			return editorFinishedMsg{err}
		})
	}

	utils.DetachFromTerminal(cmd)
	if err := cmd.Start(); err != nil {
		slog.Error("Could not start Open With application", "error", err, "path", filePath,
			"application", application.Name)
	}
	return nil
}

func newOpenWithCommand(application common.OpenWithApplication, filePath string) (*exec.Cmd, error) {
	command := strings.TrimSpace(application.Command)
	if command == "" {
		return nil, fmt.Errorf("%w: command cannot be empty", errInvalidOpenWithChoice)
	}
	if filePath == "" {
		return nil, fmt.Errorf("%w: file path cannot be empty", errInvalidOpenWithChoice)
	}

	args := append([]string(nil), application.Args...)
	args = append(args, filePath)
	return exec.Command(command, args...), nil //nolint:gosec // Command is explicitly user-configured.
}
