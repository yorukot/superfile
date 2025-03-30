package common

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pelletier/go-toml/v2"
	"github.com/yorukot/superfile/src/internal/common/utils"

	variable "github.com/yorukot/superfile/src/config"
)

// Load configurations from the configuration file. Compares the content
// with the default values and modify the config file to include default configs
// if the FixConfigFile flag is on
func LoadConfigFile() {
	hasError := utils.LoadTomlFile(variable.ConfigFile, ConfigTomlString, &Config, variable.FixConfigFile, LipglossError)
	if hasError {
		fmt.Println("To add missing fields to configuration file automatically run superfile with the --fix-config-file flag `spf --fix-config-file`")
		return
	}

	if (Config.FilePreviewWidth > 10 || Config.FilePreviewWidth < 2) && Config.FilePreviewWidth != 0 {
		utils.LogAndExit(LoadConfigError("file_preview_width"))
	}

	if Config.SidebarWidth != 0 && (Config.SidebarWidth < 3 || Config.SidebarWidth > 20) {
		utils.LogAndExit(LoadConfigError("sidebar_width"))
	}

}

// Load keybinds from the hotkeys file. Compares the content
// with the default values and modify the hotkeys if the FixHotkeys flag is on.
func LoadHotkeysFile() {
	hasError := utils.LoadTomlFile(variable.HotkeysFile, HotkeysTomlString, &Hotkeys, variable.FixHotkeys, LipglossError)
	if hasError {
		fmt.Println("To add missing fields to hotkeys file automatically run superfile with the --fix-hotkeys flag `spf --fix-hotkeys`")
		return
	}

	// Validate hotkey values
	val := reflect.ValueOf(Hotkeys)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		if value.Kind() != reflect.Slice || value.Type().Elem().Kind() != reflect.String {
			utils.LogAndExit(LoadHotkeysError(field.Name))
		}

		hotkeysList, ok := value.Interface().([]string)
		if !ok || len(hotkeysList) == 0 || hotkeysList[0] == "" {
			utils.LogAndExit(LoadHotkeysError(field.Name))
		}
	}
}

// Load configurations from theme file into &theme and return default values
// if file theme folder is empty
func LoadThemeFile() {
	themeFile := filepath.Join(variable.ThemeFolder, Config.Theme+".toml")
	data, err := os.ReadFile(themeFile)
	if err != nil {
		slog.Info("Could not read theme file", "path", themeFile, "error", err)
		data = []byte(DefaultThemeString)
	}

	err = toml.Unmarshal(data, &Theme)
	// Todo : Even if user's theme file have errors, lets not exit, but use a default theme file
	if err != nil {
		utils.LogAndExit("Error while decoding theme file( Your theme file may have errors", "error", err)
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
		if err = os.MkdirAll(variable.ThemeFolder, 0755); err != nil {
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
		// nolint:govet // Suppress err shadowing
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
