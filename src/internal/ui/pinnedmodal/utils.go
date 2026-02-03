package pinnedmodal

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Open() tea.Cmd {
	m.open = true
	m.justOpened = true
	m.textInput.SetValue("")
	_ = m.textInput.Focus()

	m.results = m.allDirs
	m.cursor = 0
	m.renderIndex = 0
	return nil
}

func (m *Model) Close() {
	m.open = false
	m.textInput.Blur()
	m.textInput.SetValue("")
	m.results = []Directory{}
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
	if width < PinnedModalMinWidth {
		slog.Warn("PinnedModal initialized with too less width", "width", width)
		width = PinnedModalMinWidth
	}
	m.width = width
	m.textInput.Width = width - 6
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < PinnedModalMinHeight {
		slog.Warn("PinnedModal initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = PinnedModalMinHeight
	}
	m.maxHeight = maxHeight
}
