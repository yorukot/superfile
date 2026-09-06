package preview

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/filesystem"
)

func TestRemoteUpdateCarriesRawTransmit(t *testing.T) {
	msg := NewRemoteUpdateMsg("/remote", "content", "kitty-clear", 10, 20, 1, "session", 2, nil)
	assert.Equal(t, "kitty-clear", msg.GetRawTransmit())
	assert.Equal(t, filesystem.SessionID("session"), msg.GetSessionID())
}
