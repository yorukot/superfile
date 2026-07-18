package quickconnect

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
	"github.com/yorukot/superfile/src/pkg/utils"
)

var quickConnectTestConfigOnce sync.Once //nolint:gochecknoglobals // Package test setup must run once.

func TestSetDimensionsClampsSmallPositiveWidth(t *testing.T) {
	model := New()
	model.SetDimensions(1, 0)
	assert.Equal(t, common.InnerPadding, model.width)

	model.SetDimensions(0, 0)
	assert.Equal(t, common.InnerPadding, model.width)
}

func TestCancelingConnectionFlowsClearsRuntimeState(t *testing.T) {
	populateQuickConnectTestConfig(t)
	tests := []struct {
		name   string
		mode   Mode
		cancel func(*Model)
	}{
		{
			name: "manual",
			mode: ModeManual,
			cancel: func(model *Model) {
				model.handleManualUpdateKey(tea.KeyPressMsg{Code: tea.KeyEscape})
			},
		},
		{
			name: "credentials",
			mode: ModeCredentials,
			cancel: func(model *Model) {
				model.handleCredentialKey(tea.KeyPressMsg{Code: tea.KeyEscape})
			},
		},
		{
			name: "host key",
			mode: ModeHostKeyConfirmation,
			cancel: func(model *Model) {
				model.handleHostKeyKey(tea.KeyPressMsg{Code: tea.KeyEscape}.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.mode = tt.mode
			model.pendingProfile = common.SSHQuickConnectProfile{Name: "old-profile"}
			model.secrets = RuntimeSecrets{Password: "old-password", IdentityPassphrase: "old-passphrase"}
			model.lastConnectErr = assert.AnError

			tt.cancel(&model)

			assert.Equal(t, ModeList, model.mode)
			assert.Empty(t, model.pendingProfile)
			assert.Empty(t, model.secrets)
			assert.NoError(t, model.lastConnectErr)
		})
	}
}

func TestOpenListsDiscoveredAliasAndSavedManualProfile(t *testing.T) {
	populateQuickConnectTestConfig(t)
	fixture := sshtest.Start(t)
	model := New()
	model.SetDiscoveryOptions(common.SSHConfigDiscoveryOptions{
		UserConfigPath:   fixture.SSHConfigPath,
		SystemConfigPath: filepath.Join(t.TempDir(), "missing_system_config"),
	})

	cfg := &common.ConfigType{SSH: common.SSHConfigSection{Profiles: []common.SSHProfileType{{
		Name:      "manual-localhost",
		Host:      "127.0.0.1",
		Port:      fixture.Port,
		User:      "manual-user",
		StartPath: "/tmp/manual-remote",
		AuthOrder: []string{common.SSHAuthMethodPassword},
	}}}}

	require.NoError(t, model.Open(cfg))
	names := profileNames(model.Profiles())
	assert.Contains(t, names, sshtest.AliasE2E)
	assert.Contains(t, names, "manual-localhost")
	assert.Contains(t, model.Render(), sshtest.AliasE2E)
	assert.Contains(t, model.Render(), "manual-localhost")
}

func TestManualProfileSaveWritesOnlyNonSecretFields(t *testing.T) {
	populateQuickConnectTestConfig(t)
	configPath := filepath.Join(t.TempDir(), "config.toml")
	model := New()
	model.SetManualFields(ManualFields{
		Name:                       "manual-localhost",
		Host:                       "127.0.0.1",
		Port:                       2200,
		User:                       "sfuser",
		StartPath:                  "/tmp/sf-remote",
		IdentityFile:               "~/.ssh/id_ed25519",
		IdentitiesOnly:             true,
		AuthPreference:             "publickey,password",
		Password:                   sshtest.TestPassword,
		IdentityPassphrase:         sshtest.TestKeyPassphrase,
		KeyboardInteractiveAnswers: []string{sshtest.TestKeyboardAnswer},
	})

	cfg := &common.ConfigType{}
	profile, err := model.SaveManualProfile(configPath, cfg)
	require.NoError(t, err)
	assert.Equal(t, "manual-localhost", profile.Name)
	assert.Equal(t, []string{common.SSHAuthMethodPublicKey, common.SSHAuthMethodPassword}, profile.AuthOrder)

	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	text := string(data)
	assert.Contains(t, text, "manual-localhost")
	assert.Contains(t, text, "identity_file")
	assert.NotContains(t, text, sshtest.TestPassword)
	assert.NotContains(t, text, sshtest.TestKeyPassphrase)
	assert.NotContains(t, text, sshtest.TestKeyboardAnswer)
	assert.NotContains(t, text, "password =")
	assert.NotContains(t, text, "passphrase")

	var raw map[string]any
	require.NoError(t, toml.Unmarshal(data, &raw))
	sshSection := raw["ssh"].(map[string]any)
	rawProfiles := sshSection["profile"].([]any)
	rawProfile := rawProfiles[0].(map[string]any)
	assert.Equal(t, "127.0.0.1", rawProfile["host"])
	assert.NotContains(t, rawProfile, "password")
	assert.NotContains(t, rawProfile, "passphrase")
	assert.NotContains(t, rawProfile, "keyboard_interactive_answers")

	model.mode = ModeManual
	rendered := model.Render()
	for _, label := range []string{"Host", "Port", "User", "Start path", "Identity file", "Auth preference"} {
		assert.Contains(t, rendered, label)
	}
}

func TestManualRenderActiveNameDoesNotPadBeforeValue(t *testing.T) {
	populateQuickConnectTestConfig(t)
	model := New()
	model.mode = ModeManual
	model.SetManualFields(ManualFields{
		Name:           "prod",
		Host:           "example.com",
		Port:           22,
		User:           "sfuser",
		StartPath:      "/",
		AuthPreference: defaultAuthPreference(),
	})

	rendered := ansi.Strip(model.Render())
	nameLine := lineContaining(t, rendered, "Name:")

	assert.Contains(t, nameLine, "Name: prod")
	assert.NotContains(t, nameLine, "Name:           prod")
	assert.Contains(t, rendered, "Host:           example.com")
	assert.Contains(t, rendered, "Auth preference:")
}

func TestManualFieldActiveRowRestylesTextAfterCursor(t *testing.T) {
	originalModalStyle := common.ModalStyle
	originalModalCursorStyle := common.ModalCursorStyle
	//nolint:reassign // Test verifies package-level theme styling.
	common.ModalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.White).
		Background(lipgloss.Black)
	//nolint:reassign // Test verifies package-level theme styling.
	common.ModalCursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Black).
		Background(lipgloss.White)
	t.Cleanup(func() {
		common.ModalStyle = originalModalStyle             //nolint:reassign // Restore shared test theme state.
		common.ModalCursorStyle = originalModalCursorStyle //nolint:reassign // Restore shared test theme state.
	})

	cursor := common.ModalCursorStyle.Render(icon.Cursor)

	row := formatManualFieldRow(cursor, "Name", "prod", true)
	styledText := common.ModalStyle.Render(" Name: prod")

	require.True(t, strings.HasPrefix(row, cursor))
	assert.Equal(t, styledText, strings.TrimPrefix(row, cursor))
	assert.Contains(t, styledText, "\x1b[")
	assert.NotEqual(t, cursor+" Name: prod", row)
	assert.Contains(t, ansi.Strip(row), "Name: prod")
	assert.Equal(t, "  Host:           example.com", formatManualFieldRow(" ", "Host", "example.com", false))
}

func TestManualEntryTreatsPrintableVimKeysAsText(t *testing.T) {
	populateQuickConnectTestConfig(t)
	model := New()
	model.mode = ModeManual

	for _, key := range []string{"h", "j", "k", "l"} {
		action := model.handleManualKey(utils.TeaRuneKeyMsg(key))
		assert.Equal(t, ActionNone, action.Type)
	}

	assert.Equal(t, "hjkl", model.manual.Name)
	assert.Equal(t, manualFieldName, model.manualCursor)
	assert.Equal(t, ModeManual, model.Mode())
}

func TestManualEntryKeepsNonPrintableControls(t *testing.T) {
	populateQuickConnectTestConfig(t)
	model := New()
	model.mode = ModeManual

	model.handleManualKey(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(t, manualFieldHost, model.manualCursor)

	model.handleManualKey(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, manualFieldName, model.manualCursor)

	model.handleManualKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	assert.Equal(t, ModeList, model.Mode())
}

func TestUnknownHostConfirmationShowsFingerprintAndWritesOnlyAfterAccept(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	populateQuickConnectTestConfig(t)
	fixture := sshtest.Start(t)
	knownHostsPath := filepath.Join(t.TempDir(), "known_hosts")
	require.NoError(t, os.WriteFile(knownHostsPath, nil, 0o600))

	model := New()
	model.SetKnownHostsPath(knownHostsPath)
	model.SetTimeout(3 * time.Second)
	model.OpenWithProfiles([]common.SSHQuickConnectProfile{profileForAlias(fixture, sshtest.AliasKey)})

	action := model.ConnectSelected(context.Background())
	assert.Equal(t, ActionNone, action.Type)
	assert.Equal(t, ModeHostKeyConfirmation, model.Mode())
	prompt := model.Render()
	assert.Contains(t, prompt, "Key type:")
	assert.Contains(t, prompt, "Fingerprint: SHA256:")
	assert.Contains(t, prompt, "known_hosts")
	beforeAccept, err := os.ReadFile(knownHostsPath)
	require.NoError(t, err)
	assert.Empty(t, beforeAccept)

	action = model.ConfirmHostKey(context.Background())
	if action.Error != nil {
		t.Logf("confirm host key reconnect error: %v", action.Error)
	}
	require.Equal(t, ActionConnected, action.Type)
	require.NotNil(t, action.Session)
	t.Cleanup(func() { _ = action.Session.Close() })
	afterAccept, err := os.ReadFile(knownHostsPath)
	require.NoError(t, err)
	assert.Contains(t, string(afterAccept), sshtest.AliasKey)
}

func TestUnnamedProfilesReceiveUniqueSessionIDs(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	populateQuickConnectTestConfig(t)
	fixture := sshtest.Start(t)
	profile := profileForAlias(fixture, sshtest.AliasKey)
	profile.Name = ""

	model := New()
	model.SetKnownHostsPath(fixture.KnownHostsPath)
	model.SetTimeout(3 * time.Second)
	model.OpenWithProfiles([]common.SSHQuickConnectProfile{profile})
	first := model.ConnectSelected(context.Background())
	require.Equal(t, ActionConnected, first.Type)
	require.NotNil(t, first.Session)
	require.NoError(t, first.Session.Close())
	require.NotEmpty(t, first.Location.SessionID)

	model.OpenWithProfiles([]common.SSHQuickConnectProfile{profile})
	second := model.ConnectSelected(context.Background())
	require.Equal(t, ActionConnected, second.Type)
	require.NotNil(t, second.Session)
	require.NoError(t, second.Session.Close())
	require.NotEmpty(t, second.Location.SessionID)
	assert.NotEqual(t, first.Location.SessionID, second.Location.SessionID)
}

func TestChangedHostKeyShowsBlockingWarningAndDoesNotOpenSession(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	populateQuickConnectTestConfig(t)
	fixture := sshtest.Start(t)
	beforeBytes, err := os.ReadFile(fixture.ChangedHostKnownHostsPath)
	require.NoError(t, err)

	model := New()
	model.SetKnownHostsPath(fixture.ChangedHostKnownHostsPath)
	model.SetTimeout(3 * time.Second)
	model.OpenWithProfiles([]common.SSHQuickConnectProfile{profileForAlias(fixture, sshtest.AliasBadKey)})

	action := model.ConnectSelected(context.Background())
	assert.Equal(t, ActionError, action.Type)
	assert.Nil(t, action.Session)
	assert.Equal(t, ModeBlockingWarning, model.Mode())
	assert.Contains(t, model.Render(), "Changed SSH host key")
	assert.Contains(t, model.Render(), "connection blocked")
	afterBytes, err := os.ReadFile(fixture.ChangedHostKnownHostsPath)
	require.NoError(t, err)
	assert.Equal(t, string(beforeBytes), string(afterBytes))
}

func TestConnectSelectedRejectsInvalidStartPathBeforeSuccess(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	populateQuickConnectTestConfig(t)
	fixture := sshtest.Start(t)

	model := New()
	model.SetKnownHostsPath(fixture.KnownHostsPath)
	model.SetTimeout(3 * time.Second)
	profile := profileForAlias(fixture, sshtest.AliasE2E)
	profile.StartPath = "/missing-start-path"
	model.OpenWithProfiles([]common.SSHQuickConnectProfile{profile})

	action := model.ConnectSelected(context.Background())
	assert.Equal(t, ActionError, action.Type)
	assert.Nil(t, action.Session)
	assert.Equal(t, ModeBlockingWarning, model.Mode())
	assert.Contains(t, model.Warning(), "missing-start-path")
}

func populateQuickConnectTestConfig(t *testing.T) {
	t.Helper()
	quickConnectTestConfigOnce.Do(func() {
		require.NoError(t, common.PopulateGlobalConfigs())
	})
}

func profileForAlias(fixture *sshtest.Fixture, aliasName string) common.SSHQuickConnectProfile {
	alias := fixture.Aliases[aliasName]
	identityFiles := []string(nil)
	if alias.IdentityFilePath != "" {
		identityFiles = []string{alias.IdentityFilePath}
	}
	return common.SSHQuickConnectProfile{
		Name:          alias.Name,
		Host:          alias.Host,
		Port:          alias.Port,
		User:          alias.User,
		StartPath:     "/",
		IdentityFile:  alias.IdentityFilePath,
		IdentityFiles: identityFiles,
		AuthOrder:     strings.Split(alias.PreferredAuthentications, ","),
		HostKeyAlias:  alias.HostKeyAlias,
	}
}

func profileNames(profiles []common.SSHQuickConnectProfile) []string {
	names := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	return names
}

func lineContaining(t *testing.T, text string, needle string) string {
	t.Helper()
	for line := range strings.SplitSeq(text, "\n") {
		if strings.Contains(line, needle) {
			return line
		}
	}
	require.Failf(t, "missing line", "expected line containing %q", needle)
	return ""
}

func TestSaveManualProfileRejectsMissingHost(t *testing.T) {
	_, err := SaveManualProfile("", &common.ConfigType{}, ManualFields{Name: "manual-localhost"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")
}

func TestManualRuntimeSecretsStayOutsideSavedProfile(t *testing.T) {
	profile, err := SaveManualProfile("", &common.ConfigType{}, ManualFields{
		Name:                       "manual-localhost",
		Host:                       "localhost",
		Port:                       22,
		User:                       "sfuser",
		Password:                   "runtime-only-secret",
		IdentityPassphrase:         sshtest.TestKeyPassphrase,
		KeyboardInteractiveAnswers: []string{sshtest.TestKeyboardAnswer},
	})
	require.NoError(t, err)
	assert.Equal(t, "manual-localhost", profile.Name)
	assert.Empty(t, profile.IdentityFile)
}
