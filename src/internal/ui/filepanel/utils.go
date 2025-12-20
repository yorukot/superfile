package filepanel

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

func DefaultFilePanel(path string, focused bool) FilePanel {
	targetFile := ""
	panelPath := path
	// If path refers to a file, switch to its parent and remember the filename
	if stat, err := os.Stat(panelPath); err == nil && !stat.IsDir() {
		targetFile = filepath.Base(panelPath)
		panelPath = filepath.Dir(panelPath)
	}

	return FilePanel{
		RenderIndex: 0,
		Cursor:      0,
		Location:    panelPath,
		SortOptions: SortOptionsModel{
			//nolint:mnd // default sort options dimensions
			Width: 20,
			//nolint:mnd // default sort options dimensions
			Height: 4,
			Open:   false,
			Cursor: common.Config.DefaultSortType,
			Data: SortOptionsModelData{
				Options: []string{
					string(sortingName), string(sortingSize),
					string(sortingDateModified), string(sortingFileType),
				},
				Selected: common.Config.DefaultSortType,
				Reversed: common.Config.SortOrderReversed,
			},
		},
		PanelMode:        BrowserMode,
		IsFocused:        focused,
		DirectoryRecords: make(map[string]DirectoryRecord),
		SearchBar:        common.GenerateSearchBar(),
		TargetFile:       targetFile,
	}
}

func (p PanelMode) String() string {
	switch p {
	case SelectMode:
		return "selectMode"
	case BrowserMode:
		return "browserMode"
	default:
		return common.InvalidTypeString
	}
}

func getOrderingFunc(elements []Element, reversed bool, sortOption string) sliceOrderFunc {
	var order func(i, j int) bool
	switch sortOption {
	case string(sortingName):
		order = func(i, j int) bool {
			// One of them is a directory, and other is not
			if elements[i].Directory != elements[j].Directory {
				return elements[i].Directory
			}
			if common.Config.CaseSensitiveSort {
				return elements[i].Name < elements[j].Name != reversed
			}
			return strings.ToLower(elements[i].Name) < strings.ToLower(elements[j].Name) != reversed
		}
	case string(sortingSize):
		order = getSizeOrderingFunc(elements, reversed)
	case string(sortingDateModified):
		order = func(i, j int) bool {
			return elements[i].Info.ModTime().After(elements[j].Info.ModTime()) != reversed
		}
	case string(sortingFileType):
		order = getTypeOrderingFunc(elements, reversed)
	}
	return order
}

func getSizeOrderingFunc(elements []Element, reversed bool) sliceOrderFunc {
	return func(i, j int) bool {
		// Directories at the top sorted by direct child count (not recursive)
		// Files sorted by size

		// One of them is a directory, and other is not
		if elements[i].Directory != elements[j].Directory {
			return elements[i].Directory
		}

		// This needs to be improved, and we should sort by actual size only
		// Repeated recursive read would be slow, so we could cache
		if elements[i].Directory && elements[j].Directory {
			filesI, err := os.ReadDir(elements[i].Location)
			// No need of early return, we only call len() on filesI, so nil would
			// just result in 0
			if err != nil {
				slog.Error("Error when reading directory during sort", "error", err)
			}
			filesJ, err := os.ReadDir(elements[j].Location)
			if err != nil {
				slog.Error("Error when reading directory during sort", "error", err)
			}
			return len(filesI) < len(filesJ) != reversed
		}
		return elements[i].Info.Size() < elements[j].Info.Size() != reversed
	}
}

func getTypeOrderingFunc(elements []Element, reversed bool) sliceOrderFunc {
	return func(i, j int) bool {
		// One of them is a directory, and the other is not
		if elements[i].Directory != elements[j].Directory {
			return elements[i].Directory
		}

		var extI, extJ string
		if !elements[i].Directory {
			extI = strings.ToLower(filepath.Ext(elements[i].Name))
		}
		if !elements[j].Directory {
			extJ = strings.ToLower(filepath.Ext(elements[j].Name))
		}

		// Compare by extension/type
		if extI != extJ {
			return (extI < extJ) != reversed
		}

		// If same type, fall back to name
		if common.Config.CaseSensitiveSort {
			return (elements[i].Name < elements[j].Name) != reversed
		}

		return (strings.ToLower(elements[i].Name) < strings.ToLower(elements[j].Name)) != reversed
	}
}

func panelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - common.PanelPadding
}

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
func ReturnDirElement(location string, displayDotFile bool, sortOptions SortOptionsModelData) []Element {
	dirEntries, err := os.ReadDir(location)
	if err != nil {
		slog.Error("Error while returning folder elements", "error", err)
		return nil
	}

	dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
		// Entries not needed to be considered
		_, err := e.Info()
		return err != nil || (strings.HasPrefix(e.Name(), ".") && !displayDotFile)
	})

	// No files/directoes to process
	if len(dirEntries) == 0 {
		return nil
	}
	return SortFileElement(sortOptions, dirEntries, location)
}

func ReturnDirElementBySearchString(location string, displayDotFile bool, searchString string,
	sortOptions SortOptionsModelData,
) []Element {
	items, err := os.ReadDir(location)
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

	return SortFileElement(sortOptions, dirElements, location)
}

func SortFileElement(sortOptions SortOptionsModelData, dirEntries []os.DirEntry, location string) []Element {
	elements := make([]Element, 0, len(dirEntries))
	for _, item := range dirEntries {
		info, err := item.Info()
		if err != nil {
			slog.Error("Error while retrieving file info during sort",
				"error", err, "path", filepath.Join(location, item.Name()))
			continue
		}

		elements = append(elements, Element{
			Name:      item.Name(),
			Directory: item.IsDir() || IsSymlinkToDir(location, info, item.Name()),
			Location:  filepath.Join(location, item.Name()),
			Info:      info,
		})
	}

	sort.Slice(elements, getOrderingFunc(elements,
		sortOptions.Reversed, sortOptions.Options[sortOptions.Selected]))

	return elements
}

// Symlinks to directories are to be identified as directories
func IsSymlinkToDir(location string, info os.FileInfo, name string) bool {
	if info.Mode()&os.ModeSymlink != 0 {
		targetInfo, errStat := os.Stat(filepath.Join(location, name))
		return errStat == nil && targetInfo.IsDir()
	}
	return false
}

func removeElementByValue(slice []string, value string) []string {
	newSlice := []string{}
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func (panel *FilePanel) GetSelectedItem() Element {
	if panel.Cursor < 0 || len(panel.Element) <= panel.Cursor {
		return Element{}
	}
	return panel.Element[panel.Cursor]
}

func (panel *FilePanel) ResetSelected() {
	panel.Selected = panel.Selected[:0]
}

// For modification. Make sure to do a nil check
func (panel *FilePanel) GetSelectedItemPtr() *Element {
	if panel.Cursor < 0 || len(panel.Element) <= panel.Cursor {
		return nil
	}
	return &panel.Element[panel.Cursor]
}

// Note : This will soon be moved to its own package.
func (panel *FilePanel) ChangeFilePanelMode() {
	switch panel.PanelMode {
	case SelectMode:
		panel.Selected = panel.Selected[:0]
		panel.PanelMode = BrowserMode
	case BrowserMode:
		panel.PanelMode = SelectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", panel.PanelMode)
	}
}

// This should be the function that is always called whenever we are updating a directory.
func (panel *FilePanel) UpdateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", panel.Location, "path", path)
	// In case non Absolute path is passed, make sure to resolve it.
	path = utils.ResolveAbsPath(panel.Location, path)

	// Ignore if its the same directory. It prevents resetting of searchBar
	if path == panel.Location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	panel.DirectoryRecords[panel.Location] = DirectoryRecord{
		DirectoryCursor: panel.Cursor,
		DirectoryRender: panel.RenderIndex,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// Switch to "path"
	panel.Location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := panel.DirectoryRecords[panel.Location]
	if hasRecord {
		panel.Cursor = curDirectoryRecord.DirectoryCursor
		panel.RenderIndex = curDirectoryRecord.DirectoryRender
	} else {
		panel.Cursor = 0
		panel.RenderIndex = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", panel.Cursor, "render", panel.RenderIndex)

	// Reset the searchbar Value
	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	panel.SearchBar.SetValue("")

	return nil
}

func (panel *FilePanel) ParentDirectory() error {
	return panel.UpdateCurrentFilePanelDir("..")
}

func (panel *FilePanel) HandleResize(height int) {
	// Min render cursor that keeps the cursor in view
	minVisibleRenderCursor := panel.Cursor - panelElementHeight(height) + 1
	// Max render cursor. This ensures all elements are rendered if there is space
	maxRenderCursor := max(len(panel.Element)-panelElementHeight(height), 0)

	if panel.RenderIndex > maxRenderCursor {
		panel.RenderIndex = maxRenderCursor
	}
	if panel.RenderIndex < minVisibleRenderCursor {
		panel.RenderIndex = minVisibleRenderCursor
	}
}

// Select the item where cursor located (only work on select mode)
func (panel *FilePanel) SingleItemSelect() {
	if len(panel.Element) > 0 && panel.Cursor >= 0 && panel.Cursor < len(panel.Element) {
		elementLocation := panel.Element[panel.Cursor].Location

		if slices.Contains(panel.Selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			panel.Selected = removeElementByValue(panel.Selected, elementLocation)
		} else {
			panel.Selected = append(panel.Selected, elementLocation)
		}
	}
}
