package filepanel

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/internal/common"
)

type fileElementRender func(index int, item Element) string
type renderGenerator func(columnWidth int) fileElementRender

const FileSizeColumnWidth = 15
const ColumnHeaderHeight = 1

// If the percentage column is smaller than this number, the additional columns will be hidden.
const FileNameRatio = 0.65

// Delimiter between columns in the file panel.
const ColumnDelimiter = "  "

type columnDefinition struct {
	Name      string
	Size      int
	Generator renderGenerator
}

func (cd *columnDefinition) GetRenderer() fileElementRender {
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
