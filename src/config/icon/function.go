package icon

func InitIcon(nerdfont bool) {
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
}
