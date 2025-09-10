package utils

import (
	"fmt"
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
			err = LoadTomlFile(filepath.Join(testdataDir, tt.configName), defaultData, &tomlVal,
				tt.fixFlag, false)
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
}

func TestLoadTomlFileIgnorer(t *testing.T) {
	_, curFilename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	testdataDir := filepath.Join(filepath.Dir(curFilename), "testdata", "load_toml", "ignorer")

	defaultDataBytes, err := os.ReadFile(filepath.Join(testdataDir, "default.toml"))
	require.NoError(t, err)

	defaultData := string(defaultDataBytes)
	var defaultTomlVal TestTOMLMissingIgnorerType
	err = toml.Unmarshal(defaultDataBytes, &defaultTomlVal)
	require.NoError(t, err)
	// This is for Ignorer Type
	testdata := []struct {
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
			configName:      "default.toml",
			fixFlag:         false,
			noError:         true,
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal,
		},
		{
			name:            "Config2 Extra Fields ignored",
			configName:      "default_extra_fields.toml",
			fixFlag:         false,
			noError:         true,
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal,
		},
		{
			name:       "Config2 Missing fields Not Ignored",
			configName: "missing_str_int.toml",
			fixFlag:    false,
			noError:    false,
			expectedError: &TomlLoadError{
				userMessage:   "missing fields: [sample_int sample_str]",
				wrappedError:  nil,
				isFatal:       false,
				missingFields: true,
			},
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal,
		},
		{
			name:            "Config2 Missing fields Ignored",
			configName:      "missing_str_ignore.toml",
			fixFlag:         false,
			noError:         true,
			checkTomlVal:    true,
			expectedTomlVal: defaultTomlVal.WithIgnoreMissing(true),
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
		{
			name:       "Config2 Invalid format",
			configName: "invalid_format.toml",
			fixFlag:    false,
			noError:    false,
			expectedError: &TomlLoadError{
				userMessage: "error decoding TOML file",
				isFatal:     true,
			},
			verifyWrappedErr: false,
			checkTomlVal:     false,
		},
		{
			name:       "Config2 Invalid Value Type",
			configName: "invalid_value_type.toml",
			fixFlag:    false,
			noError:    false,
			expectedError: &TomlLoadError{
				userMessage: "error in field at line 2 column 14",
				isFatal:     true,
			},
			verifyWrappedErr: false,
			checkTomlVal:     false,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			var tomlVal TestTOMLMissingIgnorerType
			err := LoadTomlFile(filepath.Join(testdataDir, tt.configName), defaultData, &tomlVal,
				tt.fixFlag, false)
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

	// Tests for fixing config file

	t.Run("Config2 Fixing config file", func(t *testing.T) {
		// To make sure that other values are kept.
		expectedVal1 := defaultTomlVal
		expectedVal2 := defaultTomlVal
		expectedVal1.SampleInt = -1
		expectedVal2.SampleInt = -1
		expectedVal2.IgnoreMissing = true

		tempDir := t.TempDir()

		actualTest := func(fileName string, expectedVal TestTOMLMissingIgnorerType) {
			var tomlVal TestTOMLMissingIgnorerType
			testFile := filepath.Join(testdataDir, fileName)
			orgFile := filepath.Join(tempDir, fileName)

			testContent, err := os.ReadFile(testFile)
			require.NoError(t, err)

			// Copy to temp directory first to avoid permission errors
			err = os.WriteFile(orgFile, testContent, 0644)
			require.NoError(t, err, "Error writing config file to temp directory")

			err = LoadTomlFile(orgFile, defaultData, &tomlVal, true, false)
			var tomlErr *TomlLoadError
			require.ErrorAs(t, err, &tomlErr)

			assert.True(t, tomlErr.missingFields)
			assert.Equal(t, expectedVal, tomlVal)

			pref := "config file had issues. Its fixed successfully. Original backed up to : "

			assert.True(t, strings.HasPrefix(tomlErr.userMessage, pref), "Unexpectd error : "+tomlErr.Error())

			backupFile := strings.TrimPrefix(tomlErr.userMessage, pref)

			assert.FileExists(t, backupFile)
			backupContent, err := os.ReadFile(backupFile)
			require.NoError(t, err)

			assert.Equal(t, testContent, backupContent)

			// Validate that if you Load Original File again, it loads without any errors
			err = LoadTomlFile(orgFile, defaultData, &tomlVal, true, false)
			require.NoError(t, err)

			err = os.WriteFile(orgFile, backupContent, 0644)
			require.NoError(t, err)
		}
		actualTest("missing_str2.toml", expectedVal1)
		actualTest("missing_str_ignore2.toml", expectedVal2)
	})
}

func TestReadFileContent(t *testing.T) {
	testDir := t.TempDir()
	curTestDir := filepath.Join(testDir, "TestReadFileContent")
	SetupDirectories(t, curTestDir)

	testdata := []struct {
		name          string
		content       []byte
		maxLineLength int
		previewLine   int
		expected      string
	}{
		{
			name:          "regular UTF-8 file",
			content:       []byte("line1\nline2\nline3"),
			maxLineLength: 100,
			previewLine:   5,
			expected:      "line1\nline2\nline3\n",
		},
		{
			name:          "UTF-8 BOM file",
			content:       []byte("\xEF\xBB\xBFline1\nline2\nline3"),
			maxLineLength: 100,
			previewLine:   5,
			expected:      "line1\nline2\nline3\n",
		},
		{
			name:          "limited preview lines",
			content:       []byte("line1\nline2\nline3\nline4"),
			maxLineLength: 100,
			previewLine:   2,
			expected:      "line1\nline2\n",
		},
	}

	for i, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(curTestDir, fmt.Sprintf("test_file_%d.txt", i))
			SetupFilesWithData(t, tt.content, testFile)

			result, err := ReadFileContent(testFile, tt.maxLineLength, tt.previewLine)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReadFileContentBOMHandling(t *testing.T) {
	testDir := t.TempDir()
	curTestDir := filepath.Join(testDir, "TestBOMHandling")
	SetupDirectories(t, curTestDir)

	// Write a file prefixed with UTF-8 BOM
	bomContent := []byte("\xEF\xBB\xBFHello, World!\nSecond line")
	bomFile := filepath.Join(curTestDir, "bom_file.txt")
	SetupFilesWithData(t, bomContent, bomFile)

	result, err := ReadFileContent(bomFile, 100, 10)
	require.NoError(t, err)

	// Verify BOM is removed and content is correct
	assert.True(t, strings.HasPrefix(result, "Hello, World!"),
		"Content should start with expected text, got: %q", result)
	assert.NotContains(t, result, "\uFEFF",
		"BOM character should be removed from output: %q", result)
}
