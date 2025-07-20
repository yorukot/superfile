package metadata

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
