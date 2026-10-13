package zoxide

const (
	zoxideHeadlineText = "Zoxide Navigation"

	ZoxideMinWidth  = 15
	ZoxideMinHeight = 3

	// renderOverhead is the number of non-result lines the modal box always
	// renders: top and bottom borders (2), the query input line (1), the
	// section divider below it (1), and the worst-case scroll indicators
	// (divider + 2 indicator lines). The visible results window is derived
	// from the modal's max height minus this overhead, so the modal fills
	// its allotted height on tall terminals instead of showing a fixed
	// handful of results.
	renderOverhead = 7

	// UI dimension constants for zoxide modal
	// scoreColumnWidth is width reserved for score display (including padding and separator)
	scoreColumnWidth = 13 // borders(2) + padding(2) + score(6) + separator(3)

	// modalInputPadding is total padding for modal input fields
	modalInputPadding = 6 // 2 + 1 + 2 + 1 (borders and spacing)
)
