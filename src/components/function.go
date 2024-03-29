package components

import (
	"log"
	"os"
	"os/user"
	"sort"
	"strings"
)

func getHomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return user.HomeDir
}

func getFolder() []folder {
	folders := []folder{
		{location: HomeDir, name: "󰋜 Home"},
		{location: HomeDir + "/Downloads", name: "󰏔 Downloads"},
		{location: HomeDir + "/Documents", name: "󰈙 Documents"},
		{location: HomeDir + "/Pictures", name: "󰋩 Pictures"},
		{location: HomeDir + "/Videos", name: "󰎁 Videos"},
		{location: HomeDir + "/Documents/code", name: "code"},
		{location: HomeDir + "/Documents/code/returnone", name: "returnone"},
		{location: HomeDir + "/Documents/code/returnone/website", name: "website"},
		{location: HomeDir + "/Documents/code/returnone/backend", name: "backend", endPinned: true},
		{location: HomeDir + "/Documents/code/returnone/backend/dsa", name: "Movie disk"},
		{location: HomeDir + "/Documents/code/returnone/backend/dsaa", name: "USB 3.2 SanDisk"},
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
		log.Fatal(err)
	}

	for _, item := range items {
		fileInfo, _ := item.Info()
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
    if err != nil {
        OutputLog(err)
    }
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