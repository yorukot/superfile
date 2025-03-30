package common

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/superfile/src/config/icon"
)

const WheelRunTime = 5

var (
	MinimumHeight = 24
	MinimumWidth  = 60

	// Todo : These are model object properties, not global properties
	// We are modifying them in the code many time. They need to be part of model struct.
	MinFooterHeight = 6
	ModalWidth      = 60
	ModalHeight     = 7
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
	sideBarSuperfileTitle = SidebarTitleStyle.Render("    " + icon.SuperfileIcon + " superfile")
	sideBarSuperfileTitle = ansi.Truncate(sideBarSuperfileTitle, Config.SidebarWidth, "")

	sideBarPinnedDivider = SidebarTitleStyle.Render("󰐃 Pinned") + SidebarDividerStyle.Render(" ───────────") + "\n"
	sideBarPinnedDivider = ansi.Truncate(sideBarPinnedDivider, Config.SidebarWidth, "")

	sideBarDisksDivider = SidebarTitleStyle.Render("󱇰 Disks") + SidebarDividerStyle.Render(" ────────────") + "\n"
	sideBarDisksDivider = ansi.Truncate(sideBarDisksDivider, Config.SidebarWidth, "")

	sideBarNoneText = SidebarStyle.Render(" " + icon.Error + " None")
}
