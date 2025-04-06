package icon

func InitIcon(nerdfont bool, directoryIconColor string) {
	if !nerdfont {
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
		Browser = ""
		Select = ""
		Error = ""
		Warn = ""
		Done = ""
		InOperation = ""
		Directory = ""
		Search = ""
		SortAsc = ""
		SortDesc = ""
	}

	if directoryIconColor == "" {
		directoryIconColor = "NONE" // Dark yellowish
	}
	Folders["folder"] = Style{
		Icon:  "\uf07b", // Printable Rune : ""
		Color: directoryIconColor,
	}
}
