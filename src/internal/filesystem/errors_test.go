package filesystem

import (
	"errors"
	"testing"
)

func TestUnsupportedRemoteOperationErrorCarriesMetadata(t *testing.T) {
	path := NewRemotePath("/tmp/sf-remote/archive.zip")
	err := V1CapabilityMatrix().RequireRemote(ProviderSFTP, OperationCompress, path)
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("error does not match ErrUnsupported: %v", err)
	}
	var opErr *OperationError
	if !errors.As(err, &opErr) {
		t.Fatalf("error is not *OperationError: %T", err)
	}
	if opErr.Provider != ProviderSFTP {
		t.Fatalf("Provider = %q, want %q", opErr.Provider, ProviderSFTP)
	}
	if opErr.Operation != OperationCompress {
		t.Fatalf("Operation = %q, want %q", opErr.Operation, OperationCompress)
	}
	if got, want := opErr.Path.String(), "/tmp/sf-remote/archive.zip"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
}

func TestUnsupportedRemoteShellErrorCarriesMetadata(t *testing.T) {
	err := V1CapabilityMatrix().RequireRemote(ProviderSFTP, OperationRemoteShell, RootRemotePath())
	var opErr *OperationError
	if !errors.As(err, &opErr) {
		t.Fatalf("error is not *OperationError: %T", err)
	}
	if !errors.Is(err, ErrUnsupported) || opErr.Operation != OperationRemoteShell || opErr.Provider != ProviderSFTP {
		t.Fatalf("unsupported remote shell metadata mismatch: %#v", opErr)
	}
}
