package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/utils"
)

const (
	keyWidth         = 20
	maxVersionLength = 50
)

type debugPrinter struct {
	titleColor   *color.Color
	flagColor    *color.Color
	warningColor *color.Color
	successColor *color.Color
}

func newDebugPrinter() *debugPrinter {
	return &debugPrinter{
		titleColor:   color.New(color.FgGreen, color.Bold),
		flagColor:    color.New(color.FgCyan, color.Bold),
		warningColor: color.New(color.FgRed, color.Bold),
		successColor: color.New(color.FgGreen),
	}
}

func printDebugInfo() {
	dp := newDebugPrinter()

	fmt.Println()
	dp.printHeader("Superfile")
	dp.printKeyValue("Version", variable.CurrentVersion+variable.PreReleaseSuffix)

	fmt.Println()
	dp.printHeader("System")
	dp.printKeyValue("OS", runtime.GOOS)
	dp.printKeyValue("Arch", runtime.GOARCH)
	if kernel, err := getKernelVersion(); err == nil {
		dp.printKeyValue("Kernel", kernel)
	}

	fmt.Println()
	dp.printHeader("Configuration")
	dp.printKeyValue("Config File", variable.ConfigFile)
	dp.printKeyValue("Hotkeys File", variable.HotkeysFile)
	dp.printKeyValue("Theme Folder", variable.ThemeFolder)
	dp.printKeyValue("Log File", variable.LogFile)
	dp.printKeyValue("Data Dir", variable.SuperFileDataDir)

	fmt.Println()
	dp.printHeader("Environment")
	if runtime.GOOS == utils.OsWindows {
		dp.printEnv("COMSPEC")
		dp.printEnv("APPDATA")
		dp.printEnv("LOCALAPPDATA")
	} else {
		dp.printEnv("TERM")
		dp.printEnv("TERM_PROGRAM")
		dp.printEnv("TERM_PROGRAM_VERSION")
		dp.printEnv("SHELL")
		dp.printEnv("EDITOR")
		dp.printEnv("VISUAL")
		dp.printEnv("XDG_SESSION_TYPE")
		dp.printEnv("WAYLAND_DISPLAY")
		dp.printEnv("DISPLAY")
	}

	fmt.Println()
	dp.printHeader("Dependencies")
	dp.checkDependency("ffmpeg", "-version")
	dp.checkDependency("pdftoppm", "-v")
	dp.checkDependency("exiftool", "-ver")
	dp.checkDependency("bat", "--version")
	dp.checkDependency("zoxide", "--version")
	switch runtime.GOOS {
	case utils.OsDarwin:
		dp.checkDependency("open", "")
		dp.checkDependency("pbcopy", "")
	case utils.OsWindows:
		dp.checkDependency("clip", "")
	case utils.OsLinux:
		dp.checkDependency("xdg-open", "--version")
		dp.checkDependency("wl-copy", "--version")
		dp.checkDependency("xclip", "-version")
		dp.checkDependency("xsel", "--version")
	}
}

func (dp *debugPrinter) printHeader(text string) {
	_, _ = dp.titleColor.Add(color.Underline).Println(text)
}

func (dp *debugPrinter) printKeyValue(key, value string) {
	if filepath.IsAbs(value) {
		if _, err := os.Stat(value); os.IsNotExist(err) {
			value = dp.warningColor.Sprint(value + " (Not Found)")
		}
	}

	// Use fixed width formatting for key
	keyStr := fmt.Sprintf("%-*s", keyWidth, key)
	_, _ = dp.flagColor.Print(keyStr)
	fmt.Printf(": %s\n", value)
}

func (dp *debugPrinter) printEnv(key string) {
	val := os.Getenv(key)
	if val == "" {
		val = "Not Set"
	}
	dp.printKeyValue(key, val)
}

func (dp *debugPrinter) checkDependency(name string, flag string) {
	path, err := exec.LookPath(name)
	var status string
	if err != nil {
		status = dp.warningColor.Sprint("Not Found")
	} else {
		// Try to get version
		version := "Found at " + path
		if flag != "" {
			//nolint:gosec // flags are hardcoded strings
			cmd := exec.Command(name, strings.Split(flag, " ")...)
			out, err := cmd.CombinedOutput()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				if len(lines) > 0 {
					v := strings.TrimSpace(lines[0])
					if len(v) > maxVersionLength {
						v = v[:maxVersionLength] + "..."
					}
					if v != "" {
						version = v
					}
				}
			}
		}
		status = dp.successColor.Sprint(version)
	}

	keyStr := fmt.Sprintf("%-*s", keyWidth, name)
	_, _ = dp.flagColor.Print(keyStr)
	fmt.Printf(": %s\n", status)
}

func getKernelVersion() (string, error) {
	if runtime.GOOS == utils.OsWindows {
		cmd := exec.Command("cmd", "/c", "ver")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
