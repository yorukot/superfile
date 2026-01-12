package clipboard

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
)

// The fact that its visible in UI or not, is controlled by the main model
type Model struct {
	width        int
	height       int
	items        copyItems
	needsCleanup bool
}

// Copied items
type copyItems struct {
	items []clipboardItem
	cut   bool
}

// clipboardItem stores the path plus cached metadata so render can avoid per-item stat calls.
type clipboardItem struct {
	path       string
	isDir      bool
	isLink     bool
	isSelected bool
}

func (m *Model) SetDimensions(width int, height int) {
	m.width = width
	m.height = height
}

func (m *Model) Render() string {
	r := ui.ClipboardRenderer(m.height, m.width)
	viewHeight := m.height - common.BorderPadding
	viewWidth := m.width - common.InnerPadding
	if len(m.items.items) == 0 {
		// TODO move this to a string
		r.AddLines("", common.ClipboardNoneText)
		return r.Render()
	}
	for i := 0; i < len(m.items.items) && i < viewHeight; i++ {
		if i == viewHeight-1 && i != len(m.items.items)-1 {
			r.AddLines(strconv.Itoa(len(m.items.items)-i) + " items left....")
			continue
		}

		item := m.items.items[i]
		r.AddLines(common.ClipboardPrettierName(
			item.path,
			viewWidth,
			item.isDir,
			item.isLink,
			item.isSelected,
		))
	}
	return r.Render()
}

func (m *Model) IsCut() bool {
	return m.items.cut
}

func (m *Model) Reset(cut bool) {
	m.items.cut = cut
	m.items.items = m.items.items[:0]
	m.needsCleanup = false
}

func (m *Model) Add(location string) {
	m.items.items = append(m.items.items, m.makeClipboardItem(location))
	m.needsCleanup = true
}

func (m *Model) SetItems(items []string) {
	m.items.items = make([]clipboardItem, 0, len(items))
	for _, path := range items {
		m.items.items = append(m.items.items, m.makeClipboardItem(path))
	}
	m.needsCleanup = true
}

func (m *Model) GetItems() []string {
	if m.needsCleanup {
		return m.CleanupAndGetItems()
	}
	// return a copy to prevent external mutation
	out := make([]string, len(m.items.items))
	for i, item := range m.items.items {
		out[i] = item.path
	}
	return out
}

func (m *Model) Len() int {
	return len(m.items.items)
}

func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) GetHeight() int {
	return m.height
}

func (m *Model) GetFirstItem() string {
	if len(m.items.items) == 0 {
		return ""
	}
	return m.items.items[0].path
}

func (m *Model) UpdatePath(oldpath string, newpath string) {
	if len(m.items.items) == 0 {
		return
	}
	oldpathClean := filepath.Clean(oldpath)
	newpathClean := filepath.Clean(newpath)

	for i, p := range m.items.items {
		cur := filepath.Clean(p.path)
		if cur == oldpathClean {
			m.items.items[i].path = newpathClean
			continue
		}
		if strings.HasPrefix(cur, oldpathClean+string(filepath.Separator)) {
			m.items.items[i].path = filepath.Join(
				newpathClean,
				strings.TrimPrefix(cur, oldpathClean+string(filepath.Separator)),
			)
		}
	}
	m.needsCleanup = true
}
// CleanupAndGetItems removes entries whose paths no longer exist and returns the remaining paths.
func (m *Model) CleanupAndGetItems() [] string {
	if len(m.items.items) == 0 {
		return nil
	}

	kept := m.items.items[:0]
	for _, item := range m.items.items {
		if _, err := os.Lstat(item.path); err != nil {
			continue
		}
		kept = append(kept, item)
	}
	m.items.items = kept
	m.needsCleanup = false

	if len(m.items.items) == 0 {
		return nil
	}

	out := make([]string, len(m.items.items))
	for i, item := range m.items.items {
		out[i] = item.path
	}
	return out
}

// makeClipboardItem builds a clipboardItem with cached metadata for render.
func (m *Model) makeClipboardItem(path string) clipboardItem {
	info, err := os.Lstat(path)
	if err != nil {
		return clipboardItem{path: path}
	}
	isLink := info.Mode()&os.ModeSymlink != 0
	return clipboardItem{
		path:       path,
		isDir:      info.IsDir(),
		isLink:     isLink,
		isSelected: false,
	}
}
