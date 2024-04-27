package components

var HotkeysTomlString string = `# Here is global, all global key cant conflicts with other hotkeys
quit = ['esc', 'q']
# 
list_up = ['up', 'k']
list_down = ['down', 'j']
# 
pinned_directory = ['ctrl+p', '']
# 
close_file_panel = ['ctrl+w', '']
create_new_file_panel = ['ctrl+n', '']
# 
next_file_panel = ['tab', 'L']
previous_file_panel = ['shift+left', 'H']
focus_on_process_bar = ['p', '']
focus_on_side_bar = ['b', '']
focus_on_meta_data = ['m', '']
# 
change_panel_mode = ['v', '']
open_help_menu = ['?', '']
# 
file_panel_directory_create = ['f', '']
file_panel_file_create = ['c', '']
file_panel_item_rename = ['r', '']
paste_item = ['ctrl+v', '']
extract_file = ['ctrl+e', '']
compress_file = ['ctrl+r', '']
toggle_dot_file = ['ctrl+h', '']
# 
oepn_file_with_editor = ['e', '']
open_current_directory_with_editor = ['E', '']
# 
# These hotkeys do not conflict with any other keys (including global hotkey)
cancel = ['ctrl+c', 'esc']
confirm = ['enter', '']
# 
# Here is normal mode hotkey you can conflicts with other mode (cant conflicts with global hotkey)
delete_item = ['ctrl+d', '']
select_item = ['enter', 'l']
parent_directory = ['h', 'backspace']
copy_single_item = ['ctrl+c', '']
cut_single_item = ['ctrl+x', '']
search_bar = ['ctrl+f', '']
# 
# Here is select mode hotkey you can conflicts with other mode (cant conflicts with global hotkey)
file_panel_select_mode_item_single_select = ['enter', 'l']
file_panel_select_mode_item_select_down = ['shift+down', 'J']
file_panel_select_mode_item_select_up = ['shift+up', 'K']
file_panel_select_mode_item_delete = ['ctrl+d', 'delete']
file_panel_select_mode_item_copy = ['ctrl+c', '']
file_panel_select_mode_item_cut = ['ctrl+x', '']
file_panel_select_all_item = ['ctrl+a', '']
`

var ConfigTomlString string = `# change your theme
theme = 'catpuccin'
# 
# useless for now
footer_panel_list = ['processes', 'metadata', 'clipboard']
# 
# ==========PLUGINS========== #
# 
# Show more detailed metadata, please install exiftool before enabling this plugin!
metadata = false
`

func getHelpMenuData() []helpMenuModalData {
	data := []helpMenuModalData{
		{
			subTitle: "General",
		},
		{
			hotkey:         hotkeys.Quit,
			description:    "Quit",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel navigation",
		},
		{
			hotkey:         hotkeys.PinnedDirectory,
			description:    "Pin or Unpin folder to sidebar (can be auto saved)",
			hotkeyWorkType: globalType,
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
			hotkey:         hotkeys.FocusOnSideBar,
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
			hotkey:         hotkeys.ChangePanelMode,
			description:    "Change between selection mode or normal mode",
			hotkeyWorkType: globalType,
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
			hotkey:         hotkeys.SelectItem,
			description:    "Go to folder",
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
			hotkey:         hotkeys.FilePanelSelectModeItemSelectUp,
			description:    "Select with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.SelectItem,
			description:    "Select with your course",
			hotkeyWorkType: selectType,
		},
		{
			hotkey:         hotkeys.FilePanelSelectModeItemSingleSelect,
			description:    "Select the item where the current cursor is located",
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
			subTitle: "File operations",
		},
		{
			hotkey:         hotkeys.FilePanelDirectoryCreate,
			description:    "Create a new folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FilePanelFileCreate,
			description:    "Create a new file",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.FilePanelItemRename,
			description:    "Rename file or folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.ExtractFile,
			description:    "Extract zip file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.CompressFile,
			description:    "Zip file or folder to .zip file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.DeleteItem,
			description:    "Delete file or folder (or both)",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.CopySingleItem,
			description:    "Copy file or folder (or both)",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         hotkeys.FilePanelSelectModeItemCut,
			description:    "Cut file or folder (or both)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.PasteItem,
			description:    "Paste all items in your clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.SelectItem,
			description:    "Open file with your default application",
			hotkeyWorkType: selectType,
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
		{
			subTitle: "Special",
		},
		{
			hotkey:         hotkeys.Confirm,
			description:    "Confirm rename or create item or exit search bar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         hotkeys.Cancel,
			description:    "Cancel rename or create item or exit search bar and clear search bar value",
			hotkeyWorkType: globalType,
		},
	}

	return data
}