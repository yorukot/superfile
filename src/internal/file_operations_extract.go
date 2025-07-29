package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"golift.io/xtractr"
)

func extractCompressFile(src, dest string, processBar *processbar.Model) error {
	p, err := processBar.SendAddProcessMsg(icon.ExtractFile+icon.Space+"unzip file", 1, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process : %w", err)
	}

	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
		FileMode:  0644,
		DirMode:   0755,
	}

	_, _, _, err = xtractr.ExtractFile(x)

	if err != nil {
		p.State = processbar.Failed
		slog.Error("Error extracting", "path", src, "error", err)
	} else {
		p.State = processbar.Successful
		p.Done = 1
	}

	p.DoneTime = time.Now()
	pSendErr := processBar.SendUpdateProcessNameMsg(p, true)
	if pSendErr != nil {
		slog.Error("Error sending process udpate", "error", pSendErr)
	}

	return err
}
