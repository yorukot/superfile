package filesystem

import (
	"context"
	"io"
	"os"
	"time"
)

type ProviderKind string

const (
	ProviderLocal ProviderKind = "local"
	ProviderSFTP  ProviderKind = "sftp"
)

type SessionID string

type Location struct {
	Provider  ProviderKind
	SessionID SessionID
	Path      Path
	Label     string
}

type Entry struct {
	Name string
	Path Path
	Stat Stat
}

type Stat struct {
	Name       string
	Size       int64
	Mode       os.FileMode
	ModTime    time.Time
	IsDir      bool
	IsSymlink  bool
	Owner      string
	Group      string
	Target     Path
	ProviderID string
}

type Provider interface {
	Kind() ProviderKind
	Name() string
	Capabilities() CapabilitySet
	Open(context.Context, Location) (Session, error)
}

type Session interface { //nolint:interfacebloat // Sessions expose the complete provider contract.
	ID() SessionID
	Provider() ProviderKind
	Root() Location
	Capabilities() CapabilitySet
	List(context.Context, Path) ([]Entry, error)
	Stat(context.Context, Path) (Stat, error)
	Read(context.Context, Path) (io.ReadCloser, error)
	Create(context.Context, Path, io.Reader, CreateOptions) error
	Mkdir(context.Context, Path, MkdirOptions) error
	Rename(context.Context, Path, Path, RenameOptions) error
	Delete(context.Context, Path, DeleteOptions) error
	Copy(context.Context, Path, Path, CopyOptions) error
	Move(context.Context, Path, Path, MoveOptions) error
	Chmod(context.Context, Path, os.FileMode) error
	Transfer(context.Context, TransferRequest) (Transfer, error)
	Close() error
}

type CreateOptions struct {
	Mode      os.FileMode
	Overwrite bool
}

type MkdirOptions struct {
	Mode    os.FileMode
	Parents bool
}

type RenameOptions struct {
	Overwrite bool
}

type DeleteOptions struct {
	Recursive bool
	UseTrash  bool
}

type CopyOptions struct {
	Overwrite bool
	Recursive bool
}

type MoveOptions struct {
	Overwrite bool
	Recursive bool
}
