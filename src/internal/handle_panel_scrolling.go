package internal

// ======================================== File panel controller ========================================

// Control file panel list up
func (panel *filePanel) listUp(mainPanelHeight int) {
	if len(panel.element) == 0 {
		return
	}
	if panel.cursor > 0 {
		panel.cursor--
		if panel.cursor < panel.render {
			panel.render--
		}
	} else {
		if len(panel.element) > panelElementHeight(mainPanelHeight) {
			panel.render = len(panel.element) - panelElementHeight(mainPanelHeight)
			panel.cursor = len(panel.element) - 1
		} else {
			panel.cursor = len(panel.element) - 1
		}
	}
}

// Control file panel list down
func (panel *filePanel) listDown(mainPanelHeight int) {
	if len(panel.element) == 0 {
		return
	}
	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+panelElementHeight(mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
}

func (panel *filePanel) pgUp(mainPanelHeight int) {
	panlen := len(panel.element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		panel.cursor = 0
	} else {
		if panel.cursor-panHeight <= 0 {
			panel.cursor = 0
			panel.render = 0
		} else {
			panel.cursor -= panHeight
			panel.render = panel.cursor - panCenter

			if panel.render < 0 {
				panel.render = 0
			}
		}
	}
}

func (panel *filePanel) pgDown(mainPanelHeight int) {
	panlen := len(panel.element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		panel.cursor = panlen - 1
	} else {
		if panel.cursor+panHeight >= panlen {
			panel.cursor = panlen - 1
			panel.render = panel.cursor - panCenter
		} else {
			panel.cursor += panHeight
			panel.render = panel.cursor - panCenter
		}
	}
}

// Handles the action of selecting an item in the file panel upwards. (only work on select mode)
func (panel *filePanel) itemSelectUp(mainPanelHeight int) {
	if panel.cursor > 0 {
		panel.cursor--
		if panel.cursor < panel.render {
			panel.render--
		}
	} else {
		if len(panel.element) > panelElementHeight(mainPanelHeight) {
			panel.render = len(panel.element) - panelElementHeight(mainPanelHeight)
			panel.cursor = len(panel.element) - 1
		} else {
			panel.cursor = len(panel.element) - 1
		}
	}
	selectItemIndex := panel.cursor + 1
	if selectItemIndex > len(panel.element)-1 {
		selectItemIndex = 0
	}
	if arrayContains(panel.selected, panel.element[selectItemIndex].location) {
		panel.selected = removeElementByValue(panel.selected, panel.element[selectItemIndex].location)
	} else {
		panel.selected = append(panel.selected, panel.element[selectItemIndex].location)
	}
}

// Handles the action of selecting an item in the file panel downwards. (only work on select mode)
func (panel *filePanel) itemSelectDown(mainPanelHeight int) {

	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+panelElementHeight(mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
	selectItemIndex := panel.cursor - 1
	if selectItemIndex < 0 {
		selectItemIndex = len(panel.element) - 1
	}
	if arrayContains(panel.selected, panel.element[selectItemIndex].location) {
		panel.selected = removeElementByValue(panel.selected, panel.element[selectItemIndex].location)
	} else {
		panel.selected = append(panel.selected, panel.element[selectItemIndex].location)
	}

}

// ======================================== Metadata controller ========================================

// Control metadata panel up
func (fm *fileMetadata) listUp() {
	if len(fm.metaData) == 0 {
		return
	}
	if fm.renderIndex > 0 {
		fm.renderIndex--
	} else {
		fm.renderIndex = len(fm.metaData) - 1
	}
}

// Control metadata panel down
func (fm *fileMetadata) listDown() {
	if fm.renderIndex < len(fm.metaData)-1 {
		fm.renderIndex++
	} else {
		fm.renderIndex = 0
	}
}

// ======================================== Processbar controller ========================================

// Control processbar panel list up
// There is a shadowing happening here, but it will be removed
// Once we make footerHeight part of model struct
func (p *processBarModel) listUp(footerHeight int) {
	if len(p.processList) == 0 {
		return
	}
	if p.cursor > 0 {
		p.cursor--
		if p.cursor < p.render {
			p.render--
		}
	} else {
		if len(p.processList) <= 3 || (len(p.processList) <= 2 && footerHeight < 14) {
			p.cursor = len(p.processList) - 1
		} else {
			p.render = len(p.processList) - 3
			p.cursor = len(p.processList) - 1
		}
	}
}

// Control processbar panel list down
func (p *processBarModel) listDown() {
	if len(p.processList) == 0 {
		return
	}
	if p.cursor < len(p.processList)-1 {
		p.cursor++
		if p.cursor > p.render+2 {
			p.render++
		}
	} else {
		p.render = 0
		p.cursor = 0
	}
}
