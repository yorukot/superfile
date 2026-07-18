package filesystem

import "testing"

func TestCapabilityMatrixCoversEveryV1Operation(t *testing.T) {
	matrix := V1CapabilityMatrix()
	for _, operation := range V1Operations {
		capability, ok := matrix.Capability(operation)
		if !ok {
			t.Fatalf("operation %q lacks v1 capability classification", operation)
		}
		if capability.Operation != operation {
			t.Fatalf("operation %q has mismatched capability record %q", operation, capability.Operation)
		}
		if capability.Local == "" || capability.Remote == "" {
			t.Fatalf("operation %q has incomplete capability support: %+v", operation, capability)
		}
	}
	if len(matrix) != len(V1Operations) {
		t.Fatalf("capability matrix has %d entries, V1Operations has %d", len(matrix), len(V1Operations))
	}
}

func TestCapabilityMatrixMarksRemoteDeferredOperations(t *testing.T) {
	matrix := V1CapabilityMatrix()
	for _, operation := range []Operation{
		OperationCompress,
		OperationExtract,
		OperationOpenWith,
		OperationMetadata,
		OperationZoxide,
		OperationRemoteShell,
		OperationRemoteCrossSessionMove,
	} {
		capability, ok := matrix.Capability(operation)
		if !ok {
			t.Fatalf("operation %q missing from matrix", operation)
		}
		if capability.Remote != CapabilityDeferred {
			t.Fatalf("operation %q remote support = %q, want %q", operation, capability.Remote, CapabilityDeferred)
		}
	}
}

func TestCapabilityMatrixSupportsRemotePathOperations(t *testing.T) {
	matrix := V1CapabilityMatrix()
	for _, operation := range []Operation{
		OperationList,
		OperationStat,
		OperationNavigate,
		OperationPreviewRead,
		OperationCreateFile,
		OperationMkdir,
		OperationRename,
		OperationDeleteFile,
		OperationDeleteDir,
		OperationCopy,
		OperationCutMove,
		OperationTransferLocalToRemote,
		OperationTransferRemoteToLocal,
		OperationRemoteSameSessionMove,
	} {
		if !matrix.SupportsRemote(operation) {
			t.Fatalf("operation %q should be remote-supported in v1", operation)
		}
	}
}

func TestCapabilityMatrixTreatsLocalOnlyAsLocalSupport(t *testing.T) {
	matrix := V1CapabilityMatrix()
	if !matrix.SupportsLocal(OperationChmod) {
		t.Fatal("local-only chmod should be locally supported")
	}
	if err := matrix.RequireLocal(ProviderLocal, OperationChmod, NewLocalPath("file")); err != nil {
		t.Fatalf("RequireLocal(chmod) returned %v", err)
	}
}
