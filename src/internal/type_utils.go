package internal

import (
	"errors"
	"fmt"

	"github.com/yorukot/superfile/src/internal/common"

	"github.com/yorukot/superfile/src/internal/utils"
)

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}

// ================ Model related utils =======================

// Non fatal Validations. This indicates bug / programming errors, not user configuration mistake
func (m *model) validateLayout() error { //nolint:gocognit // cumilation of validation
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

// ================ filepanel

// ================ String method for easy logging =====================

func (f focusPanelType) String() string {
	switch f {
	case nonePanelFocus:
		return "nonePanelFocus"
	case processBarFocus:
		return "processBarFocus"
	case sidebarFocus:
		return "sidebarFocus"
	case metadataFocus:
		return "metadataFocus"
	default:
		return common.InvalidTypeString
	}
}
