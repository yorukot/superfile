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

func returnFocusType(focusPanel focusPanelType) filePanelFocusType {
	if focusPanel == nonePanelFocus {
		return focus
	}
	return secondFocus
}

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
func returnDirElement(location string, displayDotFile bool, sortOptions sortOptionsModelData) []element {
	dirEntries, err := os.ReadDir(location)
	if err != nil {
		slog.Error("Error while return folder element function", "error", err)
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

	// Sort files
	sort.Slice(dirEntries, getOrderingFunc(location, dirEntries,
		sortOptions.reversed, sortOptions.options[sortOptions.selected]))
	// Preallocate for efficiency
	directoryElement := make([]element, 0, len(dirEntries))
	for _, item := range dirEntries {
		directoryElement = append(directoryElement, element{
			name:      item.Name(),
			directory: item.IsDir(),
			location:  filepath.Join(location, item.Name()),
		})
	}
	return directoryElement
}

func getOrderingFunc(location string, dirEntries []os.DirEntry, reversed bool, sortOption string) sliceOrderFunc {
	var order func(i, j int) bool
	// TODO : These strings should not be hardcoded here, but defined as constants
	switch sortOption {
	case "Name":
		order = func(i, j int) bool {
			// One of them is a directory, and other is not
			if dirEntries[i].IsDir() != dirEntries[j].IsDir() {
				return dirEntries[i].IsDir()
			}
			if common.Config.CaseSensitiveSort {
				return dirEntries[i].Name() < dirEntries[j].Name() != reversed
			}
			return strings.ToLower(dirEntries[i].Name()) < strings.ToLower(dirEntries[j].Name()) != reversed
		}
	case "Size":
		order = getSizeOrderingFunc(dirEntries, reversed, location)
	case "Date Modified":
		order = func(i, j int) bool {
			// No need for err check, we already filtered out dirEntries with err != nil in Info() call
			fileInfoI, _ := dirEntries[i].Info()
			fileInfoJ, _ := dirEntries[j].Info()
			return fileInfoI.ModTime().After(fileInfoJ.ModTime()) != reversed
		}
	case "Type":
		order = getTypeOrderingFunc(dirEntries, reversed)
	}
	return order
}

func getSizeOrderingFunc(dirEntries []os.DirEntry, reversed bool, location string) sliceOrderFunc {
	return func(i, j int) bool {
		// Directories at the top sorted by direct child count (not recursive)
		// Files sorted by size

		// One of them is a directory, and other is not
		if dirEntries[i].IsDir() != dirEntries[j].IsDir() {
			return dirEntries[i].IsDir()
		}

		// This needs to be improved, and we should sort by actual size only
		// Repeated recursive read would be slow, so we could cache
		if dirEntries[i].IsDir() && dirEntries[j].IsDir() {
			filesI, err := os.ReadDir(filepath.Join(location, dirEntries[i].Name()))
			// No need of early return, we only call len() on filesI, so nil would
			// just result in 0
			if err != nil {
				slog.Error("Error when reading directory during sort", "error", err)
			}
			filesJ, err := os.ReadDir(filepath.Join(location, dirEntries[j].Name()))
			if err != nil {
				slog.Error("Error when reading directory during sort", "error", err)
			}
			return len(filesI) < len(filesJ) != reversed
		}
		// No need for err check, we already filtered out dirEntries with err != nil in Info() call
		fileInfoI, _ := dirEntries[i].Info()
		fileInfoJ, _ := dirEntries[j].Info()
		return fileInfoI.Size() < fileInfoJ.Size() != reversed
	}
}

func getTypeOrderingFunc(dirEntries []os.DirEntry, reversed bool) sliceOrderFunc {
	return func(i, j int) bool {
		// One of them is a directory, and the other is not
		if dirEntries[i].IsDir() != dirEntries[j].IsDir() {
			return dirEntries[i].IsDir()
		}

		var extI, extJ string
		if !dirEntries[i].IsDir() {
			extI = strings.ToLower(filepath.Ext(dirEntries[i].Name()))
		}
		if !dirEntries[j].IsDir() {
			extJ = strings.ToLower(filepath.Ext(dirEntries[j].Name()))
		}

		// Compare by extension/type
		if extI != extJ {
			return (extI < extJ) != reversed
		}

		// If same type, fall back to name
		if common.Config.CaseSensitiveSort {
			return (dirEntries[i].Name() < dirEntries[j].Name()) != reversed
		}

		return (strings.ToLower(dirEntries[i].Name()) < strings.ToLower(dirEntries[j].Name())) != reversed
	}
}

func returnDirElementBySearchString(location string, displayDotFile bool, searchString string) []element {
	items, err := os.ReadDir(location)
	if err != nil {
		slog.Error("Error while return folder element function", "error", err)
		return []element{}
	}

	if len(items) == 0 {
		return []element{}
	}

	folderElementMap := map[string]element{}
	fileAndDirectories := []string{}

	for _, item := range items {
		fileInfo, err := item.Info()
		if err != nil {
			continue
		}
		if !displayDotFile && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}

		folderElementLocation := filepath.Join(location, item.Name())

		fileAndDirectories = append(fileAndDirectories, item.Name())
		folderElementMap[item.Name()] = element{
			name:      item.Name(),
			directory: item.IsDir(),
			location:  folderElementLocation,
		}
	}
	// https://github.com/reinhrst/fzf-lib/blob/main/core.go#L43
	// No sorting needed. fzf.DefaultOptions() already return values ordered on Score
	fzfResults := utils.FzfSearch(searchString, fileAndDirectories)
	dirElement := make([]element, 0, len(fzfResults))
	for _, item := range fzfResults {
		resultItem := folderElementMap[item.Key]
		dirElement = append(dirElement, resultItem)
	}

	return dirElement
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
	info, err := os.Stat(destination)
	if os.IsNotExist(err) {
		return destination, nil
	} else if err != nil {
		return "", err
	}

	if info.IsDir() {
		match := regexp.MustCompile(`\((\d+)\)$`).FindStringSubmatch(info.Name())
		if len(match) > 1 {
			number, _ := strconv.Atoi(match[1])
			for {
				number++
				newDirName := fmt.Sprintf("%s(%d)", info.Name()[:len(info.Name())-len(match[0])], number)
				newPath := filepath.Join(filepath.Dir(destination), newDirName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		} else {
			for i := 1; ; i++ {
				newDirName := fmt.Sprintf("%s(%d)", info.Name(), i)
				newPath := filepath.Join(filepath.Dir(destination), newDirName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		}
	} else {
		baseName := filepath.Base(destination)
		ext := filepath.Ext(baseName)
		fileName := baseName[:len(baseName)-len(ext)]
		match := regexp.MustCompile(`\((\d+)\)$`).FindStringSubmatch(fileName)
		if len(match) > 1 {
			number, _ := strconv.Atoi(match[1])
			for {
				number++
				newFileName := fmt.Sprintf("%s(%d)%s", fileName[:len(fileName)-len(match[0])], number, ext)
				newPath := filepath.Join(filepath.Dir(destination), newFileName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		} else {
			for i := 1; ; i++ {
				newFileName := fmt.Sprintf("%s(%d)%s", fileName, i, ext)
				newPath := filepath.Join(filepath.Dir(destination), newFileName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		}
	}
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

func isImageFile(filename string) bool {
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".svg":  true,
		".webp": true,
		".ico":  true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return imageExtensions[ext]
}

func processCmdToTeaCmd(cmd processbar.Cmd) tea.Cmd {
	if cmd == nil {
		// To prevent us from running cmd() on nil cmd
		return nil
	}
	return func() tea.Msg {
		updateMsg := cmd()
		return ProcessBarUpdateMsg{pMsg: updateMsg,
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
