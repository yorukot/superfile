package common

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var ErrInvalidOpenWithApplication = errors.New("invalid open_with application")

// OpenWithApplication describes one application available in the Open With menu.
// Command is executed directly with Args followed by the selected file path.
type OpenWithApplication struct {
	Name     string
	Command  string
	Args     []string
	Terminal bool
}

// ResolveOpenWith returns the applications configured for the file's extension.
func ResolveOpenWith(config map[string]any, filePath string) ([]OpenWithApplication, error) {
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")
	if extension == "" {
		return nil, nil
	}

	value, found := config[extension]
	if !found {
		for configuredExtension, configuredValue := range config {
			if strings.EqualFold(configuredExtension, extension) {
				value = configuredValue
				found = true
				break
			}
		}
	}
	if !found {
		return nil, nil
	}

	return parseOpenWithValue(value)
}

// ValidateOpenWith checks every configured application without requiring a file.
func ValidateOpenWith(config map[string]any) error {
	for extension, value := range config {
		if strings.TrimSpace(extension) == "" {
			return fmt.Errorf("%w: extension cannot be empty", ErrInvalidOpenWithApplication)
		}
		if _, err := parseOpenWithValue(value); err != nil {
			return fmt.Errorf("open_with.%s: %w", extension, err)
		}
	}
	return nil
}

func parseOpenWithValue(value any) ([]OpenWithApplication, error) {
	switch applications := value.(type) {
	case string:
		application, err := parseOpenWithApplication(applications)
		if err != nil {
			return nil, err
		}
		return []OpenWithApplication{application}, nil
	case []string:
		result := make([]OpenWithApplication, 0, len(applications))
		for _, applicationValue := range applications {
			application, err := parseOpenWithApplication(applicationValue)
			if err != nil {
				return nil, err
			}
			result = append(result, application)
		}
		return result, nil
	case []any:
		result := make([]OpenWithApplication, 0, len(applications))
		for _, applicationValue := range applications {
			application, err := parseOpenWithApplication(applicationValue)
			if err != nil {
				return nil, err
			}
			result = append(result, application)
		}
		return result, nil
	default:
		application, err := parseOpenWithApplication(applications)
		if err != nil {
			return nil, err
		}
		return []OpenWithApplication{application}, nil
	}
}

func parseOpenWithApplication(value any) (OpenWithApplication, error) {
	switch application := value.(type) {
	case string:
		command := strings.TrimSpace(application)
		if command == "" {
			return OpenWithApplication{}, fmt.Errorf("%w: command cannot be empty", ErrInvalidOpenWithApplication)
		}
		return OpenWithApplication{Name: filepath.Base(command), Command: command}, nil
	case OpenWithApplication:
		return normalizeOpenWithApplication(application)
	case map[string]any:
		return parseOpenWithApplicationMap(application)
	default:
		return OpenWithApplication{}, fmt.Errorf("%w: expected a string or application table", ErrInvalidOpenWithApplication)
	}
}

func parseOpenWithApplicationMap(value map[string]any) (OpenWithApplication, error) {
	application := OpenWithApplication{}
	if name, ok := value["name"].(string); ok {
		application.Name = name
	} else if _, exists := value["name"]; exists {
		return OpenWithApplication{}, fmt.Errorf("%w: name must be a string", ErrInvalidOpenWithApplication)
	}
	if command, ok := value["command"].(string); ok {
		application.Command = command
	} else {
		return OpenWithApplication{}, fmt.Errorf("%w: command must be a string", ErrInvalidOpenWithApplication)
	}
	if terminal, ok := value["terminal"].(bool); ok {
		application.Terminal = terminal
	} else if _, exists := value["terminal"]; exists {
		return OpenWithApplication{}, fmt.Errorf("%w: terminal must be a boolean", ErrInvalidOpenWithApplication)
	}
	if rawArgs, exists := value["args"]; exists {
		args, err := parseOpenWithArgs(rawArgs)
		if err != nil {
			return OpenWithApplication{}, err
		}
		application.Args = args
	}
	return normalizeOpenWithApplication(application)
}

func parseOpenWithArgs(value any) ([]string, error) {
	switch args := value.(type) {
	case []string:
		return append([]string(nil), args...), nil
	case []any:
		result := make([]string, 0, len(args))
		for _, arg := range args {
			stringArg, ok := arg.(string)
			if !ok {
				return nil, fmt.Errorf("%w: args must contain only strings", ErrInvalidOpenWithApplication)
			}
			result = append(result, stringArg)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("%w: args must be a list of strings", ErrInvalidOpenWithApplication)
	}
}

func normalizeOpenWithApplication(application OpenWithApplication) (OpenWithApplication, error) {
	application.Command = strings.TrimSpace(application.Command)
	if application.Command == "" {
		return OpenWithApplication{}, fmt.Errorf("%w: command cannot be empty", ErrInvalidOpenWithApplication)
	}
	application.Name = strings.TrimSpace(application.Name)
	if application.Name == "" {
		application.Name = filepath.Base(application.Command)
	}
	application.Args = append([]string(nil), application.Args...)
	return application, nil
}
