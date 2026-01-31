//go:build linux
// +build linux

package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSudoExecFails(t *testing.T) {
	_, _, err := ExecuteCommandInShell(nil, 2*time.Second, "/tmp", "sudo ls")
	assert.ErrorContains(t, err, "interactive mode is not allowed")
}
func TestExecHappyPath(t *testing.T) {
	_, _, err := ExecuteCommandInShell(nil, 2*time.Second, "/tmp", "ls")
	assert.NoError(t, err)
}
