package trash

import "errors"

type Backend string

const (
	BackendFreeDesktop Backend = "freedesktop"
	BackendMacOS       Backend = "macos"
	BackendWindows     Backend = "windows"
)

var ErrUnsupported = errors.New("trash is not supported on this platform")

type Result struct {
	OriginalPath     string
	TrashedPath      string
	Backend          Backend
	StrictlyRecycled bool
}
