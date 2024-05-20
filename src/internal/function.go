package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/shortuuid"
	"github.com/masatana/go-textdistance"
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

func returnFolderElement(location string, displayDotFile bool) (directoryElement []element) {

	files, err := os.ReadDir(location)
	if len(files) == 0 {
		return directoryElement
	}
	
	if err != nil {
		outPutLog("Return folder element function error", err)
	}
	
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})

	for _, item := range files {
		fileInfo, err := item.Info()
		if err != nil {
			continue
		}

		if !displayDotFile && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		if fileInfo == nil {
			continue
		}
		newElement := element{
			name:      item.Name(),
			directory: item.IsDir(),
		}
		if location == "/" {
			newElement.location = location + item.Name()
		} else {
			newElement.location = filepath.Join(location, item.Name())
		}
		directoryElement = append(directoryElement, newElement)
	}

	return directoryElement
}

func returnFolderElementBySearchString(location string, displayDotFile bool, searchString string) (folderElement []element) {

	items, err := os.ReadDir(location)
	if err != nil {
		outPutLog("Return folder element function error", err)
	}

	for _, item := range items {
		fileInfo, _ := item.Info()
		if !displayDotFile && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		if fileInfo == nil {
			continue
		}
		newElement := element{
			name:      item.Name(),
			directory: item.IsDir(),
			matchRate: textdistance.JaroWinklerDistance(item.Name(), searchString),
		}
		if location == "/" {
			newElement.location = location + item.Name()
		} else {
			newElement.location = location + "/" + item.Name()
		}
		if newElement.matchRate > 0 {
			folderElement = append(folderElement, newElement)
		}
	}

	// Sort folders and files by match rate
	sort.Slice(folderElement, func(i, j int) bool {
		return folderElement[i].matchRate > folderElement[j].matchRate
	})

	return folderElement
}

func panelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - 3
}

func bottomElementHeight(bottomElementHeight int) int {
	return bottomElementHeight - 5
}

func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func outPutLog(values ...interface{}) {
	log.SetOutput(logOutput)
	for _, value := range values {
		log.Println(value)
	}
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

func pasteFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		outPutLog("Paste file function open file error", err)
	}
	defer srcFile.Close()

	dst, err = renameIfDuplicate(dst)
	if err != nil {
		outPutLog("Paste file function rename error", err)
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		outPutLog("Paste file function create file error", err)
	}
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		outPutLog("Paste file function copy file error", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func returnMetaData(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	cursor := panel.cursor

	// Obtaining metadata will take time. If metadata is obtained for every passing file, it will cause lag.
	// Therefore, it is necessary to detect whether it is just browsing or stopping on that file or directory.
	LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), cursor}
	time.Sleep(150 * time.Millisecond)

	if LastTimeCursorMove[1] != cursor && m.focusPanel != metadataFocus {
		return m
	}

	m.fileMetaData.metaData = m.fileMetaData.metaData[:0]
	id := shortuuid.New()
	if len(panel.element) == 0 {
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m
	}
	if len(panel.element[panel.cursor].metaData) != 0 && m.focusPanel != metadataFocus {
		m.fileMetaData.metaData = panel.element[panel.cursor].metaData
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m
	}
	filePath := panel.element[panel.cursor].location

	fileInfo, err := os.Stat(filePath)

	if isSymlink(filePath) {
		if isBrokenSymlink(filePath) {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"Link file is broken!", ""})
		} else {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"This is a link file.", ""})
		}
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m

	}

	if err != nil {
		outPutLog("Return meta data function get file state error", err)
	}

	if fileInfo.IsDir() {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderName", fileInfo.Name()})
		if m.focusPanel == metadataFocus {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderSize", formatFileSize(dirSize(filePath))})
		}
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderModifyDate", fileInfo.ModTime().String()})
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m
	}

	checkIsSymlinked, err := os.Lstat(filePath)
	if err != nil {
		outPutLog("err when getting file info", err)
		return m
	}

	if Config.Metadata && checkIsSymlinked.Mode()&os.ModeSymlink == 0 {

		fileInfos := et.ExtractMetadata(filePath)

		for _, fileInfo := range fileInfos {
			if fileInfo.Err != nil {
				outPutLog("Return meta data function error", fileInfo, fileInfo.Err)
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
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, fileName, fileSize, fileModifyData)
	}

	channel <- channelMessage{
		messageId:    id,
		loadMetadata: true,
		metadata:     m.fileMetaData.metaData,
	}

	panel.element[panel.cursor].metaData = m.fileMetaData.metaData
	return m
}

// Get directory total size
func dirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			outPutLog("Dir size function error", err)
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
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

// Generate search bar for file panel
func generateSearchBar() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = footerCursorStyle
	ti.Cursor.TextStyle = footerStyle
	ti.TextStyle = filePanelStyle
	ti.Prompt = filePanelTopDirectoryIconStyle.Render(" ")
	ti.Cursor.Blink = true
	ti.PlaceholderStyle = filePanelStyle
	ti.Placeholder = "(" + hotkeys.SearchBar[0] + ") Type something"
	ti.Blur()
	ti.CharLimit = 156
	return ti
}

// Check whether is broken recursive symlinks
func isBrokenSymlink(filePath string) bool {
	linkPath, err := os.Readlink(filePath)
	if err != nil {
		return true
	}

	absLinkPath, err := filepath.Abs(linkPath)
	if err != nil {
		return true
	}

	_, err = os.Stat(absLinkPath)
	return err != nil
}

// Check whether is symlinks
func isSymlink(filePath string) bool {

	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return true
	}

	return fileInfo.Mode()&os.ModeSymlink != 0
}
