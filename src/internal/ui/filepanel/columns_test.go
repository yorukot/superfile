package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestFilePanelSize(t *testing.T) {
	common.Hotkeys.SearchBar = []string{"test"}
	testdata := []struct {
		testName             string
		location             string
		panelSize            int
		columnThreshold      int
		expectedColumns      int
		filePanelNamePercent int
	}{
		{
			testName:             "happy path",
			location:             "~/test",
			panelSize:            100,
			columnThreshold:      2,
			expectedColumns:      3,
			filePanelNamePercent: 65,
		},
		{
			testName:             "no extra columns",
			location:             "~/test",
			panelSize:            100,
			columnThreshold:      0,
			expectedColumns:      1,
			filePanelNamePercent: 65,
		},
		{
			testName:             "no space for extra columns",
			location:             "~/test",
			panelSize:            40,
			columnThreshold:      2,
			expectedColumns:      1,
			filePanelNamePercent: 65,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.testName, func(t *testing.T) {
			model := New(tt.location, sortOptionsModel{}, true, "test")
			model.SetWidth(tt.panelSize)
			actual := model.makeColumns(tt.columnThreshold, tt.filePanelNamePercent)
			assert.Len(t, actual, tt.expectedColumns)
		})
	}
}
