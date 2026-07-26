//go:build linux

package metadata

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func getFreeSpace(path string) (*DriveSize, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}
	blockSize := stat.Bsize
	if blockSize < 0 {
		return nil, fmt.Errorf("invalid disk: %s size:%d", path, stat.Bsize)
	}
	freeBytes := stat.Bfree * uint64(blockSize)
	totalBytes := stat.Blocks * uint64(blockSize)
	return &DriveSize{Total: totalBytes, Free: freeBytes}, nil
}
