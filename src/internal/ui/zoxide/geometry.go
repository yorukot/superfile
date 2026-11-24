package zoxide

// CalculateModalDimensions computes the width and maxHeight for the zoxide modal
// from the total terminal dimensions and configured percentage values.
// Percentage values outside [MinModalWidthPercent, MaxModalWidthPercent] default to
// DefaultModalWidthPercent and DefaultModalHeightPercent (50%).
// Resulting dimensions are clamped to ZoxideMinWidth and ZoxideMinHeight, and bounded
// by fullWidth and fullHeight when the terminal meets the minimum size requirements.
func CalculateModalDimensions(fullWidth, fullHeight, widthPercent, heightPercent int) (int, int) {
	if widthPercent < MinModalWidthPercent || widthPercent > MaxModalWidthPercent {
		widthPercent = DefaultModalWidthPercent
	}
	if heightPercent < MinModalHeightPercent || heightPercent > MaxModalHeightPercent {
		heightPercent = DefaultModalHeightPercent
	}

	width := fullWidth * widthPercent / PercentDivisor
	maxHeight := fullHeight * heightPercent / PercentDivisor

	if width < ZoxideMinWidth {
		width = ZoxideMinWidth
	}
	if fullWidth >= ZoxideMinWidth && width > fullWidth {
		width = fullWidth
	}

	if maxHeight < ZoxideMinHeight {
		maxHeight = ZoxideMinHeight
	}
	if fullHeight >= ZoxideMinHeight && maxHeight > fullHeight {
		maxHeight = fullHeight
	}

	return width, maxHeight
}
