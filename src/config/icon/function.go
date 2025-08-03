package icon

// InitIcon initializes the icon configuration for the application.
// It sets up different icons based on whether nerd fonts are enabled and configures directory icon colors.
//
// Parameters:
//   - nerdfont: boolean flag to determine if nerd fonts should be used
//     When false, uses simple ASCII characters for icons
//     When true, uses nerd font icons (default behavior)
//   - directoryIconColor: string representing the color for directory icons
//     If empty, defaults to "NONE" (dark yellowish)
//
// The function configures various icons for:
//   - System directories (Home, Download, Documents, etc.)
//   - File operations (Compress, Extract, Copy, Cut, Delete)
//   - UI elements (Cursor, Browser, Select, etc.)
//   - Status indicators (Error, Warn, Done, InOperation)
//   - Navigation and sorting (Directory, Search, SortAsc, SortDesc)
func InitIcon(nerdfont bool, directoryIconColor string) {
	// Make sure that these alternatives are ASCII characters only.
	// Dont place any special unicode characters here.
	if !nerdfont {
		// When nerdfont is disabled, we use simple ASCII characters
		// Space is set to empty string because we don't need special spacing
		// for ASCII characters, unlike nerd fonts which often need proper spacing
		// to display correctly
		Space = ""
		SuperfileIcon = ""

		Home = ""
		Download = ""
		Documents = ""
		Pictures = ""
		Videos = ""
		Music = ""
		Templates = ""
		PublicShare = ""

		// file operations
		CompressFile = ""
		ExtractFile = ""
		Copy = ""
		Cut = ""
		Delete = ""

		// other
		Cursor = ">"
		Browser = "B"
		Select = "S"
		Error = ""
		Warn = ""
		Done = ""
		InOperation = ""
		Directory = ""
		Search = ""
		SortAsc = "^"
		SortDesc = "v"
		Terminal = ""
		Pinned = ""
		Disk = ""
	}

	if directoryIconColor == "" {
		directoryIconColor = "NONE" // Dark yellowish
	}
	Folders["folder"] = Style{
		Icon:  "\uf07b", // Printable Rune : ""
		Color: directoryIconColor,
	}
}

func GetCopyOrCutIcon(cut bool) string {
	if cut {
		return Cut
	}
	return Copy
}
