package utils

// We have three panels, so 6 characters for border
// <---><---><--->
// Hence we have (fullWidth - 6) / 3 = fullWidth/3 - 2
func FooterWidth(fullWidth int) int {
	return fullWidth/3 - 2
}

// Including borders
func FullFooterHeight(footerHeight int, toggleFooter bool) int {
	if toggleFooter {
		return footerHeight + 2
	}
	return 0
}
