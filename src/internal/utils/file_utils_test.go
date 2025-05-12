package utils

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func TestResolveAbsPath(t *testing.T) {
	sep := string(filepath.Separator)
	dir1 := "abc"
	dir2 := "def"

	absPrefix := ""
	if runtime.GOOS == "windows" {
		absPrefix = "C:" // Windows absolute path prefix
	}
	root := absPrefix + sep

	testdata := []struct {
		name        string
		cwd         string
		path        string
		expectedRes string
	}{
		{
			name:        "Path cleaup Test 1",
			cwd:         absPrefix + sep,
			path:        absPrefix + strings.Repeat(sep, 10),
			expectedRes: absPrefix + sep,
		},
		{
			name:        "Basic test",
			cwd:         filepath.Join(root, dir1),
			path:        dir2,
			expectedRes: filepath.Join(root, dir1, dir2),
		},
		{
			name:        "Ignore cwd for abs path",
			cwd:         filepath.Join(root, dir1),
			path:        filepath.Join(root, dir2),
			expectedRes: filepath.Join(root, dir2),
		},
		{
			name:        "Path cleanup Test 2",
			cwd:         absPrefix + strings.Repeat(sep, 4) + dir1,
			path:        "." + sep + "." + sep + dir2,
			expectedRes: filepath.Join(root, dir1, dir2),
		},
		{
			name:        "Basic test with ~",
			cwd:         root,
			path:        "~",
			expectedRes: xdg.Home,
		},
		{
			name:        "~ should not be resolved if not first",
			cwd:         dir1,
			path:        filepath.Join(dir2, "~"),
			expectedRes: filepath.Join(dir1, dir2, "~"),
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedRes, ResolveAbsPath(tt.cwd, tt.path))
		})
	}
}
