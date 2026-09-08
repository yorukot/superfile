package filepanel

import "github.com/yorukot/superfile/src/internal/search"

// Clears results and cursor. Enter and Exit own Active and Root.
func (m *Model) resetSearchState() {
	m.Search.Results = nil
	m.Search.MatchCount = 0
	m.Search.UnreadableDirs = 0
	m.Search.Done = false
	m.Search.Cursor = 0
	m.Search.RenderIndex = 0
}

// Starts a session rooted at current location.
func (m *Model) EnterSearchMode() {
	m.Search.Active = true
	m.Search.Root = m.Location
	m.resetSearchState()
	m.SearchBar.SetValue("")
	m.SearchBar.Focus()
}

func (m *Model) ExitSearchMode() {
	m.Search.Active = false
	m.Search.Root = ""
	m.resetSearchState()
	m.SearchBar.Blur()
	m.SearchBar.SetValue("")
}

func (m *Model) ApplySearchProgress(p search.Progress) {
	m.Search.Results = p.Results
	m.Search.MatchCount = p.MatchCount
	m.Search.UnreadableDirs = p.UnreadableDirs
	m.Search.Done = p.Done
	m.Search.Cursor = min(m.Search.Cursor, max(0, len(m.Search.Results)-1))
	if len(m.Search.Results) == 0 {
		m.Search.Cursor = 0
		m.Search.RenderIndex = 0
	}
}

func (m *Model) GetSearchCursorResult() *search.Result {
	if len(m.Search.Results) == 0 {
		return nil
	}
	return &m.Search.Results[m.Search.Cursor]
}

func (m *Model) searchScrollToCursor(cursor int) {
	if cursor < 0 || cursor >= len(m.Search.Results) {
		return
	}
	m.Search.Cursor = cursor

	renderCount := m.PanelElementHeight()
	if m.Search.Cursor < m.Search.RenderIndex {
		m.Search.RenderIndex = m.Search.Cursor
	} else if m.Search.Cursor > m.Search.RenderIndex+renderCount-1 {
		m.Search.RenderIndex = m.Search.Cursor - renderCount + 1
	}
}

func (m *Model) searchMoveCursorBy(delta int) {
	if len(m.Search.Results) == 0 {
		return
	}
	cursor := (m.Search.Cursor + delta + len(m.Search.Results)) % len(m.Search.Results)
	m.searchScrollToCursor(cursor)
}

func (m *Model) SearchListUp() {
	m.searchMoveCursorBy(-1)
}

func (m *Model) SearchListDown() {
	m.searchMoveCursorBy(1)
}

func (m *Model) SearchPgUp() {
	m.searchPageScroll(-m.getPageScrollSize())
}

func (m *Model) SearchPgDown() {
	m.searchPageScroll(m.getPageScrollSize())
}

func (m *Model) searchPageScroll(delta int) {
	if len(m.Search.Results) == 0 {
		return
	}
	cursor := m.Search.Cursor + delta
	if cursor < 0 {
		cursor = 0
	} else if cursor >= len(m.Search.Results) {
		cursor = len(m.Search.Results) - 1
	}
	m.searchScrollToCursor(cursor)
}
