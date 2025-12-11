package internal

import (
	"fmt"
	"log/slog"
	"time"

	"golift.io/xtractr"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func extractCompressFile(src, dest string, processBar *processbar.Model) error {
	p, err := processBar.SendAddProcessMsg(icon.ExtractFile+icon.Space+"unzip file", 1, true)
	if err != nil {
		return fmt.Errorf("cannot spawn process : %w", err)
	}

	x := &xtractr.XFile{
		FilePath:  src,
		OutputDir: dest,
		FileMode:  common.ExtractedFileMode,
		DirMode:   common.ExtractedDirMode,
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
	pSendErr := processBar.SendUpdateProcessMsg(p, true)
	if pSendErr != nil {
		slog.Error("Error sending process update", "error", pSendErr)
	}

	return err
}
