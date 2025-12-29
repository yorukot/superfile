package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNaturalCompare(t *testing.T) {
	testCases := []struct {
		name          string
		a             string
		b             string
		caseSensitive bool
		expected      int // -1 for a < b, 0 for a == b, 1 for a > b
	}{
		// Basic string comparisons
		{
			name:          "equal strings",
			a:             "file",
			b:             "file",
			caseSensitive: false,
			expected:      0,
		},
		{
			name:          "lexicographic order",
			a:             "apple",
			b:             "banana",
			caseSensitive: false,
			expected:      -1,
		},

		// Natural number sorting (main feature)
		{
			name:          "file2 before file10",
			a:             "file2",
			b:             "file10",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "file10 after file2",
			a:             "file10",
			b:             "file2",
			caseSensitive: false,
			expected:      1,
		},
		{
			name:          "file1 before file2",
			a:             "file1",
			b:             "file2",
			caseSensitive: false,
			expected:      -1,
		},

		// Multiple numbers in filename
		{
			name:          "chapter1-section2 before chapter1-section10",
			a:             "chapter1-section2",
			b:             "chapter1-section10",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "chapter2-section1 after chapter1-section10",
			a:             "chapter2-section1",
			b:             "chapter1-section10",
			caseSensitive: false,
			expected:      1,
		},

		// Case sensitivity
		{
			name:          "case insensitive A equals a",
			a:             "FileA",
			b:             "filea",
			caseSensitive: false,
			expected:      0,
		},
		{
			name:          "case sensitive A before a",
			a:             "FileA",
			b:             "Filea",
			caseSensitive: true,
			expected:      -1,
		},

		// Edge cases
		{
			name:          "empty strings",
			a:             "",
			b:             "",
			caseSensitive: false,
			expected:      0,
		},
		{
			name:          "empty vs non-empty",
			a:             "",
			b:             "file",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "non-empty vs empty",
			a:             "file",
			b:             "",
			caseSensitive: false,
			expected:      1,
		},
		{
			name:          "only numbers",
			a:             "2",
			b:             "10",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "leading zeros",
			a:             "file001",
			b:             "file1",
			caseSensitive: false,
			expected:      0, // 001 == 1 numerically
		},
		{
			name:          "img1 vs img12",
			a:             "img1",
			b:             "img12",
			caseSensitive: false,
			expected:      -1,
		},

		// Realistic filenames
		{
			name:          "photo-2023-01-15 before photo-2023-02-01",
			a:             "photo-2023-01-15",
			b:             "photo-2023-02-01",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "v1.0.0 before v2.0.0",
			a:             "v1.0.0",
			b:             "v2.0.0",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "v1.9.0 before v1.10.0",
			a:             "v1.9.0",
			b:             "v1.10.0",
			caseSensitive: false,
			expected:      -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := naturalCompare(tc.a, tc.b, tc.caseSensitive)
			if tc.expected < 0 {
				assert.Less(t, result, 0, "expected %s < %s", tc.a, tc.b)
			} else if tc.expected > 0 {
				assert.Greater(t, result, 0, "expected %s > %s", tc.a, tc.b)
			} else {
				assert.Equal(t, 0, result, "expected %s == %s", tc.a, tc.b)
			}
		})
	}
}

func TestExtractNumber(t *testing.T) {
	testCases := []struct {
		name        string
		s           string
		start       int
		expectedNum uint64
		expectedEnd int
	}{
		{
			name:        "single digit",
			s:           "file1.txt",
			start:       4,
			expectedNum: 1,
			expectedEnd: 5,
		},
		{
			name:        "multiple digits",
			s:           "file123.txt",
			start:       4,
			expectedNum: 123,
			expectedEnd: 7,
		},
		{
			name:        "leading zeros",
			s:           "file001.txt",
			start:       4,
			expectedNum: 1,
			expectedEnd: 7,
		},
		{
			name:        "number at start",
			s:           "123file.txt",
			start:       0,
			expectedNum: 123,
			expectedEnd: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			num, end := extractNumber(tc.s, tc.start)
			assert.Equal(t, tc.expectedNum, num)
			assert.Equal(t, tc.expectedEnd, end)
		})
	}
}
