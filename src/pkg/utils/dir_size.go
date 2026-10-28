package utils

import (
	"context"
	"io/fs"
	"path/filepath"
)

// DirSizeWithContext returns the total size of all non-directory entries under path.
// partial is true if some entries inside path could not be read, in which case their
// sizes are not included. An error is returned if path itself cannot be read, or if
// ctx is cancelled before the walk completes.
func DirSizeWithContext(ctx context.Context, path string) (int64, bool, error) {
	var size int64
	partial := false
	walkErr := filepath.WalkDir(path, func(entryPath string, entry fs.DirEntry, err error) error {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if err != nil {
			if entryPath == path {
				return err
			}
			// Unreadable entries inside the directory are skipped
			partial = true
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			partial = true
			return nil //nolint:nilerr // Unreadable entries are skipped, and reported via partial
		}
		size += info.Size()
		return nil
	})
	if walkErr != nil {
		return 0, false, walkErr
	}
	return size, partial, nil
}
