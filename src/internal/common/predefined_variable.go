package common

import (
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
)

const (
	WheelRunTime          = 5
	DefaultCommandTimeout = 5000 * time.Millisecond
	DateModifiedOption    = "Date Modified"
	InvalidTypeString     = "InvalidType"
)

const (
	SameRenameWarnTitle   = "There is already a file or directory with that name"
	SameRenameWarnContent = "This operation will override the existing file"
)

const (
	TrashWarnTitle             = "Are you sure you want to move this to trash can"
	TrashWarnContent           = "This operation will move file or directory to trash can."
	PermanentDeleteWarnTitle   = "Are you sure you want to completely delete"
	PermanentDeleteWarnContent = "This operation cannot be undone and your data will be completely lost."
)

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
	ClipboardNoneText  string

	FilePanelTopDirectoryIcon string
	FilePanelNoneText         string

	FilePreviewNoFileInfoText               string
	FilePreviewNoContentText                string
	FilePreviewUnsupportedFormatText        string
	FilePreviewUnsupportedFileMode          string
	FilePreviewDirectoryUnreadableText      string
	FilePreviewEmptyText                    string
	FilePreviewError                        string
	FilePreviewPanelClosedText              string
	FilePreviewImagePreviewDisabledText     string
	FilePreviewUnsupportedImageFormatsText  string
	FilePreviewImageConversionErrorText     string
	FilePreviewBatNotInstalledText          string
	FilePreviewThumbnailGenerationErrorText string

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
	UnsupportedPreviewFormats = []string{".torrent"}
	ImageExtensions           = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".svg":  true,
		".webp": true,
		".ico":  true,
	}
	VideoExtensions = map[string]bool{
		".mkv":  true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".flv":  true,
		".webm": true,
		".wmv":  true,
		".m4v":  true,
		".mpeg": true,
		".3gp":  true,
		".ogv":  true,
	}
)

// No dependencies
func LoadInitialPrerenderedVariables() {
	LipglossError = lipgloss.NewStyle().Foreground(lipgloss.Color("#F93939")).Render("Error") +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFEE")).Render(" ┃ ")
}

// This should be used only after InitIcon() has been called.
func wrapFilePreviewErrorMsg(msg string) string {
	return "\n--- " + icon.Error + icon.Space + msg + " ---"
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
	ClipboardNoneText = " " + icon.Error + icon.Space + " No content in clipboard"

	FilePanelTopDirectoryIcon = FilePanelTopDirectoryIconStyle.Render(" " + icon.Directory + icon.Space)
	FilePanelNoneText = FilePanelStyle.Render(" " + icon.Error + icon.Space + "No such file or directory")

	FilePreviewNoContentText = wrapFilePreviewErrorMsg(
		"No content to preview")
	FilePreviewNoFileInfoText = wrapFilePreviewErrorMsg(
		"Could not get file info")
	FilePreviewUnsupportedFormatText = wrapFilePreviewErrorMsg(
		"Unsupported formats")
	FilePreviewUnsupportedFileMode = wrapFilePreviewErrorMsg(
		"Unsupported File Mode")
	FilePreviewDirectoryUnreadableText = wrapFilePreviewErrorMsg(
		"Cannot read directory")
	FilePreviewError = wrapFilePreviewErrorMsg(
		"Error")
	FilePreviewEmptyText = wrapFilePreviewErrorMsg(
		"Empty")
	FilePreviewPanelClosedText = wrapFilePreviewErrorMsg(
		"Preview panel is closed")
	FilePreviewImagePreviewDisabledText = wrapFilePreviewErrorMsg(
		"Image preview is disabled")
	FilePreviewUnsupportedImageFormatsText = wrapFilePreviewErrorMsg(
		"Unsupported image formats")
	FilePreviewImageConversionErrorText = wrapFilePreviewErrorMsg(
		"Error convert image to ansi")
	FilePreviewBatNotInstalledText = wrapFilePreviewErrorMsg(
		"'bat' is not installed or not found")
	FilePreviewThumbnailGenerationErrorText = wrapFilePreviewErrorMsg(
		"Thumbnail generation failed")

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
		ModalConfirm.Render(" (" + Hotkeys.ConfirmTyping[0] + ") Okay "))
	ModalConfirmInputText = ModalConfirm.Render(" (" + Hotkeys.ConfirmTyping[0] + ") Confirm ")
	ModalCancelInputText = ModalCancel.Render(" (" + Hotkeys.Quit[0] + ") Cancel ")
	ModalInputSpacingText = lipgloss.NewStyle().Background(ModalBGColor).Render("           ")
}
