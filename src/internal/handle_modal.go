package internal

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

const remoteMutationTimeout = 30 * time.Second

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
	if !m.typingModal.open || m.typingModal.submitting {
		return nil
	}

	items := []string{m.typingModal.textInput.Value()}
	location := ensureLocationLabel(m.typingModal.paneLocation)
	if location.Path.String() == "" {
		location = filepanel.NewLocalLocation(m.typingModal.location)
	}

	reqID := m.nextIoReqCnt()
	slog.Debug("Submitting create request", "id", reqID, "items cnt", len(items))
	m.typingModal.submitting = true
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
		return NewProviderCreateOperationMsg(processbar.Cancelled, location, reqID)
	}
	p, err := processBarModel.SendAddProcessMsg(filepath.Base(items[0]), processbar.OpCreate, len(items), true)
	if err != nil {
		slog.Error("Cannot spawn a new process", "error", err)
		return NewProviderCreateOperationMsg(processbar.Failed, location, reqID)
	}
	finalizer := func(state processbar.ProcessState, reqID int) tea.Msg {
		return NewProviderCreateOperationMsg(state, location, reqID)
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

	ctx, cancel := mutationContext(location)
	defer cancel()
	session, err := m.ResolveFreshSession(ctx, location)
	if err != nil {
		return err
	}
	defer session.Close()

	path := pathJoinRaw(location.Path, name)
	target := locationWithPath(location, path)
	target, err = renameLocationIfDuplicate(ctx, session, target)
	if err != nil {
		return err
	}
	isDirectory := strings.HasSuffix(name, string(filepath.Separator)) || strings.HasSuffix(name, "/")
	if isDirectory {
		return session.Mkdir(ctx, target.Path, filesystem.MkdirOptions{Mode: utils.UserDirPerm, Parents: true})
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
	m.renameOperationPending = false
}

// Confirm rename file or directory.
func (m *model) confirmRename(overwrite bool) tea.Cmd {
	if m.renameOperationPending {
		return nil
	}
	panel := m.getFocusedFilePanel()

	// Although we dont expect this to happen based on our current flow
	// Just adding it here to be safe
	if panel.Empty() {
		slog.Error("confirmRename called on empty panel")
		return nil
	}
	if err := checkFileNameValidity(panel.Rename.Value()); err != nil {
		m.notifyModel = notify.New(true, "Rename failed", err.Error(), notify.NoAction)
		return nil
	}

	panelIndex := m.fileModel.FocusedPanelIndex
	location := panel.CurrentLocation()
	itemLocation := elementLocation(location, panel.GetFocusedItem())
	newPath := pathJoinRaw(location.Path, panel.Rename.Value())
	if itemLocation.Path.String() == newPath.String() {
		m.resetRenameState(panelIndex)
		return nil
	}
	m.renameOperationPending = true
	panel.Rename.Blur()
	reqID := m.nextIoReqCnt()

	return func() tea.Msg {
		ctx, cancel := mutationContext(location)
		defer cancel()
		session, err := m.ResolveFreshSession(ctx, location)
		conflict := false
		if err == nil {
			defer session.Close()
			if !overwrite {
				conflict, err = sessionPathExists(ctx, session, newPath)
			}
			if err == nil && !conflict {
				err = session.Rename(
					ctx,
					itemLocation.Path,
					newPath,
					filesystem.RenameOptions{Overwrite: overwrite},
				)
			}
		}
		return NewFileMutationMsg(FileMutationRename, location, panelIndex, conflict, err, reqID)
	}
}

func mutationContext(location filesystem.Location) (context.Context, context.CancelFunc) {
	if location.Provider != filesystem.ProviderLocal {
		return context.WithTimeout(context.Background(), remoteMutationTimeout)
	}
	return context.WithCancel(context.Background())
}

func (m *model) resetRenameState(panelIndex int) {
	m.renameOperationPending = false
	m.fileModel.Renaming = false
	if panelIndex < 0 || panelIndex >= len(m.fileModel.FilePanels) {
		return
	}
	panel := &m.fileModel.FilePanels[panelIndex]
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
