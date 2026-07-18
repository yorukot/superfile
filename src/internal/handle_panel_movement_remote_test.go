package internal

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filemodel"
)

func TestRemoteEnterPanelAndParentDirectoryUsePaneSession(t *testing.T) {
	m := defaultTestModel(t.TempDir())
	session := newMovementRemoteSession()
	remoteLocation := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
		Label:     "ssh://user@sf-e2e",
	}

	m.fileModel.RegisterSession(filemodel.SessionState{
		ID:          "sf-e2e",
		Provider:    filesystem.ProviderSFTP,
		Label:       "ssh://user@sf-e2e",
		CurrentPath: remoteLocation.Path,
		Status:      filemodel.SessionConnected,
		Browser:     session,
	})
	require.NoError(t, m.fileModel.SetPaneLocation(0, remoteLocation))
	m.fileModel.UpdateFilePanelsIfNeeded(true)

	nestedIndex := m.getFocusedFilePanel().FindElementIndexByName("nested")
	require.NotEqual(t, -1, nestedIndex)
	m.getFocusedFilePanel().SetCursorPosition(nestedIndex)

	applyModelUpdateCommand(t, m, m.enterPanel())
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	assert.Equal(t, "/tmp/sf-remote/nested", m.getFocusedFilePanel().Location)
	assert.Equal(t, "gamma.txt", m.getFocusedFilePanel().GetFocusedItem().Name)

	applyModelUpdateCommand(t, m, m.parentDirectory())
	m.fileModel.UpdateFilePanelsIfNeeded(true)
	assert.Equal(t, "/tmp/sf-remote", m.getFocusedFilePanel().Location)
	assert.Equal(t, "nested", m.getFocusedFilePanel().GetFocusedItem().Name)
	assert.Equal(t, "ssh://user@sf-e2e:/tmp/sf-remote", m.getFocusedFilePanel().DisplayLocation())
}

func applyModelUpdateCommand(t *testing.T, m *model, cmd tea.Cmd) {
	t.Helper()
	require.NotNil(t, cmd)
	msg := cmd()
	update, ok := msg.(ModelUpdateMessage)
	require.True(t, ok)
	update.ApplyToModel(m)
}

type movementRemoteSession struct {
	id       filesystem.SessionID
	root     filesystem.Location
	stats    map[string]filesystem.Stat
	children map[string][]filesystem.Entry
}

func newMovementRemoteSession() *movementRemoteSession {
	root := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
	}
	session := &movementRemoteSession{
		id:       root.SessionID,
		root:     root,
		stats:    map[string]filesystem.Stat{},
		children: map[string][]filesystem.Entry{},
	}
	session.addDir("/tmp/sf-remote")
	session.addDir("/tmp/sf-remote/nested")
	session.addFile("/tmp/sf-remote/space name.txt", "space name\n")
	session.addFile("/tmp/sf-remote/nested/gamma.txt", "nested gamma\n")
	return session
}

func (s *movementRemoteSession) ID() filesystem.SessionID          { return s.id }
func (s *movementRemoteSession) Provider() filesystem.ProviderKind { return filesystem.ProviderSFTP }
func (s *movementRemoteSession) Root() filesystem.Location         { return s.root }
func (s *movementRemoteSession) Capabilities() filesystem.CapabilitySet {
	return filesystem.V1CapabilityMatrix()
}

func (s *movementRemoteSession) List(_ context.Context, path filesystem.Path) ([]filesystem.Entry, error) {
	clean := filesystem.NewRemotePath(path.String()).String()
	stat, ok := s.stats[clean]
	if !ok {
		return nil, filesystem.NewNotFoundError(
			filesystem.ProviderSFTP,
			filesystem.OperationList,
			filesystem.NewRemotePath(clean),
			"not found",
		)
	}
	if !stat.IsDir {
		return nil, filesystem.NewNotFoundError(
			filesystem.ProviderSFTP,
			filesystem.OperationList,
			filesystem.NewRemotePath(clean),
			"not a directory",
		)
	}
	return append([]filesystem.Entry(nil), s.children[clean]...), nil
}

func (s *movementRemoteSession) Stat(_ context.Context, path filesystem.Path) (filesystem.Stat, error) {
	clean := filesystem.NewRemotePath(path.String()).String()
	stat, ok := s.stats[clean]
	if !ok {
		return filesystem.Stat{}, filesystem.NewNotFoundError(
			filesystem.ProviderSFTP,
			filesystem.OperationStat,
			filesystem.NewRemotePath(clean),
			"not found",
		)
	}
	return stat, nil
}

func (*movementRemoteSession) Read(context.Context, filesystem.Path) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (*movementRemoteSession) Create(context.Context, filesystem.Path, io.Reader, filesystem.CreateOptions) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Mkdir(context.Context, filesystem.Path, filesystem.MkdirOptions) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Rename(
	context.Context,
	filesystem.Path,
	filesystem.Path,
	filesystem.RenameOptions,
) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Delete(context.Context, filesystem.Path, filesystem.DeleteOptions) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Copy(context.Context, filesystem.Path, filesystem.Path, filesystem.CopyOptions) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Move(context.Context, filesystem.Path, filesystem.Path, filesystem.MoveOptions) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Chmod(context.Context, filesystem.Path, os.FileMode) error {
	return errors.New("not implemented")
}

func (*movementRemoteSession) Transfer(context.Context, filesystem.TransferRequest) (filesystem.Transfer, error) {
	return nil, errors.New("not implemented")
}

func (*movementRemoteSession) Close() error { return nil }

func (s *movementRemoteSession) addDir(path string) {
	clean := filesystem.NewRemotePath(path).String()
	s.stats[clean] = filesystem.Stat{
		Name:       movementBaseName(clean),
		Mode:       os.ModeDir | 0o755,
		ModTime:    time.Unix(0, 0),
		IsDir:      true,
		ProviderID: string(filesystem.ProviderSFTP),
	}
	if _, ok := s.children[clean]; !ok {
		s.children[clean] = []filesystem.Entry{}
	}
	parent := filesystem.NewRemotePath(clean).Dir().String()
	if parent != clean {
		s.children[parent] = append(
			s.children[parent],
			filesystem.Entry{
				Name: movementBaseName(clean),
				Path: filesystem.NewRemotePath(clean),
				Stat: s.stats[clean],
			},
		)
	}
}

func (s *movementRemoteSession) addFile(path string, content string) {
	clean := filesystem.NewRemotePath(path).String()
	stat := filesystem.Stat{
		Name:       movementBaseName(clean),
		Size:       int64(len(content)),
		Mode:       0o644,
		ModTime:    time.Unix(0, 0),
		ProviderID: string(filesystem.ProviderSFTP),
	}
	s.stats[clean] = stat
	parent := filesystem.NewRemotePath(clean).Dir().String()
	s.children[parent] = append(
		s.children[parent],
		filesystem.Entry{Name: stat.Name, Path: filesystem.NewRemotePath(clean), Stat: stat},
	)
}

func movementBaseName(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "/"
	}
	return parts[len(parts)-1]
}
