package internal

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	internalssh "github.com/yorukot/superfile/src/internal/ssh"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
	"github.com/yorukot/superfile/src/internal/ui/filemodel"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/quickconnect"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestSSHQuickConnectCase(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)

	connectModel := quickconnect.New()
	connectModel.SetDiscoveryOptions(common.SSHConfigDiscoveryOptions{
		UserConfigPath:   fixture.SSHConfigPath,
		SystemConfigPath: filepath.Join(t.TempDir(), "missing_system_config"),
	})
	connectModel.SetKnownHostsPath(fixture.KnownHostsPath)
	connectModel.SetTimeout(3 * time.Second)

	cfg := &common.ConfigType{SSH: common.SSHConfigSection{Profiles: []common.SSHProfileType{{
		Name:      sshtest.AliasE2E,
		StartPath: "/",
	}}}}
	require.NoError(t, connectModel.Open(cfg))

	profile := quickConnectProfileByName(t, connectModel.Profiles(), sshtest.AliasE2E)
	connectModel.OpenWithProfiles([]common.SSHQuickConnectProfile{profile})

	action := connectModel.ConnectSelected(context.Background())
	require.Equal(t, quickconnect.ActionConnected, action.Type)
	require.NotNil(t, action.Session)
	t.Cleanup(func() {
		require.NoError(t, action.Session.Close())
	})

	localDir := t.TempDir()
	downloadDir := filepath.Join(localDir, "downloads")
	utils.SetupDirectories(t, downloadDir)
	localAlphaPath := filepath.Join(localDir, "alpha.txt")
	localAlphaContent := []byte("uploaded alpha from quick connect\n")
	utils.SetupFilesWithData(t, localAlphaContent, localAlphaPath)

	m := defaultTestModel(localDir, localDir)
	registerConnectedSessionPanel(t, m, action.Session, action.Location, 1)
	assert.Contains(t, m.fileModel.FilePanels[1].DisplayLocation(), "sf-e2e:/")

	setFocusedPanel(m, 1)
	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.BetaPath)
	m.copySingleItem(false)
	setFocusedPanel(m, 0)
	require.NoError(t, m.updateCurrentFilePanelDir(downloadDir))
	applyTeaCmd(t, m, m.getPasteItemCmd())
	assert.FileExists(t, filepath.Join(downloadDir, "beta.txt"))
	assert.Equal(t, "beta\n", string(mustReadFile(t, filepath.Join(downloadDir, "beta.txt"))))

	setFocusedPanel(m, 1)
	require.NoError(t, m.updateCurrentFilePanelDir("/"))
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), fixture.BetaPath)
	m.panelItemRename()
	m.getFocusedFilePanel().Rename.SetValue("beta-renamed.txt")
	m.confirmRename()
	assert.NoFileExists(t, remoteFixturePath(fixture, fixture.BetaPath))
	assert.FileExists(t, remoteFixturePath(fixture, "/beta-renamed.txt"))

	p := NewTestTeaProgWithEventLoop(t, m)
	setFocusedPanel(p.getModel(), 1)
	setFilePanelSelectedItemByLocation(t, p.getModel().getFocusedFilePanel(), fixture.AlphaPath)
	p.SendKey(common.Hotkeys.DeleteItems[0])
	assert.Eventually(t, p.getModel().notifyModel.IsOpen, DefaultTestTimeout, DefaultTestTick)
	assert.Equal(t, common.PermanentDeleteWarnTitle, p.getModel().notifyModel.GetTitle())
	assert.Equal(t, notify.PermanentDeleteAction, p.getModel().notifyModel.GetConfirmAction())
	p.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.Eventually(t, func() bool {
		_, err := os.Stat(remoteFixturePath(fixture, fixture.AlphaPath))
		return os.IsNotExist(err)
	}, DefaultTestTimeout, DefaultTestTick)

	m = p.getModel()
	setFocusedPanel(m, 0)
	require.NoError(t, m.updateCurrentFilePanelDir(localDir))
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), localAlphaPath)
	m.copySingleItem(false)
	setFocusedPanel(m, 1)
	require.NoError(t, m.updateCurrentFilePanelDir("/"))
	applyTeaCmd(t, m, m.getPasteItemCmd())
	assert.FileExists(t, remoteFixturePath(fixture, "/alpha.txt"))
	assert.Equal(t, string(localAlphaContent), string(mustReadFile(t, remoteFixturePath(fixture, "/alpha.txt"))))
	assert.Contains(t, stripANSI(m.getFocusedFilePanel().Render(true)), "beta-renamed.txt")
}

func TestSSHManualConnectCase(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)

	configPath := filepath.Join(t.TempDir(), "config.toml")
	knownHostsPath := filepath.Join(t.TempDir(), "known_hosts")
	require.NoError(t, os.WriteFile(knownHostsPath, nil, 0o600))
	cfg := &common.ConfigType{}
	connectModel := quickconnect.New()
	connectModel.SetKnownHostsPath(knownHostsPath)
	connectModel.SetTimeout(3 * time.Second)
	connectModel.SetManualFields(quickconnect.ManualFields{
		Name:                       "manual-localhost",
		Host:                       fixture.Host,
		Port:                       fixture.Port,
		User:                       "e2e",
		StartPath:                  "/",
		IdentityFile:               fixture.ClientKeyPath,
		IdentitiesOnly:             true,
		AuthPreference:             "publickey,password",
		Password:                   sshtest.TestPassword,
		IdentityPassphrase:         sshtest.TestKeyPassphrase,
		KeyboardInteractiveAnswers: []string{sshtest.TestKeyboardAnswer},
	})

	_, err := connectModel.SaveManualProfile(configPath, cfg)
	require.NoError(t, err)
	configBytes, err := os.ReadFile(configPath)
	require.NoError(t, err)
	configText := string(configBytes)
	assert.Contains(t, configText, "manual-localhost")
	assert.NotContains(t, configText, sshtest.TestPassword)
	assert.NotContains(t, configText, sshtest.TestKeyPassphrase)
	assert.NotContains(t, configText, sshtest.TestKeyboardAnswer)

	require.NoError(t, connectModel.Open(cfg))
	profile := quickConnectProfileByName(t, connectModel.Profiles(), "manual-localhost")
	connectModel.OpenWithProfiles([]common.SSHQuickConnectProfile{profile})
	connectModel.SetRuntimeSecrets(quickconnect.RuntimeSecrets{
		Password:                   sshtest.TestPassword,
		IdentityPassphrase:         sshtest.TestKeyPassphrase,
		KeyboardInteractiveAnswers: []string{sshtest.TestKeyboardAnswer},
	})

	action := connectModel.ConnectSelected(context.Background())
	require.Equal(t, quickconnect.ActionNone, action.Type)
	require.Equal(t, quickconnect.ModeHostKeyConfirmation, connectModel.Mode())
	action = connectModel.ConfirmHostKey(context.Background())
	require.Equal(t, quickconnect.ActionConnected, action.Type)
	require.NotNil(t, action.Session)
	t.Cleanup(func() {
		require.NoError(t, action.Session.Close())
	})
	assert.Equal(t, filesystem.SessionID("manual-localhost"), action.Location.SessionID)
	knownHostsBytes, err := os.ReadFile(knownHostsPath)
	require.NoError(t, err)
	assert.Contains(t, string(knownHostsBytes), fixture.Host)
	assert.NotContains(t, string(knownHostsBytes), "manual-localhost")

	sourcePath := filepath.Join(t.TempDir(), "alpha.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("manual overwrite attempt\n"), 0o644))
	resolver := filesystem.SessionResolverFunc(
		func(ctx context.Context, location filesystem.Location) (filesystem.Session, error) {
			if location.Provider == filesystem.ProviderLocal {
				return filesystem.NewLocalProvider().Open(ctx, location)
			}
			provider := filesystem.NewSFTPProvider(internalssh.ClientConfigRequest{
				Profile:        profile,
				KnownHostsPath: knownHostsPath,
				HostKeyAlias:   profile.HostKeyAlias,
			})
			return provider.Open(ctx, location)
		},
	)
	engine := filesystem.NewTransferEngine(resolver)

	transfer, err := engine.Start(context.Background(), filesystem.TransferRequest{
		Operation: filesystem.OperationTransferLocalToRemote,
		Source:    filesystem.Location{Provider: filesystem.ProviderLocal, Path: filesystem.NewLocalPath(sourcePath)},
		Destination: filesystem.Location{
			Provider:  filesystem.ProviderSFTP,
			SessionID: action.Location.SessionID,
			Label:     action.Location.Label,
			Path:      filesystem.NewRemotePath(fixture.AlphaPath),
		},
		Overwrite: false,
	})
	require.NoError(t, err)
	err = transfer.Wait(context.Background())
	require.Error(t, err)
	require.ErrorIs(t, err, filesystem.ErrConflict)
	assert.Equal(t, "alpha\n", string(mustReadFile(t, remoteFixturePath(fixture, fixture.AlphaPath))))
}

func TestSSHFailureModesCase(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)

	beforeBytes, err := os.ReadFile(fixture.ChangedHostKnownHostsPath)
	require.NoError(t, err)

	connectModel := quickconnect.New()
	connectModel.SetKnownHostsPath(fixture.ChangedHostKnownHostsPath)
	connectModel.SetTimeout(3 * time.Second)
	connectModel.OpenWithProfiles([]common.SSHQuickConnectProfile{sshCaseProfileForAlias(fixture, sshtest.AliasBadKey)})

	action := connectModel.ConnectSelected(context.Background())
	assert.Equal(t, quickconnect.ActionError, action.Type)
	assert.Nil(t, action.Session)
	assert.Equal(t, quickconnect.ModeBlockingWarning, connectModel.Mode())
	assert.Contains(t, connectModel.Render(), "Changed SSH host key")
	assert.Contains(t, connectModel.Render(), "connection blocked")
	afterBytes, err := os.ReadFile(fixture.ChangedHostKnownHostsPath)
	require.NoError(t, err)
	assert.Equal(t, string(beforeBytes), string(afterBytes))

	provider := filesystem.NewSFTPProvider(internalssh.ClientConfigRequest{
		Profile:        sshCaseProfileForAlias(fixture, sshtest.AliasE2E),
		KnownHostsPath: fixture.KnownHostsPath,
		HostKeyAlias:   sshtest.AliasE2E,
	})
	session, err := provider.Open(context.Background(), filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		Path:      filesystem.RootRemotePath(),
		Label:     sshtest.AliasE2E,
		SessionID: filesystem.SessionID(sshtest.AliasE2E),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, session.Close())
	})

	m := defaultTestModel(t.TempDir())
	registerConnectedSessionPanel(t, m, session, filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		Path:      filesystem.RootRemotePath(),
		Label:     sshtest.AliasE2E,
		SessionID: filesystem.SessionID(sshtest.AliasE2E),
	}, 0)
	setFocusedPanel(m, 0)
	panel := m.getFocusedFilePanel()
	setFilePanelSelectedItemByLocation(t, panel, fixture.AlphaPath)
	panel.ChangeFilePanelMode()
	panel.SetSelected(fixture.AlphaPath)
	selectedBefore := panel.GetSelectedLocationsSortedAsVisible()
	pathBefore := panel.Location

	applyTeaCmd(t, m, m.getCompressSelectedFilesCmd())
	assert.True(t, m.notifyModel.IsOpen())
	assert.Equal(
		t,
		remoteUnsupportedOperationText(filesystem.ProviderSFTP, filesystem.OperationCompress),
		m.notifyModel.GetContent(),
	)
	assert.Equal(t, pathBefore, panel.Location)
	assert.Equal(t, selectedBefore, panel.GetSelectedLocationsSortedAsVisible())
}

func quickConnectProfileByName(
	t *testing.T,
	profiles []common.SSHQuickConnectProfile,
	name string,
) common.SSHQuickConnectProfile {
	t.Helper()
	for _, profile := range profiles {
		if profile.Name == name {
			return profile
		}
	}
	t.Fatalf("profile %q not found", name)
	return common.SSHQuickConnectProfile{}
}

func sshCaseProfileForAlias(fixture *sshtest.Fixture, aliasName string) common.SSHQuickConnectProfile {
	alias := fixture.Aliases[aliasName]
	profile := common.SSHQuickConnectProfile{
		Name:           alias.Name,
		Host:           alias.Host,
		Port:           alias.Port,
		User:           alias.User,
		StartPath:      "/",
		IdentityFile:   alias.IdentityFilePath,
		IdentitiesOnly: true,
		AuthOrder:      []string{common.SSHAuthMethodPublicKey},
		HostKeyAlias:   alias.HostKeyAlias,
	}
	if alias.IdentityFilePath != "" {
		profile.IdentityFiles = []string{alias.IdentityFilePath}
	}
	return profile
}

func registerConnectedSessionPanel(
	t *testing.T,
	m *model,
	session filesystem.Session,
	location filesystem.Location,
	panelIndex int,
) {
	t.Helper()
	m.fileModel.RegisterSession(filemodel.SessionState{
		ID:          location.SessionID,
		Provider:    location.Provider,
		Label:       location.Label,
		CurrentPath: location.Path,
		Status:      filemodel.SessionConnected,
		Browser:     session,
	})
	require.NoError(t, m.fileModel.SetPaneLocation(panelIndex, location))
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	m.sessionRegistry = m.fileModel.Sessions
}

func setFocusedPanel(m *model, panelIndex int) {
	m.fileModel.FocusedPanelIndex = panelIndex
	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].IsFocused = i == panelIndex
	}
}
