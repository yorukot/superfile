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

// ElementsRequest is everything a directory listing depends on. It holds no
// reference to Model, so Read() is safe to run off the event loop, and it is
// comparable, which is how a panel tells whether a listing is the one it wants.
type ElementsRequest struct {
	Location       string
	Search         string
	SortKind       sortmodel.SortKind
	SortReversed   bool
	DisplayDotFile bool

	// generation is bumped by MarkStale. Without it a read that started before
	// superfile changed the directory would carry an identical request, and could
	// land afterwards and replace the listing read after the change.
	generation int
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

	folderElementMap := map[string]os.DirEntry{}
	fileAndDirectories := []string{}

	for _, item := range items {
		if !r.DisplayDotFile && strings.HasPrefix(item.Name(), ".") {
			continue
		}

		fileAndDirectories = append(fileAndDirectories, item.Name())
		folderElementMap[item.Name()] = item
	}
	if len(fileAndDirectories) == 0 {
		return nil
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
		generation:     m.generation,
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

// UpdateElementsIfNeeded brings the panel's elements in line with what it wants
// to display.
//
// A read the panel cannot render without - first load, directory change, new
// search text, changed sort - is done synchronously. The periodic re-read that
// only picks up outside changes is returned instead, for the caller to run off
// the event loop. Doing that one inline made every message wait on the
// filesystem, so on a network mount the cursor fell behind held keys and kept
// moving long after they were released.
func (m *Model) UpdateElementsIfNeeded(force bool, displayDotFile bool) *ElementsRequest {
	req := m.elementsRequest(displayDotFile)
	if force || m.loaded != req {
		m.ApplyElements(req, req.Read(), displayDotFile)
		return nil
	}
	if m.refreshPending || time.Since(m.lastTimeGetElement) < m.refreshInterval() {
		return nil
	}
	m.refreshPending = true
	return &req
}

// MarkStale forces the panel to re-read its listing synchronously on the next
// update, and invalidates any read already in flight. Use it when superfile
// itself changes a directory's contents, so the change is visible right away
// instead of at the next poll.
func (m *Model) MarkStale() {
	// Bumping the generation makes `loaded` differ from what the panel now wants,
	// which takes the synchronous path below, and makes a read dispatched before
	// this point fail the check in ApplyElements.
	m.generation++
}

// ApplyElements installs a listing, ignoring a result for a request the panel has
// since moved on from - a different directory, search or sort, or a generation
// bumped by MarkStale while the read was in flight.
func (m *Model) ApplyElements(req ElementsRequest, elements []Element, displayDotFile bool) {
	if req != m.elementsRequest(displayDotFile) {
		slog.Debug("Ignoring elements of a stale request", "reqLocation", req.Location,
			"location", m.Location, "reqGeneration", req.generation,
			"generation", m.generation)
		return
	}
	// Every request change goes through the synchronous path above, so this also
	// clears the flag for a pending refresh whose result will never be applied.
	m.refreshPending = false
	m.element = elements
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
