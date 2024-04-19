package components

import (
	"os"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rkoesters/xdg/basedir"
)

const (
	configFolder     string = "/config"
	themeFolder      string = "/theme"
	dataFolder       string = "/data"
	lastCheckVersion string = "/lastCheckVersion"
	pinnedFile       string = "/pinned.json"
	toggleDotFile    string = "/toggleDotFile"
	logFile          string = "/superfile.log"
	configFile       string = "/config.toml"
	hotkeysFile      string = "/hotkeys.toml"
	themeZipName     string = "/theme.zip"
)

var LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), 0}
var ListeningMessage = true

var HomeDir = basedir.Home
var SuperFileMainDir = basedir.ConfigHome + "/superfile"
var SuperFileCacheDir = basedir.CacheHome + "/superfile"
var SuperFileDataDir = basedir.DataHome + "/superfile"

var theme ThemeType
var Config ConfigType
var hotkeys HotkeysType

var logOutput *os.File
var et *exiftool.Exiftool

var channel = make(chan channelMessage, 1000)

var foreceReloadElement = false

func InitialModel(dir string) model {
	toggleDotFileBool, firstFilePanelDir := loadConfigFile(dir)

	return model{
		filePanelFocusIndex: 0,
		focusPanel:          nonePanelFocus,
		processBarModel: processBarModel{
			process: make(map[string]process),
			cursor:  0,
			render:  0,
		},
		sidebarModel: sidebarModel{
			directories: getDirectories(40),
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{
					render:          0,
					cursor:          0,
					location:        firstFilePanelDir,
					panelMode:       browserMode,
					focusType:       focus,
					directoryRecord: make(map[string]directoryRecord),
					searchBar:       generateSearchBar(),
				},
			},
			width: 10,
		},
		toggleDotFile: toggleDotFileBool,
	}
}

func listenForChannelMessage(msg chan channelMessage) tea.Cmd {
	return func() tea.Msg {
		m := <-msg
		ListeningMessage = false
		return m
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("SuperFile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
		listenForChannelMessage(channel),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	switch msg := msg.(type) {
	// check is the message by thread
	case channelMessage:
		if msg.returnWarnModal {
			m.warnModal = msg.warnModal
		} else if msg.loadMetadata {
			m.fileMetaData.metaData = msg.metadata
		} else {
			if !arrayContains(m.processBarModel.processList, msg.messageId) {
				m.processBarModel.processList = append(m.processBarModel.processList, msg.messageId)
			}
			m.processBarModel.process[msg.messageId] = msg.processNewState
		}
	// if the message by windows size change
	case tea.WindowSizeMsg:
		if msg.Height < 30 {
			footerHeight = 10
		} else {
			footerHeight = 14
		}
		m.mainPanelHeight = msg.Height - footerHeight + 1
		m.fileModel.width = (msg.Width - sidebarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		m.fileModel.maxFilePanel = (msg.Width - 20) / 24
		if m.fileModel.maxFilePanel >= 10 {
			m.fileModel.maxFilePanel = 10
		}
		return m, nil
	// if just user press key
	case tea.KeyMsg:
		// if in the create item modal
		if m.typingModal.open {
			switch msg.String() {
			case hotkeys.Cancel[0], hotkeys.Cancel[1]:
				m = cancelTypingModal(m)
			case hotkeys.Confirm[0], hotkeys.Confirm[1]:
				m = createItem(m)
			}
			// if in the renaming mode
		} else if m.warnModal.open {
			switch msg.String() {
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
			// if in the renaming mode
		} else if m.fileModel.renaming {
			switch msg.String() {
			case hotkeys.Cancel[0], hotkeys.Cancel[1]:
				m = cancelReanem(m)
			case hotkeys.Confirm[0], hotkeys.Confirm[1]:
				m = confirmRename(m)
			}
			// if search bar focus
		} else if panel.searchBar.Focused() {
			switch msg.String() {
			case hotkeys.Cancel[0], hotkeys.Cancel[1]:
				m = cancelSearch(m)
			case hotkeys.Confirm[0], hotkeys.Confirm[1], hotkeys.SearchBar[0], hotkeys.SearchBar[1]:
				m = confirmSearch(m)
			}
		} else {
			switch msg.String() {
			// return superfile
			case hotkeys.Quit[0], hotkeys.Quit[1]:
				return m, tea.Quit
			/* LIST CONTROLLER START */
			// up list
			case hotkeys.ListUp[0], hotkeys.ListUp[1]:
				if m.focusPanel == sidebarFocus {
					m = controllerSideBarListUp(m)
				} else if m.focusPanel == processBarFocus {
					m = contollerProcessBarListUp(m)
				} else if m.focusPanel == metadataFocus {
					m = controllerMetaDataListUp(m)
				} else if m.focusPanel == nonePanelFocus {
					m = controllerFilePanelListUp(m)
					m.fileMetaData.renderIndex = 0
					go func() {
						m = returnMetaData(m)
					}()
				}
			// down list
			case hotkeys.ListDown[0], hotkeys.ListDown[1]:
				if m.focusPanel == sidebarFocus {
					m = controllerSideBarListDown(m)
				} else if m.focusPanel == processBarFocus {
					m = contollerProcessBarListDown(m)
				} else if m.focusPanel == metadataFocus {
					m = controllerMetaDataListDown(m)
				} else if m.focusPanel == nonePanelFocus {
					m = controllerFilePanelListDown(m)
					m.fileMetaData.renderIndex = 0
					go func() {
						m = returnMetaData(m)
					}()
				}
			/* LIST CONTROLLER END */
			case hotkeys.ChangePanelMode[0], hotkeys.ChangePanelMode[1]:
				m = selectedMode(m)
			/* NAVIGATION CONTROLLER START */
			// change file panel
			case hotkeys.NextFilePanel[0], hotkeys.NextFilePanel[1]:
				m = nextFilePanel(m)
			// change file panel
			case hotkeys.PreviousFilePanel[0], hotkeys.PreviousFilePanel[1]:
				m = previousFilePanel(m)
			// close file panel
			case hotkeys.CloseFilePanel[0], hotkeys.CloseFilePanel[1]:
				m = closeFilePanel(m)
			// create new file panel
			case hotkeys.CreateNewFilePanel[0], hotkeys.CreateNewFilePanel[1]:
				m = createNewFilePanel(m)
			// focus to sidebar or file panel
			case hotkeys.FocusOnSideBar[0], hotkeys.FocusOnSideBar[1]:
				m = focusOnSideBar(m)
			/* NAVIGATION CONTROLLER END */
			case hotkeys.FocusOnProcessBar[0], hotkeys.FocusOnProcessBar[1]:
				m = focusOnProcessBar(m)
			case hotkeys.FocusOnMetaData[0], hotkeys.FocusOnMetaData[1]:
				m = focusOnMetaData(m)
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
				m = panelCreateNewFolder(m)
			case hotkeys.PinnedDirectory[0], hotkeys.PinnedDirectory[1]:
				m = pinnedFolder(m)
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
			default:
				// check if it's the select mode
				if m.fileModel.filePanels[m.filePanelFocusIndex].focusType == focus && m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
					switch msg.String() {
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
					// else
				} else {
					switch msg.String() {
					case hotkeys.SelectItem[0], hotkeys.SelectItem[1]:
						if m.focusPanel == sidebarFocus {
							m = sidebarSelectFolder(m)
						} else if m.focusPanel == processBarFocus {

						} else if m.focusPanel == nonePanelFocus {
							foreceReloadElement = true
							m = enterPanel(m)
						}
					case hotkeys.ParentDirectory[0], hotkeys.ParentDirectory[1]:
						foreceReloadElement = true
						m = parentFolder(m)
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

				}
			}
		}
	}
	if m.firstTextInput {
		m.firstTextInput = false
	} else if m.fileModel.renaming {
		m.fileModel.filePanels[m.filePanelFocusIndex].rename, cmd = m.fileModel.filePanels[m.filePanelFocusIndex].rename.Update(msg)
	} else if panel.searchBar.Focused() {
		m.fileModel.filePanels[m.filePanelFocusIndex].searchBar, cmd = m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Update(msg)
	} else if m.typingModal.open {
		m.typingModal.textInput, cmd = m.typingModal.textInput.Update(msg)
	}

	if m.fileModel.filePanels[m.filePanelFocusIndex].cursor < 0 {
		m.fileModel.filePanels[m.filePanelFocusIndex].cursor = 0
	}

	cmd = tea.Batch(cmd)
	m.sidebarModel.directories = getDirectories(m.fullHeight)

	if !ListeningMessage {
		cmd = tea.Batch(cmd, listenForChannelMessage(channel))
	}
	for i, filePanel := range m.fileModel.filePanels {
		var fileElenent []element
		nowTime := time.Now()
		if filePanel.focusType != noneFocus || nowTime.Sub(filePanel.lastTimeGetElement) > 3*time.Second || foreceReloadElement {
			if len(filePanel.element) < 500 || (len(filePanel.element) > 500 && (nowTime.Sub(filePanel.lastTimeGetElement) > 3*time.Second)) || foreceReloadElement {
				if filePanel.searchBar.Value() != "" {
					fileElenent = returnFolderElementBySearchString(filePanel.location, m.toggleDotFile, filePanel.searchBar.Value())
				} else {
					fileElenent = returnFolderElement(filePanel.location, m.toggleDotFile)
				}
				filePanel.element = fileElenent
				m.fileModel.filePanels[i].element = fileElenent
				m.fileModel.filePanels[i].lastTimeGetElement = nowTime
			}
			foreceReloadElement = false
		}
	}
	return m, cmd
}

func (m model) View() string {
	// check is the terminal size enough
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
		return terminalSizeWarnRender(m)
	} else if m.typingModal.open {
		return typineModalRender(m)
	} else if m.warnModal.open {
		return warnModalRender(m)
	} else {
		sidebar := sidebarRender(m)

		filePanel := filePanelRender(m)

		mainPanel := lipgloss.JoinHorizontal(0, sidebar, filePanel)

		processBar := processBarRender(m)

		metaData := metadataRender(m)

		clipboardBar := clipboardRender(m)

		footer := lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)

		// final render
		finalRender := lipgloss.JoinVertical(0, mainPanel, footer)

		return lipgloss.JoinVertical(lipgloss.Top, finalRender)
	}
}
