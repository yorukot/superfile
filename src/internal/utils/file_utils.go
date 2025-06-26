package utils

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml/v2"
)

// Utility functions related to file operations

func WriteTomlData(filePath string, data interface{}) error {
	tomlData, err := toml.Marshal(data)
	if err != nil {
		// return a wrapped error
		return fmt.Errorf("error encoding data : %w", err)
	}
	err = os.WriteFile(filePath, tomlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file : %w", err)
	}
	return nil
}

// Helper function to load and validate TOML files with field checking
// errorPrefix is appended before every error message
func LoadTomlFile(filePath string, defaultData string, target interface{}, fixFlag bool) error {
	// Initialize with default config
	_ = toml.Unmarshal([]byte(defaultData), target)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return &TomlLoadError{
			userMessage:  "config file doesn't exist",
			wrappedError: err,
		}
	}

	// Create a map to track which fields are present
	var rawData map[string]interface{}
	err = toml.Unmarshal(data, &rawData)
	if err != nil {
		return &TomlLoadError{
			userMessage:  "error decoding TOML file",
			wrappedError: err,
			isFatal:      true,
		}
	}

	// Replace default values with file values
	err = toml.Unmarshal(data, target)
	if err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			row, col := decodeErr.Position()
			return &TomlLoadError{
				userMessage:  fmt.Sprintf("error in field at line %d column %d", row, col),
				wrappedError: decodeErr,
				isFatal:      true,
			}
		}
		return &TomlLoadError{
			userMessage:  "error unmarshalling data",
			wrappedError: err,
			isFatal:      true,
		}
	}

	// Check for missing fields. Explicitly set default value to false
	ignoreMissing := false
	if config, ok := target.(MissingFieldIgnorer); ok {
		ignoreMissing = config.GetIgnoreMissingFields()
	}

	// Check for missing fields
	targetType := reflect.TypeOf(target).Elem()
	missingFields := []string{}

	for i := range targetType.NumField() {
		field := targetType.Field(i)
		var fieldName string
		tag := field.Tag.Get("toml")
		if tag != "" {
			// Discard options such as ",omitempty"
			fieldName = strings.Split(tag, ",")[0]
		} else {
			fieldName = field.Name
		}
		if _, exists := rawData[fieldName]; !exists {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) == 0 {
		return nil
	}
	if !fixFlag && ignoreMissing {
		// nil error if we dont wanna fix, and dont wanna print
		return nil
	}

	resultErr := &TomlLoadError{
		missingFields: true,
	}
	if !fixFlag {
		resultErr.userMessage = fmt.Sprintf("missing fields: %v", missingFields)
		return resultErr
	}
	// Create a unique backup of the current config file
	backupFile, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".bak-")
	if err != nil {
		resultErr.UpdateMessageAndError("failed to create backup file", err)
		return resultErr
	}

	backupPath := backupFile.Name()
	needsBackupFileRemoval := true
	defer func() {
		backupFile.Close()
		// Remove backup in case of unsuccessful write
		if needsBackupFileRemoval {
			if errRem := os.Remove(backupPath); errRem != nil {
				// Modify result Error
				resultErr.AddMessageAndError("warning: failed to remove backup file, backupPath : "+backupPath, errRem)
			}
		}
	}()
	// Copy the original file to the backup
	origFile, err := os.Open(filePath)
	if err != nil {
		resultErr.UpdateMessageAndError("failed to open original file for backup", err)
		return resultErr
	}
	defer origFile.Close()

	_, err = io.Copy(backupFile, origFile)
	if err != nil {
		resultErr.UpdateMessageAndError("failed to copy original file to backup", err)
		return resultErr
	}
	// Write the new config to a temp file
	tmpFile, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".tmp-")
	if err != nil {
		resultErr.UpdateMessageAndError("failed to create temp file for new config", err)
		return resultErr
	}
	tmpPath := tmpFile.Name()
	// Ensure cleanup via defer
	defer func() {
		// In usual case, the close would have already happened. Ignore the error here
		tmpFile.Close()
		// Cleanup file if it exists
		if _, errRem := os.Stat(tmpPath); errRem == nil {
			// File still exists
			if errRem := os.Remove(tmpPath); errRem != nil {
				resultErr.AddMessageAndError(
					fmt.Sprintf("warning: failed to remove temp config file(%s)", tmpPath), errRem)
			}
		}
	}()
	tomlData, err := toml.Marshal(target)
	if err != nil {
		resultErr.UpdateMessageAndError("failed to marshal config to TOML", err)
		return resultErr
	}
	_, err = tmpFile.Write(tomlData)

	if err != nil {
		resultErr.UpdateMessageAndError("failed to write TOML data to temp file", err)
		return resultErr
	}

	// Even though we have a defer to close, we need to close here to make
	// sure we give up the fd to tmp file, so that os.Rename() can work.
	tmpFile.Close()
	// Atomically replace the original file with the new config
	if err := os.Rename(tmpPath, filePath); err != nil {
		resultErr.UpdateMessageAndError("failed to atomically replace config file", err)
		return resultErr
	}
	// Inform user about backup location
	resultErr.userMessage = "config file had issues. Its fixed successfully. Original backed up to : " + backupPath
	// Do not remove backup; user may want to restore manually
	needsBackupFileRemoval = false

	return resultErr
}

// If path is not absolute, then append to currentDir and get absolute path
// Resolve paths starting with "~"
// currentDir should be an absolute path
func ResolveAbsPath(currentDir string, path string) string {
	if !filepath.IsAbs(currentDir) {
		slog.Warn("currentDir is not absolute", "currentDir", currentDir)
	}
	if strings.HasPrefix(path, "~") {
		// We dont use variable.HomeDir here, as the util package cannot have dependency
		// on variable package
		path = strings.Replace(path, "~", xdg.Home, 1)
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(currentDir, path)
	}
	return filepath.Clean(path)
}
