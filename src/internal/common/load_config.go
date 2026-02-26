package common

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/charmbracelet/x/ansi"
	"github.com/pelletier/go-toml/v2"

	"github.com/yorukot/superfile/src/pkg/utils"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

// Load configurations from the configuration file. Compares the content
// with the default values and modify the config file to include default configs
// if the FixConfigFile flag is on
// TODO : Fix the code duplication with LoadHotkeysFile().
func LoadConfigFile() {
	err := utils.LoadTomlFile(variable.ConfigFile, ConfigTomlString, &Config, variable.FixConfigFile, false)
	if err != nil {
		userMsg := fmt.Sprintf("%s%s", LipglossError, err.Error())

		toExit := true
		var loadError *utils.TomlLoadError
		if errors.As(err, &loadError) && loadError != nil {
			if loadError.MissingFields() && !variable.FixConfigFile {
				// Had missing fields and we did not fix
				userMsg += "\nTo add missing fields to configuration file automatically run superfile " +
					"with the --fix-config-file flag `spf --fix-config-file`"
			}
			toExit = loadError.IsFatal()
		}
		if toExit {
			utils.PrintfAndExitf("%s\n", userMsg)
		} else {
			fmt.Println(userMsg)
		}
	}

	// Even if there is a missing field, we want to validate fields that are present
	if err := ValidateConfig(&Config); err != nil {
		// If config is incorrect we cannot continue. We need to exit
		utils.PrintlnAndExit(err.Error())
	}
}

func ValidateConfig(c *ConfigType) error {
	if (c.FilePreviewWidth > 10 || c.FilePreviewWidth < 2) && c.FilePreviewWidth != 0 {
		return errors.New(
			LoadConfigError("file_preview_width", "File preview width must be 2–10, or 0 to disable preview."),
		)
	}

	if c.SidebarWidth != 0 && (c.SidebarWidth < 5 || c.SidebarWidth > 20) {
		return errors.New(LoadConfigError("sidebar_width", "Sidebar width must be 5–20, or 0 to hide the sidebar."))
	}

	for _, order := range c.SidebarSections {
		if order != utils.SidebarSectionHome &&
			order != utils.SidebarSectionPinned &&
			order != utils.SidebarSectionDisks {
			return errors.New(
				LoadConfigError(
					"sidebar_sections",
					"Sidebar sections contain an unsupported value. Allowed values are: home, pinned, disks.",
				),
			)
		}
	}

	if c.DefaultSortType < 0 || c.DefaultSortType > 4 {
		return errors.New(LoadConfigError("default_sort_type", "Default sort type must be between 0 and 4."))
	}

	if c.FilePanelNamePercent < FileNameRatioMin || c.FilePanelNamePercent > FileNameRatioMax {
		return errors.New(
			LoadConfigError("file_panel_name_percent", "File panel name percent is outside the supported range."),
		)
	}

	if ansi.StringWidth(c.BorderTop) != 1 {
		return errors.New(LoadConfigError("border_top", "Border character must be exactly one cell wide."))
	}

	return validateBorders(c)
}

func validateBorders(c *ConfigType) error {
	if ansi.StringWidth(c.BorderBottom) != 1 {
		return errors.New(LoadConfigError("border_bottom", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderLeft) != 1 {
		return errors.New(LoadConfigError("border_left", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderRight) != 1 {
		return errors.New(LoadConfigError("border_right", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderBottomLeft) != 1 {
		return errors.New(LoadConfigError("border_bottom_left", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderBottomRight) != 1 {
		return errors.New(LoadConfigError("border_bottom_right", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderTopLeft) != 1 {
		return errors.New(LoadConfigError("border_top_left", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderTopRight) != 1 {
		return errors.New(LoadConfigError("border_top_right", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderMiddleLeft) != 1 {
		return errors.New(LoadConfigError("border_middle_left", "Border character must be exactly one cell wide."))
	}
	if ansi.StringWidth(c.BorderMiddleRight) != 1 {
		return errors.New(LoadConfigError("border_middle_right", "Border character must be exactly one cell wide."))
	}

	return nil
}

// Load keybinds from the hotkeys file. Compares the content
// with the default values and modify the hotkeys if the FixHotkeys flag is on.
func LoadHotkeysFile(ignoreMissingFields bool) {
	err := utils.LoadTomlFile(
		variable.HotkeysFile,
		HotkeysTomlString,
		&Hotkeys,
		variable.FixHotkeys,
		ignoreMissingFields,
	)
	if err != nil {
		userMsg := fmt.Sprintf("%s%s", LipglossError, err.Error())

		toExit := true
		var loadError *utils.TomlLoadError
		if errors.As(err, &loadError) {
			if loadError.MissingFields() && !variable.FixHotkeys {
				// Had missing fields and we did not fix
				userMsg += "\nTo add missing fields to hotkeys file automatically run superfile " +
					"with the --fix-hotkeys flag `spf --fix-hotkeys`"
			}
			toExit = loadError.IsFatal()
		}
		if toExit {
			utils.PrintfAndExitf("%s\n", userMsg)
		} else {
			fmt.Println(userMsg)
		}
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
			utils.PrintlnAndExit(
				LoadHotkeysError(
					field.Name,
					"Hotkey value must be a list of strings.",
				),
			)
		}

		hotkeysList, ok := value.Interface().([]string)
		if !ok || len(hotkeysList) == 0 || hotkeysList[0] == "" {
			utils.PrintlnAndExit(
				LoadHotkeysError(
					field.Name,
					"Hotkey list is empty; at least one key binding is required.",
				),
			)
		}
	}
}

// LoadThemeFile : Load configurations from theme file into &theme
// set default values if we cant read user's theme file
func LoadThemeFile() {
	themeFile := filepath.Join(variable.ThemeFolder, Config.Theme+".toml")
	if err := LoadUserTheme(themeFile, &Theme); err != nil {
		slog.Error("Could not read user's theme file. Falling back to default theme", "error", err)
		err = toml.Unmarshal([]byte(DefaultThemeString), &Theme)
		if err != nil {
			utils.PrintfAndExitf("Unexpected error while reading default theme file : %v. Exiting...", err)
		}
	}

	// Validations
	if len(Theme.GradientColor) != RequiredGradientColorCount {
		utils.PrintlnAndExit(
			LoadThemeError(
				"gradient_color",
				"Gradient color must contain exactly two values.",
			),
		)
	}
}

func LoadUserTheme(themeFile string, obj *ThemeType) error {
	data, err := os.ReadFile(themeFile)
	if err != nil {
		return fmt.Errorf("could not read user's theme file(%s), err : %w", themeFile, err)
	}
	if err = toml.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("could not unmarshal user's theme file(%s) : %w", themeFile, err)
	}
	return nil
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
	if err = os.MkdirAll(filepath.Dir(variable.ThemeFileVersion), utils.ConfigDirPerm); err != nil {
		slog.Error("Error creating theme file parent directory", "error", err)
		return
	}

	err = os.WriteFile(variable.ThemeFileVersion, []byte(variable.CurrentVersion), utils.ConfigFilePerm)
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
		if err = os.MkdirAll(variable.ThemeFolder, utils.ConfigDirPerm); err != nil {
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
func PopulateGlobalConfigs() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("failed to determine source file location")
	}

	// This is src/internal/common/load_config.go
	// we want src/superfile_config
	spfConfigDir := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filename))),
		"superfile_config")

	configFilePath := filepath.Join(spfConfigDir, "config.toml")
	hotkeyFilePath := filepath.Join(spfConfigDir, "hotkeys.toml")
	themeFilePath := filepath.Join(spfConfigDir, "theme", "monokai.toml")

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

	// Populate fixed variables
	LoadInitialPrerenderedVariables()
	icon.InitIcon(Config.Nerdfont, Theme.DirectoryIconColor)
	LoadPrerenderedVariables()
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

func InitTrash() bool {
	// Create trash directories
	if runtime.GOOS != utils.OsLinux {
		return true
	}
	err := utils.CreateDirectories(
		variable.LinuxTrashDirectory,
		variable.LinuxTrashDirectoryFiles,
		variable.LinuxTrashDirectoryInfo,
	)
	if err != nil {
		slog.Warn("Failed to initialize XDG trash; falling back to permanent delete",
			"error", err, "trashDir", variable.LinuxTrashDirectory)
		return false
	}
	return true
}
