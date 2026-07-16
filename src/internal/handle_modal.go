package internal

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Cancel typing modal e.g. create file or directory
func (m *model) cancelTypingModal() {
	m.typingModal.textInput.Blur()
	m.typingModal.open = false
	m.typingModal.paneLocation = filesystem.Location{}
}

// createItem keeps the provider-aware create operation available to direct callers.
func (m *model) createItem() tea.Cmd {
	return m.getCreateCmd()
}

func (m *model) getCreateCmd() tea.Cmd {
	if !m.typingModal.open {
		return nil
	}

	items := []string{m.typingModal.textInput.Value()}
	location := ensureLocationLabel(m.typingModal.paneLocation)
	if location.Path.String() == "" {
		location = filepanel.NewLocalLocation(m.typingModal.location)
	}

	reqID := m.nextIoReqCnt()
	slog.Debug("Submitting create request", "id", reqID, "items cnt", len(items))
	m.cancelTypingModal()
	return func() tea.Msg {
		return m.createOperation(&m.processBarModel, location, items, reqID)
	}
}

func (m *model) createOperation(
	processBarModel *processbar.Model,
	location filesystem.Location,
	items []string,
	reqID int,
) tea.Msg {
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
	processor := m.makeCreateProcessor(location, p, processBarModel)
	return m.runFileProcessor(processor, finalizer, items, reqID)
}

func (m *model) makeCreateProcessor(
	location filesystem.Location,
	process processbar.Process,
	processBarModel *processbar.Model,
) processbar.FileListProcessor {
	return func(items []string) (processbar.Process, []string) {
		notProcessed := make([]string, 0)
		if len(items) == 0 {
			markProcessDone(process, processBarModel)
			return process, notProcessed
		}

		for i, item := range items {
			err := m.createItemAt(location, item)
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
			markProcessDone(process, processBarModel)
		}
		return process, notProcessed
	}
}

func (m *model) createItemAt(location filesystem.Location, name string) error {
	if err := checkFileNameValidity(name); err != nil {
		return err
	}

	ctx := context.Background()
	session, err := m.ResolveFreshSession(ctx, location)
	if err != nil {
		return err
	}
	defer session.Close()

	path := pathJoinRaw(location.Path, name)
	isDirectory := strings.HasSuffix(name, string(filepath.Separator)) || strings.HasSuffix(name, "/")
	if isDirectory {
		return session.Mkdir(ctx, path, filesystem.MkdirOptions{Mode: utils.UserDirPerm, Parents: true})
	}

	target := locationWithPath(location, path)
	target, err = renameLocationIfDuplicate(ctx, session, target)
	if err != nil {
		return err
	}
	if err := session.Mkdir(
		ctx,
		pathDir(target.Path),
		filesystem.MkdirOptions{Mode: utils.UserDirPerm, Parents: true},
	); err != nil {
		return err
	}
	return session.Create(ctx, target.Path, nil, filesystem.CreateOptions{Mode: utils.UserFilePerm})
}

// Cancel rename file or directory
func (m *model) cancelRename() {
	panel := m.getFocusedFilePanel()
	panel.Rename.Blur()
	panel.Renaming = false
	m.fileModel.Renaming = false
}

// Confirm rename file or directory.
func (m *model) confirmRename() {
	panel := m.getFocusedFilePanel()

	// Although we dont expect this to happen based on our current flow
	// Just adding it here to be safe
	if panel.Empty() {
		slog.Error("confirmRename called on empty panel")
		return
	}

	location := panel.CurrentLocation()
	itemLocation := elementLocation(location, panel.GetFocusedItem())
	newPath := pathJoinRaw(location.Path, panel.Rename.Value())
	overwrite := m.IsRenamingConflicting()
	session, err := m.ResolveSession(context.Background(), location)
	if err != nil {
		slog.Error("Error while confirmRename during session resolution", "error", err)
	} else {
		defer session.Close()
		if err = session.Rename(
			context.Background(),
			itemLocation.Path,
			newPath,
			filesystem.RenameOptions{Overwrite: overwrite},
		); err != nil {
			slog.Error("Error while confirmRename during rename", "error", err)
			// Dont return. We have to also reset the panel and model information
		}
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
