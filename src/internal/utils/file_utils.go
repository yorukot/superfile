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
func LoadTomlFile(filePath string, defaultData string, target interface{}, fixFlag bool, errorPrefix string) bool {
	// Initialize with default config
	_ = toml.Unmarshal([]byte(defaultData), target)

	data, err := os.ReadFile(filePath)
	if err != nil {
		PrintfAndExit("Config file doesn't exist. Error : %v", err)
	}

	// Create a map to track which fields are present
	var rawData map[string]interface{}
	if err := toml.Unmarshal(data, &rawData); err != nil {
		fmt.Printf("%sError decoding TOML file : %v\n", errorPrefix, err)
		return true
	}

	// Replace default values with file values
	if err := toml.Unmarshal(data, target); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			row, col := decodeErr.Position()
			fmt.Printf("%sError in field at line %d column %d: %s\n",
				errorPrefix, row, col, decodeErr.Error())
		} else {
			fmt.Printf("%sError unmarshalling data: %v", errorPrefix, err)
		}
		return true
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

	// Todo: Ideally this should behave as an internal function with no side effects
	if len(missingFields) > 0 {
		if !fixFlag {
			if !ignoreMissing {
				fmt.Printf("%sMissing fields: %v\n", errorPrefix, missingFields)
				return true
			}
			// Need to return false here if we want to ignore missing fields. If this is true
			// It would cause another print message via the callee
			return false
		}
		// Create a unique backup of the current config file
		backupFile, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".bak-")
		if err != nil {
			fmt.Printf("%sFailed to create backup file. Error : %v\n", errorPrefix, err)
			return true
		}

		backupPath := backupFile.Name()
		needsBackupFileRemoval := true
		defer func() {
			backupFile.Close()
			// Remove backup in case of unsuccessful write
			if needsBackupFileRemoval {
				if errRem := os.Remove(backupPath); errRem != nil {
					fmt.Printf("%sWarning: Failed to remove backup file, backupPath : %s, err : %v\n",
						errorPrefix, backupPath, errRem)
				}
			}
		}()
		// Copy the original file to the backup
		origFile, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("%sFailed to open original file for backup. Error : %v\n", errorPrefix, err)
			return true
		}
		defer origFile.Close()

		_, err = io.Copy(backupFile, origFile)
		if err != nil {
			fmt.Printf("%sFailed to copy original file to backup. Error : %v\n", errorPrefix, err)
			return true
		}
		// Write the new config to a temp file
		tmpFile, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".tmp-")
		if err != nil {
			fmt.Printf("%sFailed to create temp file for new config. Error : %v\n", errorPrefix, err)
			return true
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
					fmt.Printf("%sWarning: Failed to remove temp config file(%s). Error : %v\n", errorPrefix, tmpPath, errRem)
				}
			}
		}()
		tomlData, err := toml.Marshal(target)
		if err != nil {
			fmt.Printf("%sFailed to marshal config to TOML. Error : %v\n", errorPrefix, err)
			return true
		}
		_, err = tmpFile.Write(tomlData)

		if err != nil {
			fmt.Printf("%sFailed to write TOML data to temp file : %v\n", errorPrefix, err)
			return true
		}

		// Even though we have a defer to close, we need to close here to make
		// sure we give up the fd to tmp file, so that os.Rename() can work.
		tmpFile.Close()
		// Atomically replace the original file with the new config
		if err := os.Rename(tmpPath, filePath); err != nil {
			fmt.Printf("%sFailed to atomically replace config file. Error : %v\n", errorPrefix, err)
			return true
		}
		// Inform user about backup location
		fmt.Printf("%sConfig file had issues. Its fixed successfully. Original backed up to: %s\n", errorPrefix, backupPath)
		// Do not remove backup; user may want to restore manually
		needsBackupFileRemoval = false
	}

	return false
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
