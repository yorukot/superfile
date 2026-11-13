package internal

import (
	"sync"

	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/internal/ui/helpmenu"
	"github.com/yorukot/superfile/src/internal/ui/spferror"

	"github.com/yorukot/superfile/src/internal/ui/clipboard"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"

	"github.com/yorukot/superfile/src/internal/mouse"

	"github.com/yorukot/superfile/src/internal/ui/contextmenu"
	"github.com/yorukot/superfile/src/internal/ui/filemodel"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	"charm.land/bubbles/v2/textinput"

	"github.com/yorukot/superfile/src/internal/ui/prompt"
	zoxideui "github.com/yorukot/superfile/src/internal/ui/zoxide"
)

// Type representing the type of focused panel
type focusPanelType int

type modelQuitStateType int

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
	fileModel       filemodel.Model
	sidebarModel    sidebar.Model
	processBarModel processbar.Model
	clipboard       clipboard.Model
	clipboardWriter func(string) error
	focusPanel      focusPanelType

	// Modals
	notifyModel notify.Model
	typingModal typingModal
	helpMenu    helpmenu.Model
	promptModal prompt.Model
	zoxideModal zoxideui.Model
	sortModal   sortmodel.Model
	spfError    spferror.Model
	// contextMenu is the right click menu. It is an overlay anchored to the
	// pointer rather than a centred modal.
	contextMenu     contextmenu.Model
	mutexErrorModal sync.Mutex

	// Zoxide client for directory tracking
	zClient *zoxidelib.Client

	fileMetaData metadata.Model

	// no use directly for increment, use nextIoReqCnt
	ioReqCnt int32

	modelQuitState       modelQuitStateType
	firstTextInput       bool
	toggleFooter         bool
	firstLoadingComplete bool
	firstUse             bool

	// This entirely disables metadata fetching. Used in test model
	disableMetadata bool

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

	// lastLeftClick is the previous left click, used to recognise double clicks.
	lastLeftClick leftClick

	// mouseRegions maps terminal coordinates onto the widget drawn there. It is
	// rebuilt on every render pass, so it always describes the frame currently
	// on screen rather than the layout as it will be after the next update.
	mouseRegions mouse.Registry
}

type typingModal struct {
	location string
	open     bool
	// directory makes the next typed name create a directory instead of a file,
	// so that a caller can ask for a directory without the user having to know
	// that a trailing path separator is what selects one.
	directory bool
	textInput textinput.Model
}

type editorFinishedMsg struct{ err error }
