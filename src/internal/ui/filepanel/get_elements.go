package filepanel

import (
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
// getDirectoryElements returns the directory elements for the panel's current location
func (m *Model) getDirectoryElements(displayDotFile bool) []Element {
	dirEntries, err := os.ReadDir(m.Location)
	if err != nil {
		slog.Error("Error while returning folder elements", "error", err)
		return nil
	}

	dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
		// Entries not needed to be considered
		_, err := e.Info()
		return err != nil || (strings.HasPrefix(e.Name(), ".") && !displayDotFile)
	})

	// No files/directories to process
	if len(dirEntries) == 0 {
		return nil
	}
	return sortFileElement(m.SortOptions.Data, dirEntries, m.Location)
}

// getDirectoryElementsBySearch returns filtered directory elements based on search string
func (m *Model) getDirectoryElementsBySearch(displayDotFile bool) []Element {
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
		if !displayDotFile && strings.HasPrefix(fileInfo.Name(), ".") {
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
func (m *Model) shouldSkipPanelUpdate(focusPanelReRender bool,
	nowTime time.Time, updatedModelToggleDotFile bool) bool {
	// Throttle non-focused panels unless dotfile toggle changed
	if !m.IsFocused && nowTime.Sub(m.LastTimeGetElement) < 3*time.Second {
		if !updatedModelToggleDotFile {
			return true
		}
	}

	reRenderTime := int(float64(len(m.element)) / common.ReRenderChunkDivisor)
	if m.IsFocused && !focusPanelReRender &&
		nowTime.Sub(m.LastTimeGetElement) < time.Duration(reRenderTime)*time.Second {
		return true
	}
	return false
}

func (m *Model) UpdateElementsIfNeeded(focusPanelReRender bool, toggleDotFile bool, updatedToggleDotFile bool) {
	nowTime := time.Now()
	if !m.shouldSkipPanelUpdate(focusPanelReRender, nowTime, updatedToggleDotFile) {
		// Load elements for this panel (with/without search filter)
		m.element = m.getElements(toggleDotFile)
		// Update file panel list
		m.LastTimeGetElement = nowTime

		// For hover to file on first time loading
		if m.TargetFile != "" {
			m.applyTargetFileCursor()
		}
	}
}

// Retrieves elements for a panel based on search bar value and sort options.
func (m *Model) getElements(toggleDotFile bool) []Element {
	if m.SearchBar.Value() != "" {
		return m.getDirectoryElementsBySearch(toggleDotFile)
	}
	return m.getDirectoryElements(toggleDotFile)
}
