package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestFilePanelSize(t *testing.T) {
	testdata := []struct {
		testName        string
		location        string
		panelSize       int
		columnThreshold int
		expectedColumns int
	}{
		{
			testName:        "happy path",
			location:        "~/test",
			panelSize:       100,
			columnThreshold: 2,
			expectedColumns: 3,
		},
		{
			testName:        "no extra columns",
			location:        "~/test",
			panelSize:       100,
			columnThreshold: 0,
			expectedColumns: 1,
		},
		{
			testName:        "no space for extra columns",
			location:        "~/test",
			panelSize:       40,
			columnThreshold: 2,
			expectedColumns: 1,
		},
	}
	for _, tt := range testdata {
		common.Hotkeys.SearchBar = []string{"test"}
		t.Run(tt.testName, func(t *testing.T) {
			model := New(tt.location, sortOptionsModel{}, true, "test")
			model.SetWidth(tt.panelSize)
			actual := model.makeColumns(tt.columnThreshold)
			assert.Len(t, actual, tt.expectedColumns)
		})
	}
}
