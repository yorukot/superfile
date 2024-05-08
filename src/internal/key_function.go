package internal

import tea "github.com/charmbracelet/bubbletea"

func mainKey(msg string, m model, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg {

	case hotkeys.ListUp[0], hotkeys.ListUp[1]:
		if m.focusPanel == sidebarFocus {
			m = controlSideBarListUp(m)
		} else if m.focusPanel == processBarFocus {
			m = contolProcessbarListUp(m)
		} else if m.focusPanel == metadataFocus {
			m = controlMetadataListUp(m)
		} else if m.focusPanel == nonePanelFocus {
			m = controlFilePanelListUp(m)
			m.fileMetaData.renderIndex = 0
			go func() {
				m = returnMetaData(m)
			}()
		}

	case hotkeys.ListDown[0], hotkeys.ListDown[1]:
		if m.focusPanel == sidebarFocus {
			m = controlSideBarListDown(m)
		} else if m.focusPanel == processBarFocus {
			m = contolProcessbarListDown(m)
		} else if m.focusPanel == metadataFocus {
			m = controlMetadataListDown(m)
		} else if m.focusPanel == nonePanelFocus {
			m = controlFilePanelListDown(m)
			m.fileMetaData.renderIndex = 0
			go func() {
				m = returnMetaData(m)
			}()
		}

	case hotkeys.ChangePanelMode[0], hotkeys.ChangePanelMode[1]:
		m = changeFilePanelMode(m)

	case hotkeys.NextFilePanel[0], hotkeys.NextFilePanel[1]:
		m = nextFilePanel(m)

	case hotkeys.PreviousFilePanel[0], hotkeys.PreviousFilePanel[1]:
		m = previousFilePanel(m)

	case hotkeys.CloseFilePanel[0], hotkeys.CloseFilePanel[1]:
		m = closeFilePanel(m)

	case hotkeys.CreateNewFilePanel[0], hotkeys.CreateNewFilePanel[1]:
		m = createNewFilePanel(m)

	case hotkeys.FocusOnSideBar[0], hotkeys.FocusOnSideBar[1]:
		m = focusOnSideBar(m)

	case hotkeys.FocusOnProcessBar[0], hotkeys.FocusOnProcessBar[1]:
		m = focusOnProcessBar(m)

	case hotkeys.FocusOnMetaData[0], hotkeys.FocusOnMetaData[1]:
		m = focusOnMetadata(m)
		go func() {
			m = returnMetaData(m)
		}()

	case hotkeys.PasteItem[0], hotkeys.PasteItem[1]:
		go func() {
			m = pasteItem(m)
		}()

	case hotkeys.FilePanelFileCreate[0], hotkeys.FilePanelFileCreate[1]:
		m = panelCreateNewFile(m)

	case hotkeys.FilePanelDirectoryCreate[0], hotkeys.FilePanelDirectoryCreate[1]:
		m = panelCreateNewDirectory(m)

	case hotkeys.PinnedDirectory[0], hotkeys.PinnedDirectory[1]:
		m = pinnedDirectory(m)

	case hotkeys.ToggleDotFile[0], hotkeys.ToggleDotFile[1]:
		m = toggleDotFileController(m)

	case hotkeys.ExtractFile[0], hotkeys.ExtractFile[1]:
		go func() {
			m = extractFile(m)
		}()

	case hotkeys.CompressFile[0], hotkeys.CompressFile[1]:
		go func() {
			m = compressFile(m)
		}()

	case hotkeys.OpenHelpMenu[0], hotkeys.OpenHelpMenu[1]:
		m = openHelpMenu(m)

	case hotkeys.OpenFileWithEditor[0], hotkeys.OpenFileWithEditor[1]:
		cmd = openFileWithEditor(m)

	case hotkeys.OpenCurrentDirectoryWithEditor[0], hotkeys.OpenCurrentDirectoryWithEditor[1]:
		cmd = openDirectoryWithEditor(m)

	default:
		m = normalAndBrowserModeKey(msg, m)
	}

	return m, cmd
}

func normalAndBrowserModeKey(msg string, m model) model {
	// if not focus on the filepanel return
	if m.fileModel.filePanels[m.filePanelFocusIndex].focusType != focus {
		if m.focusPanel == sidebarFocus && (msg == hotkeys.SelectItem[0] || msg == hotkeys.SelectItem[1]) {
			m = sidebarSelectDirectory(m)
		}
		return m
	}
	// Check if in the select mode and focusOn filepanel
	if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
		switch msg {
		case hotkeys.FilePanelSelectModeItemSingleSelect[0], hotkeys.FilePanelSelectModeItemSingleSelect[1]:
			m = singleItemSelect(m)
		case hotkeys.FilePanelSelectModeItemSelectUp[0], hotkeys.FilePanelSelectModeItemSelectUp[1]:
			m = itemSelectUp(m)
		case hotkeys.FilePanelSelectModeItemSelectDown[0], hotkeys.FilePanelSelectModeItemSelectDown[1]:
			m = itemSelectDown(m)
		case hotkeys.FilePanelSelectModeItemDelete[0], hotkeys.FilePanelSelectModeItemDelete[1]:
			go func() {
				m = deleteMultipleItem(m)
				if !isExternalDiskPath(m.fileModel.filePanels[m.filePanelFocusIndex].location) {
					m.fileModel.filePanels[m.filePanelFocusIndex].selected = m.fileModel.filePanels[m.filePanelFocusIndex].selected[:0]
				}
			}()
		case hotkeys.FilePanelSelectModeItemCopy[0], hotkeys.FilePanelSelectModeItemCopy[1]:
			m = copyMultipleItem(m)
		case hotkeys.FilePanelSelectModeItemCut[0], hotkeys.FilePanelSelectModeItemCut[1]:
			m = cutMultipleItem(m)
		case hotkeys.FilePanelSelectAllItem[0], hotkeys.FilePanelSelectAllItem[1]:
			m = selectAllItem(m)
		}
		return m
	}

	switch msg {
	case hotkeys.SelectItem[0], hotkeys.SelectItem[1]:
		forceReloadElement = true
		m = enterPanel(m)
	case hotkeys.ParentDirectory[0], hotkeys.ParentDirectory[1]:
		forceReloadElement = true
		m = parentDirectory(m)
	case hotkeys.DeleteItem[0], hotkeys.DeleteItem[1]:
		go func() {
			m = deleteSingleItem(m)
		}()
	case hotkeys.CopySingleItem[0], hotkeys.CopySingleItem[1]:
		m = copySingleItem(m)
	case hotkeys.CutSingleItem[0], hotkeys.CutSingleItem[1]:
		m = cutSingleItem(m)
	case hotkeys.FilePanelItemRename[0], hotkeys.FilePanelItemRename[1]:
		m = panelItemRename(m)
	case hotkeys.SearchBar[0], hotkeys.SearchBar[1]:
		m = searchBarFocus(m)
	}
	return m
}

func typingModalOpenKey(msg string, m model) model {
	switch msg {
	case hotkeys.Cancel[0], hotkeys.Cancel[1]:
		m = cancelTypingModal(m)
	case hotkeys.Confirm[0], hotkeys.Confirm[1]:
		m = createItem(m)
	}
	return m
}

func warnModalOpenKey(msg string, m model) model {
	switch msg {
	case hotkeys.Cancel[0], hotkeys.Cancel[1]:
		m = cancelWarnModal(m)
	case hotkeys.Confirm[0], hotkeys.Confirm[1]:
		m.warnModal.open = false
		if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
			go func() {
				m = completelyDeleteMultipleFile(m)
				m.fileModel.filePanels[m.filePanelFocusIndex].selected = m.fileModel.filePanels[m.filePanelFocusIndex].selected[:0]
			}()
		} else {
			go func() {
				m = completelyDeleteSingleFile(m)
			}()
		}
	}
	return m
}

func renamingKey(msg string, m model) model {
	switch msg {
	case hotkeys.Cancel[0], hotkeys.Cancel[1]:
		m = cancelReanem(m)
	case hotkeys.Confirm[0], hotkeys.Confirm[1]:
		m = confirmRename(m)
	}
	return m
}

func focusOnSearchbarKey(msg string, m model) model {
	switch msg {
	case hotkeys.Cancel[0], hotkeys.Cancel[1]:
		m = cancelSearch(m)
	case hotkeys.Confirm[0], hotkeys.Confirm[1], hotkeys.SearchBar[0], hotkeys.SearchBar[1]:
		m = confirmSearch(m)
	}
	return m
}

func helpMenuKey(msg string, m model) model {
	switch msg {
	case hotkeys.ListUp[0], hotkeys.ListUp[1]:
		m = helpMenuListUp(m)
	case hotkeys.ListDown[0], hotkeys.ListDown[1]:
		m = helpMenuListDown(m)
	case hotkeys.Quit[0], hotkeys.Quit[1]:
		m = quitHelpMenu(m)
	}

	return m
}
