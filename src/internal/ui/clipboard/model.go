package clipboard

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
)

// The fact that its visible in UI or not, is controlled by the main model
type Model struct {
	width  int
	height int
	items  copyItems
}

// Copied items
type copyItems struct {
	items []string
	cut   bool
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
			r.AddLines(strconv.Itoa(len(m.items.items)-i) + " items left...")
			continue
		}
		r.AddLines(common.ClipboardPrettierName(
			m.items.items[i],
			viewWidth,
			false,
			false,
			false,
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
}

func (m *Model) Add(location string) {
	m.items.items = append(m.items.items, location)
}

func (m *Model) SetItems(items []string) {
	m.items.items = make([]string, len(items))
	copy(m.items.items, items)
}

func (m *Model) GetItems() []string {
	// return a copy to prevent external mutation
	return m.GetItemsWithExistCheck()
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
	return m.items.items[0]
}

func (m *Model) UpdatePath(oldpath string, newpath string) {
	if len(m.items.items) == 0 {
		return
	}
	oldpathClean := filepath.Clean(oldpath)
	newpathClean := filepath.Clean(newpath)

	for i, p := range m.items.items {
		cur := filepath.Clean(p)
		if cur == oldpathClean {
			m.items.items[i] = newpathClean
			continue
		}
		if strings.HasPrefix(cur, oldpathClean+string(filepath.Separator)) {
			m.items.items[i] = filepath.Join(
				newpathClean,
				strings.TrimPrefix(cur, oldpathClean+string(filepath.Separator)),
			)
		}

	}
}

func (m *Model) GetItemsWithExistCheck() []string {
	if len(m.items.items) == 0 {
		return nil
	}

	m.items.items = slices.DeleteFunc(m.items.items, func(p string) bool {
		_, err := os.Stat(p)
		return err != nil
	})
	out := make([]string, len(m.items.items))
	copy(out, m.items.items)
	return out
}
