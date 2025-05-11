package internal

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

/*
The purpose of this test file is to have the
(1) common global data for tests
(2) common setup for tests, and cleanup
(3) Basic model fuctionality tests
    - Initialization
	- Resize
	- Update
	- Quitting
*/

// Helps to have centralized cleanup
var testDir string //nolint: gochecknoglobals // One-time initialized, and then read-only global test variable

func cleanupTestDir() {
	err := os.RemoveAll(testDir)
	if err != nil {
		fmt.Printf("error while cleaning up test directory, err : %v", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	spfConfigDir := filepath.Join(filepath.Dir(filepath.Dir(filename)),
		"superfile_config")

	err := common.PopulateGlobalConfigs(
		filepath.Join(spfConfigDir, "config.toml"),
		filepath.Join(spfConfigDir, "hotkeys.toml"),
		filepath.Join(spfConfigDir, "theme", "monokai.toml"))

	if err != nil {
		fmt.Printf("error while populating config, err : %v", err)
		os.Exit(1)
	}

	// Create testDir
	testDir = filepath.Join(os.TempDir(), "spf_testdir")

	if err := os.Mkdir(testDir, 0755); err != nil {
		fmt.Printf("error while creating test directory, err : %v", err)
		os.Exit(1)
	}
	defer cleanupTestDir()

	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	m.Run()
	// Maybe catch panic
}

func TestBasic(t *testing.T) {
	curTestDir := filepath.Join(testDir, "TestBasic")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(dir1, "file1.txt")

	t.Run("Basic Checks", func(t *testing.T) {
		err := os.Mkdir(curTestDir, 0755)
		require.NoError(t, err)
		err = os.Mkdir(dir1, 0755)
		require.NoError(t, err)
		err = os.Mkdir(dir2, 0755)
		require.NoError(t, err)

		// Should permission be made lesser than this ?
		// Keep text in a const
		err = os.WriteFile(file1, SampleDataBytes, 0755)

		require.NoError(t, err)

		m := defaultTestModel(dir1)

		_, _ = TeaUpdate(&m, nil)

		// Validate the most of the data stored in model object
		// Inspect model struct to see what more can be validated.
		// 1 - File panel location, cursor, render index, etc.
		// 2 - Directory Items are listed
		// 3 - sidebar items pinned items are listed
		// 4 - process panel is empty
		// 5 - clipboard is empty
		// 6 - model's dimenstion
	})
}
