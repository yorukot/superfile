package filepreview

import "time"

// Image preview constants
const (
	// Cache configuration
	defaultImagePreviewCacheSize = 100
	defaultCacheExpiration       = 5 * time.Minute

	// Image processing
	heightScaleFactor = 2  // Factor for height scaling in terminal display
	rgbShift16        = 16 // Bit shift for red channel in RGB operations
	rgbShift8         = 8  // Bit shift for green channel in RGB operations

	// RGB color masks
	rgbMask     = 0xFF // Mask for extracting 8-bit RGB channel values
	alphaOpaque = 255  // Fully opaque alpha channel value

	maxVideoFileSizeForThumb = "104857600" // 100MB limit
	thumbOutputExt           = ".jpg"
	thumbGenerationTimeout   = 30 * time.Second
)
