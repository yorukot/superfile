package preview

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"

	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

type Model struct {
	open   bool
	width  int
	height int

	// Location denotes what is supposed to be in model.
	// Might not be always in sync with content
	location           string
	content            string
	loading            bool
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
		open:               common.Config.DefaultOpenFilePreview,
		imagePreviewer:     filepreview.NewImagePreviewer(),
		thumbnailGenerator: generator,
		// TODO:  This is an IO operation, move to async ?
		batCmd: checkBatCmd(),
	}
}
