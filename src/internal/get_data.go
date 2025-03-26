package internal

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/adrg/xdg"
	"github.com/shirou/gopsutil/v4/disk"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

// Return all sidebar directories
func getDirectories() []directory {
	return formDirctorySlice(getWellKnownDirectories(), getPinnedDirectories(), getExternalMediaFolders())
}

func formDirctorySlice(homeDirectories []directory, pinnedDirectories []directory, diskDirectories []directory) []directory {
	directories := append(homeDirectories, pinnedDividerDir)
	directories = append(directories, pinnedDirectories...)
	directories = append(directories, diskDividerDir)
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
		slog.Error("Error while read superfile data", "error", err)
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
		slog.Error("Error parsing pinned data", "error", err)
	}
	return directories
}

// Get external media directories
func getExternalMediaFolders() (disks []directory) {
	// only get physical drives
	parts, err := disk.Partitions(false)

	slog.Debug("disk.Partitions() called", "number of parts", len(parts))

	if err != nil {
		slog.Error("Error while getting external media: ", "error", err)
		return disks
	}

	for _, disk := range parts {
		// Not removing this for now, as it is helpful if we need to debug
		// any user issues related disk listing in sidebar.
		// Todo : We need to evaluate if more debug logs are a performance problem
		// even when user had set debug=false in config. We dont write those to log file
		// But we do make a functions call, and pass around some strings. So it might/might not be
		// a problem. It could be a problem in a hot path though.
		slog.Debug("Returned disk by disk.Partition()", "device", disk.Device,
			"mountpoint", disk.Mountpoint, "fstype", disk.Fstype)

		// shouldListDisk, diskName, and diskLocation, each has runtime.GOOS checks
		// We can ideally reduce it to one check only.
		if shouldListDisk(disk.Mountpoint) {

			disks = append(disks, directory{
				name:     diskName(disk.Mountpoint),
				location: diskLocation(disk.Mountpoint),
			})
		}
	}
	return disks
}

// Fuzzy search function for a list of directories.
func fuzzySearch(query string, dirs []directory) []directory {
	if len(dirs) == 0 {
		return []directory{}
	}

	var filteredDirs []directory

	// Optimization - This haystack can be kept precomputed based on directories
	// instead of re computing it in each call
	haystack := make([]string, len(dirs))
	dirMap := make(map[string]directory, len(dirs))
	for i, dir := range dirs {
		haystack[i] = dir.name
		dirMap[dir.name] = dir
	}

	for _, match := range fzfSearch(query, haystack) {
		if d, ok := dirMap[match.Key]; ok {
			filteredDirs = append(filteredDirs, d)
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
