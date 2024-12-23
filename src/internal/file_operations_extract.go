package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
	"golift.io/xtractr"
)

func extractCompressFile(src, dest string) error {
	id := shortuuid.New()

	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle

	p := process{
		name:     icon.ExtractFile + icon.Space + "unzip file",
		progress: prog,
		state:    inOperation,
		total:    1,
		done:     0,
		doneTime: time.Time{},
	}
	message := 	channelMessage{
		messageId:       id,
		messageType: sendProcess,
		processNewState: p,
	}

	if len(channel) < 5 {
	channel <- message
	}

	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
	}

	_, _, _, err := xtractr.ExtractFile(x)

	if err != nil {
		p.state = failure
		p.doneTime = time.Now()
		message.processNewState = p
		if len(channel) < 5 {
		channel <- message
		}
		outPutLog(fmt.Sprintf("Error extracting %s: %v", src, err))
		return err
	}

	p.state = successful
	p.done = 1
	p.doneTime = time.Now()
	message.processNewState = p
	if len(channel) < 5 {
	channel <- message
	}

	return nil
}

// Extract zip file
func unzip(src, dest string) error {
	id := shortuuid.New()
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			outPutLog(fmt.Sprintf("Error closing zip reader: %v", err))
		}
	}()

	totalFiles := len(r.File)
	// progressbar
	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle
	// channel message
	p := process{
		name:     icon.ExtractFile + icon.Space + "unzip file",
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
		doneTime: time.Time{},
	}

	message := channelMessage{
		messageId: id,
		messageType: sendProcess,
		processNewState: p,
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer func() {
			if err := rc.Close(); err != nil {
				outPutLog(fmt.Sprintf("Error closing file reader: %v", err))
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)

			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
		}
		return nil
	}

	for _, f := range r.File {
		p.name = icon.ExtractFile + icon.Space + f.Name
		if len(channel) < 5 {
			message.processNewState = p
			channel <- message
		}
		err := extractAndWriteFile(f)
		if err != nil {
			p.state = failure
			message.processNewState = p
			channel <- message
			outPutLog(fmt.Sprintf("Error extracting %s: %v", f.Name, err))
			p.done++
			continue
		}
		p.done++
		if len(channel) < 5 {
			message.processNewState = p
		        channel <- message
		}
	}

	p.total = totalFiles
	p.doneTime = time.Now()
	if p.done == totalFiles {
		p.state = successful
	} else {
		p.state = failure
	}
	message.processNewState = p
	if len(channel) < 5 {
	channel <- message
	}

	return nil
}