package metadata

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/yorukot/superfile/src/pkg/cache"
	"github.com/yorukot/superfile/src/pkg/utils"
)

type dirSizeEntry struct {
	size    int64
	modTime time.Time
}

var directorySizeMutex sync.RWMutex

var directorySizeCache = cache.New[dirSizeEntry](
	defaultCacheSize,
	defaultCacheExpiration,
)

var directorySizeGroup singleflight.Group

func getDirectorySize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		slog.Error(
			"failed to stat directory",
			"path", path,
			"error", err,
		)
		return 0
	}

	currentModTime := info.ModTime()

	// Fast cache lookup
	// Note: We check top-level directory ModTime as an instant freshness check.
	// Recursive changes inside subdirectories (or overwriting files without top-level mtime updates)
	// rely on directorySizeCache's defaultCacheExpiration TTL to bound staleness without expensive full-tree walks.
	directorySizeMutex.RLock()
	cached, ok := directorySizeCache.Get(path)
	directorySizeMutex.RUnlock()

	if ok && cached.modTime.Equal(currentModTime) {
		slog.Debug(
			"directory size cache hit",
			"path", path,
			"size", cached.size,
		)

		return cached.size
	}

	singleflightKey := path + ":" + currentModTime.Format(time.RFC3339Nano)
	result, err, _ := directorySizeGroup.Do(singleflightKey, func() (any, error) {

		// Check again after singleflight wait
		directorySizeMutex.RLock()
		cached, ok := directorySizeCache.Get(path)
		directorySizeMutex.RUnlock()

		if ok && cached.modTime.Equal(currentModTime) {
			return cached.size, nil
		}

		slog.Debug(
			"directory size calculating",
			"path", path,
		)

		stats, err := utils.GetDirStats(path)
		if err != nil {
			slog.Error(
				"directory size calculation failed",
				"path", path,
				"error", err,
			)
			return int64(0), err
		}

		directorySizeMutex.Lock()

		directorySizeCache.Set(path, dirSizeEntry{
			size:    stats.Size,
			modTime: currentModTime,
		})

		directorySizeMutex.Unlock()

		slog.Debug(
			"directory size calculated",
			"path", path,
			"size", stats.Size,
		)

		return stats.Size, nil
	})

	if err != nil {
		slog.Error(
			"directory size calculation failed",
			"path", path,
			"error", err,
		)
		return 0
	}

	return result.(int64)
}
