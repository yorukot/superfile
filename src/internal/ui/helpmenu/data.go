package helpmenu

import (
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Return help menu for Hotkeys
func getData() []hotkeydata { //nolint: funlen // This should be self contained
	data := []hotkeydata{
		{
			subTitle: "General",
		},
		{
			hotkey:         []string{"spf", ""},
			description:    "Open superfile",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.Confirm,
			description:    "Confirm your select or typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.Quit,
			description:    "Quit typing, modal or superfile",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CdQuit,
			description:    "Quit superfile and change directory to current folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ConfirmTyping,
			description:    "Confirm typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CancelTyping,
			description:    "Cancel typing",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenHelpMenu,
			description:    "Open help menu (hotkeylist)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenCommandLine,
			description:    "Open command line",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenSPFPrompt,
			description:    "Open SPF prompt",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenZoxide,
			description:    "Open zoxide navigation",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel navigation",
		},
		{
			hotkey:         common.Hotkeys.CreateNewFilePanel,
			description:    "Create new file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CloseFilePanel,
			description:    "Close the focused file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleFilePreviewPanel,
			description:    "Toggle file preview panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.OpenSortOptionsMenu,
			description:    "Open sort options menu",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleReverseSort,
			description:    "Toggle reverse sort",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleFooter,
			description:    "Toggle footer",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.NextFilePanel,
			description:    "Focus on the next file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PreviousFilePanel,
			description:    "Focus on the previous file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FocusOnProcessBar,
			description:    "Focus on the processbar panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FocusOnSidebar,
			description:    "Focus on the sidebar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FocusOnMetaData,
			description:    "Focus on the metadata panel",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "Panel movement",
		},
		{
			hotkey:         common.Hotkeys.ListUp,
			description:    "Up",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ListDown,
			description:    "Down",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PageUp,
			description:    "Page up",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PageDown,
			description:    "Page down",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ParentDirectory,
			description:    "Return to parent folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelSelectAllItem,
			description:    "Select all items in focused file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelSelectModeItemsSelectUp,
			description:    "Select up with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelSelectModeItemsSelectDown,
			description:    "Select down with your course",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ToggleDotFile,
			description:    "Toggle dot file display",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.SearchBar,
			description:    "Toggle active search bar",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ChangePanelMode,
			description:    "Change between selection mode or normal mode",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PinnedDirectory,
			description:    "Pin or Unpin folder to sidebar (can be auto saved)",
			hotkeyWorkType: globalType,
		},
		{
			subTitle: "File operations",
		},
		{
			hotkey:         common.Hotkeys.FilePanelItemCreate,
			description:    "Create file or folder(end with " + string(filepath.Separator) + " to create a folder)",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.FilePanelItemRename,
			description:    "Rename file or folder",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CopyItems,
			description:    "Copy selected items to the clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CutItems,
			description:    "Cut selected items to the clipboard",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PasteItems,
			description:    "Paste clipboard items into the current file panel",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.DeleteItems,
			description:    "Delete selected items",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.PermanentlyDeleteItems,
			description:    "Permanently delete selected items",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CopyPath,
			description:    "Copy current file or directory path",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.CopyPWD,
			description:    "Copy current working directory",
			hotkeyWorkType: globalType,
		},
		{
			hotkey:         common.Hotkeys.ExtractFile,
			description:    "Extract compressed file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         common.Hotkeys.CompressFile,
			description:    "Zip file or folder to .zip file",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         common.Hotkeys.OpenFileWithEditor,
			description:    "Open file with your default editor",
			hotkeyWorkType: normalType,
		},
		{
			hotkey:         common.Hotkeys.OpenCurrentDirectoryWithEditor,
			description:    "Open current directory with default editor",
			hotkeyWorkType: normalType,
		},
	}

	return data
}

func removeOrphanSections(items []hotkeydata) []hotkeydata {
	var result []hotkeydata
	// Since we can't know beforehand which section are we actually filtering
	// we may end up in a scenario where there are two sections (General, Panel navigation)
	// with no hotkeys between them, so we need to remove the section which its hotkeys was
	// completely filtered out (Orphan sections)
	for i := range items {
		if items[i].subTitle != "" {
			// Look ahead: is the next item a real hotkey?
			if i+1 < len(items) && items[i+1].subTitle == "" {
				result = append(result, items[i])
			}
			// Else: skip this subtitle because no children
		} else {
			result = append(result, items[i])
		}
	}
	return result
}

func (m *Model) filter(query string) {
	filtered := fuzzySearch(query, m.data)
	filtered = removeOrphanSections(filtered)

	m.filteredData = filtered
	if len(filtered) == 0 {
		m.cursor = 0
	} else {
		m.cursor = 1
	}
	m.renderIndex = 0
}

// Fuzzy search function for a list of helpMenuModalData.
// inspired from: sidebar/directory_utils.go
func fuzzySearch(query string, data []hotkeydata) []hotkeydata {
	if len(data) == 0 {
		return []hotkeydata{}
	}

	// Optimization - This haystack can be kept precomputed based on description
	// instead of re computing it in each call
	haystack := []string{}
	idxMap := []int{}
	for i, item := range data {
		if item.subTitle != "" {
			continue
		}
		searchText := strings.Join(item.hotkey, " ") + " " + item.description
		haystack = append(haystack, searchText)
		idxMap = append(idxMap, i)
	}

	matchedIdx := map[int]struct{}{}
	for _, match := range utils.FzfSearch(query, haystack) {
		matchedIdx[idxMap[match.HayIndex]] = struct{}{}
	}

	results := []hotkeydata{}
	for i, d := range data {
		_, isMatch := matchedIdx[i]
		if d.subTitle != "" || isMatch {
			results = append(results, d)
		}
	}

	return results
}
