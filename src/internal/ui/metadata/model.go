package metadata

import (
	"fmt"

	"github.com/yorukot/superfile/src/internal/ui"
)

type Model struct {
	// Data
	metadata Metadata

	// Render state
	renderIndex int

	// Model Dimensions, including borders
	width  int
	height int
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

func (m *Model) GetHeight() int {
	return m.height
}

func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) SetMetadata(metadata Metadata) {
	m.metadata = metadata
	// Note : Dont always reset render to 0
	// We would have update requests coming in during user scrolling through metadata
	m.ResetRenderIfInvalid()
}

func (m *Model) ResetRenderIfInvalid() {
	if m.renderIndex >= m.MetadataLen() {
		m.ResetRender()
	}
}

func (m *Model) ResetRender() {
	m.renderIndex = 0
}

func (m *Model) MetadataLen() int {
	return len(m.metadata.data)
}

// Control metadata panel up
func (m *Model) ListUp() {
	if m.MetadataLen() == 0 {
		return
	}
	if m.renderIndex > 0 {
		m.renderIndex--
	} else {
		m.renderIndex = m.MetadataLen() - 1
	}
}

// Control metadata panel down
func (m *Model) ListDown() {
	if m.renderIndex < m.MetadataLen()-1 {
		m.renderIndex++
	} else {
		m.renderIndex = 0
	}
}

func (m *Model) SetBlank() {
	m.metadata.filepath = ""
	m.metadata.data = m.metadata.data[:0]
	m.metadata.infoMsg = "No metadata present"
}

func (m *Model) IsBlank() bool {
	return m.MetadataLen() == 0 && m.metadata.infoMsg == ""
}

func (m *Model) SetInfoMsg(msg string) {
	m.metadata.infoMsg = msg
}

func (m *Model) Render(metadataFocussed bool) string {
	r := ui.MetadataRenderer(m.height, m.width, metadataFocussed)
	if m.MetadataLen() == 0 {
		r.AddLines("", " "+m.metadata.infoMsg)
		return r.Render()
	}
	keyLen, valueLen := computeRenderDimensions(m.metadata.data, m.width-2-keyValueSpacingLen)
	r.SetBorderInfoItems(fmt.Sprintf("%d/%d", m.renderIndex+1, len(m.metadata.data)))
	lines := formatMetadataLines(m.metadata.data, m.renderIndex, m.height-borderSize, keyLen, valueLen)
	r.AddLines(lines...)
	return r.Render()
}
