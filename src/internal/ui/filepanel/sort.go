package filepanel

import (
	"context"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fvbommel/sortorder"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

func getOrderingFunc(
	elements []Element,
	reversed bool,
	sortKind sortmodel.SortKind,
	directoryEntryCount func(Element) int,
) sliceOrderFunc {
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
		order = getSizeOrderingFunc(elements, reversed, directoryEntryCount)
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

func getSizeOrderingFunc(elements []Element, reversed bool, directoryEntryCount func(Element) int) sliceOrderFunc {
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
			return directoryEntryCount(elements[i]) < directoryEntryCount(elements[j]) != reversed
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

func sortFileElement(
	ctx context.Context,
	sortKind sortmodel.SortKind,
	reversed bool,
	dirEntries []filesystem.Entry,
	location filesystem.Location,
	session filesystem.Session,
) []Element {
	elements := make([]Element, 0, len(dirEntries))
	directoryCounts := map[string]int{}
	for _, item := range dirEntries {
		entryPath := item.Path
		if entryPath.String() == "" {
			entryPath = joinPath(location.Path, item.Name)
		}
		isDirectory := item.Stat.IsDir
		if !isDirectory && item.Stat.IsSymlink && item.Stat.Target.String() != "" {
			targetStat, err := session.Stat(ctx, resolveSymlinkTargetPath(entryPath, item.Stat.Target))
			if err == nil {
				isDirectory = targetStat.IsDir
			}
		}

		elements = append(elements, Element{
			Name:      item.Name,
			Directory: isDirectory,
			Location:  entryPath.String(),
			Path:      entryPath,
			Info:      item.Stat.AsFileInfo(),
		})
	}

	sort.Slice(elements, getOrderingFunc(elements, reversed, sortKind, func(element Element) int {
		if count, ok := directoryCounts[element.Location]; ok {
			return count
		}
		files, err := session.List(ctx, elementPath(element, location))
		if err != nil {
			slog.Error("Error when reading directory during sort", "error", err, "path", element.Location)
			return 0
		}
		directoryCounts[element.Location] = len(files)
		return len(files)
	}))

	return elements
}

func (m *Model) getPageScrollSize() int {
	scrollSize := common.Config.PageScrollSize
	if scrollSize <= 0 {
		// Use default full page behavior
		scrollSize = m.PanelElementHeight()
	}
	return scrollSize
}
