package filepanel

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fvbommel/sortorder"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

func getOrderingFunc(elements []Element, reversed bool, sortKind sortmodel.SortKind) sliceOrderFunc {
	var order func(i, j int) bool
	switch sortKind {
	case sortmodel.SortByName:
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
	case sortmodel.SortBySize:
		order = getSizeOrderingFunc(elements, reversed)
	case sortmodel.SortByDate:
		order = func(i, j int) bool {
			return elements[i].Info.ModTime().After(elements[j].Info.ModTime()) != reversed
		}
	case sortmodel.SortByType:
		order = getTypeOrderingFunc(elements, reversed)
	case sortmodel.SortByNatural:
		order = func(i, j int) bool {
			// One of them is a directory, and other is not
			if elements[i].Directory != elements[j].Directory {
				return elements[i].Directory
			}
			if common.Config.CaseSensitiveSort {
				return sortorder.NaturalLess(elements[i].Name, elements[j].Name) != reversed
			}
			return sortorder.NaturalLess(
				strings.ToLower(elements[i].Name),
				strings.ToLower(elements[j].Name),
			) != reversed
		}
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

func sortFileElement(sortKind sortmodel.SortKind, reversed bool, dirEntries []os.DirEntry, location string) []Element {
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
			Directory: item.IsDir() || isSymlinkToDir(location, info, item.Name()),
			Location:  filepath.Join(location, item.Name()),
			Info:      info,
		})
	}

	sort.Slice(elements, getOrderingFunc(elements, reversed, sortKind))

	return elements
}

// Symlinks to directories are to be identified as directories
func isSymlinkToDir(location string, info os.FileInfo, name string) bool {
	if info.Mode()&os.ModeSymlink != 0 {
		targetInfo, errStat := os.Stat(filepath.Join(location, name))
		return errStat == nil && targetInfo.IsDir()
	}
	return false
}

func (m *Model) getPageScrollSize() int {
	scrollSize := common.Config.PageScrollSize
	if scrollSize <= 0 {
		// Use default full page behavior
		scrollSize = m.PanelElementHeight()
	}
	return scrollSize
}
