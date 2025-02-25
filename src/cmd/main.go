package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
	variable "github.com/yorukot/superfile/src/config"
	internal "github.com/yorukot/superfile/src/internal"
)

// Run superfile app
func Run(content embed.FS) {

	internal.LoadAllDefaultConfig(content)

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
		},
		Action: func(c *cli.Context) error {
			// If no args are called along with "spf" use current dir
			path := ""
			if c.Args().Present() {
				path = c.Args().First()
			}

			// Setting the config file path
			configFileArg := c.String("config-file")

			// Validate the config file exists
			if configFileArg != "" {
				if _, err := os.Stat(variable.ConfigFile); err != nil {
					log.Fatalf("Error: While reading config file '%s' from arguement : %v", configFileArg, err)
				} else {
					variable.ConfigFile = configFileArg
				}
			}

			InitConfigFile()

			err := InitTrash()
			hasTrash := true
			if err != nil {
				hasTrash = false
			}

			variable.FixHotkeys = c.Bool("fix-hotkeys")
			variable.FixConfigFile = c.Bool("fix-config-file")
			variable.PrintLastDir = c.Bool("print-last-dir")

			firstUse := checkFirstUse()

			p := tea.NewProgram(internal.InitialModel(path, firstUse, hasTrash), tea.WithAltScreen(), tea.WithMouseCellMotion())
			if _, err := p.Run(); err != nil {
				log.Fatalf("Alas, there's been an error: %v", err)
			}

			if variable.PrintLastDir {
				fmt.Println(variable.LastDir)
			}

			CheckForUpdates()
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
		variable.PinnedFile,
		variable.ToggleDotFile,
		variable.LogFile,
		variable.ThemeFileVersion,
		variable.ToggleFooter,
	); err != nil {
		log.Fatalln("Error creating files:", err)
	}

	// Write config file
	if err := writeConfigFile(variable.ConfigFile, internal.ConfigTomlString); err != nil {
		log.Fatalln("Error writing config file:", err)
	}

	if err := writeConfigFile(variable.HotkeysFile, internal.HotkeysTomlString); err != nil {
		log.Fatalln("Error writing config file:", err)
	}
}

// We are initializing these, but not sure if we are ever using them
func InitTrash() error {
	// Create trash directories
	if runtime.GOOS != "darwin" {
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
			if err := os.MkdirAll(dir, 0755); err != nil {
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
			if err := os.WriteFile(file, nil, 0644); err != nil {
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
		if err := os.WriteFile(file, nil, 0644); err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
	}
	return firstUse
}

// Write data to the path file if it does not exists
func writeConfigFile(path, data string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(data), 0644); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", path, err)
		}
	}
	return nil
}

// Check for the need of updates if AutoCheckUpdate is on, if its the first time
// that version is checked or if has more than 24h since the last version check,
// look into the repo if  there's any more recent version
func CheckForUpdates() {
	var Config internal.ConfigType

	// Get AutoCheck flag from configuration files
	data, err := os.ReadFile(variable.ConfigFile)
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}

	err = toml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error decoding config file ( your config file may be misconfigured ): %v", err)
	}

	if !Config.AutoCheckUpdate {
		return
	}

	// Check last time the version was checked
	lastTime, err := readLastTimeCheckVersionFromFile(variable.LastCheckVersion)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading from file:", err)
		return
	}

	currentTime := time.Now().UTC()

	if lastTime.IsZero() || currentTime.Sub(lastTime) >= 24*time.Hour {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(variable.LatestVersionURL)
		if err != nil {
			slog.Error("Error checking for updates:", "error", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		type GitHubRelease struct {
			TagName string `json:"tag_name"`
		}

		var release GitHubRelease
		if err := json.Unmarshal(body, &release); err != nil {
			return
		}

		//Check if the local version is outdated
		if versionToNumber(release.TagName) > versionToNumber(variable.CurrentVersion) {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69E1")).Render("┃ ") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBA52")).Bold(true).Render("A new version ") +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFF2")).Bold(true).Italic(true).Render(release.TagName) +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBA52")).Bold(true).Render(" is available."))

			fmt.Printf(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69E1")).Render("┃ ")+"Please update.\n┏\n\n      => %s\n\n", variable.LatestVersionGithub)
			fmt.Printf("                                                               ┛\n")
		}

		timeStr := currentTime.Format(time.RFC3339)
		err = writeToFile(variable.LastCheckVersion, timeStr)
		if err != nil {
			log.Println("Error writing to file:", err)
			return
		}
	}
}

// Convert version string to number
func versionToNumber(version string) int {
	version = strings.ReplaceAll(version, "v", "")
	version = strings.ReplaceAll(version, ".", "")

	num, _ := strconv.Atoi(version)
	return num
}

// Check the last time the version file was checked
func readLastTimeCheckVersionFromFile(filename string) (time.Time, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return time.Time{}, err
	}
	if len(content) == 0 {
		return time.Time{}, nil
	}
	lastTime, err := time.Parse(time.RFC3339, string(content))
	if err != nil {
		return time.Time{}, err
	}

	// Ensure the time is in UTC
	return lastTime.UTC(), nil
}

// Write content to filename
func writeToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
