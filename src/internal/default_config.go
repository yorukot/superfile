package internal

import (
	"path/filepath"

	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/preview"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
	zoxideui "github.com/yorukot/superfile/src/internal/ui/zoxide"
)

// Generate and return model containing default configurations for interface
// Maybe we can replace slice of strings with var args - Should we ?
// TODO: Move the configuration parameters to a ModelConfig struct.
// Something like `RendererConfig` struct for `Renderer` struct in ui/renderer package
func defaultModelConfig(toggleDotFile, toggleFooter, firstUse bool,
	firstPanelPaths []string, zClient *zoxidelib.Client) *model {
	return &model{
		filePanelFocusIndex: 0,
		focusPanel:          nonePanelFocus,
		processBarModel:     processbar.New(),
		sidebarModel:        sidebar.New(),
		fileMetaData:        metadata.New(),
		fileModel: fileModel{
			filePanels:  filePanelSlice(firstPanelPaths),
			filePreview: preview.New(),
			width:       common.DefaultFilePanelWidth,
		},
		helpMenu:       newHelpMenuModal(),
		promptModal:    prompt.DefaultModel(prompt.PromptMinHeight, prompt.PromptMinWidth),
		zoxideModal:    zoxideui.DefaultModel(zoxideui.ZoxideMinHeight, zoxideui.ZoxideMinWidth, zClient),
		zClient:        zClient,
		modelQuitState: notQuitting,
		toggleDotFile:  toggleDotFile,
		toggleFooter:   toggleFooter,
		firstUse:       firstUse,
		hasTrash:       common.InitTrash(),
	}
}

func newHelpMenuModal() helpMenuModal {
	helpMenuData := getHelpMenuData()

	return helpMenuModal{
		renderIndex:  0,
		cursor:       1,
		data:         helpMenuData,
		filteredData: helpMenuData,
		open:         false,
		searchBar:    common.GenerateSearchBar(),
	}
}

// Return help menu for Hotkeys
func getHelpMenuData() []helpMenuModalData { //nolint: funlen // This should be self contained
	data := []helpMenuModalData{
		{
			subTitle: "General",
		},
		{
			hotkey:         []string{"spf", ""},
			description:    "Open superfile",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.Confirm,
			description:    "Confirm your select or typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.Quit,
			description:    "Quit typing, modal or superfile",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CdQuit,
			description:    "Quit superfile and change directory to current folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ConfirmTyping,
			description:    "Confirm typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CancelTyping,
			description:    "Cancel typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenHelpMenu,
			description:    "Open help menu (hotkeylist)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenCommandLine,
			description:    "Open command line",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenSPFPrompt,
			description:    "Open SPF prompt",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenZoxide,
			description:    "Open zoxide navigation",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel navigation",
		},
		{
			hotkey:         common.Hotkeys.CreateNewFilePanel,
			description:    "Create new file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CloseFilePanel,
			description:    "Close the focused file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleFilePreviewPanel,
			description:    "Toggle file preview panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenSortOptionsMenu,
			description:    "Open sort options menu",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleReverseSort,
			description:    "Toggle reverse sort",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleFooter,
			description:    "Toggle footer",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.NextFilePanel,
			description:    "Focus on the next file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PreviousFilePanel,
			description:    "Focus on the previous file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FocusOnProcessBar,
			description:    "Focus on the processbar panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FocusOnSidebar,
			description:    "Focus on the sidebar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FocusOnMetaData,
			description:    "Focus on the metadata panel",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel movement",
		},
		{
			hotkey:         common.Hotkeys.ListUp,
			description:    "Up",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ListDown,
			description:    "Down",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PageUp,
			description:    "Page up",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PageDown,
			description:    "Page down",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ParentDirectory,
			description:    "Return to parent folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelSelectAllItem,
			description:    "Select all items in focused file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelSelectModeItemsSelectUp,
			description:    "Select up with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelSelectModeItemsSelectDown,
			description:    "Select down with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleDotFile,
			description:    "Toggle dot file display",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.SearchBar,
			description:    "Toggle active search bar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ChangePanelMode,
			description:    "Change between selection mode or normal mode",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PinnedDirectory,
			description:    "Pin or Unpin folder to sidebar (can be auto saved)",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "File operations",
		},
		{
			hotkey:         common.Hotkeys.FilePanelItemCreate,
			description:    "Create file or folder(end with " + string(filepath.Separator) + " to create a folder)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelItemRename,
			description:    "Rename file or folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CopyItems,
			description:    "Copy selected items to the clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CutItems,
			description:    "Cut selected items to the clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PasteItems,
			description:    "Paste clipboard items into the current file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.DeleteItems,
			description:    "Delete selected items",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PermanentlyDeleteItems,
			description:    "Permanently delete selected items",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CopyPath,
			description:    "Copy current file or directory path",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CopyPWD,
			description:    "Copy current working directory",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ExtractFile,
			description:    "Extract compressed file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         common.Hotkeys.CompressFile,
			description:    "Zip file or folder to .zip file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         common.Hotkeys.OpenFileWithEditor,
			description:    "Open file with your default editor",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         common.Hotkeys.OpenCurrentDirectoryWithEditor,
			description:    "Open current directory with default editor",
			hotkeyWorkType: normalType,
		},
	}

	return data
}
