package utils

// ConfigInterface defines the interface for configuration types that can ignore missing fields
type ConfigInterface interface {
	GetIgnoreMissingFields() bool
}
