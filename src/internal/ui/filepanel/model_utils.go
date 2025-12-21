package filepanel

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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

func (m *Model) GetSelectedItem() Element {
	if m.Cursor < 0 || len(m.Element) <= m.Cursor {
		return Element{}
	}
	return m.Element[m.Cursor]
}

func (m *Model) ResetSelected() {
	m.Selected = m.Selected[:0]
}

// For modification. Make sure to do a nil check
func (m *Model) GetSelectedItemPtr() *Element {
	if m.Cursor < 0 || len(m.Element) <= m.Cursor {
		return nil
	}
	return &m.Element[m.Cursor]
}

// Note : This will soon be moved to its own package.
func (m *Model) ChangeFilePanelMode() {
	switch m.PanelMode {
	case SelectMode:
		m.Selected = m.Selected[:0]
		m.PanelMode = BrowserMode
	case BrowserMode:
		m.PanelMode = SelectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", m.PanelMode)
	}
}

// This should be the function that is always called whenever we are updating a directory.
func (m *Model) UpdateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", m.Location, "path", path)
	// In case non Absolute path is passed, make sure to resolve it.
	path = utils.ResolveAbsPath(m.Location, path)

	// Ignore if its the same directory. It prevents resetting of searchBar
	if path == m.Location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	m.DirectoryRecords[m.Location] = directoryRecord{
		directoryCursor: m.Cursor,
		directoryRender: m.RenderIndex,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// Switch to "path"
	m.Location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := m.DirectoryRecords[m.Location]
	if hasRecord {
		m.Cursor = curDirectoryRecord.directoryCursor
		m.RenderIndex = curDirectoryRecord.directoryRender
	} else {
		m.Cursor = 0
		m.RenderIndex = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", m.Cursor, "render", m.RenderIndex)

	// Reset the searchbar Value
	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	m.SearchBar.SetValue("")

	return nil
}

func (m *Model) ParentDirectory() error {
	return m.UpdateCurrentFilePanelDir("..")
}

func (m *Model) HandleResize(height int) {
	m.scrollToCursor(m.Cursor, height)
}

// Select the item where cursor located (only work on select mode)
func (m *Model) SingleItemSelect() {
	if len(m.Element) > 0 && m.Cursor >= 0 && m.Cursor < len(m.Element) {
		elementLocation := m.Element[m.Cursor].Location

		if slices.Contains(m.Selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			m.Selected = removeElementByValue(m.Selected, elementLocation)
		} else {
			m.Selected = append(m.Selected, elementLocation)
		}
	}
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

	reRenderTime := int(float64(len(m.Element)) / common.ReRenderChunkDivisor)
	if m.IsFocused && !focusPanelReRender &&
		nowTime.Sub(m.LastTimeGetElement) < time.Duration(reRenderTime)*time.Second {
		return true
	}
	return false
}

func (m *Model) UpdateElementsIfNeeded(focusPanelReRender bool, toggleDotFile bool,
	updatedToggleDotFile bool, mainPanelHeight int) {
	nowTime := time.Now()
	if !m.shouldSkipPanelUpdate(focusPanelReRender, nowTime, updatedToggleDotFile) {
		// Load elements for this panel (with/without search filter)
		m.Element = m.getElements(toggleDotFile)
		// Update file panel list
		m.LastTimeGetElement = nowTime

		// For hover to file on first time loading
		if m.TargetFile != "" {
			m.applyTargetFileCursor(mainPanelHeight)
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

// FilePanelSlice creates a slice of FilePanels from the given paths
func FilePanelSlice(paths []string) []Model {
	res := make([]Model, len(paths))
	for i := range paths {
		// Making the first panel as the focussed
		isFocus := i == 0
		res[i] = defaultFilePanel(paths[i], isFocus)
	}
	return res
}

// defaultFilePanel creates a new FilePanel with default settings
func defaultFilePanel(path string, focused bool) Model {
	targetFile := ""
	panelPath := path
	// If path refers to a file, switch to its parent and remember the filename
	if stat, err := os.Stat(panelPath); err == nil && !stat.IsDir() {
		targetFile = filepath.Base(panelPath)
		panelPath = filepath.Dir(panelPath)
	}
	sortOptions := sortOptionsModel{
		//nolint:mnd // default sort options dimensions
		Width: 20,
		//nolint:mnd // default sort options dimensions
		Height: 4,
		Open:   false,
		Cursor: common.Config.DefaultSortType,
		Data: sortOptionsModelData{
			Options: []string{
				string(sortingName), string(sortingSize),
				string(sortingDateModified), string(sortingFileType),
			},
			Selected: common.Config.DefaultSortType,
			Reversed: common.Config.SortOrderReversed,
		},
	}
	return New(panelPath, sortOptions, focused, targetFile)
}

func New(location string, sortOptions sortOptionsModel, focused bool, targetFile string) Model {
	return Model{
		Cursor:           0,
		RenderIndex:      0,
		Location:         location,
		SortOptions:      sortOptions,
		PanelMode:        BrowserMode,
		IsFocused:        focused,
		DirectoryRecords: make(map[string]directoryRecord),
		SearchBar:        common.GenerateSearchBar(),
		TargetFile:       targetFile,
	}
}

func (m *Model) ElemCount() int {
	return len(m.Element)
}

func (m *Model) Empty() bool {
	return m.ElemCount() == 0
}
