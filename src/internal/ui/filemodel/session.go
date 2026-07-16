package filemodel

import (
	"context"
	"errors"
	"fmt"

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
	Status      SessionStatus
	LastError   error
	Browser     filesystem.Session
	Reconnect   filesystem.SessionOpener
}

type SessionRegistry map[filesystem.SessionID]SessionState

func NewSessionRegistry() SessionRegistry {
	registry := SessionRegistry{}
	localSession, _ := localSessionProvider.Open(context.Background(), filepanel.NewLocalLocation(""))
	registry.Register(SessionState{
		ID:          filepanel.LocalSessionID,
		Provider:    filesystem.ProviderLocal,
		Label:       "local",
		CurrentPath: filesystem.NewLocalPath(""),
		Status:      SessionConnected,
		Browser:     localSession,
	})
	return registry
}

func (r SessionRegistry) Register(session SessionState) {
	if session.ID == "" {
		session.ID = filesystem.SessionID(session.Label)
	}
	if session.Label == "" {
		session.Label = string(session.ID)
	}
	if session.Status == "" {
		session.Status = SessionConnected
	}
	r[session.ID] = session
}

func (r SessionRegistry) UpsertLocation(location filesystem.Location) SessionState {
	session, ok := r[location.SessionID]
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
	r.Register(session)
	return session
}

func (m *Model) RegisterSession(session SessionState) {
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	m.Sessions.Register(session)
	for i := range m.FilePanels {
		if m.FilePanels[i].CurrentLocation().SessionID == session.ID {
			m.FilePanels[i].SetPaneSession(session.Browser)
			m.FilePanels[i].SetPaneConnectionStatus(string(session.Status))
		}
	}
}

func (m *Model) CloseSessions() error {
	var closeErr error
	for id, session := range m.Sessions {
		if id == filepanel.LocalSessionID || session.Browser == nil {
			continue
		}
		closeErr = errors.Join(closeErr, session.Browser.Close())
		session.Browser = nil
		session.Status = SessionDisconnected
		m.Sessions[id] = session
	}
	return closeErr
}

func (r SessionRegistry) MarkDisconnected(id filesystem.SessionID, lastErr error) error {
	session, ok := r[id]
	if !ok {
		return fmt.Errorf("unknown session %s", id)
	}
	session.Status = SessionDisconnected
	session.LastError = lastErr
	r[id] = session
	return nil
}

func (m *Model) SetPaneLocation(index int, location filesystem.Location) error {
	if index < 0 || index >= len(m.FilePanels) {
		return fmt.Errorf("panel index %d out of range", index)
	}
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	session := m.Sessions.UpsertLocation(location)
	m.FilePanels[index].SetPaneLocation(location)
	m.FilePanels[index].SetPaneSession(session.Browser)
	m.FilePanels[index].SetPaneConnectionStatus(string(session.Status))
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
	session, ok := m.Sessions[location.SessionID]
	if !ok {
		return SessionState{}, fmt.Errorf("unknown session %s", location.SessionID)
	}
	return session, nil
}

func (m *Model) MarkSessionDisconnected(id filesystem.SessionID, lastErr error) error {
	if m.Sessions == nil {
		m.Sessions = NewSessionRegistry()
	}
	return m.Sessions.MarkDisconnected(id, lastErr)
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

func ValidateV1Transfer(source, destination filesystem.Location) error {
	return filesystem.ValidateTransferTopology(source, destination)
}
