package metadata

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/yorukot/superfile/src/pkg/cache"
	"github.com/yorukot/superfile/src/pkg/utils"
)

type dirSizeEntry struct {
	size      int64
	modTime   time.Time
	fileCount int64
}

var directorySizeMutex sync.RWMutex

var directorySizeCache = cache.New[dirSizeEntry](
	defaultCacheSize,
	defaultCacheExpiration,
)

var directorySizeGroup singleflight.Group

func getDirectorySize(path string) int64 {
	fmt.Fprintln(os.Stderr, "GET DIR SIZE CALLED:", path)
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
	directorySizeMutex.RLock()
	cached, ok := directorySizeCache.Get(path)
	directorySizeMutex.RUnlock()

	if ok && cached.modTime.Equal(currentModTime) {
		slog.Info(
			"directory size cache hit",
			"path", path,
			"size", cached.size,
			"files", cached.fileCount,
		)

		return cached.size
	}

	result, err, _ := directorySizeGroup.Do(path, func() (any, error) {

		// Check again after singleflight wait
		directorySizeMutex.RLock()
		cached, ok := directorySizeCache.Get(path)
		directorySizeMutex.RUnlock()

		if ok && cached.modTime.Equal(currentModTime) {
			return cached.size, nil
		}

		slog.Info(
			"directory size calculating",
			"path", path,
		)

		stats := utils.GetDirStats(path)

		directorySizeMutex.Lock()

		directorySizeCache.Set(path, dirSizeEntry{
			size:      stats.Size,
			modTime:   currentModTime,
			fileCount: stats.FileCount,
		})

		directorySizeMutex.Unlock()

		slog.Info(
			"directory size calculated",
			"path", path,
			"size", stats.Size,
			"files", stats.FileCount,
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