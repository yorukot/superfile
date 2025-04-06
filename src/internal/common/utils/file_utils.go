package utils

import (
	"fmt"
	"os"
	"reflect"

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
		LogAndExit("Config file doesn't exist", "error", err)
	}
	errMsg := ""
	hasError := false

	// Create a map to track which fields are present
	// Have to do this manually as toml.Unmarshal does not return an error when it encounters a TOML key
	// that does not match any field in the struct.
	// Instead, it simply ignores that key and continues parsing.
	var rawData map[string]interface{}
	rawError := toml.Unmarshal(data, &rawData)
	if rawError != nil {
		hasError = true
		errMsg = errorPrefix + "Error decoding file: " + rawError.Error() + "\n"
	}

	// Replace default values with file values
	if !hasError {
		if err = toml.Unmarshal(data, target); err != nil {
			hasError = true
			//nolint: errorlint // Type assertion is better here, and we need to read data from error
			if decodeErr, ok := err.(*toml.DecodeError); ok {
				row, col := decodeErr.Position()
				errMsg = errorPrefix + fmt.Sprintf("Error in field at line %d column %d: %s\n",
					row, col, decodeErr.Error())
			} else {
				errMsg = errorPrefix + "Error unmarshalling data: " + err.Error() + "\n"
			}
		}
	}

	if !hasError {
		// Check for missing fields if no errors yet
		targetType := reflect.TypeOf(target).Elem()

		for i := range targetType.NumField() {
			field := targetType.Field(i)
			if _, exists := rawData[field.Tag.Get("toml")]; !exists {
				hasError = true
				// A field doesn't exist in the toml config file
				errMsg += errorPrefix + fmt.Sprintf("Field \"%s\" is missing\n", field.Tag.Get("toml"))
			}
		}
	}
	// File is okay
	if !hasError {
		return false
	}

	// File is bad, but we arent' allowed to fix
	// We just print error message to stdout
	// Todo : Ideally this should behave as an intenal function with no side effects
	// and the caller should do the printing to stdout
	if !fixFlag {
		fmt.Print(errMsg)
		return true
	}
	// Now we are fixing the file, we would not return hasError=true even if there was error
	// Fix the file by writing all fields
	if err := WriteTomlData(filePath, target); err != nil {
		LogAndExit("Error while writing config file", "error", err)
	}
	return false
}
