package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"golift.io/xtractr"
)

func extractCompressFile(src, dest string) error {
	id := shortuuid.New()

	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle

	p := process{
		name:     "󰛫 unzip file",
		progress: prog,
		state:    inOperation,
		total:    1,
		done:     0,
	}
	
	channel <- channelMessage{
		messageId:       id,
		processNewState: p,
	}

	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
	}

	_, _, _, err := xtractr.ExtractFile(x)

	if err != nil {
		p.state = successful
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
		return err
	}

	p.state = successful
	p.done = 1

	channel <- channelMessage{
		messageId:       id,
		processNewState: p,
	}
	
	return nil
}

// Extract zip file
func unzip(src, dest string) error {
	id := shortuuid.New()
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	totalFiles := len(r.File)
	// progessbar
	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle
	// channel message
	p := process{
		name:     "󰛫 unzip file",
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
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
				return fmt.Errorf("error open file: %s", err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)

			if err != nil {
				return fmt.Errorf("error copy file: %s", err)
			}
		}
		return nil
	}

	for _, f := range r.File {
		p.name = "󰛫 " + f.Name
		if len(channel) < 3 {
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
		}
		err := extractAndWriteFile(f)
		if err != nil {
			p.state = failure
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
			return err
		}
		p.done++
		if len(channel) < 3 {
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
		}
	}

	p.total = totalFiles
	p.state = successful
	channel <- channelMessage{
		messageId:       id,
		processNewState: p,
	}

	return nil
}