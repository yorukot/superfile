package clipboard

import (
	"log/slog"
	"os"
	"slices"
	"strconv"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui"
)

const localClipboardSession = "local"

// The fact that its visible in UI or not, is controlled by the main model
type Model struct {
	width  int
	height int
	items  copyItems
}

// Copied items
type copyItems struct {
	items     []string
	locations []filesystem.Location
	cut       bool
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
				if i < len(m.items.locations) && m.items.locations[i].Provider != filesystem.ProviderLocal {
					r.AddLines(common.ClipboardPrettierName(m.items.items[i], viewWidth, false, false, false))
					continue
				}
				// TODO: Avoid Lstat during render for performance
				// Add IsDir/IsLink information in the item type or
				// better use filepanel's Element strcut as-is
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
	m.items.locations = m.items.locations[:0]
}

func (m *Model) Add(location string) {
	m.items.items = append(m.items.items, location)
	m.items.locations = append(m.items.locations, filesystem.Location{
		Provider:  filesystem.ProviderLocal,
		SessionID: localClipboardSession,
		Path:      filesystem.NewLocalPath(location),
		Label:     localClipboardSession,
	})
}

func (m *Model) AddLocation(location filesystem.Location) {
	m.items.items = append(m.items.items, location.Path.String())
	m.items.locations = append(m.items.locations, location)
}

func (m *Model) SetItems(items []string) {
	m.items.items = make([]string, len(items))
	copy(m.items.items, items)
	m.items.locations = make([]filesystem.Location, len(items))
	for i, item := range items {
		m.items.locations[i] = filesystem.Location{
			Provider:  filesystem.ProviderLocal,
			SessionID: localClipboardSession,
			Path:      filesystem.NewLocalPath(item),
			Label:     localClipboardSession,
		}
	}
}

func (m *Model) SetLocations(locations []filesystem.Location) {
	m.items.locations = make([]filesystem.Location, len(locations))
	copy(m.items.locations, locations)
	m.items.items = make([]string, len(locations))
	for i, location := range locations {
		m.items.items[i] = location.Path.String()
	}
}

func (m *Model) pruneInaccessibleItems() {
	if len(m.items.locations) > 0 {
		m.items.locations = slices.DeleteFunc(m.items.locations, func(item filesystem.Location) bool {
			if item.Provider != filesystem.ProviderLocal {
				return false
			}
			_, err := os.Lstat(item.Path.String())
			return err != nil
		})
		m.items.items = make([]string, len(m.items.locations))
		for i, item := range m.items.locations {
			m.items.items[i] = item.Path.String()
		}
		return
	}
	m.items.items = slices.DeleteFunc(m.items.items, func(item string) bool {
		_, err := os.Lstat(item)
		return err != nil
	})
}

func (m *Model) GetItems() []string {
	// return a copy to prevent external mutation
	items := make([]string, len(m.items.items))
	copy(items, m.items.items)
	return items
}

func (m *Model) GetLocations() []filesystem.Location {
	locations := make([]filesystem.Location, len(m.items.locations))
	copy(locations, m.items.locations)
	return locations
}

// Use this to use a copy that is in sync with current state of filesystem
func (m *Model) PruneInaccessibleItemsAndGet() []string {
	// Clipboard items might becomes outdated with
	// externally/interally triggered changes
	m.pruneInaccessibleItems()
	return m.GetItems()
}

func (m *Model) PruneInaccessibleLocationsAndGet() []filesystem.Location {
	m.pruneInaccessibleItems()
	return m.GetLocations()
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
