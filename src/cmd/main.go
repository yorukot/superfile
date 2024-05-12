package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	internal "github.com/MHNightCat/superfile/src/internal"
	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var HomeDir = xdg.Home
var SuperFileMainDir = xdg.ConfigHome + "/superfile"
var SuperFileCacheDir = xdg.CacheHome + "/superfile"
var SuperFileDataDir = xdg.DataHome + "/superfile"
var SuperFileStateDir = xdg.StateHome + "/superfile"

const (
	currentVersion      string = "v1.1.3"
	latestVersionURL    string = "https://api.github.com/repos/MHNightCat/superfile/releases/latest"
	latestVersionGithub string = "github.com/MHNightCat/superfile/releases/latest"
	themeZip            string = "https://github.com/MHNightCat/superfile/raw/main/themeZip/v1.1.3/theme.zip"
)

const (
	themeFolder      string = "/theme"
	lastCheckVersion string = "/lastCheckVersion"
	themeFileVersion string = "/themeFileVersion"
	firstUseCheck    string = "/firstUseCheck"
	pinnedFile       string = "/pinned.json"
	configFile       string = "/config.toml"
	hotkeysFile      string = "/hotkeys.toml"
	toggleDotFile    string = "/toggleDotFile"
	themeZipName     string = "/theme.zip"
	logFile          string = "/superfile.log"
)

const (
	trashDirectory      string = "/Trash"
	trashDirectoryFiles string = "/Trash/files"
	trashDirectoryInfo  string = "/Trash/info"
)

func Run() {
	app := &cli.App{
		Name:        "superfile",
		Version:     currentVersion,
		Description: "Pretty fancy and modern terminal file manager ",
		ArgsUsage:   "[path]",
		Commands: []*cli.Command{
			{
				Name:  "download-theme",
				Aliases: []string{"dt"},
				Usage: "Download or update theme file from github",
				Action: func(c *cli.Context) error {
					if err := downloadAndInstallTheme(SuperFileMainDir, themeZipName, themeZip, themeFolder); err != nil {
						log.Fatalln("Error downloading theme:", err)
					}
					fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#66ff66")).Render("🎉 Successful downloaded themes!"))
					fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#66b2ff")).Render("Now you can change your superfile theme!"))
					return nil
				},
			},
			{
				Name:  "path-list",
				Aliases: []string{"pl"},
				Usage: "Print the path to the configuration and directory",
				Action: func(c *cli.Context) error {
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#66b2ff")).Render("[Configuration file path]"), filepath.Join(SuperFileMainDir, configFile))
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc66")).Render("[Hotkeys file path]"), filepath.Join(SuperFileMainDir, hotkeysFile))
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#66ff66")).Render("[Log file path]"), filepath.Join(SuperFileStateDir, logFile))
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9999")).Render("[Configuration directory path]"), SuperFileMainDir)
					fmt.Printf("%-*s %s\n", 55, lipgloss.NewStyle().Foreground(lipgloss.Color("#ff66ff")).Render("[Data directory path]"), SuperFileDataDir)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			path := ""
			if c.Args().Present() {
				path = c.Args().First()
			}

			InitConfigFile()

			firstUse := checkFirstUse()

			p := tea.NewProgram(internal.InitialModel(path, firstUse), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Fatalf("Alas, there's been an error: %v", err)
				os.Exit(1)
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

func InitConfigFile() {
	config := struct {
		MainDir      string
		DataDir      string
		StateDir     string
		PinnedFile   string
		ToggleFile   string
		LogFile      string
		ConfigFile   string
		HotkeysFile  string
	}{
		MainDir:      SuperFileMainDir,
		DataDir:      SuperFileDataDir,
		StateDir:     SuperFileStateDir,
		PinnedFile:   pinnedFile,
		ToggleFile:   toggleDotFile,
		LogFile:      logFile,
		ConfigFile:   configFile,
		HotkeysFile:  hotkeysFile,
	}

	// Create directories
	if err := createDirectories(
		config.MainDir,
		config.DataDir,
		config.StateDir,
	); err != nil {
		log.Fatalln("Error creating directories:", err)
	}

	// Create trash directories
	if runtime.GOOS != "darwin" {
		if err := createDirectories(
			xdg.DataHome+trashDirectory,
			xdg.DataHome+trashDirectoryFiles,
			xdg.DataHome+trashDirectoryInfo,
		); err != nil {
			log.Fatalln("Error creating directories:", err)
		}
	}

	// Create files
	if err := createFiles(
		config.DataDir+config.PinnedFile,
		config.DataDir+config.ToggleFile,
		config.StateDir+config.LogFile,
	); err != nil {
		log.Fatalln("Error creating files:", err)
	}

	// Write config file
	if err := writeConfigFile(config.MainDir+config.ConfigFile, internal.ConfigTomlString); err != nil {
		log.Fatalln("Error writing config file:", err)
	}

	if err := writeConfigFile(config.MainDir+config.HotkeysFile, internal.HotkeysTomlString); err != nil {
		log.Fatalln("Error writing config file:", err)
	}
}

// Helper functions
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

func checkFirstUse() bool {
	file := SuperFileDataDir + firstUseCheck
	firstUse := false
	if _, err := os.Stat(file); os.IsNotExist(err) {
		firstUse = true
		if err := os.WriteFile(file, nil, 0644); err != nil {
			log.Fatalln("failed to create file: %w", err)
		}
	}
	return firstUse
}
func writeConfigFile(path, data string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(data), 0644); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", path, err)
		}
	}
	return nil
}

func downloadAndInstallTheme(dir, zipName, zipUrl, zipFolder string) error {
	if !connected() {
		log.Fatalln("You don't have an internet connection!")
	}
	currentThemeVersion, err := readThemeVersionFromFile(SuperFileDataDir + themeFileVersion)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading from file:", err)
	}

	if _, err := os.Stat(filepath.Join(dir, zipFolder)); os.IsNotExist(err) || currentThemeVersion != currentVersion {

		err := downloadFile(filepath.Join(SuperFileMainDir, zipName), zipUrl)
		if err != nil {
			return err
		}
		err = unzip(filepath.Join(SuperFileMainDir, zipName), dir)
		if err != nil {
			return err
		} else {
			os.Remove(filepath.Join(SuperFileMainDir, zipName))
		}
		writeToFile(SuperFileDataDir+themeFileVersion, currentVersion)
	}
	return nil
}

func CheckForUpdates() {
	lastTime, err := readLastTimeCheckVersionFromFile(SuperFileDataDir + lastCheckVersion)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading from file:", err)
		return
	}

	currentTime := time.Now()

	if lastTime.IsZero() || currentTime.Sub(lastTime) >= 24*time.Hour {
		resp, err := http.Get(latestVersionURL)
		if err != nil {
			fmt.Println("Error checking for updates:", err)
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

		if versionToNumber(release.TagName) > versionToNumber(currentVersion) {
			fmt.Printf("A new version %s is available.\n", release.TagName)
			fmt.Printf("Please update.\n┏\n\n        %s\n\n", latestVersionGithub)
			fmt.Printf("                                                               ┛\n")
		}

		timeStr := currentTime.Format(time.RFC3339)
		err = writeToFile(SuperFileDataDir+lastCheckVersion, timeStr)
		if err != nil {
			log.Println("Error writing to file:", err)
			return
		}
	}
}

func versionToNumber(version string) int {
	version = strings.ReplaceAll(version, "v", "")
	version = strings.ReplaceAll(version, ".", "")

	num, _ := strconv.Atoi(version)
	return num
}

func readThemeVersionFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if len(content) == 0 {
		return "", nil
	}

	return string(content), nil
}

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

	return lastTime, nil
}

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

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func connected() (ok bool) {
    _, err := http.Get("http://clients3.google.com/generate_204")
    return err == nil
}
