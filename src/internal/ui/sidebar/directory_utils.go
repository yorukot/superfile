package sidebar

import (
	"os"
	"runtime"
	"slices"

	"github.com/adrg/xdg"

	"github.com/yorukot/superfile/src/pkg/utils"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

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

// getDirectories returns the list of directories to display in the sidebar.
func getDirectories(pinnedMgr *PinnedManager, sections []string) []directory {
	return formDirctorySlice(
		getWellKnownDirectories(),
		getPinnedDirectoriesWithIcon(pinnedMgr),
		getExternalMediaFolders(),
		sections,
	)
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

	// Add Trash directory for Linux only
	if runtime.GOOS == utils.OsLinux {
		wellKnownDirectories = append(wellKnownDirectories, directory{
			Location: variable.LinuxTrashDirectory,
			Name:     icon.Trash + icon.Space + "Trash",
		})
	}

	return slices.DeleteFunc(wellKnownDirectories, func(d directory) bool {
		_, err := os.Stat(d.Location)
		return err != nil
	})
}

func getPinnedDirectoriesWithIcon(pinnedMgr *PinnedManager) []directory {
	dirs := pinnedMgr.Load()
	for i := range dirs {
		iconInfo := common.GetElementIcon(dirs[i].Name, true, false, common.Config.Nerdfont)
		dirs[i].Name = iconInfo.Icon + icon.Space + dirs[i].Name
	}
	return dirs
}

// Get filtered directories using fuzzy search logic with three haystacks.
func getFilteredDirectories(query string, pinnedMgr *PinnedManager, sections []string) []directory {
	return formDirctorySlice(
		fuzzySearch(query, getWellKnownDirectories()),
		fuzzySearch(query, getPinnedDirectoriesWithIcon(pinnedMgr)),
		fuzzySearch(query, getExternalMediaFolders()),
		sections,
	)
}

func formDirctorySlice(homeDirectories []directory, pinnedDirectories []directory,
	diskDirectories []directory, sections []string) []directory {
	// Preallocation for efficiency
	totalCapacity := len(homeDirectories) + len(pinnedDirectories) + len(diskDirectories) + directoryCapacityForDividers
	directories := make([]directory, 0, totalCapacity)

	for _, section := range sections {
		switch section {
		case utils.SidebarSectionHome:
			if len(directories) > 0 {
				directories = append(directories, homeDividerDir)
			}
			directories = append(directories, homeDirectories...)
		case utils.SidebarSectionPinned:
			directories = append(directories, pinnedDividerDir)
			directories = append(directories, pinnedDirectories...)
		case utils.SidebarSectionDisks:
			directories = append(directories, diskDividerDir)
			directories = append(directories, diskDirectories...)
		}
	}

	return directories
}
