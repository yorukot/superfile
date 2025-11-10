package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/yorukot/superfile/src/internal/common"
)

const (
	sortingName         sortingKind = "Name"
	sortingSize         sortingKind = "Size"
	sortingDateModified sortingKind = "Date Modified"
	sortingFileType     sortingKind = "Type"
)

var suffixRegexp = regexp.MustCompile(`^(.*)\((\d+)\)$`)

// Check if the directory is external disk path
// TODO : This function should be give two directories, and it should return
// if the two share a different disk partition.
// Ideally we shouldn't even try to figure that out in our file operations, and let OS handles it.
// But at least right now its not okay. This returns if `path` is an External disk
// from perspective of `/`, but it should tell from perspective of currently open directory
// The usage of this function in cut/paste is not as expected.
func isExternalDiskPath(path string) bool {
	// This is very vague. You cannot tell if a path is belonging to an external partition
	// if you dont define the source path to compare with
	// But making this true will cause slow file operations based on current implementation
	if runtime.GOOS == utils.OsWindows {
		return false
	}

	// exclude timemachine on MacOS
	if strings.HasPrefix(path, "/Volumes/.timemachine") {
		return false
	}

	// to filter out mounted partitions like /, /boot etc
	return strings.HasPrefix(path, "/mnt") ||
		strings.HasPrefix(path, "/media") ||
		strings.HasPrefix(path, "/run/media") ||
		strings.HasPrefix(path, "/Volumes")
}

func returnFocusType(focusPanel focusPanelType) bool {
	return focusPanel == nonePanelFocus
}

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
func returnDirElement(location string, displayDotFile bool, sortOptions sortOptionsModelData) []element {
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
	return sortFileElement(sortOptions, dirEntries, location)
}

func returnDirElementBySearchString(location string, displayDotFile bool, searchString string,
	sortOptions sortOptionsModelData,
) []element {
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

	return sortFileElement(sortOptions, dirElements, location)
}

func sortFileElement(sortOptions sortOptionsModelData, dirEntries []os.DirEntry, location string) []element {
	elements := make([]element, 0, len(dirEntries))
	for _, item := range dirEntries {
		info, err := item.Info()
		if err != nil {
			slog.Error("Error while retrieving file info during sort",
				"error", err, "path", filepath.Join(location, item.Name()))
			continue
		}

		elements = append(elements, element{
			name:      item.Name(),
			directory: item.IsDir() || isSymlinkToDir(location, info, item.Name()),
			location:  filepath.Join(location, item.Name()),
			info:      info,
		})
	}

	sort.Slice(elements, getOrderingFunc(elements,
		sortOptions.reversed, sortOptions.options[sortOptions.selected]))

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

func getOrderingFunc(elements []element, reversed bool, sortOption string) sliceOrderFunc {
	var order func(i, j int) bool
	switch sortOption {
	case string(sortingName):
		order = func(i, j int) bool {
			// One of them is a directory, and other is not
			if elements[i].directory != elements[j].directory {
				return elements[i].directory
			}
			if common.Config.CaseSensitiveSort {
				return elements[i].name < elements[j].name != reversed
			}
			return strings.ToLower(elements[i].name) < strings.ToLower(elements[j].name) != reversed
		}
	case string(sortingSize):
		order = getSizeOrderingFunc(elements, reversed)
	case string(sortingDateModified):
		order = func(i, j int) bool {
			return elements[i].info.ModTime().After(elements[j].info.ModTime()) != reversed
		}
	case string(sortingFileType):
		order = getTypeOrderingFunc(elements, reversed)
	}
	return order
}

func getSizeOrderingFunc(elements []element, reversed bool) sliceOrderFunc {
	return func(i, j int) bool {
		// Directories at the top sorted by direct child count (not recursive)
		// Files sorted by size

		// One of them is a directory, and other is not
		if elements[i].directory != elements[j].directory {
			return elements[i].directory
		}

		// This needs to be improved, and we should sort by actual size only
		// Repeated recursive read would be slow, so we could cache
		if elements[i].directory && elements[j].directory {
			filesI, err := os.ReadDir(elements[i].location)
			// No need of early return, we only call len() on filesI, so nil would
			// just result in 0
			if err != nil {
				slog.Error("Error when reading directory during sort", "error", err)
			}
			filesJ, err := os.ReadDir(elements[j].location)
			if err != nil {
				slog.Error("Error when reading directory during sort", "error", err)
			}
			return len(filesI) < len(filesJ) != reversed
		}
		return elements[i].info.Size() < elements[j].info.Size() != reversed
	}
}

func getTypeOrderingFunc(elements []element, reversed bool) sliceOrderFunc {
	return func(i, j int) bool {
		// One of them is a directory, and the other is not
		if elements[i].directory != elements[j].directory {
			return elements[i].directory
		}

		var extI, extJ string
		if !elements[i].directory {
			extI = strings.ToLower(filepath.Ext(elements[i].name))
		}
		if !elements[j].directory {
			extJ = strings.ToLower(filepath.Ext(elements[j].name))
		}

		// Compare by extension/type
		if extI != extJ {
			return (extI < extJ) != reversed
		}

		// If same type, fall back to name
		if common.Config.CaseSensitiveSort {
			return (elements[i].name < elements[j].name) != reversed
		}

		return (strings.ToLower(elements[i].name) < strings.ToLower(elements[j].name)) != reversed
	}
}

func panelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - 3
}

// TODO : replace usage of this with slices.contains
func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
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

func checkFileNameValidity(name string) error {
	switch {
	case name == ".", name == "..":
		return errors.New("file name cannot be '.' or '..'")
	case strings.HasSuffix(name, fmt.Sprintf("%c.", filepath.Separator)),
		strings.HasSuffix(name, fmt.Sprintf("%c..", filepath.Separator)):
		return fmt.Errorf("file name cannot end with '%c.' or '%c..'", filepath.Separator, filepath.Separator)
	default:
		return nil
	}
}

func renameIfDuplicate(destination string) (string, error) {
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		return destination, nil
	} else if err != nil {
		return "", err
	}

	dir := filepath.Dir(destination)
	base := filepath.Base(destination)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	// Extract base name without existing suffix
	counter := 1
	if match := suffixRegexp.FindStringSubmatch(name); len(match) == 3 {
		name = match[1] // base name without (N)
		if num, err := strconv.Atoi(match[2]); err == nil {
			counter = num + 1 // start from next number
		}
	}

	// Find first available name
	for i := counter; i < 10_000; i++ {
		newName := fmt.Sprintf("%s(%d)%s", name, i, ext)
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}
	}

	return "", fmt.Errorf("could not find free name for %s after many attempts", destination)
}

// TODO : Replace all usage of "m.fileModel.filePanels[m.filePanelFocusIndex]" with this
// There are many usage
func (m *model) getFocusedFilePanel() *filePanel {
	return &m.fileModel.filePanels[m.filePanelFocusIndex]
}

// Count how many file in the directory
func countFiles(dirPath string) (int, error) {
	count := 0

	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})

	return count, err
}

func processCmdToTeaCmd(cmd processbar.Cmd) tea.Cmd {
	if cmd == nil {
		// To prevent us from running cmd() on nil cmd
		return nil
	}
	return func() tea.Msg {
		updateMsg := cmd()
		return ProcessBarUpdateMsg{
			pMsg: updateMsg,
			BaseMessage: BaseMessage{
				reqID: updateMsg.GetReqID(),
			},
		}
	}
}

func getCopyOrCutOperationName(cut bool) string {
	if cut {
		return "cut"
	}
	return "copy"
}
