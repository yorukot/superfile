package metadata

import (
	"fmt"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/utils"
)

type Model struct {
	// Data
	metadata [][2]string
	filePath string

	// Render state
	renderIndex int

	// Model Dimensions, including borders
	width  int
	height int

	// Render dimensions
	maxKeyLen  int
	sprintfLen int
	valLen     int
}

func New() Model {
	return Model{}
}

// Should be at least 2x2
// TODO : Validate this
func (m *Model) SetDimensions(width int, height int) {
	m.width = width
	m.height = height
}

func (m *Model) SetMedatada(filepath string, metadata [][2]string) {
	m.filePath = filepath
	m.metadata = metadata
	m.renderIndex = 0
}

func (m *Model) ResetRender() {
	m.renderIndex = 0
}

// Control metadata panel up
func (m *Model) ListUp() {
	if len(m.metadata) == 0 {
		return
	}
	if m.renderIndex > 0 {
		m.renderIndex--
	} else {
		m.renderIndex = len(m.metadata) - 1
	}
}

func (m *Model) SetBlank() {
	m.filePath = ""
	m.metadata = m.metadata[:0]
}

func (m *Model) IsBlank() bool {
	return len(m.metadata) == 0
}

func (m *Model) SetLoading() {
	// Note : This will cause gc of current metadata slice
	// This will cause frequent allocations and gc.
	m.metadata = [][2]string{
		{"", ""},
		{" " + icon.InOperation + icon.Space + "Loading metadata...", ""},
	}
}

// Control metadata panel down
func (m *Model) ListDown() {
	if m.renderIndex < len(m.metadata)-1 {
		m.renderIndex++
	} else {
		m.renderIndex = 0
	}
}

func (m *Model) computeRenderDimensions() {
	// Recompute dimension based values
	m.maxKeyLen = getMaxKeyLength(m.metadata)
	m.sprintfLen, m.valLen = computeMetadataWidths(m.width-2, m.maxKeyLen)
}

func (m *Model) Render(metadataFocussed bool) string {
	if len(m.metadata) == 0 {
		return ""
	}
	m.computeRenderDimensions()
	r := ui.MetadataRenderer(m.height, m.width, metadataFocussed)
	r.SetBorderInfoItems(fmt.Sprintf("%d/%d", m.renderIndex+1, len(m.metadata)))
	lines := formatMetadataLines(m.metadata, m.renderIndex, m.height-2, m.sprintfLen, m.valLen)
	r.AddLines(lines...)
	return r.Render()
}

// TODO : Simplify these mystic calculations, or add explanation comments.
// TODO : unit test and fix this mess
func formatMetadataLines(meta [][2]string, startIdx, height, sprintfLen, valueLen int) []string {
	lines := []string{}
	endIdx := min(startIdx+height, len(meta))
	for i := startIdx; i < endIdx; i++ {
		key := meta[i][0]
		value := common.TruncateMiddleText(meta[i][1], valueLen, "...")
		if utils.FooterWidth(0)-sprintfLen-3 < utils.FooterWidth(0)/2 {
			key = common.TruncateMiddleText(key, valueLen, "...")
		}
		line := fmt.Sprintf("%-*s %s", sprintfLen, key, value)
		lines = append(lines, line)
	}
	return lines
}
