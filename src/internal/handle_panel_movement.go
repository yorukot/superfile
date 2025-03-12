package internal

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	variable "github.com/yorukot/superfile/src/config"
)

// Change file panel mode (select mode or browser mode)
func (m *model) changeFilePanelMode() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.panelMode == selectMode {
		panel.selected = panel.selected[:0]
		panel.panelMode = browserMode
	} else if panel.panelMode == browserMode {
		panel.panelMode = selectMode
	}
}

// Back to parent directory
func (m *model) parentDirectory() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.directoryRecord[panel.location] = directoryRecord{
		directoryCursor: panel.cursor,
		directoryRender: panel.render,
	}
	fullPath := panel.location
	parentDir := filepath.Dir(fullPath)
	panel.location = parentDir
	directoryRecord, hasRecord := panel.directoryRecord[panel.location]
	if hasRecord {
		panel.cursor = directoryRecord.directoryCursor
		panel.render = directoryRecord.directoryRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}
}

// Enter directory or open file with default application
func (m *model) enterPanel() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if len(panel.element) == 0 {
		return
	}

	if panel.element[panel.cursor].directory {
		panel.directoryRecord[panel.location] = directoryRecord{
			directoryCursor: panel.cursor,
			directoryRender: panel.render,
		}
		panel.location = panel.element[panel.cursor].location
		directoryRecord, hasRecord := panel.directoryRecord[panel.location]
		if hasRecord {
			panel.cursor = directoryRecord.directoryCursor
			panel.render = directoryRecord.directoryRender
		} else {
			panel.cursor = 0
			panel.render = 0
		}
		panel.searchBar.SetValue("")
	} else if !panel.element[panel.cursor].directory {
		fileInfo, err := os.Lstat(panel.element[panel.cursor].location)
		if err != nil {
			slog.Error("Error while getting file info", "error", err)
			return
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			targetPath, err := filepath.EvalSymlinks(panel.element[panel.cursor].location)
			if err != nil {
				return
			}

			targetInfo, err := os.Lstat(targetPath)

			if err != nil {
				return
			}

			if targetInfo.IsDir() {
				m.fileModel.filePanels[m.filePanelFocusIndex].location = targetPath
			}

			return
		}

		openCommand := "xdg-open"
		if runtime.GOOS == "darwin" {
			openCommand = "open"
		} else if runtime.GOOS == "windows" {

			dllpath := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
			dllfile := "url.dll,FileProtocolHandler"

			cmd := exec.Command(dllpath, dllfile, panel.element[panel.cursor].location)
			err = cmd.Start()
			if err != nil {
				slog.Error("Error while open file with", "error", err)
			}

			return
		}

		cmd := exec.Command(openCommand, panel.element[panel.cursor].location)
		err = cmd.Start()
		if err != nil {
			slog.Error("Error while open file with", "error", err)
		}

	}

}

// Switch to the directory where the sidebar cursor is located
func (m *model) sidebarSelectDirectory() {
	// We can't do this when we have only divider directories
	// m.sidebarModel.directories[m.sidebarModel.cursor].location would point to a divider dir.
	if m.sidebarModel.noActualDir() {
		return
	}
	m.focusPanel = nonePanelFocus
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	panel.directoryRecord[panel.location] = directoryRecord{
		directoryCursor: panel.cursor,
		directoryRender: panel.render,
	}

	panel.location = m.sidebarModel.directories[m.sidebarModel.cursor].location
	directoryRecord, hasRecord := panel.directoryRecord[panel.location]
	if hasRecord {
		panel.cursor = directoryRecord.directoryCursor
		panel.render = directoryRecord.directoryRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}
	panel.focusType = focus
}

// Select all item in the file panel (only work on select mode)
func (m *model) selectAllItem() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	for _, item := range panel.element {
		panel.selected = append(panel.selected, item.location)
	}
}

// Select the item where cursor located (only work on select mode)
func (m *model) singleItemSelect() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex] // Access the current panel

	if len(panel.element) > 0 && panel.cursor >= 0 && panel.cursor < len(panel.element) {
		elementLocation := panel.element[panel.cursor].location

		if arrayContains(panel.selected, elementLocation) {
			panel.selected = removeElementByValue(panel.selected, elementLocation)
		} else {
			panel.selected = append(panel.selected, elementLocation)
		}
	}
}

// Toggle dotfile display or not
func (m *model) toggleDotFileController() {
	newToggleDotFile := ""
	if m.toggleDotFile {
		newToggleDotFile = "false"
		m.toggleDotFile = false
	} else {
		newToggleDotFile = "true"
		m.toggleDotFile = true
	}
	m.updatedToggleDotFile = true
	err := os.WriteFile(variable.ToggleDotFile, []byte(newToggleDotFile), 0644)
	if err != nil {
		slog.Error("Error while pinned folder function updatedData superfile data", "error", err)
	}

}

// Toggle dotfile display or not
func (m *model) toggleFooterController() {
	newToggleFooterFile := ""
	if m.toggleFooter {
		newToggleFooterFile = "false"
		m.toggleFooter = false
	} else {
		newToggleFooterFile = "true"
		m.toggleFooter = true
	}
	err := os.WriteFile(variable.ToggleFooter, []byte(newToggleFooterFile), 0644)
	if err != nil {
		slog.Error("Error while Toggle footer function updatedData superfile data", "error", err)
	}
	m.setFooterSize(m.fullHeight)

}

// Focus on search bar
func (m *model) searchBarFocus() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.searchBar.Focused() {
		panel.searchBar.Blur()
	} else {
		panel.searchBar.Focus()
		m.firstTextInput = true
	}

	// config search bar width
	panel.searchBar.Width = m.fileModel.width - 4
}

// ======================================== File panel controller ========================================

// Control file panel list up
func (panel *filePanel) listUp(mainPanelHeight int) {
	if len(panel.element) == 0 {
		return
	}
	if panel.cursor > 0 {
		panel.cursor--
		if panel.cursor < panel.render {
			panel.render--
		}
	} else {
		if len(panel.element) > panelElementHeight(mainPanelHeight) {
			panel.render = len(panel.element) - panelElementHeight(mainPanelHeight)
			panel.cursor = len(panel.element) - 1
		} else {
			panel.cursor = len(panel.element) - 1
		}
	}
}

// Control file panel list down
func (panel *filePanel) listDown(mainPanelHeight int) {
	if len(panel.element) == 0 {
		return
	}
	if panel.cursor < len(panel.element)-1 {
		panel.cursor++
		if panel.cursor > panel.render+panelElementHeight(mainPanelHeight)-1 {
			panel.render++
		}
	} else {
		panel.render = 0
		panel.cursor = 0
	}
}

func (panel *filePanel) pgUp(mainPanelHeight int) {
	panlen := len(panel.element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		panel.cursor = 0
	} else {
		if panel.cursor-panHeight <= 0 {
			panel.cursor = 0
			panel.render = 0
		} else {
			panel.cursor -= panHeight
			panel.render = panel.cursor - panCenter

			if panel.render < 0 {
				panel.render = 0
			}
		}
	}
}

func (panel *filePanel) pgDown(mainPanelHeight int) {
	panlen := len(panel.element)
	if panlen == 0 {
		return
	}

	panHeight := panelElementHeight(mainPanelHeight)
	panCenter := panHeight / 2 // For making sure the cursor is at the center of the panel

	if panHeight >= panlen {
		panel.cursor = panlen - 1
	} else {
		if panel.cursor+panHeight >= panlen {
			panel.cursor = panlen - 1
			panel.render = panel.cursor - panCenter
		} else {
			panel.cursor += panHeight
			panel.render = panel.cursor - panCenter
		}
	}
}

// Handles the action of selecting an item in the file panel upwards. (only work on select mode)
func (m *model) itemSelectUp(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		panel := m.fileModel.filePanels[m.filePanelFocusIndex]
		if panel.cursor > 0 {
			panel.cursor--
			if panel.cursor < panel.render {
				panel.render--
			}
		} else {
			if len(panel.element) > panelElementHeight(m.mainPanelHeight) {
				panel.render = len(panel.element) - panelElementHeight(m.mainPanelHeight)
				panel.cursor = len(panel.element) - 1
			} else {
				panel.cursor = len(panel.element) - 1
			}
		}
		selectItemIndex := panel.cursor + 1
		if selectItemIndex > len(panel.element)-1 {
			selectItemIndex = 0
		}
		if arrayContains(panel.selected, panel.element[selectItemIndex].location) {
			panel.selected = removeElementByValue(panel.selected, panel.element[selectItemIndex].location)
		} else {
			panel.selected = append(panel.selected, panel.element[selectItemIndex].location)
		}

		m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	}
}

// Handles the action of selecting an item in the file panel downwards. (only work on select mode)
func (m *model) itemSelectDown(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		panel := m.fileModel.filePanels[m.filePanelFocusIndex]
		if panel.cursor < len(panel.element)-1 {
			panel.cursor++
			if panel.cursor > panel.render+panelElementHeight(m.mainPanelHeight)-1 {
				panel.render++
			}
		} else {
			panel.render = 0
			panel.cursor = 0
		}
		selectItemIndex := panel.cursor - 1
		if selectItemIndex < 0 {
			selectItemIndex = len(panel.element) - 1
		}
		if arrayContains(panel.selected, panel.element[selectItemIndex].location) {
			panel.selected = removeElementByValue(panel.selected, panel.element[selectItemIndex].location)
		} else {
			panel.selected = append(panel.selected, panel.element[selectItemIndex].location)
		}

		m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	}
}

func (m *model) sidebarSearchBarFocus() {

	if m.sidebarModel.searchBar.Focused() {
		// Ideally Code should never reach here. Once sidebar is focussed, we should
		// not cause sidebarSearchBarFocus() event by pressing search key
		// Should we use Runtime panic asserts ?
		slog.Error("sidebarSearchBarFocus() called on Focussed sidebar")
		m.sidebarModel.searchBar.Blur()
	} else {
		m.sidebarModel.searchBar.Focus()
		m.firstTextInput = true
	}
}

// ======================================== Metadata controller ========================================

// Control metadata panel up
func (fm *fileMetadata) listUp() {
	if len(fm.metaData) == 0 {
		return
	}
	if fm.renderIndex > 0 {
		fm.renderIndex--
	} else {
		fm.renderIndex = len(fm.metaData) - 1
	}
}

// Control metadata panel down
func (fm *fileMetadata) listDown() {
	if fm.renderIndex < len(fm.metaData)-1 {
		fm.renderIndex++
	} else {
		fm.renderIndex = 0
	}
}

// ======================================== Processbar controller ========================================

// Control processbar panel list up
func (p *processBarModel) listUp() {
	if len(p.processList) == 0 {
		return
	}
	if p.cursor > 0 {
		p.cursor--
		if p.cursor < p.render {
			p.render--
		}
	} else {
		if len(p.processList) <= 3 || (len(p.processList) <= 2 && footerHeight < 14) {
			p.cursor = len(p.processList) - 1
		} else {
			p.render = len(p.processList) - 3
			p.cursor = len(p.processList) - 1
		}
	}
}

// Control processbar panel list down
func (p *processBarModel) listDown() {
	if len(p.processList) == 0 {
		return
	}
	if p.cursor < len(p.processList)-1 {
		p.cursor++
		if p.cursor > p.render+2 {
			p.render++
		}
	} else {
		p.render = 0
		p.cursor = 0
	}
}
