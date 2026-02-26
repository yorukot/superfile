package internal

// TODO add two new tests for sidebar, a - with only one section, and b - without any sections.
// note - we should update `testWithConfig` to take a new object of `common.ConfigType`, so that any custom config can be provided.

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

const ScrollDownCount = 10
const ScrollUpCount = 5

func TestLayout(t *testing.T) {
	// This runs 800+ tests can be skipped via go test ./... -short
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// Uncomment this to debug locally.
	// This is to prevent too many logs in CICD
	utils.SetRootLoggerToDiscarded()
	t.Cleanup(func() {
		if testing.Verbose() {
			utils.SetRootLoggerToStdout(true)
		}
	})

	baseTestDir := t.TempDir()
	subDir := filepath.Join(baseTestDir, "subdir")
	subDir2 := filepath.Join(baseTestDir, "subdir2")
	utils.SetupDirectories(t, baseTestDir, subDir, subDir2)
	utils.SetupFiles(t,
		filepath.Join(baseTestDir, "file1.txt"),
		filepath.Join(baseTestDir, "file2.txt"),
		filepath.Join(baseTestDir, "file3.txt"),
		filepath.Join(subDir, "nested1.txt"),
		filepath.Join(subDir, "nested2.txt"),
	)

	sidebarWidths := []int{0, 5, 12, 20}
	previewWidths := []int{0, 2, 3, 10}
	pWDef := common.Config.FilePreviewWidth
	sWDef := common.Config.SidebarWidth

	// Note: These cannot run in parallel for now as they share the same
	// config global variable. Later we might fix that and use parallelization
	for _, w := range sidebarWidths {
		t.Run(fmt.Sprintf("sW=%d;pW=%d", w, pWDef), func(t *testing.T) {
			testWithConfig(t, w, pWDef, false, baseTestDir)
		})
	}
	for _, w := range previewWidths {
		t.Run(fmt.Sprintf("sW=%d;pW=%d", sWDef, w), func(t *testing.T) {
			testWithConfig(t, sWDef, w, false, baseTestDir)
		})
	}

	// One test for preview border enabled
	t.Run("sW=10;pW=5;previewWithBorder", func(t *testing.T) {
		testWithConfig(t, 10, 5, true, baseTestDir)
	})
}

func testWithConfig(t *testing.T, sidebarWidth int, previewWidth int,
	previewBorderEnabled bool, testPath string) {
	// Save original config values and restore them after test
	origSidebarWidth := common.Config.SidebarWidth
	origPreviewWidth := common.Config.FilePreviewWidth
	origPreviewBorderEnabled := common.Config.EnableFilePreviewBorder
	defer func() {
		common.Config.SidebarWidth = origSidebarWidth
		common.Config.FilePreviewWidth = origPreviewWidth
		common.Config.EnableFilePreviewBorder = origPreviewBorderEnabled
	}()

	// Set test config
	common.Config.SidebarWidth = sidebarWidth
	common.Config.FilePreviewWidth = previewWidth
	common.Config.EnableFilePreviewBorder = previewBorderEnabled

	m := defaultTestModelWithFooterAndFilePreview(testPath)
	p := NewTestTeaProgWithEventLoop(t, m)

	resizeSizes := []struct {
		width, height int
	}{
		{60, 24},   // Minimum
		{80, 30},   // Small HeightBreakC (<35)
		{100, 39},  // HeightBreakC
		{130, 44},  // HeightBreakD
		{200, 60},  // Large
		{400, 120}, // Extra large
		{91, 41},   // Back to medium
		{60, 24},   // Back to minimum
	}

	// Run resize tests
	for _, size := range resizeSizes {
		t.Run(fmt.Sprintf("w=%d;h=%d", size.width, size.height), func(t *testing.T) {
			updateModelDimensionsAndValidate(t, p, size.width, size.height)
		})
	}

	t.Run("Edge cases", func(t *testing.T) {
		edgeCases := []struct {
			name          string
			width, height int
		}{
			{"Ultra-narrow", 70, 100},
			{"Ultra-wide", 500, 30},
			{"Boundary-79", 79, 30},
			{"Boundary-80", 80, 30},
			{"Boundary-81", 81, 30},
			{"Below-minimum", 59, 23},
		}

		for _, tc := range edgeCases {
			t.Run(tc.name, func(t *testing.T) {
				p.SendDirectly(tea.WindowSizeMsg{Width: tc.width, Height: tc.height})
				assertLayoutValidity(t, p.m)
			})
		}
	})
}

// Note: this will create as many panels possible and leave the model in that state
// This is to ensure that at time of resize operations, there are more panels
func updateModelDimensionsAndValidate(t *testing.T, p *TeaProg, width int, height int) {
	// Set Footer OFF, Preview OFF via model state changes
	// if p.m.toggleFooter {
	//	p.SendDirectly(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(common.Hotkeys.ToggleFooter[0])[0:1]})
	//}
	// File preview toggle - just send the key, no need to check state
	// Sending toggle key will turn it off if it's on

	require.True(t, p.m.toggleFooter)
	require.True(t, p.m.fileModel.FilePreview.IsOpen())

	testdata := []struct {
		name string
		msg  []tea.Msg
	}{
		{
			name: "Resize",
			msg:  []tea.Msg{tea.WindowSizeMsg{Width: width, Height: height}},
		},
		{
			name: "FooterOffPreviewOff",
			msg: []tea.Msg{
				utils.TeaRuneKeyMsg(common.Hotkeys.ToggleFooter[0]),
				utils.TeaRuneKeyMsg(common.Hotkeys.ToggleFilePreviewPanel[0]),
			},
		},
		{
			name: "ToggleFooterOn",
			msg:  []tea.Msg{utils.TeaRuneKeyMsg(common.Hotkeys.ToggleFooter[0])},
		},
		{
			name: "FooterOffPreviewOn",
			msg: []tea.Msg{
				utils.TeaRuneKeyMsg(common.Hotkeys.ToggleFooter[0]),
				utils.TeaRuneKeyMsg(common.Hotkeys.ToggleFilePreviewPanel[0]),
			},
		},
		{
			name: "FooterOnAgain",
			msg:  []tea.Msg{utils.TeaRuneKeyMsg(common.Hotkeys.ToggleFooter[0])},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.msg {
				p.SendDirectly(msg)
			}
			assertLayoutValidity(t, p.m)
		})
	}

	t.Run("FilePanelRemoval", func(t *testing.T) {
		for {
			initialCount := p.m.fileModel.PanelCount()

			p.SendDirectly(utils.TeaRuneKeyMsg(common.Hotkeys.CloseFilePanel[0]))
			assertLayoutValidity(t, p.m)

			if p.m.fileModel.PanelCount() == initialCount {
				break // No panel was removed
			}
			require.Positive(t, p.m.fileModel.PanelCount())
		}
	})

	testModelScrolling(t, p)

	t.Run("FilePanelCreation", func(t *testing.T) {
		for {
			initialCount := p.m.fileModel.PanelCount()
			p.SendDirectly(utils.TeaRuneKeyMsg(common.Hotkeys.CreateNewFilePanel[0]))

			assertLayoutValidity(t, p.m)

			if p.m.fileModel.PanelCount() == initialCount {
				break // No new panel created
			}

			require.LessOrEqual(t, p.m.fileModel.PanelCount(), common.FilePanelMax,
				"Panel count should not exceed maximum")
		}
	})

	assert.Equal(t, p.m.fileModel.MaxFilePanel, p.m.fileModel.PanelCount())
}

func testModelScrolling(t *testing.T, p *TeaProg) {
	// We are at Filepanel now
	testModelScrollingCore(t, p)

	// Add dummy data to ProcessBar and Metadata
	for i := range 10 {
		p.m.processBarModel.AddProcess(
			processbar.NewProcess(strconv.Itoa(i), "test", processbar.OpCopy, 1),
		)
	}
	dummyData := [][2]string{
		{"a", "b"},
		{"a", "b"},
		{"a", "b"},
		{"a", "b"},
		{"a", "b"},
	}
	p.m.fileMetaData.SetMetadata(metadata.NewMetadata(dummyData, "", ""), true)

	panels := []struct {
		name     string
		focusKey string
	}{
		{"Sidebar", common.Hotkeys.FocusOnSidebar[0]},
		{"ProcessBar", common.Hotkeys.FocusOnProcessBar[0]},
		{"Metadata", common.Hotkeys.FocusOnMetaData[0]},
	}

	for _, panel := range panels {
		t.Run(panel.name+"Scrolling", func(t *testing.T) {
			p.SendKeyDirectly(panel.focusKey)
			// TODO: Add validation that we are actually at sidebar
			testModelScrollingCore(t, p)
		})
	}
}

func testModelScrollingCore(t *testing.T, p *TeaProg) {
	for range ScrollDownCount {
		p.SendDirectly(tea.KeyMsg{Type: tea.KeyDown})
	}
	assertLayoutValidity(t, p.m)

	// Scroll up
	for range ScrollUpCount {
		p.SendDirectly(tea.KeyMsg{Type: tea.KeyUp})
	}
	assertLayoutValidity(t, p.m)
}

func assertLayoutValidity(t *testing.T, m *model) {
	// Skip for edge conditions where terminal is too small
	if m.fullHeight < common.MinimumHeight || m.fullWidth < common.MinimumWidth {
		return // Terminal too small for valid layout
	}
	if m.fileModel.SinglePanelWidth < filepanel.MinWidth {
		return // Panels too narrow for valid layout
	}

	returnFirstError := func() error {
		if err := m.validateLayout(); err != nil {
			return err
		}
		if err := m.validateComponentRender(); err != nil {
			return err
		}
		if err := m.validateFinalRender(); err != nil {
			return err
		}
		return nil
	}
	err := returnFirstError()
	// Not using assert to prevent `getLayoutInfoForDebug` getting called
	// in happy case. This is hot-path for 906 tests
	if err != nil {
		t.Errorf("validations failed, error : %v, layout info : %v",
			err, getLayoutInfoForDebug(m))
	}
}

func getLayoutInfoForDebug(m *model) string {
	firstPanel := m.fileModel.FilePanels[0]
	lastPanel := m.fileModel.FilePanels[m.fileModel.PanelCount()-1]
	location := m.getFocusedFilePanel().Location
	width := fmt.Sprintf("width=%d[sidebar=%d,filemodel=%d"+
		"[firstpanel=%d,lastpanel=%d,previewExp=%d,previewActual=%d]]"+
		"[panelCount=%d,maxPanel=%d]"+
		"[processbarWidth=%d,clipboardWidth=%d]",
		m.fullWidth, common.Config.SidebarWidth, m.fileModel.Width,
		firstPanel.GetWidth(), lastPanel.GetWidth(), m.fileModel.ExpectedPreviewWidth,
		m.fileModel.FilePreview.GetContentWidth(),
		m.fileModel.PanelCount(), m.fileModel.MaxFilePanel,
		m.processBarModel.GetWidth(), m.clipboard.GetWidth())

	height := fmt.Sprintf("height=%d[fileModel=%d[firstPanel=%d,previewActual=%d],footer=%d]",
		m.fullHeight, m.fileModel.Height, firstPanel.GetHeight(),
		m.fileModel.FilePreview.GetContentHeight(), m.footerHeight)

	return fmt.Sprintf("%s %s location=%s", width, height, location)
}
