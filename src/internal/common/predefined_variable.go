package common

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/config/icon"
)

const WheelRunTime = 5
const DefaultCommandTimeout = 5000 * time.Millisecond

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
	SideBarSuperfileTitle string
	SideBarPinnedDivider  string
	SideBarDisksDivider   string
	SideBarNoneText       string

	ProcessBarNoneText string

	FilePanelTopDirectoryIcon string
	FilePanelNoneText         string

	LipglossError string
)

// No dependencies
func LoadInitialPrerenderedVariables() {
	LipglossError = lipgloss.NewStyle().Foreground(lipgloss.Color("#F93939")).Render("Error") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFEE")).Render(" ┃ ")
}

// Dependecies - Todo We should programmatically guarantee these dependencies. And log error
// if its not satisfied.
// LoadThemeConfig() in style.go should be finished
// loadConfigFile() in config_types.go should be finished
// InitIcon() in config package in function.go should be finished
func LoadPrerenderedVariables() {
	SideBarSuperfileTitle = SidebarTitleStyle.Render(" " + icon.SuperfileIcon + " superfile")

	SideBarPinnedDivider = SidebarTitleStyle.Render("󰐃 Pinned") + SidebarDividerStyle.Render(" ───────────")

	SideBarDisksDivider = SidebarTitleStyle.Render("󱇰 Disks") + SidebarDividerStyle.Render(" ────────────")

	SideBarNoneText = SidebarStyle.Render(" " + icon.Error + " None")

	ProcessBarNoneText = icon.Error + "  No processes running"

	FilePanelTopDirectoryIcon = FilePanelTopDirectoryIconStyle.Render(" " + icon.Directory + icon.Space)
	FilePanelNoneText = FilePanelStyle.Render(" " + icon.Error + "  No such file or directory")
}
