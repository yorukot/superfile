package filesystem

import (
	"runtime"
	"strings"
	"testing"
)

func TestRemotePathUsesPOSIXSeparators(t *testing.T) {
	path := NewRemotePath(`/tmp\\sf-remote//nested/../a b.txt`)
	if got, want := path.String(), "/tmp/sf-remote/a b.txt"; got != want {
		t.Fatalf("remote path normalized with non-POSIX semantics: got %q want %q", got, want)
	}
	if strings.Contains(path.String(), `\\`) {
		t.Fatalf("remote path contains OS separator/backslash on %s: %q", runtime.GOOS, path.String())
	}
}

func TestRemotePathJoinDirAndBaseStayPOSIX(t *testing.T) {
	base := NewRemotePath("/tmp/sf-remote")
	joined := base.Join("child", `grandchild\\file.txt`)
	if got, want := joined.String(), "/tmp/sf-remote/child/grandchild/file.txt"; got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
	if got, want := joined.Dir().String(), "/tmp/sf-remote/child/grandchild"; got != want {
		t.Fatalf("Dir() = %q, want %q", got, want)
	}
	if got, want := joined.Base(), "file.txt"; got != want {
		t.Fatalf("Base() = %q, want %q", got, want)
	}
}
