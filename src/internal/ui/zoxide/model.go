package zoxide

import (
	"log/slog"
	"reflect"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func DefaultModel(maxHeight int, width int, zClient *zoxidelib.Client) Model {
	return GenerateModel(zClient, maxHeight, width)
}

func GenerateModel(zClient *zoxidelib.Client, maxHeight int, width int) Model {
	m := Model{
		headline:  icon.Search + icon.Space + zoxideHeadlineText,
		open:      false,
		textInput: common.GeneratePromptTextInput(),
		zClient:   zClient,
		results:   []zoxidelib.Result{},
	}
	m.SetMaxHeight(maxHeight)
	m.SetWidth(width)
	m.textInput.Prompt = ""
	return m
}

func (m *Model) HandleUpdate(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	slog.Debug("zoxide.Model HandleUpdate()", "msg", msg,
		"msgType", reflect.TypeOf(msg),
		"textInput", m.textInput.Value(),
		"cursorBlink", m.textInput.Cursor.Blink)
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !m.IsOpen() {
		slog.Error("HandleUpdate called on closed zoxide")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If zoxide is not available, only allow confirm/cancel to close modal
		if m.zClient == nil {
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
		// the alphanumeric keys as zoxide panel is in text input
		// mode by default
		case slices.Contains(common.Hotkeys.ListUp, msg.String()) && !isKeyAlphaNum(msg):
			m.navigateUp()
		case slices.Contains(common.Hotkeys.ListDown, msg.String()) && !isKeyAlphaNum(msg):
			m.navigateDown()
		case slices.Contains(common.Hotkeys.OpenZoxide, msg.String()) && m.justOpened:
			// Ignore the 'z' key that just opened this modal to prevent it from appearing in text input
			m.justOpened = false
		default:
			cmd = m.handleNormalKeyInput(msg)
		}
	default:
		// Non keypress updates like Cursor Blink
		// Only update text input if zoxide is available
		if m.zClient != nil {
			m.textInput, cmd = m.textInput.Update(msg)
		}
	}
	return action, cmd
}

func (m *Model) handleConfirm() common.ModelAction {
	// If we have results and a valid selection, navigate to selected result
	if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
		selectedResult := m.results[m.cursor]
		return common.CDCurrentPanelAction{
			Location: selectedResult.Path,
		}
	}

	// No results or invalid selection - close modal
	return common.NoAction{}
}

func (m *Model) handleNormalKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return tea.Batch(cmd, m.GetQueryCmd(m.textInput.Value()))
}

func (m *Model) GetQueryCmd(query string) tea.Cmd {
	if m.zClient == nil || !common.Config.ZoxideSupport {
		return nil
	}

	reqID := m.reqCnt
	m.reqCnt++

	slog.Debug("Submitting zoxide query request", "query", query, "id", reqID)

	return func() tea.Msg {
		queryFields := strings.Fields(query)
		results, err := m.zClient.QueryAll(queryFields...)
		if err != nil {
			slog.Debug("Zoxide query failed", "query", query, "error", err, "id", reqID)
			return NewUpdateMsg(query, []zoxidelib.Result{}, reqID)
		}
		return NewUpdateMsg(query, results, reqID)
	}
}

// Apply updates the zoxide modal with query results
func (msg UpdateMsg) Apply(m *Model) tea.Cmd {
	// Ignore stale results - only apply if query matches current input
	currentQuery := m.textInput.Value()
	if msg.query != currentQuery {
		slog.Debug("Ignoring stale zoxide query result",
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
