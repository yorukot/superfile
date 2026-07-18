package filesystem

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	internalssh "github.com/yorukot/superfile/src/internal/ssh"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	_ = os.Setenv("SSH_AUTH_SOCK", "")
	if err := common.PopulateGlobalConfigs(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestTransferEngineLocalToRemoteUploadChecksumsAndProcessbar(t *testing.T) {
	fixture := sshtest.Start(t)
	resolver := newTransferTestResolver(fixture)
	engine := NewTransferEngine(resolver)
	bar := newTransferProcessBar(t)

	for _, name := range []string{"alpha.txt", "beta.txt"} {
		sourcePath := filepath.Join(fixture.RemoteRootPath, name)
		destinationPath := "/uploaded-" + name

		transfer, err := engine.Start(context.Background(), TransferRequest{
			Operation:   OperationTransferLocalToRemote,
			Source:      localLocation(sourcePath),
			Destination: remoteLocation(sshtest.AliasE2E, destinationPath),
			Overwrite:   true,
		})
		require.NoError(t, err)

		process, err := TrackTransferProcess(context.Background(), bar, transfer)
		require.NoError(t, err)
		require.NoError(t, transfer.Wait(context.Background()))

		assert.Equal(t, sha256File(t, sourcePath), sha256RemoteFile(t, resolver, destinationPath))
		assert.Equal(t, mustReadFile(t, sourcePath), mustReadRemoteFile(t, resolver, destinationPath))
		waitForProcessState(t, bar, process.ID, processbar.Successful)
	}
}

func TestValidateTransferTopologyRejectsSameAndNestedPaths(t *testing.T) {
	tests := []struct {
		name        string
		source      Location
		destination Location
	}{
		{
			name:        "same local path",
			source:      localLocation(filepath.Join("tmp", "same")),
			destination: localLocation(filepath.Join("tmp", "same")),
		},
		{
			name:        "nested local path",
			source:      localLocation(filepath.Join("tmp", "source")),
			destination: localLocation(filepath.Join("tmp", "source", "nested")),
		},
		{
			name:        "nested same-session remote path",
			source:      remoteLocation("session", "/source"),
			destination: remoteLocation("session", "/source/nested"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(t, ValidateTransferTopology(tt.source, tt.destination), ErrUnsupported)
		})
	}
}

func TestChecksumSessionPathHonorsCancellationBetweenReads(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "large.bin")
	require.NoError(t, os.WriteFile(filePath, bytes.Repeat([]byte("x"), 2*1024*1024), 0o600))
	baseSession := newLocalTestSession(t, filepath.Dir(filePath))
	session := &slowReadSession{Session: baseSession, delay: 5 * time.Millisecond, chunkSize: 1024}
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(15*time.Millisecond, cancel)

	started := time.Now()
	_, err := checksumSessionPath(ctx, session, NewLocalPath(filePath), OperationCopy)
	require.ErrorIs(t, err, ErrCanceled)
	assert.Less(t, time.Since(started), time.Second)
}

func TestTransferEngineRemoteToLocalDownloadChecksumsAndProcessbar(t *testing.T) {
	fixture := sshtest.Start(t)
	resolver := newTransferTestResolver(fixture)
	engine := NewTransferEngine(resolver)
	bar := newTransferProcessBar(t)

	for _, name := range []string{"alpha.txt", "beta.txt"} {
		sourcePath := "/" + name
		destinationPath := filepath.Join(t.TempDir(), "downloaded-"+name)

		transfer, err := engine.Start(context.Background(), TransferRequest{
			Operation:   OperationTransferRemoteToLocal,
			Source:      remoteLocation(sshtest.AliasE2E, sourcePath),
			Destination: localLocation(destinationPath),
			Overwrite:   true,
		})
		require.NoError(t, err)

		process, err := TrackTransferProcess(context.Background(), bar, transfer)
		require.NoError(t, err)
		require.NoError(t, transfer.Wait(context.Background()))

		assert.Equal(t, sha256RemoteFile(t, resolver, sourcePath), sha256File(t, destinationPath))
		assert.Equal(t, mustReadRemoteFile(t, resolver, sourcePath), mustReadFile(t, destinationPath))
		waitForProcessState(t, bar, process.ID, processbar.Successful)
	}
}

func TestTransferEngineOverwriteSkipPreservesDestination(t *testing.T) {
	fixture := sshtest.Start(t)
	resolver := newTransferTestResolver(fixture)
	engine := NewTransferEngine(resolver)

	sourcePath := filepath.Join(t.TempDir(), "replacement.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("replacement\n"), 0o644))
	original := mustReadRemoteFile(t, resolver, fixture.AlphaPath)

	transfer, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationTransferLocalToRemote,
		Source:      localLocation(sourcePath),
		Destination: remoteLocation(sshtest.AliasE2E, fixture.AlphaPath),
		Overwrite:   false,
	})
	require.NoError(t, err)

	err = transfer.Wait(context.Background())
	require.Error(t, err)
	require.ErrorIs(t, err, ErrConflict)
	assert.Equal(t, original, mustReadRemoteFile(t, resolver, fixture.AlphaPath))
	assertNoTransferTemps(t, resolver)
}

func TestTransferEngineOverwriteUsesVerifiedTempAndAtomicRename(t *testing.T) {
	fixture := sshtest.Start(t)
	resolver := newTransferTestResolver(fixture)
	engine := NewTransferEngine(resolver)

	sourcePath := filepath.Join(t.TempDir(), "replacement.txt")
	content := []byte("replacement content\n")
	require.NoError(t, os.WriteFile(sourcePath, content, 0o644))

	transfer, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationTransferLocalToRemote,
		Source:      localLocation(sourcePath),
		Destination: remoteLocation(sshtest.AliasE2E, fixture.AlphaPath),
		Overwrite:   true,
	})
	require.NoError(t, err)
	require.NoError(t, transfer.Wait(context.Background()))

	assert.Equal(t, content, mustReadRemoteFile(t, resolver, fixture.AlphaPath))
	assertNoTransferTemps(t, resolver)
}

func TestTransferEngineUnsupportedAtomicReplacementPreservesOriginal(t *testing.T) {
	fixture := sshtest.Start(t)
	fixture.FailOperationOnce(
		"posixrename",
		fixture.AlphaPath,
		&sftp.StatusError{Code: uint32(sftp.ErrSSHFxOpUnsupported)},
	)
	resolver := newTransferTestResolver(fixture)
	engine := NewTransferEngine(resolver)

	sourcePath := filepath.Join(t.TempDir(), "replacement.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("unsupported replace\n"), 0o644))
	original := mustReadRemoteFile(t, resolver, fixture.AlphaPath)

	transfer, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationTransferLocalToRemote,
		Source:      localLocation(sourcePath),
		Destination: remoteLocation(sshtest.AliasE2E, fixture.AlphaPath),
		Overwrite:   true,
	})
	require.NoError(t, err)

	err = transfer.Wait(context.Background())
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupported)
	assert.Equal(t, original, mustReadRemoteFile(t, resolver, fixture.AlphaPath))
	assertNoTransferTemps(t, resolver)
}

func TestTransferEngineDisconnectCleanupAndProcessbarFailure(t *testing.T) {
	fixture := sshtest.Start(t)
	resolver := newTransferTestResolver(fixture)
	resolver.(*transferTestResolver).slowLocalReads = true
	engine := NewTransferEngine(resolver)
	bar := newTransferProcessBar(t)

	sourcePath := filepath.Join(t.TempDir(), "large.bin")
	require.NoError(t, os.WriteFile(sourcePath, bytes.Repeat([]byte("superfile-transfer"), 65536), 0o644))

	transfer, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationTransferLocalToRemote,
		Source:      localLocation(sourcePath),
		Destination: remoteLocation(sshtest.AliasE2E, "/large.bin"),
		Overwrite:   true,
	})
	require.NoError(t, err)

	process, err := TrackTransferProcess(context.Background(), bar, transfer)
	require.NoError(t, err)

	go func() {
		time.Sleep(50 * time.Millisecond)
		fixture.CloseActiveConnections()
	}()

	err = transfer.Wait(context.Background())
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDisconnected)

	_, statErr := statRemoteFile(context.Background(), resolver, "/large.bin")
	require.ErrorIs(t, statErr, ErrNotFound)
	assertNoTransferTemps(t, resolver)
	waitForProcessState(t, bar, process.ID, processbar.Failed)
}

func TestTransferEngineSameSessionRemoteCopy(t *testing.T) {
	fixture := sshtest.Start(t)
	resolver := newTransferTestResolver(fixture)
	engine := NewTransferEngine(resolver)

	transfer, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationCopy,
		Source:      remoteLocation(sshtest.AliasE2E, fixture.BetaPath),
		Destination: remoteLocation(sshtest.AliasE2E, "/same-session-copy.txt"),
		Overwrite:   true,
	})
	require.NoError(t, err)
	require.NoError(t, transfer.Wait(context.Background()))

	assert.Equal(
		t,
		mustReadRemoteFile(t, resolver, fixture.BetaPath),
		mustReadRemoteFile(t, resolver, "/same-session-copy.txt"),
	)
}

func TestTransferEngineLocalToLocalCopy(t *testing.T) {
	resolver := newTransferTestResolver(nil)
	engine := NewTransferEngine(resolver)

	sourcePath := filepath.Join(t.TempDir(), "alpha.txt")
	destinationPath := filepath.Join(t.TempDir(), "beta.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("alpha local\n"), 0o644))

	transfer, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationCopy,
		Source:      localLocation(sourcePath),
		Destination: localLocation(destinationPath),
		Overwrite:   true,
	})
	require.NoError(t, err)
	require.NoError(t, transfer.Wait(context.Background()))

	assert.Equal(t, mustReadFile(t, sourcePath), mustReadFile(t, destinationPath))
	assert.Equal(t, sha256File(t, sourcePath), sha256File(t, destinationPath))
}

func TestTransferEngineCrossSessionRemoteUnsupported(t *testing.T) {
	resolverCalled := false
	engine := NewTransferEngine(SessionResolverFunc(func(context.Context, Location) (Session, error) {
		resolverCalled = true
		return nil, errors.New("unexpected resolver call")
	}))

	_, err := engine.Start(context.Background(), TransferRequest{
		Operation:   OperationCopy,
		Source:      remoteLocation("sf-one", "/alpha.txt"),
		Destination: remoteLocation("sf-two", "/beta.txt"),
		Overwrite:   true,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupported)
	assert.False(t, resolverCalled)
}

type transferTestResolver struct {
	localProvider  *LocalProvider
	remoteFactory  func(context.Context, Location) (Session, error)
	slowLocalReads bool
}

func newTransferTestResolver(fixture *sshtest.Fixture) SessionResolver {
	resolver := &transferTestResolver{localProvider: NewLocalProvider()}
	if fixture != nil {
		resolver.remoteFactory = func(ctx context.Context, location Location) (Session, error) {
			provider := NewSFTPProvider(internalssh.ClientConfigRequest{
				Profile:        transferSFTPProfileForAlias(fixture, string(location.SessionID)),
				KnownHostsPath: fixture.KnownHostsPath,
				HostKeyAlias:   string(location.SessionID),
			})
			return provider.Open(ctx, location)
		}
	}
	return resolver
}

func (r *transferTestResolver) ResolveSession(ctx context.Context, location Location) (Session, error) {
	if location.Provider == ProviderLocal {
		session, err := r.localProvider.Open(ctx, location)
		if err != nil {
			return nil, err
		}
		if r.slowLocalReads {
			return &slowReadSession{Session: session, delay: 5 * time.Millisecond, chunkSize: 32 * 1024}, nil
		}
		return session, nil
	}
	if r.remoteFactory == nil {
		return nil, errors.New("remote factory not configured")
	}
	return r.remoteFactory(ctx, location)
}

type slowReadSession struct {
	Session

	delay     time.Duration
	chunkSize int
}

func (s *slowReadSession) Read(ctx context.Context, path Path) (io.ReadCloser, error) {
	reader, err := s.Session.Read(ctx, path)
	if err != nil {
		return nil, err
	}
	return &slowReadCloser{ReadCloser: reader, delay: s.delay, chunkSize: s.chunkSize}, nil
}

type slowReadCloser struct {
	io.ReadCloser

	delay     time.Duration
	chunkSize int
}

func (r *slowReadCloser) Read(p []byte) (int, error) {
	if len(p) > r.chunkSize {
		p = p[:r.chunkSize]
	}
	time.Sleep(r.delay)
	return r.ReadCloser.Read(p)
}

func localLocation(path string) Location {
	return Location{Provider: ProviderLocal, SessionID: "local", Path: NewLocalPath(path), Label: "local"}
}

func remoteLocation(sessionID, path string) Location {
	return Location{
		Provider:  ProviderSFTP,
		SessionID: SessionID(sessionID),
		Path:      NewRemotePath(path),
		Label:     sessionID,
	}
}

func transferSFTPProfileForAlias(fixture *sshtest.Fixture, aliasName string) common.SSHQuickConnectProfile {
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
	}
	if alias.IdentityFilePath != "" {
		profile.IdentityFiles = []string{alias.IdentityFilePath}
	}
	return profile
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func sha256File(t *testing.T, path string) string {
	t.Helper()
	sum := sha256.Sum256(mustReadFile(t, path))
	return hex.EncodeToString(sum[:])
}

func mustReadRemoteFile(t *testing.T, resolver SessionResolver, remotePath string) []byte {
	t.Helper()
	session, err := resolver.ResolveSession(context.Background(), remoteLocation(sshtest.AliasE2E, remotePath))
	require.NoError(t, err)
	defer session.Close()

	reader, err := session.Read(context.Background(), NewRemotePath(remotePath))
	require.NoError(t, err)
	defer reader.Close()

	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	return data
}

func sha256RemoteFile(t *testing.T, resolver SessionResolver, remotePath string) string {
	t.Helper()
	sum := sha256.Sum256(mustReadRemoteFile(t, resolver, remotePath))
	return hex.EncodeToString(sum[:])
}

func statRemoteFile(ctx context.Context, resolver SessionResolver, remotePath string) (Stat, error) {
	session, err := resolver.ResolveSession(ctx, remoteLocation(sshtest.AliasE2E, remotePath))
	if err != nil {
		return Stat{}, err
	}
	defer session.Close()
	return session.Stat(ctx, NewRemotePath(remotePath))
}

func assertNoTransferTemps(t *testing.T, resolver SessionResolver) {
	t.Helper()
	session, err := resolver.ResolveSession(context.Background(), remoteLocation(sshtest.AliasE2E, "/"))
	require.NoError(t, err)
	defer session.Close()

	entries, err := session.List(context.Background(), RootRemotePath())
	require.NoError(t, err)
	for _, entry := range entries {
		assert.NotContains(t, entry.Name, ".superfile-transfer-", "unexpected temp entry %s", entry.Name)
	}
}

func newTransferProcessBar(t *testing.T) *processbar.Model {
	t.Helper()
	model := processbar.New()
	model.ListenForChannelUpdates()
	t.Cleanup(func() {
		model.SendStopListeningMsgBlocking()
	})
	return &model
}

func waitForProcessState(t *testing.T, model *processbar.Model, id string, want processbar.ProcessState) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		process, ok := model.GetByID(id)
		if ok && process.State == want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	process, _ := model.GetByID(id)
	t.Fatalf("process %s state = %v, want %v", id, process.State, want)
}
