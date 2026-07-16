package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filemodel"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/pkg/utils"
)

const duplicateSuffixMatchFields = 3

type nonClosingSession struct {
	filesystem.Session
}

func (s nonClosingSession) Close() error {
	return nil
}

func (m *model) ResolveSession(ctx context.Context, location filesystem.Location) (filesystem.Session, error) {
	if location.Provider == filesystem.ProviderLocal {
		return localProvider.Open(ctx, ensureLocationLabel(location))
	}
	return m.resolveRemoteSession(ctx, location, false)
}

func (m *model) ResolveFreshSession(ctx context.Context, location filesystem.Location) (filesystem.Session, error) {
	if location.Provider == filesystem.ProviderLocal {
		return localProvider.Open(ctx, ensureLocationLabel(location))
	}
	return m.resolveRemoteSession(ctx, location, true)
}

func (m *model) resolveRemoteSession(
	ctx context.Context,
	location filesystem.Location,
	fresh bool,
) (filesystem.Session, error) {
	location = ensureLocationLabel(location)
	sessionState, err := m.remoteSessionState(location)
	if err != nil {
		return nil, err
	}

	if fresh && sessionState.Reconnect != nil {
		return sessionState.Reconnect(ctx, location)
	}

	if sessionState.Browser != nil {
		if sessionState.Status == filemodel.SessionDisconnected {
			return nil, disconnectedSessionError(location, sessionState.LastError)
		}
		return nonClosingSession{Session: sessionState.Browser}, nil
	}

	if sessionState.Reconnect != nil {
		return sessionState.Reconnect(ctx, location)
	}

	if sessionState.Status == filemodel.SessionDisconnected {
		return nil, disconnectedSessionError(location, sessionState.LastError)
	}

	return nil, unavailableSessionError(location)
}

func (m *model) remoteSessionState(location filesystem.Location) (filemodel.SessionState, error) {
	registry := m.fileModel.Sessions
	if registry == nil {
		registry = m.sessionRegistry
	}
	if registry == nil {
		return filemodel.SessionState{}, unavailableSessionError(location)
	}

	sessionState, ok := registry[location.SessionID]
	if !ok {
		return filemodel.SessionState{}, unavailableSessionError(location)
	}
	return sessionState, nil
}

func unavailableSessionError(location filesystem.Location) error {
	return filesystem.NewDisconnectedError(location.Provider, filesystem.OperationNavigate, location.Path,
		fmt.Sprintf("session %s is unavailable", location.SessionID))
}

func disconnectedSessionError(location filesystem.Location, lastErr error) error {
	message := "session is disconnected"
	if lastErr != nil {
		message = lastErr.Error()
	}
	return filesystem.NewDisconnectedError(location.Provider, filesystem.OperationNavigate, location.Path, message)
}

func ensureLocationLabel(location filesystem.Location) filesystem.Location {
	if location.Label == "" {
		if location.Provider == filesystem.ProviderLocal {
			location.Label = "local"
		} else {
			location.Label = string(location.SessionID)
		}
	}
	return location
}

func remoteUnsupportedOperationText(provider filesystem.ProviderKind, operation filesystem.Operation) string {
	providerName := strings.ToUpper(string(provider))
	if provider == filesystem.ProviderSFTP {
		providerName = "SFTP"
	}
	return fmt.Sprintf("Operation not supported for %s remote: %s", providerName, operation)
}

func (m *model) unsupportedRemoteOperationCmd(location filesystem.Location, operation filesystem.Operation) tea.Cmd {
	if location.Provider == filesystem.ProviderLocal {
		return nil
	}
	if err := filesystem.V1CapabilityMatrix().RequireRemote(location.Provider, operation, location.Path); err == nil {
		return nil
	}
	message := remoteUnsupportedOperationText(location.Provider, operation)
	reqID := m.nextIoReqCnt()
	return func() tea.Msg {
		return NewNotifyModalMsg(notify.New(true, "Unsupported remote operation", message, notify.NoAction), reqID)
	}
}

func pathJoin(base filesystem.Path, name string) filesystem.Path {
	if base.IsRemote() {
		return base.Join(name)
	}
	return filesystem.NewLocalPath(filepath.Join(base.String(), name))
}

func pathJoinRaw(base filesystem.Path, raw string) filesystem.Path {
	if base.IsRemote() {
		return base.Join(raw)
	}
	return filesystem.NewLocalPath(filepath.Join(base.String(), raw))
}

func pathDir(path filesystem.Path) filesystem.Path {
	if path.IsRemote() {
		return path.Dir()
	}
	return filesystem.NewLocalPath(filepath.Dir(path.String()))
}

func pathBase(path filesystem.Path) string {
	if path.IsRemote() {
		return path.Base()
	}
	return filepath.Base(path.String())
}

func pathFromLocation(location filesystem.Location, raw string) filesystem.Path {
	if location.Provider == filesystem.ProviderSFTP {
		if strings.HasPrefix(raw, "/") {
			return filesystem.NewRemotePath(raw)
		}
		return location.Path.Join(raw)
	}
	if filepath.IsAbs(raw) {
		return filesystem.NewLocalPath(filepath.Clean(raw))
	}
	return filesystem.NewLocalPath(utils.ResolveAbsPath(location.Path.String(), raw))
}

func sameParentDirectory(source filesystem.Location, destination filesystem.Location) bool {
	if source.Provider != destination.Provider || source.SessionID != destination.SessionID {
		return false
	}
	return pathDir(source.Path).String() == destination.Path.String()
}

func isAncestorLocation(source filesystem.Location, destination filesystem.Location) bool {
	if source.Provider != destination.Provider || source.SessionID != destination.SessionID {
		return false
	}
	if source.Path.IsRemote() {
		sourcePath := filesystem.NewRemotePath(source.Path.String()).String()
		destinationPath := filesystem.NewRemotePath(destination.Path.String()).String()
		if sourcePath == destinationPath {
			return true
		}
		if sourcePath == "/" {
			return strings.HasPrefix(destinationPath, "/")
		}
		return strings.HasPrefix(destinationPath, sourcePath+"/")
	}
	return isAncestor(source.Path.String(), destination.Path.String())
}

func renameLocationIfDuplicate(
	ctx context.Context,
	session filesystem.Session,
	location filesystem.Location,
) (filesystem.Location, error) {
	exists, err := sessionPathExists(ctx, session, location.Path)
	if err != nil || !exists {
		return location, err
	}

	parent := pathDir(location.Path)
	base := pathBase(location.Path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	start := 1
	if match := suffixRegexp.FindStringSubmatch(name); len(match) == duplicateSuffixMatchFields {
		name = match[1]
		if next, convErr := strconv.Atoi(match[2]); convErr == nil {
			start = next + 1
		}
	}

	for i := start; i < 10_000; i++ {
		candidate := location
		candidate.Path = pathJoin(parent, fmt.Sprintf("%s(%d)%s", name, i, ext))
		exists, statErr := sessionPathExists(ctx, session, candidate.Path)
		if statErr != nil {
			return filesystem.Location{}, statErr
		}
		if !exists {
			return candidate, nil
		}
	}

	return filesystem.Location{}, fmt.Errorf(
		"could not find free name for %s after many attempts",
		location.Path.String(),
	)
}

func sessionPathExists(ctx context.Context, session filesystem.Session, path filesystem.Path) (bool, error) {
	_, err := session.Stat(ctx, path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, filesystem.ErrNotFound) || os.IsNotExist(err) || errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		if errors.Is(err, filesystem.ErrNotFound) {
			return false, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return false, err
}

func locationWithPath(location filesystem.Location, path filesystem.Path) filesystem.Location {
	location = ensureLocationLabel(location)
	location.Path = path
	return location
}

func elementLocation(base filesystem.Location, element filepanel.Element) filesystem.Location {
	path := element.Path
	if path.String() == "" {
		path = pathFromLocation(base, element.Location)
	}
	return locationWithPath(base, path)
}

func focusedItemLocations(panel *filepanel.Model) []filesystem.Location {
	base := panel.CurrentLocation()
	if panel.Empty() {
		return nil
	}
	if panel.PanelMode == filepanel.SelectMode {
		items := panel.GetSelectedLocationsSortedAsVisible()
		locations := make([]filesystem.Location, 0, len(items))
		for _, item := range items {
			locations = append(locations, locationWithPath(base, pathFromLocation(base, item)))
		}
		return locations
	}
	return []filesystem.Location{elementLocation(base, panel.GetFocusedItem())}
}
