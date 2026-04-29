package internal

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/pkg/utils"
)

const termfilechooserInstructionsPrefix = "* xdg-desktop-portal-termfilechooser instructions *"

func (m *model) initializeChooserState() {
	if !m.chooser.request.Active() || !m.isSaveChooserMode() {
		return
	}

	panel := m.getFocusedFilePanel()
	fallbackDir := panel.Location
	startDir, saveName := resolveSaveChooserStartPath(m.chooser.request.SuggestedPath, fallbackDir)
	if err := panel.UpdateCurrentFilePanelDir(startDir); err != nil {
		notifyErr := err
		slog.Error("Failed to initialize save chooser directory", "path", startDir, "error", err)
		if fallbackDir != startDir {
			if fallbackErr := panel.UpdateCurrentFilePanelDir(fallbackDir); fallbackErr == nil {
				panel.EnableSaveMode(saveName)
				return
			} else {
				notifyErr = fallbackErr
				slog.Error(
					"Failed to initialize fallback save chooser directory",
					"path",
					fallbackDir,
					"error",
					fallbackErr,
				)
			}
		}
		m.notifyModel = notify.New(true, common.SaveInitErrorTitle, notifyErr.Error(), notify.NoAction)
		return
	}
	panel.EnableSaveMode(saveName)
}

func resolveSaveChooserStartPath(suggestedPath string, fallbackDir string) (string, string) {
	if suggestedPath == "" {
		return fallbackDir, ""
	}

	info, err := os.Stat(suggestedPath)
	if err == nil {
		if info.IsDir() {
			return suggestedPath, ""
		}
		return filepath.Dir(suggestedPath), filepath.Base(suggestedPath)
	}

	parentDir := filepath.Dir(suggestedPath)
	parentInfo, parentErr := os.Stat(parentDir)
	if parentErr == nil && parentInfo.IsDir() {
		return parentDir, filepath.Base(suggestedPath)
	}

	return fallbackDir, filepath.Base(suggestedPath)
}

func (m *model) hasChooserRequest() bool {
	return m.chooser.request.Active()
}

func (m *model) isOpenChooserMode() bool {
	return m.chooser.request.Mode == variable.ChooserModeOpen
}

func (m *model) isSaveChooserMode() bool {
	return m.chooser.request.Mode == variable.ChooserModeSave
}

func (m *model) chooserWriteAndQuit(paths []string) error {
	if !m.hasChooserRequest() {
		return nil
	}

	err := os.WriteFile(m.chooser.request.OutputFile, []byte(strings.Join(paths, "\n")), utils.UserFilePerm)
	if err != nil {
		return err
	}
	m.modelQuitState = quitInitiated
	return nil
}

func (m *model) openChooserWriteSelectionAndQuit() error {
	panel := m.getFocusedFilePanel()
	if panel.Empty() {
		return nil
	}

	paths := []string{panel.GetFocusedItem().Location}
	if panel.PanelMode == filepanel.SelectMode && panel.SelectedCount() > 0 {
		paths = panel.GetOrderedSelectedLocations()
	}
	return m.chooserWriteAndQuit(paths)
}

func (m *model) saveChooserConfirmFocusedItem() {
	panel := m.getFocusedFilePanel()
	if panel.Empty() {
		return
	}

	focused := panel.GetFocusedItem()
	switch {
	case focused.SaveTarget:
		m.confirmSaveChooserPath(panel.GetSaveEntryLocation())
	case focused.Directory:
		return
	default:
		m.confirmSaveChooserPath(focused.Location)
	}
}

func (m *model) saveChooserConfirmCurrentDirectory() {
	m.confirmSaveChooserPath(m.getFocusedFilePanel().GetSaveEntryLocation())
}

func (m *model) confirmSaveChooserPath(path string) {
	if path == "" {
		return
	}

	info, err := os.Stat(path)
	switch {
	case err == nil:
		m.confirmExistingSaveChooserPath(path, info)
	case os.IsNotExist(err):
		m.confirmNewSaveChooserPath(path)
	default:
		m.notifyModel = notify.New(true, common.SaveAccessErrorTitle, err.Error(), notify.NoAction)
		slog.Error("Save chooser target stat failed", "path", path, "error", err)
	}
}

func (m *model) confirmExistingSaveChooserPath(path string, info os.FileInfo) {
	if info.IsDir() {
		m.notifyModel = notify.New(true, common.SaveDirErrorTitle, common.SaveDirErrorContent, notify.NoAction)
		return
	}

	if m.isPortalSavePlaceholder(path) {
		if writeErr := m.chooserWriteAndQuit([]string{path}); writeErr != nil {
			slog.Error("Error while writing save chooser result", "error", writeErr)
		}
		return
	}

	m.warnModalForSaveOverwrite(path)
}

func (m *model) confirmNewSaveChooserPath(path string) {
	parentDir := filepath.Dir(path)
	parentInfo, statErr := os.Stat(parentDir)
	if statErr != nil || !parentInfo.IsDir() {
		slog.Error("Save chooser target parent is invalid", "path", path, "error", statErr)
		return
	}

	if createErr := createSaveChooserPlaceholder(path); createErr != nil {
		if errors.Is(createErr, os.ErrExist) {
			m.warnModalForSaveOverwrite(path)
			return
		}
		slog.Error("Error while creating save chooser placeholder", "path", path, "error", createErr)
		return
	}

	if writeErr := m.chooserWriteAndQuit([]string{path}); writeErr != nil {
		removeErr := os.Remove(path)
		if removeErr != nil {
			slog.Error(
				"Error while writing save chooser result and removing placeholder",
				"writeError",
				writeErr,
				"removeError",
				removeErr,
				"path",
				path,
			)
			return
		}
		slog.Error(
			"Error while writing save chooser result; removed placeholder",
			"writeError",
			writeErr,
			"path",
			path,
		)
	}
}

func (m *model) warnModalForSaveOverwrite(path string) {
	m.chooser.overwritePath = path
	m.notifyModel = notify.New(
		true,
		common.SaveOverwriteWarnTitle,
		common.SaveOverwriteWarnContent,
		notify.SaveOverwriteAction,
	)
}

func (m *model) confirmSaveOverwrite() {
	if m.chooser.overwritePath == "" {
		return
	}

	path := m.chooser.overwritePath
	m.chooser.overwritePath = ""
	if err := m.chooserWriteAndQuit([]string{path}); err != nil {
		slog.Error("Error while confirming save overwrite", "error", err)
	}
}

func (m *model) cancelSaveOverwrite() {
	m.chooser.overwritePath = ""
}

func createSaveChooserPlaceholder(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, utils.UserFilePerm)
	if err != nil {
		return err
	}
	return file.Close()
}

func (m *model) isPortalSavePlaceholder(path string) bool {
	if path == "" || path != m.chooser.request.SuggestedPath {
		return false
	}

	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, len(termfilechooserInstructionsPrefix))
	n, readErr := io.ReadFull(file, buf)
	if readErr != nil && readErr != io.ErrUnexpectedEOF {
		return false
	}

	return string(buf[:n]) == termfilechooserInstructionsPrefix
}
