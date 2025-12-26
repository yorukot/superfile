package metadata

import "strconv"

func (m *Model) SetMetadataCache(metadata Metadata, metadataFocused bool) {
	m.cache.Set(cacheKey(metadata.filepath, metadataFocused), metadata)
}

func (m *Model) SetMetadata(metadata Metadata, metadataFocused bool) {
	m.metadata = metadata
	m.SetMetadataLocationAndFocused(metadata.filepath, metadataFocused)
	// Note : Dont always reset render to 0
	// We would have update requests coming in during user scrolling through metadata
	m.ResetRenderIfInvalid()
}

func (m *Model) GetMetadataLocation() string {
	return m.expectedLocation
}

func (m *Model) GetMetadataExpectedFocused() bool {
	return m.expectedFocused
}

func (m *Model) SetMetadataLocationAndFocused(filepath string, metadataFocused bool) {
	m.expectedLocation = filepath
	m.expectedFocused = metadataFocused
}

func cacheKey(filePath string, metadataFocused bool) string {
	return filePath + ":" + strconv.FormatBool(metadataFocused)
}

func (m *Model) UpdateMetadataIfExistsInCache(filepath string, metadataFocused bool) bool {
	if meta, ok := m.cache.Get(cacheKey(filepath, metadataFocused)); ok {
		m.SetMetadata(meta, metadataFocused)
		return true
	}
	return false
}
