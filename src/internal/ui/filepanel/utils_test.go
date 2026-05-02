package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

func TestGetSelectedLocationsSortedAsVisible(t *testing.T) {
	testdata := []struct {
		name             string
		panel            Model
		toSelect         []string
		expectedSelected []string
	}{
		{
			name: "no any selected",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
				{Name: "file4.txt", Location: "/tmp/file4.txt"},
			}),
			expectedSelected: []string{},
		},
		{
			name: "1 item selected",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
				{Name: "file4.txt", Location: "/tmp/file4.txt"},
			}),
			toSelect:         []string{"/tmp/file2.txt"},
			expectedSelected: []string{"/tmp/file2.txt"},
		},
		{
			name: "2 item selects reverse selection order",
			panel: testModel(-1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file4.txt", Location: "/tmp/file3.txt"},
				{Name: "file5.txt", Location: "/tmp/file4.txt"},
			}),
			toSelect:         []string{"/tmp/file4.txt", "/tmp/file2.txt"},
			expectedSelected: []string{"/tmp/file2.txt", "/tmp/file4.txt"},
		},
		{
			name: "2 item selects",
			panel: testModel(-1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
				{Name: "file4.txt", Location: "/tmp/file4.txt"},
			}),
			toSelect:         []string{"/tmp/file2.txt", "/tmp/file4.txt"},
			expectedSelected: []string{"/tmp/file2.txt", "/tmp/file4.txt"},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SortKind = sortmodel.SortByName
			tt.panel.SetSelectedAll(tt.toSelect)
			assert.Equal(t, tt.expectedSelected, tt.panel.GetSelectedLocationsSortedAsVisible())
		})
	}
}
