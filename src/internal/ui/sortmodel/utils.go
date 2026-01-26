package sortmodel

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open(curSortKind SortKind) {
	m.Cursor = int(curSortKind)
	m.open = true
}

func (m *Model) Close() {
	m.open = false
	m.Cursor = 0
}

func (m *Model) GetSelectedKind() SortKind {
	return SortKind(m.Cursor)
}
