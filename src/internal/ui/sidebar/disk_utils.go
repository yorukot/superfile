package sidebar

import (
	"context"
	"log/slog"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/disks"
)

// getExternalMediaFolders retrieves a list of mounted physical drives and external media.
func getExternalMediaFolders(diskManager disks.Manager) []directory {
	disks, err := diskManager.List(context.Background())
	if err != nil {
		slog.Error("Error while getting external media: ", "error", err)
		return nil
	}

	var dirs []directory
	for _, disk := range disks {
		dirs = append(dirs, directory{
			Icon:     diskIcon(disk.Mounted),
			Name:     disk.Name,
			Location: disk.MountPoint,
		})
	}
	return dirs
}

func diskIcon(isMounted bool) string {
	if !isMounted {
		return " "
	}
	return icon.Disk
}
