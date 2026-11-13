package internal

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// doubleClickWindow is the maximum gap between two clicks on the same widget for
// them to count as a double click.
const doubleClickWindow = 400 * time.Millisecond

// Actions reported by the wheel. The wheel is dispatched on the button rather
// than on the message type, because terminals keep sending motion events with
// the wheel button held while the user keeps scrolling.
const (
	wheelUpAction   = "wheelup"
	wheelDownAction = "wheeldown"
)

// leftClick records the most recent left click so that a quick second click on
// the same widget can be recognised as a double click.
//
// This lives on the model rather than in a package level variable so that it
// cannot leak between two models or survive into an unrelated test.
type leftClick struct {
	// target is the widget the click landed on.
	target mouse.Target
	// at is when the click happened.
	at time.Time
}

// modalBlocksMouse reports whether a modal is covering the screen.
//
// The keyboard is already gated on the modal states below, so a modal only
// captures keys today. Without the same gate for the pointer, a right click
// would move the panel cursor, focus, and selection underneath a modal, and a
// click on a context menu entry would dispatch its action over the top of one.
func (m *model) modalBlocksMouse() bool {
	return m.typingModal.open ||
		m.promptModal.IsOpen() ||
		m.notifyModel.IsOpen() ||
		m.zoxideModal.IsOpen() ||
		m.sortModal.IsOpen() ||
		m.helpMenu.IsOpen() ||
		m.spfError.IsOpen() ||
		m.fileModel.Renaming ||
		m.sidebarModel.IsRenaming()
}

// handleMouseMsg resolves a mouse event against the hit-test regions published
// by the last render.
func (m *model) handleMouseMsg(msg tea.MouseMsg) tea.Cmd {
	event := msg.Mouse()

	// A modal that was open when a context menu opened does not become
	// dismissible by clicking the menu underneath it, so the menu goes first and
	// the click that closed it is consumed rather than falling through to
	// whatever it landed on.
	if m.modalBlocksMouse() {
		return m.closeContextMenu()
	}

	if event.Button == tea.MouseWheelUp {
		wheelMainAction(wheelUpAction, m)
		return nil
	}
	if event.Button == tea.MouseWheelDown {
		wheelMainAction(wheelDownAction, m)
		return nil
	}

	if motion, isMotion := msg.(tea.MouseMotionMsg); isMotion {
		m.handleMouseMotion(motion.Mouse())
		return nil
	}

	click, isClick := msg.(tea.MouseClickMsg)
	if !isClick {
		// Release carries no action of its own yet. Drag and drop reads it.
		return nil
	}
	return m.handleMouseClick(click.Mouse())
}

// handleMouseMotion tracks the pointer while no button is held.
//
// Only the open context menu reacts to this, by highlighting the entry under the
// pointer so that the mouse and the keyboard agree on what is chosen.
func (m *model) handleMouseMotion(event tea.Mouse) {
	// A wheel scroll also arrives as motion while the button is held, and that
	// is not the pointer moving.
	if event.Button != tea.MouseNone || !m.contextMenu.IsOpen() {
		return
	}

	target := m.mouseTargetAt(event.X, event.Y)
	if target.Kind == mouse.TargetContextMenuItem {
		m.contextMenu.HighlightItem(target.ItemIndex)
	}
}

// isSelectionMod reports whether the modifiers turn a click into a selection
// gesture rather than a plain pointer move.
func isSelectionMod(mod tea.KeyMod) bool {
	return mod.Contains(tea.ModCtrl) || mod.Contains(tea.ModShift)
}

// handleMouseClick acts on a completed button press.
func (m *model) handleMouseClick(event tea.Mouse) tea.Cmd {
	target := m.mouseTargetAt(event.X, event.Y)

	// Right click always reopens the menu at the pointer, which is what lets an
	// open menu be moved without being closed first.
	if event.Button == tea.MouseRight {
		return m.openContextMenu(target, event.X, event.Y)
	}

	if event.Button != tea.MouseLeft {
		return nil
	}

	// While the menu is open it takes every click: one on an entry chooses it,
	// one anywhere else just dismisses the menu.
	if m.contextMenu.IsOpen() {
		return m.handleContextMenuClick(target)
	}

	// A click carrying a selection modifier is part of building a selection
	// rather than pointing at something, so it never counts as a double click
	// and it clears any pending pair. Without this, ctrl clicking a row twice to
	// deselect it would read as a double click and open the file instead.
	var isDoubleClick bool
	if isSelectionMod(event.Mod) {
		m.lastLeftClick = leftClick{}
	} else {
		isDoubleClick = m.noteLeftClick(target)
	}

	switch target.Kind {
	case mouse.TargetFilePanelItem:
		return m.handleFilePanelClick(target, event.Mod, isDoubleClick)
	case mouse.TargetSidebarDirectory:
		m.handleSidebarDirectoryClick(target.ItemIndex)
	case mouse.TargetUnknown:
		m.handleBackgroundClick(event.X, event.Y)
	case mouse.TargetProcessBarItem, mouse.TargetMetadataItem:
		// These have no published regions yet, so a click landing on one of them
		// resolves as the background instead.
	case mouse.TargetContextMenuItem:
		// The context menu is open here, so this click was already consumed by
		// handleContextMenuClick above.
	}
	return nil
}

// noteLeftClick records a left click and reports whether it completes a double
// click on the same widget.
//
// A recognised pair is consumed so that a third click starts a fresh one rather
// than immediately entering whatever sits under the pointer.
func (m *model) noteLeftClick(target mouse.Target) bool {
	now := time.Now()
	if m.lastLeftClick.at.IsZero() {
		m.lastLeftClick = leftClick{target: target, at: now}
		return false
	}

	if target != m.lastLeftClick.target || now.Sub(m.lastLeftClick.at) >= doubleClickWindow {
		m.lastLeftClick = leftClick{target: target, at: now}
		return false
	}

	m.lastLeftClick = leftClick{}
	return true
}

// handleContextMenuClick acts on a click while the context menu is open.
func (m *model) handleContextMenuClick(target mouse.Target) tea.Cmd {
	if target.Kind != mouse.TargetContextMenuItem {
		// A click away from the menu dismisses it without acting on whatever is
		// underneath, so the click that closes a menu never has a side effect.
		m.contextMenu.Close()
		return nil
	}

	// The entry is highlighted first so that the chosen action is the one that
	// was clicked, not the one the keyboard had left highlighted.
	m.contextMenu.HighlightItem(target.ItemIndex)
	action, ok := m.contextMenu.SelectedAction()
	m.contextMenu.Close()
	if !ok {
		return nil
	}
	return m.runContextMenuAction(action)
}

// handleFilePanelClick acts on a click that landed on a file entry.
func (m *model) handleFilePanelClick(target mouse.Target, mod tea.KeyMod, isDoubleClick bool) tea.Cmd {
	m.focusFilePanelOnMouse(target.PanelIndex)
	panel := m.getFocusedFilePanel()

	if isDoubleClick {
		// A double click is what enters a directory or opens a file, matching
		// the keyboard where the cursor is moved first and entered after.
		m.enterPanel()
		return nil
	}

	if !panel.SetCursorToIndex(target.ItemIndex) {
		return nil
	}

	if !m.applyClickSelection(panel, target.ItemIndex, mod) {
		return nil
	}

	// Copy, paste and delete only read the selection while the panel is in
	// SelectMode, so a selection built with the pointer has to put the panel in
	// that mode or it would be ignored by every one of those actions.
	if panel.SelectedCount() > 0 {
		panel.SetPanelMode(filepanel.SelectMode)
	}
	return nil
}

// applyClickSelection acts on the modifiers of a click that landed on a file
// entry. It reports whether the click changed the selection.
func (m *model) applyClickSelection(panel *filepanel.Model, itemIndex int, mod tea.KeyMod) bool {
	switch {
	case mod.Contains(tea.ModCtrl):
		// Ctrl click adds the entry to the selection, or removes it when it is
		// already selected, and it becomes the anchor for a later range.
		panel.SetSelectionAnchor(itemIndex)
		panel.ToggleSelected(panel.GetElementAtIdx(itemIndex).Location)
		return true
	case mod.Contains(tea.ModShift):
		// Shift click selects everything between the anchor and this entry. The
		// anchor stays put so that a second shift click measures from the same
		// point again. The range only grows; a ctrl click is how it shrinks.
		if panel.SelectRangeToIndex(itemIndex) {
			return true
		}
		// Without a usable anchor there is no range to extend, so the click
		// points at the entry and becomes the anchor for the next one.
		panel.SetSelectionAnchor(itemIndex)
		return false
	default:
		// An unmodified click only points at the entry, but it still sets the
		// anchor so that a following shift click has a starting point.
		panel.SetSelectionAnchor(itemIndex)
		return false
	}
}

// handleSidebarDirectoryClick moves the sidebar cursor and navigates to the
// chosen directory, which is what the keyboard does when a sidebar entry is
// used.
func (m *model) handleSidebarDirectoryClick(itemIndex int) {
	m.focusSidebarOnMouse()
	if !m.sidebarModel.SetCursor(itemIndex) {
		return
	}
	m.sidebarSelectDirectory()
}

// handleBackgroundClick focuses whichever non item region was clicked, so that
// the sidebar and the footer panels can be reached with the pointer even though
// they publish no item regions yet.
func (m *model) handleBackgroundClick(x, y int) {
	// The footer occupies an exact band of rows, and the row below it is the
	// bottom border, so the band is matched rather than treated as "anything
	// past the main panel".
	footerEnd := m.mainPanelHeight + utils.FullFooterHeight(m.footerHeight, m.toggleFooter)
	if m.toggleFooter && y >= m.mainPanelHeight && y < footerEnd {
		m.focusFooterOnMouse(x)
		return
	}
	if common.Config.SidebarWidth != 0 && x < common.Config.SidebarWidth+common.BorderPadding {
		m.focusSidebarOnMouse()
	}
}

// focusFilePanelOnMouse moves focus to a file panel.
//
// This does not toggle the way the keyboard focus helpers do: clicking a panel
// that already has focus must leave it focused rather than handing focus back.
func (m *model) focusFilePanelOnMouse(panelIndex int) {
	m.focusPanel = nonePanelFocus
	m.fileModel.SetFocusedPanelIndex(panelIndex)
	// SetFocusedPanelIndex leaves an already focused panel untouched, so make
	// sure the panel really carries focus when focus is coming back from the
	// sidebar or a footer panel.
	m.getFocusedFilePanel().IsFocused = true
}

// focusSidebarOnMouse focuses the sidebar without toggling.
func (m *model) focusSidebarOnMouse() {
	if common.Config.SidebarWidth == 0 {
		return
	}
	m.focusPanel = sidebarFocus
	m.getFocusedFilePanel().IsFocused = false
}

// focusFooterOnMouse focuses the footer panel under the given column. The
// process bar sits leftmost in the footer, and the metadata panel follows it.
func (m *model) focusFooterOnMouse(x int) {
	if !m.toggleFooter {
		return
	}
	if x < m.processBarModel.GetWidth() {
		m.focusPanel = processBarFocus
	} else {
		m.focusPanel = metadataFocus
	}
	m.getFocusedFilePanel().IsFocused = false
}
