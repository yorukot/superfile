package helpmenu

func (m *Model) IsOpen() bool {
	return m.opened
}

func (m *Model) GetHeight() int {
	return m.height
}

func (m *Model) GetWidth() int {
	return m.width
}
