package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid"
	"github.com/reinhrst/fzf-lib"
	"github.com/yorukot/superfile/src/config/icon"
)

// Check if the directory is external disk path
func isExternalDiskPath(path string) bool {
	dir := filepath.Dir(path)

	// exclude timemachine
	if strings.HasPrefix(dir, "/Volumes/.timemachine") {
		return false
	}

	return strings.HasPrefix(dir, "/mnt") ||
		strings.HasPrefix(dir, "/media") ||
		strings.HasPrefix(dir, "/run/media") ||
		strings.HasPrefix(dir, "/Volumes")
}

func returnFocusType(focusPanel focusPanelType) filePanelFocusType {
	if focusPanel == nonePanelFocus {
		return focus
	}
	return secondFocus
}

func returnDirElement(location string, displayDotFile bool, sortOptions sortOptionsModelData) (directoryElement []element) {
	dirEntries, err := os.ReadDir(location)
	if err != nil {
		slog.Error("Return folder element function error", "error", err)
		return directoryElement
	}

	dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
		// Entries not needed to be considered
		_, err := e.Info()
		return err != nil || (strings.HasPrefix(e.Name(), ".") && !displayDotFile)
	})

	// No files/directoes to process
	if len(dirEntries) == 0 {
		return directoryElement
	}

	// Sort files
	var order func(i, j int) bool
	reversed := sortOptions.reversed

	// Todo : These strings should not be hardcoded here, but defined as constants
	switch sortOptions.options[sortOptions.selected] {
	case "Name":
		order = func(i, j int) bool {
			// One of them is a directory, and other is not
			if dirEntries[i].IsDir() != dirEntries[j].IsDir() {
				return dirEntries[i].IsDir()
			}
			if Config.CaseSensitiveSort {
				return dirEntries[i].Name() < dirEntries[j].Name() != reversed
			} else {
				return strings.ToLower(dirEntries[i].Name()) < strings.ToLower(dirEntries[j].Name()) != reversed
			}
		}
	case "Size":
		order = func(i, j int) bool {
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
			} else {
				// No need for err check, we already filtered out dirEntries with err != nil in Info() call
				fileInfoI, _ := dirEntries[i].Info()
				fileInfoJ, _ := dirEntries[j].Info()
				return fileInfoI.Size() < fileInfoJ.Size() != reversed
			}

		}
	case "Date Modified":
		order = func(i, j int) bool {
			// No need for err check, we already filtered out dirEntries with err != nil in Info() call
			fileInfoI, _ := dirEntries[i].Info()
			fileInfoJ, _ := dirEntries[j].Info()
			return fileInfoI.ModTime().After(fileInfoJ.ModTime()) != reversed
		}
	}

	sort.Slice(dirEntries, order)
	for _, item := range dirEntries {
		directoryElement = append(directoryElement, element{
			name:      item.Name(),
			directory: item.IsDir(),
			location:  filepath.Join(location, item.Name()),
		})
	}
	return directoryElement
}

func returnDirElementBySearchString(location string, displayDotFile bool, searchString string) (dirElement []element) {

	items, err := os.ReadDir(location)
	if err != nil {
		slog.Error("Return folder element function error", "error", err)
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
	for _, item := range fzfSearch(searchString, fileAndDirectories) {
		resultItem := folderElementMap[item.Key]
		dirElement = append(dirElement, resultItem)
	}

	return dirElement
}

// Returning a string slice causes inefficiency in current usage
func fzfSearch(query string, source []string) []fzf.MatchResult {
	fzfSearcher := fzf.New(source, fzf.DefaultOptions())
	fzfSearcher.Search(query)
	fzfResults := <-fzfSearcher.GetResultChannel()
	fzfSearcher.End()
	return fzfResults.Matches
}

func panelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - 3
}

func bottomElementHeight(bottomElementHeight int) int {
	return bottomElementHeight - 5
}

// Todo : replace usage of this with slices.contains
func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Todo : Eventually we want to remove all such usage that can result in app exiting abruptly
func LogAndExit(msg string, values ...any) {
	slog.Error(msg, values...)
	os.Exit(1)
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

// Todo : Move this model related function to the right file. Should not be in functions file.
func (m *model) returnMetaData() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	cursor := panel.cursor
	id := shortuuid.New()

	message := channelMessage{
		messageId:   id,
		messageType: sendMetadata,
		metadata:    m.fileMetaData.metaData,
	}

	// Obtaining metadata will take time. If metadata is obtained for every passing file, it will cause lag.
	// Therefore, it is necessary to detect whether it is just browsing or stopping on that file or directory.
	LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), cursor}
	time.Sleep(150 * time.Millisecond)

	if LastTimeCursorMove[1] != cursor && m.focusPanel != metadataFocus {
		return
	}

	m.fileMetaData.metaData = m.fileMetaData.metaData[:0]
	if len(panel.element) == 0 {
		message.metadata = m.fileMetaData.metaData
		channel <- message
		return
	}
	if len(panel.element[panel.cursor].metaData) != 0 && m.focusPanel != metadataFocus {
		m.fileMetaData.metaData = panel.element[panel.cursor].metaData
		message.metadata = m.fileMetaData.metaData
		channel <- message
		return
	}
	filePath := panel.element[panel.cursor].location

	fileInfo, err := os.Stat(filePath)

	if isSymlink(filePath) {
		_, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"Link file is broken!", ""})
		} else {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"This is a link file.", ""})
		}
		message.metadata = m.fileMetaData.metaData
		channel <- message
		return

	}

	if err != nil {
		slog.Error("Error while getting file state", "error", err)
	}

	if fileInfo.IsDir() {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderName", fileInfo.Name()})
		if m.focusPanel == metadataFocus {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderSize", formatFileSize(dirSize(filePath))})
		}
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderModifyDate", fileInfo.ModTime().String()})
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderPermissions", fileInfo.Mode().String()})
		message.metadata = m.fileMetaData.metaData
		channel <- message
		return
	}

	checkIsSymlinked, err := os.Lstat(filePath)
	if err != nil {
		slog.Error("Error while getting file info", "error", err)
		return
	}

	if Config.Metadata && checkIsSymlinked.Mode()&os.ModeSymlink == 0 && et != nil {

		fileInfos := et.ExtractMetadata(filePath)

		for _, fileInfo := range fileInfos {
			if fileInfo.Err != nil {
				slog.Error("Return meta data function error", "error", fileInfo.Err)
				continue
			}

			for k, v := range fileInfo.Fields {
				temp := [2]string{k, fmt.Sprintf("%v", v)}
				m.fileMetaData.metaData = append(m.fileMetaData.metaData, temp)
			}
		}
	} else {
		fileName := [2]string{"FileName", fileInfo.Name()}
		fileSize := [2]string{"FileSize", formatFileSize(fileInfo.Size())}
		fileModifyData := [2]string{"FileModifyDate", fileInfo.ModTime().String()}
		filePermissions := [2]string{"FilePermissions", fileInfo.Mode().String()}

		if Config.EnableMD5Checksum {
			// Calculate MD5 checksum
			checksum, err := calculateMD5Checksum(filePath)
			if err != nil {
				slog.Error("Error calculating MD5 checksum", "error", err)
			} else {
				md5Data := [2]string{"MD5Checksum", checksum}
				m.fileMetaData.metaData = append(m.fileMetaData.metaData, md5Data)
			}
		}

		m.fileMetaData.metaData = append(m.fileMetaData.metaData, fileName, fileSize, fileModifyData, filePermissions)
	}

	message.metadata = m.fileMetaData.metaData
	channel <- message

	panel.element[panel.cursor].metaData = m.fileMetaData.metaData
}

func calculateMD5Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5 checksum: %v", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

// Get directory total size
func dirSize(path string) int64 {
	var size int64
	err := filepath.WalkDir(path, func(_ string, entry os.DirEntry, err error) error {
		if err != nil {
			slog.Error("Error while getting directory size", "error", err)
		}
		if !entry.IsDir() {
			info, err := entry.Info()
			if err == nil {
				size += info.Size()
			}
		}
		return err
	})
	if err != nil {
		slog.Error("Errors during WalkDir", "error", err)
	}
	return size
}

// Count how many file in the directory
func countFiles(dirPath string) (int, error) {
	count := 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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

// Check whether is symlinks
func isSymlink(filePath string) bool {

	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return true
	}
	return fileInfo.Mode()&os.ModeSymlink != 0
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

func getElementIcon(file string, IsDir bool) icon.IconStyle {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	name := file

	if !Config.Nerdfont {
		return icon.IconStyle{
			Icon:  "",
			Color: theme.FilePanelFG,
		}
	}

	if IsDir {
		resultIcon := icon.Folders["folder"]
		betterIcon, hasBetterIcon := icon.Folders[name]
		if hasBetterIcon {
			resultIcon = betterIcon
		}
		return resultIcon
	} else {
		// default icon for all files. try to find a better one though...
		resultIcon := icon.Icons["file"]
		// resolve aliased extensions
		extKey := strings.ToLower(ext)
		alias, hasAlias := icon.Aliases[extKey]
		if hasAlias {
			extKey = alias
		}

		// see if we can find a better icon based on extension alone
		betterIcon, hasBetterIcon := icon.Icons[extKey]
		if hasBetterIcon {
			resultIcon = betterIcon
		}

		// now look for icons based on full names
		fullName := name

		fullName = strings.ToLower(fullName)
		fullAlias, hasFullAlias := icon.Aliases[fullName]
		if hasFullAlias {
			fullName = fullAlias
		}
		bestIcon, hasBestIcon := icon.Icons[fullName]
		if hasBestIcon {
			resultIcon = bestIcon
		}
		if resultIcon.Color == "NONE" {
			return icon.IconStyle{
				Icon:  resultIcon.Icon,
				Color: theme.FilePanelFG,
			}
		}
		return resultIcon
	}
}
