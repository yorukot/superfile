package internal

import (
	"fmt"

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

func filePanelSlice(dir []string) []filePanel {
	res := make([]filePanel, len(dir))
	for i := range dir {
		// Making the first panel as the default focus panel
		// while others remain secondFocus
		isFocus := false
		if i == 0 {
			isFocus = true
		}
		res[i] = defaultFilePanel(dir[i], isFocus)
	}
	return res
}

func defaultFilePanel(dir string, focused bool) filePanel {
	return filePanel{
		render:   0,
		cursor:   0,
		location: dir,
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

//func (f filePanelFocusType) String() string {
//switch f {
//case noneFocus:
//return "noneFocus"
//case secondFocus:
//return "secondFocus"
//case focus:
//return "focus"
//default:
//return invalidTypeString
//}
//}

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
