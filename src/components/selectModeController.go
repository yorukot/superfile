package components

import (
	"github.com/atotto/clipboard"
	"os"
	"path"
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
		for _, item := range panel.selected {
			filePath := item
			err := MoveFile(item, Config.TrashCanPath+"/"+path.Base(filePath))
			if err != nil {
				OutputLog("Error delete multiple item")
				OutputLog(err)
			}
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CopyMultipleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.items = panel.selected
	fileInfo, err := os.Stat(panel.selected[0])
	if err != nil {
		OutputLog("Can't find this file or folder")
		OutputLog(panel.selected[0])
		OutputLog(err)
	}

	if !fileInfo.IsDir() && float64(fileInfo.Size())/(1024*1024) < 250 {
		fileContent, err := os.ReadFile(panel.selected[0])

		if err != nil {
			OutputLog(err)
		}
		if err := clipboard.WriteAll(string(fileContent)); err != nil {
			OutputLog(err)
		}
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func CutMultipleItem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.copyItems.items = panel.selected
	m.copyItems.cut = true
	m.copyItems.oringnalPanel = orignalPanel{
		index:    m.filePanelFocusIndex,
		location: panel.location,
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
