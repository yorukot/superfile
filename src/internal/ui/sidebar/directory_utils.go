package sidebar

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/reinhrst/fzf-lib"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/adrg/xdg"
	"github.com/shirou/gopsutil/v4/disk"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

// Return all sidebar directories
func GetDirectories() []directory {
	return FormDirctorySlice(GetWellKnownDirectories(), GetPinnedDirectories(), GetExternalMediaFolders())
}

func FormDirctorySlice(homeDirectories []directory, pinnedDirectories []directory, diskDirectories []directory) []directory {
	// Preallocation for efficiency
	totalCapacity := len(homeDirectories) + len(pinnedDirectories) + len(diskDirectories) + 2
	directories := make([]directory, 0, totalCapacity)

	directories = append(directories, homeDirectories...)
	directories = append(directories, PinnedDividerDir)
	directories = append(directories, pinnedDirectories...)
	directories = append(directories, DiskDividerDir)
	directories = append(directories, diskDirectories...)
	return directories
}

// Return system default directory e.g. Home, Downloads, etc
func GetWellKnownDirectories() []directory {
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
func GetPinnedDirectories() []directory {
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
	return directories
}

// Get external media directories
func GetExternalMediaFolders() []directory {
	// only get physical drives
	parts, err := disk.Partitions(false)

	if err != nil {
		slog.Error("Error while getting external media: ", "error", err)
		return nil
	}
	var disks []directory
	for _, disk := range parts {
		// ShouldListDisk, DiskName, and DiskLocation, each has runtime.GOOS checks
		// We can ideally reduce it to one check only.
		if ShouldListDisk(disk.Mountpoint) {
			disks = append(disks, directory{
				Name:     DiskName(disk.Mountpoint),
				Location: DiskLocation(disk.Mountpoint),
			})
		}
	}
	return disks
}

func ShouldListDisk(mountPoint string) bool {
	if runtime.GOOS == utils.OsWindows {
		// We need to get C:, D: drive etc in the list
		return true
	}

	// Should always list the main disk
	if mountPoint == "/" {
		return true
	}

	// Todo : make a configurable field in config.yaml
	// excluded_disk_mounts = ["/Volumes/.timemachine"]
	// Mountpoints that are in subdirectory of disk_mounts
	// but still are to be excluded in disk section of sidebar
	if strings.HasPrefix(mountPoint, "/Volumes/.timemachine") {
		return false
	}

	// We avoid listing all mounted partitions (Otherwise listed disk could get huge)
	// but only a few partitions that usually corresponds to external physical devices
	// For example : mounts like /boot, /var/ will get skipped
	// This can be inaccurate based on your system setup if you mount any external devices
	// on other directories, or if you have some extra mounts on these directories
	// Todo : make a configurable field in config.yaml
	// disk_mounts = ["/mnt", "/media", "/run/media", "/Volumes"]
	// Only block devicies that are mounted on these or any subdirectory of these Mountpoints
	// Will be shown in disk sidebar
	return strings.HasPrefix(mountPoint, "/mnt") ||
		strings.HasPrefix(mountPoint, "/media") ||
		strings.HasPrefix(mountPoint, "/run/media") ||
		strings.HasPrefix(mountPoint, "/Volumes")
}

func DiskName(mountPoint string) string {
	// In windows we dont want to use filepath.Base as it returns "\" for when
	// mountPoint is any drive root "C:", "D:", etc. Hence causing same name
	// for each drive
	if runtime.GOOS == utils.OsWindows {
		return mountPoint
	}

	// This might cause duplicate names in case you mount two devices in
	// /mnt/usb and /mnt/dir2/usb . Full mountpoint is a more accurate way
	// but that results in messy UI, hence we do this.
	return filepath.Base(mountPoint)
}

func DiskLocation(mountPoint string) string {
	// In windows if you are in "C:\some\path", "cd C:" will not cd to root of C: drive
	// but "cd C:\" will
	if runtime.GOOS == utils.OsWindows {
		return filepath.Join(mountPoint, "\\")
	}
	return mountPoint
}

// Fuzzy search function for a list of directories.
func FuzzySearch(query string, dirs []directory) []directory {
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

	for _, match := range FzfSearch(query, haystack) {
		if d, ok := dirMap[match.Key]; ok {
			filteredDirs = append(filteredDirs, d)
		}
	}

	return filteredDirs
}

// Get filtered directories using fuzzy search logic with three haystacks.
func GetFilteredDirectories(query string) []directory {
	return FormDirctorySlice(
		FuzzySearch(query, GetWellKnownDirectories()),
		FuzzySearch(query, GetPinnedDirectories()),
		FuzzySearch(query, GetExternalMediaFolders()),
	)
}

// Returning a string slice causes inefficiency in current usage
func FzfSearch(query string, source []string) []fzf.MatchResult {
	fzfSearcher := fzf.New(source, fzf.DefaultOptions())
	fzfSearcher.Search(query)
	fzfResults := <-fzfSearcher.GetResultChannel()
	fzfSearcher.End()
	return fzfResults.Matches
}

// TogglePinnedDirectory adds or removes a directory from the pinned directories list
func TogglePinnedDirectory(dir string) error {
	dirs := GetPinnedDirectories()
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

	updatedData, err := json.Marshal(dirs)
	if err != nil {
		return fmt.Errorf("error marshaling pinned directories: %w", err)
	}

	err = os.WriteFile(variable.PinnedFile, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("error writing pinned directories file: %w", err)
	}

	return nil
}
