package helpmenu

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func TestNavigationKeepsCursorVisible(t *testing.T) {
	require.NoError(t, common.PopulateGlobalConfigs())

	for _, tt := range []struct {
		name   string
		height int
		up     bool
	}{
		{name: "down in short terminal", height: 22},
		{name: "up in short terminal", height: 22, up: true},
		{name: "down in tall terminal", height: 30},
		{name: "up in tall terminal", height: 30, up: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.SetDimensions(90, tt.height)
			m.Open()
			key := common.Hotkeys.ListDown[0]
			if tt.up {
				key = common.Hotkeys.ListUp[0]
			}

			for range len(m.filteredData) + 1 {
				m.HandleKey(key)
				require.Empty(t, m.filteredData[m.cursor].subTitle)
				require.Contains(t, ansi.Strip(m.Render()), icon.Cursor,
					"cursor %d is hidden with render index %d", m.cursor, m.renderIndex)
			}
		})
	}
}
