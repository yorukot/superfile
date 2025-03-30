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

// initialConfig load and handle all configuration files (spf config,Hotkeys
// themes) setted up. Returns absolute path of dir pointing to the file Panel
func initialConfig(dir string) (toggleDotFileBool bool, toggleFooter bool, firstFilePanelDir string) {
	// Open log stream

	slog.Debug("hi")
	slog.Info("hi")

	common.LoadConfigFile()

	slog.Debug("Config is", "config", common.Config)

	common.LoadHotkeysFile()

	common.LoadThemeFile()

	icon.InitIcon(common.Config.Nerdfont, common.Theme.DirectoryIconColor)

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

	if common.Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			slog.Error("Error while initial model function init exiftool error", "error", err)
		}
	}

	if dir != "" {
		firstFilePanelDir, err = filepath.Abs(dir)
	} else {
		common.Config.DefaultDirectory = strings.Replace(common.Config.DefaultDirectory, "~", variable.HomeDir, -1)
		firstFilePanelDir, err = filepath.Abs(common.Config.DefaultDirectory)
	}

	if err != nil {
		firstFilePanelDir = variable.HomeDir
	}

	slog.Debug("Runtime information", "runtime.GOOS", runtime.GOOS,
		"start directory", firstFilePanelDir)

	return toggleDotFileBool, toggleFooter, firstFilePanelDir
}
