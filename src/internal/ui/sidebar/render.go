package sidebar

import (
	"log/slog"

	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

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
		content += s.directoriesRender(mainPanelHeight, currentFilePanelLocation, isSidebarFocussed)
	}
	return common.SideBarBorderStyle(mainPanelHeight, isSidebarFocussed).Render(content)
}

func (s *Model) directoriesRender(mainPanelHeight int, curFilePanelFileLocation string, sideBarFocussed bool) string {
	// Cursor should always point to a valid directory at this point
	if s.isCursorInvalid() {
		slog.Error("Unexpected situation in sideBar Model. "+
			"Cursor is at invalid position, while there are valid directories", "cursor", s.cursor,
			"directory count", len(s.directories))
		return ""
	}

	res := ""
	totalHeight := sideBarInitialHeight
	for i := s.renderIndex; i < len(s.directories); i++ {
		if totalHeight+s.directories[i].RequiredHeight() > mainPanelHeight {
			break
		}
		res += "\n"

		totalHeight += s.directories[i].RequiredHeight()

		switch s.directories[i] {
		case pinnedDividerDir:
			res += "\n" + common.SideBarPinnedDivider
		case diskDividerDir:
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
