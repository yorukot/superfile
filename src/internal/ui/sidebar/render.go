package sidebar

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/ui"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

// Render returns the rendered sidebar string.
func (s *Model) Render(sidebarFocused bool, currentFilePanelLocation string) string {
	if s.Disabled() {
		return ""
	}

	r := ui.SidebarRenderer(s.height, s.width, sidebarFocused)

	r.AddLines(common.SideBarSuperfileTitle, "")

	if s.searchBarRendered(sidebarFocused) {
		r.AddLines(s.searchBar.View())
	}

	if s.NoActualDir() {
		r.AddLines(common.SideBarNoneText)
	} else {
		s.directoriesRender(currentFilePanelLocation, sidebarFocused, r)
	}
	return r.Render()
}

// directoriesRender handles the iterative rendering of directories within the sidebar model.
func (s *Model) directoriesRender(curFilePanelFileLocation string,
	sideBarFocused bool, r *rendering.Renderer) {
	// Cursor should always point to a valid directory at this point
	if s.isCursorInvalid() {
		slog.Error("Unexpected situation in sideBar Model. "+
			"Cursor is at invalid position, while there are valid directories", "cursor", s.cursor,
			"directory count", len(s.directories))
	}

	// TODO : This is not true when searchbar is not rendered(totalHeight is 2, not 3),
	// so we end up underutilizing one line for our render. But it wont break anything.
	for _, row := range s.layoutDirectoryRows(sideBarFocused) {
		s.renderDirectory(row.index, curFilePanelFileLocation, sideBarFocused, r)
	}
}

// renderDirectory renders the single entry identified by index, choosing the
// right form for a section divider or an ordinary directory.
func (s *Model) renderDirectory(index int, curFilePanelFileLocation string,
	sideBarFocused bool, r *rendering.Renderer) {
	switch s.directories[index] {
	case homeDividerDir:
		r.AddLines("", common.SideBarHomeDivider, "")
	case pinnedDividerDir:
		r.AddLines("", common.SideBarPinnedDivider, "")
	case diskDividerDir:
		r.AddLines("", common.SideBarDisksDivider, "")
	default:
		s.renderDirectoryEntry(index, curFilePanelFileLocation, sideBarFocused, r)
	}
}

// renderDirectoryEntry renders one selectable directory row.
func (s *Model) renderDirectoryEntry(index int, curFilePanelFileLocation string,
	sideBarFocused bool, r *rendering.Renderer) {
	if s.renaming && s.cursor == index {
		r.AddLines(s.rename.View())
		return
	}

	cursor := " "
	if s.cursor == index && sideBarFocused && !s.searchBar.Focused() {
		cursor = icon.Cursor
	}

	renderStyle := common.SidebarStyle
	if s.directories[index].Location == curFilePanelFileLocation {
		renderStyle = common.SidebarSelectedStyle
	}
	line := common.FilePanelCursorStyle.Render(cursor+" ") +
		renderStyle.Render(s.directories[index].Icon+" ") +
		renderStyle.Render(s.directories[index].Name)
	r.AddLineWithCustomTruncate(line, rendering.TailsTruncateRight)
}
