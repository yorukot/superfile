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

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func validateCompressOperation(sources []string) (int, error) {
	totalFiles := 0
	for _, src := range sources {
		stat, err := os.Stat(src)
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("source path does not exist: %s", src)
		}
		if err != nil {
			return 0, fmt.Errorf("cannot access source path %s: %w", src, err)
		}
		if !stat.IsDir() {
			if err = checkFileReadable(src); err != nil {
				slog.Error("the file is not readable", "error", err)
				return 0, fmt.Errorf("the file is not readable: %s", err.Error())
			}
		}
		count, err := countReadableFiles(src)
		if err != nil {
			slog.Error("Error while zip file count files ", "error", err)
			return 0, fmt.Errorf("error while counting files for %s: %s", src, err.Error())
		}
		totalFiles += count
	}
	return totalFiles, nil
}

func (m *model) getCompressSelectedFilesCmd() tea.Cmd {
	panel := m.getFocusedFilePanel()

	if panel.Empty() {
		return nil
	}
	var filesToCompress []string
	var firstFile string

	if panel.SelectedCount() == 0 {
		firstFile = panel.GetFocusedItem().Location
		filesToCompress = append(filesToCompress, firstFile)
	} else {
		firstFile = panel.GetFirstSelectedLocation()
		filesToCompress = panel.GetSelectedLocationsSortedAsVisible()
	}

	reqID := m.nextIoReqCnt()

	return func() tea.Msg {
		zipName, err := getZipArchiveName(filepath.Base(firstFile))
		if err != nil {
			slog.Error("Error in getZipArchiveName", "error", err)
			return NewNotifyModalMsg(notify.New(true, "Invalid zip target name", err.Error(), notify.NoAction),
				reqID)
		}
		zipPath := filepath.Join(panel.Location, zipName)
		totalFiles, err := validateCompressOperation(filesToCompress)
		if err != nil {
			return NewNotifyModalMsg(notify.New(true, "Invalid file/dir to compress", err.Error(), notify.NoAction),
				reqID)
		}
		if err := zipSources(filesToCompress, totalFiles, zipPath, &m.processBarModel); err != nil {
			slog.Error("Error in zipping files", "error", err)
			return NewCompressOperationMsg(processbar.Failed, reqID)
		}
		return NewCompressOperationMsg(processbar.Successful, reqID)
	}
}

func zipSources(sources []string, totalFiles int, target string, processBar *processbar.Model) error {
	var err error

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
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
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
