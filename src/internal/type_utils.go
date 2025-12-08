package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

const invalidTypeString = "InvalidType"

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}

// ================ Model related utils =======================

// Non fatal Validations. This indicates bug / programming errors, not user configuration mistake
func (m *model) validateLayout() error {
	if 0 < m.footerHeight && m.footerHeight < common.MinFooterHeight {
		return fmt.Errorf("footerHeight %v is too small", m.footerHeight)
	}
	if !m.toggleFooter && m.footerHeight != 0 {
		return fmt.Errorf("footer closed and footerHeight %v is non zero", m.footerHeight)
	}
	// PanelHeight + 2 lines (main border) + actual footer height
	if m.fullHeight != (m.mainPanelHeight+2)+utils.FullFooterHeight(m.footerHeight, m.toggleFooter) {
		return fmt.Errorf("invalid model layout, fullHeight : %v, mainPanelHeight : %v, footerHeight : %v",
			m.fullHeight, m.mainPanelHeight, m.footerHeight)
	}
	// TODO : Add check for width as well
	return nil
}

// ================ filepanel

func filePanelSlice(paths []string) []filePanel {
	res := make([]filePanel, len(paths))
	for i := range paths {
		// Making the first panel as the focussed
		isFocus := i == 0
		res[i] = defaultFilePanel(paths[i], isFocus)
	}
	return res
}

func defaultFilePanel(path string, focused bool) filePanel {
	targetFile := ""
	panelPath := path
	// If path refers to a file, switch to its parent and remember the filename
	if stat, err := os.Stat(panelPath); err == nil && !stat.IsDir() {
		targetFile = filepath.Base(panelPath)
		panelPath = filepath.Dir(panelPath)
	}

	return filePanel{
		render:   0,
		cursor:   0,
		location: panelPath,
		sortOptions: sortOptionsModel{
			width:  20,
			height: 4,
			open:   false,
			cursor: common.Config.DefaultSortType,
			data: sortOptionsModelData{
				options: []string{
					string(sortingName), string(sortingSize),
					string(sortingDateModified), string(sortingFileType),
				},
				selected: common.Config.DefaultSortType,
				reversed: common.Config.SortOrderReversed,
			},
		},
		panelMode:        browserMode,
		isFocused:        focused,
		directoryRecords: make(map[string]directoryRecord),
		searchBar:        common.GenerateSearchBar(),
		targetFile:       targetFile,
	}
}

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
		return invalidTypeString
	}
}

func (p panelMode) String() string {
	switch p {
	case selectMode:
		return "selectMode"
	case browserMode:
		return "browserMode"
	default:
		return invalidTypeString
	}
}
