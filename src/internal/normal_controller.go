package internal

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/shortuuid"
	"github.com/rkoesters/xdg/trash"
)

func enterPanel(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) > 0 && panel.element[panel.cursor].directory {
		panel.directoryRecord[panel.location] = directoryRecord{
			directoryCursor: panel.cursor,
			directoryRender: panel.render,
		}
		panel.location = panel.element[panel.cursor].location
		directoryRecord, hasRecord := panel.directoryRecord[panel.location]
		if hasRecord {
			panel.cursor = directoryRecord.directoryCursor
			panel.render = directoryRecord.directoryRender
		} else {
			panel.cursor = 0
			panel.render = 0
		}
	} else if len(panel.element) > 0 && !panel.element[panel.cursor].directory {
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("open", panel.element[panel.cursor].location)
			_, err := cmd.Output()
			if err != nil {
				outPutLog("err when open file with open:", err)
			}
		} else {
			cmd := exec.Command("xdg-open", panel.element[panel.cursor].location)
			_, err := cmd.Output()
			if err != nil {
				outPutLog("err when open file with xdg-open:", err)
			}
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func parentFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.directoryRecord[panel.location] = directoryRecord{
		directoryCursor: panel.cursor,
		directoryRender: panel.render,
	}
	fullPath := panel.location
	parentDir := path.Dir(fullPath)
	panel.location = parentDir
	directoryRecord, hasRecord := panel.directoryRecord[panel.location]
	if hasRecord {
		panel.cursor = directoryRecord.directoryCursor
		panel.render = directoryRecord.directoryRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func completelyDeleteSingleFile(m model) model {
	id := shortuuid.New()
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return m
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

	channel <- channelMessage{
		messageId:       id,
		processNewState: newProcess,
	}

	err := os.RemoveAll(panel.element[panel.cursor].location)
	if err != nil {
		outPutLog("Completely delete single item function remove file error", err)
	}

	if err != nil {
		p := m.processBarModel.process[id]
		p.state = failure
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
	} else {
		p := m.processBarModel.process[id]
		p.done = 1
		p.state = successful
		p.doneTime = time.Now()
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
	}
	if panel.cursor == len(panel.element)-1 {
		panel.cursor--
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func deleteSingleItem(m model) model {
	id := shortuuid.New()
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return m
	}

	if isExternalDiskPath(panel.location) {
		channel <- channelMessage{
			messageId:       id,
			returnWarnModal: true,
			warnModal: warnModal{
				open:     true,
				title:    "Are you sure you want to completely delete",
				content:  "This operation cannot be undone and your data will be completely lost.",
				warnType: confirmDeleteItem,
			},
		}
		return m
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

	channel <- channelMessage{
		messageId:       id,
		processNewState: newProcess,
	}
	var err error
	if runtime.GOOS == "darwin" {
		err := moveElement(panel.element[panel.cursor].location, HomeDir+"/.Trash/"+filepath.Base(panel.element[panel.cursor].location))
		if err != nil {
			outPutLog("Delete single item function move file to trash can error", err)
		}
	} else {
		err = trash.Trash(panel.element[panel.cursor].location)
		if err != nil {
			outPutLog("Delete single item function move file to trash can error", err)
		}
	}

	if err != nil {
		p := m.processBarModel.process[id]
		p.state = failure
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
	} else {
		p := m.processBarModel.process[id]
		p.done = 1
		p.state = successful
		p.doneTime = time.Now()
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
	}
	if panel.cursor == len(panel.element)-1 {
		panel.cursor--
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func copySingleItem(m model) model {
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
		outPutLog("Copy single item get file state error", panel.element[panel.cursor].location, err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.element[panel.cursor].location)

		if err != nil {
			outPutLog("Copy single item read file error", panel.element[panel.cursor].location, err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			outPutLog("Copy single item write file error", panel.element[panel.cursor].location, err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func cutSingleItem(m model) model {
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
		outPutLog("Cut single item get file state error", panel.element[panel.cursor].location, err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.element[panel.cursor].location)

		if err != nil {
			outPutLog("Cut single item read file error", panel.element[panel.cursor].location, err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			outPutLog("Cut single item write file error", panel.element[panel.cursor].location, err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func panelItemRename(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	ti := textinput.New()
	ti.Cursor.Style = filePanelCursorStyle
	ti.Cursor.TextStyle = filePanelStyle
	ti.TextStyle = modalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "New name"
	ti.PlaceholderStyle = modalStyle
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

func searchBarFocus(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.searchBar.Focused() {
		panel.searchBar.Blur()
	} else {
		panel.searchBar.Focus()
	}

	// config search bar width
	panel.searchBar.Width = m.fileModel.width - 4
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}
