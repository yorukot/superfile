package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	internalssh "github.com/yorukot/superfile/src/internal/ssh"
	"github.com/yorukot/superfile/src/internal/ssh/sshtest"
	"github.com/yorukot/superfile/src/pkg/utils"
)

type blockingReadCloser struct {
	unblocked chan struct{}
	started   chan struct{}
	once      sync.Once
	startOnce sync.Once
}

type callbackReader struct {
	reader io.Reader
	onRead func()
	once   sync.Once
}

func (r *callbackReader) Read(p []byte) (int, error) {
	r.once.Do(r.onRead)
	return r.reader.Read(p)
}

func (r *blockingReadCloser) Read(_ []byte) (int, error) {
	r.startOnce.Do(func() { close(r.started) })
	<-r.unblocked
	return 0, io.EOF
}

func (r *blockingReadCloser) Close() error {
	r.once.Do(func() { close(r.unblocked) })
	return nil
}

func TestSFTPProviderContract(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		entries, err := session.List(context.Background(), NewRemotePath("/"))
		require.NoError(t, err)

		names := entryNames(entries)
		assert.Contains(t, names, strings.TrimPrefix(fixture.AlphaPath, "/"))
		assert.Contains(t, names, strings.TrimPrefix(fixture.BetaPath, "/"))
		assert.Contains(t, names, strings.TrimPrefix(fixture.NestedPath, "/"))
	})

	t.Run("stat", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		stat, err := session.Stat(context.Background(), NewRemotePath(fixture.AlphaPath))
		require.NoError(t, err)

		assert.Equal(t, "alpha.txt", stat.Name)
		assert.EqualValues(t, 6, stat.Size)
		assert.False(t, stat.IsDir)
		assert.Equal(t, string(ProviderSFTP), stat.ProviderID)
	})

	t.Run("read", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		reader, err := session.Read(context.Background(), NewRemotePath(fixture.AlphaPath))
		require.NoError(t, err)
		defer reader.Close()

		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "alpha\n", string(data))
	})

	t.Run("write", func(t *testing.T) {
		_, session := newSFTPTestSession(t)
		path := NewRemotePath("/write.txt")

		err := session.Create(context.Background(), path, bytes.NewReader([]byte("writer")), CreateOptions{
			Mode:      utils.UserFilePerm,
			Overwrite: true,
		})
		require.NoError(t, err)

		reader, err := session.Read(context.Background(), path)
		require.NoError(t, err)
		defer reader.Close()
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "writer", string(data))
	})

	t.Run("chmod failure does not prevent write", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)
		path := NewRemotePath("/chmod-fails.txt")
		reader := &callbackReader{
			reader: bytes.NewReader([]byte("written first")),
			onRead: func() {
				stat, statErr := session.Stat(context.Background(), path)
				if assert.NoError(t, statErr) {
					assert.Equal(t, os.FileMode(0o600), stat.Mode.Perm())
				}
				fixture.FailOperationOnce("setstat", path.String(), os.ErrPermission)
			},
		}

		err := session.Create(context.Background(), path, reader, CreateOptions{
			Mode:      utils.UserFilePerm,
			Overwrite: true,
		})
		require.ErrorIs(t, err, ErrPermission)

		readback, readErr := session.Read(context.Background(), path)
		require.NoError(t, readErr)
		defer readback.Close()
		data, readErr := io.ReadAll(readback)
		require.NoError(t, readErr)
		assert.Equal(t, "written first", string(data))
		stat, statErr := session.Stat(context.Background(), path)
		require.NoError(t, statErr)
		assert.Equal(t, os.FileMode(0o600), stat.Mode.Perm())
	})

	t.Run("mkdir", func(t *testing.T) {
		_, session := newSFTPTestSession(t)

		err := session.Mkdir(context.Background(), NewRemotePath("/new/nested"), MkdirOptions{
			Mode:    utils.UserDirPerm,
			Parents: true,
		})
		require.NoError(t, err)

		stat, err := session.Stat(context.Background(), NewRemotePath("/new/nested"))
		require.NoError(t, err)
		assert.True(t, stat.IsDir)
	})

	t.Run("mkdir without parents fails when parent is missing", func(t *testing.T) {
		_, session := newSFTPTestSession(t)
		err := session.Mkdir(
			context.Background(),
			NewRemotePath("/missing-parent/child"),
			MkdirOptions{Parents: false},
		)
		require.Error(t, err)
	})

	t.Run("rename", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		err := session.Rename(
			context.Background(),
			NewRemotePath(fixture.AlphaPath),
			NewRemotePath("/alpha-renamed.txt"),
			RenameOptions{
				Overwrite: true,
			},
		)
		require.NoError(t, err)

		_, err = session.Stat(context.Background(), NewRemotePath(fixture.AlphaPath))
		require.ErrorIs(t, err, ErrNotFound)
		_, err = session.Stat(context.Background(), NewRemotePath("/alpha-renamed.txt"))
		assert.NoError(t, err)
	})

	t.Run("rename without overwrite preserves destination", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)
		err := session.Rename(
			context.Background(),
			NewRemotePath(fixture.AlphaPath),
			NewRemotePath(fixture.BetaPath),
			RenameOptions{Overwrite: false},
		)
		require.ErrorIs(t, err, ErrConflict)

		reader, readErr := session.Read(context.Background(), NewRemotePath(fixture.BetaPath))
		require.NoError(t, readErr)
		defer reader.Close()
		data, readErr := io.ReadAll(reader)
		require.NoError(t, readErr)
		assert.Equal(t, "beta\n", string(data))
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("file", func(t *testing.T) {
			fixture, session := newSFTPTestSession(t)

			err := session.Delete(context.Background(), NewRemotePath(fixture.AlphaPath), DeleteOptions{})
			require.NoError(t, err)

			_, err = session.Stat(context.Background(), NewRemotePath(fixture.AlphaPath))
			require.ErrorIs(t, err, ErrNotFound)
		})

		t.Run("recursive directory", func(t *testing.T) {
			fixture, session := newSFTPTestSession(t)

			err := session.Delete(
				context.Background(),
				NewRemotePath(fixture.NestedPath),
				DeleteOptions{Recursive: true},
			)
			require.NoError(t, err)

			_, err = session.Stat(context.Background(), NewRemotePath(fixture.NestedPath))
			require.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("copy", func(t *testing.T) {
		_, session := newSFTPTestSession(t)

		err := session.Copy(context.Background(), NewRemotePath("/nested"), NewRemotePath("/nested-copy"), CopyOptions{
			Overwrite: true,
			Recursive: true,
		})
		require.NoError(t, err)

		reader, err := session.Read(context.Background(), NewRemotePath("/nested-copy/gamma.txt"))
		require.NoError(t, err)
		defer reader.Close()
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "nested gamma\n", string(data))
	})

	t.Run("copy rejects nested destination before recursion", func(t *testing.T) {
		_, session := newSFTPTestSession(t)
		err := session.Copy(
			context.Background(),
			NewRemotePath("/nested"),
			NewRemotePath("/nested/copy"),
			CopyOptions{Overwrite: true, Recursive: true},
		)
		require.ErrorIs(t, err, ErrUnsupported)
		_, statErr := session.Stat(context.Background(), NewRemotePath("/nested/copy"))
		require.ErrorIs(t, statErr, ErrNotFound)
	})

	t.Run("move", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		err := session.Move(
			context.Background(),
			NewRemotePath(fixture.BetaPath),
			NewRemotePath("/beta-moved.txt"),
			MoveOptions{
				Overwrite: true,
			},
		)
		require.NoError(t, err)

		_, err = session.Stat(context.Background(), NewRemotePath(fixture.BetaPath))
		require.ErrorIs(t, err, ErrNotFound)
		reader, err := session.Read(context.Background(), NewRemotePath("/beta-moved.txt"))
		require.NoError(t, err)
		defer reader.Close()
		data, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, "beta\n", string(data))
	})
}

func TestSFTPProviderCreateAppliesRestrictiveModeBeforeFirstRead(t *testing.T) {
	_, session := newSFTPTestSession(t)
	path := NewRemotePath("/mode-during-read.txt")
	requestedMode := os.FileMode(0o640)
	var observedMode os.FileMode
	var observedModeErr error
	reader := &callbackReader{
		reader: bytes.NewReader([]byte("mode checked")),
		onRead: func() {
			stat, statErr := session.Stat(context.Background(), path)
			observedModeErr = statErr
			if statErr == nil {
				observedMode = stat.Mode.Perm()
			}
		},
	}

	err := session.Create(context.Background(), path, reader, CreateOptions{
		Mode:      requestedMode,
		Overwrite: true,
	})
	require.NoError(t, err)
	require.NoError(t, observedModeErr)
	assert.Equal(t, os.FileMode(0o600), observedMode)

	stat, statErr := session.Stat(context.Background(), path)
	require.NoError(t, statErr)
	assert.Equal(t, requestedMode, stat.Mode.Perm())
}

func TestSFTPProviderErrorMappingAndCapabilities(t *testing.T) {
	t.Run("permission denied remains typed and session usable", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		_, err := session.List(context.Background(), NewRemotePath(fixture.PermissionDeniedPath))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrPermission)

		entries, listErr := session.List(context.Background(), NewRemotePath("/"))
		require.NoError(t, listErr)
		assert.NotEmpty(t, entries)
	})

	t.Run("missing file returns typed not found", func(t *testing.T) {
		_, session := newSFTPTestSession(t)

		_, err := session.Stat(context.Background(), NewRemotePath("/missing.txt"))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)

		var opErr *OperationError
		require.ErrorAs(t, err, &opErr)
		assert.Equal(t, ProviderSFTP, opErr.Provider)
		assert.Equal(t, OperationStat, opErr.Operation)
	})

	t.Run("broken symlink lstat succeeds without following target", func(t *testing.T) {
		_, session := newSFTPTestSession(t)
		sftpSession := requireSFTPSession(t, session)
		require.NoError(t, sftpSession.client.Symlink("/broken-link.txt", "broken-target.txt"))

		stat, err := session.Stat(context.Background(), NewRemotePath("/broken-link.txt"))
		require.NoError(t, err)
		assert.True(t, stat.IsSymlink)
		assert.NotEmpty(t, stat.Target.String())
	})

	t.Run("large directory listing uses cancellable read dir context", func(t *testing.T) {
		_, session := newSFTPTestSession(t)
		require.NoError(t, session.Mkdir(context.Background(), NewRemotePath("/huge"), MkdirOptions{Parents: true}))
		for i := range 160 {
			path := NewRemotePath(fmt.Sprintf("/huge/file-%03d.txt", i))
			require.NoError(
				t,
				session.Create(
					context.Background(),
					path,
					bytes.NewReader([]byte("x")),
					CreateOptions{Overwrite: true},
				),
			)
		}

		entries, err := session.List(context.Background(), NewRemotePath("/huge"))
		require.NoError(t, err)
		assert.Len(t, entries, 160)
	})

	t.Run("canceled listing and read return typed cancellation", func(t *testing.T) {
		fixture, session := newSFTPTestSession(t)

		listCtx, cancelList := context.WithCancel(context.Background())
		cancelList()
		_, err := session.List(listCtx, NewRemotePath("/"))
		require.ErrorIs(t, err, ErrCanceled)

		readCtx, cancelRead := context.WithCancel(context.Background())
		reader, err := session.Read(readCtx, NewRemotePath(fixture.AlphaPath))
		require.NoError(t, err)
		cancelRead()
		_, err = reader.Read(make([]byte, 1))
		require.ErrorIs(t, err, ErrCanceled)
		require.NoError(t, reader.Close())
	})

	t.Run("cancel interrupts an in-flight remote read and invalidates the session", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		blockingReader := &blockingReadCloser{
			unblocked: make(chan struct{}),
			started:   make(chan struct{}),
		}
		reader := contextReadCloser{
			ctx:       ctx,
			path:      NewRemotePath("/blocked"),
			operation: OperationPreviewRead,
			reader:    blockingReader,
			cancel:    blockingReader.Close,
		}
		result := make(chan error, 1)
		go func() {
			_, readErr := reader.Read(make([]byte, 1))
			result <- readErr
		}()

		<-blockingReader.started
		cancel()
		select {
		case err := <-result:
			require.ErrorIs(t, err, ErrCanceled)
			require.ErrorIs(t, err, ErrDisconnected)
		case <-time.After(time.Second):
			t.Fatal("canceled read did not return")
		}
	})

	t.Run("unsupported remote operations carry operation metadata", func(t *testing.T) {
		_, session := newSFTPTestSession(t)

		err := session.Chmod(context.Background(), NewRemotePath("/alpha.txt"), 0o600)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnsupported)

		var opErr *OperationError
		require.ErrorAs(t, err, &opErr)
		assert.Equal(t, ProviderSFTP, opErr.Provider)
		assert.Equal(t, OperationChmod, opErr.Operation)
	})
}

func newSFTPTestSession(t *testing.T) (*sshtest.Fixture, Session) {
	t.Helper()
	t.Setenv("SSH_AUTH_SOCK", "")
	fixture := sshtest.Start(t)
	provider := NewSFTPProvider(internalssh.ClientConfigRequest{
		Profile:        sftpProfileForAlias(fixture, sshtest.AliasE2E),
		KnownHostsPath: fixture.KnownHostsPath,
		HostKeyAlias:   sshtest.AliasE2E,
	})
	session, err := provider.Open(context.Background(), Location{
		Provider:  ProviderSFTP,
		Path:      RootRemotePath(),
		Label:     sshtest.AliasE2E,
		SessionID: SessionID(sshtest.AliasE2E),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, session.Close())
	})
	return fixture, session
}

func sftpProfileForAlias(fixture *sshtest.Fixture, aliasName string) common.SSHQuickConnectProfile {
	alias := fixture.Aliases[aliasName]
	profile := common.SSHQuickConnectProfile{
		Name:          alias.Name,
		Host:          alias.Host,
		Port:          alias.Port,
		User:          alias.User,
		StartPath:     "/",
		IdentityFile:  alias.IdentityFilePath,
		IdentityFiles: nil,
		AuthOrder:     []string{common.SSHAuthMethodPublicKey},
	}
	if alias.IdentityFilePath != "" {
		profile.IdentityFiles = []string{alias.IdentityFilePath}
	}
	return profile
}

func entryNames(entries []Entry) []string {
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name)
	}
	sort.Strings(names)
	return names
}

func requireSFTPSession(t *testing.T, session Session) *SFTPSession {
	t.Helper()
	sftpSession, ok := session.(*SFTPSession)
	require.True(t, ok, "session type is %T", session)
	return sftpSession
}
