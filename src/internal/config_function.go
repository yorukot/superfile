package internal

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	varibale "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

func initialConfig(dir string) (toggleDotFileBool bool, firstFilePanelDir string) {
	var err error

	logOutput, err = os.OpenFile(varibale.LogFilea, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening superfile.log file: %v", err)
	}

	loadConfigFile()

	loadHotkeysFile()

	loadThemeFile()

	icon.InitIcon(Config.Nerdfont)

	toggleDotFileData, err := os.ReadFile(varibale.ToggleDotFilea)
	if err != nil {
		outPutLog("Error while reading toggleDotFile data error:", err)
	}
	if string(toggleDotFileData) == "true" {
		toggleDotFileBool = true
	} else if string(toggleDotFileData) == "false" {
		toggleDotFileBool = false
	}
	LoadThemeConfig()

	if Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			outPutLog("Initial model function init exiftool error", err)
		}
	}

	if dir != "" {
		firstFilePanelDir, err = filepath.Abs(dir)
	} else {
		Config.DefaultDirectory = strings.Replace(Config.DefaultDirectory, "~", varibale.HomeDir, -1)
		firstFilePanelDir, err = filepath.Abs(Config.DefaultDirectory)
	}

	if err != nil {
		firstFilePanelDir = varibale.HomeDir
	}

	return toggleDotFileBool, firstFilePanelDir
}

func loadConfigFile() {

	_ = toml.Unmarshal([]byte(ConfigTomlString), &Config)
	tempForCheckMissingConfig := ConfigType{}

	data, err := os.ReadFile(varibale.ConfigFilea)
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}

	_ = toml.Unmarshal(data, &tempForCheckMissingConfig)
	err = toml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error decoding config file ( your config file may have misconfigured ): %v", err)
	}

	if !reflect.DeepEqual(Config, tempForCheckMissingConfig) {
		tomlData, err := toml.Marshal(Config)
		if err != nil {
			log.Fatalf("Error encoding config: %v", err)
		}

		err = os.WriteFile(varibale.ConfigFilea, tomlData, 0644)
		if err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
	}
	if (Config.FilePreviewWidth > 10 || Config.FilePreviewWidth < 2) && Config.FilePreviewWidth != 0 {
		fmt.Println(loadConfigError("file_preview_width"))
		os.Exit(0)
	}

	if Config.SidebarWidth != 0 && (Config.SidebarWidth < 3 || Config.SidebarWidth > 20) {
		fmt.Println(loadConfigError("sidebar_width"))
		os.Exit(0)
	}
}

func loadHotkeysFile() {

	_ = toml.Unmarshal([]byte(HotkeysTomlString), &hotkeys)
	hotkeysFromConfig := HotkeysType{}
	data, err := os.ReadFile(varibale.HotkeysFilea)

	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}
	_ = toml.Unmarshal(data, &hotkeysFromConfig)
	err = toml.Unmarshal(data, &hotkeys)
	if err != nil {
		log.Fatalf("Error decoding hotkeys file ( your config file may have misconfigured ): %v", err)
	}

	hasMissingHotkeysInConfig := !reflect.DeepEqual(hotkeys, hotkeysFromConfig)

	if hasMissingHotkeysInConfig && !varibale.FixHotkeys {
		hotKeysConfig := reflect.ValueOf(hotkeysFromConfig)
		for i := 0; i < hotKeysConfig.NumField(); i++ {
			field := hotKeysConfig.Type().Field(i)
			value := hotKeysConfig.Field(i)
			name := field.Name
			isMissing := value.Len() == 0

			if isMissing {
				fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#F93939")).Render("Error") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFEE")).Render(" ┃ ") + 
				fmt.Sprintf("Field \"%s\" is missing in hotkeys configuration\n", name))
			}
		}
		fmt.Println("To add missing fields to hotkeys directory automaticially run Superfile with the --fix-hotkeys flag")
	}

	if hasMissingHotkeysInConfig && varibale.FixHotkeys {
		writeHotkeysFile(hotkeys)
	}

	val := reflect.ValueOf(hotkeys)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		if value.Kind() != reflect.Slice || value.Type().Elem().Kind() != reflect.String {
			fmt.Println(lodaHotkeysError(field.Name))
			os.Exit(0)
		}

		hotkeysList := value.Interface().([]string)

		if len(hotkeysList) == 0 || hotkeysList[0] == "" {
			fmt.Println(lodaHotkeysError(field.Name))
			os.Exit(0)
		}
	}

}

func writeHotkeysFile(hotkeys HotkeysType) {
	tomlData, err := toml.Marshal(hotkeys)
	if err != nil {
		log.Fatalf("Error encoding hotkeys: %v", err)
	}

	err = os.WriteFile(varibale.HotkeysFilea, tomlData, 0644)
	if err != nil {
		log.Fatalf("Error writing hotkeys file: %v", err)
	}
}

func loadThemeFile() {
	data, err := os.ReadFile(varibale.ThemeFoldera + "/" + Config.Theme + ".toml")
	if err != nil {
		data = []byte(DefaultThemeString)
	}

	err = toml.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error while decoding theme file( Your theme file may have errors ): %v", err)
	}
}

func LoadAllDefaultConfig(content embed.FS) {

	temp, err := content.ReadFile("src/superfile_config/hotkeys.toml")
	if err != nil {
		return
	}
	HotkeysTomlString = string(temp)

	temp, err = content.ReadFile("src/superfile_config/config.toml")
	if err != nil {
		return
	}
	ConfigTomlString = string(temp)

	temp, err = content.ReadFile("src/superfile_config/theme/catppuccin.toml")
	if err != nil {
		return
	}
	DefaultThemeString = string(temp)

	currentThemeVersion, err := os.ReadFile(varibale.ThemeFileVersiona)

	if err != nil && !os.IsNotExist(err) {
		outPutLog("Error reading from file:", err)
		return
	}

	_, err = os.Stat(varibale.ThemeFoldera)

	if os.IsNotExist(err) {
		err := os.MkdirAll(varibale.ThemeFoldera, 0755)
		if err != nil {
			outPutLog("error create theme direcroty", err)
			return
		}
	} else if string(currentThemeVersion) == varibale.CurrentVersion {
		return
	}

	files, err := content.ReadDir("src/superfile_config/theme")
	if err != nil {
		outPutLog("error read theme directory from embed", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		src, err := content.ReadFile(filepath.Join("src/superfile_config/theme", file.Name()))
		if err != nil {
			outPutLog("error read theme file from embed", err)
			return
		}

		file, err := os.Create(filepath.Join(varibale.ThemeFoldera, file.Name()))
		if err != nil {
			outPutLog("error create theme file from embed", err)
			return
		}
		file.Write(src)
		defer file.Close()
	}

	os.WriteFile(varibale.ThemeFileVersiona, []byte(varibale.CurrentVersion), 0644)
}
