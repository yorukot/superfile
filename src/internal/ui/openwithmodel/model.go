package openwithmodel

import (
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/internal/common"
)

const (
	defaultWidth = 32
	maxWidth     = 60
)

type Model struct {
	Width  int
	Height int
	cursor int
	open   bool
	path   string
	apps   []common.OpenWithApplication
}

func New() Model {
	return Model{Width: defaultWidth}
}

func (m *Model) Open(path string, applications []common.OpenWithApplication) {
	m.path = path
	m.apps = append([]common.OpenWithApplication(nil), applications...)
	m.cursor = 0
	m.open = true
	m.resize()
}

func (m *Model) Close() {
	m.open = false
	m.path = ""
	m.apps = nil
	m.cursor = 0
	m.resize()
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Path() string {
	return m.path
}

func (m *Model) ListUp() {
	if len(m.apps) == 0 {
		return
	}
	m.cursor = (m.cursor - 1 + len(m.apps)) % len(m.apps)
}

func (m *Model) ListDown() {
	if len(m.apps) == 0 {
		return
	}
	m.cursor = (m.cursor + 1) % len(m.apps)
}

func (m *Model) Selected() (common.OpenWithApplication, bool) {
	if m.cursor < 0 || m.cursor >= len(m.apps) {
		return common.OpenWithApplication{}, false
	}
	return m.apps[m.cursor], true
}

func (m *Model) Cursor() int {
	return m.cursor
}

func (m *Model) resize() {
	m.Width = defaultWidth
	for _, application := range m.apps {
		m.Width = max(m.Width, ansi.StringWidth(application.Name)+4)
	}
	m.Width = min(m.Width, maxWidth)
	m.Height = max(3, len(m.apps)+2)
}
