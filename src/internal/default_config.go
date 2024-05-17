package internal

var (
	HotkeysTomlString string
	ConfigTomlString  string
	DefaultThemeString string
)

func defaultModelConfig(toggleDotFileBool bool, firstFilePanelDir string) model {
	return model{
		filePanelFocusIndex: 0,
		focusPanel:          nonePanelFocus,
		processBarModel: processBarModel{
			process: make(map[string]process),
			cursor:  0,
			render:  0,
		},
		sidebarModel: sidebarModel{
			renderIndex: 0,
			directories: getDirectories(),
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{
					render:          0,
					cursor:          0,
					location:        firstFilePanelDir,
					panelMode:       browserMode,
					focusType:       focus,
					directoryRecord: make(map[string]directoryRecord),
					searchBar:       generateSearchBar(),
				},
			},
			width: 10,
		},
		helpMenu: helpMenuModal{
			renderIndex: 0,
			cursor:      1,
			data:        getHelpMenuData(),
			open:        false,
		},
		toggleDotFile: toggleDotFileBool,
	}
}

func getHelpMenuData() []helpMenuModalData {
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
			hotkey:         hotkeys.Confirm,
			description:    "Confirm your select or typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.Quit,
			description:    "Quit typing, modal or superfile",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ConfirmTyping,
			description:    "Confirm typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.CancelTyping,
			description:    "Cancel typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.OpenHelpMenu,
			description:    "Open help menu(hotkeylist)",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel navigation",
		},
		{
			hotkey:         hotkeys.CreateNewFilePanel,
			description:    "Create new file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.CloseFilePanel,
			description:    "Close the focused file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.NextFilePanel,
			description:    "Focus on the next file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.PreviousFilePanel,
			description:    "Focus on the previous file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FocusOnProcessBar,
			description:    "Focus on the processbar panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FocusOnSidebar,
			description:    "Focus on the sidebar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FocusOnMetaData,
			description:    "Focus on the metadata panel",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel movement",
		},
		{
			hotkey:         hotkeys.ListUp,
			description:    "Up",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ListDown,
			description:    "Down",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ParentDirectory,
			description:    "Return to parent folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FilePanelSelectAllItem,
			description:    "Select all items in focused file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FilePanelSelectModeItemsSelectUp,
			description:    "Select up with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FilePanelSelectModeItemsSelectDown,
			description:    "Select down with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ToggleDotFile,
			description:    "Toggle dot file display",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.SearchBar,
			description:    "Toggle active search bar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ChangePanelMode,
			description:    "Change between selection mode or normal mode",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.PinnedDirectory,
			description:    "Pin or Unpin folder to sidebar (can be auto saved)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ToggleDotFile,
			description:    "Toggle to show dotfiles or not",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "File operations",
		},
		{
			hotkey:         hotkeys.FilePanelItemCreate,
			description:    "Create file or folder(/ ends with creating a folder)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FilePanelItemRename,
			description:    "Rename file or folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.CopyItems,
			description:    "Copy selected items to the clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.PasteItems,
			description:    "Cut selected items to the clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.PasteItems,
			description:    "Paste clipboard items into the current file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.PasteItems,
			description:    "Delete selected items",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ExtractFile,
			description:    "Extract compressed file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.CompressFile,
			description:    "Zip file or folder to .zip file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.OpenFileWithEditor,
			description:    "Open file with your default editor",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.OpenCurrentDirectoryWithEditor,
			description:    "Open current directory with default editor",
			hotkeyWorkType: normalType,
		},
	}

	return data
}
