package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/barasher/go-exiftool"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// initialConfig load and handle all configuration files (spf config,Hotkeys
// themes) setted up. Processes input directories and returns toggle states.
func initialConfig(firstFilePanelDirs []string) (toggleDotFile bool, toggleFooter bool) { //nolint: nonamedreturns // This is the only usecase of named returns, distinguish between multiple return values
	// Open log stream
	file, err := os.OpenFile(variable.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// Todo : This could be improved if we want to make superfile more resilient to errors
	// For example if the log file directories have access issues.
	// we could pass a dummy object to log.SetOutput() and the app would still function.
	if err != nil {
		utils.PrintfAndExit("Error while opening superfile.log file : %v", err)
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

	for i := range firstFilePanelDirs {
		if firstFilePanelDirs[i] == "" {
			firstFilePanelDirs[i] = common.Config.DefaultDirectory
		}

		if strings.HasPrefix(firstFilePanelDirs[i], "~") {
			// We only need to replace the first ~ , not all of them
			// And only if its a prefix
			firstFilePanelDirs[i] = strings.Replace(firstFilePanelDirs[i], "~", variable.HomeDir, 1)
		}
		firstFilePanelDirs[i], err = filepath.Abs(firstFilePanelDirs[i])
		// In case of unexpected path error, fallback to home dir
		if err != nil {
			slog.Error("Unexpected error while calculating firstFilePanelDir", "error", err)
			firstFilePanelDirs[i] = variable.HomeDir
		}
	}

	slog.Debug("Runtime information", "runtime.GOOS", runtime.GOOS)
	slog.Debug("Directory configuration", "start_directories", firstFilePanelDirs)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	slog.Debug("Memory usage",
		"alloc_bytes", memStats.Alloc,
		"total_alloc_bytes", memStats.TotalAlloc,
		"heap_objects", memStats.HeapObjects,
		"sys_bytes", memStats.Sys)
	slog.Debug("Object sizes",
		"model_size_bytes", reflect.TypeOf(model{}).Size(),
		"filePanel_size_bytes", reflect.TypeOf(filePanel{}).Size(),
		"sidebarModel_size_bytes", reflect.TypeOf(sidebar.Model{}).Size(),
		"renderer_size_bytes", reflect.TypeOf(rendering.Renderer{}).Size(),
		"borderConfig_size_bytes", reflect.TypeOf(rendering.BorderConfig{}).Size())

	toggleDotFile = utils.ReadBoolFile(variable.ToggleDotFile, false)
	toggleFooter = utils.ReadBoolFile(variable.ToggleFooter, true)

	return toggleDotFile, toggleFooter
}
