package metadata

// Spacing between Key and Value while rendering
const keyValueSpacing = " "
const keyValueSpacingLen = 1

const fileStatErrorMsg = "Cannot load file stats"
const linkFileMsg = "This is a link file."
const linkFileBrokenMsg = "Link file is broken!"
const etFetchErrorMsg = "Errors while fetching metadata via exiftool"

const keyName = "Name"
const keySize = "Size"
const keyDataModified = "Date Modified"
const keyDataAccessed = "Date Accessed"
const keyPermissions = "Permissions"
const keyMd5Checksum = "MD5Checksum"

var sortPriority = map[string]int{ //nolint: gochecknoglobals // This is effectively const.
	keyName:         0,
	keySize:         1,
	keyDataModified: 2,
	keyDataAccessed: 3,
	keyPermissions:  4,
}
