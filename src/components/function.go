package components

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func getHomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("can't get home dir")
	}
	return user.HomeDir
}

func getFolder() []folder {
	var paths []string

	currentUser, err := user.Current()
	if err != nil {
		OutPutLog("Get user path error", err)
	}
	username := currentUser.Username

	folderPath := filepath.Join("/run/media", username)
	entries, err := os.ReadDir(folderPath)

	if err != nil {
		OutPutLog("Get external media error", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			paths = append(paths, filepath.Join(folderPath, entry.Name()))
		}
	}
	jsonData, err := os.ReadFile(SuperFileMainDir + pinnedFile)
	if err != nil {
		OutPutLog("Read superfile data error", err)
	}
	var pinnedFolder []string
	err = json.Unmarshal(jsonData, &pinnedFolder)
	if err != nil {
		OutPutLog("Unmarshal superfile data error", err)
	}
	folders := []folder{
		{location: HomeDir, name: "󰋜 Home"},
		{location: HomeDir + "/Downloads", name: "󰏔 Downloads"},
		{location: HomeDir + "/Documents", name: "󰈙 Documents"},
		{location: HomeDir + "/Pictures", name: "󰋩 Pictures"},
		{location: HomeDir + "/Videos", name: "󰎁 Videos"},
	}

	for i, path := range pinnedFolder {
		folderName := filepath.Base(path)
		if i == len(pinnedFolder)-1 {
			folders = append(folders, folder{location: path, name: folderName, endPinned: true})
		} else {
			folders = append(folders, folder{location: path, name: folderName})
		}
	}
	for _, path := range paths {
		folderName := filepath.Base(path)
		folders = append(folders, folder{location: path, name: folderName})
	}

	return folders
}

func repeatString(s string, count int) string {
	return strings.Repeat(s, count)
}

func returnFocusType(focusPanel focusPanelType) filePanelFocusType {
	if focusPanel == nonePanelFocus {
		return focus
	} else {
		return secondFocus
	}
}

func returnFolderElement(location string) (folderElement []element) {
	var folders []element
	var files []element

	items, err := os.ReadDir(location)
	if err != nil {
		OutPutLog("Return folder element function error", err)
	}

	for _, item := range items {
		fileInfo, _ := item.Info()
		if fileInfo == nil {
			continue
		}
		newElement := element{
			name:   item.Name(),
			folder: item.IsDir(),
		}
		if location == "/" {
			newElement.location = location + item.Name()
		} else {
			newElement.location = location + "/" + item.Name()
		}

		if item.IsDir() {
			folders = append(folders, newElement)
		} else {
			files = append(files, newElement)
		}
	}

	// Sort folders and files alphabetically
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].name < folders[j].name
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	// Concatenate folders and files
	folderElement = append(folders, files...)

	return folderElement
}

func PanelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - 3
}

func BottomElementHight(bottomElementHight int) int {
	return bottomElementHight - 5
}

func ArrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func OutPutLog(values ...interface{}) {
	log.SetOutput(logOutput)
	for _, value := range values {
		log.Println(value)
	}
}

func RemoveElementByValue(slice []string, value string) []string {
	newSlice := []string{}
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func RenameIfDuplicate(destination string) (string, error) {
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

func MoveFile(source string, destination string) error {
	destination, err := RenameIfDuplicate(destination)
	if err != nil {
		OutPutLog("Move file function error", err)
	}
	err = os.Rename(source, destination)
	if err != nil {
		OutPutLog("Move file function error", err)
	}
	return err
}

func PasteFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		OutPutLog("Paste file function open file error", err)
	}
	defer srcFile.Close()

	dst, err = RenameIfDuplicate(dst)
	if err != nil {
		OutPutLog("Paste file function rename error", err)
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		OutPutLog("Paste file function create file error", err)
	}
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		OutPutLog("Paste file function copy file error", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func PasteDir(src, dst string, id string, m model) (model, error) {
	// Check if destination directory already exists
	dst, err := RenameIfDuplicate(dst)
	if err != nil {
		return m, err
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			newPath, err = RenameIfDuplicate(newPath)
			if err != nil {
				return err
			}
			err = os.MkdirAll(newPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			p := m.processBarModel.process[id]
			if m.copyItems.cut {
				p.name = "󰆐 " + filepath.Base(path)
			} else {
				p.name = "󰆏 " + filepath.Base(path)
			}

			if len(processBarChannel) < 5 {
				processBarChannel <- processBarMessage{
					processId:       id,
					processNewState: p,
				}
			}

			err := PasteFile(path, newPath)
			if err != nil {
				p.state = failure
				processBarChannel <- processBarMessage{
					processId:       id,
					processNewState: p,
				}
				return err
			}
			p.done++
			if len(processBarChannel) < 5 {
				processBarChannel <- processBarMessage{
					processId:       id,
					processNewState: p,
				}
			}
			m.processBarModel.process[id] = p
		}

		return nil
	})

	if err != nil {
		return m, err
	}

	return m, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ReturnMetaData(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if len(panel.element) == 0 {
		return m
	}
	if len(panel.element[panel.cursor].metaData) != 0 && m.focusPanel != metaDataFocus {
		m.fileMetaData.metaData = panel.element[panel.cursor].metaData
		return m
	}
	if len(panel.element) == 0 {
		return m
	}
	m.fileMetaData.metaData = m.fileMetaData.metaData[:0]
	filePath := panel.element[panel.cursor].location

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"Link file is broken!(you can only delete this file)", ""})
		return m
	}
	if err != nil {
		OutPutLog("Return meta data function get file state error", err)
	}
	if fileInfo.IsDir() {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderName", fileInfo.Name()})
		if m.focusPanel == metaDataFocus {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderSize", FormatFileSize(DirSize(filePath))})
		}
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderModifyDate", fileInfo.ModTime().String()})
		return m
	}

	fileInfos := et.ExtractMetadata(filePath)
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			OutPutLog("Return meta data function error", fileInfo, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			temp := [2]string{k, fmt.Sprintf("%v", v)}
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, temp)
		}
	}
	panel.element[panel.cursor].metaData = m.fileMetaData.metaData
	return m
}

func FormatFileSize(size int64) string {
	units := []string{" bytes", " kB", " MB", " GB", " TB", " PB", " EB"}

	if size == 0 {
		return "0B"
	}

	unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(1024)))
	adjustedSize := float64(size) / math.Pow(1024, float64(unitIndex))

	return fmt.Sprintf("%.2f%s", adjustedSize, units[unitIndex])
}

func DirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			OutPutLog("Dir size function error", err)
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size
}

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
