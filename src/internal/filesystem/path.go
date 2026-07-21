package filesystem

import (
	"path"
	"path/filepath"
)

type PathKind string

const (
	PathKindLocal  PathKind = "local"
	PathKindRemote PathKind = "remote-posix"
)

type Path struct {
	kind  PathKind
	value string
}

func NewLocalPath(value string) Path {
	return Path{kind: PathKindLocal, value: value}
}

func NewRemotePath(value string) Path {
	return Path{kind: PathKindRemote, value: cleanRemotePath(value)}
}

func RootRemotePath() Path {
	return NewRemotePath("/")
}

func (p Path) Kind() PathKind {
	return p.kind
}

func (p Path) String() string {
	return p.value
}

func (p Path) IsRemote() bool {
	return p.kind == PathKindRemote
}

func (p Path) IsLocal() bool {
	return p.kind == PathKindLocal
}

func (p Path) Base() string {
	if p.kind == PathKindLocal {
		return filepath.Base(p.value)
	}
	return path.Base(p.value)
}

func (p Path) Dir() Path {
	if p.kind == PathKindLocal {
		return NewLocalPath(filepath.Dir(p.value))
	}
	return NewRemotePath(path.Dir(p.value))
}

func (p Path) Join(parts ...string) Path {
	if p.kind == PathKindLocal {
		segments := append([]string{p.value}, parts...)
		return NewLocalPath(filepath.Join(segments...))
	}
	segments := append([]string{p.value}, parts...)
	return NewRemotePath(path.Join(segments...))
}

func cleanRemotePath(value string) string {
	if value == "" {
		return "/"
	}
	return path.Clean(value)
}
