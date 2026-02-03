package gotointeractive

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Open() tea.Cmd {
	m.open = true
	m.justOpened = true
	m.textInput.SetValue("")
	_ = m.textInput.Focus()

	return m.GetQueryCmd("")
}

func (m *Model) Close() {
	m.open = false
	m.textInput.Blur()
	m.textInput.SetValue("")
	m.results = []string{}
	m.cursor = 0
	m.renderIndex = 0
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) GetMaxHeight() int {
	return m.maxHeight
}

func (m *Model) SetWidth(width int) {
	if width < GotoMinWidth {
		slog.Warn("Goto initialized with too less width", "width", width)
		width = GotoMinWidth
	}
	m.width = width
	m.textInput.Width = width - modalInputPadding
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < GotoMinHeight {
		slog.Warn("Goto initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = GotoMinHeight
	}
	m.maxHeight = maxHeight
}

func (m *Model) GetResults() []string {
	out := make([]string, len(m.results))
	copy(out, m.results)
	return out
}

func (m *Model) GetTextInputValue() string {
	return m.textInput.Value()
}

func (m *Model) SetCurrentPath(path string) {
	m.currentPath = path
	m.updateResults()
}

func (m *Model) GetCurrentPath() string {
	return m.currentPath
}
