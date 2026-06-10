//go:build freebsd

package metadata

import (
	"strings"
	"syscall"
)

// flags as defined in <sys/stat.h>
const (
	UF_NODUMP    = 0x00000001
	UF_IMMUTABLE = 0x00000002
	UF_APPEND    = 0x00000004
	UF_OPAQUE    = 0x00000008
	UF_NOUNLINK  = 0x00000010
	UF_SYSTEM    = 0x00000080
	UF_SPARSE    = 0x00000100
	UF_OFFLINE   = 0x00000200
	UF_REPARSE   = 0x00000400
	UF_ARCHIVE   = 0x00000800
	UF_READONLY  = 0x00001000
	UF_HIDDEN    = 0x00008000

	SF_ARCHIVED  = 0x00010000
	SF_IMMUTABLE = 0x00020000
	SF_APPEND    = 0x00040000
	SF_NOUNLINK  = 0x00100000
	SF_SNAPSHOT  = 0x00200000
)

// Mapping from flags to string.
// The order matches that in lib/lib/gen/strtofflags.c on FreeBSD.
var mapping = []struct {
	flag uint32
	name string
}{
	{SF_APPEND, "sappnd"},
	{SF_ARCHIVED, "arch"},
	{SF_IMMUTABLE, "schg"},
	{SF_NOUNLINK, "sunlnk"},
	{SF_SNAPSHOT, "snapshot"},

	{UF_APPEND, "uappnd"},
	{UF_ARCHIVE, "uarch"},
	{UF_HIDDEN, "hidden"},
	{UF_IMMUTABLE, "uchg"},
	{UF_NODUMP, "nodump"},
	{UF_NOUNLINK, "uunlnk"},
	{UF_OFFLINE, "offline"},
	{UF_OPAQUE, "opaque"},
	{UF_READONLY, "rdonly"},
	{UF_REPARSE, "reparse"},
	{UF_SPARSE, "sparse"},
	{UF_SYSTEM, "system"},
}

func getFileAttributes(path string) (string, bool) {
	var st syscall.Stat_t

	err := syscall.Stat(path, &st)
	if err != nil {
		return "", false
	}

	return fflagsToStr(st.Flags), true
}

// Format file flags as string.
// Patterned after the FreeBSD libc function fflagstostr(3).
func fflagsToStr(flags uint32) string {
	var names []string

	for _, m := range mapping {
		if flags&m.flag != 0 {
			names = append(names, m.name)
		}
	}

	return strings.Join(names, ",")
}
