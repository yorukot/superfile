package internal

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
	"github.com/yorukot/superfile/src/internal/ui/filemodel"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/quickconnect"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestQuickConnectHotkeyOpensModal(t *testing.T) {
	dir := filepath.Join(testDir, "TestQuickConnectHotkeyOpensModal")
	utils.SetupDirectories(t, dir)
	m := defaultTestModel(dir)

	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenQuickConnect[0]))

	assert.True(t, m.quickConnect.IsOpen())
	assert.True(t, m.IsOverlayModelOpen())
}

func TestRemoteStatusRendersInPanelHeaderAndSidebar(t *testing.T) {
	dir := filepath.Join(testDir, "TestRemoteStatusRendersInPanelHeaderAndSidebar")
	utils.SetupDirectories(t, dir)
	m := defaultTestModel(dir)
	location := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Label:     "ssh://e2e@127.0.0.1",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
	}
	m.fileModel.RegisterSession(filemodel.SessionState{
		ID:          location.SessionID,
		Provider:    location.Provider,
		Label:       location.Label,
		CurrentPath: location.Path,
		Status:      filemodel.SessionConnected,
	})
	require.NoError(t, m.fileModel.SetPaneLocation(m.fileModel.FocusedPanelIndex, location))

	panelRender := m.getFocusedFilePanel().Render(true)
	sidebarRender := m.sidebarRender()

	assert.Contains(t, panelRender, "sf-e2e:/tmp/sf-remote connected")
	assert.Contains(t, sidebarRender, "sf-e2e")
	assert.Contains(t, sidebarRender, "connected")
}

func TestQuickConnectBadStartPathKeepsExistingPaneState(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	dir := filepath.Join(testDir, "TestQuickConnectBadStartPathKeepsExistingPaneState")
	utils.SetupDirectories(t, dir)
	utils.SetupFilesWithData(t, []byte("local alpha\n"), filepath.Join(dir, "alpha.txt"))
	m := defaultTestModel(dir)
	m.fileModel.UpdateFilePanelsIfNeeded(true)

	originalLocation := m.getFocusedFilePanel().CurrentLocation()
	originalDisplay := m.getFocusedFilePanel().DisplayLocation()
	originalNames := panelElementNames(m.getFocusedFilePanel())

	profile := remoteOperationProfileForAlias(fixture, sshtest.AliasE2E)
	profile.StartPath = "/missing-start-path"
	m.quickConnect.SetKnownHostsPath(fixture.KnownHostsPath)
	m.quickConnect.SetTimeout(3 * time.Second)
	m.quickConnect.OpenWithProfiles([]common.SSHQuickConnectProfile{profile})

	action := m.quickConnect.ConnectSelected(context.Background())
	assert.Equal(t, quickconnect.ActionError, action.Type)
	assert.Nil(t, action.Session)
	assert.Equal(t, originalLocation, m.getFocusedFilePanel().CurrentLocation())
	assert.Equal(t, originalDisplay, m.getFocusedFilePanel().DisplayLocation())
	assert.Equal(t, originalNames, panelElementNames(m.getFocusedFilePanel()))
	_, ok := m.fileModel.Sessions[filesystem.SessionID(profile.Name)]
	assert.False(t, ok)
	assert.Contains(t, m.quickConnect.Warning(), "missing-start-path")
}

func panelElementNames(panel *filepanel.Model) []string {
	names := make([]string, 0, panel.ElemCount())
	for i := range panel.ElemCount() {
		names = append(names, panel.GetElementAtIdx(i).Name)
	}
	return names
}
