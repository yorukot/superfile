package internal

import (
	"errors"
	"log/slog"
	"os"
	"reflect"
	"runtime"

	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"

	"github.com/barasher/go-exiftool"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// initialConfig load and handle all configuration files (spf config,Hotkeys
// themes) setted up. Processes input directories and returns toggle states.

// This is the only usecase of named returns, distinguish between multiple return values
func initialConfig(firstPanelPaths []string) (toggleDotFile bool, //nolint: nonamedreturns // See above
	toggleFooter bool, zClient *zoxidelib.Client) {
	// Open log stream
	file, err := os.OpenFile(variable.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, utils.LogFilePerm)

	// TODO : This could be improved if we want to make superfile more resilient to errors
	// For example if the log file directories have access issues.
	// we could pass a dummy object to log.SetOutput() and the app would still function.
	if err != nil {
		utils.PrintfAndExitf("Error while opening superfile.log file : %v", err)
	}
	common.LoadConfigFile()

	logLevel := slog.LevelInfo
	if common.Config.Debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		file, &slog.HandlerOptions{Level: logLevel})))

	printRuntimeInfo()

	common.LoadHotkeysFile(common.Config.IgnoreMissingFields)

	common.LoadThemeFile()

	icon.InitIcon(common.Config.Nerdfont, common.Theme.DirectoryIconColor)

	common.LoadThemeConfig()
	common.LoadPrerenderedVariables()

	// TODO: Make sure to clean it up. Via et.Close()
	// Note: All the tool we use to interact with OS, should be abstracted behind a struc
	// Have exiftool manager, Zoxide Manager, OS Manager, Xtractor, Zipper, Command Executor
	if common.Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			slog.Error("Error while initial model function init exiftool error", "error", err)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("cannot get current working directory", "error", err)
		cwd = variable.HomeDir
	}

	if common.Config.ZoxideSupport {
		zClient, err = zoxidelib.New()
		if err != nil {
			slog.Error("Error initializing zoxide client", "error", err)
		}
	}

	updateFirstFilePanelPaths(firstPanelPaths, cwd, zClient)

	slog.Debug("Directory configuration", "cwd", cwd, "start_paths", firstPanelPaths)
	printRuntimeInfo()

	toggleDotFile = utils.ReadBoolFile(variable.ToggleDotFile, false)
	toggleFooter = utils.ReadBoolFile(variable.ToggleFooter, true)

	return toggleDotFile, toggleFooter, zClient
}

func updateFirstFilePanelPaths(firstPanelPaths []string, cwd string, zClient *zoxidelib.Client) {
	for i := range firstPanelPaths {
		if firstPanelPaths[i] == "" {
			firstPanelPaths[i] = common.Config.DefaultDirectory
		}
		originalPath := firstPanelPaths[i]
		firstPanelPaths[i] = utils.ResolveAbsPath(cwd, firstPanelPaths[i])
		if _, err := os.Stat(firstPanelPaths[i]); err != nil {
			slog.Error("cannot get stats", "path", firstPanelPaths[i], "error", err)
			// In case the path provided did not exist, use zoxide query
			// else, fallback to home dir
			if common.Config.ZoxideSupport && zClient != nil {
				path, err := attemptZoxideForInitPath(originalPath, zClient)
				if err != nil {
					slog.Error("Zoxide query error", "originalPath", originalPath, "error", err)
					firstPanelPaths[i] = variable.HomeDir
				} else {
					firstPanelPaths[i] = path
				}
			} else {
				firstPanelPaths[i] = variable.HomeDir
			}
		}
	}
}

func attemptZoxideForInitPath(originalPath string, zClient *zoxidelib.Client) (string, error) {
	path, err := zClient.Query(originalPath)

	if err != nil {
		return "", err
	}
	if path == "" {
		return "", errors.New("zoxide returned empty path")
	}
	if stat, statErr := os.Stat(path); statErr != nil || !stat.IsDir() {
		return "", errors.New("zoxide returned invalid path")
	}
	return path, nil
}

func printRuntimeInfo() {
	slog.Debug("Runtime information", "runtime.GOOS", runtime.GOOS)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	slog.Debug("Memory usage",
		"alloc_bytes", memStats.Alloc,
		"total_alloc_bytes", memStats.TotalAlloc,
		"heap_objects", memStats.HeapObjects,
		"sys_bytes", memStats.Sys)
	slog.Debug("Object sizes",
		"model_size_bytes", reflect.TypeOf(model{}).Size(),
		"filePanel_size_bytes", reflect.TypeOf(filepanel.Model{}).Size(),
		"sidebarModel_size_bytes", reflect.TypeOf(sidebar.Model{}).Size(),
		"renderer_size_bytes", reflect.TypeOf(rendering.Renderer{}).Size(),
		"borderConfig_size_bytes", reflect.TypeOf(rendering.BorderConfig{}).Size(),
		"process_size_bytes", reflect.TypeOf(processbar.Process{}).Size())
}
