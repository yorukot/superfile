package common

import (
	"embed"
	"errors"
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
	}
	// Even if there is a missing field, we want to validate fields that are present

	if err := ValidateConfig(&Config); err != nil {
		// If config is incorrect we cannot continue. We need to exit
		utils.LogAndExit(err.Error())
	}
}

func ValidateConfig(c *ConfigType) error {
	if (c.FilePreviewWidth > 10 || c.FilePreviewWidth < 2) && c.FilePreviewWidth != 0 {
		return errors.New(LoadConfigError("file_preview_width"))
	}

	if c.SidebarWidth != 0 && (c.SidebarWidth < 3 || c.SidebarWidth > 20) {
		return errors.New(LoadConfigError("sidebar_width"))
	}

	if c.DefaultSortType < 0 || c.DefaultSortType > 2 {
		return errors.New(LoadConfigError("default_sort_type"))
	}
	return nil
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
	for i := range val.NumField() {
		field := val.Type().Field(i)
		value := val.Field(i)

		// Although this is redundant as Hotkey is always a slice
		// This adds a layer against accidental struct modifications
		// Makes sure its always be a string slice. It's somewhat like a unit test
		if value.Kind() != reflect.Slice || value.Type().Elem().Kind() != reflect.String {
			utils.LogAndExit(LoadHotkeysError(field.Name))
		}

		hotkeysList, ok := value.Interface().([]string)
		if !ok || len(hotkeysList) == 0 || hotkeysList[0] == "" {
			utils.LogAndExit(LoadHotkeysError(field.Name))
		}
	}
}

// LoadThemeFile : Load configurations from theme file into &theme
// set default values if we cant read user's theme file
func LoadThemeFile() {
	themeFile := filepath.Join(variable.ThemeFolder, Config.Theme+".toml")
	data, err := os.ReadFile(themeFile)
	if err == nil {
		if unmarshalErr := toml.Unmarshal(data, &Theme); unmarshalErr == nil {
			return
		} else {
			slog.Error("Could not unmarshal theme file. Falling back to default theme",
				"unmarshalErr", unmarshalErr)
		}
	} else {
		slog.Error("Could not read user's theme file. Falling back to default theme", "path", themeFile, "error", err)
	}

	err = toml.Unmarshal([]byte(DefaultThemeString), &Theme)
	if err != nil {
		utils.LogAndExit("Unexpected error while reading default theme file. Exiting...", "error", err)
	}
}

// LoadAllDefaultConfig : Load all default configurations from embedded superfile_config folder into global
// configurations variables and write theme files if its needed.
func LoadAllDefaultConfig(content embed.FS) {
	err := LoadConfigStringGlobals(content)
	if err != nil {
		slog.Error("Could not load default config from embed FS", "error", err)
		return
	}

	currentThemeVersion, err := os.ReadFile(variable.ThemeFileVersion)
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Unexpected error reading from file:", "error", err)
		return
	}

	if string(currentThemeVersion) == variable.CurrentVersion {
		// We don't need to update themes as its already up to date
		return
	}

	// Write theme files to theme directory
	err = WriteThemeFiles(content)
	if err != nil {
		slog.Error("Error while writing default theme directories", "error", err)
		return
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

func LoadConfigStringGlobals(content embed.FS) error {
	hotkeyData, err := content.ReadFile(variable.EmbedHotkeysFile)
	if err != nil {
		return err
	}
	HotkeysTomlString = string(hotkeyData)

	configData, err := content.ReadFile(variable.EmbedConfigFile)
	if err != nil {
		return err
	}
	ConfigTomlString = string(configData)

	themeData, err := content.ReadFile(variable.EmbedThemeCatppuccinFile)
	if err != nil {
		return err
	}
	DefaultThemeString = string(themeData)
	return nil
}

func WriteThemeFiles(content embed.FS) error {
	_, err := os.Stat(variable.ThemeFolder)

	if os.IsNotExist(err) {
		if err = os.MkdirAll(variable.ThemeFolder, 0755); err != nil {
			slog.Error("Error creating theme directory", "error", err)
			return err
		}
	}

	files, err := content.ReadDir(variable.EmbedThemeDir)
	if err != nil {
		slog.Error("Error reading theme directory from embed", "error", err)
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// This will not break in windows. This is a relative path for Embed FS. It uses "/" only
		src, err := content.ReadFile(variable.EmbedThemeDir + "/" + file.Name())
		if err != nil {
			slog.Error("Error reading theme file from embed", "error", err)
			return err
		}

		curThemeFile, err := os.Create(filepath.Join(variable.ThemeFolder, file.Name()))
		if err != nil {
			slog.Error("Error creating theme file from embed", "error", err)
			return err
		}
		defer curThemeFile.Close()
		_, err = curThemeFile.Write(src)
		if err != nil {
			slog.Error("Error writing theme file from embed", "error", err)
			return err
		}
	}
	return nil
}

// Used only in unit tests
// Populate config variables based on given file
func PopulateGlobalConfigs(configFilePath string, hotkeyFilePath string, themeFilePath string) error {
	err := PopulateConfigFromFile(configFilePath)
	if err != nil {
		return err
	}
	err = PopulateHotkeyFromFile(hotkeyFilePath)
	if err != nil {
		return err
	}
	err = PopulateThemeFromFile(themeFilePath)
	if err != nil {
		return err
	}
	return nil
}

// No validation required
func populateFromFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(data, target)
	if err != nil {
		return err
	}
	return nil
}

func PopulateConfigFromFile(configFilePath string) error {
	return populateFromFile(configFilePath, &Config)
}

func PopulateHotkeyFromFile(hotkeyFilePath string) error {
	return populateFromFile(hotkeyFilePath, &Hotkeys)
}

func PopulateThemeFromFile(themeFilePath string) error {
	return populateFromFile(themeFilePath, &Theme)
}
