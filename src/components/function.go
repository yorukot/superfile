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
	CheckErr(err)
	username := currentUser.Username

	folderPath := filepath.Join("/run/media", username)
	entries, err := os.ReadDir(folderPath)
	CheckErr(err)

	for _, entry := range entries {
		if entry.IsDir() {
			paths = append(paths, filepath.Join(folderPath, entry.Name()))
		}
	}
	CheckErr(err)

	jsonData, err := os.ReadFile("./.superfile/data/superfile.json")
	CheckErr(err)

	var pinnedFolder []string
	err = json.Unmarshal(jsonData, &pinnedFolder)
	CheckErr(err)

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
	CheckErr(err)

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

func ArrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func OutputLog(value any) {
	log.SetOutput(logOutput)
	log.Println(value)
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

func MoveFile(source string, destination string) error {
	err := os.Rename(source, destination)
	CheckErr(err)
	return err
}

func PasteFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	CheckErr(err)
	defer srcFile.Close()

	dstDir := filepath.Dir(dst)
	baseName := filepath.Base(dst)
	ext := filepath.Ext(baseName)
	fileName := strings.TrimSuffix(baseName, ext)

	newFileName := fileName
	newDst := dst
	i := 1
	for {
		_, err := os.Stat(newDst)
		if os.IsNotExist(err) {
			break
		}

		newFileName = fileName + "(" + strconv.Itoa(i) + ")"
		newDst = filepath.Join(dstDir, newFileName+ext)
		i++
	}

	dstFile, err := os.Create(newDst)
	CheckErr(err)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	CheckErr(err)
	if err != nil {
		return err
	}
	return nil
}

func PasteDir(src, dst string, id string, m model) (model, error) {

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			_, err := os.Stat(newPath)
			if err == nil {
				for i := 1; ; i++ {
					newDir := fmt.Sprintf("%s(%d)", newPath, i)
					_, err := os.Stat(newDir)
					if err != nil {
						newPath = newDir
						break
					}
				}
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

			processBarChannel <- processBarMessage{
				processId:       id,
				processNewState: p,
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
			processBarChannel <- processBarMessage{
				processId:       id,
				processNewState: p,
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

func CheckErr(err error) {
	if err != nil {
		OutputLog(err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func returnMetaData(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	m.fileMetaData.metaData = m.fileMetaData.metaData[:0]
	filePath := panel.element[panel.cursor].location

	fileInfo, err := os.Stat(filePath)
	CheckErr(err)

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
			OutputLog(fileInfo.Err)
			OutputLog(fileInfo)
			continue
		}

		for k, v := range fileInfo.Fields {
			temp := [2]string{k, fmt.Sprintf("%v", v)}
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, temp)
		}
	}
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
			OutputLog(err)
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
