package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func maskTestPanel() Model {
	return testModel(0, 0, 12, SelectMode, []Element{
		{Name: "notes.txt", Location: "/tmp/notes.txt"},
		{Name: "README", Location: "/tmp/README"},
		{Name: "report.PDF", Location: "/tmp/report.PDF"},
		{Name: "report2.pdf", Location: "/tmp/report2.pdf"},
		{Name: "image.png", Location: "/tmp/image.png"},
		{Name: "subdir", Location: "/tmp/subdir", Directory: true},
	})
}

func TestParseMask(t *testing.T) {
	testdata := []struct {
		name             string
		mask             string
		expectedPatterns []string
		expectedErr      error
	}{
		{
			name:             "single pattern",
			mask:             "*.txt",
			expectedPatterns: []string{"*.txt"},
		},
		{
			name:             "multiple patterns separated by spaces",
			mask:             "*.txt   *.md",
			expectedPatterns: []string{"*.txt", "*.md"},
		},
		{
			name:             "patterns are lowercased for case insensitive matching",
			mask:             "*.TXT",
			expectedPatterns: []string{"*.txt"},
		},
		{
			name:             "norton commander match all is translated",
			mask:             "*.*",
			expectedPatterns: []string{"*"},
		},
		{
			name:        "empty mask",
			mask:        "   ",
			expectedErr: ErrEmptyMask,
		},
		{
			name:        "malformed pattern",
			mask:        "[a-",
			expectedErr: ErrBadMaskPattern,
		},
		{
			// Match reports a malformed pattern even for a name it cannot
			// match, as long as it reaches the malformed part
			name:        "malformed pattern behind a wildcard",
			mask:        "*.txt[a-",
			expectedErr: ErrBadMaskPattern,
		},
		{
			name:        "one malformed pattern invalidates the whole mask",
			mask:        "*.txt [a-",
			expectedErr: ErrBadMaskPattern,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			patterns, err := parseMask(tt.mask)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, patterns)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedPatterns, patterns)
		})
	}
}

func TestSelectByMask(t *testing.T) {
	testdata := []struct {
		name             string
		mask             string
		expectedMatched  int
		expectedSelected []string
	}{
		{
			name:             "extension mask",
			mask:             "*.txt",
			expectedMatched:  1,
			expectedSelected: []string{"/tmp/notes.txt"},
		},
		{
			name:            "mask is case insensitive",
			mask:            "*.pdf",
			expectedMatched: 2,
			// Sorted as visible in the panel
			expectedSelected: []string{"/tmp/report.PDF", "/tmp/report2.pdf"},
		},
		{
			name:             "prefix mask",
			mask:             "report*",
			expectedMatched:  2,
			expectedSelected: []string{"/tmp/report.PDF", "/tmp/report2.pdf"},
		},
		{
			name:             "character class",
			mask:             "report[0-9].pdf",
			expectedMatched:  1,
			expectedSelected: []string{"/tmp/report2.pdf"},
		},
		{
			name:             "multiple patterns",
			mask:             "*.png *.txt",
			expectedMatched:  2,
			expectedSelected: []string{"/tmp/notes.txt", "/tmp/image.png"},
		},
		{
			name:            "match all selects directories too",
			mask:            "*.*",
			expectedMatched: 6,
			expectedSelected: []string{"/tmp/notes.txt", "/tmp/README", "/tmp/report.PDF",
				"/tmp/report2.pdf", "/tmp/image.png", "/tmp/subdir"},
		},
		{
			name:             "no match",
			mask:             "*.go",
			expectedMatched:  0,
			expectedSelected: []string{},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			panel := maskTestPanel()
			matched, err := panel.SelectByMask(tt.mask, true)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMatched, matched)
			assert.ElementsMatch(t, tt.expectedSelected, panel.GetSelectedLocations())
		})
	}
}

func TestSelectByMaskUnselects(t *testing.T) {
	panel := maskTestPanel()

	matched, err := panel.SelectByMask("*.*", true)
	require.NoError(t, err)
	assert.Equal(t, 6, matched)
	assert.Equal(t, uint(6), panel.SelectedCount())

	matched, err = panel.SelectByMask("report*", false)
	require.NoError(t, err)
	assert.Equal(t, 2, matched)
	assert.ElementsMatch(t,
		[]string{"/tmp/notes.txt", "/tmp/README", "/tmp/image.png", "/tmp/subdir"},
		panel.GetSelectedLocations())

	// Unselecting something that was never selected is a no-op
	matched, err = panel.SelectByMask("report*", false)
	require.NoError(t, err)
	assert.Equal(t, 2, matched)
	assert.Equal(t, uint(4), panel.SelectedCount())
}

func TestSelectByMaskInvalidMaskKeepsSelection(t *testing.T) {
	panel := maskTestPanel()
	panel.SetSelected("/tmp/notes.txt")

	matched, err := panel.SelectByMask("", true)
	require.ErrorIs(t, err, ErrEmptyMask)
	assert.Equal(t, 0, matched)

	matched, err = panel.SelectByMask("[a-", true)
	require.ErrorIs(t, err, ErrBadMaskPattern)
	assert.Contains(t, err.Error(), "[a-", "the failing pattern should be in the message shown to the user")
	assert.Equal(t, 0, matched)

	assert.Equal(t, []string{"/tmp/notes.txt"}, panel.GetSelectedLocations())
}

// Mask semantics must not depend on the OS. filepath.Match would disable
// escaping on windows and treat '\' as a separator, this locks in that we
// behave the same on linux, macOS and windows.
func TestSelectByMaskIsOSIndependent(t *testing.T) {
	testdata := []struct {
		name            string
		mask            string
		expectedMatched int
	}{
		{name: "escaped bracket matches it literally", mask: `file\[1\].txt`, expectedMatched: 1},
		{name: "bracket in a character class matches it literally", mask: "file[[]1].txt", expectedMatched: 1},
		{name: "unescaped class does not match the literal bracket", mask: "file[1].txt", expectedMatched: 0},
		{name: "wildcard spans a backslash in the name", mask: "*.txt", expectedMatched: 2},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			// A backslash is a legal file name character on linux and macOS
			panel := testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file[1].txt", Location: "/tmp/file[1].txt"},
				{Name: `back\slash.txt`, Location: `/tmp/back\slash.txt`},
			})
			matched, err := panel.SelectByMask(tt.mask, true)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMatched, matched)
		})
	}
}

// Mask selection works on what the panel shows, so a search filter narrows it down
func TestSelectByMaskOnlyMatchesVisibleElements(t *testing.T) {
	panel := testModel(0, 0, 12, SelectMode, []Element{
		{Name: "report2.pdf", Location: "/tmp/report2.pdf"},
	})

	matched, err := panel.SelectByMask("*.pdf", true)
	require.NoError(t, err)
	assert.Equal(t, 1, matched)
	assert.Equal(t, []string{"/tmp/report2.pdf"}, panel.GetSelectedLocations())
}
