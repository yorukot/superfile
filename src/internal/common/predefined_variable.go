package common

import (
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
)

const WheelRunTime = 5
const DefaultCommandTimeout = 5000 * time.Millisecond
const DateModifiedOption = "Date Modified"

const SameRenameWarnTitle = "There is already a file or directory with that name"
const SameRenameWarnContent = "This operation will override the existing file"

const TrashWarnTitle = "Are you sure you want to move this to trash can"
const TrashWarnContent = "This operation will move file or directory to trash can."
const PermanentDeleteWarnTitle = "Are you sure you want to completely delete"
const PermanentDeleteWarnContent = "This operation cannot be undone and your data will be completely lost."

const (
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
	FilePreviewUnsupportedFileMode     string
	FilePreviewDirectoryUnreadableText string
	FilePreviewEmptyText               string
	FilePreviewError                   string

	CheckboxChecked        string
	CheckboxCheckedFocused string
	CheckboxEmpty          string
	CheckboxEmptyFocused   string

	ModalConfirmInputText string
	ModalCancelInputText  string
	ModalOkayInputText    string
	ModalInputSpacingText string
	LipglossError         string
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

	SideBarPinnedDivider = SidebarTitleStyle.Render(icon.Pinned+icon.Space+"Pinned") +
		SidebarDividerStyle.Render(" ───────────")

	SideBarDisksDivider = SidebarTitleStyle.Render(icon.Disk+icon.Space+"Disks") +
		SidebarDividerStyle.Render(" ────────────")

	SideBarNoneText = SidebarStyle.Render(" " + icon.Error + icon.Space + "None")

	ProcessBarNoneText = icon.Error + icon.Space + "No processes running"

	FilePanelTopDirectoryIcon = FilePanelTopDirectoryIconStyle.Render(" " + icon.Directory + icon.Space)
	FilePanelNoneText = FilePanelStyle.Render(" " + icon.Error + icon.Space + "No such file or directory")

	// TODO : This "---" being appended before and after should be done via a function
	FilePreviewNoContentText = "\n--- " + icon.Error + icon.Space + "No content to preview" + icon.Space + "---"
	FilePreviewNoFileInfoText = "\n--- " + icon.Error + icon.Space + "Could not get file info" + icon.Space + "---"
	FilePreviewUnsupportedFormatText = "\n--- " + icon.Error + icon.Space + "Unsupported formats" + icon.Space + "---"
	FilePreviewUnsupportedFileMode = "\n--- " + icon.Error + icon.Space + "Unsupported File Mode" + icon.Space + "---"
	FilePreviewDirectoryUnreadableText = "\n--- " + icon.Error + icon.Space + "Cannot read directory" + icon.Space + "---"
	FilePreviewError = "\n--- " + icon.Error + icon.Space + "Error" + icon.Space + "---"
	FilePreviewEmptyText = "\n--- Empty ---"

	CheckboxChecked = FilePanelSelectBoxStyle.
		Foreground(FilePanelBorderColor).
		Render(icon.CheckboxChecked + icon.Space)
	CheckboxCheckedFocused = FilePanelSelectBoxStyle.
		Foreground(FilePanelBorderActiveColor).
		Render(icon.CheckboxChecked + icon.Space)
	CheckboxEmpty = FilePanelSelectBoxStyle.
		Foreground(FilePanelBorderColor).
		Render(icon.CheckboxEmpty + icon.Space)
	CheckboxEmptyFocused = FilePanelSelectBoxStyle.
		Foreground(FilePanelBorderActiveColor).
		Render(icon.CheckboxEmpty + icon.Space)

	ModalOkayInputText = MainStyle.AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).Render(
		ModalConfirm.Render(" (" + Hotkeys.Confirm[0] + ") Okay "))
	ModalConfirmInputText = ModalConfirm.Render(" (" + Hotkeys.Confirm[0] + ") Confirm ")
	ModalCancelInputText = ModalCancel.Render(" (" + Hotkeys.Quit[0] + ") Cancel ")
	ModalInputSpacingText = lipgloss.NewStyle().Background(ModalBGColor).Render("           ")
}
