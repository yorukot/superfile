package metadata

import (
	"fmt"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
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

func computeMetadataWidths(metadataPanelWidth, maxKeyLen int) (int, int) {
	// Value Length = PanelLength - Key length - 2 (for border)
	valueLen := metadataPanelWidth - maxKeyLen - 2
	sprintfLen := maxKeyLen + 1
	if valueLen < metadataPanelWidth/2 {
		valueLen = metadataPanelWidth/2 - 2
		sprintfLen = valueLen
	}

	return sprintfLen, valueLen
}

// TODO : Simplify these mystic calculations, or add explanation comments.
// TODO : unit test and fix this mess
func formatMetadataLines(meta [][2]string, startIdx, height, sprintfLen, valueLen int) []string {
	lines := []string{}
	endIdx := min(startIdx+height, len(meta))
	for i := startIdx; i < endIdx; i++ {
		key := meta[i][0]
		value := common.TruncateMiddleText(meta[i][1], valueLen, "...")
		if utils.FooterWidth(0)-sprintfLen-3 < utils.FooterWidth(0)/2 {
			key = common.TruncateMiddleText(key, valueLen, "...")
		}
		line := fmt.Sprintf("%-*s %s", sprintfLen, key, value)
		lines = append(lines, line)
	}
	return lines
}

func computeRenderDimensions(metadata [][2]string, width int) (int, int) {
	// Compute dimension based values
	maxKeyLen := getMaxKeyLength(metadata)
	return computeMetadataWidths(width, maxKeyLen)
}
