package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/notify"
)

// dragStartThreshold is how far the pointer has to travel, in cells, before a
// press turns into a drag.
//
// A press alone is already a click, and a click that jitters by a cell while the
// button is down must still count as a click, so the drag waits for the pointer
// to clearly leave where it started.
const dragStartThreshold = 2

// dragState is a drag in progress, or the absence of one.
//
// The drag is only recorded between a press and the release that ends it. The
// pointer does not have to keep moving for it to stay active, so a drag can be
// held still over a target.
type dragState struct {
	// active is set once the pointer has travelled far enough. Until then a
	// press is still only a click.
	active bool
	// pending is set on a press that could still turn into a drag.
	pending bool

	pressX, pressY int

	// locations are the paths being dragged.
	locations []string
	// fromPanel is the panel the drag started in.
	fromPanel int

	// target is what the pointer is over now.
	target mouse.Target
	// destLocation is the directory the drag would move into, empty when the
	// pointer is not over somewhere valid.
	destLocation string
}

// noteLeftPress records a press that may become a drag.
//
// The click itself has already been handled by the time this is called. Only a
// press on a file entry can start a drag, since dragging anything else has no
// meaning.
func (m *model) noteLeftPress(target mouse.Target, x, y int) {
	m.drag = dragState{}

	if target.Kind != mouse.TargetFilePanelItem {
		return
	}

	panelIndex := target.PanelIndex
	if panelIndex < 0 || panelIndex >= len(m.fileModel.FilePanels) {
		return
	}
	panel := &m.fileModel.FilePanels[panelIndex]
	if target.ItemIndex < 0 || target.ItemIndex >= panel.ElemCount() {
		return
	}

	location := panel.GetElementAtIdx(target.ItemIndex).Location
	locations := []string{location}
	// Dragging from inside a selection drags the whole selection, the same way a
	// right click inside a selection keeps it.
	if panel.SelectedCount() > 0 && panel.CheckSelected(location) {
		locations = panel.GetSelectedLocationsSortedAsVisible()
	}
	if len(locations) == 0 {
		return
	}

	m.drag = dragState{
		pending:   true,
		pressX:    x,
		pressY:    y,
		locations: locations,
		fromPanel: panelIndex,
	}
}

// trackDragMotion updates the drop target while the pointer moves with a button
// held.
//
// It reports whether the drag has started, meaning the pointer has moved far
// enough from the press to be a drag rather than a click.
func (m *model) trackDragMotion(x, y int) bool {
	// Both a press that could still become a drag and a drag that already is one
	// need their target updated, so the guard is on there being no drag at all.
	if !m.drag.pending && !m.drag.active {
		return false
	}

	if !m.drag.active {
		if abs(x-m.drag.pressX)+abs(y-m.drag.pressY) < dragStartThreshold {
			return false
		}
		m.drag.active = true
		m.drag.pending = false
		// A drag must not leave a half recorded click behind, or the release
		// would look like the first half of a double click and would open the
		// file that was dragged.
		m.lastLeftClick = leftClick{}
	}

	m.updateDragAt(x, y)
	return true
}

// updateDragAt resolves what the drag would drop onto, given a terminal position.
//
// Both motion and the release call it, because the release is the one that
// decides where the files actually go and a terminal is allowed to skip motion
// events. Resolving there too means a release over a directory other than the
// last one seen still drops into the directory the pointer is really over, and a
// release over nowhere clears the destination instead of remembering an old one.
func (m *model) updateDragAt(x, y int) {
	m.drag.target = m.mouseTargetAt(x, y)
	m.drag.destLocation = m.dropLocationFor(m.drag.target, x, y)
	m.setDropTargets()
}

// dropLocationFor resolves what the drag would move into, or an empty string
// when the pointer is not over somewhere that can be dropped on.
func (m *model) dropLocationFor(target mouse.Target, x, y int) string {
	var dest string

	switch target.Kind {
	case mouse.TargetFilePanelItem:
		panelIndex := target.PanelIndex
		if panelIndex < 0 || panelIndex >= len(m.fileModel.FilePanels) {
			return ""
		}
		panel := &m.fileModel.FilePanels[panelIndex]
		if target.ItemIndex < 0 || target.ItemIndex >= panel.ElemCount() {
			return ""
		}
		element := panel.GetElementAtIdx(target.ItemIndex)
		// A file is not somewhere anything can be moved into.
		if !element.Directory {
			return ""
		}
		dest = element.Location
	case mouse.TargetSidebarDirectory:
		dest = m.sidebarModel.GetDirectoryLocation(target.ItemIndex)
	case mouse.TargetUnknown:
		// The area around the panels means the directory the panel is showing.
		if !m.inFilePanelArea(x, y) {
			return ""
		}
		panelIndex := m.filePanelIndexAtX(x)
		if panelIndex < 0 || panelIndex >= len(m.fileModel.FilePanels) {
			return ""
		}
		dest = m.fileModel.FilePanels[panelIndex].Location
	case mouse.TargetProcessBarItem, mouse.TargetMetadataItem, mouse.TargetContextMenuItem:
		// None of these are somewhere a file can be moved into.
		return ""
	}

	if dest == "" {
		return ""
	}
	if m.dragIsInvalidFor(dest) {
		return ""
	}
	return dest
}

// dragIsInvalidFor reports whether moving the dragged paths into dest cannot be
// done, so that an impossible drop is shown as no drop at all rather than as an
// error after the fact.
func (m *model) dragIsInvalidFor(dest string) bool {
	info, err := os.Stat(dest)
	if err != nil || !info.IsDir() {
		return true
	}

	for _, src := range m.drag.locations {
		if filepath.Dir(src) == dest {
			// Already in there, so there is nothing to do.
			return true
		}
		// Moving a directory into itself or into something below it would
		// destroy the directory being moved.
		if isAncestor(src, dest) {
			return true
		}
	}
	return false
}

// filePanelIndexAtX returns the panel whose border covers a terminal column, or
// -1 for a column that belongs to none of them.
//
// The column arrives in terminal coordinates, while panel origins are measured
// from the start of the panel area, so the sidebar's width comes off first. The
// full panel width is used rather than the content width, because the regions
// that resolve a pointer to a panel cover the border columns too.
func (m *model) filePanelIndexAtX(x int) int {
	x -= m.filePanelAreaX()
	if x < 0 {
		return -1
	}
	for i := range m.fileModel.FilePanels {
		origin := m.fileModel.PanelOriginX(i)
		if x >= origin && x < origin+m.fileModel.FilePanels[i].GetWidth() {
			return i
		}
	}
	return -1
}

// setDropTargets marks the row the pointer is over, so it can be drawn as the
// place the drop would land.
func (m *model) setDropTargets() {
	clearDropTargets := func() {
		for i := range m.fileModel.FilePanels {
			m.fileModel.FilePanels[i].SetDropTarget(filepanel.NoDropTarget)
		}
	}

	if !m.drag.active || m.drag.destLocation == "" {
		clearDropTargets()
		return
	}

	// Only a row can show where the drop lands. Dropping on the area around the
	// panels lands in the directory that area already shows, so there is nothing
	// to point at.
	if m.drag.target.Kind != mouse.TargetFilePanelItem {
		clearDropTargets()
		return
	}

	clearDropTargets()
	m.fileModel.FilePanels[m.drag.target.PanelIndex].
		SetDropTarget(m.drag.target.ItemIndex)
}

// finishDrag completes a release, moving the dragged paths when there is a valid
// destination.
//
// It always clears the drag, so a release that is not a drop leaves nothing
// behind.
func (m *model) finishDrag() tea.Cmd {
	drag := m.drag
	m.clearDrag()

	if !drag.active {
		// The pointer never left the press, so this was an ordinary click that
		// has already been handled.
		return nil
	}

	dest := drag.destLocation
	if dest == "" {
		slog.Debug("drag ended without a valid destination")
		return nil
	}

	slog.Debug("dropping dragged items", "count", len(drag.locations), "dest", dest)
	return m.getMoveItemsCmd(dest, drag.locations)
}

// clearDrag forgets any drag and takes the drop highlight back off the panels.
func (m *model) clearDrag() {
	m.drag = dragState{}
	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].SetDropTarget(filepanel.NoDropTarget)
	}
}

// handleDragKeyInput reports whether a key was consumed by a drag in progress.
//
// Only the cancel key is. Any other key is left to the handlers below, so that a
// release the terminal never sent cannot leave the app ignoring the keyboard.
func (m *model) handleDragKeyInput(key string) bool {
	if !m.drag.pending && !m.drag.active {
		return false
	}
	if !slices.Contains(common.Hotkeys.CancelTyping, key) {
		return false
	}
	m.cancelDrag()
	return true
}

// cancelDrag abandons a drag in progress, which is what escape does.
func (m *model) cancelDrag() bool {
	if !m.drag.pending && !m.drag.active {
		return false
	}
	m.clearDrag()
	return true
}

// getMoveItemsCmd moves paths into a directory.
//
// It goes through the same processor the cut and paste flow uses, so a move
// reports progress, refuses the same impossible destinations, and does not touch
// the clipboard the user may be relying on.
func (m *model) getMoveItemsCmd(dest string, locations []string) tea.Cmd {
	if len(locations) == 0 || dest == "" {
		return nil
	}

	reqID := m.nextIoReqCnt()
	slog.Debug("Submitting moveItems request", "id", reqID,
		"items cnt", len(locations), "dest", dest)

	return func() tea.Msg {
		err := validatePasteOperation(dest, locations, true)
		if err != nil {
			return NewNotifyModalMsg(
				notify.New(true, "Invalid move location", err.Error(), notify.NoAction), reqID)
		}
		return m.executePasteOperation(&m.processBarModel, dest, locations, true, reqID)
	}
}

// abs returns the absolute value of an int, so that a distance can be measured
// without importing math for one call.
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
