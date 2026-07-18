package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/trash"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/spferror"
	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
)

// Processes any standard (f.e. deletion) operation with a list of files
func (m *model) runFileProcessor(processor processbar.FileListProcessor,
	finalizer processbar.ProcessFinalizer,
	items []string,
	reqID int,
) tea.Msg {
	slog.Debug("Lock mutex for modal error window")
	m.mutexErrorModal.Lock()
	result, toDo := processor(items)
	if result.State == processbar.Failed && len(toDo) > 0 {
		// we unlock mutexErrorModal on dispatch SpfErrorModalUpdateMsg
		errorModel := spferror.New(true,
			"Error",
			result.ErrorMsg, spferror.NewFileListError(toDo, processor, finalizer))
		return NewSpfErrorModalMsg(errorModel, reqID)
	}
	slog.Debug("Unlock mutex for modal error window")
	m.mutexErrorModal.Unlock()
	return finalizer(result.State, reqID)
}

func markProcessDone(process processbar.Process, processBarModel *processbar.Model) {
	process.DoneTime = time.Now()
	err := processBarModel.SendUpdateProcessMsg(process, true)
	if err != nil {
		slog.Error("Failed to send final delete operation update", "error", err)
	}
}

func formatFileError(filePath string, err error) string {
	var e *os.LinkError
	if errors.As(err, &e) {
		return fmt.Sprintf("Deleting %s: \n%s", filePath, e.Err.Error())
	}
	return err.Error()
}

// Create a file in the currently focus file panel
// TODO: Fix it. It doesn't creates a new file. It just opens a file model,
// that allows you to create a file. Actual creation happens here - createItem() in handle_modal.go
func (m *model) panelCreateNewFile() {
	panel := m.getFocusedFilePanel()

	m.typingModal.location = panel.Location
	m.typingModal.paneLocation = panel.CurrentLocation()
	m.typingModal.open = true
	m.typingModal.submitting = false
	m.typingModal.textInput = common.GenerateNewFileTextInput()
	m.firstTextInput = true
}

// Rename file where the cusror is located
// TODO: Fix this. It doesn't do any rename, just opens the rename text input
// Actual rename happens at confirmRename() in handle_modal.go
func (m *model) panelItemRename() {
	panel := m.getFocusedFilePanel()
	if panel.Empty() {
		return
	}

	cursorPos := -1
	nameRunes := []rune(panel.GetFocusedItem().Name)
	nameLen := len(nameRunes)
	for i := nameLen - 1; i >= 0; i-- {
		if nameRunes[i] == '.' {
			cursorPos = i
			break
		}
	}
	if cursorPos == -1 || cursorPos == 0 && nameLen > 0 || panel.GetFocusedItem().Directory {
		cursorPos = nameLen
	}

	m.fileModel.Renaming = true
	panel.Renaming = true
	m.firstTextInput = true
	// TODO: Don't re-create a new model on each rename. Don't create
	// unnecessary gargage for collection. Reuse the existing model.
	// Maintain its state, dimensions. Update its cursor and text when needed
	panel.Rename = common.GenerateRenameTextInput(
		m.fileModel.SinglePanelWidth-common.InnerPadding,
		cursorPos,
		panel.GetFocusedItem().Name)
}

func (m *model) getDeleteCmd(permDelete bool) tea.Cmd {
	panel := m.getFocusedFilePanel()
	if panel.Empty() {
		return nil
	}
	items := focusedItemLocations(panel)
	useTrash := panel.CurrentLocation().Provider == filesystem.ProviderLocal &&
		m.hasTrash && trash.Available(panel.Location) &&
		!permDelete

	reqID := m.nextIoReqCnt()
	slog.Debug("Submitting delete request", "id", reqID, "items cnt", len(items))
	return func() tea.Msg {
		if len(items) == 0 {
			return NewDeleteOperationMsg(processbar.Cancelled, reqID)
		}
		if items[0].Provider == filesystem.ProviderLocal {
			paths := make([]string, len(items))
			for i, item := range items {
				paths[i] = item.Path.String()
			}
			return m.deleteOperation(&m.processBarModel, paths, useTrash, reqID)
		}
		return m.deleteProviderOperation(&m.processBarModel, items, reqID)
	}
}

func (m *model) deleteOperation(processBarModel *processbar.Model, items []string, useTrash bool, reqID int) tea.Msg {
	if len(items) == 0 {
		return NewDeleteOperationMsg(processbar.Cancelled, reqID)
	}
	p, err := processBarModel.SendAddProcessMsg(filepath.Base(items[0]), processbar.OpDelete, len(items), true)
	if err != nil {
		slog.Error("Cannot spawn a new process", "error", err)
		return NewDeleteOperationMsg(processbar.Failed, reqID)
	}
	finalizer := func(state processbar.ProcessState, reqID int) tea.Msg { return NewDeleteOperationMsg(state, reqID) }
	processor := makeDeleteProcessor(p, processBarModel, useTrash)
	msg := m.runFileProcessor(processor, finalizer, items, reqID)
	return msg
}

func (m *model) deleteProviderOperation(
	processBarModel *processbar.Model,
	items []filesystem.Location,
	reqID int,
) tea.Msg {
	if len(items) == 0 {
		return NewDeleteOperationMsg(processbar.Cancelled, reqID)
	}
	p, err := processBarModel.SendAddProcessMsg(pathBase(items[0].Path), processbar.OpDelete, len(items), true)
	if err != nil {
		slog.Error("Cannot spawn a new remote delete process", "error", err)
		return NewProviderDeleteOperationMsg(processbar.Failed, items[0], err, reqID)
	}
	resolveCtx, cancelResolve := mutationContext(items[0])
	session, err := m.ResolveFreshSession(resolveCtx, items[0])
	cancelResolve()
	if err != nil {
		p.State = processbar.Failed
		p.ErrorMsg = err.Error()
		markProcessDone(p, processBarModel)
		return NewProviderDeleteOperationMsg(processbar.Failed, items[0], err, reqID)
	}
	defer session.Close()
	for _, item := range items {
		deleteCtx, cancelDelete := mutationContext(item)
		deleteErr := session.Delete(
			deleteCtx,
			item.Path,
			filesystem.DeleteOptions{Recursive: true},
		)
		cancelDelete()
		if deleteErr != nil {
			p.State = processbar.Failed
			p.ErrorMsg = formatFileError(item.Path.String(), deleteErr)
			p.DoneTime = time.Now()
			_ = processBarModel.SendUpdateProcessMsg(p, true)
			return NewProviderDeleteOperationMsg(processbar.Failed, item, deleteErr, reqID)
		}
		p.CurrentFile = pathBase(item.Path)
		p.Done++
		processBarModel.TrySendingUpdateProcessMsg(p)
	}
	p.State = processbar.Successful
	markProcessDone(p, processBarModel)
	return NewProviderDeleteOperationMsg(processbar.Successful, items[0], nil, reqID)
}

func makeDeleteProcessor(process processbar.Process,
	processBarModel *processbar.Model,
	useTrash bool) processbar.FileListProcessor {
	processorFunction := func(items []string) (processbar.Process, []string) {
		notProcessed := make([]string, 0)
		if len(items) == 0 {
			markProcessDone(process, processBarModel)
			return process, notProcessed
		}
		deleteFunc := deleteElement
		if useTrash {
			deleteFunc = func(item string) error {
				_, err := trash.Move(item)
				return err
			}
		}
		for i, item := range items {
			err := deleteFunc(item)
			if err != nil {
				process.State = processbar.Failed
				slog.Error("Error in delete operation", "item", item, "useTrash", useTrash, "error", err)
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
	return processorFunction
}

func (m *model) getDeleteTriggerCmd(deletePermanent bool) tea.Cmd {
	panel := m.getFocusedFilePanel()
	if (panel.PanelMode == filepanel.SelectMode && panel.SelectedCount() == 0) ||
		(panel.PanelMode == filepanel.BrowserMode && panel.Empty()) {
		return nil
	}

	reqID := m.nextIoReqCnt()

	return func() tea.Msg {
		title := common.TrashWarnTitle
		content := common.TrashWarnContent
		action := notify.DeleteAction

		if panel.CurrentLocation().Provider != filesystem.ProviderLocal ||
			!m.hasTrash || !trash.Available(panel.Location) || deletePermanent {
			title = common.PermanentDeleteWarnTitle
			content = common.PermanentDeleteWarnContent
			action = notify.PermanentDeleteAction
		}
		return NewNotifyModalMsg(notify.New(true, title, content, action), reqID)
	}
}

// Copy directory or file's path to superfile's clipboard
// set cut to true/false accordingly
func (m *model) copySingleItem(cut bool) {
	panel := m.getFocusedFilePanel()
	m.clipboard.Reset(cut)
	if panel.Empty() {
		return
	}
	itemLocation := elementLocation(panel.CurrentLocation(), panel.GetFocusedItem())
	slog.Debug("handle_file_operations.copySingleItem", "cut", cut,
		"panel location", itemLocation.Path.String())
	m.clipboard.AddLocation(itemLocation)
}

// Copy all selected file or directory's paths to the clipboard
func (m *model) copyMultipleItem(cut bool) {
	panel := m.getFocusedFilePanel()
	m.clipboard.Reset(cut)
	if panel.SelectedCount() == 0 {
		return
	}
	items := focusedItemLocations(panel)
	slog.Debug("handle_file_operations.copyMultipleItem", "cut", cut,
		"panel selected files", items)
	m.clipboard.SetLocations(items)
}

func (m *model) getPasteItemCmd() tea.Cmd {
	copyItems := m.clipboard.PruneInaccessibleLocationsAndGet()
	cut := m.clipboard.IsCut()
	if len(copyItems) == 0 {
		return nil
	}

	reqID := m.nextIoReqCnt()
	panelLocation := m.getFocusedFilePanel().CurrentLocation()

	slog.Debug(
		"Submitting pasteItems request",
		"id",
		reqID,
		"items cnt",
		len(copyItems),
		"dest",
		panelLocation.Path.String(),
	)
	return func() tea.Msg {
		err := validatePasteOperation(panelLocation, copyItems, cut)
		if err != nil {
			title := "Invalid paste location"
			if errors.Is(err, filesystem.ErrUnsupported) {
				title = "Unsupported remote operation"
			}
			return NewNotifyModalMsg(notify.New(true, title, err.Error(), notify.NoAction), reqID)
		}
		allLocal := panelLocation.Provider == filesystem.ProviderLocal
		paths := make([]string, len(copyItems))
		for i, item := range copyItems {
			paths[i] = item.Path.String()
			allLocal = allLocal && item.Provider == filesystem.ProviderLocal
		}
		if allLocal {
			return m.executePasteOperation(&m.processBarModel, panelLocation.Path.String(), paths, cut, reqID)
		}
		return m.executeProviderPasteOperation(&m.processBarModel, panelLocation, copyItems, cut, reqID)
	}
}

func validatePasteOperation(destination filesystem.Location, copyItems []filesystem.Location, cut bool) error {
	for _, source := range copyItems {
		if err := filesystem.ValidateTransferTopology(source, destination); err != nil {
			if errors.Is(err, filesystem.ErrUnsupported) && source.Provider != filesystem.ProviderLocal {
				return errors.New(
					remoteUnsupportedOperationText(source.Provider, filesystem.OperationRemoteCrossSessionMove),
				)
			}
			return err
		}
		if cut && sameParentDirectory(source, destination) {
			return fmt.Errorf("cannot paste into parent directory of source, srcPath : %v, panelLocation : %v",
				source.Path.String(), destination.Path.String())
		}
		if cut && source.Path.String() == destination.Path.String() {
			return errors.New("cannot paste a directory into itself")
		}
		if isAncestorLocation(source, destination) {
			return fmt.Errorf("cannot %s and paste a directory into itself or its subdirectory",
				getCopyOrCutOperationName(cut))
		}
	}

	return nil
}

func makePasteProcessor(process processbar.Process,
	processBarModel *processbar.Model,
	panelLocation string, cut bool,
) processbar.FileListProcessor {
	processorFunction := func(items []string) (processbar.Process, []string) {
		notProcessed := make([]string, 0)
		if len(items) == 0 {
			markProcessDone(process, processBarModel)
			return process, notProcessed
		}
		var err error
		for i, filePath := range items {
			errMessage := "cut item error"
			if cut && !isExternalDiskPath(filePath) {
				err = moveElement(filePath, filepath.Join(panelLocation, filepath.Base(filePath)))
			} else {
				// TODO : These error cases are hard to test. We have to somehow make the paste operations fail,
				// which is time consuming and manual. We should test these with automated testcases
				// UPD: use "chattr +i" for target catalog to fail past opeations
				err = pasteDir(filePath, filepath.Join(panelLocation, filepath.Base(filePath)),
					&process, cut, processBarModel)
				if err != nil {
					errMessage = "paste item error"
				}
			}

			process.CurrentFile = filepath.Base(filePath)
			if err != nil {
				process.State = processbar.Failed
				slog.Error(errMessage, "error", err)
				slog.Debug("model.pasteItem - paste failure", "error", err,
					"current item", filePath, "errMessage", errMessage)
				process.ErrorMsg = formatFileError(filePath, err)
				notProcessed = items[i:]
				break
			}
			processBarModel.TrySendingUpdateProcessMsg(process)
		}
		if process.State != processbar.Failed {
			process.State = processbar.Successful
			process.Done = process.Total
			markProcessDone(process, processBarModel)
		}
		return process, notProcessed
	}
	return processorFunction
}

func (m *model) executePasteOperation(processBarModel *processbar.Model,
	panelLocation string, items []string, cut bool, reqID int,
) tea.Msg {
	if len(items) == 0 {
		return NewPasteOperationMsg(processbar.Cancelled, reqID)
	}
	var operation processbar.OperationType
	if cut {
		operation = processbar.OpCut
	} else {
		operation = processbar.OpCopy
	}

	p, err := processBarModel.SendAddProcessMsg(
		filepath.Base(items[0]),
		operation,
		getTotalFilesCnt(items), true)
	if err != nil {
		slog.Error("Cannot spawn a new process", "error", err)
		return NewPasteOperationMsg(processbar.Failed, reqID)
	}
	finalizer := func(state processbar.ProcessState, reqId int) tea.Msg { return NewPasteOperationMsg(state, reqId) }
	processor := makePasteProcessor(p, processBarModel, panelLocation, cut)
	msg := m.runFileProcessor(processor, finalizer, items, reqID)
	return msg
}

func (m *model) executeProviderPasteOperation(processBarModel *processbar.Model,
	destination filesystem.Location, items []filesystem.Location, cut bool, reqID int,
) tea.Msg {
	if len(items) == 0 {
		return NewPasteOperationMsg(processbar.Cancelled, reqID)
	}
	engine := filesystem.NewTransferEngine(m)
	operation := filesystem.OperationCopy
	if cut {
		operation = filesystem.OperationCutMove
	}
	refreshLocations := append(append([]filesystem.Location(nil), items...), destination)
	for i, source := range items {
		remainingSources := items[i:]
		resolveCtx, cancelResolve := mutationContext(destination)
		targetSession, err := m.ResolveFreshSession(resolveCtx, destination)
		if err != nil {
			cancelResolve()
			slog.Error("Cannot resolve destination session for provider paste", "error", err)
			return NewProviderPasteOperationMsg(
				processbar.Failed, destination, refreshLocations, remainingSources, err, reqID,
			)
		}
		target := locationWithPath(destination, pathJoin(destination.Path, pathBase(source.Path)))
		target, err = renameLocationIfDuplicate(resolveCtx, targetSession, target)
		_ = targetSession.Close()
		cancelResolve()
		if err != nil {
			slog.Error("Cannot resolve duplicate destination for provider paste", "error", err)
			return NewProviderPasteOperationMsg(
				processbar.Failed, destination, refreshLocations, remainingSources, err, reqID,
			)
		}
		transfer, err := engine.Start(context.Background(), filesystem.TransferRequest{
			Operation:   operation,
			Source:      source,
			Destination: target,
		})
		if err != nil {
			slog.Error("Cannot start provider paste transfer", "error", err)
			return NewProviderPasteOperationMsg(
				processbar.Failed,
				transferFailureLocation(err, source, target),
				refreshLocations,
				remainingSources,
				err,
				reqID,
			)
		}
		if _, err = filesystem.TrackTransferProcess(context.Background(), processBarModel, transfer); err != nil {
			slog.Error("Cannot track provider paste transfer", "error", err)
			_ = transfer.Cancel(context.Background())
			return NewProviderPasteOperationMsg(
				processbar.Failed,
				transferFailureLocation(err, source, target),
				refreshLocations,
				remainingSources,
				err,
				reqID,
			)
		}
		if err = transfer.Wait(context.Background()); err != nil {
			slog.Error("Provider paste transfer failed", "error", err)
			return NewProviderPasteOperationMsg(
				processbar.Failed,
				transferFailureLocation(err, source, target),
				refreshLocations,
				remainingSources,
				err,
				reqID,
			)
		}
	}
	return NewProviderPasteOperationMsg(
		processbar.Successful, filesystem.Location{}, refreshLocations, nil, nil, reqID,
	)
}

func transferFailureLocation(
	err error,
	source filesystem.Location,
	destination filesystem.Location,
) filesystem.Location {
	var operationErr *filesystem.OperationError
	if errors.As(err, &operationErr) {
		if operationErr.Path == destination.Path {
			return destination
		}
		if operationErr.Path == source.Path {
			return source
		}
		if source.Provider != destination.Provider {
			if operationErr.Provider == destination.Provider {
				return destination
			}
			if operationErr.Provider == source.Provider {
				return source
			}
		}
	}
	return source
}

func getTotalFilesCnt(copyItems []string) int {
	totalFiles := 0
	for _, folderPath := range copyItems {
		// TODO : Fix this. This is inefficient
		// In case of a cut operations for a directory with a lot of files
		// we are unnecessarily walking the whole directory recursively
		// while os will just perform a rename
		// So instead of few operations this will cause the cut paste
		// to read the whole directory recursively
		// we should avoid doing this.
		// Although this allows us a more detailed progress tracking
		// this make the copy/cut more inefficient
		// instead, we could just track progress based on total items in
		// copyItems
		// efficiency should be prioritized over more detailed feedback.
		count, err := countFiles(folderPath)
		if err != nil {
			slog.Error("Error in countFiles", "error", err)
			continue
		}
		totalFiles += count
	}
	return totalFiles
}

// Extract compressed file
// TODO : err should be returned and properly handled by the caller
func (m *model) getExtractFileCmd() tea.Cmd {
	panel := m.getFocusedFilePanel()
	if panel.Empty() {
		return nil
	}
	if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
		return m.unsupportedRemoteOperationCmd(panel.CurrentLocation(), filesystem.OperationExtract)
	}

	item := panel.GetFocusedItem().Location

	ext := strings.ToLower(filepath.Ext(item))
	if !common.IsExtensionExtractable(ext) {
		slog.Error("Error unexpected file", "extension type", ext, "item", item, "error", errors.ErrUnsupported)
		return nil
	}
	reqID := m.nextIoReqCnt()

	slog.Debug("Submitting Extract file request", "reqID", reqID, "item", item)

	return func() tea.Msg {
		outputDir := common.FileNameWithoutExtension(item)
		outputDir, err := renameIfDuplicate(outputDir)
		if err != nil {
			slog.Error("Error while renaming for duplicates", "error", err)
			return NewExtractOperationMsg(processbar.Failed, reqID)
		}

		err = os.MkdirAll(
			outputDir,
			utils.ExtractedDirMode,
		)
		if err != nil {
			slog.Error("Error while making directory for extracting files", "error", err)
			return NewExtractOperationMsg(processbar.Failed, reqID)
		}
		err = extractCompressFile(item, outputDir, &m.processBarModel)
		if err != nil {
			slog.Error("Error extract file", "error", err)
			return NewExtractOperationMsg(processbar.Failed, reqID)
		}
		return NewExtractOperationMsg(processbar.Successful, reqID)
	}
}

func (m *model) getCompressSelectedFilesCmd() tea.Cmd {
	panel := m.getFocusedFilePanel()

	if panel.Empty() {
		return nil
	}
	if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
		return m.unsupportedRemoteOperationCmd(panel.CurrentLocation(), filesystem.OperationCompress)
	}
	var filesToCompress []string
	var firstFile string

	if panel.SelectedCount() == 0 {
		firstFile = panel.GetFocusedItem().Location
		filesToCompress = append(filesToCompress, firstFile)
	} else {
		firstFile = panel.GetFirstSelectedLocation()
		filesToCompress = panel.GetSelectedLocationsSortedAsVisible()
	}

	reqID := m.nextIoReqCnt()

	return func() tea.Msg {
		zipName, err := getZipArchiveName(filepath.Base(firstFile))
		if err != nil {
			slog.Error("Error in getZipArchiveName", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		zipPath := filepath.Join(panel.Location, zipName)
		if err := zipSources(filesToCompress, zipPath, &m.processBarModel); err != nil {
			slog.Error("Error in zipping files", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		return NewCompressOperationMsg(processbar.Successful, reqID)
	}
}

func (m *model) chooserFileWriteAndQuit(path string) error {
	// Attempt to write to the file
	err := os.WriteFile(variable.ChooserFile, []byte(path), utils.ConfigFilePerm)
	if err != nil {
		return err
	}
	m.modelQuitState = quitInitiated
	return nil
}

// Open file with default editor
func (m *model) openFileWithEditor() tea.Cmd {
	panel := m.getFocusedFilePanel()
	// Check if panel is empty
	if panel.Empty() {
		return nil
	}
	if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
		return m.unsupportedRemoteOperationCmd(panel.CurrentLocation(), filesystem.OperationOpenWith)
	}

	if variable.ChooserFile != "" {
		err := m.chooserFileWriteAndQuit(panel.GetFocusedItem().Location)
		if err == nil {
			return nil
		}
		// Continue with preview if file is not writable
		slog.Error("Error while writing to chooser file, continuing with open via file editor", "error", err)
	}

	editor := common.Config.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	// Make sure there is an editor
	if editor == "" {
		if runtime.GOOS == utils.OsWindows {
			editor = "notepad"
		} else {
			editor = "nano"
		}
	}

	// Split the editor command into command and arguments
	parts := strings.Fields(editor)
	cmd := parts[0]

	//nolint:gocritic // appendAssign: intentionally creating a new slice
	args := append(parts[1:], panel.GetFocusedItem().Location)

	c := exec.Command(cmd, args...) //nolint:gosec // Editor command is intentionally user-configurable.

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Open directory with default editor
func (m *model) openDirectoryWithEditor() tea.Cmd {
	if m.getFocusedFilePanel().CurrentLocation().Provider != filesystem.ProviderLocal {
		return m.unsupportedRemoteOperationCmd(m.getFocusedFilePanel().CurrentLocation(), filesystem.OperationOpenWith)
	}
	if variable.ChooserFile != "" {
		err := m.chooserFileWriteAndQuit(m.getFocusedFilePanel().Location)
		if err == nil {
			return nil
		}
		// Continue with preview if file is not writable
		slog.Error("Error while writing to chooser file, continuing with open via directory editor", "error", err)
	}

	editor := common.Config.DirEditor

	if editor == "" {
		switch runtime.GOOS {
		case utils.OsWindows:
			editor = "explorer"
		case utils.OsDarwin:
			editor = "open"
		default:
			editor = "vi"
		}
	}

	// Split the editor command into command and arguments
	parts := strings.Fields(editor)
	cmd := parts[0]
	//nolint:gocritic // appendAssign: intentionally creating a new slice
	args := append(parts[1:], m.getFocusedFilePanel().Location)

	c := exec.Command(cmd, args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Copy file path
// TODO: This is also an IO operations, do it via tea.Cmd
func (m *model) copyPath() {
	pathText := m.copyPathText()
	if pathText == "" {
		return
	}

	if err := m.writeClipboard(pathText); err != nil {
		slog.Error("Error while copy path", "error", err)
	}
}

func (m *model) copyPathText() string {
	panel := m.getFocusedFilePanel()

	if panel.Empty() {
		return ""
	}

	if panel.PanelMode == filepanel.SelectMode && panel.SelectedCount() > 0 {
		return strings.Join(panel.GetSelectedLocationsSortedAsVisible(), "\n")
	}

	return panel.GetFocusedItem().Location
}

// TODO: This is also an IO operations, do it via tea.Cmd
func (m *model) copyPWD() {
	panel := m.getFocusedFilePanel()
	if err := m.writeClipboard(panel.Location); err != nil {
		slog.Error("Error while copy present working directory", "error", err)
	}
}

func (m *model) writeClipboard(text string) error {
	if m.clipboardWriter != nil {
		return m.clipboardWriter(text)
	}

	return clipboard.WriteAll(text)
}
