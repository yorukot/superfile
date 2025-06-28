package common

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/config/icon"
)

const WheelRunTime = 5
const DefaultCommandTimeout = 5000 * time.Millisecond
const DateModifiedOption = "Date Modified"

var (
	MinimumHeight = 24
	MinimumWidth  = 60

	// TODO : These are model object properties, not global properties
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

	FilePreviewNoContentText           string
	FilePreviewNoFileInfoText          string
	FilePreviewUnsupportedFormatText   string
	FilePreviewDirectoryUnreadableText string
	FilePreviewEmptyText               string
	FilePreviewError                   string

	LipglossError string
)

var (
	UnsupportedPreviewFormats = []string{".pdf", ".torrent"}
)

// No dependencies
func LoadInitialPrerenderedVariables() {
	LipglossError = lipgloss.NewStyle().Foreground(lipgloss.Color("#F93939")).Render("Error") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFEE")).Render(" ┃ ")
}

// Dependecies - TODO We should programmatically guarantee these dependencies. And log error
// if its not satisfied.
// LoadThemeConfig() in style.go should be finished
// loadConfigFile() in config_types.go should be finished
// InitIcon() in config package in function.go should be finished
func LoadPrerenderedVariables() {
	SideBarSuperfileTitle = SidebarTitleStyle.Render(" " + icon.SuperfileIcon + icon.Space + "superfile")

	SideBarPinnedDivider = SidebarTitleStyle.Render(icon.Pinned+icon.Space+"Pinned") + SidebarDividerStyle.Render(" ───────────")

	SideBarDisksDivider = SidebarTitleStyle.Render(icon.Disk+icon.Space+"Disks") + SidebarDividerStyle.Render(" ────────────")

	SideBarNoneText = SidebarStyle.Render(" " + icon.Error + icon.Space + "None")

	ProcessBarNoneText = icon.Error + icon.Space + "No processes running"

	FilePanelTopDirectoryIcon = FilePanelTopDirectoryIconStyle.Render(" " + icon.Directory + icon.Space)
	FilePanelNoneText = FilePanelStyle.Render(" " + icon.Error + icon.Space + "No such file or directory")

	FilePreviewNoContentText = "\n--- " + icon.Error + icon.Space + "No content to preview" + icon.Space + "---"
	FilePreviewNoFileInfoText = "\n--- " + icon.Error + icon.Space + "Could not get file info" + icon.Space + "---"
	FilePreviewUnsupportedFormatText = "\n--- " + icon.Error + icon.Space + "Unsupported formats" + icon.Space + "---"
	FilePreviewDirectoryUnreadableText = "\n--- " + icon.Error + icon.Space + "Cannot read directory" + icon.Space + "---"
	FilePreviewError = "\n--- " + icon.Error + icon.Space + "Error" + icon.Space + "---"
	FilePreviewEmptyText = "\n--- Empty ---"
}
