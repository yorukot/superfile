package filepanel

import (
	"fmt"
	"strings"

	"github.com/yorukot/superfile/src/internal/filesystem"
)

func NewLocalLocation(path string) filesystem.Location {
	return filesystem.Location{
		Provider:  filesystem.ProviderLocal,
		SessionID: LocalSessionID,
		Path:      filesystem.NewLocalPath(path),
		Label:     "local",
	}
}

func (m *Model) SetPaneLocation(location filesystem.Location) {
	if location.Label == "" {
		location.Label = string(location.SessionID)
	}
	if m.PaneLocation != location {
		m.InvalidateElementsLoading()
	}
	m.PaneLocation = location
	m.Location = location.Path.String()
}

func (m *Model) CurrentLocation() filesystem.Location {
	if m.PaneLocation.Path.String() == "" {
		return NewLocalLocation(m.Location)
	}
	return m.PaneLocation
}

func (m *Model) DisplayLocation() string {
	location := m.CurrentLocation()
	if location.Provider == filesystem.ProviderLocal {
		return location.Path.String()
	}
	return fmt.Sprintf("%s:%s", locationDisplayLabel(location), location.Path.String())
}

func (m *Model) SetPaneConnectionStatus(status string) {
	m.connectionStatus = strings.TrimSpace(status)
}

func (m *Model) RemoteStatusText() string {
	location := m.CurrentLocation()
	if location.Provider == filesystem.ProviderLocal {
		return ""
	}
	status := m.connectionStatus
	if status == "" {
		status = "connected"
	}
	return fmt.Sprintf("%s %s", m.DisplayLocation(), status)
}

func (m *Model) RemoteSidebarStatusText() string {
	location := m.CurrentLocation()
	if location.Provider == filesystem.ProviderLocal {
		return ""
	}
	status := m.connectionStatus
	if status == "" {
		status = "connected"
	}
	return fmt.Sprintf("%s %s", locationDisplayLabel(location), status)
}

func locationDisplayLabel(location filesystem.Location) string {
	if strings.TrimSpace(location.Label) != "" {
		return location.Label
	}
	return string(location.SessionID)
}

func (m *Model) DisplayLocationWithStatus() string {
	status := m.RemoteStatusText()
	if status == "" {
		return m.DisplayLocation()
	}
	return status
}
