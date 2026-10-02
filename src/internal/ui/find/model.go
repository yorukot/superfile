package find

import (
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func DefaultModel(maxHeight int, width int) Model {
	return GenerateModel(maxHeight, width)
}

func GenerateModel(maxHeight int, width int) Model {
	m := Model{
		headline:  icon.Search + icon.Space + findHeadlineText,
		open:      false,
		textInput: common.GeneratePromptTextInput(),
		results:   []FindResult{},
	}
	m.SetMaxHeight(maxHeight)
	m.SetWidth(width)
	m.textInput.Prompt = ""
	return m
}

func (m *Model) HandleUpdate(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	slog.Debug("find.Model HandleUpdate()", "msg", msg,
		"msgType", reflect.TypeOf(msg),
		"textInput", m.textInput.Value(),
		"cursorBlink", m.textInput.Styles().Cursor.Blink)
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !m.IsOpen() {
		slog.Error("HandleUpdate called on closed find")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// If find is not available, only allow confirm/cancel to close modal
		if !m.isAvailable() {
			switch {
			case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()),
				slices.Contains(common.Hotkeys.CancelTyping, msg.String()),
				slices.Contains(common.Hotkeys.Quit, msg.String()):
				m.Close()
			}
			return action, cmd
		}

		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			action = m.handleConfirm()
			m.Close()
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			m.Close()
		// We dont want keys like `j` and `k` to get stuck here
		// So if its a navigation key, lets specifically ignore
		// the alphanumeric keys as find panel is in text input
		// mode by default
		case slices.Contains(common.Hotkeys.ListUp, msg.String()) && !isKeyAlphaNum(msg):
			m.navigateUp()
		case slices.Contains(common.Hotkeys.ListDown, msg.String()) && !isKeyAlphaNum(msg):
			m.navigateDown()
		case slices.Contains(common.Hotkeys.FindFile, msg.String()) && m.justOpened:
			// Ignore the key that just opened this modal to prevent it from appearing in text input
			m.justOpened = false
		default:
			cmd = m.handleNormalKeyInput(msg)
		}
	default:
		// Non keypress updates like Cursor Blink
		// Only update text input if find is available
		if m.isAvailable() {
			m.textInput, cmd = m.textInput.Update(msg)
		}
	}
	return action, cmd
}

func (m *Model) isAvailable() bool {
	return common.Config.FindFileSupport && m.fdFound
}

func (m *Model) handleConfirm() common.ModelAction {
	// If we have results and a valid selection, navigate to selected result
	if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
		selectedResult := m.results[m.cursor]
		return common.GoToPathAction{
			Path:  selectedResult.Path,
			IsDir: selectedResult.IsDir,
		}
	}

	// No results or invalid selection - close modal
	return common.NoAction{}
}

func (m *Model) handleNormalKeyInput(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return tea.Batch(cmd, m.GetQueryCmd(m.textInput.Value()))
}

func (m *Model) GetQueryCmd(query string) tea.Cmd {
	if !m.isAvailable() {
		return nil
	}

	reqID := m.reqCnt
	m.reqCnt++

	slog.Debug("Submitting find query request", "query", query, "id", reqID)

	// Snapshot the search root at submit time so late completions can never
	// run against, or report results for, a different session's directory
	searchDir := m.searchDir

	return func() tea.Msg {
		retCode, output, err := m.runFindQuery(searchDir, query)
		if err != nil {
			slog.Debug("Find query failed", "query", query, "error", err, "id", reqID)
			return NewUpdateMsg(query, nil, err.Error(), reqID)
		}
		if errMsg := fdErrorMessage(retCode, output); errMsg != "" {
			return NewUpdateMsg(query, nil, errMsg, reqID)
		}
		return NewUpdateMsg(query, parseFindResults(searchDir, output), "", reqID)
	}
}

func (m *Model) runFindQuery(searchDir, query string) (int, string, error) {
	if query == "" {
		return utils.ExecuteCommand(common.DefaultCommandTimeout, searchDir,
			"fd", "--max-results", strconv.Itoa(maxResults))
	}
	// `--` ends fd's option parsing, so dash-prefixed query text is always
	// treated as a pattern and can never become an fd flag (e.g. `--exec`)
	return utils.ExecuteCommand(common.DefaultCommandTimeout, searchDir,
		"fd", "--max-results", strconv.Itoa(maxResults), "--", query)
}

func fdErrorMessage(retCode int, output string) string {
	if retCode == 0 {
		return ""
	}
	for _, line := range strings.Split(output, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func parseFindResults(searchDir string, output string) []FindResult {
	results := []FindResult{}
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		path := filepath.Join(searchDir, line)
		isDir := false
		if info, statErr := os.Stat(path); statErr == nil {
			isDir = info.IsDir()
		}
		results = append(results, FindResult{
			Path:  path,
			IsDir: isDir,
		})
	}
	return results
}

// Apply updates the find modal with query results
func (msg UpdateMsg) Apply(m *Model) tea.Cmd {
	// Ignore completions from previous sessions (before the current Open)
	if msg.reqID < m.openReqID {
		slog.Debug("Ignoring find query result from a previous session",
			"msgQuery", msg.query,
			"msgID", msg.reqID,
			"sessionID", m.openReqID)
		return nil
	}

	// Ignore stale results - only apply if query matches current input
	currentQuery := m.textInput.Value()
	if msg.query != currentQuery {
		slog.Debug("Ignoring stale find query result",
			"msgQuery", msg.query,
			"currentQuery", currentQuery,
			"id", msg.reqID)
		return nil
	}

	m.results = msg.results
	m.errMsg = msg.errMsg
	m.cursor = 0
	m.renderIndex = 0

	return nil
}
