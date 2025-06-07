package filepreview

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// prepareImageForPreview handles the complete image preparation pipeline
func prepareImageForPreview(data []byte) (image.Image, int, int, error) {
	imgReader := bytes.NewReader(data)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, 0, 0, err
	}

	// Store original dimensions
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// Adjust orientation based on EXIF data
	exifReader := bytes.NewReader(data)
	img = adjustImageOrientation(exifReader, img)

	// Limit resolution to 1080p
	img = limitImageResolution(img, originalWidth, originalHeight)

	return img, originalWidth, originalHeight, nil
}

// limitImageResolution limits image resolution to 1080p while maintaining aspect ratio
func limitImageResolution(img image.Image, originalWidth, originalHeight int) image.Image {
	const maxImageWidth = 1920
	const maxImageHeight = 1080

	// Only resize if the image is larger than 1080p
	if originalWidth > maxImageWidth || originalHeight > maxImageHeight {
		resizedImg := imaging.Fit(img, maxImageWidth, maxImageHeight, imaging.Lanczos)
		slog.Info("Image resized for 1080p limit",
			"originalSize", fmt.Sprintf("%dx%d", originalWidth, originalHeight),
			"resizedSize", fmt.Sprintf("%dx%d", resizedImg.Bounds().Dx(), resizedImg.Bounds().Dy()),
		)
		return resizedImg
	}

	return img
}

// adjustImageOrientation adjusts image orientation based on EXIF data
func adjustImageOrientation(r *bytes.Reader, img image.Image) image.Image {
	exifData, err := exif.Decode(r)
	if err != nil {
		slog.Error("exif error", "error", err)
		return img
	}
	tag, err := exifData.Get(exif.Orientation)
	if err != nil {
		slog.Error("exif orientation error", "error", err)
		return img
	}
	orientation, err := tag.Int(0)
	if err != nil {
		slog.Error("exif orientation value error", "error", err)
		return img
	}
	return adjustOrientation(img, orientation)
}

// adjustOrientation applies the specified orientation transformation to the image
func adjustOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		return img
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		slog.Error("Invalid orientation value", "error", orientation)
		return img
	}
}

// resizeForANSI resizes image specifically for ANSI rendering
func resizeForANSI(img image.Image, maxWidth, maxHeight int) image.Image {
	// Use maxHeight*2 because each terminal row represents 2 pixel rows in ANSI rendering
	return imaging.Fit(img, maxWidth, maxHeight*2, imaging.Lanczos)
}
