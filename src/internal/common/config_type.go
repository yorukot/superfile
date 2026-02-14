package common

// Theme configuration
type ThemeType struct {
	// Code syntax highlight theme
	CodeSyntaxHighlightTheme string `toml:"code_syntax_highlight"`

	// Border
	FilePanelBorder string `toml:"file_panel_border"`
	SidebarBorder   string `toml:"sidebar_border"`
	FooterBorder    string `toml:"footer_border"`

	// Border Active
	FilePanelBorderActive string `toml:"file_panel_border_active"`
	SidebarBorderActive   string `toml:"sidebar_border_active"`
	FooterBorderActive    string `toml:"footer_border_active"`
	ModalBorderActive     string `toml:"modal_border_active"`

	// Background (bg)
	FullScreenBG string `toml:"full_screen_bg"`
	FilePanelBG  string `toml:"file_panel_bg"`
	SidebarBG    string `toml:"sidebar_bg"`
	FooterBG     string `toml:"footer_bg"`
	ModalBG      string `toml:"modal_bg"`

	// Foreground (fg)
	FullScreenFG string `toml:"full_screen_fg"`
	FilePanelFG  string `toml:"file_panel_fg"`
	SidebarFG    string `toml:"sidebar_fg"`
	FooterFG     string `toml:"footer_fg"`
	ModalFG      string `toml:"modal_fg"`

	// Special Color
	Cursor  string `toml:"cursor"`
	Correct string `toml:"correct"`
	Error   string `toml:"error"`
	Hint    string `toml:"hint"`
	Cancel  string `toml:"cancel"`
	// Note: this is linked with `RequiredGradientColorCount` constant
	GradientColor      []string `toml:"gradient_color"`
	DirectoryIconColor string   `toml:"directory_icon_color"`

	// File Panel Special Items
	FilePanelTopDirectoryIcon string `toml:"file_panel_top_directory_icon"`
	FilePanelTopPath          string `toml:"file_panel_top_path"`
	FilePanelItemSelectedFG   string `toml:"file_panel_item_selected_fg"`
	FilePanelItemSelectedBG   string `toml:"file_panel_item_selected_bg"`

	// Sidebar Special Items
	SidebarTitle          string `toml:"sidebar_title"`
	SidebarItemSelectedFG string `toml:"sidebar_item_selected_fg"`
	SidebarItemSelectedBG string `toml:"sidebar_item_selected_bg"`
	SidebarDivider        string `toml:"sidebar_divider"`

	// Modal Special Items
	ModalCancelFG  string `toml:"modal_cancel_fg"`
	ModalCancelBG  string `toml:"modal_cancel_bg"`
	ModalConfirmFG string `toml:"modal_confirm_fg"`
	ModalConfirmBG string `toml:"modal_confirm_bg"`

	HelpMenuHotkey string `toml:"help_menu_hotkey"`
	HelpMenuTitle  string `toml:"help_menu_title"`
}

// Configuration settings
type ConfigType struct {
	Theme string `toml:"theme" comment:"More details are at https://superfile.dev/configure/superfile-config/\nchange your theme"`

	Editor    string `toml:"editor" comment:"\nThe editor files will be opened with. (Leave blank to use the EDITOR environment variable)."`
	DirEditor string `toml:"dir_editor" comment:"\nThe editor directories will be opened with. (Leave blank to use the default editors)."`
	// The table (map) for editor by file extension
	OpenWith map[string]string `toml:"open_with" comment:"\nCustom open commands by file extension."`

	AutoCheckUpdate        bool   `toml:"auto_check_update" comment:"\nAuto check for update"`
	CdOnQuit               bool   `toml:"cd_on_quit" comment:"\nCd on quit (For more details, please check out https://superfile.dev/configure/superfile-config/#cd_on_quit)"`
	DefaultOpenFilePreview bool   `toml:"default_open_file_preview" comment:"\nWhether to open file preview automatically every time superfile is opened."`
	ShowImagePreview       bool   `toml:"show_image_preview" comment:"\nWhether to show image preview."`
	ShowPanelFooterInfo    bool   `toml:"show_panel_footer_info" comment:"\nWhether to show additional footer info for file panel."`
	DefaultDirectory       string `toml:"default_directory" comment:"\nThe path of the first file panel when superfile is opened."`
	FileSizeUseSI          bool   `toml:"file_size_use_si" comment:"\nDisplay file sizes using powers of 1000 (kB, MB, GB) instead of powers of 1024 (KiB, MiB, GiB)."`
	DefaultSortType        int    `toml:"default_sort_type" comment:"\nDefault sort type (0: Name, 1: Size, 2: Date Modified, 3: Type)."`
	SortOrderReversed      bool   `toml:"sort_order_reversed" comment:"\nDefault sort order (false: Ascending, true: Descending)."`
	CaseSensitiveSort      bool   `toml:"case_sensitive_sort" comment:"\nCase sensitive sort by name (capital \"B\" comes before \"a\" if true)."`
	ShellCloseOnSuccess    bool   `toml:"shell_close_on_success" comment:"\nWhether to close the shell on successful command execution."`
	Debug                  bool   `toml:"debug" comment:"\nWhether to enable debug mode."`
	// IgnoreMissingFields controls whether warnings about missing TOML fields are suppressed.
	IgnoreMissingFields   bool `toml:"ignore_missing_fields" comment:"\nWhether to ignore warnings about missing fields in the config file."`
	PageScrollSize        int  `toml:"page_scroll_size" comment:"\nNumber of lines to scroll for PgUp/PgDown keys (0: full page, default behavior)."`
	FilePanelExtraColumns int  `toml:"file_panel_extra_columns" comment:"\nCount of extra columns in file panel in addition to file name. When option equal 0 then feature is disabled."`
	FilePanelNamePercent  int  `toml:"file_panel_name_percent" comment:"\nPercentage of file panel width allocated to file names (25-100). Higher values give more space to names, less to extra columns."`

	Nerdfont                bool     `toml:"nerdfont" comment:"\n================   Style =================\n\n If you don't have or don't want Nerdfont installed you can turn this off"`
	ShowSelectIcons         bool     `toml:"show_select_icons" comment:"\nShow checkbox icons in select mode (requires nerdfont)"`
	TransparentBackground   bool     `toml:"transparent_background" comment:"\nSet transparent background or not (this only work when your terminal background is transparent)"`
	FilePreviewWidth        int      `toml:"file_preview_width" comment:"\nFile preview width allow '0' (this mean same as file panel),'x' x must be less than 10 and greater than 1 (This means that the width of the file preview will be one xth of the total width.)"`
	EnableFilePreviewBorder bool     `toml:"enable_file_preview_border" comment:"\nEnable border around the file preview panel (default: false)"`
	CodePreviewer           string   `toml:"code_previewer" comment:"\nWhether to use the builtin syntax highlighting with chroma or use bat. Values: \"\" for builtin chroma, \"bat\" for bat"`
	SidebarWidth            int      `toml:"sidebar_width" comment:"\nThe length of the sidebar(excluding borders). If you don't find to display the sidebar, you can input 0 directly. If you want to display the value, please place it in the range of 5-20."`
	SidebarSections         []string `toml:"sidebar_sections" comment:"\nOrder of sidebar sections (valid values: \"home\", \"pinned\", \"disks\").\nOnly sections included in this list will be displayed."`

	BorderTop         string `toml:"border_top" comment:"\nBorder style"`
	BorderBottom      string `toml:"border_bottom"`
	BorderLeft        string `toml:"border_left"`
	BorderRight       string `toml:"border_right"`
	BorderTopLeft     string `toml:"border_top_left"`
	BorderTopRight    string `toml:"border_top_right"`
	BorderBottomLeft  string `toml:"border_bottom_left"`
	BorderBottomRight string `toml:"border_bottom_right"`
	BorderMiddleLeft  string `toml:"border_middle_left"`
	BorderMiddleRight string `toml:"border_middle_right"`

	Metadata          bool `toml:"metadata" comment:"\n==========PLUGINS========== #\nPlugins means that you need to install some external dependencies to use them.\n\nShow more detailed metadata, please install exiftool before enabling this plugin!"`
	EnableMD5Checksum bool `toml:"enable_md5_checksum" comment:"Enable MD5 checksum generation for files"`
	ZoxideSupport     bool `toml:"zoxide_support" comment:"Zoxide support for the fast navigation"`
}

// GetIgnoreMissingFields reports whether warnings about missing TOML fields should be ignored.
func (c *ConfigType) GetIgnoreMissingFields() bool {
	return c.IgnoreMissingFields
}

type HotkeysType struct {
	Confirm []string `toml:"confirm" comment:"=================================================================================================\nGlobal hotkeys (cannot conflict with other hotkeys)"`
	Quit    []string `toml:"quit"`
	CdQuit  []string `toml:"cd_quit"`

	// movement
	ListUp   []string `toml:"list_up" comment:"movement"`
	ListDown []string `toml:"list_down"`
	PageUp   []string `toml:"page_up"`
	PageDown []string `toml:"page_down"`

	CloseFilePanel         []string `toml:"close_file_panel" comment:"file panel control"`
	CreateNewFilePanel     []string `toml:"create_new_file_panel"`
	NextFilePanel          []string `toml:"next_file_panel"`
	PreviousFilePanel      []string `toml:"previous_file_panel"`
	ToggleFilePreviewPanel []string `toml:"toggle_file_preview_panel"`
	OpenSortOptionsMenu    []string `toml:"open_sort_options_menu"`
	ToggleReverseSort      []string `toml:"toggle_reverse_sort"`

	FocusOnProcessBar []string `toml:"focus_on_process_bar" comment:"change focus"`
	FocusOnSidebar    []string `toml:"focus_on_sidebar"`
	FocusOnMetaData   []string `toml:"focus_on_metadata"`

	FilePanelItemCreate []string `toml:"file_panel_item_create" comment:"create file/directory and rename "`
	FilePanelItemRename []string `toml:"file_panel_item_rename"`

	CopyItems              []string `toml:"copy_items" comment:"file operate"`
	PasteItems             []string `toml:"paste_items"`
	CutItems               []string `toml:"cut_items"`
	DeleteItems            []string `toml:"delete_items"`
	PermanentlyDeleteItems []string `toml:"permanently_delete_items"`

	ExtractFile  []string `toml:"extract_file" comment:"compress and extract"`
	CompressFile []string `toml:"compress_file"`

	OpenFileWithEditor             []string `toml:"open_file_with_editor" comment:"editor"`
	OpenCurrentDirectoryWithEditor []string `toml:"open_current_directory_with_editor"`

	PinnedDirectory []string `toml:"pinned_directory" comment:"other"`
	ToggleDotFile   []string `toml:"toggle_dot_file"`
	ChangePanelMode []string `toml:"change_panel_mode"`
	OpenHelpMenu    []string `toml:"open_help_menu"`
	OpenCommandLine []string `toml:"open_command_line"`
	OpenSPFPrompt   []string `toml:"open_spf_prompt"`
	OpenZoxide      []string `toml:"open_zoxide"`

	CopyPath []string `toml:"copy_path"`
	CopyPWD  []string `toml:"copy_present_working_directory"`

	ToggleFooter []string `toml:"toggle_footer"`

	ConfirmTyping []string `toml:"confirm_typing" comment:"=================================================================================================\nTyping hotkeys (can conflict with all hotkeys)"`
	CancelTyping  []string `toml:"cancel_typing"`

	ParentDirectory []string `toml:"parent_directory" comment:"=================================================================================================\nNormal mode hotkeys (can conflict with other modes, cannot conflict with global hotkeys)"`
	SearchBar       []string `toml:"search_bar"`

	FilePanelSelectModeItemsSelectDown []string `toml:"file_panel_select_mode_items_select_down" comment:"=================================================================================================\nSelect mode hotkeys (can conflict with other modes, cannot conflict with global hotkeys)"`
	FilePanelSelectModeItemsSelectUp   []string `toml:"file_panel_select_mode_items_select_up"`
	FilePanelSelectAllItem             []string `toml:"file_panel_select_all_items"`
}
