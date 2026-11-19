package find

const (
	findHeadlineText = "Find File/Folder"

	FindMinWidth  = 15
	FindMinHeight = 3

	// visibleOverhead is the number of non-result lines the modal box always
	// renders: top and bottom borders (2), the query input line (1), the
	// section divider below it (1), and the worst-case scroll indicators
	// (divider + 2 indicator lines). The visible results window is derived
	// from the modal's max height minus this overhead, so the modal fills
	// its allotted height on tall terminals instead of showing a fixed
	// handful of results.
	visibleOverhead = 7

	maxResults = 1000

	// UI dimension constants for find modal
	// markerColumnWidth is width reserved for type marker display (including padding and separator)
	markerColumnWidth = 8 // borders(2) + padding(2) + marker(1) + separator(3)

	// modalInputPadding is total padding for modal input fields
	modalInputPadding = 6 // 2 + 1 + 2 + 1 (borders and spacing)
)
