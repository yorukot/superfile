package utils

// MissingFieldIgnorer defines the interface for configuration types that can ignore missing fields
// during TOML file loading. Types implementing this interface can control whether missing field
// warnings are suppressed when parsing configuration files.
type MissingFieldIgnorer interface {
	GetIgnoreMissingFields() bool
}
