package zoxide

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"

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
				slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
				m.Close()
			}
			return action, cmd
		}

		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			action = m.handleConfirm()
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			m.Close()
		case slices.Contains(common.Hotkeys.ListUp, msg.String()):
			m.navigateUp()
		case slices.Contains(common.Hotkeys.ListDown, msg.String()):
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

	// Update suggestions based on current input
	m.updateSuggestions()

	return cmd
}

func (m *Model) updateSuggestions() {
	if m.zClient == nil {
		return
	}

	query := m.textInput.Value()

	// Query zoxide with the current input (empty string shows all results)
	results, err := m.zClient.QueryAll(query)
	if err != nil {
		slog.Debug("Failed to get zoxide suggestions", "query", query, "error", err)
		m.results = []zoxidelib.Result{}
		m.cursor = 0
		m.renderIndex = 0
		return
	}

	// Don't limit results here - let scrolling handle display
	m.results = results

	// Reset selection when results change
	m.cursor = 0
	m.renderIndex = 0
}

func (m *Model) Render() string {
	r := ui.ZoxideRenderer(m.maxHeight, m.width)
	r.SetBorderTitle(m.headline)

	if m.zClient == nil {
		r.AddSection()
		r.AddLines(" Zoxide not available (check zoxide_support in config)")
		return r.Render()
	}

	r.AddLines(" " + m.textInput.View())
	r.AddSection()
	if len(m.results) > 0 {
		m.renderResultList(r)
	} else {
		r.AddLines(" No zoxide results found")
	}
	return r.Render()
}

func (m *Model) renderResultList(r *rendering.Renderer) {
	// Calculate visible range
	endIndex := m.renderIndex + maxVisibleResults
	if endIndex > len(m.results) {
		endIndex = len(m.results)
	}
	// Show visible results
	m.renderVisibleResults(r, endIndex)

	// Show scroll indicators if needed
	m.renderScrollIndicators(r, endIndex)
}

func (m *Model) renderVisibleResults(r *rendering.Renderer, endIndex int) {
	for i := m.renderIndex; i < endIndex; i++ {
		result := m.results[i]
		scoreTxt := fmt.Sprintf("%.1f", result.Score)

		// Truncate path if too long (account for score, separator, and padding)
		// Available width: modal width - borders(2) - padding(2) - score(5) - separator(3) = width - 12
		availablePathWidth := m.width - 12
		path := common.TruncateTextBeginning(result.Path, availablePathWidth, "...")

		line := fmt.Sprintf(" %5s | %s", scoreTxt, path)

		// Highlight the selected item
		if i == m.cursor {
			line = common.ModalCursorStyle.Render(line)
		}
		r.AddLines(line)
	}
}

func (m *Model) renderScrollIndicators(r *rendering.Renderer, endIndex int) {
	if len(m.results) <= maxVisibleResults {
		return
	}

	if m.renderIndex > 0 {
		r.AddSection()
		r.AddLines(" ↑ More results above")
	}
	if endIndex < len(m.results) {
		if m.renderIndex == 0 {
			r.AddSection()
		}
		r.AddLines(" ↓ More results below")
	}
}
