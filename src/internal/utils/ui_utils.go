package utils

// CntFooterPanels defines the maximum number of file panels
const CntFooterPanels = 3

// BorderPaddingForFooter is the border padding for footer calculations
const BorderPaddingForFooter = 2

// Including borders
func FullFooterHeight(footerHeight int, toggleFooter bool) int {
	if toggleFooter {
		return footerHeight + BorderPaddingForFooter
	}
	return 0
}
