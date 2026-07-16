package filesystem

import "context"

type TransferDirection string

const (
	TransferLocal    TransferDirection = "local-to-local"
	TransferUpload   TransferDirection = "local-to-remote"
	TransferDownload TransferDirection = "remote-to-local"
	TransferRemote   TransferDirection = "remote-same-session"
)

type TransferID string

type SessionResolver interface {
	ResolveSession(context.Context, Location) (Session, error)
}

type FreshSessionResolver interface {
	ResolveFreshSession(context.Context, Location) (Session, error)
}

type SessionOpener func(context.Context, Location) (Session, error)

type SessionResolverFunc func(context.Context, Location) (Session, error)

func (f SessionResolverFunc) ResolveSession(ctx context.Context, location Location) (Session, error) {
	return f(ctx, location)
}

type TransferRequest struct {
	Operation   Operation
	Direction   TransferDirection
	Source      Location
	Destination Location
	Overwrite   bool
}

type Transfer interface {
	ID() TransferID
	Operation() Operation
	Direction() TransferDirection
	Progress() <-chan Progress
	Cancel(context.Context) error
	Wait(context.Context) error
}

type Progress struct {
	TransferID TransferID
	Operation  Operation
	Current    Path
	Done       int64
	Total      int64
	BytesDone  int64
	BytesTotal int64
	Err        error
}
