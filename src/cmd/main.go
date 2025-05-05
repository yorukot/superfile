package cmd

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	variable "github.com/yorukot/superfile/src/config"
	internal "github.com/yorukot/superfile/src/internal"
	"golang.org/x/mod/semver"
)

// Run superfile app
func Run(content embed.FS) {
	// Before we open log file, set all "non debug" logs to stdout
	utils.SetRootLoggerToStdout(false)

	common.LoadInitialPrerenderedVariables()
	common.LoadAllDefaultConfig(content)

	app := &cli.App{
		Name:        "superfile",
		Version:     variable.CurrentVersion,
		Description: "Pretty fancy and modern terminal file manager ",
		ArgsUsage:   "[path]",
		Commands: []*cli.Command{
			{
				Name:    "path-list",
				Aliases: []string{"pl"},
				Usage:   "Print the path to the configuration and directory",
				Action: func(_ *cli.Context) error {
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#66b2ff")).Render("[Configuration file path]"), variable.ConfigFile)
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc66")).Render("[Hotkeys file path]"), variable.HotkeysFile)
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#66ff66")).Render("[Log file path]"), variable.LogFile)
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9999")).Render("[Configuration directory path]"), variable.SuperFileMainDir)
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff66ff")).Render("[Data directory path]"), variable.SuperFileDataDir)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "fix-hotkeys",
				Aliases: []string{"fh"},
				Usage:   "Adds any missing hotkeys to the hotkey config file",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "fix-config-file",
				Aliases: []string{"fch"},
				Usage:   "Adds any missing hotkeys to the hotkey config file",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "print-last-dir",
				Aliases: []string{"pld"},
				Usage:   "Print the last dir to stdout on exit (to use for cd)",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "config-file",
				Aliases: []string{"c"},
				Usage:   "Specify the path to a different config file",
				Value:   "", // Default to the blank string indicating non-usage of flag
			},
			&cli.StringFlag{
				Name:  "hotkey-file",
				Usage: "Specify the path to a different hotkey file",
				Value: "", // Default to the blank string indicating non-usage of flag
			},
		},
		Action: func(c *cli.Context) error {
			// If no args are called along with "spf" use current dir
			firstFilePanelDirs := []string{""}
			if c.Args().Present() {
				firstFilePanelDirs = c.Args().Slice()
			}

			variable.UpdateVarFromCliArgs(c)

			InitConfigFile()

			hasTrash := true
			if err := InitTrash(); err != nil {
				hasTrash = false
			}

			firstUse := checkFirstUse()

			p := tea.NewProgram(internal.InitialModel(firstFilePanelDirs, firstUse, hasTrash), tea.WithAltScreen(), tea.WithMouseCellMotion())
			if _, err := p.Run(); err != nil {
				utils.PrintfAndExit("Alas, there's been an error: %v", err)
			}

			// This must be after calling internal.InitialModel()
			// so that we know `common.Config` is loaded
			// Should not be a goroutine, Otherwise the main
			// goroutine will exit first, and this will not be able to finish
			CheckForUpdates()

			if variable.PrintLastDir {
				fmt.Println(variable.LastDir)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		utils.PrintlnAndExit(err)
	}
}

// Create proper directories for storing configuration and write default
// configurations to Config and Hotkeys toml
func InitConfigFile() {
	// Create directories
	if err := createDirectories(
		variable.SuperFileMainDir,
		variable.SuperFileDataDir,
		variable.SuperFileStateDir,
		variable.ThemeFolder,
	); err != nil {
		utils.PrintlnAndExit("Error creating directories:", err)
	}

	// Create files
	if err := createFiles(
		variable.ToggleDotFile,
		variable.LogFile,
		variable.ThemeFileVersion,
		variable.ToggleFooter,
	); err != nil {
		utils.PrintlnAndExit("Error creating files:", err)
	}

	// Write config file
	if err := writeConfigFile(variable.ConfigFile, common.ConfigTomlString); err != nil {
		utils.PrintlnAndExit("Error writing config file:", err)
	}

	if err := writeConfigFile(variable.HotkeysFile, common.HotkeysTomlString); err != nil {
		utils.PrintlnAndExit("Error writing config file:", err)
	}

	if err := initJSONFile(variable.PinnedFile); err != nil {
		utils.PrintlnAndExit("Error initializing json file:", err)
	}
}

// We are initializing these, but not sure if we are ever using them
func InitTrash() error {
	// Create trash directories
	if runtime.GOOS != utils.OsDarwin {
		err := createDirectories(
			variable.CustomTrashDirectory,
			variable.CustomTrashDirectoryFiles,
			variable.CustomTrashDirectoryInfo,
		)
		return err
	}
	return nil
}

// Helper functions
// Create all dirs that does not already exists
func createDirectories(dirs ...string) error {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Directory doesn't exist, create it
			if err = os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		} else if err != nil {
			// Some other error occurred while checking if the directory exists
			return fmt.Errorf("failed to check directory status %s: %w", dir, err)
		}
		// else: directory already exists
	}
	return nil
}

// Create all files if they do not exists yet
func createFiles(files ...string) error {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if err = os.WriteFile(file, nil, 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", file, err)
			}
		}
	}
	return nil
}

// Check if is the first time initializing the app, if it is create
// use check file
func checkFirstUse() bool {
	file := variable.FirstUseCheck
	firstUse := false
	if _, err := os.Stat(file); os.IsNotExist(err) {
		firstUse = true
		if err = os.WriteFile(file, nil, 0644); err != nil {
			utils.PrintfAndExit("Failed to create file: %v", err)
		}
	}
	return firstUse
}

// Write data to the path file if it does not exists
func writeConfigFile(path, data string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.WriteFile(path, []byte(data), 0644); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", path, err)
		}
	}
	return nil
}

func initJSONFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.WriteFile(path, []byte("null"), 0644); err != nil {
			return fmt.Errorf("failed to initialize json file %s: %w", path, err)
		}
	}
	return nil
}

func writeLastCheckTime(t time.Time) {
	err := os.WriteFile(variable.LastCheckVersion, []byte(t.Format(time.RFC3339)), 0644)
	if err != nil {
		slog.Error("Error writing LastCheckVersion file", "error", err)
	}
}

// Check for the need of updates if AutoCheckUpdate is on, if its the first time
// that version is checked or if has more than 24h since the last version check,
// look into the repo if  there's any more recent version
func CheckForUpdates(){
	if !commonConfig.AutoCheckUpdate{
		return
	}
	
	currentTime := time.Now().UTC()
	lastCheckTime := readLastCheckTime()

	if !shouldCheckForUpdate(currentTime, lastCheckTime){
		return
	}
	
	defer writeLastCheckTime(currentTime)
	checkAndNotifyUpdate(currentTime)
}

func readLastCheckTime() time.Time{
	content, err := os.ReadFile(variable.LastCheckVersion)
	if err != nill || len(content) == 0 {
		return time.Time{}
	}

	parsedTime, parseErr := timeParse(time.RFC3339, string(content))
	if parseErr != nil {
		slog.Error("Failed to parse LastCheckVersion timestamp", "error", parseErr)
		return timeTime{}
	}

	return parsedTime.UTC()
}

func shouldCheckForUpdate(now, last time.Time) bool {
	return  last.IsZero() || now.Sub(last) >= 24*time.Hour
}

func checkAndNotifyUpdate(cureentTime time.Time) {
	ctx, cancel := context.withTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := fetchLastestRelease(ctx)
	if err != nil {
		slog.Error("Failed to fetch update", "error", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read update response","error", err)
		return
	}

	if semver.Compare(release.TagName, variable.CurrentVersion) > 0 {
		notifyUpdateAvailable(release.TagName)
	}
}

func fetchLastestRelease(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, variable.LastVersionURL, nil)
	if err != nil {
		return nil, err
	}
	return (&http.Client{}).Do(req)
}

func notifyUpdateAvailable(latest string) {
	fmt.Println(
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69E1")).Render("┃ ") +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBA52")).Bold(true).Render("A new version ") +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFF2")).Bold(true).Italic(true).Render(latest) +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBA52")).Bold(true).Render(" is available."),
	)
	fmt.Printf(
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69E1")).Render("┃ ")+"Please update.\n┏\n\n      => %s\n\n",
		variable.LatestVersionGithub,
	)
	fmt.Println("                                                               ┛")
}
