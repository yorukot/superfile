package pinnedmodal

import (
	"slices"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

func DefaultModel(maxHeight int, width int) Model {
	return GenerateModel(maxHeight, width)
}

func GenerateModel(maxHeight int, width int) Model {
	m := Model{
		headline:  pinnedModalHeadlineText,
		open:      false,
		textInput: common.GeneratePromptTextInput(),
		results:   []Directory{},
	}
	m.SetMaxHeight(maxHeight)
	m.SetWidth(width)
	m.textInput.Prompt = ""
	return m
}

func (m *Model) HandleUpdate(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !m.IsOpen() {
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			action = m.handleConfirm()
			m.Close()
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			m.Close()
		case slices.Contains(common.Hotkeys.ListUp, msg.String()) && !isKeyAlphaNum(msg):
			m.navigateUp()
		case slices.Contains(common.Hotkeys.ListDown, msg.String()) && !isKeyAlphaNum(msg):
			m.navigateDown()
		case slices.Contains(common.Hotkeys.GotoPinned, msg.String()) && m.justOpened:
			m.justOpened = false
		default:
			cmd = m.handleNormalKeyInput(msg)
		}
	default:
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return action, cmd
}

func (m *Model) handleConfirm() common.ModelAction {
	if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
		selectedDir := m.results[m.cursor]
		return common.CDCurrentPanelAction{
			Location: selectedDir.Location,
		}
	}
	return common.NoAction{}
}

func (m *Model) handleNormalKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	m.FilterPinnedDirs(m.textInput.Value())
	return cmd
}

func (m *Model) GetQueryCmd(query string) tea.Cmd {
	return func() tea.Msg {
		// Work on a copy of the model to avoid mutating shared state
		mCopy := *m
		mCopy.FilterPinnedDirs(query)
		return NewUpdateMsg(query, mCopy.results)
	}
}

func (m *Model) LoadPinnedDirs(dirs []Directory) {
	m.allDirs = dirs
	m.results = dirs
	m.cursor = 0
	m.renderIndex = 0
}

func (m *Model) FilterPinnedDirs(query string) {
	if query == "" {
		m.results = m.allDirs
		m.cursor = 0
		m.renderIndex = 0
		return
	}

	if len(m.allDirs) == 0 {
		m.results = []Directory{}
		return
	}

	var filteredDirs []Directory

	haystack := make([]string, len(m.allDirs))
	for i, dir := range m.allDirs {
		searchText := dir.Name + " " + dir.Location
		haystack[i] = searchText
	}

	for _, match := range utils.FzfSearch(query, haystack) {
		if match.HayIndex >= 0 && match.HayIndex < len(m.allDirs) {
			filteredDirs = append(filteredDirs, m.allDirs[match.HayIndex])
		}
	}

	m.results = filteredDirs
	m.cursor = 0
	m.renderIndex = 0
}

func (msg UpdateMsg) Apply(m *Model) tea.Cmd {
	currentQuery := m.textInput.Value()
	if msg.query != currentQuery {
		return nil
	}

	m.results = msg.results
	m.cursor = 0
	m.renderIndex = 0

	return nil
}

func isKeyAlphaNum(msg tea.KeyMsg) bool {
	r := []rune(msg.String())
	if len(r) != 1 {
		return false
	}
	return unicode.IsLetter(r[0]) || unicode.IsNumber(r[0])
}
