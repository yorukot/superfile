package sidebar

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// getExternalMediaFolders retrieves the list of mounted drives to display in the
// disks section, filtered according to the user's disk_mounts / excluded_disk_mounts
// configuration.
func getExternalMediaFolders() []directory {
	// List every mount (all=true) so that non-physical filesystems such as FUSE
	// (sshfs) and network mounts are available. Filtering to what the user wants
	// is then handled entirely by shouldListDisk.
	parts, err := disk.Partitions(true)

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
				Icon:     diskIcon(disk.Mountpoint),
				Name:     diskName(disk.Mountpoint),
				Location: diskLocation(disk.Mountpoint),
			})
		}
	}
	return disks
}

// shouldListDisk determines whether a given mount point should be displayed in the
// sidebar's disks section, based on the configured include/exclude mountpoint prefixes.
func shouldListDisk(mountPoint string) bool {
	return shouldListDiskWithConfig(mountPoint, common.Config.DiskMounts, common.Config.ExcludedDiskMounts)
}

// shouldListDiskWithConfig contains the mountpoint filtering logic, taking the
// include/exclude prefix lists explicitly so it can be unit tested without touching
// global config.
//
// Rules, in order:
//  1. Windows: always list (drives are enumerated; the POSIX prefixes don't apply).
//  2. The root mount "/" is always listed.
//  3. A mountpoint matching any excludePrefixes entry is hidden (blacklist wins).
//  4. An empty includePrefixes list lists every mountpoint (all filesystem types).
//  5. Otherwise the mountpoint is listed only if it matches an includePrefixes entry.
func shouldListDiskWithConfig(mountPoint string, includePrefixes, excludePrefixes []string) bool {
	if runtime.GOOS == utils.OsWindows {
		// We need to get C:, D: drive etc in the list
		return true
	}

	// Should always list the main disk
	if mountPoint == "/" {
		return true
	}

	// Blacklist takes precedence over the include list.
	if hasAnyPrefix(mountPoint, excludePrefixes) {
		return false
	}

	// An empty include list means "list everything" (all filesystem types).
	if len(includePrefixes) == 0 {
		return true
	}

	// Otherwise only list mounts under one of the included prefixes.
	return hasAnyPrefix(mountPoint, includePrefixes)
}

// hasAnyPrefix reports whether s starts with any of the given prefixes.
func hasAnyPrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func diskIcon(mountPoint string) string {
	if runtime.GOOS != utils.OsWindows && mountPoint == "/" {
		return icon.Terminal
	}

	return icon.Disk
}

// diskName generates a display name for a disk based on its mount point and the operating system.
func diskName(mountPoint string) string {
	// In windows we dont want to use filepath.Base as it returns "\" for when
	// mountPoint is any drive root "C:", "D:", etc. Hence causing same name
	// for each drive
	name := mountPoint
	if runtime.GOOS != utils.OsWindows {
		if mountPoint == "/" {
			name = "Root"
		} else {
			// This might cause duplicate names in case you mount two devices in
			// /mnt/usb and /mnt/dir2/usb . Full mountpoint is a more accurate way
			// but that results in messy UI, hence we do this.
			name = filepath.Base(mountPoint)
		}
	}

	return name
}

// diskLocation returns the normalized path for a disk's mount point.
func diskLocation(mountPoint string) string {
	// In windows if you are in "C:\some\path", "cd C:" will not cd to root of C: drive
	// but "cd C:\" will
	if runtime.GOOS == utils.OsWindows {
		return filepath.Join(mountPoint, "\\")
	}
	return mountPoint
}
