package components

import (
	"path"

	"github.com/charmbracelet/bubbles/progress"
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
	return m
}
