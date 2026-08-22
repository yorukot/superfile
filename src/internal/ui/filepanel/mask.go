package filepanel

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// ErrEmptyMask is returned when a selection mask has no pattern in it
var ErrEmptyMask = errors.New("mask is empty")

// ErrBadMaskPattern is returned when a selection mask has a malformed pattern
var ErrBadMaskPattern = errors.New("invalid mask pattern")

// matchAllMask is the Norton Commander style "everything" mask. A "*.*" pattern
// does not match extension-less names like README, but users typing it expect
// every item to be selected.
const matchAllMask = "*.*"

// validationName is matched against while validating a pattern. Match reports a
// malformed pattern even when the name does not match it, but only for the part
// of the pattern it reaches, so this must not be empty.
const validationName = "x"

// Masks are matched with path.Match and not filepath.Match on purpose. We only
// ever match against a file name, which contains no separator on any OS, and
// filepath.Match would give windows users different mask semantics than the
// other platforms - it disables escaping there, so a mask like `file\[1\].txt`
// matches on linux and macOS but never on windows.

// parseMask splits a user provided mask into individual patterns and validates
// them. Patterns are separated by whitespace, so "*.go *.md" is a valid mask
// that matches both Go and Markdown files.
// Returned patterns are lowercased, matching is case insensitive.
func parseMask(mask string) ([]string, error) {
	patterns := strings.Fields(mask)
	if len(patterns) == 0 {
		return nil, ErrEmptyMask
	}
	for i, pattern := range patterns {
		if pattern == matchAllMask {
			pattern = "*"
		}
		if _, err := path.Match(pattern, validationName); err != nil {
			return nil, fmt.Errorf("%w : %s", ErrBadMaskPattern, pattern)
		}
		patterns[i] = strings.ToLower(pattern)
	}
	return patterns, nil
}

// matchesMask reports whether name matches any of the given lowercased patterns
func matchesMask(patterns []string, name string) bool {
	name = strings.ToLower(name)
	for _, pattern := range patterns {
		// Errors are impossible here, patterns are validated by parseMask
		if matched, _ := path.Match(pattern, name); matched {
			return true
		}
	}
	return false
}

// SelectByMask selects (or unselects, when selecting is false) every element of
// the panel whose name matches the given mask. Only elements currently visible
// in the panel are considered, so an active search filter or a hidden dot file
// narrows down what a mask can match.
// It returns the number of elements matching the mask.
func (m *Model) SelectByMask(mask string, selecting bool) (int, error) {
	patterns, err := parseMask(mask)
	if err != nil {
		return 0, err
	}

	matched := 0
	for _, item := range m.element {
		if !matchesMask(patterns, item.Name) {
			continue
		}
		matched++
		if selecting {
			m.SetSelected(item.Location)
		} else {
			m.SetUnSelected(item.Location)
		}
	}
	return matched, nil
}
