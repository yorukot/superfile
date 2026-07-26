package sshtest

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisconnectingWriterAtDisconnectsAtExactThreshold(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "destination"))
	require.NoError(t, err)
	defer file.Close()
	conn, peer := net.Pipe()
	defer peer.Close()
	writer := &disconnectingWriterAt{file: file, conn: conn, remaining: 3}

	written, err := writer.WriteAt([]byte("abc"), 0)
	assert.Equal(t, 3, written)
	require.Error(t, err)
	var status *sftp.StatusError
	require.ErrorAs(t, err, &status)
	assert.Equal(t, sftp.ErrSSHFxConnectionLost, status.FxCode())
	assert.Zero(t, writer.remaining)
	assert.True(t, writer.fired)
}

func TestWaitForLogSinceIncludesDelayedMarker(t *testing.T) {
	fixture := Start(t)
	logInfo, err := os.Stat(fixture.LogPath)
	require.NoError(t, err)
	offset := logInfo.Size()
	marker := "fixture-wait-marker-" + t.Name()

	go func() {
		time.Sleep(10 * time.Millisecond)
		fixture.logf("marker=%s", marker)
	}()

	logText := fixture.WaitForLogSince(t, offset, marker)
	assert.Contains(t, logText, marker)
}
