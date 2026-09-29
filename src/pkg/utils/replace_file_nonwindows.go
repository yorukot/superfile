//go:build !windows

package utils

import "os"

func replaceFile(replacementPath, replacedPath string) error {
	return os.Rename(replacementPath, replacedPath)
}
