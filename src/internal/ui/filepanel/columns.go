package filepanel

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/internal/common"
)

type columnRenderer func(indexElement int) string
type renderGenerator func(columnWidth int) columnRenderer

const FileSizeColumnWidth = 15
const ModifyTimeSizeColumnWidth = 18
const PermissionsColumnWidth = 12
const ColumnHeaderHeight = 1

// If the percentage column is smaller than this number, the additional columns will be hidden.
// TODO: make this configurable
const FileNameRatio = 0.65

// Delimiter between columns in the file panel.
const ColumnDelimiter = "  "

type columnDefinition struct {
	Name      string
	Size      int
	Generator renderGenerator
}

func (cd *columnDefinition) GetRenderer() columnRenderer {
	return cd.Generator(cd.Size)
}

func (cd *columnDefinition) RenderHeader() string {
	return common.PrettierFixedWidthItem(
		cd.Name,
		cd.Size,
		false,
		common.FilePanelBGColor,
		lipgloss.Center,
	)
}
