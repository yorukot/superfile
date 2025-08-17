package internal

import (
	"log/slog"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
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

// This is the only usecase of named returns, distinguish between multiple return values
func initialConfig(firstFilePanelDirs []string) (toggleDotFile bool, //nolint: nonamedreturns // See above
	toggleFooter bool) {
	// Open log stream
	file, err := os.OpenFile(variable.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	// TODO : This could be improved if we want to make superfile more resilient to errors
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

	printRuntimeInfo()

	common.LoadHotkeysFile()

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

	updateFirstFilePanelDirs(firstFilePanelDirs, cwd)

	slog.Debug("Directory configuration", "cwd", cwd, "start_directories", firstFilePanelDirs)
	printRuntimeInfo()

	toggleDotFile = utils.ReadBoolFile(variable.ToggleDotFile, false)
	toggleFooter = utils.ReadBoolFile(variable.ToggleFooter, true)

	return toggleDotFile, toggleFooter
}

func updateFirstFilePanelDirs(firstFilePanelDirs []string, cwd string) {
	for i := range firstFilePanelDirs {
		if firstFilePanelDirs[i] == "" {
			firstFilePanelDirs[i] = common.Config.DefaultDirectory
		}

		if common.Config.ZoxideSupport {
			// Execute zoxide to get the absolute path
			out, err := exec.Command("zoxide", "query", firstFilePanelDirs[i]).Output()
			if err != nil {
				slog.Error("Error while executing zoxide", "error", err, "path", firstFilePanelDirs[i])
				firstFilePanelDirs[i] = utils.ResolveAbsPath(cwd, firstFilePanelDirs[i])
			} else {
				// Remove trailing newline and whitespace from the output
				trimmedOutput := strings.TrimSpace(string(out))

				if trimmedOutput == "" {
					slog.Error("Zoxide returned empty output", "path", firstFilePanelDirs[i])
					firstFilePanelDirs[i] = utils.ResolveAbsPath(cwd, firstFilePanelDirs[i])
				} else {
					firstFilePanelDirs[i] = trimmedOutput
				}
			}
		} else {
			firstFilePanelDirs[i] = utils.ResolveAbsPath(cwd, firstFilePanelDirs[i])
		}
		// In case of unexpected path error, fallback to home dir
		if _, err := os.Stat(firstFilePanelDirs[i]); err != nil {
			slog.Error("cannot get stats for firstFilePanelDir", "error", err)
			firstFilePanelDirs[i] = variable.HomeDir
		}
	}
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
		"filePanel_size_bytes", reflect.TypeOf(filePanel{}).Size(),
		"sidebarModel_size_bytes", reflect.TypeOf(sidebar.Model{}).Size(),
		"renderer_size_bytes", reflect.TypeOf(rendering.Renderer{}).Size(),
		"borderConfig_size_bytes", reflect.TypeOf(rendering.BorderConfig{}).Size(),
		"process_size_bytes", reflect.TypeOf(processbar.Process{}).Size())
}
