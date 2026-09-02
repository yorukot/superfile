package filepanel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
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
		isDir           bool
	}{
		{
			name:            "Empty dir",
			entries:         []string{},
			includeDotFiles: false,
			expectedCount:   0,
			isDir:           true,
		},
		{
			name:            "Dir with 3 files",
			entries:         []string{"file1.txt", "file2.txt", "file3.txt"},
			includeDotFiles: false,
			expectedCount:   3,
			isDir:           true,
		},
		{
			name:            "Dir with 2 non-dot files and 1 dot file, includeDotFiles false",
			entries:         []string{"file1.txt", "file2.txt", ".file3.txt"},
			includeDotFiles: false,
			expectedCount:   2,
			isDir:           true,
		},
		{
			name:            "Dir with 2 non-dot files and 1 dot file, includeDotFiles true",
			entries:         []string{"file1.txt", "file2.txt", ".file3.txt"},
			includeDotFiles: true,
			expectedCount:   3,
			isDir:           true,
		},
		{
			name:            "Dir with 3 dot files, includeDotFiles false",
			entries:         []string{".file1.txt", ".file2.txt", ".file3.txt"},
			includeDotFiles: false,
			expectedCount:   0,
			isDir:           true,
		},
		{
			name:            "Dir with 3 dot files, includeDotFiles true",
			entries:         []string{".file1.txt", ".file2.txt", ".file3.txt"},
			includeDotFiles: true,
			expectedCount:   3,
			isDir:           true,
		},
		{
			name:            "Non-dir returns 0",
			entries:         []string{},
			includeDotFiles: true,
			expectedCount:   0,
			isDir:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			for _, name := range tt.entries {
				path := filepath.Join(dir, name)
				require.NoError(t, os.WriteFile(path, []byte{}, 0644))
			}

			location := dir
			if !tt.isDir {
				location = filepath.Join(dir, "file.txt")
				require.NoError(t, os.WriteFile(location, []byte{}, 0644))
			}

			info, err := os.Stat(location)
			require.NoError(t, err)

			element := Element{
				Location: location,
				Info:     info,
			}

			assert.Equal(
				t,
				tt.expectedCount,
				element.GetChildCount(tt.includeDotFiles),
			)
		})
	}
}
