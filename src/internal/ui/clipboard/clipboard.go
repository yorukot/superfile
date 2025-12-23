package clipboard

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
)

// The fact that its visible in UI or not, is controlled by the main model
type ClipboardModel struct {
	width  int
	height int
	items  copyItems
}

// Copied items
type copyItems struct {
	items []string
	cut   bool
}

func (m *ClipboardModel) SetDimensions(width int, height int) {
	m.width = width
	m.height = height
}

func (m *ClipboardModel) Render() string {
	r := ui.ClipboardRenderer(m.height, m.width)
	viewHeight := m.height - common.BorderPadding
	viewWidth := m.width - common.InnerPadding
	if len(m.items.items) == 0 {
		// TODO move this to a string
		r.AddLines("", " "+icon.Error+"  No content in clipboard")
	} else {
		for i := 0; i < len(m.items.items) && i < viewHeight; i++ {
			if i == viewHeight-1 && i != len(m.items.items)-1 {
				// Last Entry we can render, but there are more that one left
				r.AddLines(strconv.Itoa(len(m.items.items)-i) + " item left....")
			} else {
				fileInfo, err := os.Lstat(m.items.items[i])
				if err != nil {
					slog.Error("Clipboard render function get item state ", "error", err)
				}
				if !os.IsNotExist(err) {
					isLink := fileInfo.Mode()&os.ModeSymlink != 0
					// TODO : There is an inconsistency in parameter that is being passed,
					// and its name in ClipboardPrettierName function
					r.AddLines(common.ClipboardPrettierName(m.items.items[i],
						viewWidth, fileInfo.IsDir(), isLink, false))
				}
			}
		}
	}
	return r.Render()
}

func (m *ClipboardModel) IsCut() bool {
	return m.items.cut
}

func (m *ClipboardModel) Reset(cut bool) {
	m.items.cut = cut
	m.items.items = m.items.items[:0]
}

func (m *ClipboardModel) Add(location string) {
	m.items.items = append(m.items.items, location)
}

func (m *ClipboardModel) SetItems(items []string) {
	m.items.items = items
}

func (m *ClipboardModel) GetItems() []string {
	// return a copy to prevent external mutation
	items := make([]string, len(m.items.items))
	copy(items, m.items.items)
	return items
}

func (m *ClipboardModel) Len() int {
	return len(m.items.items)
}

func (m *ClipboardModel) GetWidth() int {
	return m.width
}

func (m *ClipboardModel) GetHeight() int {
	return m.height
}

func (m *ClipboardModel) GetFirstItem() string {
	if len(m.items.items) == 0 {
		return ""
	}
	return m.items.items[0]
}
