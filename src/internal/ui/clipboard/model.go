package clipboard

import (
	"log/slog"
	"os"
	"strconv"

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
	} else {
		for i := 0; i < len(m.items.items) && i < viewHeight; i++ {
			if i == viewHeight-1 && i != len(m.items.items)-1 {
				// Last Entry we can render, but there are more that one left
				r.AddLines(strconv.Itoa(len(m.items.items)-i) + " items left....")
			} else {
				fileInfo, err := os.Lstat(m.items.items[i])
				if err != nil {
					slog.Error("Clipboard render function get item state ", "error", err)
					continue
				}
				isLink := fileInfo.Mode()&os.ModeSymlink != 0
				r.AddLines(common.ClipboardPrettierName(m.items.items[i],
					viewWidth, fileInfo.IsDir(), isLink, false))
			}
		}
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
	items := make([]string, len(m.items.items))
	copy(items, m.items.items)
	return items
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
