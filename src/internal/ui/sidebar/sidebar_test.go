package sidebar

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestPinnedItemRenameUsesStoredNameWithoutVisualIcon(t *testing.T) {
	pinnedMgr, sidebar := sidebarWithPinnedDir(t)
	require.NoError(t, pinnedMgr.Save([]directory{{
		Location: sidebar.directories[0].Location,
		Name:     "project",
	}}))

	sidebar.directories = formDirctorySlice(
		nil,
		getPinnedDirectoriesWithIcon(&pinnedMgr),
		nil,
		[]string{utils.SidebarSectionPinned},
	)
	require.Len(t, sidebar.directories, 1)
	visualIcon := sidebar.directories[0].Icon
	require.NotEmpty(t, visualIcon)

	sidebar.PinnedItemRename()

	assert.True(t, sidebar.renaming)
	assert.Equal(t, "project", sidebar.rename.Value())
	assert.NotContains(t, sidebar.rename.Value(), visualIcon)
	assert.NotContains(t, sidebar.rename.Value(), icon.Space+"project")
}

func TestConfirmSidebarRenameDoesNotPersistVisualIcon(t *testing.T) {
	pinnedMgr, sidebar := sidebarWithPinnedDir(t)
	require.NoError(t, pinnedMgr.Save([]directory{{
		Location: sidebar.directories[0].Location,
		Name:     "project",
	}}))

	sidebar.directories = formDirctorySlice(
		nil,
		getPinnedDirectoriesWithIcon(&pinnedMgr),
		nil,
		[]string{utils.SidebarSectionPinned},
	)
	require.Len(t, sidebar.directories, 1)
	visualIcon := sidebar.directories[0].Icon
	require.NotEmpty(t, visualIcon)

	const replacementName = "renamed-project"

	sidebar.PinnedItemRename()
	sidebar.rename.SetValue(replacementName)
	sidebar.ConfirmSidebarRename()

	pinnedDirs := pinnedMgr.Load()
	require.Len(t, pinnedDirs, 1)
	assert.Equal(t, replacementName, pinnedDirs[0].Name)
	assert.NotContains(t, pinnedDirs[0].Name, visualIcon)
	assert.NotContains(t, pinnedDirs[0].Name, icon.Space+replacementName)
}

func sidebarWithPinnedDir(t *testing.T) (PinnedManager, Model) {
	t.Helper()

	originalNerdfont := common.Config.Nerdfont
	common.Config.Nerdfont = true
	t.Cleanup(func() {
		common.Config.Nerdfont = originalNerdfont
	})

	tempDir := t.TempDir()
	pinnedDir := filepath.Join(tempDir, "project")
	utils.SetupDirectories(t, pinnedDir)

	pinnedMgr := PinnedManager{filePath: filepath.Join(tempDir, "pinned.json")}
	utils.SetupFilesWithData(t, []byte("[]"), pinnedMgr.filePath)

	return pinnedMgr, Model{
		directories: []directory{{
			Location: pinnedDir,
			Name:     "project",
			Section:  utils.SidebarSectionPinned,
		}},
		cursor:    0,
		pinnedMgr: &pinnedMgr,
		sections:  []string{utils.SidebarSectionPinned},
	}
}

// A sidebar refresh re-reads the pinned file, stats the well known directories
// and enumerates mounts. Doing that on every message made every keystroke wait
// on the filesystem, so it is throttled - but what the user just did still has
// to show up at once.
func TestUpdateDirectoriesIfNeededThrottlesRefresh(t *testing.T) {
	pinnedMgr, sidebar, dirA, dirB := sidebarWithPinnedSection(t)
	require.NoError(t, pinnedMgr.Save([]directory{{Location: dirA, Name: "alpha"}}))

	sidebar.UpdateDirectoriesIfNeeded(false)
	require.Equal(t, []string{"alpha"}, pinnedNames(&sidebar))

	// Another superfile instance pins a second directory
	require.NoError(t, pinnedMgr.Save([]directory{
		{Location: dirA, Name: "alpha"},
		{Location: dirB, Name: "beta"},
	}))

	sidebar.UpdateDirectoriesIfNeeded(false)
	assert.Equal(t, []string{"alpha"}, pinnedNames(&sidebar),
		"must not re-read the pinned file on every update")

	// force is for what the user just did - pinning, renaming - which cannot wait
	// out the interval
	sidebar.UpdateDirectoriesIfNeeded(true)
	assert.Equal(t, []string{"alpha", "beta"}, pinnedNames(&sidebar))

	require.NoError(t, pinnedMgr.Save([]directory{{Location: dirB, Name: "beta"}}))
	sidebar.UpdateDirectoriesIfNeeded(false)
	require.Equal(t, []string{"alpha", "beta"}, pinnedNames(&sidebar),
		"the force must have restarted the interval")

	// Past the interval, outside changes are picked up
	sidebar.lastRefresh = time.Now().Add(-2 * sidebarRefreshInterval)
	sidebar.UpdateDirectoriesIfNeeded(false)
	assert.Equal(t, []string{"beta"}, pinnedNames(&sidebar))
}

// A changed search query is something the sidebar cannot render without, so it
// is applied at once rather than at the next refresh.
func TestUpdateDirectoriesIfNeededAppliesQueryChangeAtOnce(t *testing.T) {
	pinnedMgr, sidebar, dirA, dirB := sidebarWithPinnedSection(t)
	require.NoError(t, pinnedMgr.Save([]directory{
		{Location: dirA, Name: "alpha"},
		{Location: dirB, Name: "beta"},
	}))

	sidebar.UpdateDirectoriesIfNeeded(false)
	require.Equal(t, []string{"alpha", "beta"}, pinnedNames(&sidebar))

	// Well within the interval, and not forced
	sidebar.searchBar.SetValue("alp")
	sidebar.UpdateDirectoriesIfNeeded(false)
	assert.Equal(t, []string{"alpha"}, pinnedNames(&sidebar))

	sidebar.searchBar.SetValue("")
	sidebar.UpdateDirectoriesIfNeeded(false)
	assert.Equal(t, []string{"alpha", "beta"}, pinnedNames(&sidebar))
}

// sidebarWithPinnedSection returns a sidebar showing the pinned section only, so
// what it displays depends solely on the pinned file, along with that file's
// manager and two directories available to pin.
func sidebarWithPinnedSection(t *testing.T) (PinnedManager, Model, string, string) {
	t.Helper()

	tempDir := t.TempDir()
	dirA := filepath.Join(tempDir, "alpha")
	dirB := filepath.Join(tempDir, "beta")
	utils.SetupDirectories(t, dirA, dirB)

	pinnedMgr := PinnedManager{filePath: filepath.Join(tempDir, "pinned.json")}
	utils.SetupFilesWithData(t, []byte("[]"), pinnedMgr.filePath)

	return pinnedMgr, Model{
		pinnedMgr: &pinnedMgr,
		sections:  []string{utils.SidebarSectionPinned},
	}, dirA, dirB
}

// pinnedNames returns the names of the actual directories the sidebar is showing,
// skipping section dividers.
func pinnedNames(s *Model) []string {
	names := []string{}
	for _, d := range s.directories {
		if !d.isDivider() {
			names = append(names, d.Name)
		}
	}
	return names
}
