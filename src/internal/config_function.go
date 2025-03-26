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

	icon.InitIcon(Config.Nerdfont, theme.DirectoryIconColor)

	toggleDotFileData, err := os.ReadFile(variable.ToggleDotFile)
	if err != nil {
		slog.Error("Error while reading toggleDotFile data error:", "error", err)
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
		slog.Error("Error while reading toggleFooter data error:", "error", err)
	}
	if string(toggleFooterData) == "true" {
		toggleFooter = true
	} else if string(toggleFooterData) == "false" {
		toggleFooter = false
	} else {
		toggleFooter = true
	}

	LoadThemeConfig()
	LoadPrerenderedVariables()

	if Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			slog.Error("Error while initial model function init exiftool error", "error", err)
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

func writeTomlData(filePath string, data interface{}) error {
	tomlData, err := toml.Marshal(data)
	if err != nil {
		// return a wrapped error
		return fmt.Errorf("Error encoding data : %w", err)
	}
	err = os.WriteFile(filePath, tomlData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing file : %w", err)
	}
	return nil
}

// Helper function to load and validate TOML files with field checking
func loadTomlFile(filePath string, defaultData string, target interface{}, fixFlag bool) (hasError bool) {
	// Initialize with default config
	_ = toml.Unmarshal([]byte(defaultData), target)

	data, err := os.ReadFile(filePath)
	if err != nil {
		LogAndExit("Config file doesn't exist", "error", err)
	}
	errMsg := ""

	// Create a map to track which fields are present
	// Have to do this manually as toml.Unmarshal does not return an error when it encounters a TOML key
	// that does not match any field in the struct.
	// Instead, it simply ignores that key and continues parsing.
	var rawData map[string]interface{}
	rawError := toml.Unmarshal(data, &rawData)
	if rawError != nil {
		hasError = true
		errMsg = lipglossError + "Error decoding file: " + rawError.Error() + "\n"
	}

	// Replace default values with file values
	if !hasError {
		if err = toml.Unmarshal(data, target); err != nil {
			hasError = true
			if decodeErr, ok := err.(*toml.DecodeError); ok {
				row, col := decodeErr.Position()
				errMsg = lipglossError + fmt.Sprintf("Error in field at line %d column %d: %s\n",
					row, col, decodeErr.Error())
			} else {
				errMsg = lipglossError + "Error unmarshalling data: " + err.Error() + "\n"
			}
		}
	}

	if !hasError {
		// Check for missing fields if no errors yet
		targetType := reflect.TypeOf(target).Elem()

		for i := 0; i < targetType.NumField(); i++ {
			field := targetType.Field(i)
			if _, exists := rawData[field.Tag.Get("toml")]; !exists {
				hasError = true
				// A field doesn't exist in the toml config file
				errMsg += lipglossError + fmt.Sprintf("Field \"%s\" is missing\n", field.Tag.Get("toml"))
			}
		}
	}
	// File is okay
	if !hasError {
		return false
	}

	// File is bad, but we arent' allowed to fix
	// We just print error message to stdout
	// Todo : Ideally this should behave as an intenal function with no side effects
	// and the caller should do the printing to stdout
	if !fixFlag {
		fmt.Print(errMsg)
		return true
	}
	// Now we are fixing the file, we would not return hasError=true even if there was error
	// Fix the file by writing all fields
	if err := writeTomlData(filePath, target); err != nil {
		LogAndExit("Error while writing config file", "error", err)
	}
	return false
}

// Load configurations from the configuration file. Compares the content
// with the default values and modify the config file to include default configs
// if the FixConfigFile flag is on
func loadConfigFile() {
	hasError := loadTomlFile(variable.ConfigFile, ConfigTomlString, &Config, variable.FixConfigFile)
	if hasError {
		fmt.Println("To add missing fields to configuration file automatically run superfile with the --fix-config-file flag `spf --fix-config-file`")
		return
	}

	if (Config.FilePreviewWidth > 10 || Config.FilePreviewWidth < 2) && Config.FilePreviewWidth != 0 {
		LogAndExit(loadConfigError("file_preview_width"))
	}

	if Config.SidebarWidth != 0 && (Config.SidebarWidth < 3 || Config.SidebarWidth > 20) {
		LogAndExit(loadConfigError("sidebar_width"))
	}

}

// Load keybinds from the hotkeys file. Compares the content
// with the default values and modify the hotkeys if the FixHotkeys flag is on.
func loadHotkeysFile() {
	hasError := loadTomlFile(variable.HotkeysFile, HotkeysTomlString, &hotkeys, variable.FixHotkeys)
	if hasError {
		fmt.Println("To add missing fields to hotkeys file automatically run superfile with the --fix-hotkeys flag `spf --fix-hotkeys`")
		return
	}

	// Validate hotkey values
	val := reflect.ValueOf(hotkeys)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		if value.Kind() != reflect.Slice || value.Type().Elem().Kind() != reflect.String {
			LogAndExit(loadHotkeysError(field.Name))
		}

		hotkeysList, ok := value.Interface().([]string)
		if !ok || len(hotkeysList) == 0 || hotkeysList[0] == "" {
			LogAndExit(loadHotkeysError(field.Name))
		}
	}
}

// Load configurations from theme file into &theme and return default values
// if file theme folder is empty
func loadThemeFile() {
	themeFile := filepath.Join(variable.ThemeFolder, Config.Theme+".toml")
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
		slog.Error("Error reading from embed file:", "error", err)
		return
	}
	HotkeysTomlString = string(temp)

	temp, err = content.ReadFile(variable.EmbedConfigFile)
	if err != nil {
		slog.Error("Error reading from embed file:", "error", err)
		return
	}
	ConfigTomlString = string(temp)

	temp, err = content.ReadFile(variable.EmbedThemeCatppuccinFile)
	if err != nil {
		slog.Error("Error reading from embed file:", "error", err)
		return
	}
	DefaultThemeString = string(temp)

	// Todo : We should not return here, and have a default value for this
	currentThemeVersion, err := os.ReadFile(variable.ThemeFileVersion)
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Error reading from file:", "error", err)
		return
	}

	_, err = os.Stat(variable.ThemeFolder)

	if os.IsNotExist(err) {
		err := os.MkdirAll(variable.ThemeFolder, 0755)
		if err != nil {
			slog.Error("Error creating theme directory", "error", err)
			return
		}
	} else if string(currentThemeVersion) == variable.CurrentVersion {
		return
	}

	files, err := content.ReadDir(variable.EmbedThemeDir)
	if err != nil {
		slog.Error("Error reading theme directory from embed", "error", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// This will not break in windows. This is a relative path for Embed FS. It uses "/" only
		src, err := content.ReadFile(variable.EmbedThemeDir + "/" + file.Name())
		if err != nil {
			slog.Error("Error reading theme file from embed", "error", err)
			return
		}

		file, err := os.Create(filepath.Join(variable.ThemeFolder, file.Name()))
		if err != nil {
			slog.Error("Error creating theme file from embed", "error", err)
			return
		}
		defer file.Close()
		_, err = file.Write(src)
		if err != nil {
			slog.Error("Error writing theme file from embed", "error", err)
			return
		}
	}

	// Prevent failure for first time app run by making sure parent directories exists
	if err = os.MkdirAll(filepath.Dir(variable.ThemeFileVersion), 0755); err != nil {
		slog.Error("Error creating theme file parent directory", "error", err)
		return
	}

	err = os.WriteFile(variable.ThemeFileVersion, []byte(variable.CurrentVersion), 0644)
	if err != nil {
		slog.Error("Error writing theme file version", "error", err)
	}
}
