package internal

import (
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// ThemeReloadMsg is emitted when the process receives the theme reload signal
// (SIGUSR1 on unix). External theme switchers send it so a running superfile
// picks up a changed theme file without a restart.
type ThemeReloadMsg struct{}

// listenForThemeReload blocks on the signal channel and turns the next signal
// into a ThemeReloadMsg. It must be re-issued after every message to keep
// listening.
func listenForThemeReload(signals <-chan os.Signal) tea.Cmd {
	if signals == nil {
		return nil
	}
	return func() tea.Msg {
		<-signals
		return ThemeReloadMsg{}
	}
}

// handleThemeReload re-reads the theme file and rebuilds every theme-derived
// style and prerendered string. On failure the current theme is kept.
func (m *model) handleThemeReload() tea.Cmd {
	if err := common.ReloadThemeFile(); err != nil {
		slog.Error("Theme reload failed, keeping the current theme", "error", err)
		return listenForThemeReload(m.themeReloadSignals)
	}

	icon.InitIcon(common.Config.Nerdfont, common.Theme.DirectoryIconColor)
	common.LoadThemeConfig()
	common.LoadPrerenderedVariables()
	slog.Info("Theme reloaded", "theme", common.Config.Theme)

	return tea.Batch(
		listenForThemeReload(m.themeReloadSignals),
		// The preview panel caches its rendered output, so force a redraw
		m.fileModel.GetFilePreviewCmd(true),
	)
}
