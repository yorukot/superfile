package metadata

import (
	"fmt"

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
		key := meta[i][0]
		value := common.TruncateMiddleText(meta[i][1], valueLen, "...")
		key = common.TruncateMiddleText(key, keyLen-1, "...")
		line := fmt.Sprintf("%-*s %s", keyLen, key, value)
		lines = append(lines, line)
	}
	return lines
}

func computeRenderDimensions(metadata [][2]string, viewWidth int) (int, int) {
	// Compute dimension based values
	maxKeyLen := getMaxKeyLength(metadata)
	return computeMetadataWidths(viewWidth, maxKeyLen)
}
