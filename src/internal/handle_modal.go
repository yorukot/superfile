package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Cancel typing modal e.g. create file or directory
func (m *model) cancelTypingModal() {
	m.typingModal.textInput.Blur()
	m.typingModal.open = false
}
func (m *model) closeTypingModal(process processbar.Process,
	processBarModel *processbar.Model) {
	m.typingModal.errorMesssage = ""
	markProcessDone(process, processBarModel)
}

// Confirm to create file or directory
func (m *model) createItem(item string) error {
	if err := checkFileNameValidity(item); err != nil {
		m.typingModal.errorMesssage = err.Error()
		slog.Error("Errow while createItem during item creation", "error", err)
		return err
	}
	path := filepath.Join(m.typingModal.location, item)
	if !strings.HasSuffix(item, string(filepath.Separator)) {
		path, _ = renameIfDuplicate(path)
		if err := os.MkdirAll(filepath.Dir(path), utils.UserDirPerm); err != nil {
			slog.Error("Error while createItem during directory creation", "error", err)
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			slog.Error("Error while createItem during file creation", "error", err)
			return err
		}
		defer f.Close()
	} else {
		err := os.MkdirAll(path, utils.UserDirPerm)
		if err != nil {
			slog.Error("Error while createItem during directory creation", "error", err)
			return err
		}
	}
	return nil
}

func (m *model) getCreateCmd() tea.Cmd {
	if !m.typingModal.open {
		return nil
	}

	items := []string{m.typingModal.textInput.Value()}

	reqID := m.nextIoReqCnt()
	slog.Debug("Submitting create request", "id", reqID, "items cnt", len(items))
	m.cancelTypingModal()
	return func() tea.Msg {
		return m.createOperation(&m.processBarModel, items, reqID)
	}
}

func (m *model) createOperation(processBarModel *processbar.Model, items []string, reqID int) tea.Msg {
	if len(items) == 0 {
		return NewCreateOperationMsg(processbar.Cancelled, reqID)
	}
	p, err := processBarModel.SendAddProcessMsg(filepath.Base(items[0]), processbar.OpCreate, len(items), true)
	if err != nil {
		slog.Error("Cannot spawn a new process", "error", err)
		return NewCreateOperationMsg(processbar.Failed, reqID)
	}
	finalizer := func(state processbar.ProcessState, reqID int) tea.Msg {
		return NewCreateOperationMsg(state, reqID)
	}
	processor := makeCreateProcessor(m, p, processBarModel)
	msg := m.runFileProcessor(processor, finalizer, items, reqID)
	return msg
}

func makeCreateProcessor(model *model,
	process processbar.Process,
	processBarModel *processbar.Model) processbar.FileListProcessor {
	processorFunction := func(items []string) (processbar.Process, []string) {
		notProcessed := make([]string, 0)
		if len(items) == 0 {
			model.closeTypingModal(process, processBarModel)
			return process, notProcessed
		}

		for i, item := range items {
			err := model.createItem(item)
			if err != nil {
				process.State = processbar.Failed
				slog.Error("Error in create operation", "item", item, "error", err)
				process.ErrorMsg = formatFileError(item, err)
				notProcessed = items[i:]
				break
			}
			process.CurrentFile = filepath.Base(item)
			process.Done++
			processBarModel.TrySendingUpdateProcessMsg(process)
		}

		if process.State != processbar.Failed {
			process.State = processbar.Successful
			model.closeTypingModal(process, processBarModel)
		}
		return process, notProcessed
	}
	return processorFunction
}

// Cancel rename file or directory
func (m *model) cancelRename() {
	panel := m.getFocusedFilePanel()
	panel.Rename.Blur()
	panel.Renaming = false
	m.fileModel.Renaming = false
}

// Connfirm rename file or directory
func (m *model) confirmRename() {
	panel := m.getFocusedFilePanel()

	// Although we dont expect this to happen based on our current flow
	// Just adding it here to be safe
	if panel.Empty() {
		slog.Error("confirmRename called on empty panel")
		return
	}

	oldPath := panel.GetFocusedItem().Location
	newPath := filepath.Join(panel.Location, panel.Rename.Value())

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		slog.Error("Error while confirmRename during rename", "error", err)
		// Dont return. We have to also reset the panel and model information
	}
	m.fileModel.Renaming = false
	panel.Rename.Blur()
	panel.Renaming = false
}

func (m *model) confirmSortOptions() {
	panel := m.getFocusedFilePanel()
	panel.SortKind = m.sortModal.GetSelectedKind()
	m.sortModal.Close()
}

// Cancel search, this will clear all searchbar input
func (m *model) cancelSearch() {
	panel := m.getFocusedFilePanel()
	panel.SearchBar.Blur()
	panel.SearchBar.SetValue("")
}

// Confirm search. This will exit the search bar and filter the files
func (m *model) confirmSearch() {
	panel := m.getFocusedFilePanel()
	panel.SearchBar.Blur()
}

func (m *model) getFocusedFilePanel() *filepanel.Model {
	return m.fileModel.GetFocusedFilePanel()
}
