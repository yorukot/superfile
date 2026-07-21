package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ApplyConfigOverrides applies a list of `key=value` overrides onto an already
// loaded ConfigType. Each key is the TOML field name (matching the `toml:"..."`
// struct tag, e.g. "debug"), so overrides line up with what users write in
// config.toml. The value is parsed into the field's Go type. An unknown key, a
// malformed `key=value` pair, or a value that cannot be parsed into the field's
// type returns a clear error naming the offending key.
//
// This is meant to be called from LoadConfigFile after the TOML file has been
// decoded into Config and before ValidateConfig runs, so overrides are
// validated together with the rest of the config.
func ApplyConfigOverrides(c *ConfigType, overrides []string) error {
	if len(overrides) == 0 {
		return nil
	}

	// Build a one-time lookup from TOML key -> struct field index.
	fieldByTomlKey := buildTomlFieldIndex(reflect.TypeOf(*c))

	val := reflect.ValueOf(c).Elem()
	for _, override := range overrides {
		key, rawValue, found := strings.Cut(override, "=")
		if !found {
			return fmt.Errorf("invalid config override %q: expected format 'key=value'", override)
		}
		// Trim only the key so overrides line up with config.toml keys. The raw
		// value is passed through verbatim: for string fields whitespace can be
		// the intended value (e.g. border_top = ' ' for a borderless layout, as
		// documented in config.toml), so it must survive exactly as written.
		// Numeric and bool fields trim internally where needed for parsing.
		key = strings.TrimSpace(key)

		fieldIdx, ok := fieldByTomlKey[key]
		if !ok {
			return fmt.Errorf("unknown config key %q in override %q", key, override)
		}

		if err := setFieldFromString(val.Field(fieldIdx), rawValue); err != nil {
			return fmt.Errorf("invalid value for config key %q: %w", key, err)
		}
	}

	return nil
}

// buildTomlFieldIndex maps each field's TOML key (the part before any comma in
// the `toml` tag) to its index in the struct. Fields without a toml tag are
// skipped, matching how the TOML decoder ignores them.
func buildTomlFieldIndex(t reflect.Type) map[string]int {
	index := make(map[string]int, t.NumField())
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("toml")
		if tag == "" {
			continue
		}
		key := strings.Split(tag, ",")[0]
		if key == "" || key == "-" {
			continue
		}
		index[key] = i
	}
	return index
}

// setFieldFromString parses raw into the type of field and assigns it. Only the
// kinds actually present in ConfigType are supported (string, bool, int, and
// []string); anything else is rejected so silent no-ops cannot happen.
func setFieldFromString(field reflect.Value, raw string) error {
	// Only the kinds present in ConfigType are handled; the default case rejects
	// everything else, so the remaining reflect.Kind values are intentionally
	// not enumerated.
	//exhaustive:ignore
	switch field.Kind() {
	case reflect.String:
		// Assign the raw value verbatim. Whitespace can be the intended value
		// for a string field (e.g. a single space for a borderless border), so
		// it must be preserved exactly as it would be in config.toml.
		field.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("expected a boolean (true/false), got %q", raw)
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return fmt.Errorf("expected an integer, got %q", raw)
		}
		field.SetInt(n)
	case reflect.Slice:
		if field.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice element type %s", field.Type().Elem().Kind())
		}
		// Comma separated list, mirroring how TOML arrays of strings are written
		// inline. An empty (or whitespace-only) string yields an empty slice
		// rather than [""], and each element is trimmed.
		var parts []string
		if trimmed := strings.TrimSpace(raw); trimmed != "" {
			parts = strings.Split(trimmed, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
		}
		field.Set(reflect.ValueOf(parts))
	default:
		return fmt.Errorf("unsupported config field type %s", field.Kind())
	}
	return nil
}
