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
	"github.com/yorukot/superfile/src/internal/utils"
)

func DefaultModel(maxHeight int, width int, cwd string) Model {
	return GenerateModel(cwd, maxHeight, width)
}

func GenerateModel(cwd string, maxHeight int, width int) Model {
	m := Model{
		headline:    icon.Search + icon.Space + gotoHeadlineText,
		open:        false,
		results:     []Result{},
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
		case slices.Contains(common.Hotkeys.ListUp, msg.String()) && msg.Type != tea.KeyRunes:
			m.navigateUp()
		case slices.Contains(common.Hotkeys.ListDown, msg.String()) && msg.Type != tea.KeyRunes:
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

		if selectedResult.Name == ".." {
			m.handleGoUp()
			return common.NoAction{}
		}

		fullSelectedPath := filepath.Join(m.currentPath, selectedResult.Name)

		if selectedResult.IsDir {
			return common.CDCurrentPanelAction{
				Location: fullSelectedPath,
			}
		}
		return common.OpenPanelAction{
			Location: fullSelectedPath,
		}
	}

	return common.NoAction{}
}

func (m *Model) handleTabCompletion() common.ModelAction { //nolint:unparam // Return value kept for future extensibility
	input := strings.TrimSpace(m.textInput.Value())

	if m.cursor >= 0 && m.cursor < len(m.results) {
		selected := m.results[m.cursor]

		if selected.Name == ".." {
			m.handleGoUp()
			return common.NoAction{}
		}

		fullSelectedPath := filepath.Join(m.currentPath, selected.Name)

		if selected.IsDir {
			m.currentPath = fullSelectedPath
			m.textInput.SetValue("")
			m.updateResults()
			return common.NoAction{}
		}
		m.textInput.SetValue(selected.Name)
		return common.NoAction{}
	}

	if input == ".." {
		m.handleGoUp()
		return common.NoAction{}
	}

	if input == "" {
		return common.NoAction{}
	}

	matched := filterResults(input, m.currentPath)
	if len(matched) == 1 {
		selected := matched[0]
		fullSelectedPath := filepath.Join(m.currentPath, selected.Name)

		if selected.IsDir {
			m.currentPath = fullSelectedPath
			m.textInput.SetValue("")
			m.updateResults()
			return common.NoAction{}
		}
		m.textInput.SetValue(selected.Name)
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
	path := m.currentPath
	m.reqCnt++

	slog.Debug("Submitting goto query request", "query", query, "path", path, "id", reqID)

	return func() tea.Msg {
		results := filterResults(query, path)
		return NewUpdateMsg(query, results, reqID, path)
	}
}

func (msg UpdateMsg) Apply(m *Model) tea.Cmd {
	currentQuery := m.textInput.Value()
	if msg.query != currentQuery || msg.path != m.currentPath {
		slog.Debug("Ignoring stale goto query result",
			"msgQuery", msg.query,
			"currentQuery", currentQuery,
			"msgPath", msg.path,
			"currentPath", m.currentPath,
			"msgReqID", msg.reqID,
			"currentReqID", m.reqCnt)
		return nil
	}

	m.results = msg.results
	m.cursor = 0
	m.renderIndex = 0

	return nil
}

type entryInfo struct {
	name  string
	isDir bool
}

func readDirEntries(path string) ([]entryInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	entryData := make([]entryInfo, len(entries))
	for i, entry := range entries {
		entryData[i].name = entry.Name()
		entryData[i].isDir = entry.IsDir()
	}
	return entryData, nil
}

func filterResults(query string, path string) []Result {
	entryData, err := readDirEntries(path)
	if err != nil {
		return handleReadDirError(path)
	}

	if query == "" {
		return buildAllResults(entryData, path)
	}

	return buildFilteredResults(query, entryData, path)
}

func handleReadDirError(path string) []Result {
	parentDir := filepath.Dir(path)
	if parentDir != path {
		return []Result{{Name: "..", IsDir: true}}
	}
	return []Result{}
}

func buildAllResults(entryData []entryInfo, path string) []Result {
	hasParent := filepath.Dir(path) != path

	if hasParent {
		results := make([]Result, 0, len(entryData)+1)
		results = append(results, Result{Name: "..", IsDir: true})
		for _, ed := range entryData {
			results = append(results, Result{Name: ed.name, IsDir: ed.isDir})
		}
		return results
	}

	results := make([]Result, len(entryData))
	for i, ed := range entryData {
		results[i] = Result{Name: ed.name, IsDir: ed.isDir}
	}
	return results
}

func buildFilteredResults(query string, entryData []entryInfo, path string) []Result {
	entryNames := make([]string, len(entryData))
	for i, ed := range entryData {
		entryNames[i] = ed.name
	}

	matches := utils.FzfSearch(query, entryNames)
	results := make([]Result, 0, len(matches))

	hasParent := filepath.Dir(path) != path
	if hasParent {
		queryLower := strings.ToLower(query)
		if strings.Contains(strings.ToLower(".."), queryLower) {
			results = append(results, Result{Name: "..", IsDir: true})
		}
	}

	for _, match := range matches {
		if match.HayIndex >= 0 && int(match.HayIndex) < len(entryData) {
			results = append(results, Result{
				Name:  entryData[match.HayIndex].name,
				IsDir: entryData[match.HayIndex].isDir,
			})
		}
	}
	return results
}

func (m *Model) updateResults() {
	m.results = filterResults(m.textInput.Value(), m.currentPath)
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
