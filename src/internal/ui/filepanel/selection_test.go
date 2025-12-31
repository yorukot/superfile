package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanelSelectionLifeCycle(t *testing.T) {
	panel := testModel(0, 0, 0, BrowserMode, []Element{
		{Name: "file1.txt", Location: "/tmp/file1.txt"},
		{Name: "file2.txt", Location: "/tmp/file2.txt"},
		{Name: "file3.txt", Location: "/tmp/file3.txt"},
		{Name: "file4.txt", Location: "/tmp/file4.txt"},
		{Name: "file5.txt", Location: "/tmp/file5.txt"}})
	assert.Equal(t, uint(0), panel.SelectedCount())

	// first added

	panel.SetSelected("/tmp/file1.txt", true)
	assert.Equal(t, uint(1), panel.SelectedCount())
	assert.Equal(t, map[string]int{"/tmp/file1.txt": 1}, panel.selected)

	// second added
	panel.SetSelected("/tmp/file2.txt", true)
	assert.Equal(t, map[string]int{"/tmp/file1.txt": 1, "/tmp/file2.txt": 2}, panel.selected)
	assert.Equal(t, uint(2), panel.SelectedCount())
	currentFirst := panel.GetFirstSelectedLocation()
	assert.Equal(t, "/tmp/file1.txt", currentFirst)

	// first removed
	panel.SetSelected("/tmp/file1.txt", false)
	assert.Equal(t, uint(1), panel.SelectedCount())
	assert.Equal(t, map[string]int{"/tmp/file2.txt": 2}, panel.selected)
	currentFirst = panel.GetFirstSelectedLocation()
	assert.Equal(t, "/tmp/file2.txt", currentFirst)

	// multi select
	panel.SetSelectedAll([]string{"/tmp/file3.txt", "/tmp/file4.txt"}, true)
	assert.Equal(t, map[string]int{"/tmp/file2.txt": 2, "/tmp/file3.txt": 3, "/tmp/file4.txt": 4}, panel.selected)
	assert.Equal(t, uint(3), panel.SelectedCount())

	// multi unselect
	panel.SetSelectedAll([]string{"/tmp/file2.txt", "/tmp/file4.txt"}, false)
	assert.Equal(t, map[string]int{"/tmp/file3.txt": 3}, panel.selected)
	assert.Equal(t, uint(1), panel.SelectedCount())
	currentFirst = panel.GetFirstSelectedLocation()
	assert.Equal(t, "/tmp/file3.txt", currentFirst)

	// reset selection
	panel.ResetSelected()
	assert.Equal(t, uint(0), panel.SelectedCount())
	assert.Equal(t, map[string]int{}, panel.selected)
	assert.Equal(t, 0, panel.selectOrderCounter)
}
