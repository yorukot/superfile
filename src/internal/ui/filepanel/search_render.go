package filepanel

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

// Renders result list. Last row is the status line.
func (m *Model) renderSearchResults(r *rendering.Renderer) {
	if m.SearchBar.Value() == "" {
		root := common.TruncateTextBeginning(m.Search.Root, m.GetContentWidth(), "...")
		r.AddLines(common.FilePanelNoneText)
		r.AddLines(common.FilePanelStyle.Render(" Type a query to search in " + root))
		return
	}

	if len(m.Search.Results) == 0 {
		r.AddLines(common.FilePanelNoneText)
		if m.Search.Done {
			r.AddLines(common.FilePanelStyle.Render(" No matches for \"" + m.SearchBar.Value() + "\""))
		} else {
			r.AddLines(common.FilePanelStyle.Render(" Searching..."))
		}
	}

	visibleRows := max(0, m.PanelElementHeight()-1)
	end := min(m.Search.RenderIndex+visibleRows, len(m.Search.Results))
	for itemIndex := m.Search.RenderIndex; itemIndex < end; itemIndex++ {
		r.AddLines(m.renderSearchRow(itemIndex))
	}
	r.AddLines(m.renderSearchStatusLine())
}

func (m *Model) renderSearchRow(itemIndex int) string {
	result := m.Search.Results[itemIndex]

	cursor := emptyCursor
	if itemIndex == m.Search.Cursor {
		cursor = icon.Cursor
	}
	prefix := common.FilePanelCursorStyle.Render(cursor + " ")
	contentWidth := m.GetContentWidth() - ansi.StringWidth(prefix)
	if contentWidth <= 0 {
		return ""
	}

	path, positions := truncateSearchPath(result.Path, result.Positions, contentWidth)
	itemIcon := common.GetElementIcon(path, result.Dir, false, common.Config.Nerdfont)
	renderedIcon := common.StringColorRender(lipgloss.Color(itemIcon.Color), common.FilePanelBGColor).
		Background(common.FilePanelBGColor).Render(itemIcon.Icon + " ")

	pathWidth := contentWidth - ansi.StringWidth(renderedIcon)
	return prefix + renderedIcon + renderHighlightedPath(path, positions, pathWidth)
}

// Renders path with highlight, padded to maxWidth.
func renderHighlightedPath(path string, positions []int, maxWidth int) string {
	var builder strings.Builder
	for _, segment := range splitHighlight(path, positions) {
		style := common.FilePanelStyle
		if segment.matched {
			style = common.SearchModeHighlightStyle
		}
		builder.WriteString(style.Render(segment.text))
	}
	rendered := builder.String()
	remaining := maxWidth - ansi.StringWidth(rendered)
	if remaining > 0 {
		rendered += common.FilePanelStyle.Render(strings.Repeat(" ", remaining))
	}
	return rendered
}

// Truncates from start to fit. Adjusts positions.
func truncateSearchPath(path string, positions []int, maxWidth int) (string, []int) {
	if ansi.StringWidth(path) <= maxWidth {
		return path, positions
	}
	cut := 0
	for cut < len(path) {
		_, size := utf8.DecodeRuneInString(path[cut:])
		if size == 0 {
			break
		}
		cut += size
		if ansi.StringWidth(path[cut:]) <= maxWidth {
			break
		}
	}
	if cut >= len(path) {
		return path, positions
	}
	newPositions := make([]int, 0, len(positions))
	for _, pos := range positions {
		if pos >= cut {
			newPositions = append(newPositions, pos-cut)
		}
	}
	return path[cut:], newPositions
}

type highlightSegment struct {
	text    string
	matched bool
}

// Splits path by byte-offset matches.
func splitHighlight(path string, positions []int) []highlightSegment {
	if len(positions) == 0 {
		return []highlightSegment{{text: path}}
	}
	posSet := make(map[int]struct{}, len(positions))
	for _, pos := range positions {
		posSet[pos] = struct{}{}
	}

	segments := make([]highlightSegment, 0, len(positions)+1)
	var current strings.Builder
	matched := false
	byteOffset := 0
	for _, r := range path {
		_, inMatch := posSet[byteOffset]
		byteOffset += utf8.RuneLen(r)
		if inMatch != matched {
			if current.Len() > 0 {
				segments = append(segments, highlightSegment{text: current.String(), matched: matched})
				current.Reset()
			}
			matched = inMatch
		}
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		segments = append(segments, highlightSegment{text: current.String(), matched: matched})
	}
	return segments
}

// Status shows match count, scanning state, cap and hidden binding.
// Found can exceed shown.
func (m *Model) renderSearchStatusLine() string {
	var status strings.Builder
	shown := len(m.Search.Results)
	found := max(m.Search.MatchCount, int64(shown))
	fmt.Fprintf(&status, " %d found", found)
	if !m.Search.Done {
		status.WriteString(" | scanning...")
	} else if m.Search.UnreadableDirs > 0 {
		fmt.Fprintf(&status, " | %d unreadable dir(s) skipped", m.Search.UnreadableDirs)
	}
	if m.Search.MatchCount > int64(shown) {
		fmt.Fprintf(&status, " | showing first %d", shown)
	}
	if label := searchToggleHiddenLabel(); label != "" {
		fmt.Fprintf(&status, " | toggle dotfiles: %s", label)
	}
	return common.FilePanelStyle.Render(status.String())
}

// First configured hidden-toggle binding, or "" when none.
func searchToggleHiddenLabel() string {
	for _, key := range common.Hotkeys.SearchToggleHidden {
		if key != "" {
			return key
		}
	}
	return ""
}
