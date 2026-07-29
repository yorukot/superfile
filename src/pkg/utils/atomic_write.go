package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func writeFileAtomically(
	filePath string,
	data []byte,
	mode os.FileMode,
) (resultErr error) { //nolint:nonamedreturns // Deferred close and cleanup failures must reach the caller.
	tempFile, err := os.CreateTemp(filepath.Dir(filePath), "."+filepath.Base(filePath)+".tmp-")
	if err != nil {
		return fmt.Errorf("create replacement file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFileClosed := false
	replacementPublished := false
	defer func() {
		if !tempFileClosed {
			closeErr := tempFile.Close()
			tempFileClosed = true
			if closeErr != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("close replacement file: %w", closeErr))
			}
		}
		if !replacementPublished {
			if removeErr := os.Remove(tempPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				resultErr = errors.Join(resultErr, fmt.Errorf("remove replacement file: %w", removeErr))
			}
		}
	}()

	if err = tempFile.Chmod(mode); err != nil {
		return fmt.Errorf("set replacement file permissions: %w", err)
	}

	n, err := tempFile.Write(data)
	if err != nil {
		return fmt.Errorf("write replacement file: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("write replacement file: %w", io.ErrShortWrite)
	}

	if err = tempFile.Sync(); err != nil {
		return fmt.Errorf("sync replacement file: %w", err)
	}
	err = tempFile.Close()
	tempFileClosed = true
	if err != nil {
		return fmt.Errorf("close replacement file: %w", err)
	}

	if err = replaceFile(tempPath, filePath); err != nil {
		return fmt.Errorf("publish replacement file: %w", err)
	}
	replacementPublished = true
	return nil
}
