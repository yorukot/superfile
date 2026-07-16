package filesystem

import (
	"context"
	"path/filepath"
	"time"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func TrackTransferProcess(ctx context.Context, model *processbar.Model, transfer Transfer) (processbar.Process, error) {
	first, ok := <-transfer.Progress()
	if !ok {
		waitErr := transfer.Wait(ctx)
		if waitErr != nil {
			return processbar.Process{}, waitErr
		}
		return processbar.Process{}, nil
	}

	process, err := model.SendAddProcessMsg(
		transferDisplayName(first.Current),
		processbarOperationForTransfer(transfer.Operation(), transfer.Direction()),
		int(first.Total),
		true,
	)
	if err != nil {
		return processbar.Process{}, err
	}

	process.Total = int(first.Total)
	process.Done = int(first.Done)
	process.CurrentFile = transferDisplayName(first.Current)
	model.TrySendingUpdateProcessMsg(process)

	go func(current processbar.Process, initial Progress) {
		last := initial
		apply := func(progress Progress) {
			last = progress
			current.Total = int(progress.Total)
			current.Done = int(progress.Done)
			current.CurrentFile = transferDisplayName(progress.Current)
			if progress.Err != nil {
				current.State = processbar.Failed
				current.ErrorMsg = progress.Err.Error()
			}
			model.TrySendingUpdateProcessMsg(current)
		}

		apply(initial)
		for progress := range transfer.Progress() {
			apply(progress)
		}

		if waitErr := transfer.Wait(ctx); waitErr != nil {
			current.State = processbar.Failed
			current.ErrorMsg = waitErr.Error()
		} else {
			current.State = processbar.Successful
			current.Done = current.Total
		}
		current.DoneTime = time.Now()
		if last.Current.String() != "" {
			current.CurrentFile = transferDisplayName(last.Current)
		}
		_ = model.SendUpdateProcessMsg(current, true)
	}(process, first)

	return process, nil
}

func processbarOperationForTransfer(operation Operation, direction TransferDirection) processbar.OperationType {
	if operation == OperationCutMove || operation == OperationRemoteSameSessionMove {
		return processbar.OpCut
	}
	if direction == TransferUpload {
		return processbar.OpUpload
	}
	if direction == TransferDownload {
		return processbar.OpDownload
	}
	return processbar.OpCopy
}

func transferDisplayName(path Path) string {
	if path.IsRemote() {
		return path.Base()
	}
	base := filepath.Base(path.String())
	if base == "." || base == string(filepath.Separator) || base == "" {
		return path.String()
	}
	return base
}
