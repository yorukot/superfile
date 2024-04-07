package components

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/rkoesters/xdg/trash"
	"os"
	"path/filepath"
)

func SingleItemSelect(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if ArrayContains(panel.selected, panel.element[panel.cursor].location) {
		panel.selected = RemoveElementByValue(panel.selected, panel.element[panel.cursor].location)
	} else {
		panel.selected = append(panel.selected, panel.element[panel.cursor].location)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func ItemSelectUp(m model) model {
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
	if ArrayContains(panel.selected, panel.element[panel.cursor].location) {
		panel.selected = RemoveElementByValue(panel.selected, panel.element[panel.cursor].location)
	} else {
		panel.selected = append(panel.selected, panel.element[panel.cursor].location)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func ItemSelectDown(m model) model {
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
	if ArrayContains(panel.selected, panel.element[panel.cursor].location) {
		panel.selected = RemoveElementByValue(panel.selected, panel.element[panel.cursor].location)
	} else {
		panel.selected = append(panel.selected, panel.element[panel.cursor].location)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func DeleteMultipleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.selected) != 0 {
		id := shortuuid.New()
		prog := progress.New(progress.WithScaledGradient(theme.ProcessBarGradient[0], theme.ProcessBarGradient[1]))
		newProcess := process{
			name:     "󰆴 " + filepath.Base(panel.selected[0]),
			progress: prog,
			state:    inOperation,
			total:    len(panel.selected),
			done:     0,
		}

		m.processBarModel.process[id] = newProcess

		processBarChannel <- processBarMessage{
			processId:       id,
			processNewState: newProcess,
		}

		for _, filePath := range panel.selected {

			p := m.processBarModel.process[id]
			p.name = "󰆴 " + filepath.Base(filePath)
			p.done++
			p.state = inOperation
			if len(processBarChannel) < 5 {
				processBarChannel <- processBarMessage{
					processId:       id,
					processNewState: p,
				}
			}
			err := trash.Trash(filePath)
			if err != nil {
				OutPutLog("Delete single item function move file to trash can error", err)
			}

			if err != nil {
				p.state = failure
				processBarChannel <- processBarMessage{
					processId:       id,
					processNewState: p,
				}
				OutPutLog("Delete multiple item function error", err)
				m.processBarModel.process[id] = p
				break
			} else {
				if p.done == p.total {
					p.state = successful
					processBarChannel <- processBarMessage{
						processId:       id,
						processNewState: p,
					}
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
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CopyMultipleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.cut = false
	m.copyItems.items = m.copyItems.items[:0]
	if len(panel.selected) == 0 {
		return m
	}
	m.copyItems.items = panel.selected
	fileInfo, err := os.Stat(panel.selected[0])
	if os.IsNotExist(err) {
		return m
	}
	if err != nil {
		OutPutLog("Copy multiple item function get file state error", panel.selected[0], err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.selected[0])

		if err != nil {
			OutPutLog("Copy multiple item function read file error", err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			OutPutLog("Copy multiple item function write file to clipboard error", err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CutMultipleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.cut = true
	m.copyItems.items = m.copyItems.items[:0]
	if len(panel.selected) == 0 {
		return m
	}
	m.copyItems.items = panel.selected
	fileInfo, err := os.Stat(panel.selected[0])
	if os.IsNotExist(err) {
		return m
	}
	if err != nil {
		OutPutLog("Copy multiple item function get file state error", panel.selected[0], err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.selected[0])

		if err != nil {
			OutPutLog("Copy multiple item function read file error", err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			OutPutLog("Copy multiple item function write file to clipboard error", err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func SelectAllItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	for _, item := range panel.element {
		panel.selected = append(panel.selected, item.location)
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}
