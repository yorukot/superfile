package components

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
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

    // 读取文件夹内容
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
	OutputLog(paths)
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
			name:       item.Name(),
			folder:     item.IsDir(),
			size:       fileInfo.Size(),
			updateTime: fileInfo.ModTime(),
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

func PasteFile(src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	CheckErr(err)
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	CheckErr(err)
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
