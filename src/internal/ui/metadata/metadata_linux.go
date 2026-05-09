//go:build linux

package metadata

import (
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

// see https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/fs.h
// Inode flags (FS_IOC_GETFLAGS / FS_IOC_SETFLAGS)
const (
	AttrsSupported = 16
	//nolint:staticcheck // reference name
	FS_SECRM_FL = 0x00000001 /* Secure deletion */
	//nolint:staticcheck // reference name
	FS_UNRM_FL = 0x00000002 /* Undelete */
	//nolint:staticcheck // reference name
	FS_COMPR_FL = 0x00000004 /* Compress file */
	//nolint:staticcheck // reference name
	FS_SYNC_FL = 0x00000008 /* Synchronous updates */
	//nolint:staticcheck // reference name
	FS_IMMUTABLE_FL = 0x00000010 /* Immutable file */
	//nolint:staticcheck // reference name
	FS_APPEND_FL = 0x00000020 /* writes to file may only append */
	//nolint:staticcheck // reference name
	FS_NODUMP_FL = 0x00000040 /* do not dump file */
	//nolint:staticcheck // reference name
	FS_NOATIME_FL = 0x00000080 /* do not update atime */
	/* End compression flags --- maybe not all used */
	//nolint:staticcheck // reference name
	FS_ENCRYPT_FL = 0x00000800 /* Encrypted file */
	//nolint:staticcheck // reference name
	FS_BTREE_FL = 0x00001000 /* btree format dir */
	//nolint:staticcheck // reference name
	FS_INDEX_FL = 0x00001000 /* hash-indexed directory */
	//nolint:staticcheck // reference name
	FS_IMAGIC_FL = 0x00002000 /* AFS directory */
	//nolint:staticcheck // reference name
	FS_JOURNAL_DATA_FL = 0x00004000 /* Reserved for ext3 */
	//nolint:staticcheck // reference name
	FS_NOTAIL_FL = 0x00008000 /* file tail should not be merged */
	//nolint:staticcheck // reference name
	FS_DIRSYNC_FL = 0x00010000 /* dirsync behaviour (directories only) */
	//nolint:staticcheck // reference name
	FS_TOPDIR_FL = 0x00020000 /* Top of directory hierarchies*/
	//nolint:staticcheck // reference name
	FS_HUGE_FILE_FL = 0x00040000 /* Reserved for ext4 */
	//nolint:staticcheck // reference name
	FS_EXTENT_FL = 0x00080000 /* Extents */
	//nolint:staticcheck // reference name
	FS_VERITY_FL = 0x00100000 /* Verity protected inode */
	//nolint:staticcheck // reference name
	FS_EA_INODE_FL = 0x00200000 /* Inode used for large EA */
	//nolint:staticcheck // reference name
	FS_EOFBLOCKS_FL = 0x00400000 /* Reserved for ext4 */
)

// returns file attributes
//
// Parameters:
// - path (string): file path
//
// Returns:
// - string: formated attribute string
// - bool: true in case attribute found
func getFileAttributes(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	flags, err := unix.IoctlGetUint32(int(f.Fd()), unix.FS_IOC_GETFLAGS)
	if err != nil {
		return "", false
	}

	return decodeChattr(flags), true
}
func decodeChattr(flags uint32) string {
	var attrs = make([]string, AttrsSupported)
	for i := range attrs {
		attrs[i] = "-"
	}
	if flags&FS_SECRM_FL != 0 {
		attrs[0] = "s"
	}
	if flags&FS_UNRM_FL != 0 {
		attrs[1] = "u"
	}
	if flags&FS_COMPR_FL != 0 {
		attrs[2] = "c"
	}
	if flags&FS_SYNC_FL != 0 {
		attrs[3] = "S"
	}
	if flags&FS_IMMUTABLE_FL != 0 {
		attrs[4] = "i"
	}
	if flags&FS_APPEND_FL != 0 {
		attrs[5] = "a"
	}
	if flags&FS_NODUMP_FL != 0 {
		attrs[6] = "d"
	}
	if flags&FS_NOATIME_FL != 0 {
		attrs[7] = "A"
	}
	if flags&FS_INDEX_FL != 0 {
		attrs[9] = "I"
	}
	if flags&FS_JOURNAL_DATA_FL != 0 {
		attrs[10] = "j"
	}
	if flags&FS_NOTAIL_FL != 0 {
		attrs[11] = "t"
	}
	if flags&FS_TOPDIR_FL != 0 {
		attrs[12] = "T"
	}
	if flags&FS_HUGE_FILE_FL != 0 {
		attrs[13] = "h"
	}
	if flags&FS_EXTENT_FL != 0 {
		attrs[14] = "e"
	}
	if flags&FS_DIRSYNC_FL != 0 {
		attrs[15] = "D"
	}

	return strings.Join(attrs, "")
}
