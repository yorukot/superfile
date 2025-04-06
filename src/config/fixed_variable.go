package variable

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/adrg/xdg"
)

const (
	CurrentVersion      string = "v1.2.1"
	LatestVersionURL    string = "https://api.github.com/repos/yorukot/superfile/releases/latest"
	LatestVersionGithub string = "github.com/yorukot/superfile/releases/latest"

	// This will not break in windows. This is a relative path for Embed FS. It uses "/" only
	EmbedConfigDir           string = "src/superfile_config"
	EmbedConfigFile          string = EmbedConfigDir + "/config.toml"
	EmbedHotkeysFile         string = EmbedConfigDir + "/hotkeys.toml"
	EmbedThemeDir            string = EmbedConfigDir + "/theme"
	EmbedThemeCatppuccinFile string = EmbedThemeDir + "/catppuccin.toml"

	// These are used while comparing with runtime.GOOS
	// OS_WINDOWS represents the Windows operating system identifier
	OS_WINDOWS = "windows"
	// OS_DARWIN represents the macOS (Darwin) operating system identifier
	OS_DARWIN = "darwin"
)

var (
	HomeDir           string = xdg.Home
	SuperFileMainDir  string = filepath.Join(xdg.ConfigHome, "superfile")
	SuperFileCacheDir string = filepath.Join(xdg.CacheHome, "superfile")
	SuperFileDataDir  string = filepath.Join(xdg.DataHome, "superfile")
	SuperFileStateDir string = filepath.Join(xdg.StateHome, "superfile")

	// MainDir files
	ThemeFolder string = filepath.Join(SuperFileMainDir, "theme")

	// DataDir files
	LastCheckVersion string = filepath.Join(SuperFileDataDir, "lastCheckVersion")
	ThemeFileVersion string = filepath.Join(SuperFileDataDir, "themeFileVersion")
	FirstUseCheck    string = filepath.Join(SuperFileDataDir, "firstUseCheck")
	PinnedFile       string = filepath.Join(SuperFileDataDir, "pinned.json")
	ToggleDotFile    string = filepath.Join(SuperFileDataDir, "toggleDotFile")
	ToggleFooter     string = filepath.Join(SuperFileDataDir, "toggleFooter")

	// StateDir files
	LogFile     string = filepath.Join(SuperFileStateDir, "superfile.log")
	LastDirFile string = filepath.Join(SuperFileStateDir, "lastdir")

	// Trash Directories
	DarwinTrashDirectory      string = filepath.Join(HomeDir, ".Trash")
	CustomTrashDirectory      string = filepath.Join(xdg.DataHome, "Trash")
	CustomTrashDirectoryFiles string = filepath.Join(xdg.DataHome, "Trash", "files")
	CustomTrashDirectoryInfo  string = filepath.Join(xdg.DataHome, "Trash", "info")
)

// These variables are actually not fixed, they are sometimes updated dynamically
var (
	ConfigFile  string = filepath.Join(SuperFileMainDir, "config.toml")
	HotkeysFile string = filepath.Join(SuperFileMainDir, "hotkeys.toml")

	// Other state variables
	FixHotkeys    bool   = false
	FixConfigFile bool   = false
	LastDir       string = ""
	PrintLastDir  bool   = false
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
			log.Fatalf("Error: While reading config file '%s' from argument : %v", configFileArg, err)
		} else {
			ConfigFile = configFileArg
		}
	}

	hotkeyFileArg := c.String("hotkey-file")

	if hotkeyFileArg != "" {
		if _, err := os.Stat(hotkeyFileArg); err != nil {
			log.Fatalf("Error: While reading hotkey file '%s' from argument : %v", hotkeyFileArg, err)
		} else {
			HotkeysFile = hotkeyFileArg
		}
	}

	FixHotkeys = c.Bool("fix-hotkeys")
	FixConfigFile = c.Bool("fix-config-file")
	PrintLastDir = c.Bool("print-last-dir")
}
