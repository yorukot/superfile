package internal

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/superfile/src/config/icon"
)

var (
	sideBarSuperfileTitle string
	sideBarPinnedDivider  string
	sideBarDisksDivider   string
	sideBarNoneText       string
	lipglossError         string
)

// No dependencies
func LoadInitial_PrerenderedVariables() {
	lipglossError = lipgloss.NewStyle().Foreground(lipgloss.Color("#F93939")).Render("Error") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFEE")).Render(" ┃ ")
}

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
