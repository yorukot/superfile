package filepanel

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/pkg/utils"
)

var localBrowserProvider = filesystem.NewLocalProvider() //nolint:gochecknoglobals // Shared stateless local provider.

func (m *Model) SetPaneSession(session filesystem.Session) {
	m.session = session
}

func (m *Model) paneSession() (filesystem.Session, error) {
	location := m.CurrentLocation()
	if m.session != nil {
		if m.session.Provider() == location.Provider && m.session.ID() == location.SessionID {
			return m.session, nil
		}
		m.session = nil
	}

	if location.Provider != filesystem.ProviderLocal {
		return nil, filesystem.NewDisconnectedError(location.Provider, filesystem.OperationNavigate, location.Path,
			fmt.Sprintf("session %s is unavailable", location.SessionID))
	}

	session, err := localBrowserProvider.Open(context.Background(), location)
	if err != nil {
		return nil, err
	}
	m.session = session
	return session, nil
}

func locationKey(location filesystem.Location) string {
	return fmt.Sprintf("%s|%s|%s", location.Provider, location.SessionID, location.Path.String())
}

func pathFromLocation(location filesystem.Location, raw string) filesystem.Path {
	if location.Provider == filesystem.ProviderSFTP {
		return filesystem.NewRemotePath(raw)
	}
	return filesystem.NewLocalPath(raw)
}

func resolveLocationPath(location filesystem.Location, raw string) filesystem.Path {
	if location.Provider == filesystem.ProviderSFTP {
		if strings.HasPrefix(raw, "/") {
			return filesystem.NewRemotePath(raw)
		}
		return location.Path.Join(raw)
	}
	return filesystem.NewLocalPath(utils.ResolveAbsPath(location.Path.String(), raw))
}

func joinPath(base filesystem.Path, name string) filesystem.Path {
	if base.IsRemote() {
		return base.Join(name)
	}
	return filesystem.NewLocalPath(filepath.Join(base.String(), name))
}

func parentPath(path filesystem.Path) filesystem.Path {
	if path.IsRemote() {
		return path.Dir()
	}
	return filesystem.NewLocalPath(filepath.Dir(path.String()))
}

func baseName(path filesystem.Path) string {
	if path.IsRemote() {
		return path.Base()
	}
	return filepath.Base(path.String())
}

func elementPath(element Element, location filesystem.Location) filesystem.Path {
	if element.Path.String() != "" {
		return element.Path
	}
	return pathFromLocation(location, element.Location)
}

func resolveSymlinkTargetPath(entryPath filesystem.Path, target filesystem.Path) filesystem.Path {
	targetValue := target.String()
	if targetValue == "" {
		return entryPath
	}
	base := parentPath(entryPath)
	if entryPath.IsRemote() {
		if strings.HasPrefix(targetValue, "/") {
			return filesystem.NewRemotePath(targetValue)
		}
		return base.Join(targetValue)
	}
	if filepath.IsAbs(targetValue) {
		return filesystem.NewLocalPath(filepath.Clean(targetValue))
	}
	return filesystem.NewLocalPath(utils.ResolveAbsPath(base.String(), targetValue))
}
