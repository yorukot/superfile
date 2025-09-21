package sidebar

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Load(t *testing.T) {
	tempDir := t.TempDir()
	pDirName := "pinnedDir"
	pinnedDir := filepath.Join(tempDir, pDirName)
	err := os.Mkdir(pinnedDir, 0755)
	require.NoError(t, err)

	emptyData := []directory{}
	emptyBytes, err := json.Marshal(emptyData)
	require.NoError(t, err)
	emptyPath := filepath.Join(pinnedDir, "empty.json")
	err = os.WriteFile(emptyPath, emptyBytes, 0644)
	require.NoError(t, err)

	invalidPath := filepath.Join(pinnedDir, "invalid.json")
	err = os.WriteFile(invalidPath, []byte("{ invalid json }"), 0644)
	require.NoError(t, err)

	validData := []directory{
		{
			Location: pinnedDir,
			Name:     pDirName,
		},
	}
	validBytes, err := json.Marshal(validData)
	require.NoError(t, err)
	newValidPath := filepath.Join(pinnedDir, "valid.json")
	err = os.WriteFile(newValidPath, validBytes, 0644)
	require.NoError(t, err)

	nonexistData := []directory{
		{
			Location: pinnedDir,
			Name:     pDirName,
		},
		{
			Location: filepath.Join(pinnedDir, "nonexistent9"),
			Name:     "nonexistent9",
		},
	}
	nonexistBytes, err := json.Marshal(nonexistData)
	require.NoError(t, err)
	newNonexistentPath := filepath.Join(pinnedDir, "nonexistent.json")
	err = os.WriteFile(newNonexistentPath, nonexistBytes, 0644)
	require.NoError(t, err)

	cleanDirs := []directory{
		{
			Location: pinnedDir,
			Name:     pDirName,
		},
	}

	testCases := []struct {
		name      string
		pinnedMgr PinnedManager
		expected  []directory
	}{
		{
			name:      "Empty No Pinned Directories",
			pinnedMgr: PinnedManager{filePath: emptyPath},
			expected:  []directory{},
		},
		{
			name:      "Invalid Format File",
			pinnedMgr: PinnedManager{filePath: emptyPath},
			expected:  []directory{},
		},
		{
			name:      "Valid With No Non-Existent Directories",
			pinnedMgr: PinnedManager{filePath: newValidPath},
			expected:  cleanDirs,
		},
		{
			name:      "Valid With One Non-Existent Directory",
			pinnedMgr: PinnedManager{filePath: newNonexistentPath},
			expected:  cleanDirs,
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
	pDirName := "pinnedDir"
	pinnedDir := filepath.Join(tempDir, pDirName)
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
			Name:     pDirName,
		},
	}

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
				assert.Equal(t, tt.expected, tt.pinnedMgr.Load())
			}
		})
	}
}

func Test_Toggle(t *testing.T) {
	tempDir := t.TempDir()
	pDirName := "pinnedDir"
	pinnedDir := filepath.Join(tempDir, pDirName)
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
		argDir    string
	}{
		{
			name:      "Add non existing Directory to Pinned",
			pinnedMgr: mgr,
			expected:  []directory{},
			noError:   true,
			argDir:    nonexistentDir,
		},
		{
			name:      "Add a Directory to Pinned",
			pinnedMgr: mgr,
			expected: []directory{
				{
					Location: pinnedDir,
					Name:     pDirName,
				},
			},
			noError: true,
			argDir:  pinnedDir,
		},
		{
			name:      "Remove a Directory from Pinned",
			pinnedMgr: mgr,
			expected:  []directory{},
			noError:   true,
			argDir:    pinnedDir,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pinnedMgr.Toggle(tt.argDir)
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
	pDirName := "pinnedDir"
	pinnedDir := filepath.Join(tempDir, pDirName)
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
			Name:     pDirName,
		},
	}
	badDirs := append([]directory{}, cleanDirs...)
	badDirs = append(badDirs, directory{
		Location: nonexistentDir,
		Name:     "nonexistentDir",
	})

	testCases := []struct {
		name      string
		pinnedMgr PinnedManager
		modified  bool
		expected  []directory
		argDirs   []directory
	}{
		{
			name:      "All Directories Exist",
			pinnedMgr: PinnedManager{filePath: pinnedFile},
			modified:  false,
			expected:  cleanDirs,
			argDirs:   cleanDirs,
		},
		{
			name:      "Some Directories Exist",
			pinnedMgr: PinnedManager{filePath: pinnedFile},
			modified:  true,
			expected:  cleanDirs,
			argDirs:   badDirs,
		},
		{
			name:      "Save Fails",
			pinnedMgr: PinnedManager{filePath: rOnlyPath},
			modified:  false,
			expected:  cleanDirs,
			argDirs:   badDirs,
		},
		{
			name:      "Empty Input Slice",
			pinnedMgr: PinnedManager{filePath: pinnedFile},
			modified:  false,
			expected:  []directory{},
			argDirs:   []directory{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			before, err := os.Stat(pinnedFile)
			if !errors.Is(err, fs.ErrNotExist) {
				require.NoError(t, err)
			}

			cleaned := tt.pinnedMgr.Clean(tt.argDirs)

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
