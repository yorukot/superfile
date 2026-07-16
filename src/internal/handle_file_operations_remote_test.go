package internal

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	internalssh "github.com/yorukot/superfile/src/internal/ssh"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
	"github.com/yorukot/superfile/src/internal/ui/filemodel"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestRemoteCreateMkdirAndRenameUseProviders(t *testing.T) {
	m, fixture := newRemoteOnlyOperationModel(t)

	m.panelCreateNewFile()
	m.typingModal.textInput.SetValue("created.txt")
	m.createItem()
	assert.FileExists(t, remoteFixturePath(fixture, "/created.txt"))

	m.panelCreateNewFile()
	m.typingModal.textInput.SetValue("created-dir/")
	m.createItem()
	assert.DirExists(t, remoteFixturePath(fixture, "/created-dir"))

	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.AlphaPath)
	m.panelItemRename()
	m.getFocusedFilePanel().Rename.SetValue("alpha-renamed.txt")
	m.confirmRename()

	assert.NoFileExists(t, remoteFixturePath(fixture, fixture.AlphaPath))
	assert.FileExists(t, remoteFixturePath(fixture, "/alpha-renamed.txt"))
	assert.Equal(t, filesystem.RootRemotePath().String(), m.getFocusedFilePanel().Location)
}

func TestRemoteDeleteUsesConfirmationModal(t *testing.T) {
	m, fixture := newRemoteOnlyOperationModel(t)
	p := NewTestTeaProgWithEventLoop(t, m)

	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.AlphaPath)
	p.SendKey(common.Hotkeys.DeleteItems[0])

	assert.Eventually(t, p.getModel().notifyModel.IsOpen, DefaultTestTimeout, DefaultTestTick)
	assert.Equal(t, common.PermanentDeleteWarnTitle, p.getModel().notifyModel.GetTitle())
	assert.Equal(t, notify.PermanentDeleteAction, p.getModel().notifyModel.GetConfirmAction())

	p.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.Eventually(t, func() bool {
		_, err := os.Stat(remoteFixturePath(fixture, fixture.AlphaPath))
		return os.IsNotExist(err)
	}, DefaultTestTimeout, DefaultTestTick)
	assert.Equal(t, filesystem.RootRemotePath().String(), p.getModel().getFocusedFilePanel().Location)
}

func TestRemoteSameSessionCopyAndCutMove(t *testing.T) {
	m, fixture := newRemoteOnlyOperationModel(t)

	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.AlphaPath)
	m.copySingleItem(false)
	require.NoError(t, m.updateCurrentFilePanelDir(fixture.NestedPath))
	applyTeaCmd(t, m, m.getPasteItemCmd())
	assert.FileExists(t, remoteFixturePath(fixture, path.Join(fixture.NestedPath, "alpha.txt")))
	assert.NotEmpty(t, m.clipboard.GetLocations())

	require.NoError(t, m.updateCurrentFilePanelDir("/"))
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.BetaPath)
	m.copySingleItem(true)
	require.NoError(t, m.updateCurrentFilePanelDir(fixture.NestedPath))
	applyTeaCmd(t, m, m.getPasteItemCmd())

	assert.FileExists(t, remoteFixturePath(fixture, path.Join(fixture.NestedPath, "beta.txt")))
	assert.NoFileExists(t, remoteFixturePath(fixture, fixture.BetaPath))
	assert.Empty(t, m.clipboard.GetLocations())
}

func TestRemoteUploadAndDownloadUseTransferEngine(t *testing.T) {
	localDir := t.TempDir()
	downloadDir := filepath.Join(localDir, "downloads")
	uploadPath := filepath.Join(localDir, "upload.txt")
	utils.SetupDirectories(t, downloadDir)
	utils.SetupFilesWithData(t, []byte("upload body\n"), uploadPath)

	m, fixture := newLocalRemoteOperationModel(t, localDir)

	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), uploadPath)
	m.copySingleItem(false)
	m.fileModel.NextFilePanel()
	applyTeaCmd(t, m, m.getPasteItemCmd())
	assert.FileExists(t, remoteFixturePath(fixture, "/upload.txt"))
	assert.Equal(t, "upload body\n", string(mustReadFile(t, remoteFixturePath(fixture, "/upload.txt"))))

	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.BetaPath)
	m.copySingleItem(false)
	m.fileModel.PreviousFilePanel()
	require.NoError(t, m.updateCurrentFilePanelDir(downloadDir))
	applyTeaCmd(t, m, m.getPasteItemCmd())
	assert.FileExists(t, filepath.Join(downloadDir, "beta.txt"))
	assert.Equal(t, "beta\n", string(mustReadFile(t, filepath.Join(downloadDir, "beta.txt"))))
}

func TestRemoteUploadCleanupUsesFreshSessionAfterDisconnect(t *testing.T) {
	localDir := t.TempDir()
	uploadPath := filepath.Join(localDir, "large.bin")
	utils.SetupFilesWithData(t, bytes.Repeat([]byte("superfile-transfer"), 65536), uploadPath)

	m, fixture := newLocalRemoteOperationModel(t, localDir)
	fixture.DisconnectOnAnyWriteOnce(1)

	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), uploadPath)
	m.copySingleItem(false)
	m.fileModel.NextFilePanel()
	applyTeaCmd(t, m, m.getPasteItemCmd())

	assert.NoFileExists(t, remoteFixturePath(fixture, "/large.bin"))
	assertNoRemoteTransferTemps(t, fixture)
}

func TestRemoteUnsupportedOperationsShowExactMessageAndPreserveState(t *testing.T) {
	m, fixture := newRemoteOnlyOperationModel(t)
	panel := m.getFocusedFilePanel()
	setFilePanelSelectedItemByLocation(t, panel, fixture.AlphaPath)
	panel.ChangeFilePanelMode()
	panel.SetSelected(fixture.AlphaPath)
	selectedBefore := panel.GetSelectedLocationsSortedAsVisible()
	pathBefore := panel.Location

	tests := []struct {
		name      string
		operation filesystem.Operation
		run       func(*model) tea.Cmd
	}{
		{
			name:      "compress",
			operation: filesystem.OperationCompress,
			run:       func(m *model) tea.Cmd { return m.getCompressSelectedFilesCmd() },
		},
		{
			name:      "extract",
			operation: filesystem.OperationExtract,
			run:       func(m *model) tea.Cmd { return m.getExtractFileCmd() },
		},
		{
			name:      "open-with",
			operation: filesystem.OperationOpenWith,
			run:       func(m *model) tea.Cmd { return m.openFileWithEditor() },
		},
		{
			name:      "zoxide",
			operation: filesystem.OperationZoxide,
			run:       func(m *model) tea.Cmd { return m.mainKey(common.Hotkeys.OpenZoxide[0]) },
		},
		{
			name:      "remote-shell",
			operation: filesystem.OperationRemoteShell,
			run:       func(m *model) tea.Cmd { return m.mainKey(common.Hotkeys.OpenCommandLine[0]) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyTeaCmd(t, m, tt.run(m))
			assert.True(t, m.notifyModel.IsOpen())
			assert.Equal(
				t,
				remoteUnsupportedOperationText(filesystem.ProviderSFTP, tt.operation),
				m.notifyModel.GetContent(),
			)
			assert.Equal(t, pathBefore, panel.Location)
			assert.Equal(t, selectedBefore, panel.GetSelectedLocationsSortedAsVisible())
			m.notifyModel.Close()
		})
	}

	m.disableMetadata = false
	m.focusPanel = metadataFocus
	cmd := m.getMetadataCmd()
	assert.Nil(t, cmd)
	m.fileMetaData.SetDimensions(120, 8)
	assert.Contains(
		t,
		stripANSI(m.fileMetaData.Render(true)),
		remoteUnsupportedOperationText(filesystem.ProviderSFTP, filesystem.OperationMetadata),
	)
	assert.Equal(t, pathBefore, panel.Location)
	assert.Equal(t, selectedBefore, panel.GetSelectedLocationsSortedAsVisible())
}

func newRemoteOnlyOperationModel(t *testing.T) (*model, *sshtest.Fixture) {
	t.Helper()
	localDir := t.TempDir()
	m, fixture := newLocalRemoteOperationModelWithIndex(t, []string{localDir}, 0)
	m.fileModel.FocusedPanelIndex = 0
	m.fileModel.FilePanels[0].IsFocused = true
	return m, fixture
}

func newLocalRemoteOperationModel(t *testing.T, localDir string) (*model, *sshtest.Fixture) {
	t.Helper()
	return newLocalRemoteOperationModelWithIndex(t, []string{localDir, localDir}, 1)
}

func newLocalRemoteOperationModelWithIndex(t *testing.T, dirs []string, remoteIndex int) (*model, *sshtest.Fixture) {
	t.Helper()
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	m := defaultTestModel(dirs...)
	registerRemoteOperationPanel(t, m, fixture, remoteIndex)
	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].IsFocused = i == m.fileModel.FocusedPanelIndex
	}
	return m, fixture
}

func registerRemoteOperationPanel(t *testing.T, m *model, fixture *sshtest.Fixture, panelIndex int) {
	t.Helper()
	provider := filesystem.NewSFTPProvider(internalssh.ClientConfigRequest{
		Profile:        remoteOperationProfileForAlias(fixture, sshtest.AliasE2E),
		KnownHostsPath: fixture.KnownHostsPath,
		HostKeyAlias:   sshtest.AliasE2E,
	})
	location := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: filesystem.SessionID(sshtest.AliasE2E),
		Path:      filesystem.RootRemotePath(),
		Label:     sshtest.AliasE2E,
	}
	session, err := provider.Open(context.Background(), location)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = session.Close()
	})

	m.fileModel.RegisterSession(filemodel.SessionState{
		ID:          location.SessionID,
		Provider:    location.Provider,
		Label:       location.Label,
		CurrentPath: location.Path,
		Status:      filemodel.SessionConnected,
		Browser:     session,
		Reconnect: func(ctx context.Context, loc filesystem.Location) (filesystem.Session, error) {
			return filesystem.NewSFTPProvider(internalssh.ClientConfigRequest{
				Profile:        remoteOperationProfileForAlias(fixture, sshtest.AliasE2E),
				KnownHostsPath: fixture.KnownHostsPath,
				HostKeyAlias:   sshtest.AliasE2E,
			}).Open(ctx, loc)
		},
	})
	require.NoError(t, m.fileModel.SetPaneLocation(panelIndex, location))
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	m.sessionRegistry = m.fileModel.Sessions
}

func remoteOperationProfileForAlias(fixture *sshtest.Fixture, aliasName string) common.SSHQuickConnectProfile {
	alias := fixture.Aliases[aliasName]
	profile := common.SSHQuickConnectProfile{
		Name:           alias.Name,
		Host:           alias.Host,
		Port:           alias.Port,
		User:           alias.User,
		StartPath:      "/",
		IdentityFile:   alias.IdentityFilePath,
		IdentityFiles:  nil,
		IdentitiesOnly: true,
		AuthOrder:      []string{common.SSHAuthMethodPublicKey},
	}
	if alias.IdentityFilePath != "" {
		profile.IdentityFiles = []string{alias.IdentityFilePath}
	}
	return profile
}

func remoteFixturePath(fixture *sshtest.Fixture, remotePath string) string {
	clean := strings.TrimPrefix(filesystem.NewRemotePath(remotePath).String(), "/")
	return filepath.Join(fixture.RemoteRootPath, filepath.FromSlash(clean))
}

func applyTeaCmd(t *testing.T, m *model, cmd tea.Cmd) {
	t.Helper()
	require.NotNil(t, cmd)
	msg := cmd()
	TeaUpdate(m, msg)
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func stripANSI(input string) string {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansi.ReplaceAllString(input, "")
}

func assertNoRemoteTransferTemps(t *testing.T, fixture *sshtest.Fixture) {
	t.Helper()
	entries, err := os.ReadDir(fixture.RemoteRootPath)
	require.NoError(t, err)
	for _, entry := range entries {
		assert.NotContains(t, entry.Name(), ".superfile-transfer-")
	}
}
