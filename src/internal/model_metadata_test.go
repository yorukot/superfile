package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestGetMetadataCmd_RequestTracking(t *testing.T) {
	curTestDir := filepath.Join(testDir, t.Name())
	file1 := filepath.Join(curTestDir, "file1.txt")
	file2 := filepath.Join(curTestDir, "file2.txt")

	utils.SetupDirectories(t, curTestDir)
	utils.SetupFiles(t, file1, file2)
	t.Cleanup(func() {
		_ = os.RemoveAll(curTestDir)
	})

	m := defaultTestModel(curTestDir)
	m.disableMetadata = false

	idx1 := m.getFocusedFilePanel().FindElementIndexByName("file1.txt")
	require.NotEqual(t, -1, idx1)
	m.getFocusedFilePanel().SetCursorPosition(idx1)

	// First call for fresh selection returns non-nil command
	cmd1 := m.getMetadataCmd()
	assert.NotNil(t, cmd1, "getMetadataCmd() should return non-nil command on first call for fresh selection")

	// Set expected location and focus (simulating what SetMetadata / initial tracking state sets)
	m.fileMetaData.SetMetadataLocationAndFocused(file1, false)

	// Immediate repeat call returns nil due to pending short-circuit
	cmd2 := m.getMetadataCmd()
	assert.Nil(t, cmd2, "getMetadataCmd() should return nil on immediate repeat call while request is pending")

	// Move cursor to file2 using FindElementIndexByLocation
	idx2 := m.getFocusedFilePanel().FindElementIndexByLocation(file2)
	require.NotEqual(t, -1, idx2)
	m.getFocusedFilePanel().SetCursorPosition(idx2)

	// getMetadataCmd returns non-nil again for new selection
	cmd3 := m.getMetadataCmd()
	assert.NotNil(t, cmd3, "getMetadataCmd() should return non-nil command after selection moves to another file")
}

func TestMetadataMsg_LateResponseRejection(t *testing.T) {
	curTestDir := filepath.Join(testDir, t.Name())
	file1 := filepath.Join(curTestDir, "file1.txt")
	file2 := filepath.Join(curTestDir, "file2.txt")

	utils.SetupDirectories(t, curTestDir)
	utils.SetupFiles(t, file1, file2)
	t.Cleanup(func() {
		_ = os.RemoveAll(curTestDir)
	})

	m := defaultTestModel(curTestDir)
	m.disableMetadata = false

	idx1 := m.getFocusedFilePanel().FindElementIndexByName("file1.txt")
	require.NotEqual(t, -1, idx1)
	m.getFocusedFilePanel().SetCursorPosition(idx1)

	_ = m.getMetadataCmd() // issues request for file1

	// Now move user selection to file2
	idx2 := m.getFocusedFilePanel().FindElementIndexByLocation(file2)
	require.NotEqual(t, -1, idx2)
	m.getFocusedFilePanel().SetCursorPosition(idx2)

	_ = m.getMetadataCmd() // issues request for file2

	// Construct a late MetadataMsg for file1 (user navigated away)
	staleMsg := NewMetadataMsg(metadata.NewMetadata([][2]string{{"Name", "file1.txt"}}, file1, ""), false, 1)
	staleMsg.ApplyToModel(m)

	// Confirm m.fileMetaData location did not overwrite to file1
	assert.NotEqual(t, file1, m.fileMetaData.GetMetadataLocation())

	// Construct matching MetadataMsg for file2 (current selection)
	currentMsg := NewMetadataMsg(metadata.NewMetadata([][2]string{{"Name", "file2.txt"}}, file2, ""), false, 2)
	currentMsg.ApplyToModel(m)

	// Confirm m.fileMetaData location updated to file2
	assert.Equal(t, file2, m.fileMetaData.GetMetadataLocation())
}

func TestMetadataMsg_SamePathStaleReqID(t *testing.T) {
	curTestDir := filepath.Join(testDir, t.Name())
	file1 := filepath.Join(curTestDir, "file1.txt")

	utils.SetupDirectories(t, curTestDir)
	utils.SetupFiles(t, file1)
	t.Cleanup(func() {
		_ = os.RemoveAll(curTestDir)
	})

	m := defaultTestModel(curTestDir)
	m.disableMetadata = false

	idx1 := m.getFocusedFilePanel().FindElementIndexByName("file1.txt")
	require.NotEqual(t, -1, idx1)
	m.getFocusedFilePanel().SetCursorPosition(idx1)

	// Request A (reqID 1)
	_ = m.getMetadataCmd()

	// Simulate a re-fetch / refresh on same path which yields Request B (reqID 2)
	// We clear pending request so getMetadataCmd will fire again
	m.fileMetaData.ClearPendingRequest()
	_ = m.getMetadataCmd() // reqID 2

	// Create response for Request A with data "A"
	msgA := NewMetadataMsg(metadata.NewMetadata([][2]string{{"Status", "DataA"}}, file1, ""), false, 1)
	// Create response for Request B with data "B"
	msgB := NewMetadataMsg(metadata.NewMetadata([][2]string{{"Status", "DataB"}}, file1, ""), false, 2)

	// Apply B first (newer response)
	msgB.ApplyToModel(m)

	// Apply A second (older response arriving late)
	msgA.ApplyToModel(m)

	// Set dimensions so Render formats lines properly
	m.fileMetaData.SetDimensions(80, 20)

	// Inspect rendered output to verify B's data is rendered rather than A's
	rendered := m.fileMetaData.Render(false)
	assert.Contains(t, rendered, "DataB", "Newer request B data should be retained")
	assert.NotContains(t, rendered, "DataA", "Stale request A data should not overwrite response B")
}
