package find

import (
	"log/slog"
	"os/exec"
	"unicode"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) Open(searchDir string) tea.Cmd {
	m.searchDir = searchDir
	_, lookupErr := exec.LookPath("fd")
	m.fdFound = lookupErr == nil
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
	m.results = []FindResult{}
	m.errMsg = ""
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
	if width < FindMinWidth {
		slog.Warn("Find initialized with too less width", "width", width)
		width = FindMinWidth
	}
	m.width = width
	// Excluding borders(2), SpacePadding(1), Prompt(2), and one extra character that is appended
	// by textInput.View()
	m.textInput.SetWidth(width - modalInputPadding)
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < FindMinHeight {
		slog.Warn("Find initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = FindMinHeight
	}
	m.maxHeight = maxHeight
}

func (m *Model) GetResults() []FindResult {
	out := make([]FindResult, len(m.results))
	copy(out, m.results)
	return out
}

func (m *Model) GetTextInputValue() string {
	return m.textInput.Value()
}

func isKeyAlphaNum(msg tea.KeyPressMsg) bool {
	r := []rune(msg.String())
	if len(r) != 1 {
		return false
	}
	return unicode.IsLetter(r[0]) || unicode.IsNumber(r[0])
}
