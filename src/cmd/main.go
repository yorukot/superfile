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
	"time"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
	"golang.org/x/mod/semver"

	variable "github.com/yorukot/superfile/src/config"
	internal "github.com/yorukot/superfile/src/internal"
)

// Run superfile app
func Run(content embed.FS) {
	// Enable custom colored help output
	cli.HelpPrinter = CustomHelpPrinter //nolint:reassign // Intentionally reassigning to customize help output

	// Before we open log file, set all "non debug" logs to stdout
	utils.SetRootLoggerToStdout(false)

	common.LoadInitialPrerenderedVariables()
	common.LoadAllDefaultConfig(content)

	app := &cli.Command{
		Name:        "superfile",
		Version:     variable.CurrentVersion + variable.PreReleaseSuffix,
		Description: "Pretty fancy and modern terminal file manager ",
		ArgsUsage:   "[PATH]...",
		Commands: []*cli.Command{
			{
				Name:    "path-list",
				Aliases: []string{"pl"},
				Usage:   "Print the path to the configuration and directory",
				Action: func(_ context.Context, c *cli.Command) error {
					if c.Bool("lastdir-file") {
						fmt.Println(variable.LastDirFile)
						return nil
					}
					fmt.Printf("%-*s %s\n",
						common.HelpKeyColumnWidth,
						lipgloss.NewStyle().Foreground(lipgloss.Color("#66b2ff")).Render("[Configuration file path]"),
						variable.ConfigFile,
					)
					fmt.Printf("%-*s %s\n",
						common.HelpKeyColumnWidth,
						lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc66")).Render("[Hotkeys file path]"),
						variable.HotkeysFile,
					)
					logStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#66ff66"))
					configStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9999"))
					dataStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff66ff"))
					fmt.Printf("%-*s %s\n", common.HelpKeyColumnWidth,
						logStyle.Render("[Log file path]"), variable.LogFile)
					fmt.Printf("%-*s %s\n", common.HelpKeyColumnWidth,
						configStyle.Render("[Configuration directory path]"), variable.SuperFileMainDir)
					fmt.Printf("%-*s %s\n", common.HelpKeyColumnWidth,
						dataStyle.Render("[Data directory path]"), variable.SuperFileDataDir)
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "lastdir-file",
						Aliases: []string{"ld"},
						Usage:   "Print path to lastdir file (Where last dir is written when cd_on_quit config is true)",
						Value:   false,
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug-info",
				Aliases: []string{"di"},
				Usage:   "Print debug information",
				Value:   false,
			},
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
				Name:    "hotkey-file",
				Aliases: []string{"hf"},
				Usage:   "Specify the path to a different hotkey file",
				Value:   "", // Default to the blank string indicating non-usage of flag
			},
			&cli.StringFlag{
				Name:    "chooser-file",
				Aliases: []string{"cf"},
				Usage:   "On trying to open any file, superfile will write to its path to this file, and exit",
				Value:   "", // Default to the blank string indicating non-usage of flag
			},
		},
		Action: spfAppAction,
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		utils.PrintlnAndExit(err)
	}
}

func spfAppAction(_ context.Context, c *cli.Command) error {
	variable.UpdateVarFromCliArgs(c)

	if c.Bool("debug-info") {
		printDebugInfo()
		return nil
	}
	// If no args are called along with "spf" use current dir
	firstPanelPaths := []string{""}
	if c.Args().Present() {
		firstPanelPaths = c.Args().Slice()
	}

	InitConfigFile()

	firstUse := checkFirstUse()

	p := tea.NewProgram(internal.InitialModel(firstPanelPaths, firstUse),
		tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		utils.PrintfAndExitf("Alas, there's been an error: %v", err)
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
}

// Create proper directories for storing configuration and write default
// configurations to Config and Hotkeys toml
func InitConfigFile() {
	// Create directories
	if err := utils.CreateDirectories(
		variable.SuperFileMainDir,
		variable.SuperFileDataDir,
		variable.SuperFileStateDir,
		variable.ThemeFolder,
	); err != nil {
		utils.PrintlnAndExit("Error creating directories:", err)
	}

	// Create files
	if err := utils.CreateFiles(
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
}

// Check if is the first time initializing the app, if it is create
// use check file
func checkFirstUse() bool {
	file := variable.FirstUseCheck
	firstUse := false
	if _, err := os.Stat(file); os.IsNotExist(err) {
		firstUse = true
		if err = os.WriteFile(file, nil, utils.ConfigFilePerm); err != nil {
			utils.PrintfAndExitf("Failed to create file: %v", err)
		}
	}
	return firstUse
}

// Write data to the path file if it does not exists
func writeConfigFile(path, data string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.WriteFile(path, []byte(data), utils.ConfigFilePerm); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", path, err)
		}
	}
	return nil
}

func writeLastCheckTime(t time.Time) {
	err := os.WriteFile(variable.LastCheckVersion, []byte(t.Format(time.RFC3339)), utils.ConfigFilePerm)
	if err != nil {
		slog.Error("Error writing LastCheckVersion file", "error", err)
	}
}

// Check for the need of updates if AutoCheckUpdate is on, if its the first time
// that version is checked or if has more than 24h since the last version check,
// look into the repo if  there's any more recent version
func CheckForUpdates() {
	if !common.Config.AutoCheckUpdate {
		return
	}

	currentTime := time.Now().UTC()
	lastCheckTime := readLastCheckTime()

	if !shouldCheckForUpdate(currentTime, lastCheckTime) {
		return
	}

	defer writeLastCheckTime(currentTime)
	checkAndNotifyUpdate()
}

// Default to zero time if file doesn't exist, is empty, or has errors
func readLastCheckTime() time.Time {
	content, err := os.ReadFile(variable.LastCheckVersion)
	if err != nil || len(content) == 0 {
		return time.Time{}
	}

	parsedTime, parseErr := time.Parse(time.RFC3339, string(content))
	if parseErr != nil {
		slog.Error("Failed to parse LastCheckVersion timestamp", "error", parseErr)
		return time.Time{}
	}

	return parsedTime.UTC()
}

func shouldCheckForUpdate(now, last time.Time) bool {
	return last.IsZero() || now.Sub(last) >= 24*time.Hour
}

func checkAndNotifyUpdate() {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultCLIContextTimeout)
	defer cancel()

	resp, err := fetchLatestRelease(ctx)
	if err != nil {
		slog.Error("Failed to fetch update", "error", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read update response", "error", err)
		return
	}

	type GitHubRelease struct {
		TagName string `json:"tag_name"`
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		slog.Error("Failed to parse GitHub JSON", "error", err)
		return
	}

	if semver.Compare(release.TagName, variable.CurrentVersion) > 0 {
		notifyUpdateAvailable(release.TagName)
	}
}

func fetchLatestRelease(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, variable.LatestVersionURL, nil)
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
