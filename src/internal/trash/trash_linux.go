//go:build linux

package trash

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const trashInfoDateLayout = "2006-01-02T15:04:05"

type linuxTrashDir struct {
	root     string
	files    string
	info     string
	pathBase string
	home     bool
}

func Init() error {
	return ensureTrashLayout(homeTrashDir())
}

func Available(path string) bool {
	_, err := selectTrashDir(path, true)
	return err == nil
}

func Move(path string) (Result, error) {
	srcAbs, err := filepath.Abs(path)
	if err != nil {
		return Result{OriginalPath: path, Backend: BackendFreeDesktop}, err
	}
	srcAbs = filepath.Clean(srcAbs)
	if srcAbs == string(filepath.Separator) {
		return Result{OriginalPath: srcAbs, Backend: BackendFreeDesktop},
			errors.New("refusing to trash filesystem root")
	}
	if isInsideKnownTrash(srcAbs) {
		return Result{OriginalPath: srcAbs, Backend: BackendFreeDesktop},
			fmt.Errorf("%s is already inside a trash directory", srcAbs)
	}

	td, err := selectTrashDir(srcAbs, true)
	if err != nil {
		return Result{OriginalPath: srcAbs, Backend: BackendFreeDesktop}, err
	}

	trashName, infoPath, filesPath, err := reserveTrashInfo(td, srcAbs)
	if err != nil {
		return Result{OriginalPath: srcAbs, Backend: BackendFreeDesktop}, err
	}

	if err := movePath(srcAbs, filesPath); err != nil {
		_ = os.Remove(infoPath)
		_ = os.RemoveAll(filesPath)
		return Result{OriginalPath: srcAbs, Backend: BackendFreeDesktop}, err
	}

	return Result{
		OriginalPath:     srcAbs,
		TrashedPath:      filepath.Join(td.files, trashName),
		Backend:          BackendFreeDesktop,
		StrictlyRecycled: true,
	}, nil
}

func selectTrashDir(path string, create bool) (linuxTrashDir, error) {
	srcAbs, err := filepath.Abs(path)
	if err != nil {
		return linuxTrashDir{}, err
	}
	srcAbs = filepath.Clean(srcAbs)

	home := homeTrashDir()
	if create {
		err = ensureTrashLayout(home)
		if err != nil {
			return linuxTrashDir{}, err
		}
	}

	if sameDevice(srcAbs, home.root) {
		return home, nil
	}

	topDir, err := filesystemTopDir(srcAbs)
	if err != nil {
		return linuxTrashDir{}, err
	}

	td, err := topDirectoryTrash(topDir, create)
	if err != nil {
		return linuxTrashDir{}, err
	}
	return td, nil
}

func homeTrashDir() linuxTrashDir {
	root := filepath.Join(xdgDataHome(), "Trash")
	return linuxTrashDir{
		root:     root,
		files:    filepath.Join(root, "files"),
		info:     filepath.Join(root, "info"),
		pathBase: "",
		home:     true,
	}
}

func xdgDataHome() string {
	if dataHome := os.Getenv("XDG_DATA_HOME"); filepath.IsAbs(dataHome) {
		return dataHome
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(os.TempDir(), ".local", "share")
	}
	return filepath.Join(home, ".local", "share")
}

func topDirectoryTrash(topDir string, create bool) (linuxTrashDir, error) {
	uid := strconv.Itoa(os.Getuid())
	sharedRoot := filepath.Join(topDir, ".Trash")
	if validSharedTrash(sharedRoot) {
		root := filepath.Join(sharedRoot, uid)
		td := linuxTrashDir{
			root:     root,
			files:    filepath.Join(root, "files"),
			info:     filepath.Join(root, "info"),
			pathBase: topDir,
		}
		if create {
			if err := ensureTrashLayout(td); err == nil {
				return td, nil
			}
		} else {
			return td, nil
		}
	}

	root := filepath.Join(topDir, ".Trash-"+uid)
	td := linuxTrashDir{
		root:     root,
		files:    filepath.Join(root, "files"),
		info:     filepath.Join(root, "info"),
		pathBase: topDir,
	}
	if create {
		if err := ensureTrashLayout(td); err != nil {
			return linuxTrashDir{}, err
		}
	}
	return td, nil
}

func ensureTrashLayout(td linuxTrashDir) error {
	if err := os.MkdirAll(td.root, 0o700); err != nil {
		return fmt.Errorf("failed to create trash directory %s: %w", td.root, err)
	}
	if err := ensureOwnedPrivateTrash(td.root); err != nil {
		return err
	}
	for _, dir := range []string{td.files, td.info} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("failed to create trash directory %s: %w", dir, err)
		}
		if err := ensureOwnedPrivateTrash(dir); err != nil {
			return err
		}
	}
	return nil
}

func ensureOwnedPrivateTrash(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing symlink trash directory %s", path)
	}
	if !info.IsDir() {
		return fmt.Errorf("trash path %s is not a directory", path)
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok && int(stat.Uid) != os.Getuid() {
		return fmt.Errorf("trash directory %s is not owned by current user", path)
	}
	if info.Mode().Perm()&0o077 != 0 {
		//nolint:gosec // Directory execute bit is required; trash directories are restricted to the owner.
		if err := os.Chmod(path, 0o700); err != nil {
			return fmt.Errorf("failed to restrict trash directory %s: %w", path, err)
		}
	}
	return nil
}

func validSharedTrash(path string) bool {
	info, err := os.Lstat(path)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return false
	}
	return info.Mode()&os.ModeSticky != 0
}

func reserveTrashInfo(td linuxTrashDir, srcAbs string) (string, string, string, error) {
	base := filepath.Base(srcAbs)
	if base == "." || base == string(filepath.Separator) || base == "" {
		return "", "", "", fmt.Errorf("invalid trash source %s", srcAbs)
	}
	escapedPath, err := trashInfoPath(td, srcAbs)
	if err != nil {
		return "", "", "", err
	}
	deletionDate := time.Now().Format(trashInfoDateLayout)
	content := "[Trash Info]\nPath=" + escapedPath + "\nDeletionDate=" + deletionDate + "\n"

	for i := range 10_000 {
		name := candidateTrashName(base, i)
		infoPath := filepath.Join(td.info, name+".trashinfo")
		filesPath := filepath.Join(td.files, name)
		if _, err := os.Lstat(filesPath); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return "", "", "", err
		}

		file, err := os.OpenFile(infoPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			return "", "", "", err
		}
		if _, err := file.WriteString(content); err != nil {
			_ = file.Close()
			_ = os.Remove(infoPath)
			return "", "", "", err
		}
		if err := file.Close(); err != nil {
			_ = os.Remove(infoPath)
			return "", "", "", err
		}
		return name, infoPath, filesPath, nil
	}
	return "", "", "", fmt.Errorf("could not find a free trash name for %s", srcAbs)
}

func candidateTrashName(base string, attempt int) string {
	if attempt == 0 {
		return base
	}
	return fmt.Sprintf("%s.%d.%d", base, time.Now().UnixNano(), attempt)
}

func trashInfoPath(td linuxTrashDir, srcAbs string) (string, error) {
	if td.home || td.pathBase == "" {
		return escapeTrashPath(srcAbs), nil
	}
	rel, err := filepath.Rel(td.pathBase, srcAbs)
	if err != nil {
		return "", err
	}
	if rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", fmt.Errorf("%s is not under filesystem root %s", srcAbs, td.pathBase)
	}
	return escapeTrashPath(rel), nil
}

func escapeTrashPath(path string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	for i := range len(path) {
		c := path[i]
		if c == '/' || isRFC2396Unreserved(c) {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(hex[c>>4])
		b.WriteByte(hex[c&0x0f])
	}
	return b.String()
}

func isRFC2396Unreserved(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9' ||
		strings.ContainsRune("-_.!~*'()", rune(c))
}

func sameDevice(path1, path2 string) bool {
	dev1, err := deviceID(existingDevicePath(path1))
	if err != nil {
		return false
	}
	dev2, err := deviceID(existingDevicePath(path2))
	if err != nil {
		return false
	}
	return dev1 == dev2
}

func existingDevicePath(path string) string {
	for {
		if _, err := os.Lstat(path); err == nil {
			return path
		}
		parent := filepath.Dir(path)
		if parent == path {
			return path
		}
		path = parent
	}
}

func deviceID(path string) (uint64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("cannot inspect device for %s", path)
	}
	return stat.Dev, nil
}

func filesystemTopDir(path string) (string, error) {
	start := path
	if info, err := os.Lstat(start); err == nil && !info.IsDir() {
		start = filepath.Dir(start)
	}
	start = existingDevicePath(start)
	dev, err := deviceID(start)
	if err != nil {
		return "", err
	}
	current := start
	for {
		parent := filepath.Dir(current)
		if parent == current {
			return current, nil
		}
		parentDev, err := deviceID(parent)
		if err != nil {
			return "", err
		}
		if parentDev != dev {
			return current, nil
		}
		current = parent
	}
}

func isInsideKnownTrash(path string) bool {
	path = filepath.Clean(path)
	home := homeTrashDir()
	if isPathWithin(path, home.root) {
		return true
	}
	topDir, err := filesystemTopDir(path)
	if err != nil {
		return false
	}
	uid := strconv.Itoa(os.Getuid())
	return isPathWithin(path, filepath.Join(topDir, ".Trash", uid)) ||
		isPathWithin(path, filepath.Join(topDir, ".Trash-"+uid))
}

func isPathWithin(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && (rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))))
}

func movePath(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if !errors.Is(err, syscall.EXDEV) {
		return err
	}
	if err := copyPath(src, dst); err != nil {
		return err
	}
	if err := os.RemoveAll(src); err != nil {
		_ = os.RemoveAll(dst)
		return err
	}
	return nil
}

func copyPath(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(target, dst)
	case info.IsDir():
		return copyDir(src, dst, info)
	default:
		return copyFile(src, dst, info)
	}
}

func copyDir(src, dst string, info os.FileInfo) error {
	if err := os.Mkdir(dst, 0o700); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copyPath(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
			return err
		}
	}
	if err := os.Chmod(dst, info.Mode().Perm()); err != nil {
		return err
	}
	return os.Chtimes(dst, info.ModTime(), info.ModTime())
}

func copyFile(src, dst string, info os.FileInfo) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	removeOnError := true
	defer func() {
		_ = dstFile.Close()
		if removeOnError {
			_ = os.Remove(dst)
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := dstFile.Close(); err != nil {
		return err
	}
	if err := os.Chmod(dst, info.Mode().Perm()); err != nil {
		return err
	}
	if err := os.Chtimes(dst, info.ModTime(), info.ModTime()); err != nil {
		return err
	}
	removeOnError = false
	return nil
}
