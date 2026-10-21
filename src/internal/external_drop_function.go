package internal

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/ui/notify"
)

// handlePasteMsg imports the files a file manager dropped onto the terminal.
//
// A terminal has no event of its own for a drop: it hands the paths over as
// pasted text, because that is the only way it can. So a paste is read as a
// list of paths when nothing on screen is waiting for text, and left to the
// field that has the keyboard otherwise.
func (m *model) handlePasteMsg(msg tea.PasteMsg) tea.Cmd {
	if m.pasteBelongsToSomethingElse() {
		return nil
	}

	locations := parseDroppedPaths(msg.Content)
	if len(locations) == 0 {
		return nil
	}

	// A paste carries no coordinates, so there is no pointer to aim at. The
	// files land where the keys would have typed: the focused panel.
	return m.getDroppedItemsCmd(m.getFocusedFilePanel().Location, locations)
}

// pasteBelongsToSomethingElse reports whether a paste is already spoken for.
//
// A modal holds the keys while it is open, a text field types what it is given,
// and a drag in progress owns the pointer, so none of them may be interrupted
// by files arriving from outside. Only when nothing does is a paste read as a
// drop.
//
// The modal check is the one the pointer uses, because a modal that captures
// the pointer captures the keys with it.
func (m *model) pasteBelongsToSomethingElse() bool {
	return m.modalBlocksMouse() ||
		m.contextMenu.IsOpen() ||
		m.firstTextInput ||
		m.drag.pending || m.drag.active ||
		m.sidebarModel.SearchBarFocused() ||
		m.getFocusedFilePanel().SearchBar.Focused()
}

// getDroppedItemsCmd copies paths dropped from outside into dest.
//
// It is a copy rather than a move: the files belong to whatever window they
// came from, and taking them away would be a surprise rather than an import.
// It goes through the same processor the paste flow uses, so a drop reports
// progress and refuses the same impossible destinations, and the clipboard the
// user may be relying on is left alone.
func (m *model) getDroppedItemsCmd(dest string, locations []string) tea.Cmd {
	if dest == "" {
		return nil
	}
	locations = droppableInto(dest, locations)
	if len(locations) == 0 {
		return nil
	}

	reqID := m.nextIoReqCnt()
	slog.Debug("Submitting droppedItems request", "id", reqID,
		"items cnt", len(locations), "dest", dest)

	return func() tea.Msg {
		if err := validateDropOperation(dest, locations); err != nil {
			return NewNotifyModalMsg(
				notify.New(true, "Invalid drop location", err.Error(), notify.NoAction), reqID)
		}
		return m.executePasteOperation(&m.processBarModel, dest, locations, false, reqID)
	}
}

// droppableInto keeps only the paths that copying into dest would change.
//
// A path already in dest has nowhere to go but a renamed duplicate, and a path
// listed twice would be copied twice over itself, so both are dropped. What is
// left is validated against dest as a whole, which is what refuses a directory
// being copied into itself.
func droppableInto(dest string, locations []string) []string {
	kept := make([]string, 0, len(locations))
	listed := make(map[string]struct{}, len(locations))
	for _, location := range locations {
		if _, seen := listed[location]; seen {
			continue
		}
		listed[location] = struct{}{}
		if filepath.Dir(location) == dest {
			continue
		}
		kept = append(kept, location)
	}
	return kept
}

// validateDropOperation reports why files cannot be copied into dest.
//
// The destination has to still be a directory, because a panel outlives the
// folder it was showing. Beyond that it is the check a paste makes, so a
// directory is refused as a destination for itself just as it would be by a
// paste.
func validateDropOperation(dest string, locations []string) error {
	info, err := os.Stat(dest)
	if err != nil {
		return fmt.Errorf("cannot drop into %q: %w", dest, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("cannot drop into %q: not a directory", dest)
	}
	return validatePasteOperation(dest, locations, false)
}

// parseDroppedPaths reads a pasted string as the paths a file manager handed
// over.
//
// The shape depends on the terminal: VTE sends text/uri-list with a file:// URI
// per line, kitty sends the paths bare, iTerm2 escapes the spaces in them.
// However it arrives, every path it names has to be absolute and already on
// disk for the paste to count as a drop, so prose pasted with nothing focused
// is left alone. It is all or nothing: one line that does not resolve makes the
// whole paste ordinary text.
func parseDroppedPaths(content string) []string {
	var locations []string
	listed := make(map[string]struct{})

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		// text/uri-list allows blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		resolved := droppedPathsFromLine(line)
		if len(resolved) == 0 {
			return nil
		}
		for _, location := range resolved {
			if _, seen := listed[location]; seen {
				continue
			}
			listed[location] = struct{}{}
			locations = append(locations, location)
		}
	}
	return locations
}

// droppedPathsFromLine reads one line of a paste as one or more paths.
//
// The line is tried as a whole first, so a path with spaces in it is never cut
// in half. Only when that fails is it read as several paths sat next to each
// other, which is how some terminals hand over more than one file.
func droppedPathsFromLine(line string) []string {
	if location, ok := resolveDroppedPath(line); ok {
		return []string{location}
	}

	var locations []string
	for _, candidate := range splitUnescapedFields(line) {
		location, ok := resolveDroppedPath(candidate)
		if !ok {
			return nil
		}
		locations = append(locations, location)
	}
	return locations
}

// resolveDroppedPath reads one candidate as a path to something that exists.
//
// A file manager may hand the path over as a file:// URI, with the spaces in it
// escaped, or wrapped in quotes. All three name the same file, so each way of
// reading it is tried in turn, and none is accepted unless it is an absolute
// path that is on disk.
func resolveDroppedPath(candidate string) (string, bool) {
	for _, location := range droppedPathCandidates(candidate) {
		if !filepath.IsAbs(location) {
			continue
		}
		if _, err := os.Stat(location); err != nil {
			continue
		}
		return filepath.Clean(location), true
	}
	return "", false
}

// droppedPathCandidates lists the ways one candidate may name a file, in the
// order they should be tried. The candidate as it stands comes first, so that a
// path written with backslashes on Windows is read as itself rather than as an
// attempt to escape them.
func droppedPathCandidates(candidate string) []string {
	candidates := []string{candidate}
	if location, ok := decodeFileURI(candidate); ok {
		candidates = append(candidates, location)
	}
	if location, ok := unescapeShellPath(candidate); ok {
		candidates = append(candidates, location)
	}
	if location, ok := unquotePath(candidate); ok {
		candidates = append(candidates, location)
	}
	return candidates
}

// decodeFileURI reads a token as a file:// URI and returns the path it names.
//
// The path arrives percent encoded, so a file called "my file.txt" is handed
// over as "my%20file.txt". A token that is not a file:// URI reports false,
// which leaves it to be read as the plain path it may well be.
func decodeFileURI(token string) (string, bool) {
	rest, ok := strings.CutPrefix(token, "file://")
	if !ok || rest == "" {
		return "", false
	}

	// file://localhost/tmp/x and file:///tmp/x both name /tmp/x. A host other
	// than that names a share this app does not reach into.
	if host, uriPath, found := strings.Cut(rest, "/"); found {
		if host != "" && host != "localhost" {
			return "", false
		}
		rest = "/" + uriPath
	}

	// A drive letter arrives as "/C:/dir"; the leading slash is the URI's root
	// rather than part of the path.
	if withoutRoot, ok := strings.CutPrefix(rest, "/"); ok && filepath.VolumeName(withoutRoot) != "" {
		rest = withoutRoot
	}

	// A path that is not percent encoded at all still names itself.
	if decoded, err := url.PathUnescape(rest); err == nil {
		return decoded, true
	}
	return rest, true
}

// unescapeShellPath removes the backslash a shell puts in front of a space, so
// that a path with spaces in it comes back as one path.
func unescapeShellPath(candidate string) (string, bool) {
	if !strings.Contains(candidate, `\`) {
		return "", false
	}

	var unescaped strings.Builder
	escaped := false
	for _, r := range candidate {
		switch {
		case escaped:
			unescaped.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		default:
			unescaped.WriteRune(r)
		}
	}
	if escaped {
		unescaped.WriteRune('\\')
	}
	return unescaped.String(), true
}

// unquotePath strips the quotes a terminal may wrap a path in.
func unquotePath(candidate string) (string, bool) {
	// An opening and a closing quote is the shortest thing worth unquoting.
	const quotePairLen = 2

	if len(candidate) < quotePairLen {
		return "", false
	}
	quote := candidate[0]
	if (quote != '"' && quote != '\'') || candidate[len(candidate)-1] != quote {
		return "", false
	}
	return candidate[1 : len(candidate)-1], true
}

// splitUnescapedFields cuts a line into the paths it carries, keeping a space
// that a shell escaped as part of the path in front of it.
//
// A backslash in front of anything else is left where it is, because a Windows
// path is full of them and none of them are escapes.
func splitUnescapedFields(line string) []string {
	var fields []string
	var current strings.Builder

	flush := func() {
		if current.Len() > 0 {
			fields = append(fields, current.String())
			current.Reset()
		}
	}

	pendingEscape := false
	for _, r := range line {
		if pendingEscape {
			if r == ' ' || r == '\t' {
				current.WriteRune(r)
			} else {
				current.WriteRune('\\')
				current.WriteRune(r)
			}
			pendingEscape = false
			continue
		}
		switch r {
		case '\\':
			pendingEscape = true
		case ' ', '\t':
			flush()
		default:
			current.WriteRune(r)
		}
	}
	if pendingEscape {
		current.WriteRune('\\')
	}
	flush()
	return fields
}
