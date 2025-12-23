package internal

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
)

// The fact that its visible in UI or not, is controlled by the main model
type clipboardModel struct {
	Width     int
	Height    int
	CopyItems CopyItems
}

func (m *clipboardModel) SetDimensions(width int, height int) {
	m.Width = width
	m.Height = height
}

func (m *clipboardModel) Render() string {
	r := ui.ClipboardRenderer(m.Height, m.Width)
	viewHeight := m.Height - common.BorderPadding
	viewWidth := m.Width - common.InnerPadding
	if len(m.CopyItems.items) == 0 {
		// TODO move this to a string
		r.AddLines("", " "+icon.Error+"  No content in clipboard")
	} else {
		for i := 0; i < len(m.CopyItems.items) && i < viewHeight; i++ {
			if i == viewHeight-1 && i != len(m.CopyItems.items)-1 {
				// Last Entry we can render, but there are more that one left
				r.AddLines(strconv.Itoa(len(m.CopyItems.items)-i) + " item left....")
			} else {
				fileInfo, err := os.Lstat(m.CopyItems.items[i])
				if err != nil {
					slog.Error("Clipboard render function get item state ", "error", err)
				}
				if !os.IsNotExist(err) {
					isLink := fileInfo.Mode()&os.ModeSymlink != 0
					// TODO : There is an inconsistency in parameter that is being passed,
					// and its name in ClipboardPrettierName function
					r.AddLines(common.ClipboardPrettierName(m.CopyItems.items[i],
						viewWidth, fileInfo.IsDir(), isLink, false))
				}
			}
		}
	}
	return r.Render()
}

func (m *clipboardModel) IsCut() bool {
	return m.CopyItems.cut
}

func (m *clipboardModel) Reset(cut bool) {
	m.CopyItems.cut = cut
	m.CopyItems.items = m.CopyItems.items[:0]
}

func (m *clipboardModel) Add(location string) {
	m.CopyItems.items = append(m.CopyItems.items, location)
}

func (m *clipboardModel) SetItems(items []string) {
	m.CopyItems.items = items
}

func (m *clipboardModel) GetFirstItem() string {
	if len(m.CopyItems.items) == 0 {
		return ""
	}
	return m.CopyItems.items[0]
}
