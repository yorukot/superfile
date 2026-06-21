package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func delEmptyZip(target string) error {
	r, err := zip.OpenReader(target)
	if err != nil {
		return err
	}
	if len(r.File) == 0 {
		if err := r.Close(); err != nil {
			return err
		}
		return os.Remove(target)
	}
	return r.Close()
}

func zipSources(sources []string, target string, processBar *processbar.Model) error {
	var err error

	totalFiles := 0
	for _, src := range sources {
		if _, err = os.Stat(src); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("source path does not exist: %s", src)
			}
			if os.IsPermission(err) {
				return fmt.Errorf("missing permissions: %s", src)
			}
			return fmt.Errorf("cannot access source path %s: %w", src, err)
		}

		count, e := countFiles(src)
		if e != nil {
			slog.Error("Error while zip file count files ", "error", e)
		}
		totalFiles += count
	}
	p, err := processBar.SendAddProcessMsg(filepath.Base(target), processbar.OpCompress, totalFiles, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process : %w", err)
	}
	_, err = os.Stat(target)
	if err == nil {
		p.ErrorMsg = "File already exists"
		p.State = processbar.Cancelled
		p.DoneTime = time.Now()
		pSendErr := processBar.SendUpdateProcessMsg(p, true)
		if pSendErr != nil {
			slog.Error("Error sending process update", "error", pSendErr)
		}
		return errors.New("file already exists")
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := zip.NewWriter(f)
	defer writer.Close()

	zipSourcesCore(sources, processBar, &p, writer)

	if p.State != processbar.Failed {
		// TODO: User p.SetSuccessful(), p.SetFailed()
		p.State = processbar.Successful
		p.Done = totalFiles
	}
	p.DoneTime = time.Now()
	pSendErr := processBar.SendUpdateProcessMsg(p, true)
	if pSendErr != nil {
		slog.Error("Error sending process update", "error", pSendErr)
	}
	return nil
}

func zipSourcesCore(sources []string, processBar *processbar.Model,
	p *processbar.Process, writer *zip.Writer) {
	for _, src := range sources {
		srcParentDir := filepath.Dir(src)
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			p.CurrentFile = filepath.Base(path)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(srcParentDir, path)
			if err != nil {
				return err
			}

			err = writeZipFile(path, relPath, info, writer)
			if err != nil {
				return err
			}

			p.Done++
			processBar.TrySendingUpdateProcessMsg(*p)
			return nil
		})
		if err != nil {
			slog.Error("Error while zip file", "error", err)
			p.State = processbar.Failed
			break
		}
	}
}

func writeZipFile(path string, relPath string, info os.FileInfo, writer *zip.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Method = zip.Deflate
	header.Name = relPath
	if info.IsDir() {
		header.Name += "/"
	}
	headerWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	_, err = io.Copy(headerWriter, file)
	if err != nil {
		return err
	}
	return nil
}

func getZipArchiveName(base string) (string, error) {
	if len(base) == 0 {
		return "", errors.New("empty filename to compress")
	}
	zipName := strings.TrimSuffix(base, filepath.Ext(base)) + ".zip"
	runes := []rune(base)
	if runes[0] == '.' && strings.Count(base, ".") == 1 {
		zipName = base + ".zip"
	}
	zipName, err := renameIfDuplicate(zipName)
	return zipName, err
}
