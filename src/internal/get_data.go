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
func getDirectories() []Directory {
	return formDirctorySlice(getWellKnownDirectories(), getPinnedDirectories(), getExternalMediaFolders())
}

func formDirctorySlice(homeDirectories []Directory, pinnedDirectories []Directory, diskDirectories []Directory) []Directory {
	// Preallocation for efficiency
	totalCapacity := len(homeDirectories) + len(pinnedDirectories) + len(diskDirectories) + 2
	directories := make([]Directory, 0, totalCapacity)

	directories = append(directories, homeDirectories...)
	directories = append(directories, PinnedDividerDir)
	directories = append(directories, pinnedDirectories...)
	directories = append(directories, DiskDividerDir)
	directories = append(directories, diskDirectories...)
	return directories
}

// Return system default directory e.g. Home, Downloads, etc
func getWellKnownDirectories() []Directory {
	wellKnownDirectories := []Directory{
		{Location: xdg.Home, Name: icon.Home + icon.Space + "Home"},
		{Location: xdg.UserDirs.Download, Name: icon.Download + icon.Space + "Downloads"},
		{Location: xdg.UserDirs.Documents, Name: icon.Documents + icon.Space + "Documents"},
		{Location: xdg.UserDirs.Pictures, Name: icon.Pictures + icon.Space + "Pictures"},
		{Location: xdg.UserDirs.Videos, Name: icon.Videos + icon.Space + "Videos"},
		{Location: xdg.UserDirs.Music, Name: icon.Music + icon.Space + "Music"},
		{Location: xdg.UserDirs.Templates, Name: icon.Templates + icon.Space + "Templates"},
		{Location: xdg.UserDirs.PublicShare, Name: icon.PublicShare + icon.Space + "PublicShare"},
	}

	return slices.DeleteFunc(wellKnownDirectories, func(d Directory) bool {
		_, err := os.Stat(d.Location)
		return err != nil
	})
}

// Get user pinned directories
func getPinnedDirectories() []Directory {
	directories := []Directory{}
	var paths []string

	jsonData, err := os.ReadFile(variable.PinnedFile)
	if err != nil {
		slog.Error("Error while read superfile data", "error", err)
		return directories
	}

	// Check if the data is in the old format
	// TODO: Remove this after a release 1.2.4
	if err := json.Unmarshal(jsonData, &paths); err == nil {
		for _, path := range paths {
			directoryName := filepath.Base(path)
			directories = append(directories, Directory{Location: path, Name: directoryName})
		}
	} else {
		// Check if the data is in the new format
		if err := json.Unmarshal(jsonData, &directories); err != nil {
			// If the data is in neither format, log the error
			slog.Error("Error parsing pinned data", "error", err)
		}
	}
	return directories
}

// Get external media directories
func getExternalMediaFolders() []Directory {
	// only get physical drives
	parts, err := disk.Partitions(false)

	if err != nil {
		slog.Error("Error while getting external media: ", "error", err)
		return nil
	}
	var disks []Directory
	for _, disk := range parts {
		// shouldListDisk, diskName, and diskLocation, each has runtime.GOOS checks
		// We can ideally reduce it to one check only.
		if shouldListDisk(disk.Mountpoint) {
			disks = append(disks, Directory{
				Name:     diskName(disk.Mountpoint),
				Location: diskLocation(disk.Mountpoint),
			})
		}
	}
	return disks
}

// Fuzzy search function for a list of directories.
func fuzzySearch(query string, dirs []Directory) []Directory {
	if len(dirs) == 0 {
		return []Directory{}
	}

	var filteredDirs []Directory

	// Optimization - This haystack can be kept precomputed based on directories
	// instead of re computing it in each call
	haystack := make([]string, len(dirs))
	dirMap := make(map[string]Directory, len(dirs))
	for i, dir := range dirs {
		haystack[i] = dir.Name
		dirMap[dir.Name] = dir
	}

	for _, match := range fzfSearch(query, haystack) {
		if d, ok := dirMap[match.Key]; ok {
			filteredDirs = append(filteredDirs, d)
		}
	}

	return filteredDirs
}

// Get filtered directories using fuzzy search logic with three haystacks.
func getFilteredDirectories(query string) []Directory {
	return formDirctorySlice(
		fuzzySearch(query, getWellKnownDirectories()),
		fuzzySearch(query, getPinnedDirectories()),
		fuzzySearch(query, getExternalMediaFolders()),
	)
}
