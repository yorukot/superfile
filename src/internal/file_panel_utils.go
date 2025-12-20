package internal

import (
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
					string(SortingName), string(SortingSize),
					string(SortingDateModified), string(SortingFileType),
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
		return invalidTypeString
	}
}

func GetOrderingFunc(elements []Element, reversed bool, sortOption string) SliceOrderFunc {
	var order func(i, j int) bool
	switch sortOption {
	case string(SortingName):
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
	case string(SortingSize):
		order = GetSizeOrderingFunc(elements, reversed)
	case string(SortingDateModified):
		order = func(i, j int) bool {
			return elements[i].Info.ModTime().After(elements[j].Info.ModTime()) != reversed
		}
	case string(SortingFileType):
		order = GetTypeOrderingFunc(elements, reversed)
	}
	return order
}

func GetSizeOrderingFunc(elements []Element, reversed bool) SliceOrderFunc {
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

func GetTypeOrderingFunc(elements []Element, reversed bool) SliceOrderFunc {
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

func PanelElementHeight(mainPanelHeight int) int {
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

	sort.Slice(elements, GetOrderingFunc(elements,
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
