package utils

// CntFooterPanels defines the maximum number of file panels
const CntFooterPanels = 3

// BorderPaddingForFooter is the border padding for footer calculations
const BorderPaddingForFooter = 2

// We have three panels, so 6 characters for border
// <---><---><--->
// Hence we have (fullWidth - 6) / 3 = fullWidth/3 - 2
func FooterWidth(fullWidth int) int {
	return fullWidth/CntFooterPanels - BorderPaddingForFooter
}

// Including borders
func FullFooterHeight(footerHeight int, toggleFooter bool) int {
	if toggleFooter {
		return footerHeight + BorderPaddingForFooter
	}
	return 0
}
