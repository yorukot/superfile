// Package mouse maps terminal mouse coordinates onto superfile UI targets.
//
// Rendering in superfile produces plain strings, so nothing inherently ties a
// terminal cell to the widget drawn inside it. Each component instead publishes
// the layout it used while rendering, and the registry below collects those
// layouts so that a mouse event can be resolved to a widget. Keeping the layout
// next to the code that draws it is what stops hit-testing from drifting away
// from the rendered output.
package mouse

// NoIndex marks a target that is not associated with a single panel or item.
const NoIndex = -1

// TargetKind identifies the kind of widget occupying a registered region.
type TargetKind int

const (
	// TargetUnknown means no registered region contains the coordinates.
	TargetUnknown TargetKind = iota
	// TargetFilePanelItem is a file or directory row inside a file panel.
	TargetFilePanelItem
	// TargetSidebarDirectory is a directory entry in the sidebar.
	TargetSidebarDirectory
	// TargetProcessBarItem is a running operation in the process bar.
	TargetProcessBarItem
	// TargetMetadataItem is a row of the metadata panel.
	TargetMetadataItem
	// TargetContextMenuItem is an entry of the currently open context menu.
	TargetContextMenuItem
)

// Target identifies the widget under a set of coordinates.
type Target struct {
	// Kind is the type of widget the coordinates landed on.
	Kind TargetKind
	// PanelIndex is the file panel containing the coordinates, or NoIndex when
	// the target does not belong to a file panel.
	PanelIndex int
	// ItemIndex is the index of the item within its component's item list, or
	// NoIndex when the target is not tied to a single item.
	ItemIndex int
}

// FilePanelItemTarget returns the target of a row inside a file panel.
func FilePanelItemTarget(panelIndex, itemIndex int) Target {
	return Target{Kind: TargetFilePanelItem, PanelIndex: panelIndex, ItemIndex: itemIndex}
}

// SidebarDirectoryTarget returns the target of a sidebar directory entry.
func SidebarDirectoryTarget(itemIndex int) Target {
	return Target{Kind: TargetSidebarDirectory, PanelIndex: NoIndex, ItemIndex: itemIndex}
}

// ProcessBarItemTarget returns the target of a process bar row.
func ProcessBarItemTarget(itemIndex int) Target {
	return Target{Kind: TargetProcessBarItem, PanelIndex: NoIndex, ItemIndex: itemIndex}
}

// MetadataItemTarget returns the target of a metadata row.
func MetadataItemTarget(itemIndex int) Target {
	return Target{Kind: TargetMetadataItem, PanelIndex: NoIndex, ItemIndex: itemIndex}
}

// ContextMenuItemTarget returns the target of a context menu entry.
func ContextMenuItemTarget(itemIndex int) Target {
	return Target{Kind: TargetContextMenuItem, PanelIndex: NoIndex, ItemIndex: itemIndex}
}

// Region is a half-open rectangle of terminal cells occupied by one widget.
// X grows to the right and Y grows downwards, matching tea.Mouse coordinates.
type Region struct {
	Target

	// X0 is the inclusive leftmost cell.
	X0 int
	// Y0 is the inclusive topmost cell.
	Y0 int
	// X1 is the exclusive rightmost cell.
	X1 int
	// Y1 is the exclusive bottommost cell.
	Y1 int
}

// Contains reports whether the region covers the given terminal cell.
func (r Region) Contains(x, y int) bool {
	return x >= r.X0 && x < r.X1 && y >= r.Y0 && y < r.Y1
}

// Registry holds the regions published by components during a render pass.
type Registry struct {
	regions []Region
}

// NewRegistry returns an empty registry.
func NewRegistry() Registry {
	return Registry{}
}

// Reset drops every registered region. It is called at the start of each render
// pass so that regions never outlive the layout they were computed from.
func (r *Registry) Reset() {
	r.regions = r.regions[:0]
}

// Add registers a region for a widget spanning the given cells. Regions with no
// area are dropped, which happens naturally when a panel is too small to draw
// the item that would have occupied them.
func (r *Registry) Add(target Target, x0, y0, x1, y1 int) {
	if x1 <= x0 || y1 <= y0 {
		return
	}
	r.regions = append(r.regions, Region{Target: target, X0: x0, Y0: y0, X1: x1, Y1: y1})
}

// AddRow registers a single-line region starting at the given top-left cell.
func (r *Registry) AddRow(target Target, x0, y0, width, height int) {
	r.Add(target, x0, y0, x0+width, y0+height)
}

// TargetAt returns the target of the region covering the given cell. When no
// region covers it, an unknown target is returned. Regions added later take
// precedence, which is what lets an overlay such as a context menu claim cells
// that are also covered by the widget underneath it.
func (r *Registry) TargetAt(x, y int) Target {
	for i := len(r.regions) - 1; i >= 0; i-- {
		if r.regions[i].Contains(x, y) {
			return r.regions[i].Target
		}
	}
	return unknownTarget()
}

// IsEmpty reports whether no region is registered.
func (r *Registry) IsEmpty() bool {
	return len(r.regions) == 0
}

func unknownTarget() Target {
	return Target{Kind: TargetUnknown, PanelIndex: NoIndex, ItemIndex: NoIndex}
}
