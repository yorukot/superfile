package rendering

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (r *Renderer) assertInitialized() {
	if r.contentSections == nil {
		panic("Renderer not initialized (zero value). Use NewRenderer() to create a valid instance")
	}
}

// Add lines as much as the remaining capacity allows
func (r *Renderer) AddLines(lines ...string) *Renderer {
	r.assertInitialized()
	r.contentSections[r.curSectionIdx].AddLines(lines...)
	return r
}

// Lines until now will belong to current section, and
// Any new lines will belong to a new section
func (r *Renderer) AddSection() {
	r.assertInitialized()
	// r.actualContentHeight before this point only includes sections
	// before r.curSectionIdx
	r.actualContentHeight += r.contentSections[r.curSectionIdx].CntLines()

	// Silently Fail if cannot add
	if r.contentHeight <= r.actualContentHeight {
		slog.Error("Cannot add any more sections", "name", r.name, "actualHeight", r.actualContentHeight,
			"contentHeight", r.contentHeight)
		return
	}

	// Add divider
	r.border.AddDivider(r.actualContentHeight)
	// sectionDivider should be of borderstyle
	r.sectionDividers = append(r.sectionDividers, lipgloss.NewStyle().
		Foreground(r.borderFGColor).
		Background(r.borderBGColor).
		Render(strings.Repeat(r.borderStrings.Top, r.contentWidth)))
	r.actualContentHeight++

	remainingHeight := r.contentHeight - r.actualContentHeight
	r.contentSections = append(r.contentSections,
		NewContentRenderer(remainingHeight, r.contentWidth, r.defTruncateStyle, r.name))
	// Adjust index
	r.curSectionIdx++
}

// Truncate would always preserve ansi codes.
func (r *Renderer) AddLineWithCustomTruncate(line string, truncateStyle TruncateStyle) {
	r.assertInitialized()
	r.contentSections[r.curSectionIdx].AddLineWithCustomTruncate(line, truncateStyle)
}

func (r *Renderer) AddStyleModifier(modifier StyleModifier) *Renderer {
	r.styleModifiers = append(r.styleModifiers, modifier)
	return r
}

func (r *Renderer) SetBorderTitle(title string) {
	r.border.SetTitle(title)
}

func (r *Renderer) SetBorderInfoItems(infoItems ...string) {
	r.border.SetInfoItems(infoItems...)
}

func (r *Renderer) AreInfoItemsTruncated() bool {
	return r.border.AreInfoItemsTruncated()
}

// Should not do any updates on 'r'
func (r *Renderer) Render() string {
	r.assertInitialized()
	content := strings.Builder{}
	for i := range r.contentSections {
		// After every iteration, current cursor will be on next newline
		curContent := r.contentSections[i].Render()
		content.WriteString(curContent)
		// == "" check cant differentiate between no data, vs empty line
		if r.contentSections[i].CntLines() > 0 {
			content.WriteString("\n")
		}

		if i < len(r.contentSections)-1 {
			// True for all except last section
			content.WriteString(r.sectionDividers[i])
			content.WriteString("\n")
		}
	}
	contentStr := strings.TrimSuffix(content.String(), "\n")
	res := r.Style().Render(contentStr)
	// Post rendering validations - Maybe we can return an error instead of logging
	// TODO(perf): This can be disabled to improve performance
	maxW := 0
	for line := range strings.Lines(res) {
		maxW = max(maxW, ansi.StringWidth(line))
	}

	lineCnt := strings.Count(res, "\n") + 1
	if maxW > r.totalWidth || lineCnt > r.totalHeight {
		slog.Error(
			"Rendered output data inconsistency",
			"name",
			r.name,
			"lineCnt",
			lineCnt,
			"totalHeight",
			r.totalHeight,
			"totalWidth",
			r.totalWidth,
			"maxW",
			maxW,
		)
		// lipgloss Render() doesn't always respects the "height" value,
		// so res can have more height than intended. In that case, we must truncate lines here.
		newRes := strings.Builder{}
		curCnt := 0
		// Dont use strings.Lines(), that wont allow us to have empty lines
		for line := range strings.SplitSeq(res, "\n") {
			if curCnt == r.totalHeight {
				break
			}
			newRes.WriteString(ansi.Truncate(line, r.totalWidth, ""))
			curCnt++
			if curCnt < r.totalHeight {
				newRes.WriteByte('\n')
			}
		}
		return newRes.String()
	}

	return res
}

func (r *Renderer) Style() lipgloss.Style {
	contentHeight := r.contentHeight
	if r.truncateHeight {
		contentHeight = r.actualContentHeight
	}
	s := lipgloss.NewStyle()

	for _, modifier := range r.styleModifiers {
		s = modifier(s)
	}

	s = s.Width(r.contentWidth).
		Height(contentHeight).
		Background(r.contentBGColor).
		Foreground(r.contentFGColor)

	if r.borderRequired {
		s = s.Border(r.border.GetBorder(r.borderStrings))
		s = s.BorderForeground(r.borderFGColor).
			BorderBackground(r.borderBGColor)
	}
	return s
}
