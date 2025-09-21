package sidebar

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Load(t *testing.T) {

	_, curFilename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	testDataDir := filepath.Join(filepath.Dir(curFilename), "testdata", "pinnedFile")

	// Use t.TempDir for valid and nonexistent cases,
	// because you need a valid path that is machine derived
	tempDir := t.TempDir()
	pinnedDir := filepath.Join(tempDir, "pinnedDir")
	err := os.Mkdir(pinnedDir, 0755)
	require.NoError(t, err)

	// Read sample case, replace with valid path and write a new file in pinnedDir.
	valid, err := os.ReadFile(filepath.Join(testDataDir, "valid.json"))
	require.NoError(t, err)
	newValid := strings.ReplaceAll(string(valid), "/REPLACE/ME/WITH/VALID/DIR", pinnedDir)
	newValidPath := filepath.Join(pinnedDir, "valid.json")
	err = os.WriteFile(newValidPath, []byte(newValid), 0644)
	require.NoError(t, err)

	nonexistent, err := os.ReadFile(filepath.Join(testDataDir, "nonexistent.json"))
	require.NoError(t, err)
	newNonexistent := strings.ReplaceAll(string(nonexistent), "/REPLACE/ME/WITH/VALID/DIR", pinnedDir)
	newNonexistentPath := filepath.Join(pinnedDir, "nonexistent.json")
	err = os.WriteFile(newNonexistentPath, []byte(newNonexistent), 0644)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		pinnedMgr PinnedManager
		expected  []directory
	}{
		{
			name:      "Empty No Pinned Directories",
			pinnedMgr: PinnedManager{filePath: filepath.Join(testDataDir, "empty.json")},
			expected:  []directory{},
		},
		{
			name:      "Invalid Format File",
			pinnedMgr: PinnedManager{filePath: filepath.Join(testDataDir, "invalid.json")},
			expected:  []directory{},
		},
		{
			name:      "Valid With No Non-Existent Directories",
			pinnedMgr: PinnedManager{filePath: newValidPath},
			expected: []directory{
				{
					Location: pinnedDir,
					Name:     "pinnedDir",
				},
			},
		},
		{
			name:      "Valid With One Non-Existent Directory",
			pinnedMgr: PinnedManager{filePath: newNonexistentPath},
			expected: []directory{
				{
					Location: pinnedDir,
					Name:     "pinnedDir",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pinnedMgr.Load())
		})
	}
}

func Test_Save(t *testing.T) {

	tempDir := t.TempDir()
	pinnedDir := filepath.Join(tempDir, "pinnedDir")
	err := os.Mkdir(pinnedDir, 0755)
	require.NoError(t, err)

	savePath := filepath.Join(pinnedDir, "pinned.json")

	rOnlyPath := filepath.Join(pinnedDir, "pinnedRonly.json")
	file, err := os.OpenFile(rOnlyPath, os.O_CREATE|os.O_RDONLY, 0400)
	require.NoError(t, err)
	file.Close()

	dirs := []directory{
		{
			Location: pinnedDir,
			Name:     "pinnedDir",
		},
	}

	// Marshalling fails (this failure is practically impossible with
	// the current setup as directory{} is perfectly JSON-safe)

	testCases := []struct {
		name      string
		pinnedMgr PinnedManager
		noError   bool
		toSave    []directory
		expected  []directory
	}{
		{
			name:      "Valid Normal Case",
			pinnedMgr: PinnedManager{filePath: savePath},
			noError:   true,
			toSave:    dirs,
			expected:  dirs,
		},
		{
			name:      "Empty Slice",
			pinnedMgr: PinnedManager{filePath: savePath},
			noError:   true,
			toSave:    []directory{},
			expected:  []directory{},
		},
		{
			name:      "Write Failure",
			pinnedMgr: PinnedManager{filePath: rOnlyPath},
			noError:   false,
			toSave:    dirs,
			expected:  nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pinnedMgr.Save(tt.toSave)
			if tt.noError {
				require.NoError(t, err)

				result := tt.pinnedMgr.Load()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func Test_Toggle(t *testing.T) {

	tempDir := t.TempDir()

	pinnedDir := filepath.Join(tempDir, "pinnedDir")
	err := os.Mkdir(pinnedDir, 0755)
	require.NoError(t, err)
	pinnedFile := filepath.Join(pinnedDir, "pinned.json")

	nonexistentDir := filepath.Join(tempDir, "nonExistentDir9")

	mgr := &PinnedManager{filePath: pinnedFile}

	testCases := []struct {
		name      string
		pinnedMgr *PinnedManager
		expected  []directory
		noError   bool
		arg_dir   string
	}{
		{
			name:      "Add non existing Directory to Pinned",
			pinnedMgr: mgr,
			expected:  []directory{},
			noError:   true,
			arg_dir:   nonexistentDir,
		},
		{
			name:      "Add a Directory to Pinned",
			pinnedMgr: mgr,
			expected: []directory{
				{
					Location: pinnedDir,
					Name:     "pinnedDir",
				},
			},
			noError: true,
			arg_dir: pinnedDir,
		},
		{
			name:      "Remove a Directory from Pinned",
			pinnedMgr: mgr,
			expected:  []directory{},
			noError:   true,
			arg_dir:   pinnedDir,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.pinnedMgr.Toggle(tt.arg_dir)
			if tt.noError {
				require.NoError(t, err)

				dirs := tt.pinnedMgr.Load()
				assert.Equal(t, tt.expected, dirs)
			}
		})
	}
}

func Test_Clean(t *testing.T) {

	tempDir := t.TempDir()

	pinnedDir := filepath.Join(tempDir, "pinnedDir")
	err := os.Mkdir(pinnedDir, 0755)
	require.NoError(t, err)

	nonexistentDir := filepath.Join(tempDir, "nonexistentDir")

	pinnedFile := filepath.Join(pinnedDir, "pinned.json")

	rOnlyPath := filepath.Join(pinnedDir, "pinnedRonly.json")
	file, err := os.OpenFile(rOnlyPath, os.O_CREATE|os.O_RDONLY, 0400)
	require.NoError(t, err)
	file.Close()

	cleanDirs := []directory{
		{
			Location: pinnedDir,
			Name:     "pinnedDir",
		},
	}

	testCases := []struct {
		name      string
		pinnedMgr PinnedManager
		modified  bool
		expected  []directory
		arg_dirs  []directory
	}{
		{
			name:      "All Directories Exist",
			pinnedMgr: PinnedManager{filePath: pinnedFile},
			modified:  false,
			expected:  cleanDirs,
			arg_dirs:  cleanDirs,
		},
		{
			name:      "Some Directories Exist",
			pinnedMgr: PinnedManager{filePath: pinnedFile},
			modified:  true,
			expected:  cleanDirs,
			arg_dirs: append(cleanDirs, directory{
				Location: nonexistentDir,
				Name:     "nonexistentDir",
			}),
		},
		{
			name:      "Save Fails",
			pinnedMgr: PinnedManager{filePath: rOnlyPath},
			modified:  false,
			expected:  cleanDirs,
			arg_dirs: append(cleanDirs, directory{
				Location: nonexistentDir,
				Name:     "nonexistentDir",
			}),
		},
		{
			name:      "Empty Input Slice",
			pinnedMgr: PinnedManager{filePath: pinnedFile},
			modified:  false,
			expected:  []directory{},
			arg_dirs:  []directory{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			before, err := os.Stat(pinnedFile)
			if !errors.Is(err, fs.ErrNotExist) {
				require.NoError(t, err)
			}

			cleaned := tt.pinnedMgr.Clean(tt.arg_dirs)

			after, err := os.Stat(pinnedFile)
			if !errors.Is(err, fs.ErrNotExist) {
				require.NoError(t, err)
			} else if before != nil && !tt.modified {
				require.Equal(t, before.ModTime(), after.ModTime())
			}

			assert.Equal(t, tt.expected, cleaned)
		})
	}
}
