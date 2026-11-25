package filepreview

import (
	"bytes"
	"image"
	"log/slog"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// prepareImageForPreview decodes data, shrinks it to fit a maxWidth x maxHeight
// pixel box, and applies the EXIF orientation.
func prepareImageForPreview(data []byte, maxWidth, maxHeight int) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	orientation := exifOrientation(data)
	if orientation >= exifOrientationTranspose && orientation <= exifOrientationRotate90 {
		// These orientations swap width and height, so fit the stored image
		// into the swapped box and rotate the small result.
		maxWidth, maxHeight = maxHeight, maxWidth
	}

	// Box averages every source pixel, which is the right filter for a large
	// shrink and costs much less than Lanczos on a full-size photo.
	if img.Bounds().Dx() > maxWidth || img.Bounds().Dy() > maxHeight {
		img = imaging.Fit(img, maxWidth, maxHeight, imaging.Box)
	}

	return adjustOrientation(img, orientation), nil
}

// exifOrientation returns the EXIF orientation tag of data, or 1 when it has none
func exifOrientation(data []byte) int {
	exifData, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		slog.Error("exif error", "error", err)
		return 1
	}
	tag, err := exifData.Get(exif.Orientation)
	if err != nil {
		slog.Error("exif orientation error", "error", err)
		return 1
	}
	orientation, err := tag.Int(0)
	if err != nil {
		slog.Error("exif orientation value error", "error", err)
		return 1
	}
	return orientation
}

// adjustOrientation applies the specified orientation transformation to the image
func adjustOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		return img
	case 2: //nolint:mnd // EXIF orientation: horizontal flip
		return imaging.FlipH(img)
	case 3: //nolint:mnd // EXIF orientation: 180 rotation
		return imaging.Rotate180(img)
	case 4: //nolint:mnd // EXIF orientation: vertical flip
		return imaging.FlipV(img)
	case 5: //nolint:mnd // EXIF orientation: transpose
		return imaging.Transpose(img)
	case 6: //nolint:mnd // EXIF orientation: 270 rotation
		return imaging.Rotate270(img)
	case 7: //nolint:mnd // EXIF orientation: transverse
		return imaging.Transverse(img)
	case 8: //nolint:mnd // EXIF orientation: 90 rotation
		return imaging.Rotate90(img)
	default:
		slog.Error("Invalid orientation value", "error", orientation)
		return img
	}
}

// resizeForANSI resizes image specifically for ANSI rendering
func resizeForANSI(img image.Image, maxWidth, maxHeight int) image.Image {
	// Use maxHeight*2 because each terminal row represents 2 pixel rows in ANSI rendering
	return imaging.Fit(img, maxWidth, maxHeight*heightScaleFactor, imaging.Lanczos)
}
