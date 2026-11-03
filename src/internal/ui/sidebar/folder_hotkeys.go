package sidebar

import (
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// FolderShortcutLocation uses the same filtered pin order as the sidebar,
// independent of the scroll position. Hidden sidebars still support navigation.
func (s *Model) FolderShortcutLocation(shortcut common.FolderShortcut) string {
	if shortcut.Pin == 0 {
		for _, dir := range getWellKnownDirectories() {
			if dir.Name == shortcut.Name {
				return dir.Location
			}
		}
		return ""
	}

	dirs := s.directories
	if s.Disabled() || !slices.Contains(s.sections, utils.SidebarSectionPinned) {
		if s.pinnedMgr == nil {
			return ""
		}
		dirs = getFilteredDirectories(s.searchBar.Value(), s.pinnedMgr, []string{utils.SidebarSectionPinned})
	}
	pin := 0
	for _, dir := range dirs {
		if dir.Section == utils.SidebarSectionPinned {
			pin++
			if pin == shortcut.Pin {
				return dir.Location
			}
		}
	}
	return ""
}

// folderHints maps rows to shortcuts by section and identity, not location:
// the same path can occur as several standard folders and as a pin.
func (s *Model) folderHints() map[int]string {
	hints := make(map[int]string)
	shortcuts := common.Hotkeys.FolderShortcuts()
	pin := 0
	for i, dir := range s.directories {
		if dir.Section == utils.SidebarSectionPinned {
			pin++
		}
		for _, shortcut := range shortcuts {
			isPin := dir.Section == utils.SidebarSectionPinned && shortcut.Pin == pin
			isHome := dir.Section == utils.SidebarSectionHome && shortcut.Pin == 0 && shortcut.Name == dir.Name
			if isPin || isHome {
				hints[i] = shortcut.AltHint()
				if isHome {
					hints[i] = strings.ToUpper(hints[i])
				}
				break
			}
		}
	}
	return hints
}

func appendFolderHint(line, hint string, width int) string {
	if hint == "" {
		return line
	}
	badge := folderHintBadge(hint)
	badgeWidth := ansi.StringWidth(badge)
	if badgeWidth+1 >= width {
		return line
	}
	labelWidth := width - badgeWidth - 1
	line = ansi.Truncate(line, labelWidth, "…")
	padding := strings.Repeat(" ", width-ansi.StringWidth(line)-badgeWidth)
	return line + common.SidebarStyle.Render(padding) + badge
}

// Powerline semicircles keep rounded badges on one terminal row. Use an ASCII
// keycap when Nerd Fonts are disabled.
func folderHintBadge(hint string) string {
	fill := lipgloss.Color(common.Theme.Hint)
	labelStyle := common.SidebarStyle.Background(fill).
		Foreground(lipgloss.Color(common.Theme.SidebarBG)).Bold(true)
	if !common.Config.Nerdfont {
		return labelStyle.Render("[" + hint + "]")
	}
	capStyle := common.SidebarStyle.Foreground(fill)
	return capStyle.Render("\ue0b6") + labelStyle.Render(hint) + capStyle.Render("\ue0b4")
}
