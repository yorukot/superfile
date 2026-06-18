//go:build !freebsd && !linux && !windows

package metadata

// returns file attributes
// TODO: need realisation
func getFileAttributes(_ string) (string, bool) {
	return "", false
}
