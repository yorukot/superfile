package filepanel

import (
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/utils"
)

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
// getDirectoryElements returns the directory elements for the panel's current location
func (m *Model) getDirectoryElements(displayDotFiles bool) []Element {
	dirEntries, err := os.ReadDir(m.Location)
	if err != nil {
		slog.Error("Error while returning folder elements", "error", err)
		return nil
	}

	dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
		// Entries not needed to be considered
		_, err := e.Info()
		return err != nil || (strings.HasPrefix(e.Name(), ".") && !displayDotFiles)
	})

	// No files/directories to process
	if len(dirEntries) == 0 {
		return nil
	}
	return sortFileElement(m.SortOptions.Data, dirEntries, m.Location)
}

// getDirectoryElementsBySearch returns filtered directory elements based on search string
func (m *Model) getDirectoryElementsBySearch(displayDotFiles bool) []Element {
	searchString := m.SearchBar.Value()
	items, err := os.ReadDir(m.Location)
	if err != nil {
		slog.Error("Error while return folder element function", "error", err)
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	folderElementMap := map[string]os.DirEntry{}
	fileAndDirectories := []string{}

	for _, item := range items {
		fileInfo, err := item.Info()
		if err != nil {
			continue
		}
		if !displayDotFiles && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}

		fileAndDirectories = append(fileAndDirectories, item.Name())
		folderElementMap[item.Name()] = item
	}
	// https://github.com/reinhrst/fzf-lib/blob/main/core.go#L43
	// fzf returns matches ordered by score; we subsequently sort by the chosen sort option.
	fzfResults := utils.FzfSearch(searchString, fileAndDirectories)
	dirElements := make([]os.DirEntry, 0, len(fzfResults))
	for _, item := range fzfResults {
		resultItem := folderElementMap[item.Key]
		dirElements = append(dirElements, resultItem)
	}

	return sortFileElement(m.SortOptions.Data, dirElements, m.Location)
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

func (m *Model) UpdateElementsIfNeeded(force bool, displayDotFiles bool) {
	nowTime := time.Now()
	if force || !m.shouldSkipPanelUpdate(nowTime) {
		// Load elements for this panel (with/without search filter)
		m.element = m.getElements(displayDotFiles)
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

// Retrieves elements for a panel based on search bar value and sort options.
func (m *Model) getElements(displayDotFiles bool) []Element {
	if m.SearchBar.Value() != "" {
		return m.getDirectoryElementsBySearch(displayDotFiles)
	}
	return m.getDirectoryElements(displayDotFiles)
}
