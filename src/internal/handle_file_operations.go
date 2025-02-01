package internal

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
)

// Create a file in the currently focus file panel
func (m *model) panelCreateNewFile() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	ti := textinput.New()
	ti.Cursor.Style = modalCursorStyle
	ti.Cursor.TextStyle = modalStyle
	ti.TextStyle = modalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "Add \"/\" represent Transcend folder at the end"
	ti.PlaceholderStyle = modalStyle
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = modalWidth - 10

	m.typingModal.location = panel.location
	m.typingModal.open = true
	m.typingModal.textInput = ti
	m.firstTextInput = true

}

func (m *model) IsRenamingConflicting() bool {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	oldPath := panel.element[panel.cursor].location
	newPath := panel.location + "/" + panel.rename.Value()

	if oldPath == newPath {
		return false
	}

	_, err := os.Stat(newPath)
	return err == nil
}

func (m *model) warnModalForRenaming() {
	id := shortuuid.New()
	message := channelMessage{
		messageId:   id,
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

	ti := textinput.New()
	ti.Cursor.Style = filePanelCursorStyle
	ti.Cursor.TextStyle = filePanelStyle
	ti.Prompt = filePanelCursorStyle.Render(icon.Cursor + " ")
	ti.TextStyle = modalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "New name"
	ti.PlaceholderStyle = modalStyle
	ti.SetValue(panel.element[panel.cursor].name)
	ti.SetCursor(cursorPos)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = m.fileModel.width - 4

	m.fileModel.renaming = true
	panel.renaming = true
	m.firstTextInput = true
	panel.rename = ti
}

func (m *model) deleteItemWarn() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if !((panel.panelMode == selectMode && len(panel.selected) != 0) || (panel.panelMode == browserMode)) {
		return
	}
	id := shortuuid.New()
	message := channelMessage{
		messageId:   id,
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
	} else {
		message.warnModal = warnModal{
			open:     true,
			title:    "Are you sure you want to move this to trash can",
			content:  "This operation will move file or directory to trash can.",
			warnType: confirmDeleteItem,
		}
		channel <- message
		return
	}
}

// Move file or directory to the trash can
func (m *model) deleteSingleItem() {
	id := shortuuid.New()
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}

	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle

	newProcess := process{
		name:     icon.Delete + icon.Space + panel.element[panel.cursor].name,
		progress: prog,
		state:    inOperation,
		total:    1,
		done:     0,
	}
	m.processBarModel.process[id] = newProcess

	message := channelMessage{
		messageId:       id,
		messageType:     sendProcess,
		processNewState: newProcess,
	}

	channel <- message
	err := trashMacOrLinux(panel.element[panel.cursor].location)

	if err != nil {
		p := m.processBarModel.process[id]
		p.state = failure
		message.processNewState = p
		channel <- message
	} else {
		p := m.processBarModel.process[id]
		p.done = 1
		p.state = successful
		p.doneTime = time.Now()
		message.processNewState = p
		channel <- message
	}
	if len(panel.element) == 0 {
		panel.cursor = 0
	} else {
		if panel.cursor >= len(panel.element) {
			panel.cursor = len(panel.element) - 1
		}
	}
}

// Move file or directory to the trash can
func (m *model) deleteMultipleItems() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.selected) != 0 {
		id := shortuuid.New()
		prog := progress.New(generateGradientColor())
		prog.PercentageStyle = footerStyle

		newProcess := process{
			name:     icon.Delete + icon.Space + filepath.Base(panel.selected[0]),
			progress: prog,
			state:    inOperation,
			total:    len(panel.selected),
			done:     0,
		}

		m.processBarModel.process[id] = newProcess

		message := channelMessage{
			messageId:       id,
			messageType:     sendProcess,
			processNewState: newProcess,
		}

		channel <- message

		for _, filePath := range panel.selected {

			p := m.processBarModel.process[id]
			p.name = icon.Delete + icon.Space + filepath.Base(filePath)
			p.done++
			p.state = inOperation
			if len(channel) < 5 {
				message.processNewState = p
				channel <- message
			}
			err := trashMacOrLinux(filePath)

			if err != nil {
				p.state = failure
				message.processNewState = p
				channel <- message
				outPutLog("Delete multiple item function error", err)
				m.processBarModel.process[id] = p
				break
			} else {
				if p.done == p.total {
					p.state = successful
					message.processNewState = p
					channel <- message
				}
				m.processBarModel.process[id] = p
			}
		}
	}

	if panel.cursor >= len(panel.element)-len(panel.selected)-1 {
		panel.cursor = len(panel.element) - len(panel.selected) - 1
		if panel.cursor < 0 {
			panel.cursor = 0
		}
	}
	panel.selected = panel.selected[:0]
}

// Completely delete file or folder (Not move to the trash can)
func (m *model) completelyDeleteSingleItem() {
	id := shortuuid.New()
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}

	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle

	newProcess := process{
		name:     "󰆴 " + panel.element[panel.cursor].name,
		progress: prog,
		state:    inOperation,
		total:    1,
		done:     0,
	}
	m.processBarModel.process[id] = newProcess

	message := channelMessage{
		messageId:       id,
		messageType:     sendProcess,
		processNewState: newProcess,
	}

	channel <- message

	err := os.RemoveAll(panel.element[panel.cursor].location)
	if err != nil {
		outPutLog("Completely delete single item function remove file error", err)
	}

	if err != nil {
		p := m.processBarModel.process[id]
		p.state = failure
		message.processNewState = p
		channel <- message
	} else {
		p := m.processBarModel.process[id]
		p.done = 1
		p.state = successful
		p.doneTime = time.Now()
		message.processNewState = p
		channel <- message
	}
	if len(panel.element) == 0 {
		panel.cursor = 0
	} else {
		if panel.cursor >= len(panel.element) {
			panel.cursor = len(panel.element) - 1
		}
	}
}

// Completely delete all file or folder from clipboard (Not move to the trash can)
func (m *model) completelyDeleteMultipleItems() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.selected) != 0 {
		id := shortuuid.New()
		prog := progress.New(generateGradientColor())
		prog.PercentageStyle = footerStyle

		newProcess := process{
			name:     icon.Delete + icon.Space + filepath.Base(panel.selected[0]),
			progress: prog,
			state:    inOperation,
			total:    len(panel.selected),
			done:     0,
		}

		m.processBarModel.process[id] = newProcess

		message := channelMessage{
			messageId:       id,
			messageType:     sendProcess,
			processNewState: newProcess,
		}

		channel <- message
		for _, filePath := range panel.selected {

			p := m.processBarModel.process[id]
			p.name = icon.Delete + icon.Space + filepath.Base(filePath)
			p.done++
			p.state = inOperation
			if len(channel) < 5 {
				message.processNewState = p
				channel <- message
			}
			err := os.RemoveAll(filePath)
			if err != nil {
				outPutLog("Completely delete multiple item function remove file error", err)
			}

			if err != nil {
				p.state = failure
				message.processNewState = p
				channel <- message
				outPutLog("Completely delete multiple item function error", err)
				m.processBarModel.process[id] = p
				break
			} else {
				if p.done == p.total {
					p.state = successful
					message.processNewState = p
					channel <- message
				}
				m.processBarModel.process[id] = p
			}
		}
	}

	if panel.cursor >= len(panel.element)-len(panel.selected)-1 {
		panel.cursor = len(panel.element) - len(panel.selected) - 1
		if panel.cursor < 0 {
			panel.cursor = 0
		}
	}
	panel.selected = panel.selected[:0]
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

// Paste all clipboard items
func (m *model) pasteItem() {
	if len(m.copyItems.items) == 0 {
		return
	}

	id := shortuuid.New()
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	totalFiles := 0

	for _, folderPath := range m.copyItems.items {
		// Todo : Fix this. This is inefficient
		// In case of a cut operations for a directory with a lot of files
		// we are unnecessarily walking the whole directory recursively
		// while os will just perform a rename
		// So instead of few operations this will cause the cut paste
		// to read the whole directory recursively
		// we should avoid doing this.
		// Although this allows us a more detailed progress tracking
		// this make the copy/cut more inefficient
		// instead, we could just track progress based on total items in
		// m.copyItems.items
		// efficiency should be prioritized over more detailed feedback.
		count, err := countFiles(folderPath)
		if err != nil {
			slog.Error("mode.pasteItem - Error in countFiles", "error", err)
			continue
		}
		totalFiles += count
	}

	slog.Debug("model.pasteItem", "items", m.copyItems.items, "cut", m.copyItems.cut,
		"totalFiles", totalFiles, "panel location", panel.location)

	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle

	prefixIcon := icon.Copy + icon.Space
	if m.copyItems.cut {
		prefixIcon = icon.Cut + icon.Space
	}

	newProcess := process{
		name:     prefixIcon + filepath.Base(m.copyItems.items[0]),
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}

	m.processBarModel.process[id] = newProcess

	message := channelMessage{
		messageId:       id,
		messageType:     sendProcess,
		processNewState: newProcess,
	}

	channel <- message

	p := m.processBarModel.process[id]
	for _, filePath := range m.copyItems.items {
		var err error
		if m.copyItems.cut && !isExternalDiskPath(filePath) {
			p.name = icon.Cut + icon.Space + filepath.Base(filePath)
		} else {
			if m.copyItems.cut {
				p.name = icon.Cut + icon.Space + filepath.Base(filePath)
			}
			p.name = icon.Copy + icon.Space + filepath.Base(filePath)
		}

		errMessage := "cut item error"
		if m.copyItems.cut && !isExternalDiskPath(filePath) {
			err = moveElement(filePath, filepath.Join(panel.location, filepath.Base(filePath)))
		} else {
			err = pasteDir(filePath, filepath.Join(panel.location, filepath.Base(filePath)), id, m)
			if err != nil {
				errMessage = "paste item error"
			} else {
				// Todo : These error cases are hard to test. We have to somehow make the paste operations fail,
				// which is time consuming and manual. We should test these with automated testcases
				if m.copyItems.cut {
					os.RemoveAll(filePath)
				}
			}
		}
		p = m.processBarModel.process[id]
		if err != nil {
			slog.Debug("model.pasteItem - paste failure", "error", err,
				"current item", filePath, "errMessage", errMessage)
			p.state = failure
			message.processNewState = p
			channel <- message
			outPutLog(errMessage, err)
			m.processBarModel.process[id] = p
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

	m.processBarModel.process[id] = p
	// Reset after paste is done. Only in case of cut
	// because current items in clipboard are deleted now
	if m.copyItems.cut {
		m.copyItems.reset(false)
	}
}

// Extrach compress file
func (m *model) extractFile() {
	var err error
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	ext := strings.ToLower(filepath.Ext(panel.element[panel.cursor].location))
	outputDir := fileNameWithoutExtension(panel.element[panel.cursor].location)
	outputDir, err = renameIfDuplicate(outputDir)

	if err != nil {
		outPutLog("Error extract file when create new directory", err)
	}

	switch ext {
	case ".zip":
		os.MkdirAll(outputDir, 0755)
		err = unzip(panel.element[panel.cursor].location, outputDir)
		if err != nil {
			outPutLog("Error extract file,", err)
		}
	default:
		os.MkdirAll(outputDir, 0755)
		err = extractCompressFile(panel.element[panel.cursor].location, outputDir)
		if err != nil {
			outPutLog("Error extract file,", err)
		}
	}
}

// Compress file or directory
func (m *model) compressFile() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	fileName := filepath.Base(panel.element[panel.cursor].location)

	zipName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".zip"
	zipName, err := renameIfDuplicate(zipName)

	if err != nil {
		outPutLog("Error compress file when rename duplicate", err)
	}

	zipSource(panel.element[panel.cursor].location, filepath.Join(filepath.Dir(panel.element[panel.cursor].location), zipName))
}

// Open file with default editor
func (m *model) openFileWithEditor() tea.Cmd {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	editor := Config.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	// Make sure there is an editor
	// Todo : Move hardcoded strings to constants : "windows", and editors
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "nano"
		}
	}
	c := exec.Command(editor, panel.element[panel.cursor].location)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Open directory with default editor
func (m *model) openDirectoryWithEditor() tea.Cmd {
	editor := Config.Editor
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	// Make sure there is an editor
	// Todo : Move hardcoded strings to constants : "windows", and editors
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "explorer"
		} else {
			editor = "nano"
		}
	}
	c := exec.Command(editor, m.fileModel.filePanels[m.filePanelFocusIndex].location)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// Copy file path
func (m *model) copyPath() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if err := clipboard.WriteAll(panel.element[panel.cursor].location); err != nil {
		outPutLog("Copy path error", panel.element[panel.cursor].location, err)
	}
}

func (m *model) copyPWD() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if err := clipboard.WriteAll(panel.location); err != nil {
		outPutLog("Copy present working directory error", panel.location, err)
	}
}
