package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ApplyConfigOverrides applies KEY=VALUE config overrides (from the
// --with-config-override flag) on top of the values loaded from the config
// file. Keys are matched case-insensitively against the config file's TOML
// keys, and values are parsed according to the field's type. Only scalar
// fields (bool, int, string) can be overridden; tables and lists are
// rejected, as is any unknown key.
func ApplyConfigOverrides(cfg *ConfigType, overrides []string) error {
	for _, override := range overrides {
		key, value, found := strings.Cut(override, "=")
		if !found {
			return errors.New(LoadConfigError(override, "expected KEY=VALUE format"))
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if err := applyConfigOverride(cfg, key, value); err != nil {
			return fmt.Errorf("%w (while applying override %q)", err, override)
		}
	}

	return nil
}

// applyConfigOverride sets the config field whose TOML key matches key.
func applyConfigOverride(cfg *ConfigType, key, value string) error {
	t := reflect.TypeOf(ConfigType{})

	for i := range t.NumField() {
		field := t.Field(i)
		tomlKey, _, _ := strings.Cut(field.Tag.Get("toml"), ",")
		if !strings.EqualFold(tomlKey, key) {
			continue
		}

		return setConfigFieldValue(cfg, field.Name, tomlKey, value)
	}

	return errors.New(LoadConfigError(key, "does not match any config key"))
}

func setConfigFieldValue(cfg *ConfigType, fieldName, key, value string) error {
	v := reflect.ValueOf(cfg).Elem().FieldByName(fieldName)

	switch v.Kind() { //nolint:exhaustive // unsupported kinds are rejected by the default case
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return errors.New(LoadConfigError(key, "must be a boolean (true or false)"))
		}
		v.SetBool(b)
	case reflect.Int:
		n, err := strconv.Atoi(value)
		if err != nil {
			return errors.New(LoadConfigError(key, "must be an integer"))
		}
		v.SetInt(int64(n))
	case reflect.String:
		v.SetString(value)
	default:
		// Unsupported kinds (maps, slices, structs) are rejected here on purpose
		return errors.New(LoadConfigError(key, "overriding this key from the CLI is not supported")) //nolint:exhaustive // default case covers all remaining kinds
	}

	return nil
}
