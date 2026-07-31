package filepanel

import (
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// ElementsRequest is everything a directory listing depends on. It is
// comparable, which is how a panel tells whether the listing it is showing is
// still the one it wants.
type ElementsRequest struct {
	Location       string
	Search         string
	SortKind       sortmodel.SortKind
	SortReversed   bool
	DisplayDotFile bool
}

// Read performs the filesystem IO for the request.
func (r ElementsRequest) Read() []Element {
	if r.Search != "" {
		return getDirectoryElementsBySearch(r)
	}
	return getDirectoryElements(r)
}

// TODO : Take common.Config.CaseSensitiveSort as a function parameter
// and also consider testing this caseSensitive with both true and false in
// our unit_test TestReturnDirElement
// getDirectoryElements returns the directory elements for the requested location
func getDirectoryElements(r ElementsRequest) []Element {
	dirEntries, err := os.ReadDir(r.Location)
	if err != nil {
		slog.Error("Error while returning folder elements", "error", err)
		return nil
	}

	if !r.DisplayDotFile {
		dirEntries = slices.DeleteFunc(dirEntries, func(e os.DirEntry) bool {
			return strings.HasPrefix(e.Name(), ".")
		})
	}

	// No files/directories to process
	if len(dirEntries) == 0 {
		return nil
	}
	// Entries whose Info() fails get dropped by sortFileElement
	return sortFileElement(r.SortKind, r.SortReversed, dirEntries, r.Location)
}

// getDirectoryElementsBySearch returns filtered directory elements based on search string
func getDirectoryElementsBySearch(r ElementsRequest) []Element {
	items, err := os.ReadDir(r.Location)
	if err != nil {
		slog.Error("Error while return folder element function", "error", err)
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	folderElementMap := map[string]os.DirEntry{}
	fileAndDirectories := []string{}

	for _, item := range items {
		if !r.DisplayDotFile && strings.HasPrefix(item.Name(), ".") {
			continue
		}

		fileAndDirectories = append(fileAndDirectories, item.Name())
		folderElementMap[item.Name()] = item
	}
	// https://github.com/reinhrst/fzf-lib/blob/main/core.go#L43
	// fzf returns matches ordered by score; we subsequently sort by the chosen sort option.
	fzfResults := utils.FzfSearch(r.Search, fileAndDirectories)
	dirElements := make([]os.DirEntry, 0, len(fzfResults))
	for _, item := range fzfResults {
		resultItem := folderElementMap[item.Key]
		dirElements = append(dirElements, resultItem)
	}

	return sortFileElement(r.SortKind, r.SortReversed, dirElements, r.Location)
}

// elementsRequest is the listing the panel wants to be displaying right now.
func (m *Model) elementsRequest(displayDotFile bool) ElementsRequest {
	return ElementsRequest{
		Location:       m.Location,
		Search:         m.SearchBar.Value(),
		SortKind:       m.SortKind,
		SortReversed:   m.SortReversed,
		DisplayDotFile: displayDotFile,
	}
}

// refreshInterval is how long a panel keeps a listing before re-reading it to
// pick up changes made outside superfile. Bigger directories cost more to read,
// so they get polled less often.
func (m *Model) refreshInterval() time.Duration {
	if !m.IsFocused {
		return nonFocussedPanelReRenderTime
	}
	scaled := time.Duration(m.ElemCount()/ReRenderChunkDivisor) * time.Second
	return min(max(scaled, focussedPanelReRenderTime), ReRenderMaxDelay*time.Second)
}

// UpdateElementsIfNeeded re-reads the panel's directory when it needs to.
//
// A read the panel cannot render without - first load, directory change, new
// search text, changed sort - happens right away. Otherwise the listing is only
// re-read once per refreshInterval, to pick up changes made outside superfile.
//
// The interval used to be int(elemCount / ReRenderChunkDivisor) seconds, which
// truncates to 0 for any directory with fewer than ReRenderChunkDivisor
// entries. The panel therefore re-read on every message, so every keystroke
// waited on the filesystem and navigating a network mount fell behind held keys.
func (m *Model) UpdateElementsIfNeeded(force bool, displayDotFile bool) {
	req := m.elementsRequest(displayDotFile)
	if !force && req == m.loaded && time.Since(m.lastTimeGetElement) < m.refreshInterval() {
		return
	}

	m.element = req.Read()
	m.loaded = req
	m.lastTimeGetElement = time.Now()

	// For hover to file on first time loading
	if m.TargetFile != "" {
		m.applyTargetFileCursor()
	}

	// If cursor becomes invalid due to element update, reset
	if m.ValidateCursorAndRenderIndex() != nil {
		m.scrollToCursor(0)
	}
}

// MarkStale forces the panel to re-read its listing on the next update. Use it
// when superfile itself changes a directory's contents, so the change is visible
// right away instead of at the next refresh.
func (m *Model) MarkStale() {
	m.loaded = ElementsRequest{}
}
