package zoxide

import (
	"log/slog"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	zoxidelib "github.com/lazysegtree/go-zoxide"
)

func (m *Model) Open() tea.Cmd {
	m.open = true
	m.justOpened = true
	m.textInput.SetValue("")
	_ = m.textInput.Focus()

	// Return async command for initial query instead of blocking
	return m.GetQueryCmd("")
}

func (m *Model) Close() {
	m.open = false
	m.textInput.Blur()
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
	m.textInput.Width = width - ModalInputPadding
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < ZoxideMinHeight {
		slog.Warn("Zoxide initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = ZoxideMinHeight
	}
	m.maxHeight = maxHeight
}

func (m *Model) GetResults() []zoxidelib.Result {
	out := make([]zoxidelib.Result, len(m.results))
	copy(out, m.results)
	return out
}

func (m *Model) GetTextInputValue() string {
	return m.textInput.Value()
}

func isKeyAlphaNum(msg tea.KeyMsg) bool {
	r := []rune(msg.String())
	if len(r) != 1 {
		return false
	}
	return unicode.IsLetter(r[0]) || unicode.IsNumber(r[0])
}
