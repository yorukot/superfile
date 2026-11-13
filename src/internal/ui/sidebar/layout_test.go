package sidebar

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
)

// testSidebarWidth is wide enough for a directory name to be drawn in full.
const testSidebarWidth = 20

// withRenderableBorders gives the sidebar a real border and populated labels for
// the duration of the test.
//
// Border runes and prerendered labels live in package level config that is only
// filled in by the application's startup. Without them lipgloss draws no border
// at all, which would silently shift every row and make the cross-check below
// assert the wrong offsets.
func withRenderableBorders(t *testing.T) {
	t.Helper()

	previous := common.Config
	t.Cleanup(func() {
		common.Config = previous //nolint:reassign // Restoring global config after the test.
	})

	cfg := &common.Config
	cfg.BorderTop = "─"
	cfg.BorderBottom = "─"
	cfg.BorderLeft = "│"
	cfg.BorderRight = "│"
	cfg.BorderTopLeft = "╭"
	cfg.BorderTopRight = "╮"
	cfg.BorderBottomLeft = "╰"
	cfg.BorderBottomRight = "╯"
	cfg.BorderMiddleLeft = "│"
	cfg.BorderMiddleRight = "│"

	// The prerendered labels are normally built from theme styles, which need
	// the full startup sequence. Plain text is enough here and avoids pulling
	// theme initialisation into these tests.
	setPrerenderedLabels(t)
}

// setPrerenderedLabels swaps in plain text sidebar labels for the duration of
// the test and restores whatever was there before.
func setPrerenderedLabels(t *testing.T) {
	t.Helper()

	previous := map[string]string{
		"title":  common.SideBarSuperfileTitle,
		"home":   common.SideBarHomeDivider,
		"pinned": common.SideBarPinnedDivider,
		"disks":  common.SideBarDisksDivider,
	}
	t.Cleanup(func() {
		//nolint:reassign // Restoring global labels after the test.
		common.SideBarSuperfileTitle = previous["title"]
		//nolint:reassign // Restoring global labels after the test.
		common.SideBarHomeDivider = previous["home"]
		//nolint:reassign // Restoring global labels after the test.
		common.SideBarPinnedDivider = previous["pinned"]
		//nolint:reassign // Restoring global labels after the test.
		common.SideBarDisksDivider = previous["disks"]
	})

	//nolint:reassign // Test fixture, see the doc comment.
	common.SideBarSuperfileTitle = " superfile"
	//nolint:reassign // Test fixture, see the doc comment.
	common.SideBarHomeDivider = " Home"
	//nolint:reassign // Test fixture, see the doc comment.
	common.SideBarPinnedDivider = " Pinned"
	//nolint:reassign // Test fixture, see the doc comment.
	common.SideBarDisksDivider = " Disks"
}

// TestDirectoryIndexAtRowMatchesRender cross-checks the row lookup against the
// rows the sidebar actually draws.
//
// The lookup supplies the row and the rendered frame supplies the expected
// directory name, so a disagreement in either the search bar offset or the
// scroll offset shows up as a failure instead of silently mis-targeting clicks.
func TestDirectoryIndexAtRowMatchesRender(t *testing.T) {
	testCases := []struct {
		name           string
		sidebar        Model
		sidebarFocused bool
	}{
		{
			name:           "unfocused sidebar omits the search bar row",
			sidebar:        defaultTestModel(0, 0, 20, 5, 3, 2),
			sidebarFocused: false,
		},
		{
			name:           "focused sidebar adds the search bar row",
			sidebar:        defaultTestModel(0, 0, 20, 5, 3, 2),
			sidebarFocused: true,
		},
		{
			name:           "scrolled sidebar keeps lookup aligned",
			sidebar:        defaultTestModel(8, 5, 20, 5, 3, 2),
			sidebarFocused: true,
		},
		{
			name:           "sidebar too short to show every entry",
			sidebar:        defaultTestModel(0, 0, 9, 9, 4, 4),
			sidebarFocused: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withRenderableBorders(t)
			s := tc.sidebar
			// The shared helper only sets a height, but rendering needs a width
			// or every row comes out truncated and the comparison is vacuous.
			s.width = testSidebarWidth
			rows := strings.Split(ansi.Strip(s.Render(tc.sidebarFocused, "")), "\n")

			checked := 0
			for index, dir := range s.directories {
				if dir.isDivider() {
					continue
				}

				row := rowForDirectory(t, &s, tc.sidebarFocused, index)
				if row == -1 {
					// Entry did not fit on screen, nothing to compare against.
					continue
				}

				require.Less(t, row, len(rows))
				assert.Contains(t, rows[row], dir.Name,
					"directory %q reported at row %d, which shows %q",
					dir.Name, row, strings.TrimSpace(rows[row]))
				checked++
			}
			require.Positive(t, checked, "no directory was resolved to a row")
		})
	}
}

// TestDirectoryIndexAtRowIgnoresDividers makes sure section dividers, which
// occupy three rows each, are never reported as selectable entries.
func TestDirectoryIndexAtRowIgnoresDividers(t *testing.T) {
	withRenderableBorders(t)
	s := defaultTestModel(0, 0, 20, 4, 3, 2)

	for index, dir := range s.directories {
		if !dir.isDivider() {
			continue
		}
		for y := range s.GetHeight() {
			assert.NotEqual(t, index, s.DirectoryIndexAtRow(true, y),
				"divider at index %d must not resolve to a row", index)
		}
	}
}

// TestMouseRegionsCoverOnlySelectableEntries verifies the published regions skip
// dividers, so a click in a divider's blank rows falls through to nothing
// rather than selecting a neighbouring directory.
func TestMouseRegionsCoverOnlySelectableEntries(t *testing.T) {
	s := defaultTestModel(0, 0, 20, 4, 3, 2)
	s.width = testSidebarWidth

	reg := mouse.NewRegistry()
	s.MouseRegions(&reg, false, 0, 0)

	for y := range s.GetHeight() {
		target := reg.TargetAt(0, y)
		if target.Kind != mouse.TargetSidebarDirectory {
			continue
		}
		assert.False(t, s.directories[target.ItemIndex].isDivider(),
			"a region was published for divider index %d", target.ItemIndex)
	}
}

// rowForDirectory returns the first row at which the sidebar lookup reports the
// given directory index, or -1 when the directory has no row.
func rowForDirectory(t *testing.T, s *Model, sidebarFocused bool, index int) int {
	t.Helper()

	for y := range s.GetHeight() {
		if s.DirectoryIndexAtRow(sidebarFocused, y) == index {
			return y
		}
	}
	return -1
}
