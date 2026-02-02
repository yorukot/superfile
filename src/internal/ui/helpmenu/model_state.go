package helpmenu

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/common"
)

func New() Model {
	data := getData()

	return Model{
		renderIndex:  0,
		cursor:       1,
		data:         data,
		filteredData: data,
		opened:       false,
		searchBar:    common.GenerateSearchBar(),
	}
}

// Toggle help menu
func (m *Model) Open() {
	if m.opened {
		m.searchBar.Reset()
		m.opened = false
		return
	}

	// Reset filteredData to the full data whenever the helpMenu is opened
	m.filteredData = m.data
	m.opened = true
}

// Quit help menu
func (m *Model) Close() {
	m.searchBar.Reset()
	m.opened = false
}

// Check hotkey input in help menu. Possible actions are moving up, down
// and quiting the menu
func (m *Model) HandleKey(msg string) {
	if m.searchBar.Focused() {
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg), slices.Contains(common.Hotkeys.CancelTyping, msg):
			m.searchBar.Blur()
		default:
			m.filter(m.searchBar.Value())
		}
	} else {
		m.handleNavKeys(msg)
	}
}

func (m *Model) handleNavKeys(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.ListUp()
	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.ListDown()
	case slices.Contains(common.Hotkeys.Quit, msg):
		m.Close()
	case slices.Contains(common.Hotkeys.SearchBar, msg):
		m.searchBar.Focus()
	}
}

func (m *Model) HandleTeaMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.searchBar.Focused() {
		m.searchBar, cmd = m.searchBar.Update(msg)
	}
	return cmd
}

func (m *Model) SetDimensions(width int, height int) {
	m.width = width
	m.height = height
	// 2 for border, 1 for left padding, 2 for placeholder icon of searchbar
	// 1 for additional character that View() of search bar function mysteriously adds.
	m.searchBar.Width = m.width - (common.InnerPadding + common.BorderPadding)
}
