package zoxide

import (
	"log/slog"
	"unicode"

	tea "charm.land/bubbletea/v2"
	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) Open() tea.Cmd {
	m.open = true
	m.justOpened = true
	m.textInput.SetValue("")
	_ = m.textInput.Focus()
	if common.Config.ZoxideMaxVisibleResults > 0 {
		m.maxVisibleResults = common.Config.ZoxideMaxVisibleResults
	} else if m.maxVisibleResults <= 0 {
		m.maxVisibleResults = DefaultMaxVisibleResults
	}

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
	m.textInput.SetWidth(width - modalInputPadding)
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < ZoxideMinHeight {
		slog.Warn("Zoxide initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = ZoxideMinHeight
	}
	m.maxHeight = maxHeight
}

// SetMaxVisibleResults updates the maximum visible results count.
func (m *Model) SetMaxVisibleResults(count int) {
	if count <= 0 {
		count = DefaultMaxVisibleResults
	}
	m.maxVisibleResults = count
}

// GetConfiguredMaxVisibleResults returns the configured maximum visible results.
func (m *Model) GetConfiguredMaxVisibleResults() int {
	if m.maxVisibleResults > 0 {
		return m.maxVisibleResults
	}
	if common.Config.ZoxideMaxVisibleResults > 0 {
		return common.Config.ZoxideMaxVisibleResults
	}
	return DefaultMaxVisibleResults
}

// GetMaxVisibleResults returns the effective maximum visible results count,
// taking into account the available height if maxHeight is constrained.
func (m *Model) GetMaxVisibleResults() int {
	configured := m.GetConfiguredMaxVisibleResults()
	if m.maxHeight > renderOverhead {
		available := m.maxHeight - renderOverhead
		if available > 0 && configured > available {
			return available
		}
	}
	return configured
}

func (m *Model) GetResults() []zoxidelib.Result {
	out := make([]zoxidelib.Result, len(m.results))
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
