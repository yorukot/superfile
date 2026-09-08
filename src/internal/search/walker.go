package search

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// walk visits entries under root relative to root. It follows symlinked dirs
// unless the target is an ancestor. It counts and skips unreadable dirs.
// It stops when ctx cancels.
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
			if ctx.Err() != nil {
				return
			}
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

// visitEntry reports one entry and recurses into dirs, following non-ancestor symlink targets.
func visitEntry(abs, rel, name string, info os.FileInfo, ancestors []os.FileInfo,
	visit func(relPath string, isDir bool), walkDir func(abs, rel string, ancestors []os.FileInfo)) {
	entryAbs := filepath.Join(abs, name)
	entryRel := filepath.Join(rel, name)
	switch {
	case info.IsDir():
		visit(entryRel, true)
		walkDir(entryAbs, entryRel, append(ancestors, info))
	case info.Mode()&fs.ModeSymlink != 0:
		// Broken links report as files.
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
