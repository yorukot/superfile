package variable

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/adrg/xdg"
)

const (
	CurrentVersion      = "v1.2.1"
	LatestVersionURL    = "https://api.github.com/repos/yorukot/superfile/releases/latest"
	LatestVersionGithub = "github.com/yorukot/superfile/releases/latest"

	// This will not break in windows. This is a relative path for Embed FS. It uses "/" only
	EmbedConfigDir           = "src/superfile_config"
	EmbedConfigFile          = EmbedConfigDir + "/config.toml"
	EmbedHotkeysFile         = EmbedConfigDir + "/hotkeys.toml"
	EmbedThemeDir            = EmbedConfigDir + "/theme"
	EmbedThemeCatppuccinFile = EmbedThemeDir + "/catppuccin.toml"
)

var (
	HomeDir           = xdg.Home
	SuperFileMainDir  = filepath.Join(xdg.ConfigHome, "superfile")
	SuperFileCacheDir = filepath.Join(xdg.CacheHome, "superfile")
	SuperFileDataDir  = filepath.Join(xdg.DataHome, "superfile")
	SuperFileStateDir = filepath.Join(xdg.StateHome, "superfile")

	// MainDir files
	ThemeFolder = filepath.Join(SuperFileMainDir, "theme")

	// DataDir files
	LastCheckVersion = filepath.Join(SuperFileDataDir, "lastCheckVersion")
	ThemeFileVersion = filepath.Join(SuperFileDataDir, "themeFileVersion")
	FirstUseCheck    = filepath.Join(SuperFileDataDir, "firstUseCheck")
	PinnedFile       = filepath.Join(SuperFileDataDir, "pinned.json")
	ToggleDotFile    = filepath.Join(SuperFileDataDir, "toggleDotFile")
	ToggleFooter     = filepath.Join(SuperFileDataDir, "toggleFooter")

	// StateDir files
	LogFile     = filepath.Join(SuperFileStateDir, "superfile.log")
	LastDirFile = filepath.Join(SuperFileStateDir, "lastdir")

	// Trash Directories
	DarwinTrashDirectory      = filepath.Join(HomeDir, ".Trash")
	CustomTrashDirectory      = filepath.Join(xdg.DataHome, "Trash")
	CustomTrashDirectoryFiles = filepath.Join(xdg.DataHome, "Trash", "files")
	CustomTrashDirectoryInfo  = filepath.Join(xdg.DataHome, "Trash", "info")
)

// These variables are actually not fixed, they are sometimes updated dynamically
var (
	ConfigFile  = filepath.Join(SuperFileMainDir, "config.toml")
	HotkeysFile = filepath.Join(SuperFileMainDir, "hotkeys.toml")

	// Other state variables
	FixHotkeys    = false
	FixConfigFile = false
	LastDir       = ""
	PrintLastDir  = false
)

// Still we are preventing other packages to directly modify them via reassign linter

func SetLastDir(path string) {
	LastDir = path
}

func UpdateVarFromCliArgs(c *cli.Context) {
	// Setting the config file path
	configFileArg := c.String("config-file")

	// Validate the config file exists
	if configFileArg != "" {
		if _, err := os.Stat(configFileArg); err != nil {
			utils.PrintfAndExit("Error: While reading config file '%s' from argument : %v", configFileArg, err)
		}
		ConfigFile = configFileArg
	}

	hotkeyFileArg := c.String("hotkey-file")

	if hotkeyFileArg != "" {
		if _, err := os.Stat(hotkeyFileArg); err != nil {
			utils.PrintfAndExit("Error: While reading hotkey file '%s' from argument : %v", hotkeyFileArg, err)
		}
		HotkeysFile = hotkeyFileArg
	}

	FixHotkeys = c.Bool("fix-hotkeys")
	FixConfigFile = c.Bool("fix-config-file")
	PrintLastDir = c.Bool("print-last-dir")
}
