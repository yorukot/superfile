package sidebar

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/adrg/xdg"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/utils"
)

// Return all sidebar directories
func getDirectories() []directory {
	return formDirctorySlice(getWellKnownDirectories(), getPinnedDirectories(), getExternalMediaFolders())
}

func formDirctorySlice(homeDirectories []directory, pinnedDirectories []directory,
	diskDirectories []directory) []directory {
	// Preallocation for efficiency
	totalCapacity := len(homeDirectories) + len(pinnedDirectories) + len(diskDirectories) + 2
	directories := make([]directory, 0, totalCapacity)

	directories = append(directories, homeDirectories...)
	directories = append(directories, pinnedDividerDir)
	directories = append(directories, pinnedDirectories...)
	directories = append(directories, diskDividerDir)
	directories = append(directories, diskDirectories...)
	return directories
}

// Return system default directory e.g. Home, Downloads, etc
func getWellKnownDirectories() []directory {
	wellKnownDirectories := []directory{
		{Location: xdg.Home, Name: icon.Home + icon.Space + "Home"},
		{Location: xdg.UserDirs.Download, Name: icon.Download + icon.Space + "Downloads"},
		{Location: xdg.UserDirs.Documents, Name: icon.Documents + icon.Space + "Documents"},
		{Location: xdg.UserDirs.Pictures, Name: icon.Pictures + icon.Space + "Pictures"},
		{Location: xdg.UserDirs.Videos, Name: icon.Videos + icon.Space + "Videos"},
		{Location: xdg.UserDirs.Music, Name: icon.Music + icon.Space + "Music"},
		{Location: xdg.UserDirs.Templates, Name: icon.Templates + icon.Space + "Templates"},
		{Location: xdg.UserDirs.PublicShare, Name: icon.PublicShare + icon.Space + "PublicShare"},
	}

	return slices.DeleteFunc(wellKnownDirectories, func(d directory) bool {
		_, err := os.Stat(d.Location)
		return err != nil
	})
}

// Get user pinned directories
func getPinnedDirectories() []directory {
	directories := []directory{}
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
			directories = append(directories, directory{Location: path, Name: directoryName})
		}
	} else {
		// Check if the data is in the new format
		if err := json.Unmarshal(jsonData, &directories); err != nil {
			// If the data is in neither format, log the error
			slog.Error("Error parsing pinned data", "error", err)
		}
	}

	clean := removeNotExistingDirectories(directories)

	return clean
}

// removeNotExistingDirectories removes directories that do not exist from the pinned directories list
func removeNotExistingDirectories(dirs []directory) []directory {
	cleanedDirs := make([]directory, 0, len(dirs))
	for _, dir := range dirs {
		if _, err := os.Stat(dir.Location); err == nil {
			cleanedDirs = append(cleanedDirs, dir)
		} else if !os.IsNotExist(err) {
			slog.Warn("Error while checking pinned directory", "directory", dir.Location, "error", err)
		}
	}

	// if any directory is removed, update the pinned file
	if len(cleanedDirs) != len(dirs) {
		if err := savePinnedDirectories(dirs); err != nil {
			slog.Error("error saving pinned directories", "error", err)
		}
	}

	return cleanedDirs
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
		haystack[i] = dir.Name
		dirMap[dir.Name] = dir
	}

	for _, match := range utils.FzfSearch(query, haystack) {
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

// TogglePinnedDirectory adds or removes a directory from the pinned directories list
func TogglePinnedDirectory(dir string) error {
	dirs := getPinnedDirectories()
	unPinned := false

	for i, other := range dirs {
		if other.Location == dir {
			dirs = append(dirs[:i], dirs[i+1:]...)
			unPinned = true
			break
		}
	}

	if !unPinned {
		dirs = append(dirs, directory{
			Location: dir,
			Name:     filepath.Base(dir),
		})
	}

	if err := savePinnedDirectories(dirs); err != nil {
		return fmt.Errorf("error saving pinned directories: %w", err)
	}

	return nil
}

// savePinnedDirectories marshals and writes the pinned directories to file.
func savePinnedDirectories(dirs []directory) error {
	data, err := json.Marshal(dirs)
	if err != nil {
		return fmt.Errorf("error marshaling pinned directories: %w", err)
	}

	if err := os.WriteFile(variable.PinnedFile, data, 0644); err != nil {
		return fmt.Errorf("error writing pinned directories file: %w", err)
	}

	return nil
}
