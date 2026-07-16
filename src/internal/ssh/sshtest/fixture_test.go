package sshtest

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func TestSFTPFixtureLifecycle(t *testing.T) {
	fixture := Start(t)

	assert.Equal(t, "127.0.0.1", fixture.Host)
	require.NotZero(t, fixture.Port)
	assert.DirExists(t, fixture.RemoteRootPath)
	assert.FileExists(t, fixture.SSHConfigPath)
	assert.FileExists(t, fixture.KnownHostsPath)
	assert.FileExists(t, fixture.LogPath)

	configBytes, err := os.ReadFile(fixture.SSHConfigPath)
	require.NoError(t, err)
	configText := string(configBytes)
	for _, aliasName := range []string{AliasE2E, AliasBadKey, AliasPassword, AliasKey, AliasEncryptedKey, AliasKeyboard} {
		assert.Contains(t, configText, "Host "+aliasName)
	}

	assert.FileExists(t, fixture.localPath(fixture.AlphaPath))
	assert.FileExists(t, fixture.localPath(fixture.BetaPath))
	assert.FileExists(t, fixture.localPath(fixture.ReadonlyPath))
	assert.FileExists(t, fixture.localPath(fixture.SpaceNamePath))
	assert.DirExists(t, fixture.localPath(fixture.NestedPath))
	assert.DirExists(t, fixture.localPath(fixture.PermissionDeniedPath))

	client, sftpClient := newSFTPClient(t, fixture, fixture.Aliases[AliasE2E])
	t.Cleanup(func() {
		_ = sftpClient.Close()
		_ = client.Close()
	})

	entries, err := sftpClient.ReadDir("/")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(entries), 5)

	alphaBytes, err := readRemoteFile(sftpClient, fixture.AlphaPath)
	require.NoError(t, err)
	assert.Equal(t, "alpha\n", string(alphaBytes))

	betaBytes, err := readRemoteFile(sftpClient, fixture.BetaPath)
	require.NoError(t, err)
	assert.Equal(t, "beta\n", string(betaBytes))

	if fixture.SymlinkPath != "" {
		linkTarget, readLinkErr := sftpClient.ReadLink(fixture.SymlinkPath)
		require.NoError(t, readLinkErr)
		assert.Equal(t, "alpha.txt", linkTarget)
	}

	require.NoError(t, sftpClient.Close())
	require.NoError(t, client.Close())
	require.NoError(t, fixture.Close())

	_, err = net.DialTimeout("tcp", fixture.Address, 200*time.Millisecond)
	require.Error(t, err)

	logBytes, err := os.ReadFile(fixture.LogPath)
	require.NoError(t, err)
	logText := string(logBytes)
	assert.Contains(t, logText, "event=auth method=publickey")
	assert.Contains(t, logText, "auth=publickey")
	assert.Contains(t, logText, "op=list path=/")
	assert.Contains(t, logText, "op=get path=/alpha.txt")
	assert.Contains(t, logText, "conn=1")
}

func TestSFTPFixtureChangedHostLifecycle(t *testing.T) {
	fixture := Start(t)
	alias := fixture.Aliases[AliasBadKey]

	_, err := newSSHClient(alias, fixture.ClientKeyPath)
	require.Error(t, err)

	var keyErr *knownhosts.KeyError
	require.ErrorAs(t, err, &keyErr, "expected known_hosts key mismatch, got %v", err)
	assert.NotEmpty(t, keyErr.Want)

	logBytes, readErr := os.ReadFile(fixture.LogPath)
	require.NoError(t, readErr)
	logText := string(logBytes)
	assert.NotContains(t, logText, "op=")
	assert.True(t,
		strings.Contains(logText, "event=handshake result=failed") || strings.Contains(logText, "event=accept"),
		"expected connection activity in fixture log",
	)
}

func newSFTPClient(t *testing.T, fixture *Fixture, alias Alias) (*ssh.Client, *sftp.Client) {
	t.Helper()

	client, err := newSSHClient(alias, fixture.ClientKeyPath)
	require.NoError(t, err)

	sftpClient, err := sftp.NewClient(client)
	require.NoError(t, err)

	return client, sftpClient
}

func newSSHClient(alias Alias, privateKeyPath string) (*ssh.Client, error) {
	hostKeyCallback, err := knownhosts.New(alias.KnownHostsPath)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User:            alias.User,
		HostKeyCallback: hostKeyCallback,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout:         time.Second,
	}

	netConn, err := net.DialTimeout("tcp", alias.Address(), time.Second)
	if err != nil {
		return nil, err
	}

	conn, chans, reqs, err := ssh.NewClientConn(
		netConn,
		net.JoinHostPort(alias.HostKeyAlias, strconv.Itoa(alias.Port)),
		clientConfig,
	)
	if err != nil {
		_ = netConn.Close()
		return nil, err
	}

	return ssh.NewClient(conn, chans, reqs), nil
}

func readRemoteFile(client *sftp.Client, remotePath string) ([]byte, error) {
	file, err := client.Open(remotePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
