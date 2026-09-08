//go:build !darwin && !linux && !windows

package internal

import "errors"

func renameNoReplace(_, _ string) error {
	return errors.ErrUnsupported
}
