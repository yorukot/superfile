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
	AttrsSupported = 22
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
	//nolint:staticcheck // reference name
	FS_NOCOMP_FL = 0x00000400 /* Don't compress */
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
	//nolint:staticcheck // reference name
	FS_NOCOW_FL = 0x00800000 /* Do not cow file */
	//nolint:staticcheck // reference name
	FS_DAX_FL = 0x02000000 /* Inode is DAX */
	//nolint:staticcheck // reference name
	FS_INLINE_DATA_FL = 0x10000000 /* Reserved for ext4 */
	//nolint:staticcheck // reference name
	FS_PROJINHERIT_FL = 0x20000000 /* Create with parents projid */
	//nolint:staticcheck // reference name
	FS_CASEFOLD_FL = 0x40000000 /* Folder is case insensitive */
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

// list of letter and output order see in
// lsattr sources https://github.com/tytso/e2fsprogs/blob/master/lib/e2p/pf.c
//
//nolint:gocognit // it's just mapping list of linux flags to lsattr letters
func decodeChattr(flags uint32) string {
	var attrs = make([]string, AttrsSupported)
	for i := range attrs {
		attrs[i] = "-"
	}

	if flags&FS_SECRM_FL != 0 {
		attrs[0] = "s" // Secure_Deletion
	}
	if flags&FS_UNRM_FL != 0 {
		attrs[1] = "u" // Undelete
	}
	if flags&FS_SYNC_FL != 0 {
		attrs[2] = "S" // Synchronous_Updates
	}
	if flags&FS_DIRSYNC_FL != 0 {
		attrs[3] = "D" // Synchronous_Directory_Updates
	}
	if flags&FS_IMMUTABLE_FL != 0 {
		attrs[4] = "i" // Immutable
	}
	if flags&FS_APPEND_FL != 0 {
		attrs[5] = "a" // Append_Only
	}
	if flags&FS_NODUMP_FL != 0 {
		attrs[6] = "d" // No_Dump
	}
	if flags&FS_NOATIME_FL != 0 {
		attrs[7] = "A" // No_Atime
	}
	if flags&FS_COMPR_FL != 0 {
		attrs[8] = "c" // Compression_Requested
	}
	if flags&FS_ENCRYPT_FL != 0 {
		attrs[9] = "E" // Encrypted
	}
	if flags&FS_JOURNAL_DATA_FL != 0 {
		attrs[10] = "j" // Journaled_Data
	}
	if flags&FS_INDEX_FL != 0 {
		attrs[11] = "I" // Indexed_directory
	}
	if flags&FS_NOTAIL_FL != 0 {
		attrs[12] = "t" // No_Tailmerging
	}
	if flags&FS_TOPDIR_FL != 0 {
		attrs[13] = "T" // Top_of_Directory_Hierarchies
	}
	if flags&FS_EXTENT_FL != 0 {
		attrs[14] = "e" // Extents
	}
	if flags&FS_NOCOW_FL != 0 {
		attrs[15] = "C" // No_COW
	}
	if flags&FS_DAX_FL != 0 {
		attrs[16] = "x" // DAX
	}
	if flags&FS_CASEFOLD_FL != 0 {
		attrs[17] = "F" // Casefold
	}
	if flags&FS_INLINE_DATA_FL != 0 {
		attrs[18] = "N" // Inline_Data
	}
	if flags&FS_PROJINHERIT_FL != 0 {
		attrs[19] = "P" // Project_Hierarchy
	}
	if flags&FS_VERITY_FL != 0 {
		attrs[20] = "V" // Verity
	}
	if flags&FS_NOCOMP_FL != 0 {
		attrs[21] = "m" // Dont_Compress
	}

	return strings.Join(attrs, "")
}
