package filesystem

type CapabilitySupport string

const (
	CapabilitySupported   CapabilitySupport = "supported"
	CapabilityLocalOnly   CapabilitySupport = "local-only"
	CapabilityDeferred    CapabilitySupport = "deferred"
	CapabilityUnsupported CapabilitySupport = "unsupported"
)

type Capability struct {
	Operation Operation
	Local     CapabilitySupport
	Remote    CapabilitySupport
	Notes     string
}

type CapabilitySet map[Operation]Capability

func (set CapabilitySet) Capability(operation Operation) (Capability, bool) {
	capability, ok := set[operation]
	return capability, ok
}

func (set CapabilitySet) SupportsRemote(operation Operation) bool {
	capability, ok := set[operation]
	return ok && capability.Remote == CapabilitySupported
}

func (set CapabilitySet) SupportsLocal(operation Operation) bool {
	capability, ok := set[operation]
	return ok && capability.Local == CapabilitySupported
}

func (set CapabilitySet) RequireRemote(provider ProviderKind, operation Operation, path Path) error {
	capability, ok := set[operation]
	if ok && capability.Remote == CapabilitySupported {
		return nil
	}
	message := "remote operation is not supported"
	if ok && capability.Notes != "" {
		message = capability.Notes
	}
	return NewUnsupportedError(provider, operation, path, message)
}

func (set CapabilitySet) RequireLocal(provider ProviderKind, operation Operation, path Path) error {
	capability, ok := set[operation]
	if ok && capability.Local == CapabilitySupported {
		return nil
	}
	message := "local operation is not supported"
	if ok && capability.Notes != "" {
		message = capability.Notes
	}
	return NewUnsupportedError(provider, operation, path, message)
}

func V1CapabilityMatrix() CapabilitySet { //nolint:funlen // Declarative operation matrix is clearer in one literal.
	return CapabilitySet{
		OperationList: {
			Operation: OperationList,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "directory listing",
		},
		OperationStat: {
			Operation: OperationStat,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "file metadata for panels and previews",
		},
		OperationNavigate: {
			Operation: OperationNavigate,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "change current pane path",
		},
		OperationPreviewRead: {
			Operation: OperationPreviewRead,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "read file content for preview",
		},
		OperationCreateFile: {
			Operation: OperationCreateFile,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "create or overwrite a regular file",
		},
		OperationMkdir: {
			Operation: OperationMkdir,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "create remote directories over SFTP",
		},
		OperationRename: {
			Operation: OperationRename,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "same-directory or same-session rename",
		},
		OperationDeleteFile: {
			Operation: OperationDeleteFile,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "delete regular files",
		},
		OperationDeleteDir: {
			Operation: OperationDeleteDir,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "delete directories recursively when requested",
		},
		OperationCopy: {
			Operation: OperationCopy,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "provider-local copy or transfer orchestration",
		},
		OperationCutMove: {
			Operation: OperationCutMove,
			Local:     CapabilitySupported,
			Remote:    CapabilitySupported,
			Notes:     "provider-local move or transfer orchestration",
		},
		OperationTransferLocalToRemote: {
			Operation: OperationTransferLocalToRemote,
			Local:     CapabilityUnsupported,
			Remote:    CapabilitySupported,
			Notes:     "stream upload from local provider into an SFTP session",
		},
		OperationTransferRemoteToLocal: {
			Operation: OperationTransferRemoteToLocal,
			Local:     CapabilityUnsupported,
			Remote:    CapabilitySupported,
			Notes:     "stream download from an SFTP session into local provider",
		},
		OperationRemoteSameSessionMove: {
			Operation: OperationRemoteSameSessionMove,
			Local:     CapabilityUnsupported,
			Remote:    CapabilitySupported,
			Notes:     "same-session SFTP rename/move only",
		},
		OperationRemoteCrossSessionMove: {
			Operation: OperationRemoteCrossSessionMove,
			Local:     CapabilityUnsupported,
			Remote:    CapabilityDeferred,
			Notes:     "cross-session remote to remote transfer is deferred for v1",
		},
		OperationChmod: {
			Operation: OperationChmod,
			Local:     CapabilityLocalOnly,
			Remote:    CapabilityDeferred,
			Notes:     "not exposed in current filepanel flow; remote chmod waits for explicit UI exposure",
		},
		OperationCompress: {
			Operation: OperationCompress,
			Local:     CapabilitySupported,
			Remote:    CapabilityDeferred,
			Notes:     "current archive implementation is local-only",
		},
		OperationExtract: {
			Operation: OperationExtract,
			Local:     CapabilitySupported,
			Remote:    CapabilityDeferred,
			Notes:     "current extraction implementation is local-only",
		},
		OperationOpenWith: {
			Operation: OperationOpenWith,
			Local:     CapabilitySupported,
			Remote:    CapabilityDeferred,
			Notes:     "local editor/process execution only",
		},
		OperationMetadata: {
			Operation: OperationMetadata,
			Local:     CapabilitySupported,
			Remote:    CapabilityDeferred,
			Notes:     "metadata UI currently depends on local file access",
		},
		OperationZoxide: {
			Operation: OperationZoxide,
			Local:     CapabilitySupported,
			Remote:    CapabilityDeferred,
			Notes:     "zoxide is local navigation state",
		},
		OperationRemoteShell: {
			Operation: OperationRemoteShell,
			Local:     CapabilityUnsupported,
			Remote:    CapabilityDeferred,
			Notes:     "remote command execution is out of v1 scope",
		},
	}
}
