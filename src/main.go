package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	components "github.com/MHNightCat/superfile/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var HomeDir = getHomeDir()
var SuperFileMainDir = HomeDir + "/.config/superfile"

const (
	currentVersion      string = "v1.0.1"
	latestVersionURL    string = "https://api.github.com/repos/MHNightCat/superfile/releases/latest"
	latestVersionGithub string = "github.com/MHNightCat/superfile/releases/latest"
	themeZip            string = "https://github.com/MHNightCat/superfile/raw/main/theme.zip"
)

const (
	configFolder     string = "/config"
	themeFolder      string = "/theme"
	trashFolder      string = "/trash"
	dataFolder       string = "/data"
	lastCheckVersion string = "/data/lastCheckVersion"
	pinnedFile       string = "/data/pinned.json"
	configFile       string = "/config/config.json"
	themeZipName     string = "/theme.zip"
	logFile          string = "/superfile.log"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func main() {
	app := &cli.App{
		Name:        "superfile",
		Version:     currentVersion,
		Description: "A Modern file manager with golang",
		ArgsUsage:   "[path]",
		Action: func(c *cli.Context) error {
			path := ""
			if c.Args().Present() {
				path = c.Args().First()
			}

			InitConfigFile()

			p := tea.NewProgram(components.InitialModel(path), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
			CheckForUpdates()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getHomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("can't get home dir")
	}
	return user.HomeDir
}

func InitConfigFile() {
	var err error

	// create main folder
	err = CreateFolderIfNotExist(SuperFileMainDir)
	if err != nil {
		log.Fatalln("Can't Create Superfile main config folder:", SuperFileMainDir, err)
	}
	// create data folder
	err = CreateFolderIfNotExist(SuperFileMainDir + dataFolder)
	if err != nil {
		log.Fatalln("Can't Create Superfile data folder:", SuperFileMainDir+dataFolder, err)
	}
	// create config folder
	err = CreateFolderIfNotExist(SuperFileMainDir + configFolder)
	if err != nil {
		log.Fatalln("Can't Create Superfile data folder:", SuperFileMainDir+configFolder, err)
	}
	// create pinned.json file
	err = CreateFileIfNotExist(SuperFileMainDir + pinnedFile)
	if err != nil {
		log.Fatalln("Can't Create Superfile pinned file:", SuperFileMainDir+pinnedFile, err)
	}
	// create superfile.log file
	err = CreateFileIfNotExist(SuperFileMainDir + logFile)
	if err != nil {
		log.Fatalln("Can't Create Superfile log file:", SuperFileMainDir+logFile, err)
	}
	// write config.json file
	if _, err := os.Stat(SuperFileMainDir + configFile); os.IsNotExist(err) {
		configJsonByte := []byte(configJsonString)
		err = os.WriteFile(SuperFileMainDir+configFile, configJsonByte, 0644)
		if err != nil {
			log.Fatalln("Can't Create Or Write Download Superfile config file:", SuperFileMainDir+configFile, err)
		}
	}
	if err != nil {
		log.Fatalln("Can't Create Or Write Download Superfile config file:", SuperFileMainDir+configFile, err)
	}
	// download theme and unzip it
	if _, err := os.Stat(SuperFileMainDir + themeFolder); os.IsNotExist(err) {
		fmt.Println("First initialize the superfile configuration.\nNeed to download the theme.\nPlease make sure you have internet access.\nAnd this may take some time(<10s).")
		err := DownloadFile(SuperFileMainDir+themeZipName, themeZip)
		if err != nil {
			return
		}

		err = Unzip(SuperFileMainDir+themeZipName, SuperFileMainDir)
		if err != nil {
			return
		}

		err = os.Remove(SuperFileMainDir + themeZipName)
		if err != nil {
			return
		}

	}

	if err != nil {
		log.Fatalln("Can't Auto Download Superfile theme folder:", SuperFileMainDir+themeFolder, err)
		return
	}

}

func CreateFolderIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateFileIfNotExist(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func CheckForUpdates() {
	lastTime, err := ReadFromFile(SuperFileMainDir + lastCheckVersion)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading from file:", err)
		return
	}

	currentTime := time.Now()

	if lastTime.IsZero() || currentTime.Sub(lastTime) >= 24*time.Hour {
		resp, err := http.Get(latestVersionURL)
		if err != nil {
			fmt.Println("Error checking for updates:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var release GitHubRelease
		if err := json.Unmarshal(body, &release); err != nil {
			return
		}
		if release.TagName != currentVersion {
			fmt.Printf("A new version %s is available.\n", release.TagName)
			fmt.Printf("Please update.\n┏\n\n        %s\n\n", latestVersionGithub)
			fmt.Printf("                                                               ┛\n")
		}

		timeStr := currentTime.Format(time.RFC3339)
		err = WriteToFile(SuperFileMainDir+lastCheckVersion, timeStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func ReadFromFile(filename string) (time.Time, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return time.Time{}, err
	}
	if len(content) == 0 {
		return time.Time{}, nil
	}
	lastTime, err := time.Parse(time.RFC3339, string(content))
	if err != nil {
		return time.Time{}, err
	}

	return lastTime, nil
}

func WriteToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const configJsonString string = `{
	"theme": "gruvbox",
	"terminal": "",
	"terminalWorkDirFlag": "",
  
	"_COMMIT_HOTKEY": "",
  
	"_COMMIT_global_hotkey": "Here is global, all global key cant conflicts with other hotkeys",
	"reload": ["ctrl+r", ""],
	"quit": ["esc", "q"],
  
	"listUp": ["up", "k"],
	"listDown": ["down", "j"],
  
	"openTerminal": ["ctrl+t", ""],
  
	"pinnedFolder": ["ctrl+p", ""],
  
	"closeFilePanel": ["ctrl+w", ""],
	"createNewFilePanel": ["ctrl+n", ""],
  
	"nextFilePanel": ["tab", ""],
	"previousFilePanel": ["shift+left", ""],
	"focusOnProcessBar": ["p", ""],
	"focusOnSideBar": ["b", ""],
	"focusOnMetaData": ["m", ""],
  
	"changePanelMode": ["v", ""],
  
	"filePanelFolderCreate": ["f", ""],
	"filePanelFileCreate": ["c", ""],
	"filePanelItemRename": ["r", ""],
	"pasteItem": ["ctrl+v", ""],
  
	"_COMMIT_special_hotkey": "These hotkeys do not conflict with any other keys (including global hotkey)",
	"cancel": ["ctrl+c", "esc"],
	"confirm": ["enter", ""],
  
	"_COMMIT_normal_mode_hotkey": "Here is normal mode hotkey you can conflicts with other mode (cant conflicts with global hotkey)",
	"deleteItem": ["ctrl+d", ""],
	"selectItem": ["enter", "l"],
	"parentFolder": ["h", "backspace"],
	"copySingleItem": ["ctrl+c", ""],
	"cutSingleItem": ["ctrl+x", ""],
  
	"_COMMIT_select_mode_hotkey": "Here is select mode hotkey you can conflicts with other mode (cant conflicts with global hotkey)",
	"filePanelSelectModeItemSingleSelect": ["enter", "l"],
	"filePanelSelectModeItemSelectDown": ["shift+down", "J"],
	"filePanelSelectModeItemSelectUp": ["shift+up", "K"],
	"filePanelSelectModeItemDelete": ["ctrl+d", "delete"],
	"filePanelSelectModeItemCopy": ["ctrl+c", ""],
	"filePanelSelectModeItemCut": ["ctrl+x", ""],
	"filePanelSelectAllItem": ["ctrl+a", ""],
  
	"_COMMIT_process_bar_hotkey": "Here is process bar panel hotkey you can conflicts with other mode (cant conflicts global hotkey)"
  }`
