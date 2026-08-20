package filesystem

type Operation string

const (
	OperationList                   Operation = "list"
	OperationStat                   Operation = "stat"
	OperationNavigate               Operation = "navigate"
	OperationPreviewRead            Operation = "preview-read"
	OperationCreateFile             Operation = "create-file"
	OperationMkdir                  Operation = "mkdir"
	OperationRename                 Operation = "rename"
	OperationDeleteFile             Operation = "delete-file"
	OperationDeleteDir              Operation = "delete-dir"
	OperationCopy                   Operation = "copy"
	OperationCutMove                Operation = "cut-move"
	OperationTransferLocalToRemote  Operation = "transfer-local-to-remote"
	OperationTransferRemoteToLocal  Operation = "transfer-remote-to-local"
	OperationRemoteSameSessionMove  Operation = "remote-same-session-move"
	OperationRemoteCrossSessionMove Operation = "remote-cross-session-move"
	OperationChmod                  Operation = "chmod"
	OperationCompress               Operation = "compress"
	OperationExtract                Operation = "extract"
	OperationOpenWith               Operation = "open-with"
	OperationMetadata               Operation = "metadata"
	OperationZoxide                 Operation = "zoxide"
	OperationRemoteShell            Operation = "remote-shell"
)

var V1Operations = []Operation{ //nolint:gochecknoglobals // Exported immutable operation catalog.
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
	OperationRemoteCrossSessionMove,
	OperationChmod,
	OperationCompress,
	OperationExtract,
	OperationOpenWith,
	OperationMetadata,
	OperationZoxide,
	OperationRemoteShell,
}
