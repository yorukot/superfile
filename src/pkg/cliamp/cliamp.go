// Package cliamp is a small client for the cliamp music player's headless
// daemon. cliamp exposes a newline-delimited JSON protocol over a Unix socket
// (see the "cliamp remote" documentation). This package speaks the V2 envelope
// so superfile can query now-playing state, poll the spectrum visualizer, and
// tell cliamp to play a local file.
package cliamp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/yorukot/superfile/src/pkg/utils"
)

const (
	// protocolVersion is the required V2 envelope version.
	protocolVersion = 2

	defaultTimeout = 700 * time.Millisecond
	jobPollEvery   = 10 * time.Millisecond
	jobPollTimeout = 2 * time.Second
)

// ErrNotRunning is returned when the cliamp daemon socket cannot be reached.
var ErrNotRunning = errors.New("cliamp daemon is not running")

// SocketPath resolves the path of the cliamp daemon's Unix socket using the
// same config-directory rules as the cliamp binary itself.
func SocketPath() string {
	if dir := os.Getenv("CLIAMP_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "cliamp.sock")
	}
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "cliamp", "cliamp.sock")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "cliamp", "cliamp.sock")
}

// Track is a single track as reported by cliamp. Only Title and Path are
// guaranteed; artist/album are present for tagged files.
type Track struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Path   string `json:"path"`
	Stream bool   `json:"stream"`
}

// Snapshot is the now-playing state returned by the state.get method.
type Snapshot struct {
	Revision         int     `json:"revision"`
	PlaylistRevision int     `json:"playlist_revision"`
	State            string  `json:"state"`
	Track            Track   `json:"track"`
	LogicalTrack     Track   `json:"logical_track"`
	Position         float64 `json:"position"`
	Duration         float64 `json:"duration"`
	Seekable         bool    `json:"seekable"`
	Total            int     `json:"total"`
	Shuffle          bool    `json:"shuffle"`
	Repeat           string  `json:"repeat"`
	Mono             bool    `json:"mono"`
	Speed            float64 `json:"speed"`
	Visualizer       string  `json:"visualizer"`
}

// IsPlaying reports whether cliamp is currently outputting audio.
func (s Snapshot) IsPlaying() bool {
	return s.State == "playing"
}

type trackResult struct {
	Total  int     `json:"total"`
	Tracks []Track `json:"tracks"`
}

type spectrumResult struct {
	OK         bool      `json:"ok"`
	Visualizer string    `json:"visualizer"`
	Bands      []float64 `json:"bands"`
}

type job struct {
	ID        string          `json:"id"`
	Operation string          `json:"operation"`
	State     string          `json:"state"`
	Error     string          `json:"error"`
	Result    json.RawMessage `json:"result"`
	Snapshot  *Snapshot       `json:"snapshot"`
}

func (j *job) terminal() bool {
	return j.State == "succeeded" || j.State == "failed" || j.State == "canceled"
}

type request struct {
	Version   int            `json:"version"`
	ID        string         `json:"id"`
	Method    string         `json:"method"`
	Operation string         `json:"operation,omitempty"`
	Params    map[string]any `json:"params,omitempty"`
	JobID     string         `json:"job_id,omitempty"`
}

type response struct {
	Version  int             `json:"version"`
	ID       string          `json:"id"`
	OK       bool            `json:"ok"`
	Code     string          `json:"code"`
	Error    string          `json:"error"`
	Snapshot *Snapshot       `json:"snapshot"`
	Result   json.RawMessage `json:"result"`
	Job      *job            `json:"job"`
}

// Client talks to a cliamp daemon over its Unix socket.
type Client struct {
	Path    string
	Timeout time.Duration
}

// New returns a Client using the default socket path and timeout.
func New() *Client {
	return &Client{Path: SocketPath(), Timeout: defaultTimeout}
}

func (c *Client) timeout() time.Duration {
	if c.Timeout <= 0 {
		return defaultTimeout
	}
	return c.Timeout
}

// do sends a single request and reads a single response line.
func (c *Client) do(req request) (*response, error) {
	if c.Path == "" {
		return nil, ErrNotRunning
	}
	dialer := net.Dialer{Timeout: c.timeout()}
	conn, err := dialer.Dial("unix", c.Path)
	if err != nil {
		return nil, err
	}
	defer conn.Close() //nolint:errcheck
	_ = conn.SetDeadline(time.Now().Add(c.timeout()))

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	payload = append(payload, '\n')
	if _, err := conn.Write(payload); err != nil {
		return nil, err
	}

	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil && len(line) == 0 {
		return nil, err
	}
	var resp response
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("cliamp: invalid response: %w", err)
	}
	if !resp.OK {
		if resp.Error != "" {
			return &resp, fmt.Errorf("cliamp: %s", resp.Error)
		}
		return &resp, errors.New("cliamp: request failed")
	}
	return &resp, nil
}

// State returns the current now-playing snapshot.
func (c *Client) State() (Snapshot, error) {
	resp, err := c.do(request{Version: protocolVersion, ID: "state", Method: "state.get"})
	if err != nil {
		return Snapshot{}, err
	}
	if resp.Snapshot == nil {
		return Snapshot{}, errors.New("cliamp: state response had no snapshot")
	}
	return *resp.Snapshot, nil
}

// Spectrum returns the current visualizer bands (each in the range 0..1).
func (c *Client) Spectrum() ([]float64, error) {
	resp, err := c.do(request{Version: protocolVersion, ID: "spectrum", Method: "spectrum.get"})
	if err != nil {
		return nil, err
	}
	var res spectrumResult
	if err := json.Unmarshal(resp.Result, &res); err != nil {
		return nil, fmt.Errorf("cliamp: invalid spectrum result: %w", err)
	}
	return res.Bands, nil
}

// IsRunning reports whether a cliamp daemon is reachable.
func (c *Client) IsRunning() bool {
	_, err := c.State()
	return err == nil
}

// submit sends an operation and waits until its job reaches a terminal state.
// A nil return means the operation was accepted (or completed successfully).
func (c *Client) submit(operation string, params map[string]any) error {
	resp, err := c.do(request{
		Version:   protocolVersion,
		ID:        operation,
		Method:    "operation.submit",
		Operation: operation,
		Params:    params,
	})
	if err != nil {
		return err
	}
	if resp.Job == nil || resp.Job.terminal() {
		return jobError(resp.Job)
	}

	deadline := time.Now().Add(jobPollTimeout)
	for time.Now().Before(deadline) {
		time.Sleep(jobPollEvery)
		got, err := c.do(request{Version: protocolVersion, ID: "job", Method: "job.get", JobID: resp.Job.ID})
		if err != nil {
			return err
		}
		if got.Job == nil {
			return nil
		}
		if got.Job.terminal() {
			return jobError(got.Job)
		}
	}
	return nil
}

// fire submits an operation without waiting for its job to finish. Operations
// are processed in submission order, so a later waited-on operation is still
// guaranteed to run after an earlier fired one.
func (c *Client) fire(operation string, params map[string]any) error {
	_, err := c.do(request{
		Version:   protocolVersion,
		ID:        operation,
		Method:    "operation.submit",
		Operation: operation,
		Params:    params,
	})
	return err
}

func jobError(j *job) error {
	if j == nil {
		return nil
	}
	if j.State == "failed" || j.State == "canceled" {
		msg := j.Error
		if msg == "" {
			msg = j.State
		}
		return fmt.Errorf("cliamp: %s: %s", j.Operation, msg)
	}
	return nil
}

// ApplyQueue replaces the current playlist with a single path and plays it.
// This is the sequence superfile uses when a file is opened.
func (c *Client) ApplyQueue(path string) error {
	// Best effort: clearing an already-empty queue should not abort playback,
	// and we don't wait for it because the queued operations run in order.
	_ = c.fire("queue.clear", map[string]any{})
	if err := c.submit("queue", map[string]any{"path": path}); err != nil {
		return err
	}
	return c.submit("play", map[string]any{})
}

// Play is an alias for ApplyQueue to match cliamp's "play a file" intent.
func (c *Client) Play(path string) error {
	return c.ApplyQueue(path)
}

// EnsureDaemon starts a detached cliamp daemon if none is running, then waits
// up to a few seconds for its socket to become reachable.
func (c *Client) EnsureDaemon() error {
	if c.IsRunning() {
		return nil
	}
	bin, err := exec.LookPath("cliamp")
	if err != nil {
		return fmt.Errorf("cliamp not found in PATH: %w", err)
	}
	cmd := exec.Command(bin, "-d")
	utils.DetachFromTerminal(cmd)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cliamp daemon: %w", err)
	}
	_ = cmd.Process.Release()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if c.IsRunning() {
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return errors.New("timed out waiting for cliamp daemon")
}
