package components

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
	m.sideBarFocus = false
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

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
	return m
}

func PreviousFilePanel(m model) model {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
	return m
}

func CloseFilePanel(m model) model {
	if len(m.fileModel.filePanels) != 1 {
		m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

		if m.filePanelFocusIndex != 0 {
			m.filePanelFocusIndex--
		}
		m.fileModel.width = (m.fullWidth - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
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
		m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.sideBarFocus)
		m.fileModel.width = (m.fullWidth - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.filePanelFocusIndex++
	}
	return m
}

func FocusOnSideBar(m model) model {
	if m.sideBarFocus {
		m.sideBarFocus = false
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.sideBarFocus = true
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
	return m
}
