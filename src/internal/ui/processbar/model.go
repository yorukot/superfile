package processbar

import (
	"fmt"
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
)

// Model for process bar internal
type Model struct {
	renderIndex int
	cursor      int

	// Including borders
	height int
	width  int

	// TODO: Fix this. No mechanism to remove completed processes from memory
	// processes map grows indefinitely
	// Maybe, TTL or cleanup mechanism for successful/failed processes
	processes map[string]Process
	msgChan   chan UpdateMsg
	reqCnt    int
}

func New() Model {
	return NewModelWithOptions(minWidth, minHeight)
}

// Note: We should considering our internal models, they
// should be returning pointer object, and implement tea.Model
func NewModelWithOptions(width int, height int) Model {
	m := Model{
		renderIndex: 0,
		cursor:      0,
		processes:   make(map[string]Process),
		msgChan:     make(chan UpdateMsg, msgChannelSize),
		reqCnt:      0,
	}
	m.SetDimensions(width, height)
	return m
}

func (m *Model) SetDimensions(width int, height int) {
	if width < minWidth {
		slog.Warn("Invalid width, using minimum", "provided", width, "minimum", minWidth)
		width = minWidth
	}
	if height < minHeight {
		slog.Warn("Invalid height, using minimum", "provided", height, "minimum", minHeight)
		height = minHeight
	}
	m.width = width
	m.height = height
}

func (m *Model) AddProcess(p Process) error {
	if _, ok := m.processes[p.ID]; ok {
		return &ProcessAlreadyExistsError{id: p.ID}
	}
	m.processes[p.ID] = p
	return nil
}

func (m *Model) AddOrUpdateProcess(p Process) {
	m.processes[p.ID] = p
}

func (m *Model) UpdateExistingProcess(p Process) error {
	if _, ok := m.processes[p.ID]; !ok {
		return &NoProcessFoundError{id: p.ID}
	}
	m.processes[p.ID] = p
	return nil
}

func (m *Model) GetByID(id string) (Process, bool) {
	p, ok := m.processes[id]
	return p, ok
}

func (m *Model) HasRunningProcesses() bool {
	for _, data := range m.processes {
		if data.State == InOperation && data.Done != data.Total {
			return true
		}
	}
	return false
}

func (m *Model) Render(processBarFocussed bool) string {
	r := ui.ProcessBarRenderer(m.height, m.width, processBarFocussed)
	if !m.isValid() {
		slog.Error("processBar in invalid state", "render", m.renderIndex,
			"cursor", m.cursor, "Height", m.height)
		r.AddLines("Invalid state")
		return r.Render()
	}
	if m.cntProcesses() == 0 {
		r.AddLines("", " "+common.ProcessBarNoneText)
		return r.Render()
	}

	r.SetBorderInfoItems(fmt.Sprintf("%d/%d", m.cursor+1, m.cntProcesses()))

	renderedHeight := 0
	processes := m.getSortedProcesses()
	for i := m.renderIndex; i < len(processes); i++ {
		// We allow rendering of a process if we have at least 2 lines left
		if m.viewHeight() < renderedHeight+2 {
			break
		}
		renderedHeight += 3

		// Note : We will be updating this on each Render, although harmless from performance
		// perspective. We are rendering modified version of the data.
		// TODO: We could, save pointer of process in map and update progressbar of each
		// map on each SetWidth. This would be cleaner and more efficient.
		curProcess := processes[i]
		curProcess.Progress.Width = m.viewWidth() - ProgressBarRightPadding

		// TODO : get them via a separate function.
		var cursor string
		if i == m.cursor {
			// TODO : Prerender it.
			cursor = common.FooterCursorStyle.Render("┃ ")
		} else {
			cursor = common.FooterCursorStyle.Render("  ")
		}

		r.AddLines(cursor + common.FooterStyle.Render(
			common.TruncateText(curProcess.Name, m.viewWidth()-ProcessNameTruncatePadding, "...")+" ") +
			curProcess.State.Icon())

		// calculate progress percentage
		// if the total is 0, that means the process only have directory
		// so we can set the progress to 100%
		if curProcess.Total != 0 {
			progressPercentage := float64(curProcess.Done) / float64(curProcess.Total)
			r.AddLines(cursor+curProcess.Progress.ViewAs(progressPercentage), "")
		} else {
			r.AddLines(cursor + curProcess.Progress.ViewAs(1))
		}
	}

	return r.Render()
}
