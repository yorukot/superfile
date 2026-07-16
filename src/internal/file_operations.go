package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

var localProvider = filesystem.NewLocalProvider() //nolint:gochecknoglobals // Shared stateless local provider.

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
	return withLocalSession(src, func(ctx context.Context, session filesystem.Session) error {
		return session.Move(ctx,
			filesystem.NewLocalPath(src),
			filesystem.NewLocalPath(dst),
			filesystem.MoveOptions{Overwrite: true, Recursive: true},
		)
	})
}

func deleteElement(src string) error {
	return withLocalSession(src, func(ctx context.Context, session filesystem.Session) error {
		return session.Delete(ctx,
			filesystem.NewLocalPath(src),
			filesystem.DeleteOptions{Recursive: true},
		)
	})
}

// copyLinkFile is equivalent to cp -P for a symbolic link.
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

// pasteDir handles directory copying with progress tracking
func pasteDir(src, dst string, p *processbar.Process, cut bool, processBarModel *processbar.Model) error {
	dst, err := renameIfDuplicate(dst)
	if err != nil {
		return err
	}

	// Check if we can do a fast move within the same partition
	sameDev, err := isSamePartition(src, dst)
	if err == nil && sameDev && cut {
		// For cut operations on same partition, try fast rename first
		err = os.Rename(src, dst)
		if err == nil {
			return nil
		}
		// If rename fails, fall back to manual copy
	}

	session, err := openLocalSession(src)
	if err != nil {
		return err
	}
	defer session.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(dst, relPath)
		return actualPasteOperation(session, info, path, newPath, cut, sameDev, p, processBarModel)
	})

	if err != nil {
		return err
	}

	// If this was a cut operation and we had to do a manual copy, remove the source
	if cut && !sameDev {
		err = deleteElement(src)
		if err != nil {
			return fmt.Errorf("failed to remove source after move: %w", err)
		}
	}

	return nil
}

func actualPasteOperation(
	session filesystem.Session,
	info os.FileInfo,
	path string,
	newPath string,
	cut bool,
	sameDev bool,
	p *processbar.Process,
	processBarModel *processbar.Model,
) error {
	var err error
	if info.IsDir() {
		// TODO - this is likely not needed because we did
		// dst, err := renameIfDuplicate(dst) above
		newPath, err = renameIfDuplicate(newPath)
		if err != nil {
			return err
		}
		err = session.Mkdir(context.Background(), filesystem.NewLocalPath(newPath), filesystem.MkdirOptions{
			Mode:    info.Mode(),
			Parents: true,
		})
		return err
	}
	if isSymlink(info) {
		return copyLinkFile(path, newPath)
	}

	// File
	p.CurrentFile = filepath.Base(path)
	if cut && sameDev {
		err = session.Move(context.Background(),
			filesystem.NewLocalPath(path),
			filesystem.NewLocalPath(newPath),
			filesystem.MoveOptions{Overwrite: true},
		)
	} else {
		err = session.Copy(context.Background(),
			filesystem.NewLocalPath(path),
			filesystem.NewLocalPath(newPath),
			filesystem.CopyOptions{Overwrite: true},
		)
	}

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

func openLocalSession(root string) (filesystem.Session, error) {
	return localProvider.Open(context.Background(), filesystem.Location{
		Provider: filesystem.ProviderLocal,
		Path:     filesystem.NewLocalPath(root),
		Label:    "local",
	})
}

func withLocalSession(root string, action func(context.Context, filesystem.Session) error) error {
	session, err := openLocalSession(root)
	if err != nil {
		return err
	}
	defer session.Close()

	return action(context.Background(), session)
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
