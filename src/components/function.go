package components

import "strings"

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
