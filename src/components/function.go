package components

import (
	"log"
	"os"
	"os/user"
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
		{location: "~/", name: "󰋜 Home"},
		{location: "~/Downloads", name: "󰏔 Downloads"},
		{location: "~/Documents", name: "󰈙 Documents"},
		{location: "~/Pictures", name: "󰋩 Pictures"},
		{location: "~/Videos", name: "󰎁 Videos"},
		{location: "~/Documents/code", name: "code"},
		{location: "~/Documents/code/returnone", name: "returnone"},
		{location: "~/Documents/code/returnone/website", name: "website"},
		{location: "~/Documents/code/returnone/backend", name: "backend", endPinned: true},
		{location: "~/Documents/code/returnone/backend/dsa", name: "Movie disk"},
		{location: "~/Documents/code/returnone/backend/dsaa", name: "USB 3.2 SanDisk"},
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
	files, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		folderElement = append(folderElement, element{
			name: file.Name(), 
			location: location + file.Name(),
			folder: file.IsDir(),
		})
	}
	return folderElement
}
