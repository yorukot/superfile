package components

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/shortuuid"
	"github.com/rkoesters/xdg/trash"
)

/* CURSOR CONTROLLER START */
func ControllerSideBarListUp(m model) model {
	if m.sideBarModel.cursor > 0 {
		m.sideBarModel.cursor--
	} else {
		m.sideBarModel.cursor = len(m.sideBarModel.pinnedModel.folder) - 1
	}
	return m
}

func ControllerSideBarListDown(m model) model {
	if m.sideBarModel.cursor < len(m.sideBarModel.pinnedModel.folder)-1 {
		m.sideBarModel.cursor++
	} else {
		m.sideBarModel.cursor = 0
	}
	return m
}

func ControllerMetaDataListUp(m model) model {
	if m.fileMetaData.renderIndex > 0 {
		m.fileMetaData.renderIndex--
	} else {
		m.fileMetaData.renderIndex = len(m.fileMetaData.metaData) - 1
	}
	return m
}

func ControllerMetaDataListDown(m model) model {
	if m.fileMetaData.renderIndex < len(m.fileMetaData.metaData)-1 {
		m.fileMetaData.renderIndex++
	} else {
		m.fileMetaData.renderIndex = 0
	}
	return m
}

func ControllerFilePanelListUp(m model) model {
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
		if len(panel.element) > PanelElementHeight(m.mainPanelHeight) {
			panel.render = len(panel.element) - PanelElementHeight(m.mainPanelHeight)
			panel.cursor = len(panel.element) - 1
		} else {
			panel.cursor = len(panel.element) - 1
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func ControllerFilePanelListDown(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+PanelElementHeight(m.mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func ContollerProcessBarListUp(m model) model {
	if len(m.processBarModel.processList) == 0 {
		return m
	}
	if m.processBarModel.cursor > 0 {
		m.processBarModel.cursor--
		if m.processBarModel.cursor < m.processBarModel.render {
			m.processBarModel.render--
		}
	} else {
		if len(m.processBarModel.processList) <= 3 {
			m.processBarModel.cursor = len(m.processBarModel.processList) - 1
		} else {
			m.processBarModel.render = len(m.processBarModel.processList) - 3
			m.processBarModel.cursor = len(m.processBarModel.processList) - 1
		}
	}

	return m
}

func ContollerProcessBarListDown(m model) model {
	if len(m.processBarModel.processList) == 0 {
		return m
	}
	if m.processBarModel.cursor < len(m.processBarModel.processList)-1 {
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

func SideBarSelectFolder(m model) model {
	m.focusPanel = nonePanelFocus
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	panel.folderRecord[panel.location] = folderRecord{
		folderCursor: panel.cursor,
		folderRender: panel.render,
	}

	panel.location = m.sideBarModel.pinnedModel.folder[m.sideBarModel.cursor].location
	folderRecord, hasRecord := panel.folderRecord[panel.location]
	if hasRecord {
		panel.cursor = folderRecord.folderCursor
		panel.render = folderRecord.folderRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}
	panel.focusType = focus
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func SelectedMode(m model) model {
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

func NextFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	return m
}

func PreviousFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	return m
}

func CloseFilePanel(m model) model {
	if len(m.fileModel.filePanels) != 1 {
		m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

		if m.filePanelFocusIndex != 0 {
			m.filePanelFocusIndex--
		}
		m.fileModel.width = (m.fullWidth - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	}
	return m
}
func CreateNewFilePanel(m model) model {
	if len(m.fileModel.filePanels) != 4 {
		m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{
			location:     HomeDir,
			panelMode:    browserMode,
			focusType:    secondFocus,
			folderRecord: make(map[string]folderRecord),
		})

		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
		m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.focusPanel)
		m.fileModel.width = (m.fullWidth - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.filePanelFocusIndex++
	}
	return m
}

func FocusOnSideBar(m model) model {
	if m.focusPanel == sideBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = sideBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

func FocusOnProcessBar(m model) model {
	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

func FocusOnMetaData(m model) model {
	if m.focusPanel == metaDataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = metaDataFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}

func PasteItem(m model) model {
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
	prog := progress.New(progress.WithScaledGradient(theme.ProcessBarGradient[0], theme.ProcessBarGradient[1]))
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

	processBarChannel <- processBarMessage{
		processId:       id,
		processNewState: newProcess,
	}

	for _, filePath := range m.copyItems.items {
		p := m.processBarModel.process[id]
		if m.copyItems.cut {
			p.name = "󰆐 " + filepath.Base(filePath)
		} else {
			p.name = "󰆏 " + filepath.Base(filePath)
		}

		newModel, err := PasteDir(filePath, panel.location+"/"+path.Base(filePath), id, m)
		m = newModel
		p = m.processBarModel.process[id]
		if err != nil {
			p.state = failure
			processBarChannel <- processBarMessage{
				processId:       id,
				processNewState: p,
			}
			OutPutLog("Pasted item error", err)
			m.processBarModel.process[id] = p
			break
		} else {
			if p.done == p.total {
				p.state = successful
				p.done = totalFiles
				p.doneTime = time.Now()
				processBarChannel <- processBarMessage{
					processId:       id,
					processNewState: p,
				}
			}
			m.processBarModel.process[id] = p
		}
	}
	if m.copyItems.cut {
		for _, item := range m.copyItems.items {
			err := trash.Trash(item)
			if err != nil {
				OutPutLog("Paste item function move file to trash can error", err)
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

func PanelCreateNewFile(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	ti := textinput.New()
	ti.Placeholder = "File name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = modalWidth - 10

	m.createNewItem.location = panel.location
	m.createNewItem.itemType = newFile
	m.createNewItem.open = true
	m.createNewItem.textInput = ti
	m.firstTextInput = true

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func PanelCreateNewFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	ti := textinput.New()
	ti.Placeholder = "Folder name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = modalWidth - 10

	m.createNewItem.location = panel.location
	m.createNewItem.itemType = newFolder
	m.createNewItem.open = true
	m.createNewItem.textInput = ti
	m.firstTextInput = true

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel

	return m
}

func PinnedFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	unPinned := false

	jsonData, err := os.ReadFile(SuperFileMainDir + pinnedFile)
	if err != nil {
		OutPutLog("Pinned folder function read superfile data error", err)
	}
	var pinnedFolder []string
	err = json.Unmarshal(jsonData, &pinnedFolder)
	if err != nil {
		OutPutLog("Pinned folder function unmarshal superfile data error", err)
	}
	for i, other := range pinnedFolder {
		if other == panel.location {
			pinnedFolder = append(pinnedFolder[:i], pinnedFolder[i+1:]...)
			unPinned = true
		}
	}

	if !contains(pinnedFolder, panel.location) && !unPinned {
		pinnedFolder = append(pinnedFolder, panel.location)
	}

	updatedData, err := json.Marshal(pinnedFolder)
	if err != nil {
		OutPutLog("Pinned folder function updatedData superfile data error", err)
	}

	err = os.WriteFile(SuperFileMainDir+pinnedFile, updatedData, 0644)
	if err != nil {
		OutPutLog("Pinned folder function updatedData superfile data error", err)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}


// I can't test all of the os system so if you have any problem or want to add support please create new pull request!
func OpenTerminal(m model) model {

	currentDir := m.fileModel.filePanels[m.filePanelFocusIndex].location
	terminal := Config.Terminal
	workDirSet := Config.TerminalWorkDir
	if terminal != "" {
		cmd := exec.Command(terminal, workDirSet+currentDir)
		err := cmd.Start()
		if err != nil {
			OutPutLog("Error opening"+terminal+":", err)
		}
		return m
	}

	if runtime.GOOS == "darwin" {
		terminal = "Terminal.app"
        workDirSet = "--working-directory="
		cmd := exec.Command(terminal, workDirSet+currentDir)
		err := cmd.Start()
		if err != nil {
			OutPutLog("Error opening"+terminal+":", err)
		}
		return m
    }
	
	desktopEnv := os.Getenv("XDG_CURRENT_DESKTOP")
	switch desktopEnv {
	case "GNOME":
		terminal = "gnome-terminal"
		workDirSet = "--working-directory="
	case "KDE":
		terminal = "konsole"
		workDirSet = "--workdir="
	case "XFCE":
		terminal = "xfce4-terminal"
		workDirSet = "--working-directory="
	case "LXDE":
		terminal = "lxterminal"
		workDirSet = "--working-directory="
	case "CINNAMON":
		terminal = "gnome-terminal"
		workDirSet = "--working-directory="
	case "MATE":
		terminal = "mate-terminal"
		workDirSet = "--working-directory="
	case "LXQT":
		terminal = "qterminal"
		workDirSet = "--working-directory="
	case "BUDGIE":
		terminal = "gnome-terminal"
		workDirSet = "--working-directory="
	case "PANTHEON":
		terminal = "pantheon-terminal"
		workDirSet = "--working-directory="
	case "DEEPIN":
		terminal = "deepin-terminal"
		workDirSet = "--working-directory="
	case "ENLIGHTENMENT":
		terminal = "terminology"
		workDirSet = "--working-directory="
	case "UNITY":
		terminal = "gnome-terminal"
		workDirSet = "--working-directory="
	default:
		log.Fatalf("We can't find your default terminal please go to ~/.config/superfile/config/config.json setting your default terminal and terminalWorkDirFlag!")
	}
	
	cmd := exec.Command(terminal, workDirSet+currentDir)
	err := cmd.Start()
	if err != nil {
		OutPutLog("Error opening"+terminal+":", err)
	}

	return m
}
