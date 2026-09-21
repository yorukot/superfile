package sidebar

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"charm.land/bubbles/v2/textinput"
	"github.com/adrg/xdg"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

//nolint:reassign // Isolate global settings and restore them with t.Cleanup.
func TestStandardFolderShortcutLocations(
	t *testing.T,
) {
	originalHome, originalDirs, originalTrash := xdg.Home, xdg.UserDirs, variable.LinuxTrashDirectory
	t.Cleanup(func() {
		xdg.Home, xdg.UserDirs, variable.LinuxTrashDirectory = originalHome, originalDirs, originalTrash
	})
	paths := map[string]*string{
		common.FolderHome: &xdg.Home, "Desktop": &xdg.UserDirs.Desktop, common.FolderDownloads: &xdg.UserDirs.Download,
		"Documents": &xdg.UserDirs.Documents, "Pictures": &xdg.UserDirs.Pictures, "Music": &xdg.UserDirs.Music,
		"Videos": &xdg.UserDirs.Videos, "Templates": &xdg.UserDirs.Templates,
		"PublicShare": &xdg.UserDirs.PublicShare, common.FolderTrash: &variable.LinuxTrashDirectory,
	}
	s := Model{disabled: true}
	for name, path := range paths {
		*path = t.TempDir()
		if name == common.FolderTrash && runtime.GOOS != utils.OsLinux {
			continue
		}
		assert.Equal(t, *path, s.FolderShortcutLocation(common.FolderShortcut{Name: name}), name)
	}
	xdg.UserDirs.Download = filepath.Join(t.TempDir(), "missing")
	assert.Empty(t, s.FolderShortcutLocation(common.FolderShortcut{Name: common.FolderDownloads}))
}

func TestFolderShortcutsFollowFilteredPins(t *testing.T) {
	mgr, s := sidebarWithPinnedDir(t)
	root := t.TempDir()
	alpha, beta := filepath.Join(root, "alpha"), filepath.Join(root, "beta")
	utils.SetupDirectories(t, alpha, beta)
	require.NoError(t, mgr.Save([]directory{{Location: alpha, Name: "alpha"}, {Location: beta, Name: "beta"}}))
	s.searchBar = textinput.New()
	s.searchBar.SetValue("beta")
	s.UpdateDirectories()
	require.Len(t, s.directories, 1)
	assert.Equal(t, beta, s.FolderShortcutLocation(common.FolderShortcut{Pin: 1}))
	assert.Empty(t, s.FolderShortcutLocation(common.FolderShortcut{Pin: 2}))
	s.searchBar.SetValue("")
	s.UpdateDirectories()
	s.renderIndex = 1
	assert.Equal(t, alpha, s.FolderShortcutLocation(common.FolderShortcut{Pin: 1}))
	assert.Equal(t, beta, s.FolderShortcutLocation(common.FolderShortcut{Pin: 2}))
}

//nolint:reassign // Isolate global settings and restore them with t.Cleanup.
func TestFolderHintsIdentifyRowsAndKeepPinNumbersWhileScrolling(
	t *testing.T,
) {
	original := common.Hotkeys
	common.Hotkeys = common.HotkeysType{
		GoToHome: []string{"alt+h"}, GoToDownloads: []string{"alt+d"},
		GoToPin1: []string{"alt+1"}, GoToPin2: []string{"alt+2"},
	}
	t.Cleanup(func() { common.Hotkeys = original })
	s := Model{directories: []directory{
		{Name: common.FolderHome, Location: "/same", Section: utils.SidebarSectionHome},
		{Name: common.FolderDownloads, Location: "/same", Section: utils.SidebarSectionHome},
		pinnedDividerDir,
		{Name: common.FolderHome, Location: "/same", Section: utils.SidebarSectionPinned},
		{Name: "second", Location: "/other", Section: utils.SidebarSectionPinned},
	}, renderIndex: 4}
	assert.Equal(t, map[int]string{0: "H", 1: "D", 3: "1", 4: "2"}, s.folderHints())
	common.Hotkeys.GoToPin2 = nil
	assert.Empty(t, s.folderHints()[4])
}

func TestFolderHintsFitSidebarWidth(t *testing.T) {
	originalNerdfont := common.Config.Nerdfont
	t.Cleanup(func() { common.Config.Nerdfont = originalNerdfont })
	for _, nerdfont := range []bool{false, true} {
		common.Config.Nerdfont = nerdfont
		expectedBadge := " [H]"
		if nerdfont {
			expectedBadge = " \ue0b6H\ue0b4"
		}
		for _, width := range []int{5, 10, 20} {
			for _, name := range []string{common.FolderHome, "a very long folder name", "文件資料夾", "\x1b[31mColored folder\x1b[0m"} {
				line := appendFolderHint(name, "H", width)
				assert.Equal(t, width, ansi.StringWidth(line))
				assert.True(t, strings.HasSuffix(ansi.Strip(line), expectedBadge))
			}
		}
	}
	assert.Equal(t, "folder", appendFolderHint("folder", "", 20))
	assert.Equal(t, "folder", appendFolderHint("folder", "shift+f12", 5))
}
