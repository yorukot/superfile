package utils

// MissingFieldIgnorer defines the interface for configuration types that can ignore missing fields
type MissingFieldIgnorer interface {
	GetIgnoreMissingFields() bool
}
