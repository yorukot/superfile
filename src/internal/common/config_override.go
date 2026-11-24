package common

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ApplyConfigOverrides applies TOML overrides to the loaded config before validation.
func ApplyConfigOverrides(c *ConfigType, overrides []string) error {
	for _, override := range overrides {
		decoder := toml.NewDecoder(strings.NewReader(override))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(c); err != nil {
			return fmt.Errorf("invalid config override %q: %w", override, err)
		}
	}
	return nil
}
