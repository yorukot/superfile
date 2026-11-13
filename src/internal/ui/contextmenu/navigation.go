package contextmenu

// ListUp highlights the previous entry, wrapping around at the top.
func (m *Model) ListUp() {
	if !m.open || len(m.items) == 0 {
		return
	}
	m.cursor = (m.cursor - 1 + len(m.items)) % len(m.items)
}

// ListDown highlights the next entry, wrapping around at the bottom.
func (m *Model) ListDown() {
	if !m.open || len(m.items) == 0 {
		return
	}
	m.cursor = (m.cursor + 1) % len(m.items)
}

// HighlightItem moves the cursor onto a specific entry. A click reports the entry
// it landed on through the same cursor the keyboard uses, so that activating by
// mouse and activating by keyboard cannot disagree.
func (m *Model) HighlightItem(index int) bool {
	if !m.open || index < 0 || index >= len(m.items) {
		return false
	}
	m.cursor = index
	return true
}

// HasItems reports whether the menu was opened with at least one entry.
func (m *Model) HasItems() bool {
	return len(m.items) > 0
}

// clamp limits v to the inclusive range lo..hi.
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
