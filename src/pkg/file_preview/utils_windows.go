//go:build windows
// +build windows

package filepreview

// getTerminalCellSizeViaIoctl is not supported on Windows, so always return false
func getTerminalCellSizeViaIoctl() (TerminalCellSize, bool) {
	return TerminalCellSize{}, false
}
