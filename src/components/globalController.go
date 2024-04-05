package components

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/shortuuid"
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

/* CURSOR CONTROLLER END */

/* LIST CONTROLLER START */

func SideBarSelectFolder(m model) model {
	m.sideBarModel.pinnedModel.selected = m.sideBarModel.pinnedModel.folder[m.sideBarModel.cursor].location
	m.fileModel.filePanels[m.filePanelFocusIndex].location = m.sideBarModel.pinnedModel.selected
	m.focusPanel = nonePanelFocus
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	m.fileModel.filePanels[m.filePanelFocusIndex].cursor = 0
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
		OutputLog(filePath)
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
			OutputLog("Error delete multiple item")
			OutputLog(err)
			m.processBarModel.process[id] = p
			break
		} else {
			if p.done == p.total {
				p.state = successful
				p.done = totalFiles
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
			filePath := item
			err := MoveFile(item, Config.TrashCanPath+"/"+path.Base(filePath))
			CheckErr(err)
		}
		if m.fileModel.filePanels[m.copyItems.oringnalPanel.index].location == m.copyItems.oringnalPanel.location {
			m.fileModel.filePanels[m.copyItems.oringnalPanel.index].selected = panel.selected[:0]
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

	jsonData, err := os.ReadFile("./.superfile/data/superfile.json")
	CheckErr(err)

	var pinnedFolder []string
	err = json.Unmarshal(jsonData, &pinnedFolder)
	CheckErr(err)
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
	CheckErr(err)

	err = os.WriteFile("./.superfile/data/superfile.json", updatedData, 0644)
	CheckErr(err)

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}
