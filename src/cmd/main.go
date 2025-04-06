package cmd

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/common/utils"

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

	common.LoadInitial_PrerenderedVariables()
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
				Action: func(c *cli.Context) error {
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
			path := ""
			if c.Args().Present() {
				path = c.Args().First()
			}

			variable.UpdateVarFromCliArgs(c)

			InitConfigFile()

			hasTrash := true
			if err := InitTrash(); err != nil {
				hasTrash = false
			}

			firstUse := checkFirstUse()

			p := tea.NewProgram(internal.InitialModel(path, firstUse, hasTrash), tea.WithAltScreen(), tea.WithMouseCellMotion())
			if _, err := p.Run(); err != nil {
				log.Fatalf("Alas, there's been an error: %v", err)
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
		log.Fatalln(err)
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
		log.Fatalln("Error creating directories:", err)
	}

	// Create files
	if err := createFiles(
		variable.ToggleDotFile,
		variable.LogFile,
		variable.ThemeFileVersion,
		variable.ToggleFooter,
	); err != nil {
		log.Fatalln("Error creating files:", err)
	}

	// Write config file
	if err := writeConfigFile(variable.ConfigFile, common.ConfigTomlString); err != nil {
		log.Fatalln("Error writing config file:", err)
	}

	if err := writeConfigFile(variable.HotkeysFile, common.HotkeysTomlString); err != nil {
		log.Fatalln("Error writing config file:", err)
	}

	if err := initJsonFile(variable.PinnedFile); err != nil {
		log.Fatalln("Error initializing json file:", err)
	}
}

// We are initializing these, but not sure if we are ever using them
func InitTrash() error {
	// Create trash directories
	if runtime.GOOS != variable.OS_DARWIN {
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
			log.Fatalf("Failed to create file: %v", err)
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

func initJsonFile(path string) error {
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
// Todo : This is too big of a function. Refactor it to displayUpdateNotification, fetchLatestVersion,
// shouldCheckForUpdates, chucks
func CheckForUpdates() {
	if !common.Config.AutoCheckUpdate {
		return
	}

	// Get current time in UTC
	currentTime := time.Now().UTC()

	// Check last time the version was checked
	content, err := os.ReadFile(variable.LastCheckVersion)

	// Default to zero time if file doesn't exist, is empty, or has errors
	lastTime := time.Time{}

	if err == nil && len(content) > 0 {
		parsedTime, parseErr := time.Parse(time.RFC3339, string(content))
		if parseErr == nil {
			lastTime = parsedTime.UTC()
		} else {
			// Let the time stay as zero initialized value
			slog.Error("Error parsing time from LastCheckVersion file. Setting last time to zero", "error", parseErr)
		}
	}

	if lastTime.IsZero() || currentTime.Sub(lastTime) >= 24*time.Hour {
		// We would make sure to update the file in all return paths
		defer func() {
			writeLastCheckTime(currentTime)
		}()
		// Create a context with a 5-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // Cancel the context when done
		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, variable.LatestVersionURL, nil)
		if err != nil {
			slog.Error("Error forming updates request:", "error", err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				slog.Error("Update check request timed out")
				return
			}
			slog.Error("Error checking for updates:", "error", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err)
			return
		}

		type GitHubRelease struct {
			TagName string `json:"tag_name"`
		}

		var release GitHubRelease
		if err := json.Unmarshal(body, &release); err != nil {
			// Update the timestamp file even if JSON parsing fails
			slog.Error("Error parsing JSON from Github", "error", err)
			return
		}

		// Check if the local version is outdated
		if semver.Compare(release.TagName, variable.CurrentVersion) > 0 {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69E1")).Render("┃ ") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBA52")).Bold(true).Render("A new version ") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFF2")).Bold(true).Italic(true).Render(release.TagName) +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBA52")).Bold(true).Render(" is available."))

			fmt.Printf(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69E1")).Render("┃ ")+"Please update.\n┏\n\n      => %s\n\n", variable.LatestVersionGithub)
			fmt.Printf("                                                               ┛\n")
		}
	}
}
