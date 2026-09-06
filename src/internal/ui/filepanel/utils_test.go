package filepanel

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestGetSelectedLocationsSortedAsVisible(t *testing.T) {
	testdata := []struct {
		name             string
		panel            Model
		toSelect         []string
		expectedSelected []string
	}{
		{
			name: "no any selected",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
				{Name: "file4.txt", Location: "/tmp/file4.txt"},
			}),
			expectedSelected: []string{},
		},
		{
			name: "1 item selected",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
				{Name: "file4.txt", Location: "/tmp/file4.txt"},
			}),
			toSelect:         []string{"/tmp/file2.txt"},
			expectedSelected: []string{"/tmp/file2.txt"},
		},
		{
			name: "2 item selects reverse selection order",
			panel: testModel(-1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file4.txt", Location: "/tmp/file3.txt"},
				{Name: "file5.txt", Location: "/tmp/file4.txt"},
			}),
			toSelect:         []string{"/tmp/file4.txt", "/tmp/file2.txt"},
			expectedSelected: []string{"/tmp/file2.txt", "/tmp/file4.txt"},
		},
		{
			name: "2 item selects",
			panel: testModel(-1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
				{Name: "file4.txt", Location: "/tmp/file4.txt"},
			}),
			toSelect:         []string{"/tmp/file2.txt", "/tmp/file4.txt"},
			expectedSelected: []string{"/tmp/file2.txt", "/tmp/file4.txt"},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SortKind = sortmodel.SortByName
			tt.panel.SetSelectedAll(tt.toSelect)
			assert.Equal(t, tt.expectedSelected, tt.panel.GetSelectedLocationsSortedAsVisible())
		})
	}
}

func TestGetChildCount(t *testing.T) {
	tests := []struct {
		name            string
		entries         []string
		includeDotFiles bool
		expectedCount   int
	}{
		{
			name:            "Empty dir",
			entries:         []string{},
			includeDotFiles: false,
			expectedCount:   0,
		},
		{
			name:            "Dir with 3 files",
			entries:         []string{"file1.txt", "file2.txt", "file3.txt"},
			includeDotFiles: false,
			expectedCount:   3,
		},
		{
			name:            "Dir with 2 non-dot files and 1 dot file, includeDotFiles false",
			entries:         []string{"file1.txt", "file2.txt", ".file3.txt"},
			includeDotFiles: false,
			expectedCount:   2,
		},
		{
			name:            "Dir with 2 non-dot files and 1 dot file, includeDotFiles true",
			entries:         []string{"file1.txt", "file2.txt", ".file3.txt"},
			includeDotFiles: true,
			expectedCount:   3,
		},
		{
			name:            "Dir with 3 dot files, includeDotFiles false",
			entries:         []string{".file1.txt", ".file2.txt", ".file3.txt"},
			includeDotFiles: false,
			expectedCount:   0,
		},
		{
			name:            "Dir with 3 dot files, includeDotFiles true",
			entries:         []string{".file1.txt", ".file2.txt", ".file3.txt"},
			includeDotFiles: true,
			expectedCount:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			files := make([]string, 0, len(tt.entries))
			for _, name := range tt.entries {
				files = append(files, filepath.Join(dir, name))
			}
			utils.SetupFiles(t, files...)

			count, err := getChildCount(dir, tt.includeDotFiles)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, count)
		})
	}
}

func TestGetChildCountReturnsReadError(t *testing.T) {
	count, err := getChildCount(filepath.Join(t.TempDir(), "missing"), true)

	assert.Zero(t, count)
	require.Error(t, err)
}

func TestDirectorySymlinkChildCount(t *testing.T) {
	parent := t.TempDir()
	target := t.TempDir()
	utils.SetupFiles(t,
		filepath.Join(target, "visible.txt"),
		filepath.Join(target, ".hidden.txt"),
	)

	symlink := filepath.Join(parent, "directory-link")
	if err := os.Symlink(target, symlink); err != nil {
		t.Skipf("directory symlinks are unavailable: %v", err)
	}
	dirEntries, err := os.ReadDir(parent)
	require.NoError(t, err)
	require.Len(t, dirEntries, 1)

	for _, tt := range []struct {
		name            string
		includeDotFiles bool
		expectedCount   int
	}{
		{name: "Exclude dotfiles", includeDotFiles: false, expectedCount: 1},
		{name: "Include dotfiles", includeDotFiles: true, expectedCount: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			panel := testModel(0, 0, 12, BrowserMode, nil)
			panel.Location = parent
			panel.SortKind = sortmodel.SortByName
			elements := panel.sortFileElements(dirEntries, tt.includeDotFiles)
			require.Len(t, elements, 1)
			assert.True(t, elements[0].Directory)
			assert.False(t, elements[0].Info.IsDir())
			assert.Equal(t, tt.expectedCount, elements[0].ChildCount)
			assert.NoError(t, elements[0].ChildCountErr)
		})
	}
}

func TestRenderFileSizeUsesPopulatedDirectoryCount(t *testing.T) {
	dir := t.TempDir()
	info, err := os.Stat(dir)
	require.NoError(t, err)

	tests := []struct {
		name     string
		element  Element
		expected string
	}{
		{
			name: "Singular count",
			element: Element{
				Directory:  true,
				Info:       info,
				ChildCount: 1,
			},
			expected: "1 item",
		},
		{
			name: "Plural count",
			element: Element{
				Directory:  true,
				Info:       info,
				ChildCount: 2,
			},
			expected: "2 items",
		},
		{
			name: "Read error",
			element: Element{
				Directory:     true,
				Info:          info,
				ChildCountErr: errors.New("read failed"),
			},
			expected: "(Error)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panel := testModel(0, 0, 12, BrowserMode, []Element{tt.element})
			assert.Contains(t, panel.renderFileSize(0, FileSizeColumnWidth), tt.expected)
		})
	}
}
