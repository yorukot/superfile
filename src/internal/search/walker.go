package search

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// walk traverses root and visits every visible entry with its path relative
// to root. It descends into symlinked directories but never follows a link
// whose target is already on the current ancestor chain, so symlink cycles
// cannot cause infinite recursion. unreadable receives the count of
// directories that could not be read (permission denied and the like);
// unreadable entries are skipped, not treated as errors.
func walk(ctx context.Context, root string, includeHidden bool, unreadable *int,
	visit func(relPath string, isDir bool)) {
	rootInfo, err := os.Stat(root)
	if err != nil {
		*unreadable++
		return
	}
	chain := []os.FileInfo{rootInfo}

	var walkDir func(abs, rel string, ancestors []os.FileInfo)
	walkDir = func(abs, rel string, ancestors []os.FileInfo) {
		if ctx.Err() != nil {
			return
		}
		entries, err := os.ReadDir(abs)
		if err != nil {
			*unreadable++
			return
		}
		for _, entry := range entries {
			if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			visitEntry(abs, rel, entry.Name(), info, ancestors, visit, walkDir)
		}
	}

	walkDir(root, "", chain)
}

// visitEntry reports one directory entry to visit and recurses into
// directories (following symlinks to directories, provided the target is not
// an ancestor, which would loop forever).
func visitEntry(abs, rel, name string, info os.FileInfo, ancestors []os.FileInfo,
	visit func(relPath string, isDir bool), walkDir func(abs, rel string, ancestors []os.FileInfo)) {
	entryAbs := filepath.Join(abs, name)
	entryRel := filepath.Join(rel, name)
	switch {
	case info.IsDir():
		visit(entryRel, true)
		walkDir(entryAbs, entryRel, ancestors)
	case info.Mode()&fs.ModeSymlink != 0:
		// Resolve the link; broken links are reported as files.
		target, err := os.Stat(entryAbs)
		if err != nil {
			visit(entryRel, false)
			return
		}
		if !target.IsDir() {
			visit(entryRel, false)
			return
		}
		visit(entryRel, true)
		if containsSameFile(target, ancestors) {
			return
		}
		walkDir(entryAbs, entryRel, append(ancestors, target))
	default:
		visit(entryRel, false)
	}
}

func containsSameFile(target os.FileInfo, ancestors []os.FileInfo) bool {
	for _, ancestor := range ancestors {
		if os.SameFile(target, ancestor) {
			return true
		}
	}
	return false
}
