package filepanel

import "math"

func (m *Model) GetCursor() int {
	return m.cursor
}

func (m *Model) GetRenderIndex() int {
	return m.renderIndex
}

func (m *Model) GetFocusedItem() Element {
	return m.GetElementAtIdx(m.GetCursor())
}

func (m *Model) GetElementAtIdx(idx int) Element {
	if idx < 0 || m.ElemCount() <= idx {
		return Element{}
	}
	return m.element[idx]
}

func (m *Model) GetFirstElement() Element {
	return m.GetElementAtIdx(0)
}

func (m *Model) ResetSelected() {
	m.selectOrderCounter = 0
	m.selected = make(map[string]int)
}

// For modification. Make sure to do a nil check
func (m *Model) GetFocusedItemPtr() *Element {
	if m.GetCursor() < 0 || m.ElemCount() <= m.GetCursor() {
		return nil
	}
	return &m.element[m.GetCursor()]
}

// Note : If this is called on an already selected element
// it will make its order last. This is expected behaviour
func (m *Model) SetSelected(location string) {
	m.selectOrderCounter++
	m.selected[location] = m.selectOrderCounter
}

func (m *Model) SetUnSelected(location string) {
	if m.CheckSelected(location) {
		delete(m.selected, location)
	}
}

func (m *Model) ToggleSelected(location string) {
	if m.CheckSelected(location) {
		delete(m.selected, location)
		return
	}
	m.SetSelected(location)
}

// Only used in tests, including tests outside this package
func (m *Model) SetSelectedAll(locations []string) {
	for _, location := range locations {
		m.SetSelected(location)
	}
}

func (m *Model) CheckSelected(location string) bool {
	_, isSelected := m.selected[location]
	return isSelected
}

// Returns an unordered list of selected locations
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
	if !m.EmptyOrInvalid() {
		m.ToggleSelected(m.GetFocusedItem().Location)
	}
}

func (m *Model) ElemCount() int {
	return len(m.element)
}

func (m *Model) SelectedCount() uint {
	return uint(len(m.selected))
}

func (m *Model) Empty() bool {
	return m.ElemCount() == 0
}

func (m *Model) EmptyOrInvalid() bool {
	return m.Empty() || m.ValidateCursorAndRenderIndex() != nil
}

func (m *Model) ToggleReverseSort() {
	m.SortReversed = !m.SortReversed
}

// SetCursorPosition sets cursor and updates renderIndex accordingly.
// Note: Intended for test utilities only!!!!!
func (m *Model) SetCursorPosition(cursor int) {
	m.scrollToCursor(cursor)
}

func (m *Model) FindElementIndexByName(name string) int {
	for i, elem := range m.element {
		if elem.Name == name {
			return i
		}
	}
	return -1
}

func (m *Model) FindElementIndexByLocation(location string) int {
	for i, elem := range m.element {
		if elem.Location == location {
			return i
		}
	}
	return -1
}
