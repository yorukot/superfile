package ssh

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cryptossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
)

func TestBuildClientConfigAuthMethods(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	agentProfile := profileForAlias(fixture, sshtest.AliasKey)
	agentProfile.IdentityFile = fixture.EncryptedClientKeyPath
	agentProfile.IdentityFiles = []string{fixture.EncryptedClientKeyPath}
	passwordProfile := profileForAlias(fixture, sshtest.AliasPassword)
	passwordProfile.IdentityFile = fixture.ClientKeyPath
	passwordProfile.IdentityFiles = []string{fixture.ClientKeyPath}
	keyboardProfile := profileForAlias(fixture, sshtest.AliasKeyboard)
	keyboardProfile.IdentityFile = fixture.ClientKeyPath
	keyboardProfile.IdentityFiles = []string{fixture.ClientKeyPath}
	configuredSigner, err := publicKeySigner(fixture.EncryptedClientKeyPath, fixture.KeyPassphrase)
	require.NoError(t, err)
	agentSigner, err := publicKeySigner(fixture.ClientKeyPath, "")
	require.NoError(t, err)

	tests := []struct {
		name                string
		request             ClientConfigRequest
		wantAuth            string
		wantAuthSequence    []string
		wantKeyFingerprint  string
		rejectedFingerprint string
		wantBuildErr        string
		wantDialErr         string
	}{
		{
			name: "agent auth precedes configured identity files",
			request: ClientConfigRequest{
				Profile:                  agentProfile,
				AgentSocket:              startAgent(t, fixture.ClientKeyPath, ""),
				KnownHostsPath:           fixture.KnownHostsPath,
				ManualIdentityPassphrase: fixture.KeyPassphrase,
			},
			wantAuth:            "publickey",
			wantAuthSequence:    []string{"publickey"},
			wantKeyFingerprint:  cryptossh.FingerprintSHA256(agentSigner.PublicKey()),
			rejectedFingerprint: cryptossh.FingerprintSHA256(configuredSigner.PublicKey()),
		},
		{
			name: "configured identity file still works when agent is empty",
			request: ClientConfigRequest{
				Profile:        profileForAlias(fixture, sshtest.AliasKey),
				AgentSocket:    startEmptyAgent(t),
				KnownHostsPath: fixture.KnownHostsPath,
			},
			wantAuth: "publickey",
		},
		{
			name: "unencrypted configured identity file",
			request: ClientConfigRequest{
				Profile:        profileForAlias(fixture, sshtest.AliasKey),
				KnownHostsPath: fixture.KnownHostsPath,
			},
			wantAuth: "publickey",
		},
		{
			name: "encrypted configured identity file with passphrase",
			request: ClientConfigRequest{
				Profile:                  profileForAlias(fixture, sshtest.AliasEncryptedKey),
				KnownHostsPath:           fixture.KnownHostsPath,
				ManualIdentityPassphrase: fixture.KeyPassphrase,
			},
			wantAuth: "publickey",
		},
		{
			name: "password auth after public key methods",
			request: ClientConfigRequest{
				Profile:        passwordProfile,
				KnownHostsPath: fixture.KnownHostsPath,
				Password:       fixture.Password,
			},
			wantAuth:         "password",
			wantAuthSequence: []string{"publickey", "password"},
		},
		{
			name: "keyboard interactive auth last",
			request: ClientConfigRequest{
				Profile:        keyboardProfile,
				KnownHostsPath: fixture.KnownHostsPath,
				Password:       fixture.Password,
				KeyboardInteractive: func(_ string, _ string, questions []string, _ []bool) ([]string, error) {
					return []string{fixture.KeyboardAnswer}, nil
				},
			},
			wantAuth:         "keyboard-interactive",
			wantAuthSequence: []string{"publickey", "password", "keyboard-interactive"},
		},
		{
			name: "wrong passphrase",
			request: ClientConfigRequest{
				Profile:                  profileForAlias(fixture, sshtest.AliasEncryptedKey),
				KnownHostsPath:           fixture.KnownHostsPath,
				ManualIdentityPassphrase: "wrong-passphrase",
			},
			wantBuildErr: "parse encrypted ssh identity file",
		},
		{
			name: "wrong password",
			request: ClientConfigRequest{
				Profile:        profileForAlias(fixture, sshtest.AliasPassword),
				KnownHostsPath: fixture.KnownHostsPath,
				Password:       "wrong-password",
			},
			wantDialErr: "unable to authenticate",
		},
		{
			name: "no usable auth method",
			request: ClientConfigRequest{
				Profile:        withoutIdentity(profileForAlias(fixture, sshtest.AliasPassword)),
				KnownHostsPath: fixture.KnownHostsPath,
			},
			wantBuildErr: "no usable authentication method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeOffset := fixtureLogSize(t, fixture.LogPath)
			bundle, err := BuildClientConfig(tt.request)
			if tt.wantBuildErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantBuildErr)
				assertNoSecretLeak(t, err.Error())
				return
			}
			require.NoError(t, err)
			defer bundle.Close()

			client, err := bundle.Dial()
			if tt.wantDialErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantDialErr)
				assertNoSecretLeak(t, err.Error())
				return
			}
			require.NoError(t, err)
			require.NoError(t, client.Close())

			logText := fixtureLogSince(t, fixture.LogPath, beforeOffset)
			assert.Contains(t, logText, "auth="+tt.wantAuth)
			if len(tt.wantAuthSequence) > 0 {
				assert.Equal(t, tt.wantAuthSequence, authAttemptSequence(logText))
			}
			if tt.wantKeyFingerprint != "" {
				assert.Contains(t, logText, "fingerprint="+tt.wantKeyFingerprint)
			}
			if tt.rejectedFingerprint != "" {
				assert.NotContains(t, logText, "fingerprint="+tt.rejectedFingerprint)
			}
			assertNoSecretLeak(t, logText)
		})
	}
}

func authAttemptSequence(logText string) []string {
	sequence := make([]string, 0)
	for _, line := range strings.Split(logText, "\n") {
		const marker = "event=auth method="
		markerIndex := strings.Index(line, marker)
		if markerIndex < 0 || strings.Contains(line, " result=") {
			continue
		}
		method := strings.Fields(line[markerIndex+len(marker):])[0]
		if len(sequence) == 0 || sequence[len(sequence)-1] != method {
			sequence = append(sequence, method)
		}
	}
	return sequence
}

func TestBuildClientConfigHonorsProfileAuthOrder(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)

	t.Run("password before keyboard interactive skips later keyboard prompt after success", func(t *testing.T) {
		profile := withoutIdentity(profileForAlias(fixture, sshtest.AliasPassword))
		profile.AuthOrder = []string{common.SSHAuthMethodPassword, common.SSHAuthMethodKeyboardInteractive}
		beforeOffset := fixtureLogSize(t, fixture.LogPath)

		bundle, err := BuildClientConfig(ClientConfigRequest{
			Profile:        profile,
			KnownHostsPath: fixture.KnownHostsPath,
			Password:       fixture.Password,
			KeyboardInteractive: func(_ string, _ string, _ []string, _ []bool) ([]string, error) {
				return []string{fixture.KeyboardAnswer}, nil
			},
		})
		require.NoError(t, err)
		defer bundle.Close()

		client, err := bundle.Dial()
		require.NoError(t, err)
		require.NoError(t, client.Close())

		logText := fixtureLogSince(t, fixture.LogPath, beforeOffset)
		assert.Contains(t, logText, "event=auth method=password")
		assert.NotContains(t, logText, "event=auth method=keyboard-interactive")
		assert.Contains(t, logText, "auth=password")
		assertNoSecretLeak(t, logText)
	})

	t.Run("excluding password omits supplied password auth", func(t *testing.T) {
		profile := withoutIdentity(profileForAlias(fixture, sshtest.AliasPassword))
		profile.AuthOrder = []string{common.SSHAuthMethodKeyboardInteractive}
		beforeOffset := fixtureLogSize(t, fixture.LogPath)

		bundle, err := BuildClientConfig(ClientConfigRequest{
			Profile:        profile,
			KnownHostsPath: fixture.KnownHostsPath,
			Password:       fixture.Password,
			KeyboardInteractive: func(_ string, _ string, _ []string, _ []bool) ([]string, error) {
				return []string{fixture.KeyboardAnswer}, nil
			},
		})
		require.NoError(t, err)
		defer bundle.Close()

		_, err = bundle.Dial()
		require.Error(t, err)
		logText := fixtureLogSince(t, fixture.LogPath, beforeOffset)
		assert.Contains(t, logText, "event=auth method=keyboard-interactive")
		assert.NotContains(t, logText, "event=auth method=password")
		assert.NotContains(t, logText, "event=auth method=publickey")
		assertNoSecretLeak(t, logText)
	})

	t.Run("keyboard interactive only excludes password and publickey", func(t *testing.T) {
		profile := profileForAlias(fixture, sshtest.AliasKeyboard)
		profile.IdentityFile = fixture.ClientKeyPath
		profile.IdentityFiles = []string{fixture.ClientKeyPath}
		profile.AuthOrder = []string{common.SSHAuthMethodKeyboardInteractive}
		beforeOffset := fixtureLogSize(t, fixture.LogPath)

		bundle, err := BuildClientConfig(ClientConfigRequest{
			Profile:        profile,
			KnownHostsPath: fixture.KnownHostsPath,
			Password:       fixture.Password,
			KeyboardInteractive: func(_ string, _ string, _ []string, _ []bool) ([]string, error) {
				return []string{fixture.KeyboardAnswer}, nil
			},
		})
		require.NoError(t, err)
		defer bundle.Close()

		client, err := bundle.Dial()
		require.NoError(t, err)
		require.NoError(t, client.Close())

		logText := fixtureLogSince(t, fixture.LogPath, beforeOffset)
		assert.Contains(t, logText, "event=auth method=keyboard-interactive")
		assert.NotContains(t, logText, "event=auth method=password")
		assert.NotContains(t, logText, "event=auth method=publickey")
		assert.Contains(t, logText, "auth=keyboard-interactive")
		assertNoSecretLeak(t, logText)
	})
}

func TestHostKeyUnknownRejectAcceptAndChangedReject(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)

	t.Run("unknown host reject returns typed confirmation request without write", func(t *testing.T) {
		knownHostsPath := filepath.Join(t.TempDir(), "known_hosts")
		require.NoError(t, os.WriteFile(knownHostsPath, nil, 0o600))
		bundle := buildHostKeyTestBundle(t, fixture, sshtest.AliasKey, knownHostsPath)
		defer bundle.Close()

		client, err := bundle.Dial()
		require.Error(t, err)
		assert.Nil(t, client)
		var unknownHost *UnknownHostKeyError
		require.ErrorAs(t, err, &unknownHost, "got %T %v", err, err)
		assert.Equal(t, sshtest.AliasKey, unknownHost.Host)
		assert.NotEmpty(t, unknownHost.Address)
		assert.NotEmpty(t, unknownHost.KeyType)
		assert.True(t, strings.HasPrefix(unknownHost.Fingerprint, "SHA256:"))
		assert.Equal(t, knownHostsPath, unknownHost.KnownHostsPath)
		knownHostsBytes, readErr := os.ReadFile(knownHostsPath)
		require.NoError(t, readErr)
		assert.Empty(t, knownHostsBytes)
	})

	t.Run("unknown host accept persists to injected known hosts", func(t *testing.T) {
		knownHostsPath := filepath.Join(t.TempDir(), "known_hosts")
		require.NoError(t, os.WriteFile(knownHostsPath, nil, 0o600))
		bundle := buildHostKeyTestBundle(t, fixture, sshtest.AliasKey, knownHostsPath)
		defer bundle.Close()

		_, err := bundle.Dial()
		require.Error(t, err)
		require.NoError(t, AcceptUnknownHostKey(err))
		knownHostsBytes, readErr := os.ReadFile(knownHostsPath)
		require.NoError(t, readErr)
		assert.Contains(t, string(knownHostsBytes), sshtest.AliasKey)

		acceptedBundle := buildHostKeyTestBundle(t, fixture, sshtest.AliasKey, knownHostsPath)
		defer acceptedBundle.Close()
		client, dialErr := acceptedBundle.Dial()
		require.NoError(t, dialErr)
		require.NoError(t, client.Close())
	})

	t.Run("changed host key rejects without write", func(t *testing.T) {
		beforeBytes, err := os.ReadFile(fixture.ChangedHostKnownHostsPath)
		require.NoError(t, err)
		bundle := buildHostKeyTestBundle(t, fixture, sshtest.AliasBadKey, fixture.ChangedHostKnownHostsPath)
		defer bundle.Close()

		client, dialErr := bundle.Dial()
		require.Error(t, dialErr)
		assert.Nil(t, client)
		var unknownHost *UnknownHostKeyError
		assert.NotErrorAs(t, dialErr, &unknownHost, "changed host must not become an accept prompt")
		afterBytes, readErr := os.ReadFile(fixture.ChangedHostKnownHostsPath)
		require.NoError(t, readErr)
		assert.Equal(t, string(beforeBytes), string(afterBytes))
	})
}

func TestInjectedKnownHostsPathDoesNotUseDefaultKnownHosts(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	fakeHome := t.TempDir()
	defaultKnownHostsPath := filepath.Join(fakeHome, ".ssh", "known_hosts")
	require.NoError(t, os.MkdirAll(filepath.Dir(defaultKnownHostsPath), 0o700))
	defaultContents := []byte("poison default known hosts must stay untouched\n")
	require.NoError(t, os.WriteFile(defaultKnownHostsPath, defaultContents, 0o600))
	t.Setenv("HOME", fakeHome)
	t.Setenv("USERPROFILE", fakeHome)

	injectedKnownHostsPath := filepath.Join(t.TempDir(), "injected_known_hosts")
	copyFile(t, fixture.KnownHostsPath, injectedKnownHostsPath)

	bundle := buildHostKeyTestBundle(t, fixture, sshtest.AliasE2E, injectedKnownHostsPath)
	defer bundle.Close()
	client, err := bundle.Dial()
	require.NoError(t, err)
	require.NoError(t, client.Close())

	afterDefault, err := os.ReadFile(defaultKnownHostsPath)
	require.NoError(t, err)
	assert.Equal(t, string(defaultContents), string(afterDefault))
}

func TestAcceptUnknownHostKeyUsesDefaultKnownHostsWhenNotInjected(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)
	t.Setenv("USERPROFILE", fakeHome)
	defaultKnownHostsPath := filepath.Join(fakeHome, ".ssh", "known_hosts")
	require.NoError(t, os.MkdirAll(filepath.Dir(defaultKnownHostsPath), 0o700))
	require.NoError(t, os.WriteFile(defaultKnownHostsPath, nil, 0o600))

	request := ClientConfigRequest{
		Profile:      profileForAlias(fixture, sshtest.AliasKey),
		HostKeyAlias: sshtest.AliasKey,
	}
	bundle, err := BuildClientConfig(request)
	require.NoError(t, err)
	defer bundle.Close()

	_, err = bundle.Dial()
	require.Error(t, err)
	require.NoError(t, AcceptUnknownHostKey(err))
	knownHostsBytes, readErr := os.ReadFile(defaultKnownHostsPath)
	require.NoError(t, readErr)
	assert.Contains(t, string(knownHostsBytes), sshtest.AliasKey)
}

func TestRedactionHelpersRemoveSecretsAndPrivateKeys(t *testing.T) {
	message := "secret-password secret-passphrase -----BEGIN OPENSSH PRIVATE KEY-----"
	redactedMessage := RedactString(message)
	assert.NotContains(t, redactedMessage, "secret-password")
	assert.NotContains(t, redactedMessage, "secret-passphrase")
	assert.NotContains(t, redactedMessage, "-----BEGIN OPENSSH PRIVATE KEY-----")
}

func profileForAlias(fixture *sshtest.Fixture, aliasName string) common.SSHQuickConnectProfile {
	alias := fixture.Aliases[aliasName]
	profile := common.SSHQuickConnectProfile{
		Name:          alias.Name,
		Host:          alias.Host,
		Port:          alias.Port,
		User:          alias.User,
		IdentityFile:  alias.IdentityFilePath,
		IdentityFiles: nil,
	}
	if alias.IdentityFilePath != "" {
		profile.IdentityFiles = []string{alias.IdentityFilePath}
	}
	return profile
}

func withoutIdentity(profile common.SSHQuickConnectProfile) common.SSHQuickConnectProfile {
	profile.IdentityFile = ""
	profile.IdentityFiles = nil
	return profile
}

func buildHostKeyTestBundle(
	t *testing.T,
	fixture *sshtest.Fixture,
	aliasName string,
	knownHostsPath string,
) *ClientConfigBundle {
	t.Helper()
	req := ClientConfigRequest{
		Profile:        profileForAlias(fixture, aliasName),
		HostKeyAlias:   aliasName,
		KnownHostsPath: knownHostsPath,
	}
	bundle, err := BuildClientConfig(req)
	require.NoError(t, err)
	return bundle
}

func startAgent(t *testing.T, identityPath string, passphrase string) string {
	t.Helper()
	keyBytes, err := os.ReadFile(identityPath)
	require.NoError(t, err)
	var privateKey any
	if passphrase == "" {
		privateKey, err = cryptossh.ParseRawPrivateKey(keyBytes)
	} else {
		privateKey, err = cryptossh.ParseRawPrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
	}
	require.NoError(t, err)

	keyring := agent.NewKeyring()
	require.NoError(t, keyring.Add(agent.AddedKey{PrivateKey: privateKey}))
	//nolint:usetesting // Short /tmp path avoids Unix socket length limits.
	socketDir, err := os.MkdirTemp(
		"/tmp",
		"sf-agent-",
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(socketDir) })
	socketPath := filepath.Join(socketDir, "a.sock")
	listener, err := net.Listen("unix", socketPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go func() { _ = agent.ServeAgent(keyring, conn) }()
		}
	}()

	return socketPath
}

func startEmptyAgent(t *testing.T) string {
	t.Helper()
	keyring := agent.NewKeyring()
	//nolint:usetesting // Short /tmp path avoids Unix socket length limits.
	socketDir, err := os.MkdirTemp(
		"/tmp",
		"sf-agent-empty-",
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(socketDir) })
	socketPath := filepath.Join(socketDir, "a.sock")
	listener, err := net.Listen("unix", socketPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go func() { _ = agent.ServeAgent(keyring, conn) }()
		}
	}()

	return socketPath
}

func startTrackedEmptyAgent(t *testing.T) (string, <-chan struct{}) {
	t.Helper()
	keyring := agent.NewKeyring()
	//nolint:usetesting // Short /tmp path avoids Unix socket length limits.
	socketDir, err := os.MkdirTemp("/tmp", "sf-agent-tracked-empty-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(socketDir) })
	socketPath := filepath.Join(socketDir, "a.sock")
	listener, err := net.Listen("unix", socketPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	closed := make(chan struct{})
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			close(closed)
			return
		}
		_ = agent.ServeAgent(keyring, conn)
		_ = conn.Close()
		close(closed)
	}()
	return socketPath, closed
}

func fixtureLogSize(t *testing.T, path string) int64 {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err)
	return info.Size()
}

func fixtureLogSince(t *testing.T, path string, offset int64) string {
	t.Helper()
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.Seek(offset, 0)
	require.NoError(t, err)
	bytes, err := os.ReadFile(path)
	require.NoError(t, err)
	if offset >= int64(len(bytes)) {
		return ""
	}
	return string(bytes[offset:])
}

func copyFile(t *testing.T, source string, destination string) {
	t.Helper()
	bytes, err := os.ReadFile(source)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(destination, bytes, 0o600))
}

func assertNoSecretLeak(t *testing.T, text string) {
	t.Helper()
	assert.NotContains(t, text, sshtest.TestPassword)
	assert.NotContains(t, text, sshtest.TestKeyPassphrase)
	assert.NotContains(t, text, "PRIVATE KEY")
}

func TestClientConfigBundleDialRejectsNilBundle(t *testing.T) {
	var bundle *ClientConfigBundle
	_, err := bundle.Dial()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestStrictHostKeyCallbackCreatesMissingKnownHostsSecurely(t *testing.T) {
	knownHostsPath := filepath.Join(t.TempDir(), ".ssh", "known_hosts")
	_, err := StrictHostKeyCallback(knownHostsPath)
	require.NoError(t, err)

	info, err := os.Stat(knownHostsPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	directoryInfo, err := os.Stat(filepath.Dir(knownHostsPath))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), directoryInfo.Mode().Perm())
}

func TestBuildClientConfigSkipsMissingIdentityAndStaleAgent(t *testing.T) {
	knownHostsPath := filepath.Join(t.TempDir(), "known_hosts")
	profile := common.SSHQuickConnectProfile{
		Host:          "example.com",
		User:          "user",
		IdentityFiles: []string{filepath.Join(t.TempDir(), "missing-key")},
		AuthOrder:     []string{common.SSHAuthMethodPublicKey, common.SSHAuthMethodPassword},
	}
	bundle, err := BuildClientConfig(ClientConfigRequest{
		Profile:        profile,
		KnownHostsPath: knownHostsPath,
		AgentSocket:    filepath.Join(t.TempDir(), "stale-agent.sock"),
		Password:       "runtime-password",
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = bundle.Close() })
	assert.Len(t, bundle.Config.Auth, 1)
}

func TestBuildClientConfigClosesAgentConnectionOnBuildFailure(t *testing.T) {
	tests := []struct {
		name         string
		identityFile func(*testing.T) string
		errorText    string
	}{
		{
			name:      "no usable authentication methods",
			errorText: "no usable authentication method",
		},
		{
			name: "identity parsing failure",
			identityFile: func(t *testing.T) string {
				path := filepath.Join(t.TempDir(), "invalid-key")
				require.NoError(t, os.WriteFile(path, []byte("not a private key"), 0o600))
				return path
			},
			errorText: "parse ssh identity file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socketPath, agentClosed := startTrackedEmptyAgent(t)
			profile := common.SSHQuickConnectProfile{
				Host:      "example.com",
				User:      "user",
				AuthOrder: []string{common.SSHAuthMethodPublicKey},
			}
			if tt.identityFile != nil {
				profile.IdentityFiles = []string{tt.identityFile(t)}
			}
			_, err := BuildClientConfig(ClientConfigRequest{
				Profile:        profile,
				KnownHostsPath: filepath.Join(t.TempDir(), "known_hosts"),
				AgentSocket:    socketPath,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorText)

			select {
			case <-agentClosed:
			case <-time.After(time.Second):
				t.Fatal("SSH agent connection remained open after client-config failure")
			}
		})
	}
}

func TestBuildClientConfigRejectsMissingRequiredProfileFields(t *testing.T) {
	_, err := BuildClientConfig(ClientConfigRequest{Profile: common.SSHQuickConnectProfile{User: "u"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")

	_, err = BuildClientConfig(ClientConfigRequest{Profile: common.SSHQuickConnectProfile{Host: "h"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user is required")
}

func TestWrongPasswordErrorIsRedacted(t *testing.T) {
	err := RedactError(errors.New("auth failed for secret-password"))
	require.Error(t, err)
	assert.NotContains(t, err.Error(), sshtest.TestPassword)
}

func TestBuildClientConfigTimeoutDefault(t *testing.T) {
	fixture := sshtest.Start(t)
	bundle, err := BuildClientConfig(ClientConfigRequest{
		Profile:        profileForAlias(fixture, sshtest.AliasKey),
		KnownHostsPath: fixture.KnownHostsPath,
	})
	require.NoError(t, err)
	defer bundle.Close()
	assert.Equal(t, defaultSSHTimeout, bundle.Config.Timeout)

	customTimeout := 250 * time.Millisecond
	bundle, err = BuildClientConfig(ClientConfigRequest{
		Profile:        profileForAlias(fixture, sshtest.AliasKey),
		KnownHostsPath: fixture.KnownHostsPath,
		Timeout:        customTimeout,
	})
	require.NoError(t, err)
	defer bundle.Close()
	assert.Equal(t, customTimeout, bundle.Config.Timeout)
}

func TestIdentityFileTildePathExpandsToUserHome(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	sshDir := filepath.Join(homeDir, ".ssh")
	require.NoError(t, os.MkdirAll(sshDir, 0o700))
	identityPath := filepath.Join(sshDir, "id_ed25519")
	keyBytes, err := os.ReadFile(fixture.ClientKeyPath)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(identityPath, keyBytes, 0o600))

	profile := common.SSHQuickConnectProfile{
		Name:           sshtest.AliasKey,
		Host:           fixture.Host,
		Port:           fixture.Port,
		User:           "key",
		IdentityFile:   "~/.ssh/id_ed25519",
		IdentityFiles:  []string{"~/.ssh/id_ed25519"},
		IdentitiesOnly: true,
		AuthOrder:      []string{common.SSHAuthMethodPublicKey},
	}

	bundle, err := BuildClientConfig(ClientConfigRequest{
		Profile:        profile,
		KnownHostsPath: fixture.KnownHostsPath,
		HostKeyAlias:   sshtest.AliasKey,
	})
	require.NoError(t, err)
	defer bundle.Close()

	client, err := bundle.Dial()
	require.NoError(t, err)
	require.NoError(t, client.Close())

	assert.Equal(t, []string{identityPath}, identityFiles(profile))
}
