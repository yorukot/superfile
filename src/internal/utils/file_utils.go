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
		slog.Error(errorPrefix+"Error decoding TOML file", "err", err)
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
		// Create a unique backup of the current config file
		backupFile, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".bak-")
		if err != nil {
			slog.Error(errorPrefix+"Failed to create backup file", "err", err)
			return true
		}
		backupPath := backupFile.Name()
		// Copy the original file to the backup
		origFile, err := os.Open(filePath)
		if err != nil {
			slog.Error(errorPrefix+"Failed to open original file for backup", "err", err)
			backupFile.Close()
			os.Remove(backupPath)
			return true
		}
		_, err = io.Copy(backupFile, origFile)
		origFile.Close()
		backupFile.Close()
		if err != nil {
			slog.Error(errorPrefix+"Failed to copy original file to backup", "err", err)
			os.Remove(backupPath)
			return true
		}
		// Write the new config to a temp file
		tmpFile, err := os.CreateTemp(filepath.Dir(filePath), filepath.Base(filePath)+".tmp-")
		if err != nil {
			slog.Error(errorPrefix+"Failed to create temp file for new config", "err", err)
			return true
		}
		tmpPath := tmpFile.Name()
		tomlData, err := toml.Marshal(target)
		if err != nil {
			slog.Error(errorPrefix+"Failed to marshal config to TOML", "err", err)
			tmpFile.Close()
			os.Remove(tmpPath)
			return true
		}
		_, err = tmpFile.Write(tomlData)
		tmpFile.Close()
		if err != nil {
			slog.Error(errorPrefix+"Failed to write TOML data to temp file", "err", err)
			os.Remove(tmpPath)
			return true
		}
		// Atomically replace the original file with the new config
		if err := os.Rename(tmpPath, filePath); err != nil {
			slog.Error(errorPrefix+"Failed to atomically replace config file", "err", err)
			// Do not remove backup; user may want to restore manually
			return true
		}
		// Remove backup after successful write
		if err := os.Remove(backupPath); err != nil {
			slog.Error(errorPrefix+"Warning: Failed to remove backup file", "backupPath", backupPath, "err", err)
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
