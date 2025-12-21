package filepanel

func (m *Model) scrollToCursor(mainPanelHeight int) {
	if m.Cursor < 0 || m.Cursor >= len(m.Element) {
		m.Cursor = 0
		m.RenderIndex = 0
		return
	}

	renderCount := panelElementHeight(mainPanelHeight)
	if m.Cursor < m.RenderIndex {
		m.RenderIndex = max(0, m.Cursor-renderCount+1)
	} else if m.Cursor > m.RenderIndex+renderCount-1 {
		m.RenderIndex = m.Cursor - renderCount + 1
	}
}

// Control file panel list up
func (m *Model) ListUp(mainPanelHeight int) {
	if len(m.Element) == 0 {
		return
	}
	if m.Cursor > 0 {
		m.Cursor--
		if m.Cursor < m.RenderIndex {
			m.RenderIndex--
		}
	} else {
		if len(m.Element) > panelElementHeight(mainPanelHeight) {
			m.RenderIndex = len(m.Element) - panelElementHeight(mainPanelHeight)
			m.Cursor = len(m.Element) - 1
		} else {
			m.Cursor = len(m.Element) - 1
		}
	}
}

// Control file panel list down
func (m *Model) ListDown(mainPanelHeight int) {
	if len(m.Element) == 0 {
		return
	}
	if m.Cursor < len(m.Element)-1 {
		m.Cursor++
		if m.Cursor > m.RenderIndex+panelElementHeight(mainPanelHeight)-1 {
			m.RenderIndex++
		}
	} else {
		m.RenderIndex = 0
		m.Cursor = 0
	}
}

func (m *Model) PgUp(mainPanelHeight int) {
	panlen := len(m.Element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 //nolint:mnd // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		m.Cursor = 0
	} else {
		if m.Cursor-panHeight <= 0 {
			m.Cursor = 0
			m.RenderIndex = 0
		} else {
			m.Cursor -= panHeight
			m.RenderIndex = m.Cursor - panCenter

			if m.RenderIndex < 0 {
				m.RenderIndex = 0
			}
		}
	}
}

func (m *Model) PgDown(mainPanelHeight int) {
	panlen := len(m.Element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 //nolint:mnd // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		m.Cursor = panlen - 1
	} else {
		if m.Cursor+panHeight >= panlen {
			m.Cursor = panlen - 1
			m.RenderIndex = m.Cursor - panCenter
		} else {
			m.Cursor += panHeight
			m.RenderIndex = m.Cursor - panCenter
		}
	}
}

// Handles the action of selecting an item in the file panel upwards. (only work on select mode)
// This basically just toggles the "selected" status of element that is pointed by the cursor
// and then moves the cursor up
// TODO : Add unit tests for ItemSelectUp and singleItemSelect
func (m *Model) ItemSelectUp(mainPanelHeight int) {
	m.SingleItemSelect()
	m.ListUp(mainPanelHeight)
}

// Handles the action of selecting an item in the file panel downwards. (only work on select mode)
func (m *Model) ItemSelectDown(mainPanelHeight int) {
	m.SingleItemSelect()
	m.ListDown(mainPanelHeight)
}

// Applies targetFile cursor positioning, if configured for the panel.
func (m *Model) applyTargetFileCursor() {
	for idx, el := range m.Element {
		if el.Name == m.TargetFile {
			m.Cursor = idx
			break
		}
	}
	m.TargetFile = ""
}
