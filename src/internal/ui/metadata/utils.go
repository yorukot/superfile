package metadata

import (
	"crypto/md5" //nolint:gosec // MD5 used for file checksum display only, not for security
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/yorukot/superfile/src/internal/common"
)

func getMaxKeyLength(meta [][2]string) int {
	maxLen := 0
	for _, pair := range meta {
		if len(pair[0]) > maxLen {
			maxLen = len(pair[0])
		}
	}
	return maxLen
}

func computeMetadataWidths(viewWidth, maxKeyLen int) (int, int) {
	keyLen := maxKeyLen
	valueLen := viewWidth - keyLen
	if valueLen < viewWidth/2 {
		//nolint:mnd // standard halving for center split
		valueLen = viewWidth / 2
		keyLen = viewWidth - valueLen
	}

	return keyLen, valueLen
}

// TODO : Simplify these mystic calculations, or add explanation comments.
// TODO : unit test and fix this mess
func formatMetadataLines(meta [][2]string, startIdx, height, keyLen, valueLen int) []string {
	lines := []string{}
	endIdx := min(startIdx+height, len(meta))
	for i := startIdx; i < endIdx; i++ {
		value := common.TruncateMiddleText(meta[i][1], valueLen, "...")
		key := common.TruncateMiddleText(meta[i][0], keyLen, "...")
		line := fmt.Sprintf("%-*s%s%s", keyLen, key, keyValueSpacing, value)
		lines = append(lines, line)
	}
	return lines
}

func computeRenderDimensions(metadata [][2]string, viewWidth int) (int, int) {
	// Compute dimension based values
	maxKeyLen := getMaxKeyLength(metadata)
	return computeMetadataWidths(viewWidth, maxKeyLen)
}

// TODO : Unit test this
func calculateMD5Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New() //nolint:gosec // MD5 is sufficient for file integrity display, not used for security
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5 checksum: %w", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}
