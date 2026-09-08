package internal

import (
	"context"
	"log/slog"
	"path/filepath"
	"time"
	"unicode"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/search"

	variable "github.com/yorukot/superfile/src/config"
)

// searchDebounceDelay is the quiet period after the last query change before
// the walk for the new query starts.
const searchDebounceDelay = 120 * time.Millisecond

// searchChannelSize bounds the progress snapshots buffered between the walk
// goroutine and the pump command.
const searchChannelSize = 4

// searchModeEnter starts a search session in the focused panel.
func (m *model) searchModeEnter() tea.Cmd {
	panel := m.getFocusedFilePanel()
	panel.EnterSearchMode()
	panel.SearchBar.SetWidth(m.fileModel.SinglePanelWidth - common.InnerPadding)
	// Consume the key that opened the session so it does not become the
	// first character of the query.
	m.firstTextInput = true
	return m.restartSearchCmd()
}

// searchModeExit ends the search session in the focused panel and cancels
// any running search.
func (m *model) searchModeExit() {
	if m.searchCancel != nil {
		m.searchCancel()
		m.searchCancel = nil
	}
	m.searchReqID = m.nextIoReqCnt()
	m.getFocusedFilePanel().ExitSearchMode()
	m.fileModel.UpdateFilePanelsIfNeeded(true)
}

// searchModeConfirm opens the result under the search cursor and leaves
// search mode. Opening a directory makes it the new focused directory;
// opening a file navigates to its containing directory and focuses it.
func (m *model) searchModeConfirm() {
	panel := m.getFocusedFilePanel()
	if panel.SearchBar.Value() == "" {
		m.searchModeExit()
		return
	}
	root := panel.Search.Root
	result := panel.GetSearchCursorResult()
	m.searchModeExit()
	if result == nil {
		return
	}
	targetPath := filepath.Join(root, result.Path)

	if !result.Dir && variable.ChooserFile != "" {
		if err := m.chooserFileWriteAndQuit(targetPath); err == nil {
			return
		} else {
			// Continue with navigation if the chooser file is not writable
			slog.Error("Error while writing to chooser file, continuing with navigation", "error", err)
		}
	}

	if result.Dir {
		err := m.updateCurrentFilePanelDir(targetPath)
		if err != nil {
			slog.Error("Error while opening search result directory", "error", err)
		}
		return
	}

	err := m.updateCurrentFilePanelDir(filepath.Dir(targetPath))
	if err != nil {
		slog.Error("Error while opening search result file", "error", err)
		return
	}
	m.getFocusedFilePanel().TargetFile = filepath.Base(targetPath)
}

// searchModeKey handles keys while a search session is active. Letter keys
// always reach the query input: navigation reuses the list_up/list_down/
// page_up/page_down bindings, but only when they cannot be query text (see
// isSearchTypingKey), so aliases like j/k keep typing instead of moving.
func (m *model) searchModeKey(msg tea.KeyPressMsg) tea.Cmd {
	panel := m.getFocusedFilePanel()
	switch {
	case searchKeyMatchesAction(msg, common.Hotkeys.CancelTyping):
		m.searchModeExit()
		return nil
	case searchKeyMatchesAction(msg, common.Hotkeys.ConfirmTyping):
		m.searchModeConfirm()
		return nil
	case searchKeyMatchesAction(msg, common.Hotkeys.SearchToggleHidden):
		return m.toggleDotFileController()
	case searchKeyMatchesAction(msg, common.Hotkeys.ListUp):
		panel.SearchListUp()
	case searchKeyMatchesAction(msg, common.Hotkeys.ListDown):
		panel.SearchListDown()
	case searchKeyMatchesAction(msg, common.Hotkeys.PageUp):
		panel.SearchPgUp()
	case searchKeyMatchesAction(msg, common.Hotkeys.PageDown):
		panel.SearchPgDown()
	}
	return nil
}

// isSearchTypingKey reports whether msg is query input rather than an
// action: a single printable rune with no modifier beyond shift/caps-lock
// (e.g. "k", "K"). Such keys must type even when they appear in a hotkey
// list; actions in search mode require a real modifier, a multi-rune key,
// or a non-printable key.
func isSearchTypingKey(msg tea.KeyPressMsg) bool {
	r := []rune(msg.Text)
	if len(r) != 1 || !unicode.IsPrint(r[0]) {
		return false
	}
	mods := msg.Mod &^ (tea.ModShift | tea.ModCapsLock)
	return mods == 0
}

// searchKeyMatchesAction reports whether msg triggers one of bindings while
// a search session is active. Typing keys (isSearchTypingKey) never match,
// so bare letters stay in the query. Empty "" padding slots are skipped.
// Besides the verbatim String() form, the Text-less keystroke form is
// compared so terminals that report modifier combos with Text set (e.g.
// ctrl+. with Text=".") still match.
func searchKeyMatchesAction(msg tea.KeyPressMsg, bindings []string) bool {
	if isSearchTypingKey(msg) {
		return false
	}
	for _, binding := range bindings {
		if binding == "" {
			continue
		}
		if msg.String() == binding {
			return true
		}
		if (tea.KeyPressMsg{Code: msg.Code, Mod: msg.Mod}).String() == binding {
			return true
		}
	}
	return false
}

// restartSearchCmd cancels the running search (if any) and starts a new one
// for the current query. It is called on every query change and on session
// entry. The walk itself is delayed by searchDebounceDelay so that fast
// typing does not restart it per keystroke.
func (m *model) restartSearchCmd() tea.Cmd {
	panel := m.getFocusedFilePanel()
	if !panel.Search.Active {
		return nil
	}
	if m.searchCancel != nil {
		m.searchCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.searchCancel = cancel
	reqID := m.nextIoReqCnt()
	m.searchReqID = reqID

	ch := make(chan search.Progress, searchChannelSize)
	m.searchChan = ch

	root := panel.Search.Root
	query := panel.SearchBar.Value()
	includeHidden := m.fileModel.DisplayDotFiles

	// Clear the previous session's results immediately: the new walk is
	// debounced, and confirming during the quiet period must not open a
	// result from the old query.
	panel.Search.Done = false
	panel.Search.Results = nil
	panel.Search.MatchCount = 0
	panel.Search.UnreadableDirs = 0
	panel.Search.Cursor = 0
	panel.Search.RenderIndex = 0
	go func() {
		defer close(ch)
		select {
		case <-time.After(searchDebounceDelay):
		case <-ctx.Done():
			return
		}
		search.Run(ctx, root, query, includeHidden, func(p search.Progress) {
			select {
			case ch <- p:
			case <-ctx.Done():
			}
		})
	}()

	return m.searchPumpCmd(ch, reqID)
}

// searchPumpCmd forwards the next progress snapshot from a search channel as
// a message. It returns no message when the channel closes, e.g. because the
// search was cancelled by a restart.
func (m *model) searchPumpCmd(ch <-chan search.Progress, reqID int) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return nil
		}
		return NewSearchProgressMsg(p, reqID)
	}
}

// applySearchProgress merges a progress snapshot into the focused panel and
// keeps the pump alive until the search completes. Stale messages from
// cancelled sessions are dropped.
func (m *model) applySearchProgress(msg SearchProgressMsg) tea.Cmd {
	if msg.GetReqID() != m.searchReqID {
		return nil
	}
	panel := m.getFocusedFilePanel()
	if !panel.Search.Active {
		return nil
	}
	panel.ApplySearchProgress(msg.progress)
	if msg.progress.Done {
		if m.searchCancel != nil {
			m.searchCancel()
			m.searchCancel = nil
		}
		return nil
	}
	return m.searchPumpCmd(m.searchChan, msg.GetReqID())
}
