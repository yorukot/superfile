package filepanel

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

var filepanelTestConfigOnce sync.Once //nolint:gochecknoglobals // Package test setup must run once.

func TestRemotePanelListingSearchAndDisplayLocation(t *testing.T) {
	panel := newRemoteTestPanel(t)
	panel.UpdateDimensions(120, 14)
	panel.UpdateElementsIfNeeded(true, false)

	assert.Equal(
		t,
		[]string{"nested", "permission-denied", "alpha.txt", "beta.txt", "readonly.txt", "space name.txt"},
		elementNames(panel.element),
	)
	assert.Contains(t, panel.Render(true), "sf-e2e:/tmp/sf-remote")

	panel.SearchBar.SetValue("space")
	results, err := panel.getElements(false)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "space name.txt", results[0].Name)
	assert.Equal(t, "/tmp/sf-remote/space name.txt", results[0].Location)
}

func TestRemotePanelNavigationPreservesCursorAcrossParentTraversal(t *testing.T) {
	panel := newRemoteTestPanel(t)
	panel.UpdateDimensions(120, 14)
	panel.UpdateElementsIfNeeded(true, false)

	nestedIndex := panel.FindElementIndexByName("nested")
	require.NotEqual(t, -1, nestedIndex)
	panel.SetCursorPosition(nestedIndex)

	require.NoError(t, panel.UpdateCurrentFilePanelDir("nested"))
	panel.UpdateElementsIfNeeded(true, false)
	assert.Equal(t, "/tmp/sf-remote/nested", panel.Location)
	assert.Equal(t, []string{"gamma.txt"}, elementNames(panel.element))

	require.NoError(t, panel.ParentDirectory())
	panel.UpdateElementsIfNeeded(true, false)
	assert.Equal(t, "/tmp/sf-remote", panel.Location)
	assert.Equal(t, "nested", panel.GetFocusedItem().Name)
}

func TestRemotePanelNavigationErrorsKeepPreviousLocationAndElements(t *testing.T) {
	panel := newRemoteTestPanel(t)
	panel.UpdateDimensions(120, 14)
	panel.UpdateElementsIfNeeded(true, false)
	previousElements := append([]Element(nil), panel.element...)

	tests := []struct {
		name        string
		target      string
		expectedErr error
	}{
		{name: "permission denied", target: "permission-denied", expectedErr: filesystem.ErrPermission},
		{name: "missing path", target: "missing", expectedErr: filesystem.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := panel.UpdateCurrentFilePanelDir(tt.target)
			require.Error(t, err)
			require.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, "/tmp/sf-remote", panel.Location)
			assert.Equal(t, elementNames(previousElements), elementNames(panel.element))
		})
	}
}

func newRemoteTestPanel(t *testing.T) *Model {
	t.Helper()
	filepanelTestConfigOnce.Do(func() {
		require.NoError(t, common.PopulateGlobalConfigs())
	})
	session := newStaticRemoteSession()
	panel := New("/tmp/sf-remote", true, "", sortmodel.SortByName, false)
	panel.SetPaneLocation(filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
		Label:     "ssh://user@sf-e2e",
	})
	panel.SetPaneSession(session)
	return &panel
}

type staticRemoteSession struct {
	id           filesystem.SessionID
	root         filesystem.Location
	stats        map[string]filesystem.Stat
	children     map[string][]filesystem.Entry
	inaccessible map[string]error
}

func newStaticRemoteSession() *staticRemoteSession {
	root := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
		Label:     "ssh://user@sf-e2e",
	}
	session := &staticRemoteSession{
		id:           root.SessionID,
		root:         root,
		stats:        map[string]filesystem.Stat{},
		children:     map[string][]filesystem.Entry{},
		inaccessible: map[string]error{},
	}
	session.addDir("/tmp/sf-remote")
	session.addDir("/tmp/sf-remote/nested")
	session.addDir("/tmp/sf-remote/permission-denied")
	session.addFile("/tmp/sf-remote/alpha.txt", "alpha\n")
	session.addFile("/tmp/sf-remote/beta.txt", "beta\n")
	session.addFile("/tmp/sf-remote/readonly.txt", "readonly\n")
	session.addFile("/tmp/sf-remote/space name.txt", "space name\n")
	session.addFile("/tmp/sf-remote/nested/gamma.txt", "nested gamma\n")
	session.inaccessible["/tmp/sf-remote/permission-denied"] = filesystem.NewPermissionError(
		filesystem.ProviderSFTP,
		filesystem.OperationNavigate,
		filesystem.NewRemotePath("/tmp/sf-remote/permission-denied"),
		"permission denied",
	)
	return session
}

func (s *staticRemoteSession) ID() filesystem.SessionID          { return s.id }
func (s *staticRemoteSession) Provider() filesystem.ProviderKind { return filesystem.ProviderSFTP }
func (s *staticRemoteSession) Root() filesystem.Location         { return s.root }
func (s *staticRemoteSession) Capabilities() filesystem.CapabilitySet {
	return filesystem.V1CapabilityMatrix()
}

func (s *staticRemoteSession) List(_ context.Context, path filesystem.Path) ([]filesystem.Entry, error) {
	clean := filesystem.NewRemotePath(path.String()).String()
	if err, ok := s.inaccessible[clean]; ok {
		return nil, err
	}
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
	entries := append([]filesystem.Entry(nil), s.children[clean]...)
	return entries, nil
}

func (s *staticRemoteSession) Stat(_ context.Context, path filesystem.Path) (filesystem.Stat, error) {
	clean := filesystem.NewRemotePath(path.String()).String()
	if err, ok := s.inaccessible[clean]; ok {
		return filesystem.Stat{}, err
	}
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

func (*staticRemoteSession) Read(context.Context, filesystem.Path) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (*staticRemoteSession) Create(context.Context, filesystem.Path, io.Reader, filesystem.CreateOptions) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Mkdir(context.Context, filesystem.Path, filesystem.MkdirOptions) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Rename(context.Context, filesystem.Path, filesystem.Path, filesystem.RenameOptions) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Delete(context.Context, filesystem.Path, filesystem.DeleteOptions) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Copy(context.Context, filesystem.Path, filesystem.Path, filesystem.CopyOptions) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Move(context.Context, filesystem.Path, filesystem.Path, filesystem.MoveOptions) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Chmod(context.Context, filesystem.Path, os.FileMode) error {
	return errors.New("not implemented")
}

func (*staticRemoteSession) Transfer(context.Context, filesystem.TransferRequest) (filesystem.Transfer, error) {
	return nil, errors.New("not implemented")
}

func (*staticRemoteSession) Close() error { return nil }

func (s *staticRemoteSession) addDir(path string) {
	clean := filesystem.NewRemotePath(path).String()
	s.stats[clean] = filesystem.Stat{
		Name:       baseRemoteName(clean),
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
		s.children[parent] = append(s.children[parent], filesystem.Entry{
			Name: baseRemoteName(clean),
			Path: filesystem.NewRemotePath(clean),
			Stat: s.stats[clean],
		})
	}
}

func (s *staticRemoteSession) addFile(path string, content string) {
	clean := filesystem.NewRemotePath(path).String()
	stat := filesystem.Stat{
		Name:       baseRemoteName(clean),
		Size:       int64(len(content)),
		Mode:       0o644,
		ModTime:    time.Unix(0, 0),
		ProviderID: string(filesystem.ProviderSFTP),
	}
	s.stats[clean] = stat
	parent := filesystem.NewRemotePath(clean).Dir().String()
	s.children[parent] = append(s.children[parent], filesystem.Entry{
		Name: stat.Name,
		Path: filesystem.NewRemotePath(clean),
		Stat: stat,
	})
}

func baseRemoteName(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "/"
	}
	return parts[len(parts)-1]
}

func elementNames(elements []Element) []string {
	names := make([]string, 0, len(elements))
	for _, element := range elements {
		names = append(names, element.Name)
	}
	return names
}
