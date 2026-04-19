package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func TestRenderElementNameForSaveTargetUsesDownloadIcon(t *testing.T) {
	previousNerdfont := common.Config.Nerdfont
	common.Config.Nerdfont = true
	t.Cleanup(func() {
		common.Config.Nerdfont = previousNerdfont
	})

	panel := testModel(0, 0, 12, BrowserMode, nil)

	rendered := panel.renderElementName(Element{
		Name:       "download.txt",
		Location:   "/tmp/download.txt",
		SaveTarget: true,
	}, 40, false)

	assert.Contains(t, rendered, icon.Download)
	assert.Contains(t, rendered, "download.txt")
	assert.NotContains(t, rendered, "save download.txt")
}

func TestGetPanelModeInfoForSaveModeUsesDownloadIcon(t *testing.T) {
	previousNerdfont := common.Config.Nerdfont
	common.Config.Nerdfont = true
	t.Cleanup(func() {
		common.Config.Nerdfont = previousNerdfont
	})

	panel := testModel(0, 0, 12, BrowserMode, nil)
	panel.SaveMode = true

	label, iconValue := panel.getPanelModeInfo(0)
	assert.Equal(t, "", label)
	assert.Equal(t, icon.Download, iconValue)

	panel.PanelMode = SelectMode
	label, iconValue = panel.getPanelModeInfo(2)
	assert.Equal(t, "(2)", label)
	assert.Equal(t, icon.Download, iconValue)
}
