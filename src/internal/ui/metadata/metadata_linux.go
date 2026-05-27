//go:build linux

package metadata

import (
	"os"

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

type attrEntry struct {
	flag   uint32
	letter byte
}

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
	// list of letter and output order see in
	// lsattr sources https://github.com/tytso/e2fsprogs/blob/master/lib/e2p/pf.c
	var attrOrder = [...]attrEntry{
		{FS_SECRM_FL, 's'},        // Secure_Deletion
		{FS_UNRM_FL, 'u'},         // Undelete
		{FS_SYNC_FL, 'S'},         // Synchronous_Updates
		{FS_DIRSYNC_FL, 'D'},      // Synchronous_Directory_Updates
		{FS_IMMUTABLE_FL, 'i'},    // Immutable
		{FS_APPEND_FL, 'a'},       // Append_Only
		{FS_NODUMP_FL, 'd'},       // No_Dump
		{FS_NOATIME_FL, 'A'},      // No_Atime
		{FS_COMPR_FL, 'c'},        // Compression_Requested
		{FS_ENCRYPT_FL, 'E'},      // Encrypted
		{FS_JOURNAL_DATA_FL, 'j'}, // Journaled_Data
		{FS_INDEX_FL, 'I'},        // Indexed_directory
		{FS_NOTAIL_FL, 't'},       // No_Tailmerging
		{FS_TOPDIR_FL, 'T'},       // Top_of_Directory_Hierarchies
		{FS_EXTENT_FL, 'e'},       // Extents
		{FS_NOCOW_FL, 'C'},        // No_COW
		{FS_DAX_FL, 'x'},          // DAX
		{FS_CASEFOLD_FL, 'F'},     // Casefold
		{FS_INLINE_DATA_FL, 'N'},  // Inline_Data
		{FS_PROJINHERIT_FL, 'P'},  // Project_Hierarchy
		{FS_VERITY_FL, 'V'},       // Verity
		{FS_NOCOMP_FL, 'm'},       // Dont_Compress
	}

	attrs := make([]byte, len(attrOrder))
	for i := range attrs {
		attrs[i] = '-'
	}

	for i, e := range attrOrder {
		if flags&e.flag != 0 {
			attrs[i] = e.letter
		}
	}

	return string(attrs)
}
