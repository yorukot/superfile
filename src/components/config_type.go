package components

import "github.com/charmbracelet/lipgloss"

var (
	filePanelBorderColor lipgloss.Color
	sidebarBorderColor lipgloss.Color
	footerBorderColor lipgloss.Color
	modalBorderColor lipgloss.Color

	filePanelBorderActiveColor lipgloss.Color
	sidebarBorderActiveColor lipgloss.Color
	footerBorderActiveColor lipgloss.Color
	modalBorderActiveColor lipgloss.Color

	fullScreenBGColor lipgloss.Color
	filePanelBGColor lipgloss.Color
	sidebarBGColor lipgloss.Color
	footerBGColor lipgloss.Color
	modalBGColor lipgloss.Color

	fullScreenFGColor lipgloss.Color
	filePanelFGColor lipgloss.Color
	sidebarFGColor lipgloss.Color
	footerFGColor lipgloss.Color
	modalFGColor lipgloss.Color

	cursorColor lipgloss.Color
	correctColor lipgloss.Color
	errorColor lipgloss.Color
	hintColor lipgloss.Color
	cancelColor lipgloss.Color
	warnColor lipgloss.Color
	
	filePanelTopDirectoryIconColor lipgloss.Color
	filePanelTopPathColor lipgloss.Color
	filePanelItemSelectedFGColor lipgloss.Color
	filePanelItemSelectedBGColor lipgloss.Color

	sidebarTitleColor lipgloss.Color
	sidebarItemSelectedFGColor lipgloss.Color
	sidebarItemSelectedBGColor lipgloss.Color
	
	ModalCancelFGColor lipgloss.Color
	ModalCancelBGColor lipgloss.Color
	ModalConfirmFGColor lipgloss.Color
	ModalConfirmBGColor lipgloss.Color
)
// Theme configuration
type ThemeType struct {
	// Border
    FilePanelBorder       string `toml:"file_panel_border"`
    SidebarBorder         string `toml:"sidebar_border"`
    FooterBorder          string `toml:"footer_border"`
    ModalBorder           string `toml:"modal_border"`

    // Border Active
    FilePanelBorderActive string `toml:"file_panel_border_active"`
    SidebarBorderActive   string `toml:"sidebar_border_active"`
    FooterBorderActive    string `toml:"footer_border_active"`
    ModalBorderActive     string `toml:"modal_border_active"`

    // Background (bg)
    FullScreenBG  string `toml:"full_screen_bg"`
    FilePanelBG   string `toml:"file_panel_bg"`
    SidebarBG     string `toml:"sidebar_bg"`
    FooterBG      string `toml:"footer_bg"`
    ModalBG       string `toml:"modal_bg"`

    // Foreground (fg)
    FullScreenFG  string `toml:"full_screen_fg"`
    FilePanelFG   string `toml:"file_panel_fg"`
    SidebarFG     string `toml:"sidebar_fg"`
    FooterFG      string `toml:"footer_fg"`
    ModalFG       string `toml:"modal_fg"`

    // Special Color
    Cursor        string `toml:"cursor"`
    Correct       string `toml:"correct"`
    Error         string `toml:"error"`
    Hint          string `toml:"hint"`
    Cancel        string `toml:"cancel"`
    Warn          string `toml:"warn"`
    GradientColor []string `toml:"gradient_color"`

    // File Panel Special Items
    FilePanelTopDirectoryIcon string `toml:"file_panel_top_directory_icon"`
    FilePanelTopPath          string `toml:"file_panel_top_path"`
    FilePanelItemSelectedFG   string `toml:"file_panel_item_selected_fg"`
    FilePanelItemSelectedBG   string `toml:"file_panel_item_selected_bg"`

    // Sidebar Special Items
    SidebarTitle             string `toml:"sidebar_title"`
    SidebarItemSelectedFG    string `toml:"sidebar_item_selected_fg"`
    SidebarItemSelectedBG    string `toml:"sidebar_item_selected_bg"`

    // Modal Special Items
    ModalCancelFG            string `toml:"modal_cancel_fg"`
    ModalCancelBG            string `toml:"modal_cancel_bg"`
    ModalConfirmFG           string `toml:"modal_confirm_fg"`
    ModalConfirmBG           string `toml:"modal_confirm_bg"`

	Border         string `toml:"border"`
	MainBackground string `toml:"main_background"`

	TerminalTooSmallError string `toml:"terminal_too_small_error"`
	TerminalSizeCorrect   string `toml:"terminal_size_correct"`

	SidebarItem     string `toml:"sidebar_item"`
	SidebarSelected string `toml:"sidebar_selected"`
	SidebarFocus    string `toml:"sidebar_focus"`

	FilePanelFocus            string `toml:"file_panel_focus"`
	FilePanelItem             string `toml:"file_panel_item"`

	FilePanelItemSelected string    `toml:"file_panel_item_selected"`
	FooterFocus           string    `toml:"footer_focus"`
	ProcessBarGradient    [2]string `toml:"process_bar_gradient"`
	InOperation           string    `toml:"in_operation"`
	Done                  string    `toml:"done"`
	Fail                  string    `toml:"fail"`
	ModalForeground       string    `toml:"modal_foreground"`
	ModalCancel           string    `toml:"modal_cancel"`
	ModalConfirm          string    `toml:"modal_confirm"`
}

// Configuration settings
type ConfigType struct {
	Theme string `toml:"theme"`
	FooterPanelList []string `toml:"footer_panel_list"`
	Metadata bool `toml:"metadata"`
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
