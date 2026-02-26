package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

var suffixRegexp = regexp.MustCompile(`^(.*)\((\d+)\)$`)

// Check if the directory is external disk path
// TODO : This function should be give two directories, and it should return
// if the two share a different disk partition.
// Ideally we shouldn't even try to figure that out in our file operations, and let OS handles it.
// But at least right now its not okay. This returns if `path` is an External disk
// from perspective of `/`, but it should tell from perspective of currently open directory
// The usage of this function in cut/paste is not as expected.
func isExternalDiskPath(path string) bool {
	// This is very vague. You cannot tell if a path is belonging to an external partition
	// if you dont define the source path to compare with
	// But making this true will cause slow file operations based on current implementation
	if runtime.GOOS == utils.OsWindows {
		return false
	}

	// exclude timemachine on macOS
	if strings.HasPrefix(path, "/Volumes/.timemachine") {
		return false
	}

	// to filter out mounted partitions like /, /boot etc
	return strings.HasPrefix(path, "/mnt") ||
		strings.HasPrefix(path, "/media") ||
		strings.HasPrefix(path, "/run/media") ||
		strings.HasPrefix(path, "/Volumes")
}

func checkFileNameValidity(name string) error {
	switch {
	case name == ".", name == "..":
		return errors.New("file name cannot be '.' or '..'")
	case strings.HasSuffix(name, fmt.Sprintf("%c.", filepath.Separator)),
		strings.HasSuffix(name, fmt.Sprintf("%c..", filepath.Separator)):
		return fmt.Errorf("file name cannot end with '%c.' or '%c..'", filepath.Separator, filepath.Separator)
	default:
		return nil
	}
}

func renameIfDuplicate(destination string) (string, error) {
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		return destination, nil
	} else if err != nil {
		return "", err
	}

	dir := filepath.Dir(destination)
	base := filepath.Base(destination)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	// Extract base name without existing suffix
	counter := 1
	//nolint:mnd // 3 = full match + 2 capture groups
	if match := suffixRegexp.FindStringSubmatch(name); len(match) == 3 {
		name = match[1] // base name without (N)
		if num, err := strconv.Atoi(match[2]); err == nil {
			counter = num + 1 // start from next number
		}
	}

	// Find first available name
	for i := counter; i < 10_000; i++ {
		newName := fmt.Sprintf("%s(%d)%s", name, i, ext)
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, nil
		}
	}

	return "", fmt.Errorf("could not find free name for %s after many attempts", destination)
}

// Count how many file in the directory
func countFiles(dirPath string) (int, error) {
	count := 0

	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})

	return count, err
}

func processCmdToTeaCmd(cmd processbar.Cmd) tea.Cmd {
	if cmd == nil {
		// To prevent us from running cmd() on nil cmd
		return nil
	}
	return func() tea.Msg {
		updateMsg := cmd()
		return ProcessBarUpdateMsg{
			pMsg: updateMsg,
			BaseMessage: BaseMessage{
				reqID: updateMsg.GetReqID(),
			},
		}
	}
}

func getCopyOrCutOperationName(cut bool) string {
	if cut {
		return "cut"
	}
	return "copy"
}
