package internal

import (
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
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// Create a file in the currently focus file panel
// TODO: Fix it. It doesn't creates a new file. It just opens a file model,
// that allows you to create a file. Actual creation happens here - createItem() in handle_modal.go
func (m *model) panelCreateNewFile() {
	panel := m.getFocusedFilePanel()

	m.typingModal.location = panel.Location
	m.typingModal.open = true
	m.typingModal.textInput = common.GenerateNewFileTextInput()
	m.firstTextInput = true
}

// TODO : This function does not needs the entire model. Only pass the panel object
func (m *model) IsRenamingConflicting() bool {
	// TODO : Replace this with m.getCurrentFilePanel() everywhere
	panel := m.getFocusedFilePanel()
	if panel.ElemCount() == 0 {
		slog.Error("IsRenamingConflicting() being called on empty panel")
		return false
	}
	oldPath := panel.GetFocusedItem().Location
	newPath := filepath.Join(panel.Location, panel.Rename.Value())

	if oldPath == newPath {
		return false
	}

	_, err := os.Stat(newPath)
	return err == nil
}

// TODO: Remove channel messaging and use tea.Cmd
func (m *model) warnModalForRenaming() tea.Cmd {
	reqID := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting rename notify model request", "reqID", reqID)
	res := func() tea.Msg {
		notifyModel := notify.New(true,
			common.SameRenameWarnTitle,
			common.SameRenameWarnContent,
			notify.RenameAction)
		return NewNotifyModalMsg(notifyModel, reqID)
	}
	return res
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

	var items []string
	if panel.PanelMode == filepanel.SelectMode {
		items = panel.GetSelectedLocations()
	} else {
		items = []string{panel.GetFocusedItem().Location}
	}

	useTrash := m.hasTrash && !isExternalDiskPath(panel.Location) && !permDelete

	reqID := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting delete request", "id", reqID, "items cnt", len(items))
	return func() tea.Msg {
		state := deleteOperation(&m.processBarModel, items, useTrash)
		return NewDeleteOperationMsg(state, reqID)
	}
}

func deleteOperation(processBarModel *processbar.Model, items []string, useTrash bool) processbar.ProcessState {
	if len(items) == 0 {
		return processbar.Cancelled
	}
	p, err := processBarModel.SendAddProcessMsg(filepath.Base(items[0]), processbar.OpDelete, len(items), true)
	if err != nil {
		slog.Error("Cannot spawn a new process", "error", err)
		return processbar.Failed
	}

	deleteFunc := os.RemoveAll
	if useTrash {
		deleteFunc = moveToTrash
	}
	for _, item := range items {
		err = deleteFunc(item)
		if err != nil {
			p.State = processbar.Failed
			slog.Error("Error in delete operation", "item", item, "useTrash", useTrash, "error", err)
			break
		}
		p.CurrentFile = filepath.Base(item)
		p.Done++
		processBarModel.TrySendingUpdateProcessMsg(p)
	}

	if p.State != processbar.Failed {
		p.State = processbar.Successful
	}
	p.DoneTime = time.Now()
	err = processBarModel.SendUpdateProcessMsg(p, true)
	if err != nil {
		slog.Error("Failed to send final delete operation update", "error", err)
	}
	return p.State
}

func (m *model) getDeleteTriggerCmd(deletePermanent bool) tea.Cmd {
	panel := m.getFocusedFilePanel()
	if (panel.PanelMode == filepanel.SelectMode && panel.SelectedCount() == 0) ||
		(panel.PanelMode == filepanel.BrowserMode && panel.Empty()) {
		return nil
	}

	reqID := m.ioReqCnt
	m.ioReqCnt++

	return func() tea.Msg {
		title := common.TrashWarnTitle
		content := common.TrashWarnContent
		action := notify.DeleteAction

		if !m.hasTrash || isExternalDiskPath(panel.Location) || deletePermanent {
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
	slog.Debug("handle_file_operations.copySingleItem", "cut", cut,
		"panel location", panel.GetFocusedItem().Location)
	m.clipboard.Add(panel.GetFocusedItem().Location)
}

// Copy all selected file or directory's paths to the clipboard
func (m *model) copyMultipleItem(cut bool) {
	panel := m.getFocusedFilePanel()
	m.clipboard.Reset(cut)
	if panel.SelectedCount() == 0 {
		return
	}
	slog.Debug("handle_file_operations.copyMultipleItem", "cut", cut,
		"panel selected files", panel.GetSelectedLocations())
	m.clipboard.SetItems(panel.GetSelectedLocations())
}

func (m *model) getPasteItemCmd() tea.Cmd {
	copyItems := m.clipboard.PruneInaccessibleItemsAndGet()
	cut := m.clipboard.IsCut()
	if len(copyItems) == 0 {
		return nil
	}

	// TODO: Do it via m.getNewReqID()
	// TODO: Have an IO Req Management, collecting info about pending IO Req too
	reqID := m.ioReqCnt
	m.ioReqCnt++
	panelLocation := m.getFocusedFilePanel().Location

	slog.Debug("Submitting pasteItems request", "id", reqID, "items cnt", len(copyItems), "dest", panelLocation)
	return func() tea.Msg {
		err := validatePasteOperation(panelLocation, copyItems, cut)
		if err != nil {
			return NewNotifyModalMsg(notify.New(true, "Invalid paste location", err.Error(), notify.NoAction),
				reqID)
		}
		state := executePasteOperation(&m.processBarModel, panelLocation, copyItems, cut)
		return NewPasteOperationMsg(state, reqID)
	}
}

func validatePasteOperation(panelLocation string, copyItems []string, cut bool) error {
	// Check if trying to paste into source or subdirectory for both cut and copy operations
	for _, srcPath := range copyItems {
		// Check if trying to cut and paste into the same directory - this would be a no-op
		// and could potentially cause issues, so we prevent it
		if filepath.Dir(srcPath) == panelLocation && cut {
			return fmt.Errorf("cannot paste into parent directory of source, srcPath : %v, panelLocation : %v",
				srcPath, panelLocation)
		}
		if cut && srcPath == panelLocation {
			return errors.New("cannot paste a directory into itself")
		}

		if isAncestor(srcPath, panelLocation) {
			return fmt.Errorf("cannot %s and paste a directory into itself or its subdirectory",
				getCopyOrCutOperationName(cut))
		}
	}

	return nil
}

// new func to check and return an error that will go in m.content
// create a new error type

// Paste all clipboard items
func executePasteOperation(processBarModel *processbar.Model,
	panelLocation string, copyItems []string, cut bool,
) processbar.ProcessState {
	slog.Debug("executePasteOperation", "items", copyItems, "cut", cut, "panel location", panelLocation)

	var operation processbar.OperationType
	if cut {
		operation = processbar.OpCut
	} else {
		operation = processbar.OpCopy
	}

	p, err := processBarModel.SendAddProcessMsg(
		filepath.Base(copyItems[0]),
		operation,
		getTotalFilesCnt(copyItems), true)
	if err != nil {
		slog.Error("Cannot spawn a new process", "error", err)
		return processbar.Failed
	}

	for _, filePath := range copyItems {
		errMessage := "cut item error"
		if cut && !isExternalDiskPath(filePath) {
			err = moveElement(filePath, filepath.Join(panelLocation, filepath.Base(filePath)))
		} else {
			// TODO : These error cases are hard to test. We have to somehow make the paste operations fail,
			// which is time consuming and manual. We should test these with automated testcases
			err = pasteDir(filePath, filepath.Join(panelLocation, filepath.Base(filePath)), &p, cut, processBarModel)
			if err != nil {
				errMessage = "paste item error"
			}
		}

		p.CurrentFile = filepath.Base(filePath)
		if err != nil {
			slog.Debug("model.pasteItem - paste failure", "error", err,
				"current item", filePath, "errMessage", errMessage)
			p.State = processbar.Failed
			slog.Error(errMessage, "error", err)
			break
		}
		processBarModel.TrySendingUpdateProcessMsg(p)
	}

	if p.State != processbar.Failed {
		p.State = processbar.Successful
		p.Done = p.Total
	}
	p.DoneTime = time.Now()
	err = processBarModel.SendUpdateProcessMsg(p, true)
	if err != nil {
		slog.Error("Could not send final update for process Bar", "error", err)
	}

	return p.State
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

	item := panel.GetFocusedItem().Location

	ext := strings.ToLower(filepath.Ext(item))
	if !common.IsExtensionExtractable(ext) {
		slog.Error("Error unexpected file", "extension type", ext, "item", item, "error", errors.ErrUnsupported)
		return nil
	}
	reqID := m.ioReqCnt
	m.ioReqCnt++

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
	var filesToCompress []string
	var firstFile string

	if panel.SelectedCount() == 0 {
		firstFile = panel.GetFocusedItem().Location
		filesToCompress = append(filesToCompress, firstFile)
	} else {
		firstFile = panel.GetFirstSelectedLocation()
		filesToCompress = panel.GetSelectedLocations()
	}

	reqID := m.ioReqCnt
	m.ioReqCnt++

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

	c := exec.Command(cmd, args...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Open directory with default editor
func (m *model) openDirectoryWithEditor() tea.Cmd {
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
	panel := m.getFocusedFilePanel()

	if panel.Empty() {
		return
	}

	if err := clipboard.WriteAll(panel.GetFocusedItem().Location); err != nil {
		slog.Error("Error while copy path", "error", err)
	}
}

// TODO: This is also an IO operations, do it via tea.Cmd
func (m *model) copyPWD() {
	panel := m.getFocusedFilePanel()
	if err := clipboard.WriteAll(panel.Location); err != nil {
		slog.Error("Error while copy present working directory", "error", err)
	}
}
