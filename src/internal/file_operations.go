package internal

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yorukot/superfile/src/internal/utils"

	trash_win "github.com/hymkor/trash-go"
	"github.com/rkoesters/xdg/trash"
	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/config/icon"
)

// isSamePartition checks if two paths are on the same filesystem partition
func isSamePartition(path1, path2 string) (bool, error) {
	// Get the absolute path to handle relative paths
	absPath1, err := filepath.Abs(path1)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path of the first path: %w", err)
	}

	absPath2, err := filepath.Abs(path2)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path of the second path: %w", err)
	}

	if runtime.GOOS == utils.OsWindows {
		// On Windows, we can check if both paths are on the same drive (same letter)
		drive1 := getDriveLetter(absPath1)
		drive2 := getDriveLetter(absPath2)
		return drive1 == drive2, nil
	}

	// For Unix-like systems, we use the same path to check the root partition
	return filepath.VolumeName(absPath1) == filepath.VolumeName(absPath2), nil
}

// getDriveLetter extracts the drive letter from a Windows path
func getDriveLetter(path string) string {
	// Windows paths are usually like "C:\path\to\file"
	// So we need to extract the drive letter (e.g., "C")
	return strings.ToUpper(string(path[0]))
}

// moveElement moves a file or directory efficiently
func moveElement(src, dst string) error {
	// Check if source and destination are on the same partition
	sameDev, err := isSamePartition(src, dst)
	if err != nil {
		return fmt.Errorf("failed to check partitions: %w", err)
	}

	// If on the same partition, attempt to rename (which will use the same inode)
	if sameDev {
		if err = os.Rename(src, dst); err == nil {
			return nil
		}
		// If rename fails, fall back to copy+delete
	}

	// If on different partitions or rename failed, fall back to copy+delete
	err = copyElement(src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	err = os.RemoveAll(src)
	if err != nil {
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

// Move file to trash can and can auto switch macos trash can or linux trash can
func trashMacOrLinux(src string) error {
	var err error
	switch runtime.GOOS {
	case utils.OsDarwin:
		err = moveElement(src, filepath.Join(variable.DarwinTrashDirectory, filepath.Base(src)))
	case utils.OsWindows:
		err = trash_win.Throw(src)
	default:
		err = trash.Trash(src)
	}
	if err != nil {
		slog.Error("Error while deleting single item, in function to move file to trash can", "error", err)
	}
	return err
}

// pasteDir handles directory copying with progress tracking
// model would only have changes in m.processBarModel.process[id]
func pasteDir(src, dst string, id string, m *model) error {
	dst, err := renameIfDuplicate(dst)
	if err != nil {
		return err
	}

	// Check if we can do a fast move within the same partition
	sameDev, err := isSamePartition(src, dst)
	if err == nil && sameDev && m.copyItems.cut {
		// For cut operations on same partition, try fast rename first
		err = os.Rename(src, dst)
		if err == nil {
			return nil
		}
		// If rename fails, fall back to manual copy
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

		if info.IsDir() {
			newPath, err = renameIfDuplicate(newPath)
			if err != nil {
				return err
			}
			err = os.MkdirAll(newPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			p := m.processBarModel.process[id]
			message := channelMessage{
				messageID:       id,
				messageType:     sendProcess,
				processNewState: p,
			}

			if m.copyItems.cut {
				p.name = icon.Cut + icon.Space + filepath.Base(path)
			} else {
				p.name = icon.Copy + icon.Space + filepath.Base(path)
			}

			if len(channel) < 5 {
				message.processNewState = p
				channel <- message
			}

			var err error
			if m.copyItems.cut && sameDev {
				err = os.Rename(path, newPath)
			} else {
				err = copyFile(path, newPath, info)
			}

			if err != nil {
				p.state = failure
				message.processNewState = p
				channel <- message
				return err
			}

			p.done++
			if len(channel) < 5 {
				message.processNewState = p
				channel <- message
			}
			m.processBarModel.process[id] = p
		}
		return nil
	})

	if err != nil {
		return err
	}

	// If this was a cut operation and we had to do a manual copy, remove the source
	if m.copyItems.cut && !sameDev {
		err = os.RemoveAll(src)
		if err != nil {
			return fmt.Errorf("failed to remove source after move: %w", err)
		}
	}

	return nil
}
