package disks

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	bytesPerTB = 1_000_000_000_000
	bytesPerGB = 1_000_000_000
	bytesPerMB = 1_000_000
)

// linuxDiskFromObject builds a Disk from one udisks2 block device object.
// The bool is false when the object should not be shown (not a mountable
// block device, hidden, no media, malformed, or an in-use internal partition).
func diskFromObject(path dbus.ObjectPath,
	ifaces map[string]map[string]dbus.Variant,
) (Disk, bool) {
	block, isBlockDevice := ifaces["org.freedesktop.UDisks2.Block"]
	fs, hasFS := ifaces["org.freedesktop.UDisks2.Filesystem"]

	// Must be a block device with a Filesystem (i.e., mountable.
	// Filesystem should not be confused with a Partition)
	if !isBlockDevice || !hasFS {
		return Disk{}, false
	}

	// The disk should be hidden from the user, as per the udisks2 docs
	hintIgnore, ok := block["HintIgnore"].Value().(bool)
	if hintIgnore || !ok {
		return Disk{}, false
	}

	size, ok := block["Size"].Value().(uint64)
	if size == 0 || !ok {
		return Disk{}, false
	}

	// If true then it means the disk is an internal disk
	hintSystem, ok := block["HintSystem"].Value().(bool)
	if !ok {
		return Disk{}, false
	}

	mountPoints, ok := fs["MountPoints"].Value().([][]byte)
	if !ok {
		return Disk{}, false
	}

	mountPoint := diskMountPoint(mountPoints)

	// If a user has a partitioned internal disk, e.g., if they are
	// dual-booting Linux and Windows, then we want to show the partition
	// that is not running their current OS
	if hintSystem && len(mountPoints) > 0 {
		return Disk{}, false
	}

	idLabel, _ := block["IdLabel"].Value().(string) // blank IdLabel is fine

	return Disk{
		Name:       diskName(idLabel, size),
		MountPoint: mountPoint,
		Mounted:    len(mountPoints) > 0,
		Internal:   hintSystem,
		ObjectPath: string(path),
	}, true
}

// If a block device has a IdLabel we return that
// otherwise we fallback to using the disks size.
func diskName(idLabel string, size uint64) string {
	if idLabel != "" {
		return idLabel
	}

	switch {
	case size < bytesPerGB:
		return fmt.Sprintf("%d MB Volume", size/bytesPerMB)
	case size < bytesPerTB:
		return fmt.Sprintf("%d GB Volume", size/bytesPerGB)
	default:
		return fmt.Sprintf("%.1f TB Volume", float64(size)/bytesPerTB)
	}
}

// Returns the first mount point in mountPoints as a string,
// with the trailing NUL byte removed. If mountPoints is empty,
// it returns an empty string.
func diskMountPoint(mountPoints [][]byte) string {
	if len(mountPoints) == 0 {
		return ""
	}
	return strings.TrimRight(string(mountPoints[0]), "\x00")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
