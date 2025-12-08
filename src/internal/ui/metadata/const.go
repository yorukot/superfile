package metadata

// Spacing between Key and Value while rendering
const keyValueSpacing = " "
const keyValueSpacingLen = 1

const fileStatErrorMsg = "Cannot load file stats"
const linkFileBrokenMsg = "Link file is broken!"
const etFetchErrorMsg = "Errors while fetching metadata via exiftool"

const keyName = "Name"
const keySize = "Size"
const keyDataModified = "Date Modified"
const keyDataAccessed = "Date Accessed"
const keyPermissions = "Permissions"
const keyMd5Checksum = "MD5Checksum"
const keyOwner = "Owner"
const keyGroup = "Group"
const keyPath = "Path"

var sortPriority = map[string]int{ //nolint: gochecknoglobals // This is effectively const.
	// Metadata field priority indices for display ordering
	keyName:         0,
	keySize:         1,
	keyDataModified: 2,
	keyDataAccessed: 3, //nolint:mnd // display order index
	keyPermissions:  4, //nolint:mnd // display order index
	keyOwner:        5, //nolint:mnd // display order index
	keyGroup:        6, //nolint:mnd // display order index
	keyPath:         7, //nolint:mnd // display order index
}
