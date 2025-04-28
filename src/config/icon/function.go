package icon

func InitIcon(nerdfont bool, directoryIconColor string) {
	// Make sure that these alternatives are ASCII characters only.
	// Dont place any special unicode characters here.
	if !nerdfont {
		// Do we need this to be empty ? Maybe it should just be normal space ?
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
