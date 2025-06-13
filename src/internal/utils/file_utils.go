package utils

import (
	"errors"
	"fmt"
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
		fmt.Print(errorPrefix + "Error decoding file: " + err.Error() + "\n")
		return true
	}

	// Replace default values with file values
	if err := toml.Unmarshal(data, target); err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			row, col := decodeErr.Position()
			fmt.Print(errorPrefix + fmt.Sprintf("Error in field at line %d column %d: %s\n",
				row, col, decodeErr.Error()))
		} else {
			fmt.Print(errorPrefix + "Error unmarshalling data: " + err.Error() + "\n")
		}
		return true
	}

	// Check for missing fields
	var ignoreMissing bool
	if config, ok := target.(ConfigInterface); ok {
		ignoreMissing = config.GetIgnoreMissingFields()
	}

	// Check for missing fields
	targetType := reflect.TypeOf(target).Elem()
	missingFields := []string{}

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldName := field.Tag.Get("toml")
		if fieldName == "" {
			fieldName = field.Name
		}
		if _, exists := rawData[fieldName]; !exists && !ignoreMissing {
			missingFields = append(missingFields, fieldName)
		}
	}

	// Todo: Ideally this should behave as an internal function with no side effects
	if len(missingFields) > 0 {
		if !fixFlag {
			fmt.Print(errorPrefix + fmt.Sprintf("Missing fields: %v\n", missingFields))
			return true
		}
		// Create a backup of the current config file
		backupPath := filePath + ".bak"
		if err := os.Rename(filePath, backupPath); err != nil {
			fmt.Print(errorPrefix + fmt.Sprintf("Failed to create backup of config file: %v\n", err))
			return true
		}
		// Fix the file by writing all fields
		if err := WriteTomlData(filePath, target); err != nil {
			// Restore from backup if write fails
			if restoreErr := os.Rename(backupPath, filePath); restoreErr != nil {
				fmt.Print(errorPrefix + fmt.Sprintf("Failed to restore config from backup: %v\n", restoreErr))
			}
			PrintfAndExit("Error while writing config file : %v", err)
		}
		// Remove backup after successful write
		if err := os.Remove(backupPath); err != nil {
			fmt.Print(errorPrefix + fmt.Sprintf("Warning: Failed to remove backup file %s: %v\n", backupPath, err))
		}
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
