package filepanel

func (m *Model) scrollToCursor(cursor int, mainPanelHeight int) {
	if cursor < 0 || cursor >= m.ElemCount() {
		return
	}
	m.Cursor = cursor

	// Modify renderIndex if needed
	renderCount := panelElementHeight(mainPanelHeight)
	if m.Cursor < m.RenderIndex {
		// Due to size change, when last element is selected, we might have
		// empty space (RenderIndex ... ElemCount()-1 spans less then renderCount)
		// Even with >0 RenderIndex
		m.RenderIndex = m.Cursor
	} else if m.Cursor > m.RenderIndex+renderCount-1 {
		m.RenderIndex = m.Cursor - renderCount + 1
	}
}

func (m *Model) moveCursorBy(delta int, mainPanelHeight int) {
	if m.Empty() {
		return
	}
	// Wrap cursor
	cursor := (m.Cursor + delta + m.ElemCount()) % m.ElemCount()
	m.scrollToCursor(cursor, mainPanelHeight)
}

// Control file panel list up
func (m *Model) ListUp(mainPanelHeight int) {
	m.moveCursorBy(-1, mainPanelHeight)
}

// Control file panel list down
func (m *Model) ListDown(mainPanelHeight int) {
	m.moveCursorBy(1, mainPanelHeight)
}

func (m *Model) PgUp(mainPanelHeight int) {
	m.moveCursorBy(-getScrollSize(mainPanelHeight), mainPanelHeight)
}

func (m *Model) PgDown(mainPanelHeight int) {
	m.moveCursorBy(getScrollSize(mainPanelHeight), mainPanelHeight)
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
func (m *Model) applyTargetFileCursor(mainPanelHeight int) {
	for idx, el := range m.Element {
		if el.Name == m.TargetFile {
			m.scrollToCursor(idx, mainPanelHeight)
			break
		}
	}
	m.TargetFile = ""
}
