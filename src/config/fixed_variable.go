package variable

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/adrg/xdg"
)

const (
	CurrentVersion = "v1.5.0"
	// Allowing pre-releases with non production version
	// Set this to "" for production releases
	PreReleaseSuffix = ""

	// This gives most recent non-prerelease, non-draft release
	LatestVersionURL    = "https://api.github.com/repos/yorukot/superfile/releases/latest"
	LatestVersionGithub = "github.com/yorukot/superfile/releases/latest"

	// This will not break in windows. This is a relative path for Embed FS. It uses "/" only
	EmbedConfigDir           = "src/superfile_config"
	EmbedConfigFile          = EmbedConfigDir + "/config.toml"
	EmbedHotkeysFile         = EmbedConfigDir + "/hotkeys.toml"
	EmbedThemeDir            = EmbedConfigDir + "/theme"
	EmbedThemeCatppuccinFile = EmbedThemeDir + "/catppuccin-mocha.toml"
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
	DarwinTrashDirectory = filepath.Join(HomeDir, ".Trash")

	// These are used by github.com/rkoesters/xdg/trash package
	// We need to make sure that these directories exist
	LinuxTrashDirectory      = filepath.Join(xdg.DataHome, "Trash")
	LinuxTrashDirectoryFiles = filepath.Join(xdg.DataHome, "Trash", "files")
	LinuxTrashDirectoryInfo  = filepath.Join(xdg.DataHome, "Trash", "info")
)

// These variables are actually not fixed, they are sometimes updated dynamically
var (
	ConfigFile  = filepath.Join(SuperFileMainDir, "config.toml")
	HotkeysFile = filepath.Join(SuperFileMainDir, "hotkeys.toml")

	// ChooserFile is the path where superfile will write the file's path, which is to be
	// opened, before exiting
	ChooserFile = ""

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

func SetChooserFile(path string) {
	ChooserFile = path
}

func UpdateVarFromCliArgs(c *cli.Command) {
	// Setting the config file path
	configFileArg := c.String("config-file")

	// Validate the config file exists
	if configFileArg != "" {
		if _, err := os.Stat(configFileArg); err != nil {
			utils.PrintfAndExitf("Error: While reading config file '%s' from argument : %v", configFileArg, err)
		}
		ConfigFile = configFileArg
	}

	hotkeyFileArg := c.String("hotkey-file")

	if hotkeyFileArg != "" {
		if _, err := os.Stat(hotkeyFileArg); err != nil {
			utils.PrintfAndExitf("Error: While reading hotkey file '%s' from argument : %v", hotkeyFileArg, err)
		}
		HotkeysFile = hotkeyFileArg
	}

	// It could be non existent. We are writing to the file. If file doesn't exists, we would attempt to create it.
	SetChooserFile(c.String("chooser-file"))

	FixHotkeys = c.Bool("fix-hotkeys")
	FixConfigFile = c.Bool("fix-config-file")
	PrintLastDir = c.Bool("print-last-dir")
}
