package components

// Theme configuration
type ThemeType struct {
	Border         string `toml:"border"`
	Cursor         string `toml:"cursor"`
	MainBackground string `toml:"main_background"`

	TerminalTooSmallError string `toml:"terminal_too_small_error"`
	TerminalSizeCorrect   string `toml:"terminal_size_correct"`

	SidebarTitle    string `toml:"sidebar_title"`
	SidebarItem     string `toml:"sidebar_item"`
	SidebarSelected string `toml:"sidebar_selected"`
	SidebarFocus    string `toml:"sidebar_focus"`

	FilePanelFocus            string `toml:"file_panel_focus"`
	FilePanelTopDirectoryIcon string `toml:"file_panel_top_directory_icon"`
	FilePanelTopPath          string `toml:"file_panel_top_path"`
	FilePanelItem             string `toml:"file_panel_item"`

	FilePanelItemSelected string    `toml:"file_panel_item_selected"`
	FooterFocus           string    `toml:"footer_focus"`
	ProcessBarGradient    [2]string `toml:"process_bar_gradient"`
	InOperation           string    `toml:"in_operation"`
	Done                  string    `toml:"done"`
	Fail                  string    `toml:"fail"`
	Cancel                string    `toml:"cancel"`
	ModalForeground       string    `toml:"modal_foreground"`
	ModalCancel           string    `toml:"modal_cancel"`
	ModalConfirm          string    `toml:"modal_confirm"`
}

// Configuration settings
type ConfigType struct {
	Theme string `toml:"theme"`
}

type HotkeysType struct {
	Quit []string `toml:"quit"`

	ListUp   []string `toml:"list_up"`
	ListDown []string `toml:"list_down"`

	PinnedDirectory []string `toml:"pinned_directory"`

	CloseFilePanel           []string `toml:"close_file_panel"`
	CreateNewFilePanel       []string `toml:"create_new_file_panel"`
	NextFilePanel            []string `toml:"next_file_panel"`
	PreviousFilePanel        []string `toml:"previous_file_panel"`
	FocusOnProcessBar        []string `toml:"focus_on_process_bar"`
	FocusOnSideBar           []string `toml:"focus_on_side_bar"`
	FocusOnMetaData          []string `toml:"focus_on_meta_data"`
	ChangePanelMode          []string `toml:"change_panel_mode"`
	FilePanelDirectoryCreate []string `toml:"file_panel_directory_create"`
	FilePanelFileCreate      []string `toml:"file_panel_file_create"`
	FilePanelItemRename      []string `toml:"file_panel_item_rename"`
	PasteItem                []string `toml:"paste_item"`
	ExtractFile              []string `toml:"extract_file"`
	CompressFile             []string `toml:"compress_file"`
	ToggleDotFile            []string `toml:"toggle_dot_file"`

	Cancel  []string `toml:"cancel"`
	Confirm []string `toml:"confirm"`

	DeleteItem      []string `toml:"delete_item"`
	SelectItem      []string `toml:"select_item"`
	ParentDirectory []string `toml:"parent_directory"`
	CopySingleItem  []string `toml:"copy_single_item"`
	CutSingleItem   []string `toml:"cut_single_item"`
	SearchBar		[]string `toml:"search_bar"`

	FilePanelSelectModeItemSingleSelect []string `toml:"file_panel_select_mode_item_single_select"`
	FilePanelSelectModeItemSelectDown   []string `toml:"file_panel_select_mode_item_select_down"`
	FilePanelSelectModeItemSelectUp     []string `toml:"file_panel_select_mode_item_select_up"`
	FilePanelSelectModeItemDelete       []string `toml:"file_panel_select_mode_item_delete"`
	FilePanelSelectModeItemCopy         []string `toml:"file_panel_select_mode_item_copy"`
	FilePanelSelectModeItemCut          []string `toml:"file_panel_select_mode_item_cut"`
	FilePanelSelectAllItem              []string `toml:"file_panel_select_all_item"`
}
