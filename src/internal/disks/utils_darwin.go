package disks

import (
	"path/filepath"
	"strings"
)

func shouldListDisk(mountPoint string) bool {
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

	// TODO : make a configurable field in config.yaml
	// disk_mounts = ["/mnt", "/media", "/run/media", "/Volumes"]
	// Only block devicies that are mounted on these or any subdirectory of these Mountpoints
	// Will be shown in disk sidebar
	return strings.HasPrefix(mountPoint, "/Volumes")
}

func diskName(mountPoint string) string {
	name := mountPoint
	if mountPoint == "/" {
		name = "Root"
	} else {
		// This might cause duplicate names in case you mount two devices in
		// /mnt/usb and /mnt/dir2/usb . Full mountpoint is a more accurate way
		// but that results in messy UI, hence we do this.
		name = filepath.Base(mountPoint)
	}

	return name
}
