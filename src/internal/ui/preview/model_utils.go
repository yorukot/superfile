package preview

import "log/slog"

func (m *Model) GetContent() string {
	return m.content
}

func (m *Model) GetContentWidth() int {
	return m.contentWidth
}

func (m *Model) GetContentHeight() int {
	return m.contentHeight
}

func (m *Model) GetLocation() string {
	return m.location
}

func (m *Model) SetOpen(open bool) {
	m.open = open
}

func (m *Model) SetLocation(location string) {
	m.location = location
}

func (m *Model) SetLoading() {
	m.loading = true
}

// All content change happen via this only, to ensure the sync between
// content and width x height, and the loading variable reset
func (m *Model) setContent(content string, width int, height int, location string) {
	m.content = content
	m.contentWidth = width
	m.contentHeight = height
	m.location = location
	m.loading = false
}

func (m *Model) SetEmptyWithDimensions(width int, height int) {
	m.setContent(m.RenderTextWithDimension("", height, width), width, height, "")
}

func (m *Model) IsLoading() bool {
	return m.loading
}

func (m *Model) ToggleOpen() {
	m.open = !m.open
}

func (m *Model) CleanUp() {
	if m.thumbnailGenerator != nil {
		err := m.thumbnailGenerator.CleanUp()
		if err != nil {
			slog.Error("Error While cleaning up TempDirectory", "error", err)
		}
	}
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) Close() {
	m.open = false
}
