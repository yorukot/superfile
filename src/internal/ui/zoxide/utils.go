package zoxide

import (
	"log/slog"

	zoxidelib "github.com/lazysegtree/go-zoxide"
)

func (m *Model) Open() {
	m.open = true
	m.justOpened = true
	m.textInput.SetValue("") // Clear any unwanted characters
	_ = m.textInput.Focus()
	m.updateSuggestions()
}

func (m *Model) Close() {
	m.open = false
	m.textInput.SetValue("")
	m.results = []zoxidelib.Result{}
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
	if width < ZoxideMinWidth {
		slog.Warn("Zoxide initialized with too less width", "width", width)
		width = ZoxideMinWidth
	}
	m.width = width
	// Excluding borders(2), SpacePadding(1), Prompt(2), and one extra character that is appended
	// by textInput.View()
	m.textInput.Width = width - 2 - 1 - 2 - 1
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < ZoxideMinHeight {
		slog.Warn("Zoxide initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = ZoxideMinHeight
	}
	m.maxHeight = maxHeight
}

func (m *Model) GetResults() []zoxidelib.Result {
	return m.results
}
