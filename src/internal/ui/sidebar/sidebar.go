package sidebar

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/term/ansi"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

/* SIDE BAR internal TYPE START*/
// Model for sidebar internal
type Model struct {
	directories []common.Directory
	renderIndex int
	cursor      int
	rename      textinput.Model
	renaming    bool
	searchBar   textinput.Model
}

// True if only dividers are in directories slice,
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *Model) NoActualDir() bool {
	for _, d := range s.directories {
		if !d.IsDivider() {
			return false
		}
	}
	return true
}

func (s *Model) IsCursorInvalid() bool {
	return s.cursor < 0 || s.cursor >= len(s.directories) || s.directories[s.cursor].IsDivider()
}

func (s *Model) ResetCursor() {
	s.cursor = 0
	// Move to first non Divider dir
	for i, d := range s.directories {
		if !d.IsDivider() {
			s.cursor = i
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
	for i := startIndex; i < len(s.directories); i++ {
		curHeight += s.directories[i].RequiredHeight()
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
	if endIndex >= len(s.directories) {
		return endIndex + 1
	}

	curHeight := common.SideBarInitialHeight
	startIndex := endIndex + 1
	for i := endIndex; i >= 0; i-- {
		curHeight += s.directories[i].RequiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		startIndex = i
	}
	return startIndex
}

func (s *Model) UpdateRenderIndex(mainPanelHeight int) {
	// Case I : New cursor moved above current renderable range
	if s.cursor < s.renderIndex {
		// We will start rendering from there
		s.renderIndex = s.cursor
		return
	}

	curEndIndex := s.LastRenderedIndex(mainPanelHeight, s.renderIndex)

	// Case II : new cursor also comes in range of rendered directories
	// Taking this case later avoid extra lastRenderedIndex() call
	if s.renderIndex <= s.cursor && s.cursor <= curEndIndex {
		// no need to update s.renderIndex
		return
	}

	// Case III : New cursor is too below
	if curEndIndex < s.cursor {
		s.renderIndex = s.FirstRenderedIndex(mainPanelHeight, s.cursor)
		return
	}

	// Code should never reach here
	slog.Error("Unexpected situation in updateRenderIndex", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
}

// ======================================== Sidebar controller ========================================

func (s *Model) ListUp(mainPanelHeight int) {
	slog.Debug("controlListUp called", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
	if s.NoActualDir() {
		return
	}
	if s.cursor > 0 {
		// Not at the top, can safely decrease
		s.cursor--
	} else {
		// We are at the top. Move to the bottom
		s.cursor = len(s.directories) - 1
	}
	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.UpdateRenderIndex(mainPanelHeight)
	if s.directories[s.cursor].IsDivider() {
		// cause another listUp trigger to move up.
		s.ListUp(mainPanelHeight)
	}
}

func (s *Model) ListDown(mainPanelHeight int) {
	slog.Debug("controlListDown called", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
	if s.NoActualDir() {
		return
	}
	if s.cursor < len(s.directories)-1 {
		// Not at the bottom, can safely increase
		s.cursor++
	} else {
		// We are at the bottom. Move to the top
		s.cursor = 0
	}

	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.UpdateRenderIndex(mainPanelHeight)

	// Move below special divider directories
	if s.directories[s.cursor].IsDivider() {
		// cause another listDown trigger to move down.
		s.ListDown(mainPanelHeight)
	}
}

func (s *Model) DirectoriesRender(mainPanelHeight int, curFilePanelFileLocation string, sideBarFocussed bool) string {
	// Cursor should always point to a valid directory at this point
	if s.IsCursorInvalid() {
		slog.Error("Unexpected situation in sideBar Model. "+
			"Cursor is at invalid position, while there are valide directories", "cursor", s.cursor,
			"directory count", len(s.directories))
		return ""
	}

	res := ""
	totalHeight := common.SideBarInitialHeight
	for i := s.renderIndex; i < len(s.directories); i++ {
		if totalHeight+s.directories[i].RequiredHeight() > mainPanelHeight {
			break
		}
		res += "\n"

		totalHeight += s.directories[i].RequiredHeight()

		switch s.directories[i] {
		case common.PinnedDividerDir:
			res += "\n" + common.SideBarPinnedDivider
		case common.DiskDividerDir:
			res += "\n" + common.SideBarDisksDivider
		default:
			cursor := " "
			if s.cursor == i && sideBarFocussed && !s.searchBar.Focused() {
				cursor = icon.Cursor
			}
			if s.renaming && s.cursor == i {
				res += s.rename.View()
			} else {
				renderStyle := common.SidebarStyle
				if s.directories[i].Location == curFilePanelFileLocation {
					renderStyle = common.SidebarSelectedStyle
				}
				res += common.FilePanelCursorStyle.Render(cursor+" ") +
					renderStyle.Render(common.TruncateText(s.directories[i].Name, common.Config.SidebarWidth-2, "..."))
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
	for i, d := range s.directories {
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
	if s.cursor < pinnedBegin || s.cursor > pinnedEnd {
		return
	}

	nameLen := len(s.directories[s.cursor].Name)
	cursorPos := nameLen

	s.renaming = true
	s.rename = common.GeneratePinnedRenameTextInput(cursorPos, s.directories[s.cursor].Name)
}

// Cancel rename pinned directory
func (s *Model) CancelSidebarRename() {
	s.rename.Blur()
	s.renaming = false
}

// Confirm rename pinned directory
func (s *Model) ConfirmSidebarRename() {
	itemLocation := s.directories[s.cursor].Location
	newItemName := s.rename.Value()
	// This is needed to update the current pinned directory data loaded into memory
	s.directories[s.cursor].Name = newItemName

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
	if s.renaming {
		s.rename, cmd = s.rename.Update(msg)
	} else if s.searchBar.Focused() {
		s.searchBar, cmd = s.searchBar.Update(msg)
	}

	if s.cursor < 0 {
		s.cursor = 0
	}
	return cmd
}

// HandleSearchBarKey handles key events for the sidebar search bar
func (s *Model) HandleSearchBarKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		s.SearchBarBlur()
		s.searchBar.SetValue("")
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		s.SearchBarBlur()
		s.ResetCursor()
	}
}

// SearchBarFocused returns whether the search bar is focused
func (s *Model) SearchBarFocused() bool {
	return s.searchBar.Focused()
}

// SearchBarBlur removes focus from the search bar
func (s *Model) SearchBarBlur() {
	s.searchBar.Blur()
}

// SearchBarFocus sets focus on the search bar
func (s *Model) SearchBarFocus() {
	s.searchBar.Focus()
}

// UpdateDirectories updates the directories list based on search value
// This is a bit inefficient, as we already had the directories when we
// initialized the sidebar. We call the directory fetching logic many times
// which is a disk heavy operation.
func (s *Model) UpdateDirectories() {
	if s.searchBar.Value() != "" {
		s.directories = common.GetFilteredDirectories(s.searchBar.Value())
	} else {
		s.directories = common.GetDirectories()
	}
	// This is needed, as due to filtering, the cursor might be invalid
	if s.IsCursorInvalid() {
		s.ResetCursor()
	}
}

// Render returns the rendered sidebar string
func (s *Model) Render(mainPanelHeight int, isSidebarFocussed bool, currentFilePanelLocation string) string {
	if common.Config.SidebarWidth == 0 {
		return ""
	}
	slog.Debug("Rendering sidebar.", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "dirs count", len(s.directories),
		"sidebar focused", isSidebarFocussed)

	content := common.SideBarSuperfileTitle + "\n"

	if s.searchBar.Focused() || s.searchBar.Value() != "" || isSidebarFocussed {
		s.searchBar.Placeholder = "(" + common.Hotkeys.SearchBar[0] + ")" + " Search"
		content += "\n" + ansi.Truncate(s.searchBar.View(), common.Config.SidebarWidth-2, "...")
	}

	if s.NoActualDir() {
		content += "\n" + common.SideBarNoneText
	} else {
		content += s.DirectoriesRender(mainPanelHeight, currentFilePanelLocation, isSidebarFocussed)
	}
	return common.SideBarBorderStyle(mainPanelHeight, isSidebarFocussed).Render(content)
}

// GetCurrentDirectoryLocation returns the location of the currently selected directory
func (s *Model) GetCurrentDirectoryLocation() string {
	if s.IsCursorInvalid() || s.NoActualDir() {
		return ""
	}
	return s.directories[s.cursor].Location
}

// New creates a new sidebar model with the given parameters
func New(directories []common.Directory, searchBar textinput.Model) Model {
	return Model{
		renderIndex: 0,
		directories: directories,
		searchBar:   searchBar,
	}
}

// IsRenaming returns whether the sidebar is currently in renaming mode
func (s *Model) IsRenaming() bool {
	return s.renaming
}
