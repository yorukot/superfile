package utils

const (
	TrueString  = "true"
	FalseString = "false"
	// These are used while comparing with runtime.GOOS
	// OsWindows represents the Windows operating system identifier
	OsWindows = "windows"
	// OsDarwin represents the macOS (Darwin) operating system identifier
	OsDarwin = "darwin"
	OsLinux  = "linux"

	// File permissions
	ConfigFilePerm = 0600 // configuration files (owner read/write only)
	UserFilePerm   = 0644 // user-created files (owner rw, others r)
	LogFilePerm    = 0600 // log files (owner read/write only)

	// Directory permissions
	ConfigDirPerm = 0700 // configuration directories (owner only)
	UserDirPerm   = 0755 // user-created directories (owner rwx, others rx)

	// Extracted file permissions (from archives)
	ExtractedFileMode = 0644 // extracted files
	ExtractedDirMode  = 0755 // extracted directories

	// Sidebar sections
	SidebarSectionHome   = "home"
	SidebarSectionPinned = "pinned"
	SidebarSectionDisks  = "disks"
)
