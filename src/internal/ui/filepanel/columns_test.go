package filepanel

import (
	"io/fs"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

type FakeFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	sys     interface{}
}

func (f FakeFileInfo) Name() string       { return f.name }
func (f FakeFileInfo) Size() int64        { return f.size }
func (f FakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f FakeFileInfo) ModTime() time.Time { return f.modTime }
func (f FakeFileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f FakeFileInfo) Sys() interface{}   { return f.sys }

func NewFakeFileInfo(name string, size int64, mode fs.FileMode, modTime time.Time) fs.FileInfo {
	return FakeFileInfo{
		name:    name,
		size:    size,
		mode:    mode,
		modTime: modTime,
		sys:     nil,
	}
}

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

func TestRenderFileName(t *testing.T) {
	is := assert.New(t)
	t.Run("Regular file without cursor or selection", func(t *testing.T) {
		expectedWigth := 30
		panel := testModel(0, 0, 12, BrowserMode, []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt", Info: NewFakeFileInfo("file1.txt", 1024, 0644, time.Now())},
			{Name: "file2.txt", Location: "/tmp/file2.txt", Info: NewFakeFileInfo("file2.txt", 1024, 0644, time.Now())},
			{Name: "file3verylong-long-filename-for-testing-long-file-name.txt",
				Location: "/tmp/file3verylong-long-filename-for-testing-long-file-name.txt",
				Info: NewFakeFileInfo("ffile3verylong-long-filename-for-testing-long-file-name.txt",
					1024, 0644, time.Now())},
		})
		renderedStr := panel.renderFileName(1, expectedWigth)
		is.Equal(expectedWigth, ansi.StringWidth(renderedStr))
		is.Equal("   file2.txt                  ", renderedStr)
		renderedStr = panel.renderFileName(2, expectedWigth)
		is.Equal(expectedWigth, ansi.StringWidth(renderedStr))
		is.Equal("   file3verylong-long-filen...", renderedStr)
	})
	t.Run("File with cursor", func(t *testing.T) {
		expectedWigth := 32
		panel := testModel(1, 0, 12, BrowserMode, []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt", Info: NewFakeFileInfo("file1.txt", 1024, 0644, time.Now())},
			{Name: "file2.txt", Location: "/tmp/file2.txt", Info: NewFakeFileInfo("file2.txt", 1024, 0644, time.Now())},
			{Name: "file3.txt", Location: "/tmp/file3.txt", Info: NewFakeFileInfo("file3.txt", 1024, 0644, time.Now())},
		})
		renderedStr := panel.renderFileName(1, expectedWigth)
		is.Equal(expectedWigth, ansi.StringWidth(renderedStr))
		is.Equal("\uf054  file2.txt                    ", renderedStr)
	})
	t.Run("Selected file", func(t *testing.T) {
		expected := "  F\uf15c file3.txt                "
		panel := testModel(1, 0, 12, SelectMode, []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt", Info: NewFakeFileInfo("file1.txt", 1024, 0644, time.Now())},
			{Name: "file2.txt", Location: "/tmp/file2.txt", Info: NewFakeFileInfo("file2.txt", 1024, 0644, time.Now())},
			{Name: "file3.txt", Location: "/tmp/file3.txt", Info: NewFakeFileInfo("file3.txt", 1024, 0644, time.Now())},
		})
		origShowSelectIcons := common.Config.ShowSelectIcons
		origNerdfont := common.Config.Nerdfont
		origCheckbox := common.CheckboxChecked
		common.Config.ShowSelectIcons = true
		common.Config.Nerdfont = true
		//nolint:reassign // made for test
		common.CheckboxChecked = "F"
		defer func() {
			common.Config.ShowSelectIcons = origShowSelectIcons
			common.Config.Nerdfont = origNerdfont
			//nolint:reassign // rolled back test value
			common.CheckboxChecked = origCheckbox
		}()

		panel.selected = map[string]int{"/tmp/file3.txt": 1}
		renderedStr := panel.renderFileName(2, 30)
		is.Equal(30, ansi.StringWidth(renderedStr))
		is.Equal(expected, renderedStr)
	})
}

func TestRenderFileSize(t *testing.T) {
	is := assert.New(t)
	t.Run("Regular file without cursor or selection", func(t *testing.T) {
		panel := testModel(0, 0, 12, BrowserMode, []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt",
				Info: NewFakeFileInfo("file1.txt", 1024, 0644, time.Now())},
			{Name: "file2.txt", Location: "/tmp/file2.txt",
				Info: NewFakeFileInfo("file2.txt", 102400000, 0644, time.Now())},
		})
		renderedStr := panel.renderFileSize(0, FileSizeColumnWidth)
		is.Equal(FileSizeColumnWidth, ansi.StringWidth(renderedStr))
		renderedStr = panel.renderFileSize(1, FileSizeColumnWidth)
		is.Equal(FileSizeColumnWidth, ansi.StringWidth(renderedStr))
	})
}

func TestRenderModifyTime(t *testing.T) {
	is := assert.New(t)
	t.Run("Regular file without cursor or selection", func(t *testing.T) {
		panel := testModel(0, 0, 12, BrowserMode, []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt", Info: NewFakeFileInfo("file1.txt", 1024, 0644, time.Now())},
		})
		renderedStr := panel.renderModifyTime(0, ModifyTimeSizeColumnWidth)
		is.Equal(ModifyTimeSizeColumnWidth, ansi.StringWidth(renderedStr))
	})
}

func TestRenderPermissions(t *testing.T) {
	is := assert.New(t)
	t.Run("Regular file without cursor or selection", func(t *testing.T) {
		panel := testModel(0, 0, 12, BrowserMode, []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt", Info: NewFakeFileInfo("file1.txt", 1024, 0644, time.Now())},
		})
		renderedStr := panel.renderPermissions(0, PermissionsColumnWidth)
		is.Equal(PermissionsColumnWidth, ansi.StringWidth(renderedStr))
	})
}
