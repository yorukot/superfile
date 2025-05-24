package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func zipSources(sources []string, target string) error {
	f, err := os.Create(target)
	if err != nil {
		return err
	}

	defer f.Close()
	writer := zip.NewWriter(f)
	defer writer.Close()

	for _, src := range sources {
		srcParentDir := filepath.Dir(src)
		err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(srcParentDir, path)
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Method = zip.Deflate
			header.Name = relPath
			if info.IsDir() {
				header.Name += "/"
			}

			hw, err := writer.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(hw, file)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func getZipArchiveName(base string) (string, error) {
	zipName := strings.TrimSuffix(base, filepath.Ext(base)) + ".zip"
	zipName, err := renameIfDuplicate(zipName)
	return zipName, err

}
