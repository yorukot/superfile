package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// We cannot use ConfigType here, as that is not accessible by "utils" package
func TestLoadTomlFile(t *testing.T) {
	_, curFilename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	testdataDir := filepath.Join(filepath.Dir(curFilename), "testdata", "load_toml")

	defaultDataBytes, err := os.ReadFile(filepath.Join(testdataDir, "default.toml"))
	require.NoError(t, err)

	defaultData := string(defaultDataBytes)
	var defaultTomlVal TestTOMLType
	err = toml.Unmarshal(defaultDataBytes, &defaultTomlVal)
	require.NoError(t, err)

	testdata := []struct {
		name string
		// Relative to corr
		configName string
		fixFlag    bool
		noError    bool

		// If we have error. It should be TomlLoadError
		expectedError *TomlLoadError

		// For checking the result value
		checkTomlVal    bool
		expectedTomlVal TestTOMLType
	}{
		{
			name:            "Config1 Load Default",
			configName:      "default.toml",
			fixFlag:         false,
			noError:         true,
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal,
		},
		{
			name:       "Config1 Missing fields",
			configName: "missing_str.toml",
			fixFlag:    false,
			noError:    false,
			expectedError: &TomlLoadError{
				userMessage:   "missing fields: [sample_str]",
				wrappedError:  nil,
				isFatal:       false,
				missingFields: true,
			},
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			var tomlVal TestTOMLType
			err = LoadTomlFile(filepath.Join(testdataDir, tt.configName), defaultData, &tomlVal, tt.fixFlag)
			if tt.noError {
				require.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}

			if tt.checkTomlVal {
				assert.Equal(t, tt.expectedTomlVal, tomlVal)
			}
		})
	}

	defaultDataBytes2, err := os.ReadFile(filepath.Join(testdataDir, "default2.toml"))
	require.NoError(t, err)

	defaultData2 := string(defaultDataBytes)
	var defaultTomlVal2 TestTOMLMissingIgnorerType
	err = toml.Unmarshal(defaultDataBytes2, &defaultTomlVal2)
	require.NoError(t, err)
	// This is for Ignorer Type
	testdata2 := []struct {
		name string
		// Relative to corr
		configName string
		fixFlag    bool
		noError    bool

		// If we have error. It should be TomlLoadError
		expectedError    *TomlLoadError
		verifyWrappedErr bool
		// For checking the result value
		checkTomlVal    bool
		expectedTomlVal TestTOMLMissingIgnorerType
	}{
		{
			name:            "Config2 Load Default",
			configName:      "default2.toml",
			fixFlag:         false,
			noError:         true,
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal2,
		},
		{
			name:       "Config2 Missing fields Not Ignored",
			configName: "missing_str_int2.toml",
			fixFlag:    false,
			noError:    false,
			expectedError: &TomlLoadError{
				userMessage:   "missing fields: [sample_int sample_str]",
				wrappedError:  nil,
				isFatal:       false,
				missingFields: true,
			},
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal2,
		},
		{
			name:            "Config2 Missing fields Ignored",
			configName:      "missing_str2_ignore.toml",
			fixFlag:         false,
			noError:         true,
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal2.WithIgnoreMissing(true),
		},
		{
			name:       "Config2 Non Existent config",
			configName: "non_existent_config.toml",
			fixFlag:    false,
			noError:    false,
			expectedError: &TomlLoadError{
				userMessage: "config file doesn't exist",
			},
			verifyWrappedErr: false,
			checkTomlVal:     false,
		},
	}

	for _, tt := range testdata2 {
		t.Run(tt.name, func(t *testing.T) {
			var tomlVal TestTOMLMissingIgnorerType
			err := LoadTomlFile(filepath.Join(testdataDir, tt.configName), defaultData2, &tomlVal, tt.fixFlag)
			if tt.noError {
				require.NoError(t, err)
			} else {
				var tomlErr *TomlLoadError
				require.ErrorAs(t, err, &tomlErr)
				if tt.verifyWrappedErr {
					assert.Equal(t, tt.expectedError, tomlErr)
				} else {
					assert.Equal(t, tt.expectedError.userMessage, tomlErr.userMessage)
					assert.Equal(t, tt.expectedError.isFatal, tomlErr.isFatal)
					assert.Equal(t, tt.expectedError.missingFields, tomlErr.missingFields)
				}
			}

			if tt.checkTomlVal {
				assert.Equal(t, tt.expectedTomlVal, tomlVal)
			}
		})
	}
}
