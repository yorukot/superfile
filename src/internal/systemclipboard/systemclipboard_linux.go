//go:build linux

package systemclipboard

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MIME types used to exchange file references with Linux file managers.
//
// x-special/gnome-copied-files is understood by the GTK file managers
// (Nautilus, Nemo, Caja, PCManFM, Thunar, ...) and, unlike text/uri-list, it
// also encodes whether the operation was a copy or a cut. We use it as the
// primary format and fall back to text/uri-list when reading.
const (
	gnomeCopiedFilesMime = "x-special/gnome-copied-files"
	uriListMime          = "text/uri-list"
)

// linuxTool describes how to drive an external clipboard helper. On Linux the
// clipboard is served by the owning process, so a short-lived TUI cannot hold a
// selection itself. Both wl-copy and xclip fork into the background and keep
// serving the selection after we return, which is exactly what we need.
type linuxTool struct {
	name      string
	copyArgs  func(mime string) []string
	pasteArgs func(mime string) []string
}

func waylandTool() linuxTool {
	return linuxTool{
		name: "wl-clipboard",
		copyArgs: func(mime string) []string {
			return []string{"wl-copy", "--type", mime}
		},
		pasteArgs: func(mime string) []string {
			return []string{"wl-paste", "--no-newline", "--type", mime}
		},
	}
}

func xclipTool() linuxTool {
	return linuxTool{
		name: "xclip",
		copyArgs: func(mime string) []string {
			return []string{"xclip", "-selection", "clipboard", "-target", mime, "-in"}
		},
		pasteArgs: func(mime string) []string {
			return []string{"xclip", "-selection", "clipboard", "-target", mime, "-out"}
		},
	}
}

// detectTool picks the best available clipboard helper for the current session.
func detectTool() (linuxTool, error) {
	waylandReady := hasBinary("wl-copy") && hasBinary("wl-paste")

	// Prefer Wayland tooling when we are in a Wayland session.
	if os.Getenv("WAYLAND_DISPLAY") != "" && waylandReady {
		return waylandTool(), nil
	}

	if hasBinary("xclip") {
		return xclipTool(), nil
	}

	// Fall back to wl-clipboard even without WAYLAND_DISPLAY (e.g. XWayland).
	if waylandReady {
		return waylandTool(), nil
	}

	return linuxTool{}, fmt.Errorf(
		"%w: install wl-clipboard (Wayland) or xclip (X11) to enable system clipboard file transfer",
		ErrUnsupported)
}

func hasBinary(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Available reports whether a supported clipboard helper is present.
func Available() bool {
	_, err := detectTool()
	return err == nil
}

// CopyFiles places the given paths on the system clipboard using the
// gnome-copied-files format, tagged as a cut or copy operation.
func CopyFiles(paths []string, cut bool) error {
	tool, err := detectTool()
	if err != nil {
		return err
	}
	abs, err := absPaths(paths)
	if err != nil {
		return err
	}
	if len(abs) == 0 {
		return ErrNoFiles
	}

	payload := buildGnomeCopiedFiles(abs, cut)
	return runCopy(tool, gnomeCopiedFilesMime, []byte(payload))
}

// PasteFiles reads file references from the system clipboard.
func PasteFiles() ([]string, bool, error) {
	tool, err := detectTool()
	if err != nil {
		return nil, false, err
	}

	// Primary: gnome-copied-files (carries the cut/copy flag).
	if out, perr := runPaste(tool, gnomeCopiedFilesMime); perr == nil {
		if paths, cut, ok := parseGnomeCopiedFiles(out); ok {
			return paths, cut, nil
		}
	}

	// Fallback: plain text/uri-list (copy semantics only).
	if out, perr := runPaste(tool, uriListMime); perr == nil {
		if paths := parseURIList(out); len(paths) > 0 {
			return paths, false, nil
		}
	}

	return nil, false, ErrNoFiles
}

func runCopy(tool linuxTool, mime string, payload []byte) error {
	// A buffer would make os/exec wait for EOF on a pipe inherited by the
	// background clipboard owner. A file lets Run return when the parent exits.
	stderr, err := os.CreateTemp("", "superfile-clipboard-stderr-*")
	if err != nil {
		return fmt.Errorf("%s copy failed: %w", tool.name, err)
	}
	defer os.Remove(stderr.Name())
	defer stderr.Close()

	args := tool.copyArgs(mime)
	// #nosec G204 -- args come from a fixed allow-list, only the payload is dynamic.
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = bytes.NewReader(payload)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		diagnostic, _ := os.ReadFile(stderr.Name())
		return fmt.Errorf("%s copy failed: %w: %s", tool.name, err, strings.TrimSpace(string(diagnostic)))
	}
	return nil
}

func runPaste(tool linuxTool, mime string) ([]byte, error) {
	args := tool.pasteArgs(mime)
	// #nosec G204 -- args come from a fixed allow-list.
	cmd := exec.Command(args[0], args[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

// --- payload formatting helpers (pure functions, unit tested) ---

func absPaths(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", p, err)
		}
		out = append(out, abs)
	}
	return out, nil
}

// pathToURI converts an absolute filesystem path into a file:// URI with proper
// percent-encoding.
func pathToURI(abs string) string {
	u := url.URL{Scheme: "file", Path: abs}
	return u.String()
}

// uriToPath converts a file:// URI (or a bare path) back into a filesystem path.
func uriToPath(uri string) string {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return ""
	}
	u, err := url.Parse(uri)
	if err != nil || u.Scheme == "" {
		// Not a URI; treat it as a raw path.
		return uri
	}
	if u.Scheme != "file" {
		return ""
	}
	return u.Path
}

func buildGnomeCopiedFiles(paths []string, cut bool) string {
	op := "copy"
	if cut {
		op = "cut"
	}
	lines := make([]string, 0, len(paths)+1)
	lines = append(lines, op)
	for _, p := range paths {
		lines = append(lines, pathToURI(p))
	}
	return strings.Join(lines, "\n")
}

func parseGnomeCopiedFiles(data []byte) (paths []string, cut bool, ok bool) {
	text := strings.TrimRight(string(data), "\x00\r\n")
	if text == "" {
		return nil, false, false
	}
	lines := strings.Split(text, "\n")
	op := strings.TrimSpace(lines[0])
	switch op {
	case "cut":
		cut = true
	case "copy":
		cut = false
	default:
		return nil, false, false
	}
	for _, line := range lines[1:] {
		if p := uriToPath(line); p != "" {
			paths = append(paths, p)
		}
	}
	return paths, cut, len(paths) > 0
}

func parseURIList(data []byte) []string {
	var paths []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if p := uriToPath(line); p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}
