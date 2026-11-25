package filemodel

import (
	"log/slog"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
)

// sixelDrawDelay is how long to wait before drawing a Sixel image. Bubble
// Tea writes raw output before each frame, so the image must be sent after
// the frame that puts the marker cells on screen, or that frame would paint
// over it. A few frames' worth at 60fps.
const sixelDrawDelay = 50 * time.Millisecond

// SixelDrawMsg asks the model to draw the preview's Sixel image. It is
// dropped if anything changed since it was scheduled.
type SixelDrawMsg struct{ gen int }

// sixelState tracks the Sixel image drawn in the preview panel.
type sixelState struct {
	// gen identifies the latest scheduled draw. Older SixelDrawMsgs are stale.
	gen int
	// drawn is true when a draw is scheduled or done for the current content
	// and the image has not been erased since.
	drawn bool
	// hidden is true while something (an overlay) is drawn over the preview.
	hidden bool
}

// InvalidateSixel marks the Sixel image as erased, e.g. after the whole
// screen was repainted (resize, returning from an external program), so the
// next SyncSixel draws it again.
func (m *Model) InvalidateSixel() {
	m.sixel.drawn = false
}

// SyncSixel hides or shows the preview's Sixel image and schedules a draw
// when it needs one. Call it at the end of every Update, after the preview
// content is up to date. hidden must be true while anything is drawn over
// the preview panel, such as an overlay.
func (m *Model) SyncSixel(hidden bool) tea.Cmd {
	m.sixel.hidden = hidden
	m.FilePreview.SetSixelHidden(hidden)

	if !m.sixelVisible() {
		// The marker cells are gone from the view (or will be at the next
		// frame), and rewriting them erases the image.
		m.sixel.drawn = false
		return nil
	}
	if m.sixel.drawn {
		return nil
	}

	m.sixel.drawn = true
	m.sixel.gen++
	gen := m.sixel.gen
	return tea.Tick(sixelDrawDelay, func(time.Time) tea.Msg {
		return SixelDrawMsg{gen: gen}
	})
}

// sixelPlacement is the raw output drawing a Sixel image. It carries the
// generation it was drawn for, so HandleSixelRaw can tell if it went stale.
type sixelPlacement struct {
	gen int
	seq string
}

// String returns the escape sequences. tea.Raw writes the String() value.
func (p sixelPlacement) String() string {
	return p.seq
}

// HandleSixelDraw returns the command drawing the Sixel image, or nil if msg
// is stale.
func (m *Model) HandleSixelDraw(msg SixelDrawMsg) tea.Cmd {
	if !m.sixelDrawCurrent(msg.gen) {
		slog.Debug("Dropping stale Sixel draw", "gen", msg.gen, "currentGen", m.sixel.gen)
		return nil
	}
	x, y := m.previewContentOrigin()
	m.FilePreview.MarkSixelDrawn()
	return tea.Raw(sixelPlacement{gen: msg.gen, seq: m.FilePreview.GetSixel().Placement(x, y)})
}

// HandleSixelRaw checks a Sixel image that was just written to the terminal.
// Bubble Tea writes raw output before the model sees it, and the command
// returned by HandleSixelDraw runs asynchronously, so the preview may have
// changed in between, with the image painted over newer content. Then it
// returns a command that repaints the whole screen, erasing it.
func (m *Model) HandleSixelRaw(msg tea.RawMsg) tea.Cmd {
	placement, ok := msg.Msg.(sixelPlacement)
	if !ok || m.sixelDrawCurrent(placement.gen) {
		return nil
	}
	slog.Debug("Stale Sixel image was drawn, repainting", "gen", placement.gen, "currentGen", m.sixel.gen)
	m.InvalidateSixel()
	return tea.ClearScreen
}

// sixelDrawCurrent reports whether a draw scheduled as gen is still the
// latest, and the image is still visible.
func (m *Model) sixelDrawCurrent(gen int) bool {
	return gen == m.sixel.gen && m.sixel.drawn && m.sixelVisible()
}

// sixelVisible reports whether the preview panel currently shows its Sixel
// image's marker cells.
func (m *Model) sixelVisible() bool {
	return !m.sixel.hidden &&
		m.FilePreview.IsOpen() &&
		m.FilePreview.GetSixel() != nil &&
		!m.FilePreview.IsLoading() &&
		m.FilePreview.GetContentWidth() == m.ExpectedPreviewWidth &&
		m.FilePreview.GetContentHeight() == m.Height
}

// previewContentOrigin returns the 0-based screen position of the preview
// panel's top-left content cell.
func (m *Model) previewContentOrigin() (int, int) {
	x := m.fullModalWidth() - m.ExpectedPreviewWidth
	y := 0
	if common.Config.EnableFilePreviewBorder {
		x++
		y++
	}
	return x, y
}

// fullModalWidth returns the width of the sidebar and file model together.
// HACK!!!. fileModel must not be aware of other dimensions. but...
// Unfortunately, previewPanel isn't completely 'under' fileModel
func (m *Model) fullModalWidth() int {
	width := m.Width + common.Config.SidebarWidth
	if common.Config.SidebarWidth != 0 {
		width += common.BorderPadding
	}
	return width
}
