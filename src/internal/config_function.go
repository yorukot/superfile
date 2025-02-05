package internal

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

// initialConfig load and handle all configuration files (spf config,hotkeys
// themes) setted up. Returns absolute path of dir pointing to the file Panel
func initialConfig(dir string) (toggleDotFileBool bool, toggleFooter bool, firstFilePanelDir string) {
    // Open log stream
	file, err := os.OpenFile(variable.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	
	// Todo : This could be improved if we want to make superfile more resilient to errors
	// For example if the log file directories have access issues.
	// we could pass a dummy object to log.SetOutput() and the app would still function.
	if err != nil {
		// At this point, it will go to stdout since log file is not initilized
		LogAndExit("Error while opening superfile.log file", "error", err)
	}

	loadConfigFile()

	logLevel := slog.LevelInfo
	if Config.Debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		file, &slog.HandlerOptions{Level: logLevel})))
	
	
	loadHotkeysFile()

	loadThemeFile()

	icon.InitIcon(Config.Nerdfont)

	toggleDotFileData, err := os.ReadFile(variable.ToggleDotFile)
	if err != nil {
		outPutLog("Error while reading toggleDotFile data error:", err)
	}
	if string(toggleDotFileData) == "true" {
		toggleDotFileBool = true
	} else if string(toggleDotFileData) == "false" {
		toggleDotFileBool = false
	} else {
		toggleDotFileBool = false
	}

	toggleFooterData, err := os.ReadFile(variable.ToggleFooter)
	if err != nil {
		outPutLog("Error while reading toggleFooter data error:", err)
	}
	if string(toggleFooterData) == "true" {
		toggleFooter = true
	} else if string(toggleFooterData) == "false" {
		toggleFooter = false
	} else {
		toggleFooter = true
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
		Config.DefaultDirectory = strings.Replace(Config.DefaultDirectory, "~", variable.HomeDir, -1)
		firstFilePanelDir, err = filepath.Abs(Config.DefaultDirectory)
	}

	if err != nil {
		firstFilePanelDir = variable.HomeDir
	}

	slog.Debug("Runtime information", "runtime.GOOS", runtime.GOOS,
		"start directory", firstFilePanelDir)

	return toggleDotFileBool, toggleFooter, firstFilePanelDir
}

// Load configurations from the configuration file. Compares the content
// with the default values and modify the config file to include default configs
// if the FixConfigFile flag is on
func loadConfigFile() {

	//Initialize default configs
	_ = toml.Unmarshal([]byte(ConfigTomlString), &Config)
	//Initialize empty configs
	tempForCheckMissingConfig := ConfigType{}

	data, err := os.ReadFile(variable.ConfigFile)
	if err != nil {
		LogAndExit("Config file doesn't exist", "error", err)
	}

	// Insert data present in the config file inside temp variable
	_ = toml.Unmarshal(data, &tempForCheckMissingConfig)
	// Replace default values for values specifieds in config file
	err = toml.Unmarshal(data, &Config)
	if err != nil && !variable.FixConfigFile {
		fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#F93939")).Render("Error") +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFEE")).Render(" ┃ ") +
			"Error decoding configuration file\n")
		fmt.Println("To add missing fields to hotkeys directory automaticially run Superfile with the --fix-config-file flag `spf --fix-config-file`")
	}

	// If data is different and FixConfigFile option is on, then fullfill then
	// fullfill the config file with the default values
	if !reflect.DeepEqual(Config, tempForCheckMissingConfig) && variable.FixConfigFile {
		tomlData, err := toml.Marshal(Config)
		if err != nil {
			LogAndExit("Error encoding config", "error", err)
		}

		err = os.WriteFile(variable.ConfigFile, tomlData, 0644)
		if err != nil {
			LogAndExit("Error writing config file", "error", err)
		}
	}

	if (Config.FilePreviewWidth > 10 || Config.FilePreviewWidth < 2) && Config.FilePreviewWidth != 0 {
		LogAndExit(loadConfigError("file_preview_width"))
	}

	if Config.SidebarWidth != 0 && (Config.SidebarWidth < 3 || Config.SidebarWidth > 20) {
		LogAndExit(loadConfigError("sidebar_width"))
	}
}

// Load keybinds from the hotkeys file. Compares the content
// with the default values and modify the hotkeys  if the FixHotkeys flag is on.
// If is off check if all hotkeys are properly setted
func loadHotkeysFile() {

    // load default Hotkeys configs
	_ = toml.Unmarshal([]byte(HotkeysTomlString), &hotkeys)
	hotkeysFromConfig := HotkeysType{}
	data, err := os.ReadFile(variable.HotkeysFile)

	if err != nil {
		LogAndExit("Config file doesn't exist", "error", err)
	}
    // Load data from hotkeys file
	_ = toml.Unmarshal(data, &hotkeysFromConfig)
    // Override default hotkeys with the ones from the file
	err = toml.Unmarshal(data, &hotkeys)
	if err != nil {
		LogAndExit("Error decoding hotkeys file ( your config file may have misconfigured", "error", err)
	}

	hasMissingHotkeysInConfig := !reflect.DeepEqual(hotkeys, hotkeysFromConfig)

    // If FixHotKeys is not on then check if every needed hotkey is properly setted
	if hasMissingHotkeysInConfig && !variable.FixHotkeys {
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
		fmt.Println("To add missing fields to hotkeys directory automaticially run Superfile with the --fix-hotkeys flag `spf --fix-hotkeys`")
	}

    // Override hotkey files with default configs if the Fix flag is on
	if hasMissingHotkeysInConfig && variable.FixHotkeys {
		writeHotkeysFile(hotkeys)
	}

	val := reflect.ValueOf(hotkeys)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		if value.Kind() != reflect.Slice || value.Type().Elem().Kind() != reflect.String {
			LogAndExit(loadHotkeysError(field.Name))
		}

		hotkeysList := value.Interface().([]string)

		if len(hotkeysList) == 0 || hotkeysList[0] == "" {
			LogAndExit(loadHotkeysError(field.Name))
		}
	}

}

// Write hotkeys inside the hotkeys toml file
func writeHotkeysFile(hotkeys HotkeysType) {
	tomlData, err := toml.Marshal(hotkeys)
	if err != nil {
		LogAndExit("Error encoding hotkeys", "error", err)
	}

	err = os.WriteFile(variable.HotkeysFile, tomlData, 0644)
	if err != nil {
		LogAndExit("Error writing hotkeys file", "error", err)
	}
}

// Load configurations from theme file into &theme and return default values 
// if file theme folder is empty
func loadThemeFile() {
	themeFile := filepath.Join(variable.ThemeFolder, Config.Theme + ".toml")
	data, err := os.ReadFile(themeFile)
	if err != nil {
		slog.Info("Could not read theme file", "path", themeFile, "error", err)
		data = []byte(DefaultThemeString)
	}

	err = toml.Unmarshal(data, &theme)
	// Todo : Even if user's theme file have errors, lets not exit, but use a default theme file
	if err != nil {
		LogAndExit("Error while decoding theme file( Your theme file may have errors", "error", err)
	}
}

// Load all default configurations from superfile_config folder into global 
// configurations variables
func LoadAllDefaultConfig(content embed.FS) {

	temp, err := content.ReadFile(variable.EmbedHotkeysFile)
	if err != nil {
		outPutLog("Error reading from embed file:", err)
		return
	}
	HotkeysTomlString = string(temp)

	temp, err = content.ReadFile(variable.EmbedConfigFile)
	if err != nil {
		outPutLog("Error reading from embed file:", err)
		return
	}
	ConfigTomlString = string(temp)

	temp, err = content.ReadFile(variable.EmbedThemeCatppuccinFile)
	if err != nil {
		outPutLog("Error reading from embed file:", err)
		return
	}
	DefaultThemeString = string(temp)

	// Todo : We should not return here, and have a default value for this
	currentThemeVersion, err := os.ReadFile(variable.ThemeFileVersion)
	if err != nil && !os.IsNotExist(err) {
		outPutLog("Error reading from file:", err)
		return
	}

	_, err = os.Stat(variable.ThemeFolder)

	if os.IsNotExist(err) {
		err := os.MkdirAll(variable.ThemeFolder, 0755)
		if err != nil {
			outPutLog("error create theme directory", err)
			return
		}
	} else if string(currentThemeVersion) == variable.CurrentVersion {
		return
	}

	files, err := content.ReadDir(variable.EmbedThemeDir)
	if err != nil {
		outPutLog("error read theme directory from embed", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// This will not break in windows. This is a relative path for Embed FS. It uses "/" only 
		src, err := content.ReadFile(variable.EmbedThemeDir + "/" + file.Name())
		if err != nil {
			outPutLog("error read theme file from embed", err)
			return
		}

		file, err := os.Create(filepath.Join(variable.ThemeFolder, file.Name()))
		if err != nil {
			outPutLog("error create theme file from embed", err)
			return
		}
		file.Write(src)
		defer file.Close()
	}

	os.WriteFile(variable.ThemeFileVersion, []byte(variable.CurrentVersion), 0644)
}
