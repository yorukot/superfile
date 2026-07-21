package internal

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func TestProviderPasteFailurePreservesUnprocessedCutSources(t *testing.T) {
	m := defaultTestModel(t.TempDir())
	sources := []filesystem.Location{
		providerPasteTestLocation("one", "/one"),
		providerPasteTestLocation("two", "/two"),
		providerPasteTestLocation("three", "/three"),
	}
	m.clipboard.Reset(true)
	m.clipboard.SetLocations(sources)

	msg := NewProviderPasteOperationMsg(
		processbar.Failed,
		sources[1],
		sources,
		sources[1:],
		errors.New("transfer failed"),
		1,
	)
	_ = msg.ApplyToModel(m)

	assert.True(t, m.clipboard.IsCut())
	assert.Equal(t, sources[1:], m.clipboard.GetLocations())
}

func TestProviderPasteSuccessfulCutClearsClipboard(t *testing.T) {
	m := defaultTestModel(t.TempDir())
	source := providerPasteTestLocation("one", "/one")
	m.clipboard.Reset(true)
	m.clipboard.SetLocations([]filesystem.Location{source})

	msg := NewProviderPasteOperationMsg(
		processbar.Successful,
		filesystem.Location{},
		[]filesystem.Location{source},
		nil,
		nil,
		1,
	)
	_ = msg.ApplyToModel(m)

	assert.False(t, m.clipboard.IsCut())
	assert.Empty(t, m.clipboard.GetLocations())
}

func TestProviderPasteRefreshesLocalPanels(t *testing.T) {
	dir := t.TempDir()
	m := defaultTestModel(dir)
	const fileName = "pasted.txt"
	require.Equal(t, -1, m.fileModel.FilePanels[0].FindElementIndexByName(fileName))
	require.NoError(t, os.WriteFile(filepath.Join(dir, fileName), []byte("pasted"), 0o600))

	msg := NewProviderPasteOperationMsg(
		processbar.Successful,
		filesystem.Location{},
		[]filesystem.Location{{Provider: filesystem.ProviderLocal, Path: filesystem.NewLocalPath(dir)}},
		nil,
		nil,
		1,
	)
	cmd := msg.ApplyToModel(m)

	assert.Nil(t, cmd)
	assert.NotEqual(t, -1, m.fileModel.FilePanels[0].FindElementIndexByName(fileName))
}

func TestCreateSubmissionGuardClearsOnCompletion(t *testing.T) {
	m := defaultTestModel(t.TempDir())
	m.panelCreateNewFile()
	m.typingModal.textInput.SetValue("created.txt")

	cmd := m.getCreateCmd()
	require.NotNil(t, cmd)
	assert.True(t, m.typingModal.submitting)
	m.typingModal.open = true
	assert.Nil(t, m.getCreateCmd())

	msg, ok := cmd().(CreateOperationMsg)
	require.True(t, ok)
	_ = msg.ApplyToModel(m)
	assert.False(t, m.typingModal.submitting)
}

func providerPasteTestLocation(sessionID string, path string) filesystem.Location {
	return filesystem.Location{
		Provider:  filesystem.ProviderSFTP,
		SessionID: filesystem.SessionID(sessionID),
		Label:     sessionID,
		Path:      filesystem.NewRemotePath(path),
	}
}
