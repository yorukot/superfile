package internal

import (
	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/yorukot/superfile/src/internal/ui/preview"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
	zoxideui "github.com/yorukot/superfile/src/internal/ui/zoxide"
)

// Type representing the type of focused panel
type focusPanelType int

type hotkeyType int

type modelQuitStateType int

const (
	globalType hotkeyType = iota
	normalType
	selectType
)

// Constants for panel with no focus
const (
	nonePanelFocus focusPanelType = iota
	processBarFocus
	sidebarFocus
	metadataFocus
)

const (
	notQuitting modelQuitStateType = iota
	quitInitiated
	quitConfirmationInitiated
	quitConfirmationReceived
	quitDone
)

// Main model
// TODO : We could consider using *model as tea.Model, instead of model.
// for reducing re-allocations. The struct is 20K bytes. But this could lead to
// issues like race conditions and whatnot, which are hidden since we are creating
// new model in each tea update.
type model struct {
	// Main Panels
	fileModel       fileModel
	sidebarModel    sidebar.Model
	processBarModel processbar.Model
	focusPanel      focusPanelType
	copyItems       copyItems

	// Modals
	notifyModel notify.Model
	typingModal typingModal
	helpMenu    helpMenuModal
	promptModal prompt.Model
	zoxideModal zoxideui.Model

	// Zoxide client for directory tracking
	zClient *zoxidelib.Client

	fileMetaData         metadata.Model
	ioReqCnt             int
	modelQuitState       modelQuitStateType
	firstTextInput       bool
	toggleDotFile        bool
	updatedToggleDotFile bool
	toggleFooter         bool
	firstLoadingComplete bool
	firstUse             bool

	// This entirely disables metadata fetching. Used in test model
	disableMetadata     bool
	filePanelFocusIndex int

	// Height in number of lines of actual viewport of
	// main panel and sidebar excluding border
	mainPanelHeight int

	// Height in number of lines of actual viewport of
	// footer panels - process/metadata/clipboard - excluding border
	footerHeight int
	fullWidth    int
	fullHeight   int

	// whether usable trash directory exists or not
	hasTrash bool
}

// Modal
type helpMenuModal struct {
	height       int
	width        int
	open         bool
	renderIndex  int
	cursor       int
	data         []helpMenuModalData
	filteredData []helpMenuModalData
	searchBar    textinput.Model
}

type helpMenuModalData struct {
	hotkey         []string
	description    string
	hotkeyWorkType hotkeyType
	subTitle       string
}

type typingModal struct {
	location      string
	open          bool
	textInput     textinput.Model
	errorMesssage string
}

// Copied items
type copyItems struct {
	items []string
	cut   bool
}

/* FILE WINDOWS TYPE START*/
// Model for file windows
type fileModel struct {
	filePanels   []FilePanel
	width        int
	renaming     bool
	maxFilePanel int
	filePreview  preview.Model
}

/* FILE WINDOWS TYPE END*/

type editorFinishedMsg struct{ err error }

