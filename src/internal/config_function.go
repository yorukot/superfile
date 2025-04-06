package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yorukot/superfile/src/internal/common/utils"

	"github.com/barasher/go-exiftool"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// initialConfig load and handle all configuration files (spf config,Hotkeys
// themes) setted up. Returns absolute path of dir pointing to the file Panel
func initialConfig(dir string) (bool, bool, string) {
	// Open log stream
	file, err := os.OpenFile(variable.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// Todo : This could be improved if we want to make superfile more resilient to errors
	// For example if the log file directories have access issues.
	// we could pass a dummy object to log.SetOutput() and the app would still function.
	if err != nil {
		// At this point, it will go to stdout since log file is not initilized
		utils.LogAndExit("Error while opening superfile.log file", "error", err)
	}
	common.LoadConfigFile()

	logLevel := slog.LevelInfo
	if common.Config.Debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		file, &slog.HandlerOptions{Level: logLevel})))

	common.LoadHotkeysFile()

	common.LoadThemeFile()

	icon.InitIcon(common.Config.Nerdfont, common.Theme.DirectoryIconColor)

	common.LoadThemeConfig()
	common.LoadPrerenderedVariables()

	if common.Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			slog.Error("Error while initial model function init exiftool error", "error", err)
		}
	}

	// spf was not called with an argument
	firstFilePanelDir := dir
	if firstFilePanelDir == "" {
		firstFilePanelDir = common.Config.DefaultDirectory
	}

	if strings.HasPrefix(firstFilePanelDir, "~") {
		// We only need to replace the first ~ , not all of them
		// And only if its a prefix
		firstFilePanelDir = strings.Replace(firstFilePanelDir, "~", variable.HomeDir, 1)
	}
	firstFilePanelDir, err = filepath.Abs(firstFilePanelDir)
	// In case of unexpected path error, fallback to home dir
	if err != nil {
		slog.Error("Unexpected error while calculating firstFilePanelDir", "error", err)
		firstFilePanelDir = variable.HomeDir
	}

	slog.Debug("Runtime information", "runtime.GOOS", runtime.GOOS,
		"start directory", firstFilePanelDir)

	toggleDotFile := utils.ReadBoolFile(variable.ToggleDotFile, false)
	toggleFooter := utils.ReadBoolFile(variable.ToggleFooter, true)

	return toggleDotFile, toggleFooter, firstFilePanelDir
}
