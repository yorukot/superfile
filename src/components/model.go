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
	configFile       string = "/config.json"
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

var logOutput *os.File
var et *exiftool.Exiftool

var channel = make(chan channelMessage, 1000)

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
		sideBarModel: sideBarModel{
			directories: getDirectories(),
			// wellKnownModel: getWellKnownDirectories(),
			// pinnedModel:    getPinnedDirectories(),
			// disksModel:     getExternalMediaFolders(),
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
		m.mainPanelHeight = msg.Height - bottomBarHeight + 1
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		return m, nil
	// if just user press key
	case tea.KeyMsg:
		// if in the create item modal
		if m.typingModal.open {
			switch msg.String() {
			case Config.Cancel[0], Config.Cancel[1]:
				m = cancelTypingModal(m)
			case Config.Confirm[0], Config.Confirm[1]:
				m = createItem(m)
			}
			// if in the renaming mode
		} else if m.warnModal.open {
			switch msg.String() {
			case Config.Cancel[0], Config.Cancel[1]:
				m = cancelWarnModal(m)
			case Config.Confirm[0], Config.Confirm[1]:
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
			case Config.Cancel[0], Config.Cancel[1]:
				m = cancelReanem(m)
			case Config.Confirm[0], Config.Confirm[1]:
				m = confirmRename(m)
			}
		} else {
			switch msg.String() {
			// return superfile
			case Config.Quit[0], Config.Quit[1]:
				return m, tea.Quit
			/* LIST CONTROLLER START */
			// up list
			case Config.ListUp[0], Config.ListUp[1]:
				if m.focusPanel == sideBarFocus {
					m = controllerSideBarListUp(m)
				} else if m.focusPanel == processBarFocus {
					m = contollerProcessBarListUp(m)
				} else if m.focusPanel == metaDataFocus {
					m = controllerMetaDataListUp(m)
				} else if m.focusPanel == nonePanelFocus {
					m = controllerFilePanelListUp(m)
					m.fileMetaData.renderIndex = 0
					go func() {
						m = returnMetaData(m)
					}()
				}
			// down list
			case Config.ListDown[0], Config.ListDown[1]:
				if m.focusPanel == sideBarFocus {
					m = controllerSideBarListDown(m)
				} else if m.focusPanel == processBarFocus {
					m = contollerProcessBarListDown(m)
				} else if m.focusPanel == metaDataFocus {
					m = controllerMetaDataListDown(m)
				} else if m.focusPanel == nonePanelFocus {
					m = controllerFilePanelListDown(m)
					m.fileMetaData.renderIndex = 0
					go func() {
						m = returnMetaData(m)
					}()
				}
			/* LIST CONTROLLER END */
			case Config.ChangePanelMode[0], Config.ChangePanelMode[1]:
				m = selectedMode(m)
			/* NAVIGATION CONTROLLER START */
			// change file panel
			case Config.NextFilePanel[0], Config.NextFilePanel[1]:
				m = nextFilePanel(m)
			// change file panel
			case Config.PreviousFilePanel[0], Config.PreviousFilePanel[1]:
				m = previousFilePanel(m)
			// close file panel
			case Config.CloseFilePanel[0], Config.CloseFilePanel[1]:
				m = closeFilePanel(m)
			// create new file panel
			case Config.CreateNewFilePanel[0], Config.CreateNewFilePanel[1]:
				m = createNewFilePanel(m)
			// focus to sidebar or file panel
			case Config.FocusOnSideBar[0], Config.FocusOnSideBar[1]:
				m = focusOnSideBar(m)
			/* NAVIGATION CONTROLLER END */
			case Config.FocusOnProcessBar[0], Config.FocusOnProcessBar[1]:
				m = focusOnProcessBar(m)
			case Config.FocusOnMetaData[0], Config.FocusOnMetaData[1]:
				m = focusOnMetaData(m)
				go func() {
					m = returnMetaData(m)
				}()
			case Config.PasteItem[0], Config.PasteItem[1]:
				go func() {
					m = pasteItem(m)
				}()
			case Config.FilePanelFileCreate[0], Config.FilePanelFileCreate[1]:
				m = panelCreateNewFile(m)
			case Config.FilePanelDirectoryCreate[0], Config.FilePanelDirectoryCreate[1]:
				m = panelCreateNewFolder(m)
			case Config.PinnedDirectory[0], Config.PinnedDirectory[1]:
				m = pinnedFolder(m)
			case Config.ToggleDotFile[0], Config.ToggleDotFile[1]:
				m = toggleDotFileController(m)
			case Config.ExtractFile[0], Config.ExtractFile[1]:
				go func() {
					m = extractFile(m)
				}()
			case Config.CompressFile[0], Config.CompressFile[1]:
				go func() {
				 	m = compressFile(m)
				}()
			default:
				// check if it's the select mode
				if m.fileModel.filePanels[m.filePanelFocusIndex].focusType == focus && m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
					switch msg.String() {
					case Config.FilePanelSelectModeItemSingleSelect[0], Config.FilePanelSelectModeItemSingleSelect[1]:
						m = singleItemSelect(m)
					case Config.FilePanelSelectModeItemSelectUp[0], Config.FilePanelSelectModeItemSelectUp[1]:
						m = itemSelectUp(m)
					case Config.FilePanelSelectModeItemSelectDown[0], Config.FilePanelSelectModeItemSelectDown[1]:
						m = itemSelectDown(m)
					case Config.FilePanelSelectModeItemDelete[0], Config.FilePanelSelectModeItemDelete[1]:
						go func() {
							m = deleteMultipleItem(m)
							if !isExternalDiskPath(m.fileModel.filePanels[m.filePanelFocusIndex].location) {
								m.fileModel.filePanels[m.filePanelFocusIndex].selected = m.fileModel.filePanels[m.filePanelFocusIndex].selected[:0]
							}
						}()
					case Config.FilePanelSelectModeItemCopy[0], Config.FilePanelSelectModeItemCopy[1]:
						m = copyMultipleItem(m)
					case Config.FilePanelSelectModeItemCut[0], Config.FilePanelSelectModeItemCut[1]:
						m = cutMultipleItem(m)
					case Config.FilePanelSelectAllItem[0], Config.FilePanelSelectAllItem[1]:
						m = selectAllItem(m)
					}
					// else
				} else {
					switch msg.String() {
					case Config.SelectItem[0], Config.SelectItem[1]:
						if m.focusPanel == sideBarFocus {
							m = sideBarSelectFolder(m)
						} else if m.focusPanel == processBarFocus {

						} else if m.focusPanel == nonePanelFocus {
							m = enterPanel(m)
						}
					case Config.ParentDirectory[0], Config.ParentDirectory[1]:
						m = parentFolder(m)
					case Config.DeleteItem[0], Config.DeleteItem[1]:
						go func() {
							m = deleteSingleItem(m)
						}()
					case Config.CopySingleItem[0], Config.CopySingleItem[1]:
						m = copySingleItem(m)
					case Config.CutSingleItem[0], Config.CutSingleItem[1]:
						m = cutSingleItem(m)
					case Config.FilePanelItemRename[0], Config.FilePanelItemRename[1]:
						m = panelItemRename(m)
					}

				}
			}
		}
	}
	if m.firstTextInput {
		m.firstTextInput = false
	} else if m.fileModel.renaming {
		m.fileModel.filePanels[m.filePanelFocusIndex].rename, cmd = m.fileModel.filePanels[m.filePanelFocusIndex].rename.Update(msg)
	} else {
		m.typingModal.textInput, cmd = m.typingModal.textInput.Update(msg)
	}

	if m.fileModel.filePanels[m.filePanelFocusIndex].cursor < 0 {
		m.fileModel.filePanels[m.filePanelFocusIndex].cursor = 0
	}

	cmd = tea.Batch(cmd)
	m.sideBarModel.directories = getDirectories()

	if !ListeningMessage {
		cmd = tea.Batch(cmd, listenForChannelMessage(channel))
	}
	return m, cmd
}

func (m model) View() string {
	// check is the terminal size enough
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
		return TerminalSizeWarnRender(m)
	} else if m.typingModal.open {
		return TypineModalRender(m)
	} else if m.warnModal.open {
		return WarnModalRender(m)
	} else {
		sideBar := SideBarRender(m)

		filePanel := FilePanelRender(m)

		mainPanel := lipgloss.JoinHorizontal(0, sideBar, filePanel)

		processBar := ProcessBarRender(m)

		metaData := MetaDataRender(m)

		clipboardBar := ClipboardRender(m)

		bottomBar := lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)

		// final render
		finalRender := lipgloss.JoinVertical(0, mainPanel, bottomBar)

		return lipgloss.JoinVertical(lipgloss.Top, finalRender)
	}
}
