package filemodel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

var localSessionProvider = filesystem.NewLocalProvider() //nolint:gochecknoglobals // Shared stateless local provider.

type SessionStatus string

const (
	SessionConnected    SessionStatus = "connected"
	SessionConnecting   SessionStatus = "connecting"
	SessionDisconnected SessionStatus = "disconnected"
)

type SessionState struct {
	ID          filesystem.SessionID
	Provider    filesystem.ProviderKind
	Label       string
	CurrentPath filesystem.Path
	Generation  uint64
	Status      SessionStatus
	LastError   error
	Browser     filesystem.Session
	Reconnect   filesystem.SessionOpener
}

type SessionRegistry struct {
	mu       sync.RWMutex
	sessions map[filesystem.SessionID]SessionState
}

func NewSessionRegistry() *SessionRegistry {
	registry := &SessionRegistry{sessions: make(map[filesystem.SessionID]SessionState)}
	localSession, _ := localSessionProvider.Open(context.Background(), filepanel.NewLocalLocation(""))
	_, _ = registry.Register(SessionState{
		ID:          filepanel.LocalSessionID,
		Provider:    filesystem.ProviderLocal,
		Label:       "local",
		CurrentPath: filesystem.NewLocalPath(""),
		Status:      SessionConnected,
		Browser:     localSession,
	})
	return registry
}

func normalizeSessionState(session SessionState) SessionState {
	if session.ID == "" {
		session.ID = filesystem.SessionID(session.Label)
	}
	if session.Label == "" {
		session.Label = string(session.ID)
	}
	if session.Status == "" {
		session.Status = SessionConnected
	}
	return session
}

func (r *SessionRegistry) Register(session SessionState) (SessionState, bool) {
	session = normalizeSessionState(session)
	r.mu.Lock()
	defer r.mu.Unlock()
	previous, replaced := r.sessions[session.ID]
	session.Generation = 1
	if replaced {
		session.Generation = previous.Generation + 1
	}
	r.sessions[session.ID] = session
	return previous, replaced
}

func (r *SessionRegistry) Get(id filesystem.SessionID) (SessionState, bool) {
	if r == nil {
		return SessionState{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[id]
	return session, ok
}

func (r *SessionRegistry) UpsertLocation(location filesystem.Location) SessionState {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[location.SessionID]
	if !ok {
		session = SessionState{
			ID:       location.SessionID,
			Provider: location.Provider,
			Label:    location.Label,
			Status:   SessionConnected,
		}
	}
	if session.ID == "" {
		session.ID = location.SessionID
	}
	if session.Provider == "" {
		session.Provider = location.Provider
	}
	if location.Label != "" {
		session.Label = location.Label
	}
	session.CurrentPath = location.Path
	session = normalizeSessionState(session)
	r.sessions[session.ID] = session
	return session
}

func (r *SessionRegistry) Remove(id filesystem.SessionID) (SessionState, bool) {
	if r == nil {
		return SessionState{}, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[id]
	if ok {
		delete(r.sessions, id)
	}
	return session, ok
}

func (r *SessionRegistry) Disconnect(id filesystem.SessionID, lastErr error) (filesystem.Session, error) {
	browser, _, err := r.disconnectIfGeneration(id, 0, lastErr)
	return browser, err
}

func (r *SessionRegistry) disconnectIfGeneration(
	id filesystem.SessionID,
	generation uint64,
	lastErr error,
) (filesystem.Session, bool, error) {
	if r == nil {
		return nil, false, fmt.Errorf("unknown session %s", id)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[id]
	if !ok {
		return nil, false, fmt.Errorf("unknown session %s", id)
	}
	if generation != 0 && session.Generation != generation {
		return nil, false, nil
	}
	browser := session.Browser
	session.Browser = nil
	session.Status = SessionDisconnected
	session.LastError = lastErr
	r.sessions[id] = session
	return browser, true, nil
}

func (r *SessionRegistry) DisconnectAll() []filesystem.Session {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	browsers := make([]filesystem.Session, 0, len(r.sessions))
	for id, session := range r.sessions {
		if id == filepanel.LocalSessionID || session.Browser == nil {
			continue
		}
		browsers = append(browsers, session.Browser)
		session.Browser = nil
		session.Status = SessionDisconnected
		r.sessions[id] = session
	}
	return browsers
}

func (m *Model) RegisterSession(session SessionState) {
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	previous, replaced := m.Sessions.Register(session)
	if replaced && previous.Browser != nil && previous.Browser != session.Browser {
		if err := previous.Browser.Close(); err != nil {
			slog.Warn("failed to close replaced filesystem session", "session", session.ID, "error", err)
		}
	}
	for i := range m.FilePanels {
		if m.FilePanels[i].CurrentLocation().SessionID == session.ID {
			m.FilePanels[i].InvalidateElementsLoading()
			m.FilePanels[i].SetPaneSession(session.Browser)
			m.FilePanels[i].SetPaneConnectionStatus(string(session.Status))
		}
	}
}

func (m *Model) CloseSessions() error {
	var closeErr error
	for _, browser := range m.Sessions.DisconnectAll() {
		closeErr = errors.Join(closeErr, browser.Close())
	}
	return closeErr
}

func (m *Model) SetPaneLocation(index int, location filesystem.Location) error {
	if index < 0 || index >= len(m.FilePanels) {
		return fmt.Errorf("panel index %d out of range", index)
	}
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	previousLocation := m.FilePanels[index].CurrentLocation()
	session := m.Sessions.UpsertLocation(location)
	m.FilePanels[index].SetPaneLocation(location)
	m.FilePanels[index].SetPaneSession(session.Browser)
	m.FilePanels[index].SetPaneConnectionStatus(string(session.Status))
	if previousLocation.SessionID != location.SessionID {
		if err := m.closeSessionIfUnused(previousLocation.SessionID); err != nil {
			slog.Warn("failed to close unused filesystem session", "session", previousLocation.SessionID, "error", err)
		}
	}
	return nil
}

func (m *Model) PaneLocation(index int) (filesystem.Location, error) {
	if index < 0 || index >= len(m.FilePanels) {
		return filesystem.Location{}, fmt.Errorf("panel index %d out of range", index)
	}
	return m.FilePanels[index].CurrentLocation(), nil
}

func (m *Model) PaneSession(index int) (SessionState, error) {
	location, err := m.PaneLocation(index)
	if err != nil {
		return SessionState{}, err
	}
	session, ok := m.Sessions.Get(location.SessionID)
	if !ok {
		return SessionState{}, fmt.Errorf("unknown session %s", location.SessionID)
	}
	return session, nil
}

func (m *Model) MarkSessionDisconnected(id filesystem.SessionID, lastErr error) error {
	return m.markSessionDisconnectedIfGeneration(id, 0, lastErr)
}

func (m *Model) MarkSessionDisconnectedIfCurrent(
	id filesystem.SessionID,
	generation uint64,
	lastErr error,
) error {
	return m.markSessionDisconnectedIfGeneration(id, generation, lastErr)
}

func (m *Model) markSessionDisconnectedIfGeneration(
	id filesystem.SessionID,
	generation uint64,
	lastErr error,
) error {
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	browser, matched, err := m.Sessions.disconnectIfGeneration(id, generation, lastErr)
	if err != nil {
		return err
	}
	if !matched {
		return nil
	}
	if browser != nil {
		err = browser.Close()
	}
	for i := range m.FilePanels {
		if m.FilePanels[i].CurrentLocation().SessionID == id {
			m.FilePanels[i].SetPaneSession(nil)
			m.FilePanels[i].SetPaneConnectionStatus(string(SessionDisconnected))
		}
	}
	return err
}

func (m *Model) RequirePaneConnected(index int, operation filesystem.Operation) error {
	session, err := m.PaneSession(index)
	if err != nil {
		return err
	}
	if session.Status != SessionDisconnected {
		return nil
	}
	return filesystem.NewDisconnectedError(session.Provider, operation, session.CurrentPath, "session is disconnected")
}

func (m *Model) SyncPaneSessionLocations() {
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	for i := range m.FilePanels {
		location := m.FilePanels[i].CurrentLocation()
		session := m.Sessions.UpsertLocation(location)
		m.FilePanels[i].SetPaneSession(session.Browser)
		m.FilePanels[i].SetPaneConnectionStatus(string(session.Status))
	}
}

func (m *Model) closeSessionIfUnused(id filesystem.SessionID) error {
	if id == "" || id == filepanel.LocalSessionID || m.Sessions == nil {
		return nil
	}
	for i := range m.FilePanels {
		if m.FilePanels[i].CurrentLocation().SessionID == id {
			return nil
		}
	}
	session, ok := m.Sessions.Remove(id)
	if !ok || session.Browser == nil {
		return nil
	}
	return session.Browser.Close()
}

func ValidateV1Transfer(source, destination filesystem.Location) error {
	return filesystem.ValidateTransferTopology(source, destination)
}
