package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

func TestRemoteStatusText(t *testing.T) {
	filepanelTestConfigOnce.Do(func() {
		require.NoError(t, common.PopulateGlobalConfigs())
	})
	tests := []struct {
		name         string
		remote       bool
		status       string
		expected     string
		expectedSide string
	}{
		{
			name:         "local",
			expected:     "",
			expectedSide: "",
		},
		{
			name:         "remote defaults to connected",
			remote:       true,
			expected:     "remote:/tmp/remote connected",
			expectedSide: "remote connected",
		},
		{
			name:         "remote disconnected",
			remote:       true,
			status:       "  disconnected  ",
			expected:     "remote:/tmp/remote disconnected",
			expectedSide: "remote disconnected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panel := New("/tmp/local", false, "", sortmodel.SortByName, false)
			if tt.remote {
				panel.SetPaneLocation(filesystem.Location{
					Provider:  filesystem.ProviderSFTP,
					SessionID: "remote",
					Path:      filesystem.NewRemotePath("/tmp/remote"),
					Label:     "remote",
				})
			}
			panel.SetPaneConnectionStatus(tt.status)

			assert.Equal(t, tt.expected, panel.RemoteStatusText())
			assert.Equal(t, tt.expectedSide, panel.RemoteSidebarStatusText())
		})
	}
}
