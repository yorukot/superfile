package filesystem

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestRemotePathPreservesBackslashesInPOSIXNames(t *testing.T) {
	path := NewRemotePath(`/tmp\\sf-remote//nested/../a b.txt`)
	if got, want := path.String(), `/tmp\\sf-remote/a b.txt`; got != want {
		t.Fatalf("remote path normalized with non-POSIX semantics: got %q want %q", got, want)
	}
}

func TestRemotePathEmptyValueDefaultsToRoot(t *testing.T) {
	if got, want := NewRemotePath("").String(), "/"; got != want {
		t.Fatalf("NewRemotePath(\"\") = %q, want %q", got, want)
	}
}

func TestRemotePathJoinDirAndBaseStayPOSIX(t *testing.T) {
	base := NewRemotePath("/tmp/sf-remote")
	joined := base.Join("child", `grandchild\\file.txt`)
	if got, want := joined.String(), `/tmp/sf-remote/child/grandchild\\file.txt`; got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
	if got, want := joined.Dir().String(), "/tmp/sf-remote/child"; got != want {
		t.Fatalf("Dir() = %q, want %q", got, want)
	}
	if got, want := joined.Base(), `grandchild\\file.txt`; got != want {
		t.Fatalf("Base() = %q, want %q", got, want)
	}
}

func TestRemotePathPreservesLeadingParentComponents(t *testing.T) {
	remotePath := NewRemotePath("../../parent/../file.txt")
	if got, want := remotePath.String(), "../../file.txt"; got != want {
		t.Fatalf("remote path = %q, want %q", got, want)
	}
	if got, want := remotePath.Dir().String(), "../.."; got != want {
		t.Fatalf("Dir() = %q, want %q", got, want)
	}
	if got, want := NewRemotePath("../parent").Join("..", "..", "child").String(), "../../child"; got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
}

func TestLocalPathUsesPlatformPathSemantics(t *testing.T) {
	base := NewLocalPath(filepath.Join("tmp", "parent"))
	joined := base.Join("child", "file.txt")
	if got, want := joined.String(), filepath.Join("tmp", "parent", "child", "file.txt"); got != want {
		t.Fatalf("Join() on %s = %q, want %q", runtime.GOOS, got, want)
	}
	if got, want := joined.Dir().String(), filepath.Join("tmp", "parent", "child"); got != want {
		t.Fatalf("Dir() = %q, want %q", got, want)
	}
	if got := joined.Base(); got != "file.txt" {
		t.Fatalf("Base() = %q, want file.txt", got)
	}
}
