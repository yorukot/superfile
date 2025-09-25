//go:build windows

package metadata

import (
	"os"
)

func getOwnerAndGroup(_ os.FileInfo) (string, string) {
	return "", ""
}
