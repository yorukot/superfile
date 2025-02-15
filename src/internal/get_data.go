package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"

	"github.com/adrg/xdg"
	"github.com/reinhrst/fzf-lib"
	"github.com/shirou/gopsutil/disk"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

// Return all sidebar directories
func getDirectories() []directory {
	return formDirctorySlice(getWellKnownDirectories(), getPinnedDirectories(), getExternalMediaFolders())
}

func formDirctorySlice(homeDirectories []directory, pinnedDirectories []directory, diskDirectories []directory) []directory {
	directories := append(homeDirectories, pinnedDivider)
	directories = append(directories, pinnedDirectories...)
	directories = append(directories, diskDivider)
	directories = append(directories, diskDirectories...)
	return directories
}

// Return system default directory e.g. Home, Downloads, etc
func getWellKnownDirectories() []directory {
	wellKnownDirectories := []directory{
		{location: xdg.Home, name: icon.Home + icon.Space + "Home"},
		{location: xdg.UserDirs.Download, name: icon.Download + icon.Space + "Downloads"},
		{location: xdg.UserDirs.Documents, name: icon.Documents + icon.Space + "Documents"},
		{location: xdg.UserDirs.Pictures, name: icon.Pictures + icon.Space + "Pictures"},
		{location: xdg.UserDirs.Videos, name: icon.Videos + icon.Space + "Videos"},
		{location: xdg.UserDirs.Music, name: icon.Music + icon.Space + "Music"},
		{location: xdg.UserDirs.Templates, name: icon.Templates + icon.Space + "Templates"},
		{location: xdg.UserDirs.PublicShare, name: icon.PublicShare + icon.Space + "PublicShare"},
	}

	return slices.DeleteFunc(wellKnownDirectories, func(d directory) bool {
		_, err := os.Stat(d.location)
		return err != nil
	})
}

// Get user pinned directories
func getPinnedDirectories() []directory {
	directories := []directory{}
	var paths []string
	var pinnedDirs []struct {
		Location string `json:"location"`
		Name     string `json:"name"`
	}

	jsonData, err := os.ReadFile(variable.PinnedFile)
	if err != nil {
		outPutLog("Read superfile data error", err)
		return directories
	}

	// Check if the data is in the old format
	if err := json.Unmarshal(jsonData, &paths); err == nil {
		for _, path := range paths {
			directoryName := filepath.Base(path)
			directories = append(directories, directory{location: path, name: directoryName})
		}
		// Check if the data is in the new format
	} else if err := json.Unmarshal(jsonData, &pinnedDirs); err == nil {
		// Todo : we can optimize this. pinnedDirs and directories have exact same struct format
		// we are just copying data needlessly. We should directly unmarshal to 'directories' 
		for _, pinnedDir := range pinnedDirs {
			directories = append(directories, directory{location: pinnedDir.Location, name: pinnedDir.Name})
		}
		// If the data is in neither format, log the error
	} else {
		outPutLog("Error parsing pinned data", err)
	}
	return directories
}

// Get external media directories
func getExternalMediaFolders() (disks []directory) {
	parts, err := disk.Partitions(true)

	if err != nil {
		outPutLog("Error while getting external media: ", err)
		return disks
	}
	for _, disk := range parts {
		if isExternalDiskPath(disk.Mountpoint) {
			disks = append(disks, directory{
				name:     filepath.Base(disk.Mountpoint),
				location: disk.Mountpoint,
			})
		}
	}
	return disks
}

// Fuzzy search function for a list of directories.
func fuzzySearch(query string, dirs []directory) []directory {
	var filteredDirs []directory
	if len(dirs) > 0 {
		haystack := make([]string, len(dirs))
		dirMap := make(map[string]directory, len(dirs))

		for i, dir := range dirs {
			haystack[i] = dir.name
			dirMap[dir.name] = dir
		}

		options := fzf.DefaultOptions()
		fzfNone := fzf.New(haystack, options)
		fzfNone.Search(query)
		result := <-fzfNone.GetResultChannel()
		fzfNone.End()

		for _, match := range result.Matches {
			if d, ok := dirMap[match.Key]; ok {
				filteredDirs = append(filteredDirs, d)
			}
		}
	}
	return filteredDirs
}

// Get filtered directories using fuzzy search logic with three haystacks.
func getFilteredDirectories(query string) []directory {
	return formDirctorySlice(
		fuzzySearch(query, getWellKnownDirectories()),
		fuzzySearch(query, getPinnedDirectories()),
		fuzzySearch(query, getExternalMediaFolders()),
	)
}
