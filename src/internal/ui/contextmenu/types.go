// Package contextmenu implements the menu shown on right click.
//
// The menu is an overlay rather than a centred modal: it is anchored to the
// pointer and moves with it, so it has to know where it was placed in order to
// both draw itself and publish hit-test regions that agree with that drawing.
package contextmenu

// Action identifies what a menu entry does.
//
// The menu deliberately does not know how to perform any of them. It reports
// which action was chosen and lets the caller run the matching handler, which
// keeps this package free of any dependency on file operations.
type Action int

const (
	// ActionNone is a placeholder for an entry that does nothing.
	ActionNone Action = iota
	// ActionOpen enters a directory or opens a file.
	ActionOpen
	// ActionCopy copies the target to the clipboard.
	ActionCopy
	// ActionCut copies the target to the clipboard and marks it for moving.
	ActionCut
	// ActionPaste writes the clipboard into the target directory.
	ActionPaste
	// ActionRename starts renaming the target.
	ActionRename
	// ActionDelete moves the target to the trash.
	ActionDelete
	// ActionCopyPath copies the path of the target.
	ActionCopyPath
	// ActionNewFile starts creating a file in the target directory.
	ActionNewFile
	// ActionNewFolder starts creating a directory in the target directory.
	ActionNewFolder
)

// Item is a single row of the menu.
type Item struct {
	// Label is the text shown for the entry.
	Label string
	// Shortcut is the hotkey that performs the same action, shown right aligned
	// so the menu doubles as a reminder. It may be empty.
	Shortcut string
	// Action is what the caller should run when this entry is chosen.
	Action Action
}

// Model is the state of the context menu.
//
// The zero value is a closed menu, which is what the application starts with.
type Model struct {
	// open reports whether the menu is currently shown. Entries can only be
	// chosen while it is.
	open bool
	// originX and originY are the terminal coordinates of the menu's top left
	// corner after clamping it to fit on screen.
	originX int
	originY int
	// cursor is the index of the highlighted entry.
	cursor int
	// items are the entries currently offered, in display order.
	items []Item
	// width is the total width of the rendered menu, borders included.
	width int
	// height is the total height of the rendered menu, borders included.
	height int
}
