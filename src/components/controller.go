package components

import (
	"path"
	"strconv"
)

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
	m.test = strconv.Itoa(len(panel.element) - PanelElementHeight(m.mainPanelHeight))
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
	m.test = strconv.Itoa(len(panel.element) - PanelElementHeight(m.mainPanelHeight))
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func SideBarSelectFolder(m model) model {
	m.sideBarModel.pinnedModel.selected = m.sideBarModel.pinnedModel.folder[m.sideBarModel.cursor].location
	m.fileModel.filePanels[m.filePanelFocusIndex].location = m.sideBarModel.pinnedModel.selected
	m.sideBarFocus = false
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	return m
}

func EnterPanel(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) > 0 && panel.element[panel.cursor].folder {
		panel.location = panel.element[panel.cursor].location
		m.fileModel.filePanels[m.filePanelFocusIndex] = panel
		m.fileModel.filePanels[m.filePanelFocusIndex].cursor = 0
		m.fileModel.filePanels[m.filePanelFocusIndex].render = 0
	}
	return m
}

func ParentFolder(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	fullPath := panel.location
	parentDir := path.Dir(fullPath)
	panel.location = parentDir
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}
