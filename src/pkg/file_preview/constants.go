package filepreview

import "time"

// Image preview constants
const (
	// Cache configuration
	defaultThumbnailCacheSize = 100 // Default number of thumbnails to cache
	defaultCacheExpiration    = 5 * time.Minute

	// Image processing
	heightScaleFactor = 2  // Factor for height scaling in terminal display
	rgbShift16        = 16 // Bit shift for red channel in RGB operations
	rgbShift8         = 8  // Bit shift for green channel in RGB operations

	// Kitty protocol
	kittyHashSeed      = 42     // Seed for kitty image ID hashing
	kittyHashPrime     = 31     // Prime multiplier for hash calculation
	kittyMaxID         = 0xFFFF // Maximum ID value for kitty images
	kittyNonZeroOffset = 1000   // Offset to ensure non-zero IDs

	// RGB color masks
	rgbMask     = 0xFF // Mask for extracting 8-bit RGB channel values
	alphaOpaque = 255  // Fully opaque alpha channel value
)
