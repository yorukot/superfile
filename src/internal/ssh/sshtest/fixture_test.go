package sshtest

import (
	"net"
	"os"
	"path/filepath"
	"sync"
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

func TestDisconnectingWriterAtSerializesConcurrentWrites(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "destination"))
	require.NoError(t, err)
	defer file.Close()
	conn, peer := net.Pipe()
	defer peer.Close()
	countingConn := &closeCountingConn{Conn: conn}
	const (
		budget      = int64(32)
		writeSize   = 8
		writerCount = 8
	)
	writer := &disconnectingWriterAt{
		file:      file,
		conn:      countingConn,
		remaining: budget,
	}

	start := make(chan struct{})
	type result struct {
		written int
		err     error
	}
	results := make(chan result, writerCount)
	var wg sync.WaitGroup
	for i := range writerCount {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start
			written, err := writer.WriteAt([]byte("12345678"), int64(index*writeSize))
			results <- result{written: written, err: err}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)

	totalWritten := 0
	thresholdWrites := 0
	connectionLost := 0
	for result := range results {
		totalWritten += result.written
		if result.err == nil {
			continue
		}
		connectionLost++
		var status *sftp.StatusError
		require.ErrorAs(t, result.err, &status)
		assert.Equal(t, sftp.ErrSSHFxConnectionLost, status.FxCode())
		if result.written == writeSize {
			thresholdWrites++
		}
	}

	assert.LessOrEqual(t, totalWritten, int(budget))
	assert.Equal(t, int(budget), totalWritten)
	assert.Equal(t, 1, thresholdWrites)
	assert.Equal(t, writerCount-3, connectionLost)
	assert.Zero(t, writer.remaining)
	assert.True(t, writer.fired)
	assert.Equal(t, 1, countingConn.closes())

	for range 2 {
		written, err := writer.WriteAt([]byte("12345678"), 0)
		assert.Zero(t, written)
		var status *sftp.StatusError
		require.ErrorAs(t, err, &status)
		assert.Equal(t, sftp.ErrSSHFxConnectionLost, status.FxCode())
	}
	assert.Equal(t, 1, countingConn.closes())
}

type closeCountingConn struct {
	net.Conn

	mu         sync.Mutex
	closeCount int
}

func (c *closeCountingConn) Close() error {
	c.mu.Lock()
	c.closeCount++
	c.mu.Unlock()
	return c.Conn.Close()
}

func (c *closeCountingConn) closes() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closeCount
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
