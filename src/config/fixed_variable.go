package variable

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	CurrentVersion      string = "v1.1.7.1"
	LatestVersionURL    string = "https://api.github.com/repos/yorukot/superfile/releases/latest"
	LatestVersionGithub string = "github.com/yorukot/superfile/releases/latest"
)

var (
	HomeDir                   string = xdg.Home
	SuperFileMainDir          string = filepath.Join(xdg.ConfigHome, "superfile")
	SuperFileCacheDir         string = filepath.Join(xdg.CacheHome, "superfile")
	SuperFileDataDir          string = filepath.Join(xdg.DataHome, "superfile")
	SuperFileStateDir         string = filepath.Join(xdg.StateHome, "superfile")
	
	ThemeFolder               string = filepath.Join(SuperFileMainDir, "theme")
	LastCheckVersion          string = filepath.Join(SuperFileDataDir, "lastCheckVersion")
	ThemeFileVersion          string = filepath.Join(SuperFileDataDir, "themeFileVersion")
	FirstUseCheck             string = filepath.Join(SuperFileDataDir, "firstUseCheck")
	PinnedFile                string = filepath.Join(SuperFileDataDir, "pinned.json")
	ConfigFile                string = filepath.Join(SuperFileMainDir, "config.toml")
	HotkeysFile               string = filepath.Join(SuperFileMainDir, "hotkeys.toml")
	ToggleDotFile             string = filepath.Join(SuperFileDataDir, "toggleDotFile")
	ToggleFooter              string = filepath.Join(SuperFileDataDir, "toggleFooter")
	LogFile                   string = filepath.Join(SuperFileStateDir, "superfile.log")
	
	FixHotkeys                bool   = false
	FixConfigFile             bool   = false
	LastDir                   string = ""
	PrintLastDir              bool   = false

	DarwinTrashDirectory      string = filepath.Join(HomeDir, ".Trash")
	CustomTrashDirectory      string = filepath.Join(xdg.DataHome, "Trash")
	CustomTrashDirectoryFiles string = filepath.Join(xdg.DataHome, "Trash", "files")
	CustomTrashDirectoryInfo  string = filepath.Join(xdg.DataHome, "Trash", "info")
)
