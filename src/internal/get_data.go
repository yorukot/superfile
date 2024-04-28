package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rkoesters/xdg/userdirs"
	"github.com/shirou/gopsutil/disk"
)

// Return all sidebar directories
func getDirectories(height int) []directory {
	directories := []directory{}

	directories = append(directories, getWellKnownDirectories()...)
	if height > 30 {
		directories = append(directories, getPinnedDirectories()...)
		directories = append(directories, getExternalMediaFolders()...)
	}
	return directories
}

// Return system default directory e.g. Home, Downloads, etc
func getWellKnownDirectories() []directory {
	directories := []directory{}
	wellKnownDirectories := []directory{
		{location: HomeDir, name: "󰋜 Home"},
		{location: userdirs.Download, name: "󰏔 Downloads"},
		{location: userdirs.Documents, name: "󰈙 Documents"},
		{location: userdirs.Pictures, name: "󰋩 Pictures"},
		{location: userdirs.Videos, name: "󰎁 Videos"},
		{location: userdirs.Music, name: "♬ Music"},
		{location: userdirs.Templates, name: "󰏢 Templates"},
		{location: userdirs.PublicShare, name: " PublicShare"},
	}

	if runtime.GOOS == "darwin" {
		wellKnownDirectories[1].location = HomeDir + "/Downloads/"
		wellKnownDirectories[2].location = HomeDir + "/Documents/"
		wellKnownDirectories[3].location = HomeDir + "/Pictures/"
		wellKnownDirectories[4].location = HomeDir + "/Movies/"
		wellKnownDirectories[5].location = HomeDir + "/Music/"
		wellKnownDirectories[7].location = HomeDir + "/Public/"
	}

	for _, dir := range wellKnownDirectories {
		if _, err := os.Stat(dir.location); !os.IsNotExist(err) {
			// Directory exists
			directories = append(directories, dir)
		}
	}
	return directories
}

// Get user pinned directories
func getPinnedDirectories() []directory {
	directories := []directory{}
	var paths []string

	jsonData, err := os.ReadFile(SuperFileDataDir + pinnedFile)
	if err != nil {
		outPutLog("Read superfile data error", err)
	}

	json.Unmarshal(jsonData, &paths)

	for _, path := range paths {
		directoryName := filepath.Base(path)
		directories = append(directories, directory{location: path, name: directoryName})
	}
	return directories
}

// Get external media directories 
func getExternalMediaFolders() (disks []directory) {
	parts, err := disk.Partitions(true)

	if err != nil {
		outPutLog("Error while getting external media: ", err)
	}
	for _, disk := range parts {
		if isExternalDiskPath(disk.Mountpoint) {
			disks = append(disks, directory{
				name:     filepath.Base(disk.Mountpoint),
				location: disk.Mountpoint,
			})
		}
	}
	if err != nil {
		outPutLog("Error while getting external media: ", err)
	}
	return disks
}
