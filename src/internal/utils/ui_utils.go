package utils

const CntFooterPanels = 3

const BorderPaddingForFooter = 2

// analog of unix.Winsize, so uint16 used
type Winsize struct {
	Row uint16
	Col uint16
}

func safeIntToUint16(value int) uint16 {
	const maxUint16 = 0xFFFF // 65535
	if value < 0 {
		return 0
	}
	if value > maxUint16 {
		return maxUint16
	}
	return uint16(value)
}
func NewWinSize(rowCount int, colCount int) *Winsize {
	return &Winsize{Row: safeIntToUint16(rowCount), Col: safeIntToUint16(colCount)}
}

// Including borders
func FullFooterHeight(footerHeight int, toggleFooter bool) int {
	if toggleFooter {
		return footerHeight + BorderPaddingForFooter
	}
	return 0
}
