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

// Quiet period after query change before starting the walk.
const searchDebounceDelay = 120 * time.Millisecond

// Buffers progress snapshots between walk and pump.
const searchChannelSize = 4

// Opens recursive search on the focused panel and starts the first walk.
func (m *model) searchModeEnter() tea.Cmd {
	panel := m.getFocusedFilePanel()
	panel.EnterSearchMode()
	panel.SearchBar.SetWidth(m.fileModel.SinglePanelWidth - common.InnerPadding)
	// Keeps the opening key out of the query.
	m.firstTextInput = true
	return m.restartSearchCmd()
}

// Closes recursive search and restores the panel listing.
func (m *model) searchModeExit() {
	if m.searchCancel != nil {
		m.searchCancel()
		m.searchCancel = nil
	}
	m.searchReqID = m.nextIoReqCnt()
	m.getFocusedFilePanel().ExitSearchMode()
	m.fileModel.UpdateFilePanelsIfNeeded(true)
}

// Opens cursor result. Dirs become focused dir. Files focus in parent dir.
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

// Routes search keys. Single printables stay as query text. See isSearchTypingKey.
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

// Reports navigation keys so model.go can keep them out of the query input.
// Typing keys never match: searchKeyMatchesAction excludes them.
func searchNavKey(msg tea.KeyPressMsg) bool {
	return searchKeyMatchesAction(msg, common.Hotkeys.ListUp) ||
		searchKeyMatchesAction(msg, common.Hotkeys.ListDown) ||
		searchKeyMatchesAction(msg, common.Hotkeys.PageUp) ||
		searchKeyMatchesAction(msg, common.Hotkeys.PageDown)
}

// isSearchTypingKey reports query text. Single printable with no real modifier.
func isSearchTypingKey(msg tea.KeyPressMsg) bool {
	r := []rune(msg.Text)
	if len(r) != 1 || !unicode.IsPrint(r[0]) {
		return false
	}
	mods := msg.Mod &^ (tea.ModShift | tea.ModCapsLock)
	return mods == 0
}

// Reports binding match. Typing keys stay in the query. It also compares the
// Text-less form for terminals that report modifiers with Text set.
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

// Cancels running search and starts a debounced one for current query.
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

	// Clear stale results before the next walk, keeping the selection.
	panel.ResetSearchStateKeepingSelection()
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

// Forwards next progress snapshot. Returns nil when the channel closes.
func (m *model) searchPumpCmd(ch <-chan search.Progress, reqID int) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			return nil
		}
		return NewSearchProgressMsg(p, reqID)
	}
}

// Merges progress. Drops stale sessions and keeps pump alive.
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
