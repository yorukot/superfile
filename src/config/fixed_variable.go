package varibale

import "github.com/adrg/xdg"

var HomeDir = xdg.Home
var SuperFileMainDir = xdg.ConfigHome + "/superfile"
var SuperFileCacheDir = xdg.CacheHome + "/superfile"
var SuperFileDataDir = xdg.DataHome + "/superfile"
var SuperFileStateDir = xdg.StateHome + "/superfile"

const (
	CurrentVersion      string = "v1.1.4"
	LatestVersionURL    string = "https://api.github.com/repos/yorukot/superfile/releases/latest"
	LatestVersionGithub string = "github.com/yorukot/superfile/releases/latest"
)

var (
	ThemeFoldera      string = SuperFileMainDir + "/theme"
	LastCheckVersiona string = SuperFileDataDir + "/lastCheckVersion"
	ThemeFileVersiona string = SuperFileDataDir + "/themeFileVersion"
	FirstUseChecka    string = SuperFileDataDir + "/firstUseCheck"
	PinnedFilea       string = SuperFileDataDir + "/pinned.json"
	ConfigFilea       string = SuperFileMainDir + "/config.toml"
	HotkeysFilea      string = SuperFileMainDir + "/hotkeys.toml"
	ToggleDotFilea    string = SuperFileDataDir + "/toggleDotFile"
	LogFilea          string = SuperFileStateDir + "/superfile.log"
)

const (
	TrashDirectory      string = "/Trash"
	TrashDirectoryFiles string = "/Trash/files"
	TrashDirectoryInfo  string = "/Trash/info"
)
