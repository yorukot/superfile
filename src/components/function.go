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

func returnFocusType(sideBarFocus bool) filePanelFocusType {
	if sideBarFocus {
		return secondFocus
	} else {
		return focus
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
