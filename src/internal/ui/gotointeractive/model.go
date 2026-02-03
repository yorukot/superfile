package gotointeractive

import (
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func DefaultModel(maxHeight int, width int, cwd string) Model {
	return GenerateModel(cwd, maxHeight, width)
}

func GenerateModel(cwd string, maxHeight int, width int) Model {
	m := Model{
		headline:    icon.Search + icon.Space + gotoHeadlineText,
		open:        false,
		results:     []string{},
		currentPath: cwd,
	}
	m.SetMaxHeight(maxHeight)
	m.SetWidth(width)
	m.textInput = common.GeneratePromptTextInput()
	m.textInput.Prompt = getGotoPrompt()
	return m
}

func (m *Model) HandleUpdate(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	slog.Debug("gotointeractive.Model HandleUpdate()", "msg", msg,
		"msgType", reflect.TypeOf(msg),
		"textInput", m.textInput.Value(),
		"cursorBlink", m.textInput.Cursor.Blink)
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !m.IsOpen() {
		slog.Error("HandleUpdate called on closed goto")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			action = m.handleConfirm()
			if _, ok := action.(common.NoAction); !ok {
				m.Close()
			}
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()),
			slices.Contains(common.Hotkeys.Quit, msg.String()):
			m.Close()
		case msg.String() == "tab":
			action = m.handleTabCompletion()
		case msg.Type == tea.KeyBackspace && m.textInput.Value() == "":
			m.handleGoUp()
		case slices.Contains(common.Hotkeys.ListUp, msg.String()):
			m.navigateUp()
		case slices.Contains(common.Hotkeys.ListDown, msg.String()):
			m.navigateDown()
		case slices.Contains(common.Hotkeys.PageUp, msg.String()):
			m.navigatePageUp()
		case slices.Contains(common.Hotkeys.PageDown, msg.String()):
			m.navigatePageDown()
		case slices.Contains(common.Hotkeys.OpenGotoInteractive, msg.String()) && m.justOpened:
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
	input := strings.TrimSpace(m.textInput.Value())

	if input == ".." {
		m.handleGoUp()
		return common.NoAction{}
	}

	if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
		selectedResult := m.results[m.cursor]

		if selectedResult == ".." {
			m.handleGoUp()
			return common.NoAction{}
		}

		fullSelectedPath := filepath.Join(m.currentPath, selectedResult)

		fileInfo, err := os.Stat(fullSelectedPath)
		if err == nil {
			if fileInfo.IsDir() {
				return common.CDCurrentPanelAction{
					Location: fullSelectedPath,
				}
			}
			return common.OpenPanelAction{
				Location: fullSelectedPath,
			}
		}
	}

	return common.NoAction{}
}

func (m *Model) handleTabCompletion() common.ModelAction { //nolint:unparam // Return value kept for future extensibility
	input := strings.TrimSpace(m.textInput.Value())

	if m.cursor >= 0 && m.cursor < len(m.results) {
		selected := m.results[m.cursor]

		if selected == ".." {
			m.handleGoUp()
			return common.NoAction{}
		}

		fullSelectedPath := filepath.Join(m.currentPath, selected)

		fileInfo, err := os.Stat(fullSelectedPath)
		if err == nil {
			if fileInfo.IsDir() {
				m.currentPath = fullSelectedPath
				m.textInput.SetValue("")
				m.updateResults()
				return common.NoAction{}
			}
			m.textInput.SetValue(selected)
			return common.NoAction{}
		}
	}

	if input == ".." {
		m.handleGoUp()
		return common.NoAction{}
	}

	if input == "" {
		return common.NoAction{}
	}

	matched := m.filterResults(input)
	if len(matched) == 1 {
		selected := matched[0]
		fullSelectedPath := filepath.Join(m.currentPath, selected)

		fileInfo, err := os.Stat(fullSelectedPath)
		if err == nil {
			if fileInfo.IsDir() {
				m.currentPath = fullSelectedPath
				m.textInput.SetValue("")
				m.updateResults()
				return common.NoAction{}
			}
			m.textInput.SetValue(selected)
		}
	}

	return common.NoAction{}
}

func (m *Model) handleGoUp() {
	parentPath := filepath.Dir(m.currentPath)
	if parentPath != m.currentPath {
		m.currentPath = parentPath
		m.textInput.SetValue("")
		m.updateResults()
	}
}

func (m *Model) handleNormalKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return tea.Batch(cmd, m.GetQueryCmd(m.textInput.Value()))
}

func (m *Model) GetQueryCmd(query string) tea.Cmd {
	reqID := m.reqCnt
	m.reqCnt++

	slog.Debug("Submitting goto query request", "query", query, "id", reqID)

	return func() tea.Msg {
		results := m.filterResults(query)
		return NewUpdateMsg(query, results, reqID)
	}
}

func (msg UpdateMsg) Apply(m *Model) tea.Cmd {
	currentQuery := m.textInput.Value()
	if msg.query != currentQuery {
		slog.Debug("Ignoring stale goto query result",
			"msgQuery", msg.query,
			"currentQuery", currentQuery,
			"id", msg.reqID)
		return nil
	}

	m.results = msg.results
	m.cursor = 0
	m.renderIndex = 0

	return nil
}

func (m *Model) filterResults(query string) []string {
	parentDir := filepath.Dir(m.currentPath)
	if parentDir != m.currentPath {
		hasParent := true
		if query != "" {
			queryLower := strings.ToLower(query)
			hasParent = strings.Contains(strings.ToLower(".."), queryLower)
		}

		entries, err := os.ReadDir(m.currentPath)
		if err != nil {
			if hasParent {
				return []string{".."}
			}
			return []string{}
		}

		results := make([]string, 0, len(entries)+1)
		if hasParent {
			results = append(results, "..")
		}

		for _, entry := range entries {
			name := entry.Name()
			if query == "" || strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
				results = append(results, name)
			}
		}
		return results
	}

	entries, err := os.ReadDir(m.currentPath)
	if err != nil {
		return []string{}
	}

	queryLower := strings.ToLower(query)
	results := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if strings.Contains(strings.ToLower(name), queryLower) {
			results = append(results, name)
		}
	}
	return results
}

func (m *Model) updateResults() {
	m.results = m.filterResults(m.textInput.Value())
	m.cursor = 0
	m.renderIndex = 0
}

func (m *Model) navigateUp() {
	if len(m.results) == 0 {
		return
	}
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.renderIndex {
			m.renderIndex = m.cursor
		}
	} else {
		m.cursor = len(m.results) - 1
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	}
}

func (m *Model) navigateDown() {
	if len(m.results) == 0 {
		return
	}
	if m.cursor < len(m.results)-1 {
		m.cursor++
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	} else {
		m.cursor = 0
		m.renderIndex = 0
	}
}

func (m *Model) navigatePageUp() {
	if len(m.results) == 0 {
		return
	}
	scrollAmount := maxVisibleResults - 1
	if scrollAmount < 1 {
		scrollAmount = 1
	}

	if m.cursor-scrollAmount >= 0 {
		m.cursor -= scrollAmount
		if m.cursor < m.renderIndex {
			m.renderIndex = m.cursor
		}
	} else {
		m.cursor = 0
		m.renderIndex = 0
	}
}

func (m *Model) navigatePageDown() {
	if len(m.results) == 0 {
		return
	}
	scrollAmount := maxVisibleResults - 1
	if scrollAmount < 1 {
		scrollAmount = 1
	}

	if m.cursor+scrollAmount < len(m.results) {
		m.cursor += scrollAmount
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	} else {
		m.cursor = len(m.results) - 1
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	}
}
