package contentpanel

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
)

// checkEditable determines whether a file can be safely edited in the content
// panel. Returns true and empty reason if editable, or false and a short
// human-readable reason why not.
func checkEditable(path string, info os.FileInfo, maxSize int64) (bool, string) {
	// Must be a regular file (not directory, symlink-to-dir, device, etc.)
	if !info.Mode().IsRegular() {
		return false, "Not a regular file"
	}

	// Must be writable by the current user
	if info.Mode().Perm()&0o200 == 0 {
		// Check effective access too — the file mode bits may not tell the
		// whole story on some systems
		f, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return false, "No write permission"
		}
		f.Close()
	}

	// Size limit
	if maxSize > 0 && info.Size() > maxSize {
		return false, "File too large"
	}

	// Must be a text file, not binary
	file, err := os.Open(path)
	if err != nil {
		return false, "Cannot open file"
	}
	defer file.Close()

	// Read up to 8KB to detect binary content
	buf := make([]byte, min(info.Size(), 8192))
	n, err := file.Read(buf)
	if err != nil && n == 0 {
		return false, "Cannot read file"
	}
	buf = buf[:n]

	// Use the common text-detection logic
	isText, err := common.IsTextFile(path)
	if err != nil {
		return false, "Cannot determine file type"
	}
	if !isText {
		return false, "Binary file"
	}

	return true, ""
}

// Ensure we can re-read the full file for editing.
func readFileContent(path string, maxSize int64) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if maxSize > 0 && info.Size() > maxSize {
		return "", fs.ErrPermission // reuse as "too large"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsCodeFile checks whether the file's name or extension matches the
// configurable code-extensions list (e.g. .go, .py, *_test.go).
func IsCodeFile(path string, codeExts []string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	for _, pattern := range codeExts {
		// Support "*_test.go" style patterns
		if strings.Contains(pattern, "*") {
			matched, err := filepath.Match(pattern, base)
			if err == nil && matched {
				return true
			}
		} else if ext == pattern || strings.HasSuffix(base, pattern) {
			return true
		}
	}
	return false
}
