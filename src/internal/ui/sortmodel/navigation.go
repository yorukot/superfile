package sortmodel

func (m *Model) ListUp() {
	m.Cursor = (m.Cursor - 1 + SortTypeCount) % SortTypeCount
}

func (m *Model) ListDown() {
	m.Cursor = (m.Cursor + 1 + SortTypeCount) % SortTypeCount
}
