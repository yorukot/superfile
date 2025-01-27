package internal

import "github.com/charmbracelet/lipgloss"

var (
	wheelRunTime = 5
)

var (
	filePanelBorderColor lipgloss.Color
	sidebarBorderColor   lipgloss.Color
	footerBorderColor    lipgloss.Color

	filePanelBorderActiveColor lipgloss.Color
	sidebarBorderActiveColor   lipgloss.Color
	footerBorderActiveColor    lipgloss.Color
	modalBorderActiveColor     lipgloss.Color

	fullScreenBGColor lipgloss.Color
	filePanelBGColor  lipgloss.Color
	sidebarBGColor    lipgloss.Color
	footerBGColor     lipgloss.Color
	modalBGColor      lipgloss.Color

	fullScreenFGColor lipgloss.Color
	filePanelFGColor  lipgloss.Color
	sidebarFGColor    lipgloss.Color
	footerFGColor     lipgloss.Color
	modalFGColor      lipgloss.Color

	cursorColor  lipgloss.Color
	correctColor lipgloss.Color
	errorColor   lipgloss.Color
	hintColor    lipgloss.Color
	cancelColor  lipgloss.Color

	filePanelTopDirectoryIconColor lipgloss.Color
	filePanelTopPathColor          lipgloss.Color
	filePanelItemSelectedFGColor   lipgloss.Color
	filePanelItemSelectedBGColor   lipgloss.Color

	sidebarTitleColor          lipgloss.Color
	sidebarItemSelectedFGColor lipgloss.Color
	sidebarItemSelectedBGColor lipgloss.Color
	sidebarDividerColor        lipgloss.Color

	modalCancelFGColor  lipgloss.Color
	modalCancelBGColor  lipgloss.Color
	modalConfirmFGColor lipgloss.Color
	modalConfirmBGColor lipgloss.Color

	helpMenuHotkeyColor lipgloss.Color
	helpMenuTitleColor  lipgloss.Color
)

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
	Cursor        string   `toml:"cursor"`
	Correct       string   `toml:"correct"`
	Error         string   `toml:"error"`
	Hint          string   `toml:"hint"`
	Cancel        string   `toml:"cancel"`
	GradientColor []string `toml:"gradient_color"`

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
	Theme string `toml:"theme" comment:"More details are at https://superfile.netlify.app/configure/superfile-config/\nchange your theme"`

	Editor                 string `toml:"editor" comment:"\nThe editor files/directories will be opened with. (leave blank to use the EDITOR environment variable)."`
	AutoCheckUpdate        bool   `toml:"auto_check_update" comment:"\nAuto check for update"`
	CdOnQuit               bool   `toml:"cd_on_quit" comment:"\nCd on quit (For more details, please check out https://superfile.netlify.app/configure/superfile-config/#cd_on_quit)"`
	DefaultOpenFilePreview bool   `toml:"default_open_file_preview" comment:"\nWhether to open file preview automatically every time superfile is opened."`
	DefaultDirectory       string `toml:"default_directory" comment:"\nThe path of the first file panel when superfile is opened."`
	FileSizeUseSI          bool   `toml:"file_size_use_si" comment:"\nDisplay file sizes using powers of 1000 (kB, MB, GB) instead of powers of 1024 (KiB, MiB, GiB)."`
	DefaultSortType        int    `toml:"default_sort_type" comment:"\nDefault sort type (0: Name, 1: Size, 2: Date Modified)."`
	SortOrderReversed      bool   `toml:"sort_order_reversed" comment:"\nDefault sort order (false: Ascending, true: Descending)."`
	CaseSensitiveSort      bool   `toml:"case_sensitive_sort" comment:"\nCase sensitive sort by name (captal \"B\" comes before \"a\" if true)."`
	Debug                  bool   `toml:"debug" comment:"\nWhether to enable debug mode."`

	Nerdfont              bool `toml:"nerdfont" comment:"\n================   Style =================\n\n If you don't have or don't want Nerdfont installed you can turn this off"`
	TransparentBackground bool `toml:"transparent_background" comment:"\nSet transparent background or not (this only work when your terminal background is transparent)"`
	FilePreviewWidth      int  `toml:"file_preview_width" comment:"\nFile preview width allow '0' (this mean same as file panel),'x' x must be less than 10 and greater than 1 (This means that the width of the file preview will be one xth of the total width.)"`
	SidebarWidth          int  `toml:"sidebar_width" comment:"\nThe length of the sidebar. If you don't want to display the sidebar, you can input 0 directly. If you want to display the value, please place it in the range of 3-20."`

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

	Metadata          bool `toml:"metadata" comment:"\n==========PLUGINS========== #\n\nShow more detailed metadata, please install exiftool before enabling this plugin!"`
	EnableMD5Checksum bool `toml:"enable_md5_checksum" comment:"Enable MD5 checksum generation for files"`
}

type HotkeysType struct {
	Confirm []string `toml:"confirm" comment:"=================================================================================================\nGlobal hotkeys (cannot conflict with other hotkeys)"`
	Quit    []string `toml:"quit"`
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

	CopyItems   []string `toml:"copy_items" comment:"file operate"`
	PasteItems  []string `toml:"paste_items"`
	CutItems    []string `toml:"cut_items"`
	DeleteItems []string `toml:"delete_items"`

	ExtractFile  []string `toml:"extract_file" comment:"compress and extract"`
	CompressFile []string `toml:"compress_file"`

	OpenFileWithEditor             []string `toml:"open_file_with_editor" comment:"editor"`
	OpenCurrentDirectoryWithEditor []string `toml:"open_current_directory_with_editor"`

	PinnedDirectory []string `toml:"pinned_directory" comment:"other"`
	ToggleDotFile   []string `toml:"toggle_dot_file"`
	ChangePanelMode []string `toml:"change_panel_mode"`
	OpenHelpMenu    []string `toml:"open_help_menu"`
	OpenCommandLine []string `toml:"open_command_line"`

	CopyPath []string `toml:"copy_path"`
	CopyPWD  []string `toml:"copy_present_working_directory"`

	ToggleFooter []string `toml:"toggle_footer"`

	ConfirmTyping []string `toml:"confirm_typing" comment:"=================================================================================================\nTyping hotkeys (can conflict with all hotkeys)"`
	CancelTyping  []string `toml:"cancel_typing"`

	ParentDirectory []string `toml:"parent_directory" comment:"=================================================================================================\nNormal mode hotkeys (can conflict with other modes, cannot conflict with global hotkeys)"`
	SearchBar       []string `toml:"search_bar"`

	FilePanelSelectModeItemsSelectDown []string `toml:"file_panel_select_mode_items_select_down" comment:"=================================================================================================\nSelect mode hotkeys (can conflict with other modes, cananot conflict with global hotkeys)"`
	FilePanelSelectModeItemsSelectUp   []string `toml:"file_panel_select_mode_items_select_up"`
	FilePanelSelectAllItem             []string `toml:"file_panel_select_all_items"`
}
