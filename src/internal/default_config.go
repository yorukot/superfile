package internal

import (
	zoxidelib "github.com/lazysegtree/go-zoxide"

	"github.com/yorukot/superfile/src/internal/ui/helpmenu"

	"github.com/yorukot/superfile/src/internal/ui/filemodel"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
	zoxideui "github.com/yorukot/superfile/src/internal/ui/zoxide"
)

// Generate and return model containing default configurations for interface
// Maybe we can replace slice of strings with var args - Should we ?
// TODO: Move the configuration parameters to a ModelConfig struct.
// Something like `RendererConfig` struct for `Renderer` struct in ui/renderer package
// Or even better API like varargs lambda function opts
// which can be WithFooter(), WithXYZ()
// Lots of improvements are waiting on it
//   - Allow Sending thumbnailGeneratorNeeded as false to preview.New()
//     to prevent noise in test logs. Same with imagePreviewer
func defaultModelConfig(toggleDotFile, toggleFooter, firstUse bool,
	firstPanelPaths []string, zClient *zoxidelib.Client) *model {
	return &model{
		focusPanel:      nonePanelFocus,
		processBarModel: processbar.New(),
		sidebarModel:    sidebar.New(),
		fileMetaData:    metadata.New(),
		fileModel:       filemodel.New(firstPanelPaths, toggleDotFile),
		helpMenu:        helpmenu.New(),
		promptModal:     prompt.DefaultModel(prompt.PromptMinHeight, prompt.PromptMinWidth),
		zoxideModal:     zoxideui.DefaultModel(zoxideui.ZoxideMinHeight, zoxideui.ZoxideMinWidth, zClient),
		sortModal:       sortmodel.New(),
		zClient:         zClient,
		modelQuitState:  notQuitting,
		toggleFooter:    toggleFooter,
		firstUse:        firstUse,
		hasTrash:        common.InitTrash(),
	}
}
