package utils

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"os"

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

func TestLoadTomlFile_MissingFieldsWithOmitEmpty(t *testing.T) {
	type testConfig struct {
		FieldA string `toml:"field_a,omitempty"`
		FieldB string `toml:"field_b"`
	}

	defaultData := `field_a = "A"
field_b = "B"`

	// Only field_b present in file
	tomlData := `field_b = "B"`

	tmpFile, err := os.CreateTemp("", "testconfig-*.toml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write([]byte(tomlData))
	tmpFile.Close()
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	var cfg testConfig
	missing := LoadTomlFile(tmpFile.Name(), defaultData, &cfg, false, "test: ")
	if !missing {
		t.Errorf("expected missing fields to be detected when field_a is omitted")
	}
}
