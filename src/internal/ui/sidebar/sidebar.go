package sidebar

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

/* SIDE BAR internal TYPE START*/
// Model for sidebar internal
type Model struct {
	Directories []common.Directory
	RenderIndex int
	Cursor      int
	Rename      textinput.Model
	Renaming    bool
	SearchBar   textinput.Model
}

// True if only dividers are in directories slice,
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *Model) NoActualDir() bool {
	for _, d := range s.Directories {
		if !d.IsDivider() {
			return false
		}
	}
	return true
}

func (s *Model) IsCursorInvalid() bool {
	return s.Cursor < 0 || s.Cursor >= len(s.Directories) || s.Directories[s.Cursor].IsDivider()
}

func (s *Model) ResetCursor() {
	s.Cursor = 0
	// Move to first non Divider dir
	for i, d := range s.Directories {
		if !d.IsDivider() {
			s.Cursor = i
			return
		}
	}
	// If all directories are divider, code will reach here. and s.cursor will stay 0
	// Or s.directories is empty
}

// Return till what indexes we will render, if we start from startIndex
// if returned value is `startIndex - 1`, that means nothing can be rendered
// This could be made constant time by keeping Indexes ot special directories saved,
// but that too much.
func (s *Model) LastRenderedIndex(mainPanelHeight int, startIndex int) int {
	curHeight := common.SideBarInitialHeight
	endIndex := startIndex - 1
	for i := startIndex; i < len(s.Directories); i++ {
		curHeight += s.Directories[i].RequiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		endIndex = i
	}
	return endIndex
}

// Return what will be the startIndex, if we end at endIndex
// if returned value is `endIndex + 1`, that means nothing can be rendered
func (s *Model) FirstRenderedIndex(mainPanelHeight int, endIndex int) int {
	// This should ideally never happen. Maybe we should panic ?
	if endIndex >= len(s.Directories) {
		return endIndex + 1
	}

	curHeight := common.SideBarInitialHeight
	startIndex := endIndex + 1
	for i := endIndex; i >= 0; i-- {
		curHeight += s.Directories[i].RequiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		startIndex = i
	}
	return startIndex
}

func (s *Model) UpdateRenderIndex(mainPanelHeight int) {
	// Case I : New cursor moved above current renderable range
	if s.Cursor < s.RenderIndex {
		// We will start rendering from there
		s.RenderIndex = s.Cursor
		return
	}

	curEndIndex := s.LastRenderedIndex(mainPanelHeight, s.RenderIndex)

	// Case II : new cursor also comes in range of rendered directories
	// Taking this case later avoid extra lastRenderedIndex() call
	if s.RenderIndex <= s.Cursor && s.Cursor <= curEndIndex {
		// no need to update s.renderIndex
		return
	}

	// Case III : New cursor is too below
	if curEndIndex < s.Cursor {
		s.RenderIndex = s.FirstRenderedIndex(mainPanelHeight, s.Cursor)
		return
	}

	// Code should never reach here
	slog.Error("Unexpected situation in updateRenderIndex", "cursor", s.Cursor,
		"renderIndex", s.RenderIndex, "directory count", len(s.Directories))
}

// ======================================== Sidebar controller ========================================

func (s *Model) ListUp(mainPanelHeight int) {
	slog.Debug("controlListUp called", "cursor", s.Cursor,
		"renderIndex", s.RenderIndex, "directory count", len(s.Directories))
	if s.NoActualDir() {
		return
	}
	if s.Cursor > 0 {
		// Not at the top, can safely decrease
		s.Cursor--
	} else {
		// We are at the top. Move to the bottom
		s.Cursor = len(s.Directories) - 1
	}
	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.UpdateRenderIndex(mainPanelHeight)
	if s.Directories[s.Cursor].IsDivider() {
		// cause another listUp trigger to move up.
		s.ListUp(mainPanelHeight)
	}
}

func (s *Model) ListDown(mainPanelHeight int) {
	slog.Debug("controlListDown called", "cursor", s.Cursor,
		"renderIndex", s.RenderIndex, "directory count", len(s.Directories))
	if s.NoActualDir() {
		return
	}
	if s.Cursor < len(s.Directories)-1 {
		// Not at the bottom, can safely increase
		s.Cursor++
	} else {
		// We are at the bottom. Move to the top
		s.Cursor = 0
	}

	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.UpdateRenderIndex(mainPanelHeight)

	// Move below special divider directories
	if s.Directories[s.Cursor].IsDivider() {
		// cause another listDown trigger to move down.
		s.ListDown(mainPanelHeight)
	}
}

func (s *Model) DirectoriesRender(mainPanelHeight int, curFilePanelFileLocation string, sideBarFocussed bool) string {
	// Cursor should always point to a valid directory at this point
	if s.IsCursorInvalid() {
		slog.Error("Unexpected situation in sideBar Model. "+
			"Cursor is at invalid position, while there are valide directories", "cursor", s.Cursor,
			"directory count", len(s.Directories))
		return ""
	}

	res := ""
	totalHeight := common.SideBarInitialHeight
	for i := s.RenderIndex; i < len(s.Directories); i++ {
		if totalHeight+s.Directories[i].RequiredHeight() > mainPanelHeight {
			break
		}
		res += "\n"

		totalHeight += s.Directories[i].RequiredHeight()

		switch s.Directories[i] {
		case common.PinnedDividerDir:
			res += "\n" + common.SideBarPinnedDivider
		case common.DiskDividerDir:
			res += "\n" + common.SideBarDisksDivider
		default:
			cursor := " "
			if s.Cursor == i && sideBarFocussed && !s.SearchBar.Focused() {
				cursor = icon.Cursor
			}
			if s.Renaming && s.Cursor == i {
				res += s.Rename.View()
			} else {
				renderStyle := common.SidebarStyle
				if s.Directories[i].Location == curFilePanelFileLocation {
					renderStyle = common.SidebarSelectedStyle
				}
				res += common.FilePanelCursorStyle.Render(cursor+" ") +
					renderStyle.Render(common.TruncateText(s.Directories[i].Name, common.Config.SidebarWidth-2, "..."))
			}
		}
	}
	return res
}

func (s *Model) PinnedIndexRange() (int, int) {
	// pinned directories start after well-known directories and the divider
	// Can't use getPinnedDirectories() here, as if we are in search mode, we would be showing
	// and having less directories in sideBar.directories slice

	// Todo : This is inefficient to iterate each time for this.
	// This information can be kept precomputed
	pinnedDividerIdx := -1
	diskDividerIdx := -1
	for i, d := range s.Directories {
		if d == common.PinnedDividerDir {
			pinnedDividerIdx = i
		}
		if d == common.DiskDividerDir {
			diskDividerIdx = i
			break
		}
	}
	return pinnedDividerIdx + 1, diskDividerIdx - 1
}

// Rename file where the cursor is located
func (s *Model) PinnedItemRename() {
	pinnedBegin, pinnedEnd := s.PinnedIndexRange()
	// We have not selected a pinned directory, rename is not allowed
	if s.Cursor < pinnedBegin || s.Cursor > pinnedEnd {
		return
	}

	nameLen := len(s.Directories[s.Cursor].Name)
	cursorPos := nameLen

	s.Renaming = true
	s.Rename = common.GeneratePinnedRenameTextInput(cursorPos, s.Directories[s.Cursor].Name)
}

// Cancel rename pinned directory
func (s *Model) CancelSidebarRename() {
	s.Rename.Blur()
	s.Renaming = false
}

// Confirm rename pinned directory
func (s *Model) ConfirmSidebarRename() {
	itemLocation := s.Directories[s.Cursor].Location
	newItemName := s.Rename.Value()
	// This is needed to update the current pinned directory data loaded into memory
	s.Directories[s.Cursor].Name = newItemName

	// recover the state of rename
	s.CancelSidebarRename()

	pinnedDirs := common.GetPinnedDirectories()
	// Call getPinnedDirectories, instead of using what is stored in sidebar.directories
	// sidebar.directories could have less directories in case a search filter is used
	for i := range pinnedDirs {
		// Considering the situation when many
		if pinnedDirs[i].Location == itemLocation {
			pinnedDirs[i].Name = newItemName
		}
	}

	jsonData, err := json.Marshal(pinnedDirs)
	if err != nil {
		slog.Error("Error marshaling pinned directories data", "error", err)
	}

	err = os.WriteFile(variable.PinnedFile, jsonData, 0644)
	if err != nil {
		slog.Error("Error updating pinned directories data", "error", err)
	}
}

// UpdateState handles the sidebar's state updates
func (s *Model) UpdateState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if s.Renaming {
		s.Rename, cmd = s.Rename.Update(msg)
	} else if s.SearchBar.Focused() {
		s.SearchBar, cmd = s.SearchBar.Update(msg)
	}

	if s.Cursor < 0 {
		s.Cursor = 0
	}
	return cmd
}

// HandleSearchBarKey handles key events for the sidebar search bar
func (s *Model) HandleSearchBarKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		s.SearchBar.Blur()
		s.SearchBar.SetValue("")
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		s.SearchBar.Blur()
		s.ResetCursor()
	}
}
