package common

type Directory struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

// ================ Sidebar related utils =====================
// Hopefully compiler inlines it
func (d Directory) IsDivider() bool {
	return d == PinnedDividerDir || d == DiskDividerDir
}
func (d Directory) RequiredHeight() int {
	if d.IsDivider() {
		return 3
	}
	return 1
}
