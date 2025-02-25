package internal

import (
	"log/slog"
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
	message := channelMessage{
		messageId:       id,
		messageType:     sendProcess,
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
		slog.Error("Error extracting", "path", src, "error", err)
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
