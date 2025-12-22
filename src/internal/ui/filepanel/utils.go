package filepanel

import (
	"slices"
)

func (m *Model) GetSelectedItem() Element {
	if m.Cursor < 0 || len(m.Element) <= m.Cursor {
		return Element{}
	}
	return m.Element[m.Cursor]
}

func (m *Model) ResetSelected() {
	m.Selected = m.Selected[:0]
}

// For modification. Make sure to do a nil check
func (m *Model) GetSelectedItemPtr() *Element {
	if m.Cursor < 0 || len(m.Element) <= m.Cursor {
		return nil
	}
	return &m.Element[m.Cursor]
}

// Select the item where cursor located (only work on select mode)
func (m *Model) SingleItemSelect() {
	if len(m.Element) > 0 && m.Cursor >= 0 && m.Cursor < len(m.Element) {
		elementLocation := m.Element[m.Cursor].Location

		if slices.Contains(m.Selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			m.Selected = removeElementByValue(m.Selected, elementLocation)
		} else {
			m.Selected = append(m.Selected, elementLocation)
		}
	}
}

func (m *Model) ElemCount() int {
	return len(m.Element)
}

func (m *Model) Empty() bool {
	return m.ElemCount() == 0
}
