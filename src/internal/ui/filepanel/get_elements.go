package filepanel

import (
	"context"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/pkg/utils"
)

const panelIOTimeout = 10 * time.Second

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
// getDirectoryElements returns the directory elements for the panel's current location
func (m *Model) listDirectoryEntries(ctx context.Context) ([]filesystem.Entry, error) {
	session, err := m.paneSession()
	if err != nil {
		return nil, err
	}
	return session.List(ctx, m.CurrentLocation().Path)
}

func (m *Model) getDirectoryElements(displayDotFile bool) ([]Element, error) {
	ctx, cancel := context.WithTimeout(context.Background(), panelIOTimeout)
	defer cancel()
	return m.getDirectoryElementsWithContext(ctx, displayDotFile)
}

func (m *Model) getDirectoryElementsWithContext(ctx context.Context, displayDotFile bool) ([]Element, error) {
	dirEntries, err := m.listDirectoryEntries(ctx)
	if err != nil {
		return nil, err
	}
	return m.ElementsFromDirectoryEntries(ctx, dirEntries, displayDotFile)
}

func (m *Model) ElementsFromDirectoryEntries(
	ctx context.Context,
	dirEntries []filesystem.Entry,
	displayDotFile bool,
) ([]Element, error) {
	dirEntries = slices.Clone(dirEntries)
	dirEntries = slices.DeleteFunc(dirEntries, func(e filesystem.Entry) bool {
		// Entries not needed to be considered
		return strings.HasPrefix(e.Name, ".") && !displayDotFile
	})

	// No files/directories to process
	if len(dirEntries) == 0 {
		return nil, nil
	}
	session, err := m.paneSession()
	if err != nil {
		return nil, err
	}
	return sortFileElement(ctx, m.SortKind, m.SortReversed, dirEntries, m.CurrentLocation(), session), nil
}

// getDirectoryElementsBySearch returns filtered directory elements based on search string
func (m *Model) getDirectoryElementsBySearch(displayDotFile bool) ([]Element, error) {
	ctx, cancel := context.WithTimeout(context.Background(), panelIOTimeout)
	defer cancel()
	return m.getDirectoryElementsBySearchWithContext(ctx, displayDotFile)
}

func (m *Model) getDirectoryElementsBySearchWithContext(ctx context.Context, displayDotFile bool) ([]Element, error) {
	searchString := m.SearchBar.Value()
	items, err := m.listDirectoryEntries(ctx)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	folderElementMap := map[string]filesystem.Entry{}
	fileAndDirectories := []string{}

	for _, item := range items {
		if !displayDotFile && strings.HasPrefix(item.Name, ".") {
			continue
		}

		fileAndDirectories = append(fileAndDirectories, item.Name)
		folderElementMap[item.Name] = item
	}
	// https://github.com/reinhrst/fzf-lib/blob/main/core.go#L43
	// fzf returns matches ordered by score; we subsequently sort by the chosen sort option.
	fzfResults := utils.FzfSearch(searchString, fileAndDirectories)
	dirElements := make([]filesystem.Entry, 0, len(fzfResults))
	for _, item := range fzfResults {
		resultItem := folderElementMap[item.Key]
		dirElements = append(dirElements, resultItem)
	}

	session, err := m.paneSession()
	if err != nil {
		return nil, err
	}
	return sortFileElement(ctx, m.SortKind, m.SortReversed, dirElements, m.CurrentLocation(), session), nil
}

// Helper to decide whether to skip updating a panel this tick.
func (m *Model) shouldSkipPanelUpdate(nowTime time.Time) bool {
	if !m.IsFocused {
		return nowTime.Sub(m.LastTimeGetElement) < nonFocussedPanelReRenderTime
	}

	reRenderTime := int(float64(m.ElemCount()) / ReRenderChunkDivisor)
	reRenderTime = min(reRenderTime, ReRenderMaxDelay)
	return !m.NeedsReRender() &&
		nowTime.Sub(m.LastTimeGetElement) < time.Duration(reRenderTime)*time.Second
}

func (m *Model) UpdateElementsIfNeeded(force bool, displayDotFile bool) {
	nowTime := time.Now()
	if force || !m.shouldSkipPanelUpdate(nowTime) {
		// Load elements for this panel (with/without search filter)
		elements, err := m.getElements(displayDotFile)
		if err != nil {
			slog.Error("Error while loading folder elements", "error", err, "location", m.DisplayLocation())
			return
		}
		m.element = elements
		// Update file panel list
		m.LastTimeGetElement = nowTime

		// For hover to file on first time loading
		if m.TargetFile != "" {
			m.applyTargetFileCursor()
		}

		// If cursor becomes invalid due to element update, reset
		if m.ValidateCursorAndRenderIndex() != nil {
			m.scrollToCursor(0)
		}
	}
}

func (m *Model) ShouldUpdateElements(force bool, now time.Time) bool {
	if m.elementsLoading {
		if force {
			m.elementsRefreshPending = true
		}
		return false
	}
	if force {
		return true
	}
	if m.CurrentLocation().Provider != filesystem.ProviderLocal {
		refreshInterval := remoteFocusedPanelRefreshTime
		if !m.IsFocused {
			refreshInterval = nonFocussedPanelReRenderTime
		}
		return m.LastTimeGetElement.IsZero() || now.Sub(m.LastTimeGetElement) >= refreshInterval
	}
	return !m.shouldSkipPanelUpdate(now)
}

func (m *Model) BeginElementsLoading(force bool, now time.Time) (uint64, bool) {
	if !m.ShouldUpdateElements(force, now) {
		return 0, false
	}
	m.elementsRequestID++
	m.elementsLoading = true
	m.elementsRefreshPending = false
	m.elementsContext, m.elementsCancel = context.WithTimeout(context.Background(), panelIOTimeout)
	m.LastTimeGetElement = now
	return m.elementsRequestID, true
}

func (m *Model) FinishElementsLoading(requestID uint64) (bool, bool) {
	if !m.elementsLoading || requestID != m.elementsRequestID {
		return false, false
	}
	if m.elementsCancel != nil {
		m.elementsCancel()
	}
	m.elementsContext = nil
	m.elementsCancel = nil
	m.elementsLoading = false
	pending := m.elementsRefreshPending
	m.elementsRefreshPending = false
	return true, pending
}

func (m *Model) InvalidateElementsLoading() {
	if m.elementsCancel != nil {
		m.elementsCancel()
	}
	m.elementsContext = nil
	m.elementsCancel = nil
	m.LastTimeGetElement = time.Time{}
	m.elementsLoading = false
	m.elementsRefreshPending = false
	m.elementsRequestID++
}

func (m *Model) LoadElements(displayDotFile bool) ([]Element, error) {
	if m.elementsContext != nil {
		return m.getElementsWithContext(m.elementsContext, displayDotFile)
	}
	return m.getElements(displayDotFile)
}

func (m *Model) ApplyLoadedElements(elements []Element, loadedAt time.Time) {
	m.element = elements
	m.LastTimeGetElement = loadedAt
	if m.TargetFile != "" {
		m.applyTargetFileCursor()
	}
	if m.ValidateCursorAndRenderIndex() != nil {
		m.scrollToCursor(0)
	}
}

// Retrieves elements for a panel based on search bar value and sort options.
func (m *Model) getElements(displayDotFile bool) ([]Element, error) {
	ctx, cancel := context.WithTimeout(context.Background(), panelIOTimeout)
	defer cancel()
	return m.getElementsWithContext(ctx, displayDotFile)
}

func (m *Model) getElementsWithContext(ctx context.Context, displayDotFile bool) ([]Element, error) {
	if m.SearchBar.Value() != "" {
		return m.getDirectoryElementsBySearchWithContext(ctx, displayDotFile)
	}
	return m.getDirectoryElementsWithContext(ctx, displayDotFile)
}
