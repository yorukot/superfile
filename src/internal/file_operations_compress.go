package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
)

func zipSource(source, target string) error {
	id := shortuuid.New()
	prog := progress.New()
	prog.PercentageStyle = footerStyle

	totalFiles, err := countFiles(source)

	if err != nil {
		outPutLog("Zip file count files error: ", err)
	}

	p := process{
		name:     "zip files",
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}

	message := channelMessage{
		messageId:       id,
		messageType:     sendProcess,
		processNewState: p,
	}

	_, err = os.Stat(target)
	if os.IsExist(err) {
		p.name = icon.CompressFile + icon.Space + "File already exist"
		message.processNewState = p
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

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		p.name = icon.CompressFile+ icon.Space + filepath.Base(path)
		if len(channel) < 5 {
			message.processNewState = p
			channel <- message
		}

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
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

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		if err != nil {
			return err
		}
		p.done++
		if len(channel) < 5 {

			message.processNewState = p
			channel <- message
		}
		return nil
	})

	if err != nil {
		outPutLog("Error while zip file:", err)
		p.state = failure
		message.processNewState = p
		channel <- message
	}
	p.state = successful
	p.done = totalFiles

	message.processNewState = p
	channel <- message

	return nil
}
