package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

// Non fatal Validations. This indicates bug / programming errors, not user configuration mistake
func (m *model) validateLayout() error { //nolint:gocognit // cumilation of validations
	// Validate footer height
	if 0 < m.footerHeight && m.footerHeight < common.MinFooterHeight {
		return fmt.Errorf("footerHeight %v is too small", m.footerHeight)
	}
	if !m.toggleFooter && m.footerHeight != 0 {
		return fmt.Errorf("footer closed and footerHeight %v is non zero", m.footerHeight)
	}
	if m.toggleFooter && m.footerHeight == 0 {
		return errors.New("footer open but footerHeight is 0")
	}

	// PanelHeight + 2 lines (main border) + actual footer height
	if m.fullHeight != (m.mainPanelHeight+common.BorderPadding)+utils.FullFooterHeight(m.footerHeight, m.toggleFooter) {
		return fmt.Errorf(
			"invalid model layout, total height doesn't sum correctly, fullHeight : %v, mainPanelHeight : %v, footerHeight : %v",
			m.fullHeight,
			m.mainPanelHeight,
			m.footerHeight,
		)
	}

	// Validate width constraints
	if m.fullWidth < common.MinimumWidth {
		return fmt.Errorf("fullWidth %v is below minimum %v", m.fullWidth, common.MinimumWidth)
	}

	// Check that file panel width is positive if we have file panels
	if m.fileModel.PanelCount() == 0 {
		return errors.New("file model is empty")
	}

	// Check total width calculation consistency
	if m.fullWidth != common.Config.SidebarWidth+common.BorderPadding+m.fileModel.Width {
		return fmt.Errorf(
			"width layout inconsistent: fullWidth=%v, sidebar=%v filemodel=%v",
			m.fullWidth, common.Config.SidebarWidth, m.fileModel.Width)
	}

	// Check file panels count
	if m.fileModel.PanelCount() > m.fileModel.MaxFilePanel {
		return fmt.Errorf(
			"too many file panels: %v exceeds maximum %v",
			m.fileModel.PanelCount(), m.fileModel.MaxFilePanel)
	}

	totalFileModelWidth := 0
	// Check preview panel dimensions if open
	if m.fileModel.FilePreview.IsOpen() {
		if m.fileModel.FilePreview.GetWidth() <= 0 {
			return fmt.Errorf("preview panel is open but width is %v", m.fileModel.FilePreview.GetWidth())
		}
		if m.fileModel.FilePreview.GetHeight() <= 0 {
			return fmt.Errorf("preview panel is open but height is %v", m.fileModel.FilePreview.GetHeight())
		}
		totalFileModelWidth += m.fileModel.FilePreview.GetWidth()
	}

	// Check each file panel has correct dimensions set
	for i, panel := range m.fileModel.FilePanels {
		totalFileModelWidth += panel.GetWidth()
		if panel.GetHeight() != m.fileModel.Height {
			return fmt.Errorf(
				"file panel %v height mismatch: expected %v, got %v",
				i, m.mainPanelHeight, panel.GetHeight())
		}

		// Validate search bar width matches panel width minus padding
		if panel.SearchBar.Value() != "" && panel.SearchBar.Width != m.fileModel.Width-common.InnerPadding {
			return fmt.Errorf("file panel %v search bar width mismatch: expected %v, got %v",
				i, m.fileModel.Width-common.InnerPadding, panel.SearchBar.Width)
		}
	}

	if m.fileModel.Width != totalFileModelWidth {
		return fmt.Errorf(
			"file model width mismatch: expected %v, got %v",
			m.fileModel.Width, totalFileModelWidth)
	}

	// Validate focus panel index is within valid range
	if m.fileModel.FocusedPanelIndex < 0 || m.fileModel.FocusedPanelIndex >= m.fileModel.PanelCount() {
		return fmt.Errorf("FocusedPanelIndex %v is out of range [0, %v)",
			m.fileModel.FocusedPanelIndex, m.fileModel.PanelCount())
	}

	// Validate overlay panels have less width and height than total
	if m.helpMenu.open {
		if m.helpMenu.width >= m.fullWidth {
			return fmt.Errorf("help menu width %v exceeds full width %v", m.helpMenu.width, m.fullWidth)
		}
		if m.helpMenu.height >= m.fullHeight {
			return fmt.Errorf("help menu height %v exceeds full height %v", m.helpMenu.height, m.fullHeight)
		}
	}

	if m.promptModal.IsOpen() {
		if m.promptModal.GetWidth() >= m.fullWidth {
			return fmt.Errorf("prompt modal width %v exceeds full width %v", m.promptModal.GetWidth(), m.fullWidth)
		}
	}

	return nil
}

// validateRender validates rendered output dimensions and border
func validateRender(
	out string,
	height int,
	width int,
	border bool,
) error { //nolint:gocognit // cumilation of validations
	// Strip ANSI escape sequences for accurate width measurement
	strippedOut := ansi.Strip(out)

	// Split into lines
	lines := strings.Split(strippedOut, "\n")

	// Check number of lines
	if len(lines) != height {
		return fmt.Errorf("render height mismatch: expected %v lines, got %v", height, len(lines))
	}

	// Check each line width
	for i, line := range lines {
		lineWidth := ansi.StringWidth(line)
		if lineWidth != width {
			return fmt.Errorf("render line %v width mismatch: expected %v, got %v", i, width, lineWidth)
		}
	}

	// Check for border if required
	if border && len(lines) > 2 {
		// Check first line starts with TopLeft and ends with TopRight
		if !strings.HasPrefix(lines[0], common.Config.BorderTopLeft) {
			return fmt.Errorf("render missing top left border, expected %q", common.Config.BorderTopLeft)
		}
		if !strings.HasSuffix(lines[0], common.Config.BorderTopRight) {
			return fmt.Errorf("render missing top right border, expected %q", common.Config.BorderTopRight)
		}

		// Check last line starts with BottomLeft and ends with BottomRight
		lastLine := lines[len(lines)-1]
		if !strings.HasPrefix(lastLine, common.Config.BorderBottomLeft) {
			return fmt.Errorf("render missing bottom left border, expected %q", common.Config.BorderBottomLeft)
		}
		if !strings.HasSuffix(lastLine, common.Config.BorderBottomRight) {
			return fmt.Errorf("render missing bottom right border, expected %q", common.Config.BorderBottomRight)
		}

		// Check middle lines wrapped with BorderLeft and BorderRight
		for i := 1; i < len(lines)-1; i++ {
			if !strings.HasPrefix(lines[i], common.Config.BorderLeft) &&
				!strings.HasPrefix(lines[i], common.Config.BorderMiddleLeft) {
				return fmt.Errorf("render line '%v' missing left border", lines[i])
			}
			if !strings.HasSuffix(lines[i], common.Config.BorderRight) &&
				!strings.HasSuffix(lines[i], common.Config.BorderMiddleRight) {
				return fmt.Errorf("render line '%v' missing right border", lines[i])
			}
		}

		// Check top line contains BorderTop
		if !strings.Contains(lines[0], common.Config.BorderTop) {
			return fmt.Errorf("render missing top border character %q", common.Config.BorderTop)
		}

		// Check bottom line contains BorderBottom
		if !strings.Contains(lastLine, common.Config.BorderBottom) {
			return fmt.Errorf("render missing bottom border character %q", common.Config.BorderBottom)
		}
	}

	return nil
}

// validateComponentRender validates render output of all components
func (m *model) validateComponentRender() error { //nolint:gocognit // cumilation of validations
	// Skip validation if dimensions are not set
	if m.mainPanelHeight <= 0 || m.fullWidth <= 0 {
		return nil
	}

	// Validate sidebar render
	if common.Config.SidebarWidth > 0 && m.mainPanelHeight > 0 {
		sidebarRender := m.sidebarRender()
		// Sidebar render includes borders, so add BorderPadding to expected height
		// and width includes border padding as well
		if err := validateRender(sidebarRender, m.mainPanelHeight+common.BorderPadding, common.Config.SidebarWidth+common.BorderPadding, true); err != nil {
			return fmt.Errorf("sidebar render validation failed: %w", err)
		}
	}

	// Validate file panels
	if m.fileModel.PanelCount() > 0 && m.mainPanelHeight > 0 {
		for i := range m.fileModel.FilePanels {
			panel := &m.fileModel.FilePanels[i]
			if panel.GetWidth() > 0 && panel.GetHeight() > 0 {
				panelRender := panel.Render(i == m.fileModel.FocusedPanelIndex)
				if err := validateRender(panelRender, panel.GetHeight(), panel.GetWidth(), true); err != nil {
					return fmt.Errorf("file panel %v render validation failed: %w", i, err)
				}
			}
		}
	}

	// Validate footer components if visible
	if m.toggleFooter && m.footerHeight > 0 {
		// Validate process bar
		processBarRender := m.processBarRender()
		processBarWidth := utils.FooterWidth(m.fullWidth) + common.BorderPadding
		processBarHeight := m.footerHeight + common.BorderPadding
		if processBarWidth > 0 && processBarHeight > 0 {
			if err := validateRender(processBarRender, processBarHeight, processBarWidth, true); err != nil {
				return fmt.Errorf("process bar render validation failed: %w", err)
			}
		}
	}

	return nil
}
