package preview

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"

	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

type Model struct {
	open bool

	// Location denotes what is supposed to be in model.
	// Might not be always in sync with content
	location      string
	content       string
	contentWidth  int
	contentHeight int

	loading bool

	// sixel is the Sixel image shown by content, if any. content then has
	// marker cells where the image goes. See filepreview.SixelMarker.
	sixel *filepreview.SixelImage
	// sixelAlt is true when content uses SixelMarkerAlt instead of SixelMarker.
	sixelAlt bool
	// drawnSixelAlt is the sixelAlt of the image most recently drawn on the
	// terminal. See setSixel.
	drawnSixelAlt bool
	// sixelHidden hides the Sixel image, e.g. while an overlay is open, by
	// showing hiddenContent (content with plain spaces instead of markers).
	sixelHidden   bool
	hiddenContent string

	imagePreviewer     *filepreview.ImagePreviewer
	batCmd             string
	thumbnailGenerator *filepreview.ThumbnailGenerator
}

func New() Model {
	generator, err := filepreview.NewThumbnailGenerator()
	if err != nil {
		slog.Error("Could not NewThumbnailGenerator object", "error", err)
	}

	return Model{
		open: common.Config.DefaultOpenFilePreview,
		// TODO: This causes unnecessary terminal cell size detection
		// logs in tests, we should not be initializing it in tests
		// when  `DefaultOpenFilePreview` is false
		// And only initialize these objects when Open() is called
		// Still them being nil should be handled well, right we don't
		// have code with good defensive programming
		// Some of these processes are IO operations so maybe it should
		// be done via an Init() function
		imagePreviewer:     filepreview.NewImagePreviewer(),
		thumbnailGenerator: generator,
		// TODO:  This is an IO operation, move to async ?
		batCmd: checkBatCmd(),
	}
}
