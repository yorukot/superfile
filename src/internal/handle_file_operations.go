package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/yorukot/superfile/src/internal/common"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

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

// Bulk rename selected files
func (m *model) panelBulkRename() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	// Check if there are selected files
	if panel.panelMode != selectMode || len(panel.selected) == 0 {
		return
	}

	// Initialize bulk rename modal
	m.bulkRenameModal.open = true
	m.bulkRenameModal.renameType = 0 // Default to find/replace
	m.bulkRenameModal.cursor = 0
	m.bulkRenameModal.startNumber = 1
	m.bulkRenameModal.caseType = 0
	m.bulkRenameModal.errorMessage = ""
	m.firstTextInput = true

	// Initialize text inputs
	m.bulkRenameModal.findInput = common.GenerateBulkRenameTextInput("Find text")
	m.bulkRenameModal.replaceInput = common.GenerateBulkRenameTextInput("Replace with")
	m.bulkRenameModal.prefixInput = common.GenerateBulkRenameTextInput("Add prefix")
	m.bulkRenameModal.suffixInput = common.GenerateBulkRenameTextInput("Add suffix")

	// Focus the first input based on rename type
	m.bulkRenameModal.findInput.Focus()
}

func (m *model) getDeleteCmd(permDelete bool) tea.Cmd {
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

	useTrash := m.hasTrash && !isExternalDiskPath(panel.location) && !permDelete

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
	p, err := processBarModel.SendAddProcessMsg(icon.Delete+icon.Space+filepath.Base(items[0]), len(items), true)
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
		p.Name = icon.Delete + icon.Space + filepath.Base(item)
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
	if (panel.panelMode == selectMode && len(panel.selected) == 0) ||
		(panel.panelMode == browserMode && len(panel.element) == 0) {
		return nil
	}

	reqID := m.ioReqCnt
	m.ioReqCnt++

	return func() tea.Msg {
		title := common.TrashWarnTitle
		content := common.TrashWarnContent
		action := notify.DeleteAction

		if !m.hasTrash || isExternalDiskPath(panel.location) || deletePermanent {
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
	panelLocation string, copyItems []string, cut bool) processbar.ProcessState {
	slog.Debug("executePasteOperation", "items", copyItems, "cut", cut, "panel location", panelLocation)

	p, err := processBarModel.SendAddProcessMsg(
		icon.GetCopyOrCutIcon(cut)+icon.Space+filepath.Base(copyItems[0]),
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
			} else if cut {
				//TODO: Fix unhandled error
				os.RemoveAll(filePath)
			}
		}

		p.Name = icon.GetCopyOrCutIcon(cut) + icon.Space + filepath.Base(filePath)
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
	if len(panel.element) == 0 {
		return nil
	}

	item := panel.getSelectedItem().location

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
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}

		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			slog.Error("Error while making directory for extracting files", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		err = extractCompressFile(item, outputDir, &m.processBarModel)
		if err != nil {
			slog.Error("Error extract file", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		return NewCompressOperationMsg(processbar.Successful, reqID)
	}
}

func (m *model) getCompressSelectedFilesCmd() tea.Cmd {
	panel := m.getFocusedFilePanel()

	if len(panel.element) == 0 {
		return nil
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

	reqID := m.ioReqCnt
	m.ioReqCnt++

	return func() tea.Msg {
		zipName, err := getZipArchiveName(filepath.Base(firstFile))
		if err != nil {
			slog.Error("Error in getZipArchiveName", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		zipPath := filepath.Join(panel.location, zipName)
		if err := zipSources(filesToCompress, zipPath, &m.processBarModel); err != nil {
			slog.Error("Error in zipping files", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		return NewCompressOperationMsg(processbar.Successful, reqID)
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
// TODO: This is also an IO operations, do it via tea.Cmd
func (m *model) copyPath() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}

	if err := clipboard.WriteAll(panel.element[panel.cursor].location); err != nil {
		slog.Error("Error while copy path", "error", err)
	}
}

// TODO: This is also an IO operations, do it via tea.Cmd
func (m *model) copyPWD() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if err := clipboard.WriteAll(panel.location); err != nil {
		slog.Error("Error while copy present working directory", "error", err)
	}
}

// Bulk rename helper functions

// applyFindReplace applies find and replace to a filename
func applyFindReplace(filename, find, replace string) string {
	if find == "" {
		return filename
	}
	return strings.ReplaceAll(filename, find, replace)
}

// applyPrefix adds a prefix to a filename (before extension)
func applyPrefix(filename, prefix string) string {
	if prefix == "" {
		return filename
	}
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return prefix + nameWithoutExt + ext
}

// applySuffix adds a suffix to a filename (before extension)
func applySuffix(filename, suffix string) string {
	if suffix == "" {
		return filename
	}
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + suffix + ext
}

// applyNumbering adds a number to a filename
func applyNumbering(filename string, number int) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + "_" + strconv.Itoa(number) + ext
}

// applyCaseConversion converts filename case (0: lowercase, 1: uppercase, 2: title case)
func applyCaseConversion(filename string, caseType int) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	switch caseType {
	case 0: // lowercase
		return strings.ToLower(nameWithoutExt) + ext
	case 1: // uppercase
		return strings.ToUpper(nameWithoutExt) + ext
	case 2: // title case
		// Simple title case implementation (capitalize first letter of each word)
		words := strings.Fields(strings.ToLower(nameWithoutExt))
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(string(word[0])) + word[1:]
			}
		}
		return strings.Join(words, " ") + ext
	default:
		return filename
	}
}

// generateBulkRenamePreview generates preview for bulk rename
func (m *model) generateBulkRenamePreview() []renamePreview {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	previews := make([]renamePreview, 0, len(panel.selected))

	for i, itemPath := range panel.selected {
		oldName := filepath.Base(itemPath)
		newName := oldName
		var err string

		// Apply transformation based on rename type
		switch m.bulkRenameModal.renameType {
		case 0: // Find & Replace
			newName = applyFindReplace(newName, m.bulkRenameModal.findInput.Value(), m.bulkRenameModal.replaceInput.Value())
		case 1: // Prefix
			newName = applyPrefix(newName, m.bulkRenameModal.prefixInput.Value())
		case 2: // Suffix
			newName = applySuffix(newName, m.bulkRenameModal.suffixInput.Value())
		case 3: // Numbering
			newName = applyNumbering(newName, m.bulkRenameModal.startNumber+i)
		case 4: // Case conversion
			newName = applyCaseConversion(newName, m.bulkRenameModal.caseType)
		}

		// Validate new name
		if newName == "" {
			err = "Empty filename"
		} else if newName == oldName {
			err = "No change"
		} else {
			// Check if new name would cause conflict
			newPath := filepath.Join(filepath.Dir(itemPath), newName)
			if _, statErr := os.Stat(newPath); statErr == nil && newPath != itemPath {
				err = "File already exists"
			}
		}

		previews = append(previews, renamePreview{
			oldName: oldName,
			newName: newName,
			error:   err,
		})
	}

	return previews
}
