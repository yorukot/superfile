package metadata

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	"github.com/barasher/go-exiftool"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

type Metadata struct {
	// Stores key value pairs
	data     [][2]string
	infoMsg  string
	filepath string
}

func NewMetadata(data [][2]string, filepath string, infoMsg string) Metadata {
	return Metadata{
		data:     data,
		filepath: filepath,
		infoMsg:  infoMsg,
	}
}

func (m Metadata) GetPath() string {
	return m.filepath
}

func (m Metadata) GetData() [][2]string {
	return m.data
}

func (m Metadata) GetValue(key string) (string, error) {
	for _, pair := range m.data {
		if pair[0] == key {
			return pair[1], nil
		}
	}

	return "", fmt.Errorf("key %s not found", key)
}

// Note : We dont use map[string]string, as metadata
// 1 -> We dont need to support get(key) yet. Only usage is via iterating the whole list
// 2 -> We need custom ordering

func sortMetadata(meta [][2]string) {
	sort.SliceStable(meta, func(i, j int) bool {
		pi, iOkay := sortPriority[meta[i][0]]
		pj, jOkay := sortPriority[meta[j][0]]

		// Both are priority fields
		if iOkay && jOkay {
			return pi < pj
		}
		// i is a priority field, and j is not
		if iOkay {
			return true
		}

		// j is a priority field, and i is not
		if jOkay {
			return false
		}

		// None of them are priority fields, sort with name
		return meta[i][0] < meta[j][0]
	})
}

func GetMetadata(filePath string, metadataFocussed bool, et *exiftool.Exiftool) Metadata {
	meta := getMetaDataUnsorted(filePath, metadataFocussed, et)
	sortMetadata(meta.data)
	return meta
}

func getMetaDataUnsorted(filePath string, metadataFocussed bool, et *exiftool.Exiftool) Metadata {
	res := Metadata{
		filepath: filePath,
	}

	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		res.infoMsg = fileStatErrorMsg
		return res
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		_, symlinkErr := filepath.EvalSymlinks(filePath)
		if symlinkErr != nil {
			res.infoMsg = linkFileBrokenMsg
		} else {
			res.infoMsg = linkFileMsg
		}
		return res
	}
	// Add basic metadata information irrespective of what is fetched from exiftool
	// Note : we prioritize these while sorting Metadata
	name := [2]string{keyName, fileInfo.Name()}
	size := [2]string{keySize, common.FormatFileSize(fileInfo.Size())}
	modifyDate := [2]string{keyDataModified, fileInfo.ModTime().String()}
	permissions := [2]string{keyPermissions, fileInfo.Mode().String()}
	ownerVal, groupVal := getOwnerAndGroup(fileInfo)
	owner := [2]string{keyOwner, ownerVal}
	group := [2]string{keyGroup, groupVal}

	if fileInfo.IsDir() && metadataFocussed {
		// TODO : Calling dirSize() could be expensive for large directories, as it recursively
		// walks the entire tree. For now we have async approach of loading metadata,
		// and its only loaded when metadata panel is focussed.
		size = [2]string{keySize, common.FormatFileSize(utils.DirSize(filePath))}
	}
	res.data = append(res.data, name, size, modifyDate, permissions, owner, group)

	updateExiftoolMetadata(filePath, et, &res)

	if common.Config.EnableMD5Checksum {
		// Calculate MD5 checksum
		checksum, err := calculateMD5Checksum(filePath)
		if err != nil {
			slog.Error("Error calculating MD5 checksum", "error", err)
		} else {
			md5Data := [2]string{keyMd5Checksum, checksum}
			res.data = append(res.data, md5Data)
		}
	}

	return res
}

func updateExiftoolMetadata(filePath string, et *exiftool.Exiftool, res *Metadata) {
	if !common.Config.Metadata || et == nil {
		return
	}
	fileInfos := et.ExtractMetadata(filePath)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			res.infoMsg = etFetchErrorMsg
			slog.Error("Error while return metadata function", "fileInfo", fileInfo, "error", fileInfo.Err)
			continue
		}
		for k, v := range fileInfo.Fields {
			res.data = append(res.data, [2]string{k, fmt.Sprintf("%v", v)})
		}
	}
}
