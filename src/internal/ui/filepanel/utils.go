package filepanel

import "math"

func (m *Model) GetFocusedItem() Element {
	if m.Cursor < 0 || len(m.Element) <= m.Cursor {
		return Element{}
	}
	return m.Element[m.Cursor]
}

func (m *Model) ResetSelected() {
	m.selectOrderCounter = 0
	m.selected = make(map[string]int)
}

// For modification. Make sure to do a nil check
func (m *Model) GetFocusedItemPtr() *Element {
	if m.Cursor < 0 || len(m.Element) <= m.Cursor {
		return nil
	}
	return &m.Element[m.Cursor]
}

func (m *Model) CheckSelected(location string) bool {
	_, isSelected := m.selected[location]
	return isSelected
}

func (m *Model) GetSelectedLocations() []string {
	result := make([]string, 0, len(m.selected))
	for k := range m.selected {
		result = append(result, k)
	}
	return result
}

func (m *Model) GetFirstSelectedLocation() string {
	if len(m.selected) == 0 {
		return ""
	}
	result := ""
	minOrder := math.MaxInt
	for location, order := range m.selected {
		if minOrder > order {
			result = location
			minOrder = order
		}
	}
	return result
}

// Select the item where cursor located (only work on select mode)
func (m *Model) SingleItemSelect() {
	if len(m.Element) > 0 && m.Cursor >= 0 && m.Cursor < len(m.Element) {
		elementLocation := m.Element[m.Cursor].Location

		m.SetSelected(elementLocation, !m.CheckSelected(elementLocation))
	}
}

func (m *Model) ElemCount() int {
	return len(m.Element)
}

func (m *Model) SelectedCount() uint {
	if m.selected == nil {
		return 0
	}
	return uint(len(m.selected))
}

func (m *Model) Empty() bool {
	return m.ElemCount() == 0
}

func (m *Model) EmptyOrInvalid() bool {
	return m.Empty() || m.ValidateCursorAndRenderIndex() != nil
}
