package sidebar

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/yorukot/superfile/src/pkg/utils"
)

// Get external media directories
func getExternalMediaFolders() []directory {
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
		if shouldListDisk(disk.Mountpoint) {
			disks = append(disks, directory{
				Name:     diskName(disk.Mountpoint),
				Location: diskLocation(disk.Mountpoint),
			})
		}
	}
	return disks
}

func shouldListDisk(mountPoint string) bool {
	if runtime.GOOS == utils.OsWindows {
		// We need to get C:, D: drive etc in the list
		return true
	}

	// Should always list the main disk
	if mountPoint == "/" {
		return true
	}

	// TODO : make a configurable field in config.yaml
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
	// TODO : make a configurable field in config.yaml
	// disk_mounts = ["/mnt", "/media", "/run/media", "/Volumes"]
	// Only block devicies that are mounted on these or any subdirectory of these Mountpoints
	// Will be shown in disk sidebar
	return strings.HasPrefix(mountPoint, "/mnt") ||
		strings.HasPrefix(mountPoint, "/media") ||
		strings.HasPrefix(mountPoint, "/run/media") ||
		strings.HasPrefix(mountPoint, "/Volumes")
}

func diskName(mountPoint string) string {
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

func diskLocation(mountPoint string) string {
	// In windows if you are in "C:\some\path", "cd C:" will not cd to root of C: drive
	// but "cd C:\" will
	if runtime.GOOS == utils.OsWindows {
		return filepath.Join(mountPoint, "\\")
	}
	return mountPoint
}
