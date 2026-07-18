package filemodel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

type closeTrackingSession struct {
	filesystem.Session

	mu         sync.Mutex
	closeCount int
}

type listCountingSession struct {
	closeTrackingSession

	listCount int
}

func (s *listCountingSession) List(_ context.Context, _ filesystem.Path) ([]filesystem.Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listCount++
	return nil, nil
}

func (s *listCountingSession) ID() filesystem.SessionID {
	return "refresh-once"
}

func (s *listCountingSession) Provider() filesystem.ProviderKind {
	return filesystem.ProviderSFTP
}

func (s *listCountingSession) Root() filesystem.Location {
	return filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: s.ID(),
		Path:      filesystem.RootRemotePath(),
	}
}

func (s *listCountingSession) Capabilities() filesystem.CapabilitySet {
	return filesystem.V1CapabilityMatrix()
}

func (s *listCountingSession) lists() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listCount
}

func (s *closeTrackingSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeCount++
	return nil
}

func (s *closeTrackingSession) closes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closeCount
}

var _ filesystem.Session = (*closeTrackingSession)(nil)
var _ filesystem.Session = (*listCountingSession)(nil)

func TestMain(m *testing.M) {
	if err := common.PopulateGlobalConfigs(); err != nil {
		fmt.Printf("error while populating config, err : %v", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestPaneLocationsCanHoldLocalAndFakeRemotePaths(t *testing.T) {
	model := New([]string{"/tmp/sf-local", "/tmp/sf-unused"}, false)
	remote := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
		Label:     "ssh://user@sf-e2e",
	}

	if err := model.SetPaneLocation(1, remote); err != nil {
		t.Fatalf("set remote pane location: %v", err)
	}

	localLocation, err := model.PaneLocation(0)
	if err != nil {
		t.Fatalf("get local pane location: %v", err)
	}
	if localLocation.Provider != filesystem.ProviderLocal {
		t.Fatalf("local provider = %q, want %q", localLocation.Provider, filesystem.ProviderLocal)
	}
	if got := localLocation.Path.String(); got != "/tmp/sf-local" {
		t.Fatalf("local path = %q, want /tmp/sf-local", got)
	}

	remoteLocation, err := model.PaneLocation(1)
	if err != nil {
		t.Fatalf("get remote pane location: %v", err)
	}
	if remoteLocation.Provider != filesystem.ProviderSFTP {
		t.Fatalf("remote provider = %q, want %q", remoteLocation.Provider, filesystem.ProviderSFTP)
	}
	if remoteLocation.SessionID != "sf-e2e" {
		t.Fatalf("remote session = %q, want sf-e2e", remoteLocation.SessionID)
	}
	if got := remoteLocation.Path.String(); got != "/tmp/sf-remote" {
		t.Fatalf("remote path = %q, want /tmp/sf-remote", got)
	}
	if got := model.FilePanels[1].DisplayLocation(); got != "ssh://user@sf-e2e:/tmp/sf-remote" {
		t.Fatalf("remote display location = %q, want ssh://user@sf-e2e:/tmp/sf-remote", got)
	}

	model.NextFilePanel()
	model.PreviousFilePanel()

	localLocation, _ = model.PaneLocation(0)
	remoteLocation, _ = model.PaneLocation(1)
	if got := localLocation.Path.String(); got != "/tmp/sf-local" {
		t.Fatalf("local path after focus switches = %q, want /tmp/sf-local", got)
	}
	if got := remoteLocation.Path.String(); got != "/tmp/sf-remote" {
		t.Fatalf("remote path after focus switches = %q, want /tmp/sf-remote", got)
	}
}

func TestDisconnectedRemotePanePreservesPathAndReturnsTypedState(t *testing.T) {
	model := New([]string{"/tmp/sf-local", "/tmp/sf-unused"}, false)
	remote := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/sf-remote"),
		Label:     "ssh://user@sf-e2e",
	}

	if err := model.SetPaneLocation(1, remote); err != nil {
		t.Fatalf("set remote pane location: %v", err)
	}
	if err := model.MarkSessionDisconnected("sf-e2e", errors.New("network closed")); err != nil {
		t.Fatalf("mark remote disconnected: %v", err)
	}

	session, err := model.PaneSession(1)
	if err != nil {
		t.Fatalf("get remote session: %v", err)
	}
	if session.Status != SessionDisconnected {
		t.Fatalf("session status = %q, want %q", session.Status, SessionDisconnected)
	}
	if got := session.CurrentPath.String(); got != "/tmp/sf-remote" {
		t.Fatalf("disconnected path = %q, want /tmp/sf-remote", got)
	}

	err = model.RequirePaneConnected(1, filesystem.OperationList)
	if !errors.Is(err, filesystem.ErrDisconnected) {
		t.Fatalf("RequirePaneConnected error = %v, want ErrDisconnected", err)
	}
	var operationErr *filesystem.OperationError
	if !errors.As(err, &operationErr) {
		t.Fatalf("RequirePaneConnected error type = %T, want *filesystem.OperationError", err)
	}
	if operationErr.Provider != filesystem.ProviderSFTP {
		t.Fatalf("operation provider = %q, want %q", operationErr.Provider, filesystem.ProviderSFTP)
	}
	if got := operationErr.Path.String(); got != "/tmp/sf-remote" {
		t.Fatalf("operation path = %q, want /tmp/sf-remote", got)
	}
}

func TestCrossSessionRemoteTransferUnsupportedInV1(t *testing.T) {
	source := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-one",
		Path:      filesystem.NewRemotePath("/tmp/source"),
		Label:     "ssh://user@one",
	}
	destination := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-two",
		Path:      filesystem.NewRemotePath("/tmp/destination"),
		Label:     "ssh://user@two",
	}

	err := ValidateV1Transfer(source, destination)
	if !errors.Is(err, filesystem.ErrUnsupported) {
		t.Fatalf("ValidateV1Transfer error = %v, want ErrUnsupported", err)
	}
	var operationErr *filesystem.OperationError
	if !errors.As(err, &operationErr) {
		t.Fatalf("ValidateV1Transfer error type = %T, want *filesystem.OperationError", err)
	}
	if operationErr.Operation != filesystem.OperationRemoteCrossSessionMove {
		t.Fatalf("operation = %q, want %q", operationErr.Operation, filesystem.OperationRemoteCrossSessionMove)
	}
}

func TestSameSessionRemoteTransferAllowedByTopology(t *testing.T) {
	source := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/source"),
	}
	destination := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "sf-e2e",
		Path:      filesystem.NewRemotePath("/tmp/destination"),
	}

	if err := ValidateV1Transfer(source, destination); err != nil {
		t.Fatalf("same-session transfer topology error: %v", err)
	}
}

func TestSessionRegistrySupportsConcurrentReadersAndWriters(t *testing.T) {
	registry := NewSessionRegistry()
	location := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "concurrent",
		Path:      filesystem.RootRemotePath(),
		Label:     "concurrent",
	}

	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			workerLocation := location
			for iteration := range 500 {
				workerLocation.Path = filesystem.NewRemotePath(fmt.Sprintf("/%d/%d", worker, iteration))
				registry.UpsertLocation(workerLocation)
				_, _ = registry.Get(workerLocation.SessionID)
			}
		}(worker)
	}
	wg.Wait()

	_, ok := registry.Get(location.SessionID)
	if !ok {
		t.Fatal("concurrent session was not registered")
	}
}

func TestUnusedRemoteSessionClosesWhenLastPaneStopsUsingIt(t *testing.T) {
	model := New([]string{"/tmp/sf-one", "/tmp/sf-two"}, false)
	remote := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "shared",
		Path:      filesystem.RootRemotePath(),
		Label:     "shared",
	}
	browser := &closeTrackingSession{}
	model.RegisterSession(SessionState{
		ID:       remote.SessionID,
		Provider: remote.Provider,
		Status:   SessionConnected,
		Browser:  browser,
	})

	if err := model.SetPaneLocation(0, remote); err != nil {
		t.Fatalf("set first remote pane: %v", err)
	}
	if err := model.SetPaneLocation(1, remote); err != nil {
		t.Fatalf("set second remote pane: %v", err)
	}

	if err := model.SetPaneLocation(0, filepanel.NewLocalLocation("/tmp/sf-one")); err != nil {
		t.Fatalf("switch first pane local: %v", err)
	}
	if got := browser.closes(); got != 0 {
		t.Fatalf("session close count with one remaining pane = %d, want 0", got)
	}
	if err := model.SetPaneLocation(1, filepanel.NewLocalLocation("/tmp/sf-two")); err != nil {
		t.Fatalf("switch second pane local: %v", err)
	}
	if got := browser.closes(); got != 1 {
		t.Fatalf("session close count after last pane = %d, want 1", got)
	}
	if _, ok := model.Sessions.Get(remote.SessionID); ok {
		t.Fatal("unused session remained registered")
	}
}

func TestMarkSessionDisconnectedClosesBrowserAndUpdatesPane(t *testing.T) {
	model := New([]string{"/tmp/sf-local"}, false)
	remote := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "disconnect-me",
		Path:      filesystem.RootRemotePath(),
		Label:     "disconnect-me",
	}
	browser := &closeTrackingSession{}
	model.RegisterSession(SessionState{
		ID:       remote.SessionID,
		Provider: remote.Provider,
		Status:   SessionConnected,
		Browser:  browser,
	})
	if err := model.SetPaneLocation(0, remote); err != nil {
		t.Fatalf("set remote pane: %v", err)
	}

	disconnectErr := errors.New("network closed")
	if err := model.MarkSessionDisconnected(remote.SessionID, disconnectErr); err != nil {
		t.Fatalf("mark session disconnected: %v", err)
	}
	session, ok := model.Sessions.Get(remote.SessionID)
	if !ok {
		t.Fatal("disconnected session was removed")
	}
	if session.Browser != nil {
		t.Fatal("disconnected session retained its browser")
	}
	if session.Status != SessionDisconnected {
		t.Fatalf("status = %q, want %q", session.Status, SessionDisconnected)
	}
	if got := browser.closes(); got != 1 {
		t.Fatalf("browser close count = %d, want 1", got)
	}
	if got := model.FilePanels[0].RemoteStatusText(); !strings.Contains(got, string(SessionDisconnected)) {
		t.Fatalf("remote status %q does not show disconnected", got)
	}
}

func TestStaleDisconnectDoesNotCloseReplacementSession(t *testing.T) {
	model := New([]string{"/tmp/sf-local"}, false)
	remote := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "replace-session",
	}
	oldBrowser := &closeTrackingSession{}
	model.RegisterSession(SessionState{
		ID:       remote.SessionID,
		Provider: remote.Provider,
		Status:   SessionConnected,
		Browser:  oldBrowser,
	})
	oldState, ok := model.Sessions.Get(remote.SessionID)
	if !ok {
		t.Fatal("old session was not registered")
	}

	newBrowser := &closeTrackingSession{}
	model.RegisterSession(SessionState{
		ID:       remote.SessionID,
		Provider: remote.Provider,
		Status:   SessionConnected,
		Browser:  newBrowser,
	})
	if err := model.MarkSessionDisconnectedIfCurrent(
		remote.SessionID,
		oldState.Generation,
		errors.New("stale disconnect"),
	); err != nil {
		t.Fatalf("mark stale session disconnected: %v", err)
	}

	current, ok := model.Sessions.Get(remote.SessionID)
	if !ok {
		t.Fatal("replacement session was removed")
	}
	if current.Browser != newBrowser || current.Status != SessionConnected {
		t.Fatalf("replacement session changed after stale disconnect: %#v", current)
	}
	if got := oldBrowser.closes(); got != 1 {
		t.Fatalf("old browser close count = %d, want 1", got)
	}
	if got := newBrowser.closes(); got != 0 {
		t.Fatalf("replacement browser close count = %d, want 0", got)
	}
}

func TestSessionRegistryGenerationsRemainMonotonicAfterRemoval(t *testing.T) {
	registry := NewSessionRegistry()
	first, _, _ := registry.Register(SessionState{ID: "first", Label: "first"})
	registry.Remove(first.ID)
	second, _, _ := registry.Register(SessionState{ID: "second", Label: "second"})
	if second.Generation <= first.Generation {
		t.Fatalf("generation after removal = %d, want greater than %d", second.Generation, first.Generation)
	}
}

func TestRegisterSessionUsesNormalizedIDWhenAttachingPanel(t *testing.T) {
	model := New([]string{"/tmp/sf-local"}, false)
	location := filesystem.Location{
		Provider: filesystem.ProviderSFTP,
		Label:    "normalized-session",
		Path:     filesystem.RootRemotePath(),
	}
	if err := model.SetPaneLocation(0, location); err != nil {
		t.Fatalf("set pane location: %v", err)
	}
	browser := &closeTrackingSession{}
	model.RegisterSession(SessionState{
		Provider: filesystem.ProviderSFTP,
		Label:    location.Label,
		Status:   SessionConnected,
		Browser:  browser,
	})

	panelLocation := model.FilePanels[0].CurrentLocation()
	if panelLocation.SessionID != filesystem.SessionID(location.Label) {
		t.Fatalf("panel session ID = %q, want %q", panelLocation.SessionID, location.Label)
	}
	state, err := model.PaneSession(0)
	if err != nil {
		t.Fatalf("get normalized pane session: %v", err)
	}
	if state.Browser != browser {
		t.Fatal("normalized session browser was not attached to panel")
	}
}

func TestRemotePanelCompletionDoesNotImmediatelyScheduleAnotherList(t *testing.T) {
	model := New([]string{"/tmp/sf-local"}, false)
	remote := filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: "refresh-once",
		Path:      filesystem.RootRemotePath(),
		Label:     "refresh-once",
	}
	browser := &listCountingSession{}
	model.RegisterSession(SessionState{
		ID:       remote.SessionID,
		Provider: remote.Provider,
		Status:   SessionConnected,
		Browser:  browser,
	})
	if err := model.SetPaneLocation(0, remote); err != nil {
		t.Fatalf("set remote pane: %v", err)
	}

	cmd := model.GetRemoteFilePanelUpdateCmd(false)
	if cmd == nil {
		t.Fatal("initial remote refresh command is nil")
	}
	rawMsg := cmd()
	msg, ok := rawMsg.(PanelUpdateMsg)
	if !ok {
		t.Fatalf("refresh command returned %T, want PanelUpdateMsg", rawMsg)
	}
	if next := model.ApplyPanelUpdate(msg); next != nil {
		t.Fatal("completed empty listing immediately scheduled another refresh")
	}
	if next := model.GetRemoteFilePanelUpdateCmd(false); next != nil {
		t.Fatal("remote panel refreshed again before the minimum interval")
	}
	if got := browser.lists(); got != 1 {
		t.Fatalf("remote list count = %d, want 1", got)
	}
}
