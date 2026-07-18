//go:build !darwin && !linux && !windows

package filesystem

func renameNoReplace(_, _ string) error {
	return errNoReplaceUnsupported
}
