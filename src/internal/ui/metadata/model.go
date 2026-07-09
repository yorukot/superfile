package metadata

import (
	"fmt"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/pkg/cache"
)

type Model struct {
	metadata Metadata // current metadata
	cache    *cache.Cache[Metadata]

	// It tells what the metadata should have. Its used to prevent additional requests
	// if one is already underway
	expectedLocation string
	expectedFocused  bool

	// Render state
	renderIndex int

	// Model Dimensions, including borders
	width  int
	height int
}

func New() Model {
	return Model{
		cache: cache.New[Metadata](defaultCacheSize, defaultCacheExpiration),
	}
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

// Move renderIndex by delta rows, wrapping at both ends.
func (m *Model) moveRenderIndexBy(delta int) {
	l := m.MetadataLen()
	if l == 0 {
		return
	}
	m.renderIndex = ((m.renderIndex+delta)%l + l) % l
}

func (m *Model) getPageScrollSize() int {
	scrollSize := common.Config.PageScrollSize
	if scrollSize <= 0 {
		// Use default full page behavior
		scrollSize = m.height - borderSize
	}
	// height can be tiny on small terminals, so keep moving at least one row
	if scrollSize < 1 {
		scrollSize = 1
	}
	return scrollSize
}

// Control metadata panel up
func (m *Model) ListUp() {
	m.moveRenderIndexBy(-1)
}

// Control metadata panel down
func (m *Model) ListDown() {
	m.moveRenderIndexBy(1)
}

// Control metadata panel page up
func (m *Model) PgUp() {
	m.moveRenderIndexBy(-m.getPageScrollSize())
}

// Control metadata panel page down
func (m *Model) PgDown() {
	m.moveRenderIndexBy(m.getPageScrollSize())
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

func (m *Model) Render(metadataFocused bool) string {
	r := ui.MetadataRenderer(m.height, m.width, metadataFocused)
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
