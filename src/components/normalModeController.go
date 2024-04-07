package components

import (
	"os"
	"path"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/shortuuid"
	"github.com/rkoesters/xdg/trash"
	"os/exec"
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
	} else if len(panel.element) > 0 && !panel.element[panel.cursor].folder {
		cmd := exec.Command("xdg-open", panel.element[panel.cursor].location)
		_, err := cmd.Output()
		if err != nil {
			OutPutLog("err when open file with xdg-open:", err)
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
	id := shortuuid.New()
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	prog := progress.New(progress.WithScaledGradient(theme.ProcessBarGradient[0], theme.ProcessBarGradient[1]))
	newProcess := process{
		name:     "󰆴 " + panel.element[panel.cursor].name,
		progress: prog,
		state:    inOperation,
		total:    1,
		done:     0,
	}
	m.processBarModel.process[id] = newProcess

	processBarChannel <- processBarMessage{
		processId:       id,
		processNewState: newProcess,
	}

	err := trash.Trash(panel.element[panel.cursor].location)
	if err != nil {
		OutPutLog("Delete single item function move file to trash can error", err)
	}

	if err != nil {
		p := m.processBarModel.process[id]
		p.state = failure
		processBarChannel <- processBarMessage{
			processId:       id,
			processNewState: p,
		}
	} else {
		p := m.processBarModel.process[id]
		p.done = 1
		p.state = successful
		p.doneTime = time.Now()
		processBarChannel <- processBarMessage{
			processId:       id,
			processNewState: p,
		}
	}
	if panel.cursor == len(panel.element)-1 {
		panel.cursor--
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CopySingleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.cut = false
	m.copyItems.items = m.copyItems.items[:0]
	if len(panel.element) == 0 {
		return m
	}
	m.copyItems.items = append(m.copyItems.items, panel.element[panel.cursor].location)
	fileInfo, err := os.Stat(panel.element[panel.cursor].location)
	if os.IsNotExist(err) {
		m.copyItems.items = m.copyItems.items[:0]
		return m
	}
	if err != nil {
		OutPutLog("Copy single item get file state error", panel.element[panel.cursor].location, err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.element[panel.cursor].location)

		if err != nil {
			OutPutLog("Copy single item read file error", panel.element[panel.cursor].location, err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			OutPutLog("Copy single item write file error", panel.element[panel.cursor].location, err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CutSingleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.cut = true
	m.copyItems.items = m.copyItems.items[:0]
	if len(panel.element) == 0 {
		return m
	}
	m.copyItems.items = append(m.copyItems.items, panel.element[panel.cursor].location)
	fileInfo, err := os.Stat(panel.element[panel.cursor].location)
	if os.IsNotExist(err) {
		m.copyItems.items = m.copyItems.items[:0]
		return m
	}
	if err != nil {
		OutPutLog("Cut single item get file state error", panel.element[panel.cursor].location, err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.element[panel.cursor].location)

		if err != nil {
			OutPutLog("Cut single item read file error", panel.element[panel.cursor].location, err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			OutPutLog("Cut single item write file error", panel.element[panel.cursor].location, err)
		}
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
