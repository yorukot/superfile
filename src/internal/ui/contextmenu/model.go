package contextmenu

// Open shows the menu at the given terminal coordinates with the given entries.
//
// The requested position is clamped so the menu stays fully on screen, which is
// what keeps the drawn menu and the published hit-test regions in agreement: both
// read the clamped origin stored here rather than recomputing it.
//
// A menu with no entries is not opened at all, since there would be nothing to
// show and nothing to choose.
func (m *Model) Open(x, y, screenWidth, screenHeight int, items []Item) {
	if len(items) == 0 {
		m.Close()
		return
	}

	m.open = true
	m.items = items
	m.cursor = 0
	m.width = contentWidth(items)
	m.height = len(items) + borderRows

	m.originX = clamp(x, 0, max(0, screenWidth-m.width))
	m.originY = clamp(y, 0, max(0, screenHeight-m.height))
}

// Close hides the menu and forgets its entries, so a later reopen starts from a
// clean state rather than from whatever was offered last time.
func (m *Model) Close() {
	m.open = false
	m.items = nil
	m.cursor = 0
}

// IsOpen reports whether the menu is currently shown.
func (m *Model) IsOpen() bool {
	return m.open
}

// OriginX returns the leftmost terminal column of the menu.
func (m *Model) OriginX() int {
	return m.originX
}

// OriginY returns the topmost terminal row of the menu.
func (m *Model) OriginY() int {
	return m.originY
}

// Width returns the total width of the rendered menu.
func (m *Model) Width() int {
	return m.width
}

// Height returns the total height of the rendered menu.
func (m *Model) Height() int {
	return m.height
}

// Items returns the entries currently offered.
func (m *Model) Items() []Item {
	return m.items
}

// SelectedAction returns the action of the highlighted entry. It reports false
// when the menu is closed or when the cursor is not on a real entry.
func (m *Model) SelectedAction() (Action, bool) {
	item, ok := m.itemAt(m.cursor)
	if !ok {
		return ActionNone, false
	}
	return item.Action, true
}

// ChosenAction returns the action of the entry at the given index. It reports
// false when the index is out of range, which is how a stale hit-test region is
// rejected.
func (m *Model) ChosenAction(index int) (Action, bool) {
	item, ok := m.itemAt(index)
	if !ok {
		return ActionNone, false
	}
	return item.Action, true
}

func (m *Model) itemAt(index int) (Item, bool) {
	if !m.open || index < 0 || index >= len(m.items) {
		return Item{}, false
	}
	return m.items[index], true
}
