package sidebar

import (
	"path/filepath"
	"testing"

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

	sidebar.PinnedItemRename()
	sidebar.ConfirmSidebarRename()

	pinnedDirs := pinnedMgr.Load()
	require.Len(t, pinnedDirs, 1)
	assert.Equal(t, "project", pinnedDirs[0].Name)
	assert.NotContains(t, pinnedDirs[0].Name, visualIcon)
	assert.NotContains(t, pinnedDirs[0].Name, icon.Space+"project")
}

func sidebarWithPinnedDir(t *testing.T) (PinnedManager, Model) {
	t.Helper()

	originalConfig := common.Config
	common.Config.Nerdfont = true
	t.Cleanup(func() {
		common.Config = originalConfig
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
