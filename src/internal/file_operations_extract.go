package internal

import (
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"golift.io/xtractr"
)

func extractCompressFile(src, dest string) error {
	id := shortuuid.New()

	prog := progress.New(common.GenerateGradientColor())
	prog.PercentageStyle = common.FooterStyle

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
		FileMode:  0644,
		DirMode:   0755,
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
