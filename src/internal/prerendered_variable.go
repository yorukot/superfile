package internal

import (
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/superfile/src/config/icon"
)

var (
	sideBarSuperfileTitle string
	sideBarPinnedDivider  string
	sideBarDisksDivider   string
	sideBarNoneText       string
)

// Dependecies
// LoadThemeConfig() in style.go should be finished
// loadConfigFile() in config_types.go should be finished
// InitIcon() in config package in function.go should be finished
func LoadPrerenderedVariables() {
	sideBarSuperfileTitle = sidebarTitleStyle.Render("    " + icon.SuperfileIcon + " superfile")
	sideBarSuperfileTitle = ansi.Truncate(sideBarSuperfileTitle, Config.SidebarWidth, "")

	sideBarPinnedDivider = sidebarTitleStyle.Render("󰐃 Pinned") + sidebarDividerStyle.Render(" ───────────") + "\n"
	sideBarPinnedDivider = ansi.Truncate(sideBarPinnedDivider, Config.SidebarWidth, "")

	sideBarDisksDivider = sidebarTitleStyle.Render("󱇰 Disks") + sidebarDividerStyle.Render(" ────────────") + "\n"
	sideBarDisksDivider = ansi.Truncate(sideBarDisksDivider, Config.SidebarWidth, "")

	sideBarNoneText = sidebarStyle.Render(" " + icon.Error + " None")
}
