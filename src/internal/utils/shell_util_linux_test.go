//go:build linux
// +build linux

package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// this test need start "not root" user

func TestSudoExecFails(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root; sudo may succeed without prompting")
	} else {
		_, _, err := ExecuteCommandInShell(nil, 2*time.Second, "/tmp", "sudo ls")
		assert.ErrorContains(t, err, "interactive mode is not allowed")
	}
}
func TestExecHappyPath(t *testing.T) {
	_, _, err := ExecuteCommandInShell(nil, 2*time.Second, "/tmp", "ls")
	assert.NoError(t, err)
}
