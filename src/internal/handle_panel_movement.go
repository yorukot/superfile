package internal

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	variable "github.com/yorukot/superfile/src/config"
)

// Change file panel mode (select mode or browser mode)
func (m *model) changeFilePanelMode() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.panelMode == selectMode {
		panel.selected = panel.selected[:0]
		panel.panelMode = browserMode
	} else if panel.panelMode == browserMode {
		panel.panelMode = selectMode
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Back to parent directory
func (m *model) parentDirectory() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
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
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Enter directory or open file with default application
func (m *model) enterPanel() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

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
			outPutLog("err when getting file info", err)
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
				outPutLog(fmt.Sprintf("err when open file with %s", openCommand), err)
			}

			return
		}

		cmd := exec.Command(openCommand, panel.element[panel.cursor].location)
		err = cmd.Start()
		if err != nil {
			outPutLog(fmt.Sprintf("err when open file with %s", openCommand), err)
		}

	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Switch to the directory where the sidebar cursor is located
func (m *model) sidebarSelectDirectory() {
	m.focusPanel = nonePanelFocus
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

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
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Select all item in the file panel (only work on select mode)
func (m *model) selectAllItem() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	for _, item := range panel.element {
		panel.selected = append(panel.selected, item.location)
	}
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Select the item where cursor located (only work on select mode)
func (m *model) singleItemSelect() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex] // Access the current panel

	if len(panel.element) > 0 && panel.cursor >= 0 && panel.cursor < len(panel.element) {
		elementLocation := panel.element[panel.cursor].location

		if arrayContains(panel.selected, elementLocation) {
			panel.selected = removeElementByValue(panel.selected, elementLocation)
		} else {
			panel.selected = append(panel.selected, elementLocation)
		}

		m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	} else {
		outPutLog("No elements to select or cursor out of bounds.")
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
		outPutLog("Pinned folder function updatedData superfile data error", err)
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
		outPutLog("Toggle footer function updatedData superfile data error", err)
	}
	m.setFooterSize(m.fullHeight)

}

// Focus on search bar
func (m *model) searchBarFocus() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	slog.Debug("model.searchBarFocus", "focused", panel.searchBar.Focused())
	if panel.searchBar.Focused() {
		panel.searchBar.Blur()
	} else {
		panel.searchBar.Focus()
	}

	m.firstTextInput = true

	// config search bar width
	panel.searchBar.Width = m.fileModel.width - 4
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// ======================================== File panel controller ========================================

// Control file panel list up
func (m *model) controlFilePanelListUp(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		panel := m.fileModel.filePanels[m.filePanelFocusIndex]
		if len(panel.element) == 0 {
			return
		}
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

		m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	}
}

// Control file panel list down
func (m *model) controlFilePanelListDown(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		panel := m.fileModel.filePanels[m.filePanelFocusIndex]
		if len(panel.element) == 0 {
			return
		}
		if panel.cursor < len(panel.element)-1 {
			panel.cursor++
			if panel.cursor > panel.render+panelElementHeight(m.mainPanelHeight)-1 {
				panel.render++
			}
		} else {
			panel.render = 0
			panel.cursor = 0
		}
		m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	}

}

func (m *model) controlFilePanelPgUp(){
  panel := m.fileModel.filePanels[m.filePanelFocusIndex]
  panlen := len(panel.element)
  panHeight := panelElementHeight(m.mainPanelHeight)
  panCenter := panHeight / 2     // For making sure the cursor is at the center of the panel

  if panlen == 0 {
    return
  }

  if panHeight >= panlen {
    panel.cursor = 0
  } else {
    if panel.cursor - panHeight <= 0 {
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

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

func (m *model) controlFilePanelPgDown(){
  panel := m.fileModel.filePanels[m.filePanelFocusIndex]
  panlen := len(panel.element)
  panHeight := panelElementHeight(m.mainPanelHeight)
  panCenter := panHeight / 2     // For making sure the cursor is at the center of the panel

  if panlen == 0 {
    return
  }

  if panHeight >= panlen {
    panel.cursor = panlen - 1
  } else {
    if panel.cursor + panHeight >= panlen {
      panel.cursor = panlen - 1
      panel.render = panel.cursor - panCenter
    } else {
      panel.cursor += panHeight
      panel.render = panel.cursor - panCenter
    }
  }

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
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

// ======================================== Sidebar controller ========================================

// Yorukot: P.S God bless me, this sidebar controller code is really ugly...

// Control sidebar panel list up
func (m *model) controlSideBarListUp(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		if m.sidebarModel.cursor > 0 {
			m.sidebarModel.cursor--
		} else {
			m.sidebarModel.cursor = len(m.sidebarModel.directories) - 1
		}
		newDirectory := m.sidebarModel.directories[m.sidebarModel.cursor].location

		for newDirectory == "Pinned+-*/=?" || newDirectory == "Disks+-*/=?" {
			m.sidebarModel.cursor--
			newDirectory = m.sidebarModel.directories[m.sidebarModel.cursor].location
		}
		changeToPlus := false
		cursorRender := false
		for !cursorRender {
			totalHeight := 2
			for i := m.sidebarModel.renderIndex; i < len(m.sidebarModel.directories); i++ {
				if totalHeight >= m.mainPanelHeight {
					break
				}
				directory := m.sidebarModel.directories[i]

				if directory.location == "Pinned+-*/=?" {
					totalHeight += 3
					continue
				}

				if directory.location == "Disks+-*/=?" {
					if m.mainPanelHeight-totalHeight <= 2 {
						break
					}
					totalHeight += 3
					continue
				}

				totalHeight++
				if m.sidebarModel.cursor == i && m.focusPanel == sidebarFocus {
					cursorRender = true
				}
			}

			if changeToPlus {
				m.sidebarModel.renderIndex++
				continue
			}

			if !cursorRender {
				m.sidebarModel.renderIndex--
			}
			if m.sidebarModel.renderIndex < 0 {
				changeToPlus = true
				m.sidebarModel.renderIndex++
			}
		}

		if changeToPlus {
			m.sidebarModel.renderIndex--
		}
	}
}

// Control sidebar panel list down
func (m *model) controlSideBarListDown(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		lenDirs := len(m.sidebarModel.directories)
		if m.sidebarModel.cursor < lenDirs-1 {
			m.sidebarModel.cursor++
		} else {
			m.sidebarModel.cursor = 0
		}

		newDirectory := m.sidebarModel.directories[m.sidebarModel.cursor].location
		for newDirectory == "Pinned+-*/=?" || newDirectory == "Disks+-*/=?" {
			m.sidebarModel.cursor++
			if m.sidebarModel.cursor+1 > len(m.sidebarModel.directories) {
				m.sidebarModel.cursor = 0
			}
			newDirectory = m.sidebarModel.directories[m.sidebarModel.cursor].location
		}
		cursorRender := false
		for !cursorRender {
			totalHeight := 2
			for i := m.sidebarModel.renderIndex; i < len(m.sidebarModel.directories); i++ {
				if totalHeight >= m.mainPanelHeight {
					break
				}

				directory := m.sidebarModel.directories[i]

				if directory.location == "Pinned+-*/=?" {
					totalHeight += 3
					continue
				}

				if directory.location == "Disks+-*/=?" {
					if m.mainPanelHeight-totalHeight <= 2 {
						break
					}
					totalHeight += 3
					continue
				}

				totalHeight++
				if m.sidebarModel.cursor == i && m.focusPanel == sidebarFocus {
					cursorRender = true
				}
			}

			if !cursorRender {
				m.sidebarModel.renderIndex++
			}
			if m.sidebarModel.renderIndex > m.sidebarModel.cursor {
				m.sidebarModel.renderIndex = 0
			}
		}
	}
}

// ======================================== Metadata controller ========================================

// Control metadata panel up
func (m *model) controlMetadataListUp(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	if len(m.fileMetaData.metaData) == 0 {
		return
	}

	for i := 0; i < runTime; i++ {
		if m.fileMetaData.renderIndex > 0 {
			m.fileMetaData.renderIndex--
		} else {
			m.fileMetaData.renderIndex = len(m.fileMetaData.metaData) - 1
		}
	}
}

// Control metadata panel down
func (m *model) controlMetadataListDown(wheel bool) {
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		if m.fileMetaData.renderIndex < len(m.fileMetaData.metaData)-1 {
			m.fileMetaData.renderIndex++
		} else {
			m.fileMetaData.renderIndex = 0
		}
	}
}

// ======================================== Processbar controller ========================================

// Control processbar panel list up
func (m *model) controlProcessbarListUp(wheel bool) {
	if len(m.processBarModel.processList) == 0 {
		return
	}
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		if m.processBarModel.cursor > 0 {
			m.processBarModel.cursor--
			if m.processBarModel.cursor < m.processBarModel.render {
				m.processBarModel.render--
			}
		} else {
			if len(m.processBarModel.processList) <= 3 || (len(m.processBarModel.processList) <= 2 && footerHeight < 14) {
				m.processBarModel.cursor = len(m.processBarModel.processList) - 1
			} else {
				m.processBarModel.render = len(m.processBarModel.processList) - 3
				m.processBarModel.cursor = len(m.processBarModel.processList) - 1
			}
		}
	}
}

// Control processbar panel list down
func (m *model) controlProcessbarListDown(wheel bool) {
	if len(m.processBarModel.processList) == 0 {
		return
	}

	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}

	for i := 0; i < runTime; i++ {
		if m.processBarModel.cursor < len(m.processBarModel.processList)-1 {
			m.processBarModel.cursor++
			if m.processBarModel.cursor > m.processBarModel.render+2 {
				m.processBarModel.render++
			}
		} else {
			m.processBarModel.render = 0
			m.processBarModel.cursor = 0
		}
	}
}
