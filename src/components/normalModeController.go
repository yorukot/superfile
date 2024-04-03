package components

import (
	"os"
	"path"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
)

func EnterPanel(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) > 0 && panel.element[panel.cursor].folder {
		panel.folderRecord[panel.location] = folderRecord{
			folderCursor: panel.cursor,
			folderRender: panel.render,
		}
		panel.location = panel.element[panel.cursor].location
		folderRecord, hasRecord := panel.folderRecord[panel.location]
		if hasRecord {
			panel.cursor = folderRecord.folderCursor
			panel.render = folderRecord.folderRender
		} else {
			panel.cursor = 0
			panel.render = 0
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func ParentFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.folderRecord[panel.location] = folderRecord{
		folderCursor: panel.cursor,
		folderRender: panel.render,
	}
	fullPath := panel.location
	parentDir := path.Dir(fullPath)
	panel.location = parentDir
	folderRecord, hasRecord := panel.folderRecord[panel.location]
	if hasRecord {
		panel.cursor = folderRecord.folderCursor
		panel.render = folderRecord.folderRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func DeleteSingleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	prog := progress.New(progress.WithScaledGradient(theme.ProcessBarGradient[0], theme.ProcessBarGradient[1]))
	m.processBarModel.process = append(m.processBarModel.process, process{
		name:     "󰆴 " + panel.element[panel.cursor].name,
		progress: prog,
		state:    deleting,
	})
	err := MoveFile(panel.element[panel.cursor].location, Config.TrashCanPath+"/"+panel.element[panel.cursor].name)
	if err != nil {
		m.processBarModel.process[0].state = failure
	} else {
		m.processBarModel.process[0].state = successful
		m.processBarModel.process[0].progress.IncrPercent(1)
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CopySingleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	m.copyItems.items = append(m.copyItems.items, panel.element[panel.cursor].location)
	fileInfo, err := os.Stat(panel.element[panel.cursor].location)
	if err != nil {
		OutputLog("Can't find this file or folder")
		OutputLog(panel.element[panel.cursor].location)
		OutputLog(err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.element[panel.cursor].location)

		CheckErr(err)

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			OutputLog(err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CutSingleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	m.copyItems.items = append(m.copyItems.items, panel.element[panel.cursor].location)
	m.copyItems.cut = true
	m.copyItems.oringnalPanel = orignalPanel{
		index:    m.filePanelFocusIndex,
		location: panel.location,
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func PanelItemRename(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	ti := textinput.New()
	ti.Placeholder = "New name"
	ti.SetValue(panel.element[panel.cursor].name)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = m.fileModel.width - 4

	m.fileModel.renaming = true
	panel.renaming = true
	m.firstTextInput = true
	panel.rename = ti
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}
