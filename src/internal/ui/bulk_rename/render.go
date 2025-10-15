package bulkrename

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

const (
	modalWidth      = 80
	modalHeight     = 25
	leftColWidth    = 20
	rightColWidth   = 56 // modalWidth - 4 (borders) - leftColWidth - 2 (separator)
	columnHeight    = 6
	maxPreviewItems = 3
)

// Render renders the bulk rename modal
func (m *Model) Render() string {
	if !m.open {
		return ""
	}

	r := ui.BulkRenameRenderer(modalHeight, modalWidth)

	m.renderTitle(r)
	r.AddSection()
	m.renderTypeOptionsAndInputs(r)
	r.AddSection()
	m.renderPreview(r)
	r.AddSection()
	m.renderTips(r)

	if m.errorMessage != "" {
		r.AddSection()
		m.renderError(r)
	}

	return r.Render()
}

func (m *Model) renderTitle(r *rendering.Renderer) {
	count := len(m.selectedFiles)
	title := common.ModalTitleStyle.Render("  Bulk Rename") +
		common.ModalStyle.Render(fmt.Sprintf(" (%d files selected)", count))
	r.AddLines(title)
}

func (m *Model) renderTypeOptionsAndInputs(r *rendering.Renderer) {
	typeOptions := m.renderTypeOptions()

	inputs := m.renderInputs()
	leftStyle := lipgloss.NewStyle().
		Width(leftColWidth).
		Height(columnHeight).
		Background(common.ModalBGColor)

	rightStyle := lipgloss.NewStyle().
		Width(rightColWidth).
		Height(columnHeight).
		Background(common.ModalBGColor)

	separator := lipgloss.NewStyle().
		Width(2).
		Height(columnHeight).
		Background(common.ModalBGColor).
		Render("  ")

	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(typeOptions),
		separator,
		rightStyle.Render(inputs),
	)

	r.AddLines(combined)
}

func (m *Model) renderTypeOptions() string {
	types := []string{
		"Find & Replace",
		"Add Prefix",
		"Add Suffix",
		"Add Numbering",
		"Change Case",
		"Editor Mode",
	}

	cursorColor := GetCursorColor()
	typeStyle := lipgloss.NewStyle().
		Width(leftColWidth).
		Background(common.ModalBGColor).
		Foreground(common.ModalFGColor)

	var result string
	for i, typeName := range types {
		cursorIcon := icon.Cursor
		if !common.Config.Nerdfont {
			cursorIcon = ">"
		}

		line := "   " + typeName
		style := typeStyle
		if i == int(m.renameType) {
			line = " " + cursorIcon + " " + typeName
			style = typeStyle.Foreground(cursorColor)
		}
		if i > 0 {
			result += "\n"
		}
		result += style.Render(line)
	}
	return result
}

func (m *Model) renderInputs() string {
	inputStyle := lipgloss.NewStyle().
		Width(rightColWidth).
		Background(common.ModalBGColor)

	labelStyle := lipgloss.NewStyle().
		Background(common.ModalBGColor).
		Foreground(common.ModalFGColor)

	activeLabelStyle := lipgloss.NewStyle().
		Background(common.ModalBGColor).
		Foreground(GetCursorColor())

	var content string

	switch m.renameType {
	case FindReplace:
		findStyle := labelStyle
		replaceStyle := labelStyle
		if m.cursor == 0 {
			findStyle = activeLabelStyle
		}
		if m.cursor == 1 {
			replaceStyle = activeLabelStyle
		}

		findLine := findStyle.Render("Find:    ") + m.findInput.View()
		replaceLine := replaceStyle.Render("Replace: ") + m.replaceInput.View()
		content = inputStyle.Render(findLine) + "\n" + inputStyle.Render(replaceLine)

	case AddPrefix:
		content = inputStyle.Render(activeLabelStyle.Render("Prefix: ") + m.prefixInput.View())

	case AddSuffix:
		content = inputStyle.Render(activeLabelStyle.Render("Suffix: ") + m.suffixInput.View())

	case AddNumbering:
		numberText := fmt.Sprintf("Start number: %d\n(Use ↑/↓ to adjust)", m.startNumber)
		content = inputStyle.Render(labelStyle.Render(numberText))

	case ChangeCase:
		content = m.renderCaseOptions(inputStyle, labelStyle)

	case EditorMode:
		content = inputStyle.Render(labelStyle.Render("Opens your $EDITOR\nwith list of filenames"))
	}

	return content
}

func (m *Model) renderCaseOptions(inputStyle, labelStyle lipgloss.Style) string {
	caseTypes := []string{"lowercase", "UPPERCASE", "Title Case"}
	cursorColor := GetCursorColor()
	var result string

	for i, caseType := range caseTypes {
		style := labelStyle
		cursorIcon := icon.Cursor
		if !common.Config.Nerdfont {
			cursorIcon = ">"
		}

		line := "   " + caseType
		if i == int(m.caseType) {
			line = " " + cursorIcon + " " + caseType
			style = labelStyle.Foreground(cursorColor)
		}
		result += inputStyle.Render(style.Render(line)) + "\n"
	}
	return result
}

func (m *Model) renderPreview(r *rendering.Renderer) {
	if len(m.preview) == 0 {
		m.preview = m.GeneratePreview()
	}

	previewCount := min(maxPreviewItems, len(m.preview))
	if previewCount == 0 {
		return
	}

	previewTitleStyle := lipgloss.NewStyle().
		Background(common.ModalBGColor).
		Foreground(common.ModalFGColor)

	r.AddLines(previewTitleStyle.Render("  Preview:"))

	for i := 0; i < previewCount; i++ {
		preview := m.preview[i]
		availableWidth := modalWidth - 6
		truncatedName := common.TruncateText(preview.NewName, availableWidth, "...")

		lineStyle := lipgloss.NewStyle().
			Background(common.ModalBGColor).
			Foreground(common.ModalFGColor)

		if preview.Error != "" {
			errorStyle := lipgloss.NewStyle().
				Background(common.ModalBGColor).
				Foreground(lipgloss.Color(common.Theme.Error))

			r.AddLines(errorStyle.Render("  " + truncatedName))
			r.AddLines(errorStyle.Render("  " + fmt.Sprintf("(%s)", preview.Error)))
		} else {
			r.AddLines(lineStyle.Render("  " + truncatedName))
		}
	}

	if len(m.preview) > previewCount {
		moreText := fmt.Sprintf(" ... and %d more files", len(m.preview)-previewCount)
		moreStyle := lipgloss.NewStyle().
			Background(common.ModalBGColor).
			Foreground(common.ModalFGColor)
		r.AddLines(moreStyle.Render("  " + moreText))
	}
}

func (m *Model) renderTips(r *rendering.Renderer) {
	tips := " Tab/Shift+Tab: Change type | ↑/↓: Navigate | Enter: Rename | Esc: Cancel"
	tipsStyle := lipgloss.NewStyle().
		Background(common.ModalBGColor).
		Foreground(common.ModalFGColor)
	r.AddLines(tipsStyle.Render(tips))
}

func (m *Model) renderError(r *rendering.Renderer) {
	errorStyle := lipgloss.NewStyle().
		Background(common.ModalBGColor).
		Foreground(lipgloss.Color(common.Theme.Error))
	r.AddLines(errorStyle.Render("  " + m.errorMessage))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
