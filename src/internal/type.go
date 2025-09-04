package internal

import (
	"time"

	"github.com/lazysegtree/go-zoxide"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/yorukot/superfile/src/internal/ui/preview"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
)

// Type representing the mode of the panel
type panelMode uint

// Type representing the focus type of the file panel
type filePanelFocusType uint

// Type representing the type of focused panel
type focusPanelType int

type hotkeyType int

type modelQuitStateType int

// TODO: Convert to integer enum
type sortingKind string

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

// Constants for file panel with no focus
const (
	noneFocus filePanelFocusType = iota
	secondFocus
	focus
)

// Constants for select mode or browser mode
const (
	selectMode panelMode = iota
	browserMode
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
	disableMetatdata    bool
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

	// Executors for plugins
	zClient *zoxide.Client
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
	filePanels   []filePanel
	width        int
	renaming     bool
	maxFilePanel int
	filePreview  preview.Model
}

// Panel representing a file
type filePanel struct {
	cursor             int
	render             int
	focusType          filePanelFocusType
	location           string
	sortOptions        sortOptionsModel
	panelMode          panelMode
	selected           []string
	element            []element
	directoryRecords   map[string]directoryRecord
	rename             textinput.Model
	renaming           bool
	searchBar          textinput.Model
	lastTimeGetElement time.Time
}

// Sort options
type sortOptionsModel struct {
	width  int
	height int
	open   bool
	cursor int
	data   sortOptionsModelData
}

type sortOptionsModelData struct {
	options  []string
	selected int
	reversed bool
}

// Record for directory navigation
type directoryRecord struct {
	directoryCursor int
	directoryRender int
}

// Element within a file panel
type element struct {
	name      string
	location  string
	directory bool
	metaData  [][2]string
}

/* FILE WINDOWS TYPE END*/

type editorFinishedMsg struct{ err error }

type sliceOrderFunc func(i, j int) bool
