package components

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/shortuuid"
	"github.com/rkoesters/xdg/trash"
)

/* CURSOR CONTROLLER START */
func controllerSideBarListUp(m model) model {
	if m.sidebarModel.cursor > 0 {
		m.sidebarModel.cursor--
	} else {
		m.sidebarModel.cursor = len(m.sidebarModel.directories) - 1
	}
	return m
}

func controllerSideBarListDown(m model) model {
	lenDirs := len(m.sidebarModel.directories)
	if m.sidebarModel.cursor < lenDirs-1 {
		m.sidebarModel.cursor++
	} else {
		m.sidebarModel.cursor = 0
	}
	return m
}

func controllerMetaDataListUp(m model) model {
	if m.fileMetaData.renderIndex > 0 {
		m.fileMetaData.renderIndex--
	} else {
		m.fileMetaData.renderIndex = len(m.fileMetaData.metaData) - 1
	}
	return m
}

func controllerMetaDataListDown(m model) model {
	if m.fileMetaData.renderIndex < len(m.fileMetaData.metaData)-1 {
		m.fileMetaData.renderIndex++
	} else {
		m.fileMetaData.renderIndex = 0
	}
	return m
}

func controllerFilePanelListUp(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	if panel.cursor > 0 {
		panel.cursor--
		if panel.cursor < panel.render {
			panel.render--
		}
	} else {
		if len(panel.element) > panelElementHeight(m.mainPanelHeight) {
			panel.render = len(panel.element) - panelElementHeight(m.mainPanelHeight)
			panel.cursor = len(panel.element) - 1
		} else {
			panel.cursor = len(panel.element) - 1
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func controllerFilePanelListDown(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+panelElementHeight(m.mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func contollerProcessBarListUp(m model) model {
	if len(m.processBarModel.processList) == 0 {
		return m
	}
	if m.processBarModel.cursor > 0 {
		m.processBarModel.cursor--
		if m.processBarModel.cursor < m.processBarModel.render {
			m.processBarModel.render--
		}
	} else {
		if len(m.processBarModel.processList) <= 3 || (len(m.processBarModel.processList) <= 2 && footerHeight < 14) {
			m.processBarModel.cursor = len(m.processBarModel.processList) - 1
		} else {
			m.processBarModel.render = len(m.processBarModel.processList) - 3
			m.processBarModel.cursor = len(m.processBarModel.processList) - 1
		}
	}

	return m
}

func contollerProcessBarListDown(m model) model {
	if len(m.processBarModel.processList) == 0 {
		return m
	}
	if m.processBarModel.cursor < len(m.processBarModel.processList)-1  {
		m.processBarModel.cursor++
		if m.processBarModel.cursor > m.processBarModel.render+2 {
			m.processBarModel.render++
		}
	} else { 
		m.processBarModel.render = 0
		m.processBarModel.cursor = 0
	}
	return m
}

/* CURSOR CONTROLLER END */

/* LIST CONTROLLER START */

func sidebarSelectFolder(m model) model {
	m.focusPanel = nonePanelFocus
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	panel.directoryRecord[panel.location] = directoryRecord{
		directoryCursor: panel.cursor,
		directoryRender: panel.render,
	}

	panel.location = m.sidebarModel.directories[m.sidebarModel.cursor].location
	directoryRecord, hasRecord := panel.directoryRecord[panel.location]
	if hasRecord {
		panel.cursor = directoryRecord.directoryCursor
		panel.render = directoryRecord.directoryRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}
	panel.focusType = focus
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func selectedMode(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.panelMode == selectMode {
		panel.selected = panel.selected[:0]
		panel.panelMode = browserMode
	} else if panel.panelMode == browserMode {
		panel.panelMode = selectMode
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

/* LIST CONTROLLER END */

func nextFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	return m
}

func previousFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	return m
}

func closeFilePanel(m model) model {
	if len(m.fileModel.filePanels) != 1 {
		m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

		if m.filePanelFocusIndex != 0 {
			m.filePanelFocusIndex--
		}
		m.fileModel.width = (m.fullWidth - sidebarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	}
	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
	return m
}
func createNewFilePanel(m model) model {
	if len(m.fileModel.filePanels) != m.fileModel.maxFilePanel {
		m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{
			location:        HomeDir,
			panelMode:       browserMode,
			focusType:       secondFocus,
			directoryRecord: make(map[string]directoryRecord),
			searchBar:       generateSearchBar(),
		})

		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
		m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.focusPanel)
		m.fileModel.width = (m.fullWidth - sidebarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.filePanelFocusIndex++
	}
	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
	return m
}

func focusOnSideBar(m model) model {
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = sidebarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

func focusOnProcessBar(m model) model {
	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

func focusOnMetaData(m model) model {
	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = metadataFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

func pasteItem(m model) model {
	id := shortuuid.New()
	if len(m.copyItems.items) == 0 {
		return m
	}
	totalFiles := 0
	for _, folderPath := range m.copyItems.items {
		count, err := countFiles(folderPath)
		if err != nil {
			continue
		}
		totalFiles += count
	}

	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle

	newProcess := process{}

	if m.copyItems.cut {
		newProcess = process{
			name:     "󰆐 " + filepath.Base(m.copyItems.items[0]),
			progress: prog,
			state:    inOperation,
			total:    totalFiles,
			done:     0,
		}
	} else {
		newProcess = process{
			name:     "󰆏 " + filepath.Base(m.copyItems.items[0]),
			progress: prog,
			state:    inOperation,
			total:    totalFiles,
			done:     0,
		}
	}

	m.processBarModel.process[id] = newProcess

	channel <- channelMessage{
		messageId:       id,
		processNewState: newProcess,
	}

	for _, filePath := range m.copyItems.items {
		p := m.processBarModel.process[id]
		if m.copyItems.cut {
			p.name = "󰆐 " + filepath.Base(filePath)
		} else {
			p.name = "󰆏 " + filepath.Base(filePath)
		}

		newModel, err := pasteDir(filePath, panel.location+"/"+path.Base(filePath), id, m)
		m = newModel
		p = m.processBarModel.process[id]
		if err != nil {
			p.state = failure
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
			outPutLog("Pasted item error", err)
			m.processBarModel.process[id] = p
			break
		} else {
			if p.done == p.total {
				p.state = successful
				p.done = totalFiles
				p.doneTime = time.Now()
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
			}
			m.processBarModel.process[id] = p
		}
	}
	if m.copyItems.cut {
		for _, item := range m.copyItems.items {
			if runtime.GOOS == "darwin" {
				err := moveElement(item,  HomeDir + "/.Trash/" + filepath.Base(item))
				if err != nil {
					outPutLog("Delete single item function move file to trash can error", err)
				}
			} else {
				err := trash.Trash(item)
				if err != nil {
					outPutLog("Paste item function move file to trash can error", err)
				}
			}

		}
		if m.fileModel.filePanels[m.copyItems.originalPanel.index].location == m.copyItems.originalPanel.location {
			m.fileModel.filePanels[m.copyItems.originalPanel.index].selected = panel.selected[:0]
		}
	}
	m.copyItems.cut = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func panelCreateNewFile(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	ti := textinput.New()
	ti.Cursor.Style = modalCursorStyle
	ti.Cursor.TextStyle = modalStyle
	ti.TextStyle = modalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "File name"
	ti.PlaceholderStyle = modalStyle
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = modalWidth - 10

	m.typingModal.location = panel.location
	m.typingModal.itemType = newFile
	m.typingModal.open = true
	m.typingModal.textInput = ti
	m.firstTextInput = true

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func panelCreateNewFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	ti := textinput.New()
	ti.Cursor.Style = modalCursorStyle
	ti.Cursor.TextStyle = modalStyle
	ti.TextStyle = modalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "Folder name"
	ti.PlaceholderStyle = modalStyle
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = modalWidth - 10

	m.typingModal.location = panel.location
	m.typingModal.itemType = newDirectory
	m.typingModal.open = true
	m.typingModal.textInput = ti
	m.firstTextInput = true

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func pinnedFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	unPinned := false

	jsonData, err := os.ReadFile(SuperFileDataDir + pinnedFile)
	if err != nil {
		outPutLog("Pinned folder function read superfile data error", err)
	}
	var pinnedFolder []string
	err = json.Unmarshal(jsonData, &pinnedFolder)
	if err != nil {
		outPutLog("Pinned folder function unmarshal superfile data error", err)
	}
	for i, other := range pinnedFolder {
		if other == panel.location {
			pinnedFolder = append(pinnedFolder[:i], pinnedFolder[i+1:]...)
			unPinned = true
		}
	}

	if !arrayContains(pinnedFolder, panel.location) && !unPinned {
		pinnedFolder = append(pinnedFolder, panel.location)
	}

	updatedData, err := json.Marshal(pinnedFolder)
	if err != nil {
		outPutLog("Pinned folder function updatedData superfile data error", err)
	}

	err = os.WriteFile(SuperFileDataDir+pinnedFile, updatedData, 0644)
	if err != nil {
		outPutLog("Pinned folder function updatedData superfile data error", err)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func toggleDotFileController(m model) model {
	newToggleDotFile := ""
	if m.toggleDotFile {
		newToggleDotFile = "false"
		m.toggleDotFile = false
	} else {
		newToggleDotFile = "true"
		m.toggleDotFile = true
	}
	err := os.WriteFile(SuperFileDataDir+toggleDotFile, []byte(newToggleDotFile), 0644)
	if err != nil {
		outPutLog("Pinned folder function updatedData superfile data error", err)
	}

	return m
}

func extractFile(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	unzip(panel.element[panel.cursor].location, filepath.Dir(panel.element[panel.cursor].location))
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func compressFile(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	fileName := filepath.Base(panel.element[panel.cursor].location)

	zipName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".zip"
	zipSource(panel.element[panel.cursor].location, filepath.Join(filepath.Dir(panel.element[panel.cursor].location), zipName))
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func openFileWithEditor(m model) tea.Cmd {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	editor := os.Getenv("EDITOR")
	m.editorMode = true
	if editor == "" {
		editor = "nano"
	}
	c := exec.Command(editor, panel.element[panel.cursor].location)
	
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func openDirectoryWithEditor(m model) tea.Cmd {
	editor := os.Getenv("EDITOR")
	m.editorMode = true
	if editor == "" {
		editor = "nano"
	}
	c := exec.Command(editor, m.fileModel.filePanels[m.filePanelFocusIndex].location)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func openHelpMenu(m model) model {
	if m.helpMenu.open {
		m.helpMenu.open = false
		return m
	}

	m.helpMenu.open = true
	return m
}