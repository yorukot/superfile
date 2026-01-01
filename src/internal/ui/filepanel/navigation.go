package filepanel

import (
	"fmt"
)

func (m *Model) scrollToCursor(cursor int) {
	if cursor < 0 || cursor >= m.ElemCount() {
		return
	}
	m.Cursor = cursor

	// Modify renderIndex if needed
	renderCount := m.PanelElementHeight()
	if m.Cursor < m.RenderIndex {
		// Due to size change, when last element is selected, we might have
		// empty space (RenderIndex ... ElemCount()-1 spans less then renderCount)
		// Even with >0 RenderIndex
		m.RenderIndex = m.Cursor
	} else if m.Cursor > m.RenderIndex+renderCount-1 {
		m.RenderIndex = m.Cursor - renderCount + 1
	}
}

func (m *Model) moveCursorBy(delta int) {
	if m.Empty() {
		return
	}
	// Wrap cursor
	cursor := (m.Cursor + delta + m.ElemCount()) % m.ElemCount()
	m.scrollToCursor(cursor)
}

// Control file panel list up
func (m *Model) ListUp() {
	m.moveCursorBy(-1)
}

// Control file panel list down
func (m *Model) ListDown() {
	m.moveCursorBy(1)
}

func (m *Model) PgUp() {
	m.moveCursorBy(-m.getPageScrollSize())
}

func (m *Model) PgDown() {
	m.moveCursorBy(m.getPageScrollSize())
}

// Handles the action of selecting an item in the file panel upwards. (only work on select mode)
// This basically just toggles the "selected" status of element that is pointed by the cursor
// and then moves the cursor up
// TODO : Add unit tests for ItemSelectUp and singleItemSelect
func (m *Model) ItemSelectUp() {
	m.SingleItemSelect()
	m.ListUp()
}

// Handles the action of selecting an item in the file panel downwards. (only work on select mode)
func (m *Model) ItemSelectDown() {
	m.SingleItemSelect()
	m.ListDown()
}

// Applies targetFile cursor positioning, if configured for the panel.
func (m *Model) applyTargetFileCursor() {
	for idx, el := range m.Element {
		if el.Name == m.TargetFile {
			m.scrollToCursor(idx)
			break
		}
	}
	m.TargetFile = ""
}

func (m *Model) ValidateCursorAndRenderIndex() error {
	if m.Cursor < 0 || m.ElemCount() <= m.Cursor {
		return fmt.Errorf("invalid cursor : %d, element count : %d", m.Cursor, m.ElemCount())
	}
	renderCount := m.PanelElementHeight()
	if (m.Cursor < m.RenderIndex) || (m.Cursor > m.RenderIndex+renderCount-1) {
		return fmt.Errorf("invalid renderIndex : %d, cursor : %d, renderCount : %d",
			m.RenderIndex, m.Cursor, renderCount)
	}
	return nil
}
