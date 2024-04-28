package internal

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/rkoesters/xdg/trash"
)

func singleItemSelect(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if arrayContains(panel.selected, panel.element[panel.cursor].location) {
		panel.selected = removeElementByValue(panel.selected, panel.element[panel.cursor].location)
	} else {
		panel.selected = append(panel.selected, panel.element[panel.cursor].location)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func itemSelectUp(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
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
	if arrayContains(panel.selected, panel.element[panel.cursor].location) {
		panel.selected = removeElementByValue(panel.selected, panel.element[panel.cursor].location)
	} else {
		panel.selected = append(panel.selected, panel.element[panel.cursor].location)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func itemSelectDown(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+panelElementHeight(m.mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
	if arrayContains(panel.selected, panel.element[panel.cursor].location) {
		panel.selected = removeElementByValue(panel.selected, panel.element[panel.cursor].location)
	} else {
		panel.selected = append(panel.selected, panel.element[panel.cursor].location)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func completelyDeleteMultipleFile(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.selected) != 0 {
		id := shortuuid.New()
		prog := progress.New(generateGradientColor())
		prog.PercentageStyle = footerStyle

		newProcess := process{
			name:     "󰆴 " + filepath.Base(panel.selected[0]),
			progress: prog,
			state:    inOperation,
			total:    len(panel.selected),
			done:     0,
		}

		m.processBarModel.process[id] = newProcess

		channel <- channelMessage{
			messageId:       id,
			processNewState: newProcess,
		}
		for _, filePath := range panel.selected {

			p := m.processBarModel.process[id]
			p.name = "󰆴 " + filepath.Base(filePath)
			p.done++
			p.state = inOperation
			if len(channel) < 5 {
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
			}
			err := os.RemoveAll(filePath)
			if err != nil {
				outPutLog("Completely delete multiple item function remove file error", err)
			}

			if err != nil {
				p.state = failure
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
				outPutLog("Completely delete multiple item function error", err)
				m.processBarModel.process[id] = p
				break
			} else {
				if p.done == p.total {
					p.state = successful
					channel <- channelMessage{
						messageId:       id,
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

func deleteMultipleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.selected) != 0 {
		id := shortuuid.New()
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
			name:     "󰆴 " + filepath.Base(panel.selected[0]),
			progress: prog,
			state:    inOperation,
			total:    len(panel.selected),
			done:     0,
		}

		m.processBarModel.process[id] = newProcess

		channel <- channelMessage{
			messageId:       id,
			processNewState: newProcess,
		}

		for _, filePath := range panel.selected {

			p := m.processBarModel.process[id]
			p.name = "󰆴 " + filepath.Base(filePath)
			p.done++
			p.state = inOperation
			if len(channel) < 5 {
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
			}
			var err error
			if runtime.GOOS == "darwin" {
				err := moveElement(panel.element[panel.cursor].location, HomeDir+"/.Trash/"+filepath.Base(panel.element[panel.cursor].location))
				if err != nil {
					outPutLog("Delete single item function move file to trash can error", err)
				}
			} else {
				err = trash.Trash(filePath)
				if err != nil {
					outPutLog("Delete single item function move file to trash can error", err)
				}
			}

			if err != nil {
				p.state = failure
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
				outPutLog("Delete multiple item function error", err)
				m.processBarModel.process[id] = p
				break
			} else {
				if p.done == p.total {
					p.state = successful
					channel <- channelMessage{
						messageId:       id,
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

func copyMultipleItem(m model) model {
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
		outPutLog("Copy multiple item function get file state error", panel.selected[0], err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.selected[0])

		if err != nil {
			outPutLog("Copy multiple item function read file error", err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			outPutLog("Copy multiple item function write file to clipboard error", err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func cutMultipleItem(m model) model {
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
		outPutLog("Copy multiple item function get file state error", panel.selected[0], err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.selected[0])

		if err != nil {
			outPutLog("Copy multiple item function read file error", err)
		}

		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			outPutLog("Copy multiple item function write file to clipboard error", err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func selectAllItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	for _, item := range panel.element {
		panel.selected = append(panel.selected, item.location)
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}
