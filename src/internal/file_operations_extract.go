package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
	"golift.io/xtractr"
)

func getDefaultFileMode() os.FileMode {
	if runtime.GOOS == "windows" {
		return 0666
	}
	return 0644
}

func shouldSkipFile(name string) bool {
	// Skip system files across platforms
	return strings.HasPrefix(name, "__MACOSX/") ||
		strings.EqualFold(name, "Thumbs.db") ||
		strings.EqualFold(name, "desktop.ini")
}

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

		// Cross-platform path security check
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		fileMode := f.Mode()
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, fileMode)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			return nil
		}

		// Create directory structure
		if err := os.MkdirAll(filepath.Dir(path), fileMode); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// Try default permissions first
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, getDefaultFileMode())
		if err != nil {
			// Fall back to original file permissions
			outFile, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
		}
		defer func() {
			if err := outFile.Close(); err != nil {
				outPutLog(fmt.Sprintf("Error closing output file %s: %v", path, err))
			}
		}()

		if _, err := io.Copy(outFile, rc); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		return nil
	}

	for _, f := range r.File {
		p.name = icon.ExtractFile + icon.Space + f.Name
		if len(channel) < 5 {
			message.processNewState = p
			channel <- message
		}

		if shouldSkipFile(f.Name) {
			p.done++
			continue
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