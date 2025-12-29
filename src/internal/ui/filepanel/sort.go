package filepanel

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

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
			if common.Config.NaturalSort {
				cmp := naturalCompare(elements[i].Name, elements[j].Name, common.Config.CaseSensitiveSort)
				return (cmp < 0) != reversed
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
		if common.Config.NaturalSort {
			cmp := naturalCompare(elements[i].Name, elements[j].Name, common.Config.CaseSensitiveSort)
			return (cmp < 0) != reversed
		}
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

// naturalCompare compares two strings using natural sorting.
// Returns negative if a < b, 0 if a == b, positive if a > b.
// Natural sorting treats numeric sequences as numbers (e.g., "file2" < "file10").
func naturalCompare(a, b string, caseSensitive bool) int {
	if !caseSensitive {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}

	i, j := 0, 0
	for i < len(a) && j < len(b) {
		aIsDigit := unicode.IsDigit(rune(a[i]))
		bIsDigit := unicode.IsDigit(rune(b[j]))

		if aIsDigit && bIsDigit {
			// Extract numeric sequences
			aNum, aEnd := extractNumber(a, i)
			bNum, bEnd := extractNumber(b, j)

			// Compare numeric values
			if aNum != bNum {
				if aNum < bNum {
					return -1
				}
				return 1
			}

			// If numeric values are equal, move to next segment
			i = aEnd
			j = bEnd
		} else {
			// Compare characters lexicographically
			if a[i] != b[j] {
				if a[i] < b[j] {
					return -1
				}
				return 1
			}
			i++
			j++
		}
	}

	// Handle remaining characters
	if i < len(a) {
		return 1
	}
	if j < len(b) {
		return -1
	}
	return 0
}

// extractNumber extracts a numeric sequence from string s starting at position start.
// Returns the numeric value and the position after the last digit.
func extractNumber(s string, start int) (uint64, int) {
	end := start
	for end < len(s) && unicode.IsDigit(rune(s[end])) {
		end++
	}

	// Parse the number (handle leading zeros by just parsing the value)
	var num uint64
	for i := start; i < end; i++ {
		num = num*10 + uint64(s[i]-'0')
	}
	return num, end
}
