package sortmodel

func (m *Model) ListUp() {
	m.Cursor = (m.Cursor - 1 + len(SortOptionsStr)) % len(SortOptionsStr)
}

func (m *Model) ListDown() {
	m.Cursor = (m.Cursor + 1 + len(SortOptionsStr)) % len(SortOptionsStr)
}
