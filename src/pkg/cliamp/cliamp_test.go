package cliamp

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocketPath(t *testing.T) {
	t.Setenv("CLIAMP_CONFIG_DIR", "/tmp/cliampcfg")
	assert.Equal(t, filepath.Join("/tmp/cliampcfg", "cliamp.sock"), SocketPath())

	t.Setenv("CLIAMP_CONFIG_DIR", "")
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg")
	assert.Equal(t, filepath.Join("/tmp/xdg", "cliamp", "cliamp.sock"), SocketPath())

	t.Setenv("XDG_CONFIG_HOME", "")
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(home, ".config", "cliamp", "cliamp.sock"), SocketPath())
}

// startFakeServer serves the cliamp Unix socket protocol with a caller supplied
// handler. It returns the socket path.
func startFakeServer(t *testing.T, handler func(req map[string]any) map[string]any) string {
	t.Helper()
	dir := t.TempDir()
	sock := filepath.Join(dir, "cliamp.sock")
	ln, err := net.Listen("unix", sock)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close() //nolint:errcheck
				sc := bufio.NewScanner(c)
				for sc.Scan() {
					var req map[string]any
					if err := json.Unmarshal(sc.Bytes(), &req); err != nil {
						return
					}
					data, err := json.Marshal(handler(req))
					if err != nil {
						return
					}
					data = append(data, '\n')
					if _, err := c.Write(data); err != nil {
						return
					}
				}
			}(conn)
		}
	}()
	return sock
}

func TestStateAndSpectrum(t *testing.T) {
	sock := startFakeServer(t, func(req map[string]any) map[string]any {
		switch req["method"] {
		case "state.get":
			return map[string]any{
				"version": 2, "id": req["id"], "ok": true,
				"snapshot": map[string]any{
					"state":    "playing",
					"track":    map[string]any{"title": "tone", "path": "/music/tone.wav"},
					"position": 1.5,
					"duration": 5.0,
				},
			}
		case "spectrum.get":
			return map[string]any{
				"version": 2, "id": req["id"], "ok": true,
				"result": map[string]any{
					"ok": true, "visualizer": "Bars",
					"bands": []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 1.0},
				},
			}
		}
		return map[string]any{"version": 2, "id": req["id"], "ok": false, "error": "unknown"}
	})

	c := &Client{Path: sock}
	state, err := c.State()
	require.NoError(t, err)
	assert.True(t, state.IsPlaying())
	assert.Equal(t, "tone", state.Track.Title)
	assert.Equal(t, 5.0, state.Duration)

	bands, err := c.Spectrum()
	require.NoError(t, err)
	require.Len(t, bands, 10)
	assert.Equal(t, 1.0, bands[9])
}

func TestPlaySubmitsOperationsInOrder(t *testing.T) {
	var operations []string
	sock := startFakeServer(t, func(req map[string]any) map[string]any {
		if req["method"] == "operation.submit" {
			op, _ := req["operation"].(string)
			operations = append(operations, op)
			return map[string]any{
				"version": 2, "id": req["id"], "ok": true,
				"job": map[string]any{"id": "j", "operation": op, "state": "succeeded"},
			}
		}
		return map[string]any{"version": 2, "id": req["id"], "ok": true}
	})

	c := &Client{Path: sock}
	require.NoError(t, c.Play("/music/song.mp3"))
	assert.Equal(t, []string{"queue.clear", "queue", "play"}, operations)
}

func TestSubmitPollsJobUntilTerminal(t *testing.T) {
	polls := 0
	sock := startFakeServer(t, func(req map[string]any) map[string]any {
		switch req["method"] {
		case "operation.submit":
			return map[string]any{
				"version": 2, "id": req["id"], "ok": true,
				"job": map[string]any{"id": "job1", "operation": "queue", "state": "running"},
			}
		case "job.get":
			polls++
			state := "running"
			if polls >= 2 {
				state = "succeeded"
			}
			return map[string]any{
				"version": 2, "id": req["id"], "ok": true,
				"job": map[string]any{"id": "job1", "operation": "queue", "state": state},
			}
		}
		return map[string]any{"version": 2, "id": req["id"], "ok": true}
	})

	c := &Client{Path: sock}
	require.NoError(t, c.submit("queue", map[string]any{"path": "x"}))
	assert.GreaterOrEqual(t, polls, 2)
}

func TestJobError(t *testing.T) {
	assert.NoError(t, jobError(&job{State: "succeeded"}))
	assert.Error(t, jobError(&job{State: "failed", Operation: "queue", Error: "boom"}))
}

func TestIsRunningFalseForMissingSocket(t *testing.T) {
	c := &Client{Path: filepath.Join(t.TempDir(), "nope.sock")}
	assert.False(t, c.IsRunning())
}
