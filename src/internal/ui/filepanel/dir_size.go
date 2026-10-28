package filepanel

import (
	"context"
	"os"
	"time"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// StartDirSizeCalc marks the focused directory as being calculated and returns a
// function that computes its size. The returned function does filesystem IO and
// must be run asynchronously, e.g. via a tea.Cmd.
// It returns nil if the focused item is not a directory, is a symlink, or is already being calculated.
// Pressing it again on a calculated directory recalculates it.
func (m *Model) StartDirSizeCalc(reqID int) func() DirSizeResult {
	elem := m.GetFocusedItemPtr()
	if elem == nil || !elem.Directory || elem.Info == nil {
		return nil
	}
	// Symlinks to directories are skipped. Their size column shows the link itself,
	// and the link's Info cannot tell us when the target directory changes.
	if elem.Info.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	if entry, ok := m.dirSizes[elem.Location]; ok && entry.calculating {
		return nil
	}

	if m.dirSizes == nil {
		m.dirSizes = make(map[string]dirSizeEntry)
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.dirSizes[elem.Location] = dirSizeEntry{
		calculating: true,
		reqID:       reqID,
		cancel:      cancel,
		startTime:   time.Now(),
	}
	location := elem.Location
	return func() DirSizeResult {
		size, partial, err := utils.DirSizeWithContext(ctx, location)
		return DirSizeResult{
			Location: location,
			ReqID:    reqID,
			Size:     size,
			Partial:  partial,
			Err:      err,
		}
	}
}

// ApplyDirSizeResult stores the result of a calculation started by StartDirSizeCalc.
// It returns false if the result is outdated, e.g. calculation was cancelled or restarted.
func (m *Model) ApplyDirSizeResult(res DirSizeResult) bool {
	entry, ok := m.dirSizes[res.Location]
	if !ok || !entry.calculating || entry.reqID != res.ReqID {
		return false
	}
	entry.cancel()
	if res.Err != nil {
		delete(m.dirSizes, res.Location)
		return true
	}
	entry.calculating = false
	entry.cancel = nil
	entry.size = res.Size
	entry.partial = res.Partial
	m.dirSizes[res.Location] = entry
	return true
}

// CancelDirSizeCalcs cancels all in-progress calculations. Finished results are kept.
func (m *Model) CancelDirSizeCalcs() {
	for location, entry := range m.dirSizes {
		if entry.calculating {
			entry.cancel()
			delete(m.dirSizes, location)
		}
	}
}

// Returns the text for size column of a directory. Does no IO
func (m *Model) dirSizeText(elem Element) string {
	entry, ok := m.dirSizes[elem.Location]
	if !ok {
		return ""
	}
	if entry.calculating {
		return dirSizeCalculatingText
	}
	// Directory changed since calculation. Result is stale
	if elem.Info == nil || elem.Info.ModTime().After(entry.startTime) {
		return ""
	}
	sizeText := common.FormatFileSize(entry.size)
	if entry.partial {
		return dirSizePartialPrefix + sizeText
	}
	return sizeText
}
