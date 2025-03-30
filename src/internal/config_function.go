package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/barasher/go-exiftool"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
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

	common.LoadConfigFile()

	logLevel := slog.LevelInfo
	if Config.Debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		file, &slog.HandlerOptions{Level: logLevel})))

	common.LoadHotkeysFile()

	common.LoadThemeFile()

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

	common.LoadThemeConfig()
	common.LoadPrerenderedVariables()

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
