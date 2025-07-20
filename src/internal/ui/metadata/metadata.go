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

func GetMetadata(filePath string, metadataFocussed bool, et *exiftool.Exiftool) [][2]string {
	meta := getMetaDataUnsorted(filePath, metadataFocussed, et)
	sortMetadata(meta)
	return meta
}

func getMetaDataUnsorted(filePath string, metadataFocussed bool, et *exiftool.Exiftool) [][2]string {
	var res [][2]string
	fileInfo, err := os.Stat(filePath)

	if utils.IsSymlink(filePath) {
		_, symlinkErr := filepath.EvalSymlinks(filePath)
		if symlinkErr != nil {
			res = append(res, [2]string{"Link file is broken!", ""})
		} else {
			res = append(res, [2]string{"This is a link file.", ""})
		}
		return res
	}

	if err != nil {
		slog.Error("Error while getting file state in getMetadata", "error", err)
		return res
	}

	if fileInfo.IsDir() {
		res = append(res, [2]string{"Name", fileInfo.Name()})
		if metadataFocussed {
			// TODO : Calling dirSize() could be expensive for large directories, as it recursively
			// walks the entire tree. For now we have async approach of loading metadata,
			// and its only loaded when metadata panel is focussed.
			res = append(res, [2]string{"Size", common.FormatFileSize(utils.DirSize(filePath))})
		}
		res = append(res,
			[2]string{"Date Modified", fileInfo.ModTime().String()},
			[2]string{"Permissions", fileInfo.Mode().String()})
		return res
	}

	checkIsSymlinked, err := os.Lstat(filePath)
	if err != nil {
		slog.Error("Error when getting file info", "error", err)
		return res
	}

	if common.Config.Metadata && checkIsSymlinked.Mode()&os.ModeSymlink == 0 && et != nil {
		fileInfos := et.ExtractMetadata(filePath)

		for _, fileInfo := range fileInfos {
			if fileInfo.Err != nil {
				slog.Error("Error while return metadata function", "fileInfo", fileInfo, "error", fileInfo.Err)
				continue
			}
			for k, v := range fileInfo.Fields {
				res = append(res, [2]string{k, fmt.Sprintf("%v", v)})
			}
		}
	} else {
		fileName := [2]string{"Name", fileInfo.Name()}
		fileSize := [2]string{"Size", common.FormatFileSize(fileInfo.Size())}
		fileModifyData := [2]string{"Date Modified", fileInfo.ModTime().String()}
		filePermissions := [2]string{"Permissions", fileInfo.Mode().String()}

		if common.Config.EnableMD5Checksum {
			// Calculate MD5 checksum
			checksum, err := utils.CalculateMD5Checksum(filePath)
			if err != nil {
				slog.Error("Error calculating MD5 checksum", "error", err)
			} else {
				md5Data := [2]string{"MD5Checksum", checksum}
				res = append(res, md5Data)
			}
		}

		res = append(res, fileName, fileSize, fileModifyData, filePermissions)
	}
	return res
}
