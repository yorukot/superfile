package internal

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/utils"
)

func TestFilePreviewRenderWithDimensions(t *testing.T) {
	// Test that
	// 1 - we can truncate width and height
	// 2 - We add extra whitespace to make up for width and height
	// 3 - Emojis and special unicodes characters can be rendered and Special characters - ~!@#$%^&*()_+-={}\""
	// 4 - File with spaces, tabs, unicode spaces, etc, is rendered correctly
	// 5 - File with problematic characters like ascii control char, invalid unicodes etc,
	//     is cleaned up

	// Additional tests
	// 1 - File with ascii color sequences can be rendered correctly
	// 2 - Test all cases - unsupported file, non text file
	curTestDir := filepath.Join(testDir, "TestFilePreviewRender")

	// Cleanup is taken care by TestMain()
	utils.SetupDirectories(t, curTestDir)

	testdata := []struct {
		name            string
		fileContent     string
		fileName        string
		height          int
		width           int
		expectedPreview string
	}{
		{
			name: "Basic test",
			fileContent: "" +
				"abcd\n" +
				"1234",
			fileName: "basic.txt",
			height:   2,
			width:    4,
			expectedPreview: "" +
				"abcd\n" +
				"1234",
		},
		{
			name: "Width and height truncation",
			fileContent: "" +
				"abcd\n" +
				"1234\n" +
				"WXYZ",
			fileName: "truncate.txt",
			height:   2,
			width:    3,
			expectedPreview: "" +
				"abc\n" +
				"123",
		},
		{
			name: "Whitespace filling",
			fileContent: "" +
				"abc\n" +
				"123",
			fileName: "fill.txt",
			height:   3,
			width:    4,
			expectedPreview: "" +
				"abc \n" +
				"123 \n" +
				"    ",
		},
		{
			name: "Special char, Emojies and special unicodes",
			fileContent: "" +
				"✅\uf410\U000f0868abcdABCD0123~\n" +
				"!@#$%^&*()_+-={}|:\"<>?,./;'[]",
			fileName: "special.txt",
			height:   2,
			width:    30,
			expectedPreview: "" +
				"✅\uf410\U000f0868abcdABCD0123~             \n" +
				"!@#$%^&*()_+-={}|:\"<>?,./;'[] ",
		},
		{
			// Contains various Unicode whitespace characters:
			// U+00A0 (NO-BREAK SPACE)
			// U+202F (NARROW NO-BREAK SPACE)
			// U+205F (MEDIUM MATHEMATICAL SPACE)
			// U+2029 (PARAGRAPH SEPARATOR)
			name: "Whitespace handling",
			fileContent: "" +
				"\n" +
				"\t1\t\t2\t\n" +
				"0\u00a01\u00a02\u202f3\u205f4\u20295\u202f6\u205f7\u2029\n" +
				"0\u30001\u30002",
			fileName: "whitespace.txt",
			height:   5,
			width:    12,
			expectedPreview: "" +
				"            \n" +
				"    1       \n" +
				"0\u00a01\u00a02 3 4 5 \n" +
				"0 1 2       \n" +
				"            ",
		},
		{
			// Contains control characters:
			// \x0b (Vertical Tab)
			// \x0d (Carriage Return)
			// \x00 (Null)
			// \x05 (Enquiry)
			// \x0f (Shift In)
			// \x7f (Delete)
			// \xa0 (Non-breaking space)
			// \ufffd (Replacement character)
			name: "Invalid character cleanup",
			fileContent: "" +
				"\x0b\x0d\x00\x05\x0f\x7f\xa0\ufffd",
			fileName: "invalid.txt",
			height:   2,
			width:    10,
			expectedPreview: "" +
				"          \n" +
				"          ",
		},
	}

	for i, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			curDir := filepath.Join(curTestDir, "dir"+strconv.Itoa(i))
			utils.SetupDirectories(t, curDir)
			filePath := filepath.Join(curDir, tt.fileName)
			err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
			require.NoError(t, err)

			m := defaultTestModel(curDir)
			m.fileModel.filePreview.SetWidth(tt.width)
			m.fileModel.filePreview.SetHeight(tt.height)
			res := ansi.Strip(m.fileModel.filePreview.RenderWithPath(filePath, m.fullWidth))

			assert.Equal(t, tt.expectedPreview, res, "filePath = %s", filePath)
		})
	}

	// To prevent "normalizeOutput function unused" error.
	_ = normalizeOutput("")
}

// normalizeOutput removes leading empty lines and normalizes line endings
func normalizeOutput(output string) string {
	// Split the output into lines
	lines := strings.Split(output, "\n")

	// Filter out empty lines at the beginning and end
	var filteredLines []string
	startFound := false
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" || startFound {
			startFound = true
			filteredLines = append(filteredLines, line)
		}
	}

	// Remove trailing empty lines
	for len(filteredLines) > 0 && strings.TrimSpace(filteredLines[len(filteredLines)-1]) == "" {
		filteredLines = filteredLines[:len(filteredLines)-1]
	}

	// Join the lines back together
	return strings.Join(filteredLines, "\n")
}
