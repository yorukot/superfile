package components

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	configFolder     string = "/config"
	themeFolder      string = "/theme"
	dataFolder       string = "/data"
	lastCheckVersion string = "/data/lastCheckVersion"
	pinnedFile       string = "/data/pinned.json"
	logFile          string = "/data/superfile.log"
	configFile       string = "/config/config.json"
	themeZipName     string = "/theme.zip"
)

var HomeDir = getHomeDir()
var SuperFileMainDir = HomeDir + "/.config/superfile"

var theme ThemeType
var Config ConfigType

var logOutput *os.File
var et *exiftool.Exiftool

var processBarChannel = make(chan processBarMessage, 1000)

func InitialModel(dir string) model {
	var err error
	logOutput, err = os.OpenFile(SuperFileMainDir+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening superfile.log file: %v", err)
	}

	data, err := os.ReadFile(SuperFileMainDir + configFile)
	if err != nil {
		log.Fatalf("config file not exist: %v", err)
	}

	err = json.Unmarshal(data, &Config)

	if err != nil {
		log.Fatalf("Error decoding config json(your config  file may be errors): %v", err)
	}

	data, err = os.ReadFile(SuperFileMainDir + themeFolder + "/" + Config.Theme + ".json")

	if err != nil {
		log.Fatalf("Theme file not exist: %v", err)
	}

	err = json.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error decoding theme json(your theme file may be errors): %v", err)
	}

	LoadThemeConfig()
	et, err = exiftool.NewExiftool()
	if err != nil {
		OutPutLog("Initia model function init exiftool error", err)
	}
	firstFilePanelDir := HomeDir
	if dir != "" {
		firstFilePanelDir, err = filepath.Abs(dir)
		if err != nil {
			firstFilePanelDir = HomeDir
		}
	}
	return model{
		filePanelFocusIndex: 0,
		focusPanel:          nonePanelFocus,
		processBarModel: processBarModel{
			process: make(map[string]process),
			cursor:  0,
			render:  0,
		},
		sideBarModel: sideBarModel{
			pinnedModel: pinnedModel{
				folder: getFolder(),
			},
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{
					render:       0,
					cursor:       0,
					location:     firstFilePanelDir,
					panelMode:    browserMode,
					focusType:    focus,
					folderRecord: make(map[string]folderRecord),
				},
			},
			width: 10,
		},
	}
}

func listenForProcessBarMessage(msg chan processBarMessage) tea.Cmd {
	return func() tea.Msg {
		select {
		case m := <-msg:
			return m // 返回从通道中接收到的消息
		default:
			return nil // 没有消息可用时返回 nil
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("SuperFile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
		listenForProcessBarMessage(processBarChannel),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case processBarMessage:
		if !contains(m.processBarModel.processList, msg.processId) {
			m.processBarModel.processList = append(m.processBarModel.processList, msg.processId)
		}
		m.processBarModel.process[msg.processId] = msg.processNewState
	case tea.WindowSizeMsg:
		m.mainPanelHeight = msg.Height - bottomBarHeight
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		return m, nil
	case tea.KeyMsg:
		// if in the create item modal
		if m.createNewItem.open {
			switch msg.String() {
			case Config.Cancel[0], Config.Cancel[1]:
				m = CancelModal(m)
			case Config.Confirm[0], Config.Confirm[1]:
				m = CreateItem(m)
			}
			// if in the renaming mode
		} else if m.fileModel.renaming {
			switch msg.String() {
			case Config.Cancel[0], Config.Cancel[1]:
				m = CancelReanem(m)
			case Config.Confirm[0], Config.Confirm[1]:
				m = ConfirmRename(m)
			}
		} else {
			switch msg.String() {
			// return superfile
			case Config.Reload[0], Config.Reload[1]:
				//just do nothing
			case Config.Quit[0], Config.Quit[1]:
				return m, tea.Quit
			/* LIST CONTROLLER START */
			// up list
			case Config.ListUp[0], Config.ListUp[1]:
				if m.focusPanel == sideBarFocus {
					m = ControllerSideBarListUp(m)
				} else if m.focusPanel == processBarFocus {
					m = ContollerProcessBarListUp(m)
				} else if m.focusPanel == metaDataFocus {
					m = ControllerMetaDataListUp(m)
				} else if m.focusPanel == nonePanelFocus {
					m = ControllerFilePanelListUp(m)
					m = ReturnMetaData(m)
				}
			// down list
			case Config.ListDown[0], Config.ListDown[1]:
				if m.focusPanel == sideBarFocus {
					m = ControllerSideBarListDown(m)
				} else if m.focusPanel == processBarFocus {
					m = ContollerProcessBarListDown(m)
				} else if m.focusPanel == metaDataFocus {
					m = ControllerMetaDataListDown(m)
				} else if m.focusPanel == nonePanelFocus {
					m = ControllerFilePanelListDown(m)
					m = ReturnMetaData(m)
				}
			/* LIST CONTROLLER END */
			case Config.ChangePanelMode[0], Config.ChangePanelMode[1]:
				m = SelectedMode(m)
			/* NAVIGATION CONTROLLER START */
			// change file panel
			case Config.NextFilePanel[0], Config.NextFilePanel[1]:
				m = NextFilePanel(m)
			// change file panel
			case Config.PreviousFilePanel[0], Config.PreviousFilePanel[1]:
				m = PreviousFilePanel(m)
			// close file panel
			case Config.CloseFilePanel[0], Config.CloseFilePanel[1]:
				m = CloseFilePanel(m)
			// create new file panel
			case Config.CreateNewFilePanel[0], Config.CreateNewFilePanel[1]:
				m = CreateNewFilePanel(m)
			// focus to sidebar or file panel
			case Config.FocusOnSideBar[0], Config.FocusOnSideBar[1]:
				m = FocusOnSideBar(m)
			/* NAVIGATION CONTROLLER END */
			case Config.FocusOnProcessBar[0], Config.FocusOnProcessBar[1]:
				m = FocusOnProcessBar(m)
			case Config.FocusOnMetaData[0], Config.FocusOnMetaData[1]:
				m = FocusOnMetaData(m)
				m = ReturnMetaData(m)
			case Config.PasteItem[0], Config.PasteItem[1]:
				go func() {
					m = PasteItem(m)
				}()
			case Config.FilePanelFileCreate[0], Config.FilePanelFileCreate[1]:
				m = PanelCreateNewFile(m)
			case Config.FilePanelFolderCreate[0], Config.FilePanelFolderCreate[1]:
				m = PanelCreateNewFolder(m)
			case Config.PinnedFolder[0], Config.PinnedFolder[1]:
				m = PinnedFolder(m)
			case Config.OpenTerminal[0], Config.OpenTerminal[1]:
				m = OpenTerminal(m)
			default:
				// check if it's the select mode
				if m.fileModel.filePanels[m.filePanelFocusIndex].focusType == focus && m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
					switch msg.String() {
					case Config.FilePanelSelectModeItemSingleSelect[0], Config.FilePanelSelectModeItemSingleSelect[1]:
						m = SingleItemSelect(m)
					case Config.FilePanelSelectModeItemSelectUp[0], Config.FilePanelSelectModeItemSelectUp[1]:
						m = ItemSelectUp(m)
					case Config.FilePanelSelectModeItemSelectDown[0], Config.FilePanelSelectModeItemSelectDown[1]:
						m = ItemSelectDown(m)
					case Config.FilePanelSelectModeItemDelete[0], Config.FilePanelSelectModeItemDelete[1]:
						go func() {
							m = DeleteMultipleItem(m)
							m.fileModel.filePanels[m.filePanelFocusIndex].selected = m.fileModel.filePanels[m.filePanelFocusIndex].selected[:0]
						}()
					case Config.FilePanelSelectModeItemCopy[0], Config.FilePanelSelectModeItemCopy[1]:
						m = CopyMultipleItem(m)
					case Config.FilePanelSelectModeItemCut[0], Config.FilePanelSelectModeItemCut[1]:
						m = CutMultipleItem(m)
					case Config.FilePanelSelectAllItem[0], Config.FilePanelSelectAllItem[1]:
						m = SelectAllItem(m)
					}
					// else
				} else {
					switch msg.String() {
					case Config.SelectItem[0], Config.SelectItem[1]:
						if m.focusPanel == sideBarFocus {
							m = SideBarSelectFolder(m)
						} else if m.focusPanel == processBarFocus {

						} else if m.focusPanel == nonePanelFocus {
							m = EnterPanel(m)
						}
					case Config.ParentFolder[0], Config.ParentFolder[1]:
						m = ParentFolder(m)
					case Config.DeleteItem[0], Config.DeleteItem[1]:
						go func() {
							m = DeleteSingleItem(m)
						}()
					case Config.CopySingleItem[0], Config.CopySingleItem[1]:
						m = CopySingleItem(m)
					case Config.CutSingleItem[0], Config.CutSingleItem[1]:
						m = CutSingleItem(m)
					case Config.FilePanelItemRename[0], Config.FilePanelItemRename[1]:
						m = PanelItemRename(m)
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
		m.createNewItem.textInput, cmd = m.createNewItem.textInput.Update(msg)
	}

	if m.fileModel.filePanels[m.filePanelFocusIndex].cursor < 0 {
		m.fileModel.filePanels[m.filePanelFocusIndex].cursor = 0
	}

	cmd = tea.Batch(cmd, listenForProcessBarMessage(processBarChannel))
	m.sideBarModel.pinnedModel.folder = getFolder()
	return m, cmd
}

func (m model) View() string {
	// check is the terminal size enough
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
		return TerminalSizeWarnRender(m)
	} else if m.createNewItem.open {
		return ModalRender(m)
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
