package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var SampleDataBytes = []byte("This is sample") //nolint: gochecknoglobals // Effectively const

type TestTOMLType struct {
	SampleBool  bool     `toml:"sample_bool"`
	SampleInt   int      `toml:"sample_int"`
	SampleStr   string   `toml:"sample_str"`
	SampleSlice []string `toml:"sample_slice"`
}

type TestTOMLMissingIgnorerType struct {
	SampleBool    bool     `toml:"sample_bool"`
	SampleInt     int      `toml:"sample_int"`
	SampleStr     string   `toml:"sample_str"`
	SampleSlice   []string `toml:"sample_slice"`
	IgnoreMissing bool     `toml:"ignore_missing"`
}

func (t TestTOMLMissingIgnorerType) GetIgnoreMissingFields() bool {
	return t.IgnoreMissing
}

func (t TestTOMLMissingIgnorerType) WithIgnoreMissing(val bool) TestTOMLMissingIgnorerType {
	t.IgnoreMissing = val
	return t
}

func SetupDirectories(t *testing.T, dirs ...string) {
	t.Helper()
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		err := os.MkdirAll(dir, UserDirPerm)
		require.NoError(t, err)
	}
}

func SetupFilesWithData(t *testing.T, data []byte, files ...string) {
	t.Helper()
	for _, file := range files {
		err := os.WriteFile(file, data, UserFilePerm)
		require.NoError(t, err)
	}
}

func SetupFiles(t *testing.T, files ...string) {
	SetupFilesWithData(t, SampleDataBytes, files...)
}
