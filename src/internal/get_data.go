package internal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/shirou/gopsutil/disk"
	varibale "github.com/yorukot/superfile/src/config"
)

// Return all sidebar directories
func getDirectories() []directory {
	directories := []directory{}

	directories = append(directories, getWellKnownDirectories()...)
	directories = append(directories, directory{
		// Just make sure no one owns the hard drive or directory named this path
		location: "Pinned+-*/=?",
	})
	directories = append(directories, getPinnedDirectories()...)
	directories = append(directories, directory{
		// Just make sure no one owns the hard drive or directory named this path
		location: "Disks+-*/=?",
	})
	directories = append(directories, getExternalMediaFolders()...)

	return directories
}

// Return system default directory e.g. Home, Downloads, etc
func getWellKnownDirectories() []directory {
	directories := []directory{}
	wellKnownDirectories := []directory{
		{location: xdg.Home, name: "󰋜 Home"},
		{location: xdg.UserDirs.Download, name: "󰏔 Downloads"},
		{location: xdg.UserDirs.Documents, name: "󰈙 Documents"},
		{location: xdg.UserDirs.Pictures, name: "󰋩 Pictures"},
		{location: xdg.UserDirs.Videos, name: "󰎁 Videos"},
		{location: xdg.UserDirs.Music, name: "♬ Music"},
		{location: xdg.UserDirs.Templates, name: "󰏢 Templates"},
		{location: xdg.UserDirs.PublicShare, name: " PublicShare"},
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

	jsonData, err := os.ReadFile(varibale.PinnedFilea)
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
