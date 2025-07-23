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

// Note : We dont use map[string]string, as metadata
// 1 -> We dont need to support get(key) yet. Only usage is via iterating the whole list
// 2 -> We need custom ordering

func sortMetadata(meta [][2]string) {
	priority := map[string]int{
		"Name":          0,
		"Size":          1,
		"Date Modified": 2,
		"Date Accessed": 3,
	}

	sort.SliceStable(meta, func(i, j int) bool {
		pi, iOkay := priority[meta[i][0]]
		pj, jOkay := priority[meta[j][0]]

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
		res.infoMsg = "Cannot load file stats"
		return res
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		_, symlinkErr := filepath.EvalSymlinks(filePath)
		if symlinkErr != nil {
			res.infoMsg = "Link file is broken!"
		} else {
			res.infoMsg = "This is a link file."
		}
		return res
	}
	// Add basic metadata information irrespective of what is fetched from exiftool
	// Note : we prioritize these while sorting Metadata
	name := [2]string{"Name", fileInfo.Name()}
	size := [2]string{"Size", common.FormatFileSize(fileInfo.Size())}
	modifyDate := [2]string{"Date Modified", fileInfo.ModTime().String()}
	permissions := [2]string{"Permissions", fileInfo.Mode().String()}

	if fileInfo.IsDir() && metadataFocussed {
		// TODO : Calling dirSize() could be expensive for large directories, as it recursively
		// walks the entire tree. For now we have async approach of loading metadata,
		// and its only loaded when metadata panel is focussed.
		size = [2]string{"Size", common.FormatFileSize(utils.DirSize(filePath))}
	}
	res.data = append(res.data, name, size, modifyDate, permissions)

	if common.Config.Metadata && et != nil {
		fileInfos := et.ExtractMetadata(filePath)

		for _, fileInfo := range fileInfos {
			if fileInfo.Err != nil {
				slog.Error("Error while return metadata function", "fileInfo", fileInfo, "error", fileInfo.Err)
				continue
			}
			for k, v := range fileInfo.Fields {
				res.data = append(res.data, [2]string{k, fmt.Sprintf("%v", v)})
			}
		}
	}

	if common.Config.EnableMD5Checksum {
		// Calculate MD5 checksum
		checksum, err := utils.CalculateMD5Checksum(filePath)
		if err != nil {
			slog.Error("Error calculating MD5 checksum", "error", err)
		} else {
			md5Data := [2]string{"MD5Checksum", checksum}
			res.data = append(res.data, md5Data)
		}
	}

	return res
}
