package sidebar

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/ui"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

// Render returns the rendered sidebar string
func (s *Model) Render(mainPanelHeight int, sidebarFocussed bool, currentFilePanelLocation string) string {
	if common.Config.SidebarWidth == 0 {
		return ""
	}
	slog.Debug("Rendering sidebar.", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "dirs count", len(s.directories),
		"sidebar focused", sidebarFocussed)

	r := ui.SidebarRenderer(mainPanelHeight+2, common.Config.SidebarWidth+2, sidebarFocussed)

	r.AddLines(common.SideBarSuperfileTitle, "")

	if s.searchBar.Focused() || s.searchBar.Value() != "" || sidebarFocussed {
		r.AddLines(s.searchBar.View())
	}

	if s.NoActualDir() {
		r.AddLines(common.SideBarNoneText)
	} else {
		s.directoriesRender(mainPanelHeight, currentFilePanelLocation, sidebarFocussed, r)
	}
	return r.Render()
}

func (s *Model) directoriesRender(mainPanelHeight int, curFilePanelFileLocation string, sideBarFocussed bool, r *rendering.Renderer) {
	// Cursor should always point to a valid directory at this point
	if s.isCursorInvalid() {
		slog.Error("Unexpected situation in sideBar Model. "+
			"Cursor is at invalid position, while there are valid directories", "cursor", s.cursor,
			"directory count", len(s.directories))
	}

	// Todo : This is not true when searchbar is not rendered(totalHeight is 2, not 3),
	// so we end up underutilizing one line for our render. But it wont break anything.
	totalHeight := sideBarInitialHeight
	for i := s.renderIndex; i < len(s.directories); i++ {
		if totalHeight+s.directories[i].RequiredHeight() > mainPanelHeight {
			break
		}

		totalHeight += s.directories[i].RequiredHeight()

		switch s.directories[i] {
		case pinnedDividerDir:
			r.AddLines("", common.SideBarPinnedDivider, "")
		case diskDividerDir:
			r.AddLines("", common.SideBarDisksDivider, "")
		default:
			cursor := " "
			if s.cursor == i && sideBarFocussed && !s.searchBar.Focused() {
				cursor = icon.Cursor
			}
			if s.renaming && s.cursor == i {
				r.AddLines(s.rename.View())
			} else {
				renderStyle := common.SidebarStyle
				if s.directories[i].Location == curFilePanelFileLocation {
					renderStyle = common.SidebarSelectedStyle
				}
				line := common.FilePanelCursorStyle.Render(cursor+" ") + renderStyle.Render(s.directories[i].Name)
				r.AddLineWithCustomTruncate(line, rendering.TailsTruncateRight)
			}
		}
	}
}
