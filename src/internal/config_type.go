package internal

import "github.com/charmbracelet/lipgloss"

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
	Theme string `toml:"theme" comment:"change your theme"`

	FooterPanelList []string `toml:"footer_panel_list" comment:"\nuseless for now"`

	TransparentBackground bool `toml:"transparent_backgroun" comment:"\n================   Style =================\n\nSet transparent background or not (this only work when your terminal background is transparent)"`

	BorderTop         string `toml:"border_top" comment:"\n\nBorder style"`
	BorderBottom      string `toml:"border_bottom"`
	BorderLeft        string `toml:"border_left"`
	BorderRight       string `toml:"border_right"`
	BorderTopLeft     string `toml:"border_top_left"`
	BorderTopRight    string `toml:"border_top_right"`
	BorderBottomLeft  string `toml:"border_bottom_left"`
	BorderBottomRight string `toml:"border_bottom_right"`
	BorderMiddleLeft  string `toml:"border_middle_left"`
	BorderMiddleRight string `toml:"border_middle_right"`

	Metadata bool `toml:"metadata" comment:"\n==========PLUGINS========== #\n\nShow more detailed metadata, please install exiftool before enabling this plugin!"`
}

type HotkeysType struct {
	Quit []string `toml:"quit" comment:"Here is global, all global key cant conflicts with other hotkeys"`

	ListUp   []string `toml:"list_up" comment:"\n"`
	ListDown []string `toml:"list_down"`

	PinnedDirectory []string `toml:"pinned_directory" comment:"\n"`

	CloseFilePanel     []string `toml:"close_file_panel" comment:"\n"`
	CreateNewFilePanel []string `toml:"create_new_file_panel"`

	NextFilePanel     []string `toml:"next_file_panel" comment:"\n"`
	PreviousFilePanel []string `toml:"previous_file_panel"`
	FocusOnProcessBar []string `toml:"focus_on_process_bar"`
	FocusOnSideBar    []string `toml:"focus_on_side_bar"`
	FocusOnMetaData   []string `toml:"focus_on_metadata"`

	ChangePanelMode []string `toml:"change_panel_mode" comment:"\n"`
	OpenHelpMenu    []string `toml:"open_help_menu"`

	FilePanelDirectoryCreate []string `toml:"file_panel_directory_create" comment:"\n"`
	FilePanelFileCreate      []string `toml:"file_panel_file_create"`
	FilePanelItemRename      []string `toml:"file_panel_item_rename"`
	PasteItem                []string `toml:"paste_item"`
	ExtractFile              []string `toml:"extract_file"`
	CompressFile             []string `toml:"compress_file"`
	ToggleDotFile            []string `toml:"toggle_dot_file"`

	OpenFileWithEditor             []string `toml:"oepn_file_with_editor" comment:"\n"`
	OpenCurrentDirectoryWithEditor []string `toml:"open_current_directory_with_editor"`

	Cancel  []string `toml:"cancel" comment:"\nThese hotkeys do not conflict with any other keys (including global hotkey)"`
	Confirm []string `toml:"confirm"`

	DeleteItem      []string `toml:"delete_item" comment:"\nHere is normal mode hotkey you can conflicts with other mode (cant conflicts with global hotkey)"`
	SelectItem      []string `toml:"select_item"`
	ParentDirectory []string `toml:"parent_directory"`
	CopySingleItem  []string `toml:"copy_single_item"`
	CutSingleItem   []string `toml:"cut_single_item"`
	SearchBar       []string `toml:"search_bar"`
	CommandLine     []string `toml:"command_line"`

	FilePanelSelectModeItemSingleSelect []string `toml:"file_panel_select_mode_item_single_select" comment:"\nHere is select mode hotkey you can conflicts with other mode (cant conflicts with global hotkey)"`
	FilePanelSelectModeItemSelectDown   []string `toml:"file_panel_select_mode_item_select_down"`
	FilePanelSelectModeItemSelectUp     []string `toml:"file_panel_select_mode_item_select_up"`
	FilePanelSelectModeItemDelete       []string `toml:"file_panel_select_mode_item_delete"`
	FilePanelSelectModeItemCopy         []string `toml:"file_panel_select_mode_item_copy"`
	FilePanelSelectModeItemCut          []string `toml:"file_panel_select_mode_item_cut"`
	FilePanelSelectAllItem              []string `toml:"file_panel_select_all_item"`
}
