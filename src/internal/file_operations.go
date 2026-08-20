package internal

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

// moveElement moves a file or directory efficiently
func moveElement(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	if err := copyElement(src, dst); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	if err := os.RemoveAll(src); err != nil {
		return fmt.Errorf("failed to remove source after copy: %w", err)
	}

	return nil
}

// copyElement handles copying of both files and directories
func copyElement(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst, srcInfo)
	}
	return copyFile(src, dst, srcInfo)
}

// copyDir recursively copies a directory
func copyDir(src, dst string, srcInfo os.FileInfo) error {
	err := os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to get entry info: %w", err)
		}

		if entryInfo.IsDir() {
			err = copyDir(srcPath, dstPath, entryInfo)
		} else {
			err = copyFile(srcPath, dstPath, entryInfo)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// is an equivalent of "cp -P" command
func copyLinkFile(src, dst string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(target, dst)
}

func isSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

// copyFile copies a single file
func copyFile(src, dst string, srcInfo os.FileInfo) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	return nil
}

// pasteDir handles directory copying with progress tracking
func pasteDir(src, dst string, p *processbar.Process, cut bool, processBarModel *processbar.Model) error {
	dst, err := renameIfDuplicate(dst)
	if err != nil {
		return err
	}

	if cut {
		if renameErr := os.Rename(src, dst); renameErr == nil {
			return nil
		}
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(dst, relPath)
		return actualPasteOperation(info, path, newPath, p, processBarModel)
	})

	if err != nil {
		return err
	}

	if cut {
		if err := os.RemoveAll(src); err != nil {
			return fmt.Errorf("failed to remove source after move: %w", err)
		}
	}

	return nil
}

func actualPasteOperation(info os.FileInfo, path string, newPath string,
	p *processbar.Process, processBarModel *processbar.Model) error {
	var err error
	if info.IsDir() {
		// TODO - this is likely not needed because we did
		// dst, err := renameIfDuplicate(dst) above
		newPath, err = renameIfDuplicate(newPath)
		if err != nil {
			return err
		}
		err = os.MkdirAll(newPath, info.Mode())
		return err
	}
	if isSymlink(info) {
		return copyLinkFile(path, newPath)
	}

	// File
	p.CurrentFile = filepath.Base(path)
	err = copyFile(path, newPath, info)

	if err != nil {
		p.State = processbar.Failed
		pSendErr := processBarModel.SendUpdateProcessMsg(*p, true)
		if pSendErr != nil {
			slog.Error("Error sending process update", "error", pSendErr)
		}
		return err
	}

	p.Done++
	processBarModel.TrySendingUpdateProcessMsg(*p)
	return nil
}

// isAncestor checks if dst is the same as src or a subdirectory of src.
// It handles symlinks by resolving them and applies case-insensitive comparison on Windows.
func isAncestor(src, dst string) bool {
	// Resolve symlinks for both paths
	srcResolved, err := filepath.EvalSymlinks(src)
	if err != nil {
		// If we can't resolve symlinks, fall back to original path
		srcResolved = src
	}

	dstResolved, err := filepath.EvalSymlinks(dst)
	if err != nil {
		// If we can't resolve symlinks, fall back to original path
		dstResolved = dst
	}

	// Get absolute paths. Abs() also Cleans paths to normalize separators and resolve . and ..
	srcAbs, err := filepath.Abs(srcResolved)
	if err != nil {
		return false
	}

	dstAbs, err := filepath.Abs(dstResolved)
	if err != nil {
		return false
	}

	// On Windows, perform case-insensitive comparison
	if runtime.GOOS == "windows" {
		srcAbs = strings.ToLower(srcAbs)
		dstAbs = strings.ToLower(dstAbs)
	}

	// Check if dst is the same as src
	if srcAbs == dstAbs {
		return true
	}

	// Check if dst is a subdirectory of src
	// Use filepath.Rel to check the relationship
	rel, err := filepath.Rel(srcAbs, dstAbs)
	if err != nil {
		return false
	}

	// If rel is "." or doesn't start with "..", then dst is inside src
	return rel == "." || !strings.HasPrefix(rel, "..")
}
