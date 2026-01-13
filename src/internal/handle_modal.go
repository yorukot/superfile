package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/utils"
)

// Cancel typing modal e.g. create file or directory
func (m *model) cancelTypingModal() {
	m.typingModal.textInput.Blur()
	m.typingModal.open = false
}

// Confirm to create file or directory
func (m *model) createItem() {
	if err := checkFileNameValidity(m.typingModal.textInput.Value()); err != nil {
		m.typingModal.errorMesssage = err.Error()
		slog.Error("Errow while createItem during item creation", "error", err)

		return
	}

	defer func() {
		m.typingModal.errorMesssage = ""
		m.typingModal.open = false
		m.typingModal.textInput.Blur()
	}()

	path := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())
	if !strings.HasSuffix(m.typingModal.textInput.Value(), string(filepath.Separator)) {
		path, _ = renameIfDuplicate(path)
		if err := os.MkdirAll(filepath.Dir(path), utils.UserDirPerm); err != nil {
			slog.Error("Error while createItem during directory creation", "error", err)
			return
		}
		f, err := os.Create(path)
		if err != nil {
			slog.Error("Error while createItem during file creation", "error", err)
			return
		}
		defer f.Close()
	} else {
		err := os.MkdirAll(path, utils.UserDirPerm)
		if err != nil {
			slog.Error("Error while createItem during directory creation", "error", err)
			return
		}
	}
}

// Cancel rename file or directory
func (m *model) cancelRename() {
	panel := m.getFocusedFilePanel()
	panel.Rename.Blur()
	panel.Renaming = false
	m.fileModel.Renaming = false
}

// Connfirm rename file or directory
func (m *model) confirmRename() {
	panel := m.getFocusedFilePanel()

	// Although we dont expect this to happen based on our current flow
	// Just adding it here to be safe
	if panel.Empty() {
		slog.Error("confirmRename called on empty panel")
		return
	}

	oldPath := panel.GetFocusedItem().Location
	newPath := filepath.Join(panel.Location, panel.Rename.Value())

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		slog.Error("Error while confirmRename during rename", "error", err)
		// Dont return. We have to also reset the panel and model information
	}

	m.clipboard.UpdatePath(oldPath, newPath)
	m.fileModel.Renaming = false
	panel.Rename.Blur()
	panel.Renaming = false
}

func (m *model) openSortOptionsMenu() {
	panel := m.getFocusedFilePanel()
	panel.SortOptions.Open = true
}

func (m *model) cancelSortOptions() {
	panel := m.getFocusedFilePanel()
	panel.SortOptions.Cursor = panel.SortOptions.Data.Selected
	panel.SortOptions.Open = false
}

func (m *model) confirmSortOptions() {
	panel := m.getFocusedFilePanel()
	panel.SortOptions.Data.Selected = panel.SortOptions.Cursor
	panel.SortOptions.Open = false
}

// Move the cursor up in the sort options menu
func (m *model) sortOptionsListUp() {
	panel := m.getFocusedFilePanel()
	if panel.SortOptions.Cursor > 0 {
		panel.SortOptions.Cursor--
	} else {
		panel.SortOptions.Cursor = len(panel.SortOptions.Data.Options) - 1
	}
}

// Move the cursor down in the sort options menu
func (m *model) sortOptionsListDown() {
	panel := m.getFocusedFilePanel()
	if panel.SortOptions.Cursor < len(panel.SortOptions.Data.Options)-1 {
		panel.SortOptions.Cursor++
	} else {
		panel.SortOptions.Cursor = 0
	}
}

func (m *model) toggleReverseSort() {
	panel := m.getFocusedFilePanel()
	panel.SortOptions.Data.Reversed = !panel.SortOptions.Data.Reversed
}

// Cancel search, this will clear all searchbar input
func (m *model) cancelSearch() {
	panel := m.getFocusedFilePanel()
	panel.SearchBar.Blur()
	panel.SearchBar.SetValue("")
}

// Confirm search. This will exit the search bar and filter the files
func (m *model) confirmSearch() {
	panel := m.getFocusedFilePanel()
	panel.SearchBar.Blur()
}

// Help menu panel list up
func (m *model) helpMenuListUp() {
	if m.helpMenu.cursor > 1 {
		m.helpMenu.cursor--
		if m.helpMenu.cursor < m.helpMenu.renderIndex {
			m.helpMenu.renderIndex = m.helpMenu.cursor
		}
		if m.helpMenu.filteredData[m.helpMenu.cursor].subTitle != "" {
			m.helpMenu.cursor--
		}
	} else {
		// Set the cursor to the last item in the list.
		// We use max(..., 0) as a safeguard to prevent a negative cursor index
		// in case the filtered list is empty.
		m.helpMenu.cursor = max(len(m.helpMenu.filteredData)-1, 0)

		// Adjust the render index to show the bottom of the list.
		// Similarly, we use max(..., 0) to ensure the renderIndex doesn't become negative,
		// which can happen if the number of items is less than the view height.
		// This prevents a potential out-of-bounds panic during rendering.
		m.helpMenu.renderIndex = max(len(m.helpMenu.filteredData)-(m.helpMenu.height-common.InnerPadding), 0)
	}
}

// Help menu panel list down
func (m *model) helpMenuListDown() {
	if len(m.helpMenu.filteredData) == 0 {
		return
	}

	if m.helpMenu.cursor < len(m.helpMenu.filteredData)-1 {
		// Compute the next selectable row (skip subtitles).
		next := m.helpMenu.cursor + 1
		for next < len(m.helpMenu.filteredData) && m.helpMenu.filteredData[next].subTitle != "" {
			next++
		}
		if next >= len(m.helpMenu.filteredData) {
			// Wrap if no more selectable rows.
			m.helpMenu.cursor = 1
			m.helpMenu.renderIndex = 0
			return
		}
		m.helpMenu.cursor = next

		// Scroll down if cursor moved past the viewport.
		if m.helpMenu.cursor > m.helpMenu.renderIndex+m.helpMenu.height-5 {
			m.helpMenu.renderIndex++
		}
		// Clamp renderIndex to bottom.
		bottom := len(m.helpMenu.filteredData) - (m.helpMenu.height - common.InnerPadding)
		if bottom < 0 {
			bottom = 0
		}
		if m.helpMenu.renderIndex > bottom {
			m.helpMenu.renderIndex = bottom
		}
	} else {
		m.helpMenu.cursor = 1
		m.helpMenu.renderIndex = 0
	}
}

func removeOrphanSections(items []helpMenuModalData) []helpMenuModalData {
	var result []helpMenuModalData
	// Since we can't know beforehand which section are we actually filtering
	// we may end up in a scenario where there are two sections (General, Panel navigation)
	// with no hotkeys between them, so we need to remove the section which its hotkeys was
	// completely filtered out (Orphan sections)
	for i := range items {
		if items[i].subTitle != "" {
			// Look ahead: is the next item a real hotkey?
			if i+1 < len(items) && items[i+1].subTitle == "" {
				result = append(result, items[i])
			}
			// Else: skip this subtitle because no children
		} else {
			result = append(result, items[i])
		}
	}
	return result
}

func (m *model) filterHelpMenu(query string) {
	filtered := fuzzySearch(query, m.helpMenu.data)
	filtered = removeOrphanSections(filtered)

	m.helpMenu.filteredData = filtered
	if len(filtered) == 0 {
		m.helpMenu.cursor = 0
	} else {
		m.helpMenu.cursor = 1
	}
	m.helpMenu.renderIndex = 0
}

// Fuzzy search function for a list of helpMenuModalData.
// inspired from: sidebar/directory_utils.go
func fuzzySearch(query string, data []helpMenuModalData) []helpMenuModalData {
	if len(data) == 0 {
		return []helpMenuModalData{}
	}

	// Optimization - This haystack can be kept precomputed based on description
	// instead of re computing it in each call
	haystack := []string{}
	idxMap := []int{}
	for i, item := range data {
		if item.subTitle != "" {
			continue
		}
		searchText := strings.Join(item.hotkey, " ") + " " + item.description
		haystack = append(haystack, searchText)
		idxMap = append(idxMap, i)
	}

	matchedIdx := map[int]struct{}{}
	for _, match := range utils.FzfSearch(query, haystack) {
		matchedIdx[idxMap[match.HayIndex]] = struct{}{}
	}

	results := []helpMenuModalData{}
	for i, d := range data {
		_, isMatch := matchedIdx[i]
		if d.subTitle != "" || isMatch {
			results = append(results, d)
		}
	}

	return results
}

// Toggle help menu
func (m *model) openHelpMenu() {
	if m.helpMenu.open {
		m.helpMenu.searchBar.Reset()
		m.helpMenu.open = false
		return
	}

	// Reset filteredData to the full data whenever the helpMenu is opened
	m.helpMenu.filteredData = m.helpMenu.data
	m.helpMenu.open = true
}

// Quit help menu
func (m *model) quitHelpMenu() {
	m.helpMenu.searchBar.Reset()
	m.helpMenu.open = false
}

func (m *model) getFocusedFilePanel() *filepanel.Model {
	return m.fileModel.GetFocusedFilePanel()
}
