package filepanel

func (panel *FilePanel) scrollToCursor(mainPanelHeight int) {
	if panel.Cursor < 0 || panel.Cursor >= len(panel.Element) {
		panel.Cursor = 0
		panel.RenderIndex = 0
		return
	}

	renderCount := panelElementHeight(mainPanelHeight)
	if panel.Cursor < panel.RenderIndex {
		panel.RenderIndex = max(0, panel.Cursor-renderCount+1)
	} else if panel.Cursor > panel.RenderIndex+renderCount-1 {
		panel.RenderIndex = panel.Cursor - renderCount + 1
	}
}

// Control file panel list up
func (panel *FilePanel) ListUp(mainPanelHeight int) {
	if len(panel.Element) == 0 {
		return
	}
	if panel.Cursor > 0 {
		panel.Cursor--
		if panel.Cursor < panel.RenderIndex {
			panel.RenderIndex--
		}
	} else {
		if len(panel.Element) > panelElementHeight(mainPanelHeight) {
			panel.RenderIndex = len(panel.Element) - panelElementHeight(mainPanelHeight)
			panel.Cursor = len(panel.Element) - 1
		} else {
			panel.Cursor = len(panel.Element) - 1
		}
	}
}

// Control file panel list down
func (panel *FilePanel) ListDown(mainPanelHeight int) {
	if len(panel.Element) == 0 {
		return
	}
	if panel.Cursor < len(panel.Element)-1 {
		panel.Cursor++
		if panel.Cursor > panel.RenderIndex+panelElementHeight(mainPanelHeight)-1 {
			panel.RenderIndex++
		}
	} else {
		panel.RenderIndex = 0
		panel.Cursor = 0
	}
}

func (panel *FilePanel) PgUp(mainPanelHeight int) {
	panlen := len(panel.Element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 //nolint:mnd // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		panel.Cursor = 0
	} else {
		if panel.Cursor-panHeight <= 0 {
			panel.Cursor = 0
			panel.RenderIndex = 0
		} else {
			panel.Cursor -= panHeight
			panel.RenderIndex = panel.Cursor - panCenter

			if panel.RenderIndex < 0 {
				panel.RenderIndex = 0
			}
		}
	}
}

func (panel *FilePanel) PgDown(mainPanelHeight int) {
	panlen := len(panel.Element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 //nolint:mnd // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		panel.Cursor = panlen - 1
	} else {
		if panel.Cursor+panHeight >= panlen {
			panel.Cursor = panlen - 1
			panel.RenderIndex = panel.Cursor - panCenter
		} else {
			panel.Cursor += panHeight
			panel.RenderIndex = panel.Cursor - panCenter
		}
	}
}

// Handles the action of selecting an item in the file panel upwards. (only work on select mode)
// This basically just toggles the "selected" status of element that is pointed by the cursor
// and then moves the cursor up
// TODO : Add unit tests for ItemSelectUp and singleItemSelect
func (panel *FilePanel) ItemSelectUp(mainPanelHeight int) {
	panel.SingleItemSelect()
	panel.ListUp(mainPanelHeight)
}

// Handles the action of selecting an item in the file panel downwards. (only work on select mode)
func (panel *FilePanel) ItemSelectDown(mainPanelHeight int) {
	panel.SingleItemSelect()
	panel.ListDown(mainPanelHeight)
}

// Applies targetFile cursor positioning, if configured for the panel.
func (panel *FilePanel) applyTargetFileCursor() {
	for idx, el := range panel.Element {
		if el.Name == panel.TargetFile {
			panel.Cursor = idx
			break
		}
	}
	panel.TargetFile = ""
}
