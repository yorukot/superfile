package variable

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

var HomeDir = xdg.Home
var SuperFileMainDir = filepath.Join(xdg.ConfigHome, "superfile")
var SuperFileCacheDir = filepath.Join(xdg.CacheHome, "superfile")
var SuperFileDataDir = filepath.Join(xdg.DataHome, "superfile")
var SuperFileStateDir = filepath.Join(xdg.StateHome, "superfile")

const (
	CurrentVersion      string = "v1.1.7.1"
	LatestVersionURL    string = "https://api.github.com/repos/yorukot/superfile/releases/latest"
	LatestVersionGithub string = "github.com/yorukot/superfile/releases/latest"
)

var (
	ThemeFolder         string = filepath.Join(SuperFileMainDir, "theme")
	LastCheckVersion    string = filepath.Join(SuperFileDataDir, "lastCheckVersion")
	ThemeFileVersion    string = filepath.Join(SuperFileDataDir, "themeFileVersion")
	FirstUseCheck       string = filepath.Join(SuperFileDataDir, "firstUseCheck")
	PinnedFile          string = filepath.Join(SuperFileDataDir, "pinned.json")
	ConfigFile          string = filepath.Join(SuperFileMainDir, "config.toml")
	HotkeysFile         string = filepath.Join(SuperFileMainDir, "hotkeys.toml")
	ToggleDotFile       string = filepath.Join(SuperFileDataDir, "toggleDotFile")
	ToggleFooter        string = filepath.Join(SuperFileDataDir, "toggleFooter")
	LogFile             string = filepath.Join(SuperFileStateDir, "superfile.log")
	FixHotkeys          bool   = false
	FixConfigFile       bool   = false
	LastDir             string = ""
	PrintLastDir        bool   = false
	TrashDirectory      string = filepath.Join("/", "Trash")
	TrashDirectoryFiles string = filepath.Join("/", "Trash", "files")
	TrashDirectoryInfo  string = filepath.Join("/", "Trash", "info")
)
