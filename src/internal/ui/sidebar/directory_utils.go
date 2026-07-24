package sidebar

import (
	"os"
	"runtime"
	"slices"

	"github.com/adrg/xdg"

	"github.com/yorukot/superfile/src/internal/common"

	"github.com/yorukot/superfile/src/pkg/utils"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
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
		{Location: xdg.Home, Icon: icon.Home, Name: "Home"},
		{Location: xdg.UserDirs.Desktop, Icon: icon.Desktop, Name: "Desktop"},
		{Location: xdg.UserDirs.Download, Icon: icon.Download, Name: "Downloads"},
		{Location: xdg.UserDirs.Documents, Icon: icon.Documents, Name: "Documents"},
		{Location: xdg.UserDirs.Pictures, Icon: icon.Pictures, Name: "Pictures"},
		{Location: xdg.UserDirs.Videos, Icon: icon.Videos, Name: "Videos"},
		{Location: xdg.UserDirs.Music, Icon: icon.Music, Name: "Music"},
		{Location: xdg.UserDirs.Templates, Icon: icon.Templates, Name: "Templates"},
		{Location: xdg.UserDirs.PublicShare, Icon: icon.PublicShare, Name: "PublicShare"},
	}

	// Add Trash directory for Linux only
	if runtime.GOOS == utils.OsLinux {
		wellKnownDirectories = append(wellKnownDirectories, directory{
			Location: variable.LinuxTrashDirectory,
			Icon:     icon.Trash,
			Name:     "Trash",
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
		dirs[i].Icon = common.GetDirectoryIcon(dirs[i].Location, dirs[i].Name, common.Config.Nerdfont)
	}
	return dirs
}

// getFilteredDirectories returns a list of directories that match the fuzzy search query across all sections.
func getFilteredDirectories(query string, pinnedMgr *PinnedManager, sections []string) []directory {
	return formDirctorySlice(
		fuzzySearch(query, getWellKnownDirectories()),
		fuzzySearch(query, getPinnedDirectoriesWithIcon(pinnedMgr)),
		fuzzySearch(query, getExternalMediaFolders()),
		sections,
	)
}

// formDirctorySlice assembles the final list of directories for the sidebar based on the configured sections.
// It ensures that dividers are only added between non-empty sections.
func formDirctorySlice(homeDirectories []directory, pinnedDirectories []directory,
	diskDirectories []directory, sections []string) []directory {
	totalCapacity := len(homeDirectories) + len(pinnedDirectories) + len(diskDirectories) + directoryCapacityForDividers
	directories := make([]directory, 0, totalCapacity)

	for _, section := range sections {
		switch section {
		case utils.SidebarSectionHome:
			directories = appendSection(directories, utils.SidebarSectionHome, homeDividerDir, homeDirectories)
		case utils.SidebarSectionPinned:
			directories = appendSection(directories, utils.SidebarSectionPinned, pinnedDividerDir, pinnedDirectories)
		case utils.SidebarSectionDisks:
			directories = appendSection(directories, utils.SidebarSectionDisks, diskDividerDir, diskDirectories)
		}
	}

	return directories
}

func appendSection(dirs []directory, sectionName string, divider directory, items []directory) []directory {
	if len(items) == 0 {
		return dirs
	}

	if len(dirs) > 0 {
		dirs = append(dirs, divider)
	}

	for i := range items {
		items[i].Section = sectionName
		dirs = append(dirs, items[i])
	}

	return dirs
}
