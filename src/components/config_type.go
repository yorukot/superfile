package components

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
}

// Configuration settings
type ConfigType struct {
	Theme           string   `mapstructure:"theme"`
	FooterPanelList []string `mapstructure:"footer_panel_list"`
	Metadata        bool     `mapstructure:"metadata"`
}

type HotkeysType struct {
	Quit []string `mapstructure:"quit"`

	ListUp   []string `mapstructure:"list_up"`
	ListDown []string `mapstructure:"list_down"`

	PinnedDirectory []string `mapstructure:"pinned_directory"`

	CloseFilePanel           []string `mapstructure:"close_file_panel"`
	CreateNewFilePanel       []string `mapstructure:"create_new_file_panel"`
	NextFilePanel            []string `mapstructure:"next_file_panel"`
	PreviousFilePanel        []string `mapstructure:"previous_file_panel"`
	FocusOnProcessBar        []string `mapstructure:"focus_on_process_bar"`
	FocusOnSideBar           []string `mapstructure:"focus_on_side_bar"`
	FocusOnMetaData          []string `mapstructure:"focus_on_meta_data"`
	ChangePanelMode          []string `mapstructure:"change_panel_mode"`
	FilePanelDirectoryCreate []string `mapstructure:"file_panel_directory_create"`
	FilePanelFileCreate      []string `mapstructure:"file_panel_file_create"`
	FilePanelItemRename      []string `mapstructure:"file_panel_item_rename"`
	PasteItem                []string `mapstructure:"paste_item"`
	ExtractFile              []string `mapstructure:"extract_file"`
	CompressFile             []string `mapstructure:"compress_file"`
	ToggleDotFile            []string `mapstructure:"toggle_dot_file"`

	OpenFileWithEditor             []string `mapstructure:"oepn_file_with_editor"`
	OpenCurrentDirectoryWithEditor []string `mapstructure:"open_current_directory_with_editor"`

	Cancel  []string `mapstructure:"cancel"`
	Confirm []string `mapstructure:"confirm"`

	DeleteItem      []string `mapstructure:"delete_item"`
	SelectItem      []string `mapstructure:"select_item"`
	ParentDirectory []string `mapstructure:"parent_directory"`
	CopySingleItem  []string `mapstructure:"copy_single_item"`
	CutSingleItem   []string `mapstructure:"cut_single_item"`
	SearchBar       []string `mapstructure:"search_bar"`

	FilePanelSelectModeItemSingleSelect []string `mapstructure:"file_panel_select_mode_item_single_select"`
	FilePanelSelectModeItemSelectDown   []string `mapstructure:"file_panel_select_mode_item_select_down"`
	FilePanelSelectModeItemSelectUp     []string `mapstructure:"file_panel_select_mode_item_select_up"`
	FilePanelSelectModeItemDelete       []string `mapstructure:"file_panel_select_mode_item_delete"`
	FilePanelSelectModeItemCopy         []string `mapstructure:"file_panel_select_mode_item_copy"`
	FilePanelSelectModeItemCut          []string `mapstructure:"file_panel_select_mode_item_cut"`
	FilePanelSelectAllItem              []string `mapstructure:"file_panel_select_all_item"`
}
