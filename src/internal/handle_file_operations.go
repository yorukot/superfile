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
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/yorukot/superfile/src/internal/common"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
)

// Create a file in the currently focus file panel
// TODO: Fix it. It doesn't creates a new file. It just opens a file model,
// that allows you to create a file. Actual creation happens here - createItem() in handle_modal.go
func (m *model) panelCreateNewFile() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	m.typingModal.location = panel.location
	m.typingModal.open = true
	m.typingModal.textInput = common.GenerateNewFileTextInput()
	m.firstTextInput = true
}

// TODO : This function does not needs the entire model. Only pass the panel object
func (m *model) IsRenamingConflicting() bool {
	// TODO : Replace this with m.getCurrentFilePanel() everywhere
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		slog.Error("IsRenamingConflicting() being called on empty panel")
		return false
	}

	oldPath := panel.element[panel.cursor].location
	newPath := filepath.Join(panel.location, panel.rename.Value())

	if oldPath == newPath {
		return false
	}

	_, err := os.Stat(newPath)
	return err == nil
}

// TODO: Remove channel messaging and use tea.Cmd
func (m *model) warnModalForRenaming() {
	id := shortuuid.New()
	message := channelMessage{
		messageID:   id,
		messageType: sendWarnModal,
	}

	message.warnModal = warnModal{
		open:     true,
		title:    "There is already a file or directory with that name",
		content:  "This operation will override the existing file",
		warnType: confirmRenameItem,
	}
	channel <- message
}

// Rename file where the cusror is located
// TODO: Fix this. It doesn't do any rename, just opens the rename text input
// Actual rename happens at confirmRename() in handle_modal.go
func (m *model) panelItemRename() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return
	}

	cursorPos := strings.LastIndex(panel.element[panel.cursor].name, ".")
	nameLen := len(panel.element[panel.cursor].name)
	if cursorPos == -1 || cursorPos == 0 && nameLen > 0 || panel.element[panel.cursor].directory {
		cursorPos = nameLen
	}

	m.fileModel.renaming = true
	panel.renaming = true
	m.firstTextInput = true
	panel.rename = common.GenerateRenameTextInput(m.fileModel.width-4, cursorPos, panel.element[panel.cursor].name)
}

func (m *model) getDeleteCmd() tea.Cmd {
	panel := m.getFocusedFilePanel()
	if len(panel.element) == 0 {
		return nil
	}

	var items []string
	if panel.panelMode == selectMode {
		items = panel.selected
	} else {
		items = []string{panel.getSelectedItem().location}
	}

	useTrash := hasTrash && !isExternalDiskPath(panel.location)

	reqID := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting delete request", "id", reqID, "items cnt", len(items))
	return func() tea.Msg {
		state := deleteOperation(&m.processBarModel, items, useTrash)
		return NewDeleteOperationMsg(state, reqID)
	}
}

func deleteOperation(processBarModel *processBarModel, items []string, useTrash bool) processState {
	if len(items) == 0 {
		return cancel
	}
	id := shortuuid.New()
	prog := common.GenerateDefaultProgress()

	p := process{
		name:     icon.Delete + icon.Space + filepath.Base(items[0]),
		progress: prog,
		state:    inOperation,
		total:    len(items),
		done:     0,
	}
	processBarModel.process[id] = p
	message := channelMessage{
		messageID:       id,
		messageType:     sendProcess,
		processNewState: p,
	}
	channel <- message
	deleteFunc := os.RemoveAll
	if useTrash {
		deleteFunc = trashMacOrLinux
	}
	for _, item := range items {
		err := deleteFunc(item)

		if err != nil {
			p.state = failure
			slog.Error("Error in delete operation", "item", item, "useTrash", useTrash, "error", err)
			break
		}
		p.name = icon.Delete + icon.Space + filepath.Base(item)
		p.done++
		if p.done == p.total {
			p.state = successful
			p.doneTime = time.Now()
		} else {
			p.state = inOperation
			if len(channel) < 5 {
				message.processNewState = p
				channel <- message
			}
		}
	}
	message.processNewState = p
	channel <- message
	processBarModel.process[id] = p
	return p.state
}

func (m *model) deleteItemWarn() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if panel.panelMode == browserMode && len(panel.element) == 0 || panel.panelMode == selectMode && len(panel.selected) == 0 {
		return
	}

	id := shortuuid.New()
	message := channelMessage{
		messageID:   id,
		messageType: sendWarnModal,
	}

	if !hasTrash || isExternalDiskPath(panel.location) {
		message.warnModal = warnModal{
			open:     true,
			title:    "Are you sure you want to completely delete",
			content:  "This operation cannot be undone and your data will be completely lost.",
			warnType: confirmDeleteItem,
		}
		channel <- message
		return
	}
	message.warnModal = warnModal{
		open:     true,
		title:    "Are you sure you want to move this to trash can",
		content:  "This operation will move file or directory to trash can.",
		warnType: confirmDeleteItem,
	}
	channel <- message
}

// Copy directory or file's path to superfile's clipboard
// set cut to true/false accordingly
func (m *model) copySingleItem(cut bool) {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.reset(cut)
	if len(panel.element) == 0 {
		return
	}
	slog.Debug("handle_file_operations.copySingleItem", "cut", cut,
		"panel location", panel.element[panel.cursor].location)
	m.copyItems.items = append(m.copyItems.items, panel.element[panel.cursor].location)
}

// Copy all selected file or directory's paths to the clipboard
func (m *model) copyMultipleItem(cut bool) {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.reset(cut)
	if len(panel.selected) == 0 {
		return
	}
	slog.Debug("handle_file_operations.copyMultipleItem", "cut", cut,
		"panel selected files", panel.selected)
	m.copyItems.items = panel.selected
}

func (m *model) getPasteItemCmd() tea.Cmd {
	copyItems := m.copyItems.items
	cut := m.copyItems.cut
	if len(copyItems) == 0 {
		return nil
	}

	//TODO: Do it via m.getNewReqID()
	//TODO: Have an IO Req Management, collecting info about pending IO Req too
	reqID := m.ioReqCnt
	m.ioReqCnt++
	panelLocation := m.getFocusedFilePanel().location

	slog.Debug("Submitting pasteItems request", "id", reqID, "items cnt", len(copyItems), "dest", panelLocation)
	return func() tea.Msg {
		state := executePasteOperation(&m.processBarModel, panelLocation, copyItems, cut)
		return NewPasteOperationMsg(state, reqID)
	}
}

// Paste all clipboard items
func executePasteOperation(processBarModel *processBarModel,
	panelLocation string, copyItems []string, cut bool) processState {
	if len(copyItems) == 0 {
		return cancel
	}

	id := shortuuid.New()
	totalFiles := 0

	// Check if trying to paste into source or subdirectory for both cut and copy operations
	for _, srcPath := range copyItems {
		// Check if trying to cut and paste into the same directory - this would be a no-op
		// and could potentially cause issues, so we prevent it
		if filepath.Dir(srcPath) == panelLocation && cut {
			slog.Error("Cannot paste into parent directory of source", "src", srcPath, "dst", panelLocation)
			message := channelMessage{
				messageID:   id,
				messageType: sendNotifyModal,
				notifyModal: notifyModal{
					open:    true,
					title:   "Invalid paste location",
					content: "Cannot paste into parent directory of source",
				},
			}
			channel <- message
			return cancel
		}

		slog.Debug("model.pasteItem", "srcPath", srcPath, "panel location", panelLocation)

		if cut && srcPath == panelLocation {
			slog.Error("Cannot paste a directory into itself", "operation", "cut", "src", srcPath, "dst", panelLocation)
			message := channelMessage{
				messageID:   id,
				messageType: sendNotifyModal,
				notifyModal: notifyModal{
					open:    true,
					title:   "Invalid paste location",
					content: "Cannot paste a directory into itself",
				},
			}
			channel <- message
			return cancel
		}

		if isAncestor(srcPath, panelLocation) {
			operation := "copy"
			if cut {
				operation = "cut"
			}

			slog.Error("Cannot paste a directory into itself or its subdirectory", "operation", operation, "src", srcPath, "dst", panelLocation)
			message := channelMessage{
				messageID:   id,
				messageType: sendNotifyModal,
				notifyModal: notifyModal{
					open:    true,
					title:   "Invalid paste location",
					content: fmt.Sprintf("Cannot %s and paste a directory into itself or its subdirectory", operation),
				},
			}
			channel <- message
			return cancel
		}
	}

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
			slog.Error("mode.pasteItem - Error in countFiles", "error", err)
			continue
		}
		totalFiles += count
	}

	slog.Debug("model.pasteItem", "items", copyItems, "cut", cut,
		"totalFiles", totalFiles, "panel location", panelLocation)

	prog := common.GenerateDefaultProgress()

	prefixIcon := icon.Copy + icon.Space
	if cut {
		prefixIcon = icon.Cut + icon.Space
	}

	newProcess := process{
		name:     prefixIcon + filepath.Base(copyItems[0]),
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}

	processBarModel.process[id] = newProcess

	message := channelMessage{
		messageID:       id,
		messageType:     sendProcess,
		processNewState: newProcess,
	}

	channel <- message

	//TODO: This is wrong. We are mutating this model object, that will
	// get changed after Update()
	// Mutation of model should not happen in this function
	p := processBarModel.process[id]
	for _, filePath := range copyItems {
		var err error
		prefixIcon := icon.Copy
		if cut {
			prefixIcon = icon.Cut
		}
		p.name = prefixIcon + icon.Space + filepath.Base(filePath)

		errMessage := "cut item error"
		if cut && !isExternalDiskPath(filePath) {
			err = moveElement(filePath, filepath.Join(panelLocation, filepath.Base(filePath)))
		} else {
			// TODO : These error cases are hard to test. We have to somehow make the paste operations fail,
			// which is time consuming and manual. We should test these with automated testcases
			err = pasteDir(filePath, filepath.Join(panelLocation, filepath.Base(filePath)), id, cut, processBarModel)
			if err != nil {
				errMessage = "paste item error"
			} else if cut {
				//TODO: Fix unhandled error
				os.RemoveAll(filePath)
			}
		}
		p = processBarModel.process[id]
		if err != nil {
			slog.Debug("model.pasteItem - paste failure", "error", err,
				"current item", filePath, "errMessage", errMessage)
			p.state = failure
			message.processNewState = p
			channel <- message
			slog.Error(errMessage, "error", err)
			processBarModel.process[id] = p
			break
		}
	}

	if p.state != failure {
		p.state = successful
		p.done = totalFiles
		p.doneTime = time.Now()
	}
	message.processNewState = p
	channel <- message

	processBarModel.process[id] = p
	return p.state
}

// Extract compressed file
// TODO : err should be returned and properly handled by the caller
func (m *model) extractFile() {
	var err error
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}

	ext := strings.ToLower(filepath.Ext(panel.element[panel.cursor].location))
	if !common.IsExtensionExtractable(ext) {
		slog.Error(fmt.Sprintf("Error unexpected file extension type: %s", ext), "error", errors.ErrUnsupported)
		return
	}

	outputDir := common.FileNameWithoutExtension(panel.element[panel.cursor].location)
	outputDir, err = renameIfDuplicate(outputDir)
	if err != nil {
		slog.Error("Error extract file when create new directory", "error", err)
		return
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		slog.Error("Error while making directory for extracting files", "error", err)
		return
	}
	err = extractCompressFile(panel.element[panel.cursor].location, outputDir)
	if err != nil {
		slog.Error("Error extract file", "error", err)
		return
	}
}

func (m *model) compressSelectedFiles() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}
	var filesToCompress []string
	var firstFile string

	if len(panel.selected) == 0 {
		firstFile = panel.element[panel.cursor].location
		filesToCompress = append(filesToCompress, firstFile)
	} else {
		firstFile = panel.selected[0]
		filesToCompress = panel.selected
	}
	zipName, err := getZipArchiveName(filepath.Base(firstFile))
	if err != nil {
		slog.Error("Error in getZipArchiveName", "error", err)
		return
	}
	zipPath := filepath.Join(panel.location, zipName)
	if err := zipSources(filesToCompress, zipPath); err != nil {
		slog.Error("Error in zipping files", "error", err)
		return
	}
}

func (m *model) chooserFileWriteAndQuit(path string) error {
	// Attempt to write to the file
	err := os.WriteFile(variable.ChooserFile, []byte(path), 0644)
	if err != nil {
		return err
	}
	m.modelQuitState = quitInitiated
	return nil
}

// Open file with default editor
func (m *model) openFileWithEditor() tea.Cmd {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	// Check if panel is empty
	if len(panel.element) == 0 {
		return nil
	}

	if variable.ChooserFile != "" {
		err := m.chooserFileWriteAndQuit(panel.element[panel.cursor].location)
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
	args := append(parts[1:], panel.element[panel.cursor].location)

	c := exec.Command(cmd, args...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Open directory with default editor
func (m *model) openDirectoryWithEditor() tea.Cmd {
	if variable.ChooserFile != "" {
		err := m.chooserFileWriteAndQuit(m.fileModel.filePanels[m.filePanelFocusIndex].location)
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
	args := append(parts[1:], m.fileModel.filePanels[m.filePanelFocusIndex].location)

	c := exec.Command(cmd, args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Copy file path
func (m *model) copyPath() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}

	if err := clipboard.WriteAll(panel.element[panel.cursor].location); err != nil {
		slog.Error("Error while copy path", "error", err)
	}
}

func (m *model) copyPWD() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if err := clipboard.WriteAll(panel.location); err != nil {
		slog.Error("Error while copy present working directory", "error", err)
	}
}
