package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCursorOutOfBoundsAfterDirectorySwitch(t *testing.T) {
	// Create two directories with different file counts
	tempDir := t.TempDir()
	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir2")
	require.NoError(t, os.MkdirAll(dir1, 0755))
	require.NoError(t, os.MkdirAll(dir2, 0755))

	// Create 10 files in dir1
	for i := 0; i < 10; i++ {
		file := filepath.Join(dir1, string('a'+rune(i))+".txt")
		require.NoError(t, os.WriteFile(file, []byte("test"), 0644))
	}

	// Create 5 files in dir2  
	for i := 0; i < 5; i++ {
		file := filepath.Join(dir2, string('a'+rune(i))+".txt")
		require.NoError(t, os.WriteFile(file, []byte("test"), 0644))
	}

	// Start with dir1
	m := defaultTestModel(dir1)
	p := NewTestTeaProgWithEventLoop(t, m)

	// Wait for initial load
	assert.Eventually(t, func() bool {
		panel := m.getFocusedFilePanel()
		return len(panel.Element) == 10
	}, DefaultTestTimeout, DefaultTestTick, "Should load 10 files in dir1")

	// Move cursor to position 8 (near end of list)
	panel := m.getFocusedFilePanel()
	for i := 0; i < 8; i++ {
		p.Send(tea.KeyMsg{Type: tea.KeyDown})
	}

	// Verify cursor is at position 8
	assert.Eventually(t, func() bool {
		return m.getFocusedFilePanel().Cursor == 8
	}, DefaultTestTimeout, DefaultTestTick, "Cursor should be at position 8")
	t.Logf("Cursor at position %d with %d elements", panel.Cursor, len(panel.Element))

	// Navigate to dir2 (this saves cursor=8 in directoryRecords)
	err := m.updateCurrentFilePanelDir(dir2)
	require.NoError(t, err)

	// Force update to load dir2 files
	m.getFilePanelItems()

	// Verify we're in dir2 with 5 files
	assert.Eventually(t, func() bool {
		panel := m.getFocusedFilePanel()
		return panel.Location == dir2 && len(panel.Element) == 5
	}, DefaultTestTimeout, DefaultTestTick, "Should be in dir2 with 5 files")

	// Delete files from dir1 externally (simulating external changes)
	for i := 4; i < 10; i++ {
		file := filepath.Join(dir1, string('a'+rune(i))+".txt")
		os.Remove(file)
	}
	t.Log("Deleted 6 files from dir1 externally")

	// Navigate back to dir1 (this restores cursor=8 from cache)
	err = m.updateCurrentFilePanelDir(dir1)
	require.NoError(t, err)
	t.Logf("After switching back, cursor restored to: %d", m.getFocusedFilePanel().Cursor)

	// Force update to load the new file list
	m.getFilePanelItems()
	time.Sleep(50 * time.Millisecond) // Give time for update

	panel = m.getFocusedFilePanel()
	t.Logf("After update: cursor=%d, elements=%d", panel.Cursor, len(panel.Element))

	// BUG: Cursor is 8 but only 4 elements exist!
	assert.Equal(t, 8, panel.Cursor, "Cursor restored from directoryRecords cache")
	assert.Equal(t, 4, len(panel.Element), "Only 4 files after external deletion")

	// This assertion FAILS, proving the bug
	assert.Less(t, panel.Cursor, len(panel.Element),
		"BUG: Cursor(%d) >= len(Element)(%d) - No cursor validation after UpdateElementsIfNeeded!",
		panel.Cursor, len(panel.Element))

	// Attempting to press Enter here would panic
	t.Logf("🔴 Pressing Enter now would cause panic: index out of range [%d] with length %d",
		panel.Cursor, len(panel.Element))

	// This is exactly what would happen in executeOpenCommand (line 73):
	// filePath := panel.Element[panel.Cursor].Location // PANIC!
}
