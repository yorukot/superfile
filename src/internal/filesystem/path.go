package filesystem

import (
	"path/filepath"
	"strings"
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
	trimmed := strings.TrimRight(p.value, "/")
	if trimmed == "" {
		return "/"
	}
	idx := strings.LastIndex(trimmed, "/")
	if idx < 0 {
		return trimmed
	}
	return trimmed[idx+1:]
}

func (p Path) Dir() Path {
	if p.kind == PathKindLocal {
		return NewLocalPath(filepath.Dir(p.value))
	}
	trimmed := strings.TrimRight(p.value, "/")
	if trimmed == "" || trimmed == "/" {
		return RootRemotePath()
	}
	idx := strings.LastIndex(trimmed, "/")
	if idx <= 0 {
		return RootRemotePath()
	}
	return NewRemotePath(trimmed[:idx])
}

func (p Path) Join(parts ...string) Path {
	if p.kind == PathKindLocal {
		segments := append([]string{p.value}, parts...)
		return NewLocalPath(filepath.Join(segments...))
	}
	segments := append([]string{p.value}, parts...)
	return NewRemotePath(strings.Join(segments, "/"))
}

func cleanRemotePath(value string) string {
	if value == "" {
		return "/"
	}
	absolute := strings.HasPrefix(value, "/")
	parts := strings.Split(value, "/")
	stack := make([]string, 0, len(parts))
	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			stack = append(stack, part)
		}
	}
	cleaned := strings.Join(stack, "/")
	if absolute {
		cleaned = "/" + cleaned
	}
	if cleaned == "" {
		return "/"
	}
	return cleaned
}
