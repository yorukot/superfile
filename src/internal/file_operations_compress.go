package internal

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func zipSource(sources []string, target string) error {
	id := shortuuid.New()
	prog := progress.New()
	prog.PercentageStyle = common.FooterStyle

	totalFiles := len(sources) // Count of selected files
	p := process{
		name:     "zip files",
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}

	message := channelMessage{
		messageID:       id,
		messageType:     sendProcess,
		processNewState: p,
	}

	_, err := os.Stat(target)
	if !os.IsNotExist(err) {
		p.name = icon.CompressFile + icon.Space + "File already exists"
		message.processNewState = p
		p.done = 100
		channel <- message

		return nil
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	for _, path := range sources {
		info, err := os.Stat(path)
		if err != nil {
			slog.Error("Error getting file info:", "Error", err)
			continue // Skip to the next file
		}

		p.name = icon.CompressFile + icon.Space + filepath.Base(path)
		if len(channel) < 5 {
			message.processNewState = p
			channel <- message
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		// Set the header name to the file's name
		header.Name = filepath.Base(path)

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(headerWriter, f)
			if err != nil {
				return err
			}
		}

		p.done++
		if len(channel) < 5 {
			message.processNewState = p
			channel <- message
		}
	}

	p.state = successful
	message.processNewState = p
	channel <- message

	return nil
}
