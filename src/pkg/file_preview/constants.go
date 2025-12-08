package filepreview

import "time"

// Image preview constants
const (
	// Cache configuration
	DefaultThumbnailCacheSize = 100 // Default number of thumbnails to cache
	DefaultCacheExpiration    = 5 * time.Minute
	// Image processing
	HeightScaleFactor = 2  // Factor for height scaling in terminal display
	RGBShift16        = 16 // Bit shift for red channel in RGB operations
	RGBShift8         = 8  // Bit shift for green channel in RGB operations

	// Kitty protocol
	KittyHashSeed      = 42     // Seed for kitty image ID hashing
	KittyHashPrime     = 31     // Prime multiplier for hash calculation
	KittyMaxID         = 0xFFFF // Maximum ID value for kitty images
	KittyNonZeroOffset = 1000   // Offset to ensure non-zero IDs

	// RGB color masks
	RGBMask     = 0xFF // Mask for extracting 8-bit RGB channel values
	AlphaOpaque = 255  // Fully opaque alpha channel value
)
