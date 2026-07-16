package filemodel

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
)

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
	if got := model.FilePanels[1].DisplayLocation(); got != "sf-e2e:/tmp/sf-remote" {
		t.Fatalf("remote display location = %q, want sf-e2e:/tmp/sf-remote", got)
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
