package metadata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestModel_IsPending(t *testing.T) {
	tests := []struct {
		name            string
		setPath         string
		setFocus        bool
		setReqID        int
		checkPath       string
		checkFocus      bool
		expectedPending bool
	}{
		{
			name:            "matches on exact path and focus",
			setPath:         "/path/to/file1",
			setFocus:        true,
			setReqID:        1,
			checkPath:       "/path/to/file1",
			checkFocus:      true,
			expectedPending: true,
		},
		{
			name:            "false on path mismatch",
			setPath:         "/path/to/file1",
			setFocus:        true,
			setReqID:        1,
			checkPath:       "/path/to/file2",
			checkFocus:      true,
			expectedPending: false,
		},
		{
			name:            "false on focus mismatch",
			setPath:         "/path/to/file1",
			setFocus:        true,
			setReqID:        1,
			checkPath:       "/path/to/file1",
			checkFocus:      false,
			expectedPending: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.SetPendingRequest(tt.setPath, tt.setFocus, tt.setReqID)
			assert.Equal(t, tt.expectedPending, m.IsPending(tt.checkPath, tt.checkFocus))
		})
	}
}

func TestModel_MatchPendingRequest(t *testing.T) {
	tests := []struct {
		name          string
		setPath       string
		setFocus      bool
		setReqID      int
		checkPath     string
		checkFocus    bool
		checkReqID    int
		expectedMatch bool
	}{
		{
			name:          "matches on exact path, focus, and reqID",
			setPath:       "/path/to/file1",
			setFocus:      true,
			setReqID:      100,
			checkPath:     "/path/to/file1",
			checkFocus:    true,
			checkReqID:    100,
			expectedMatch: true,
		},
		{
			name:          "false when reqID differs but path/focus match (stale response for superseded request)",
			setPath:       "/path/to/file1",
			setFocus:      true,
			setReqID:      101,
			checkPath:     "/path/to/file1",
			checkFocus:    true,
			checkReqID:    100,
			expectedMatch: false,
		},
		{
			name:          "false when path differs",
			setPath:       "/path/to/file1",
			setFocus:      true,
			setReqID:      100,
			checkPath:     "/path/to/file2",
			checkFocus:    true,
			checkReqID:    100,
			expectedMatch: false,
		},
		{
			name:          "false when focus differs",
			setPath:       "/path/to/file1",
			setFocus:      true,
			setReqID:      100,
			checkPath:     "/path/to/file1",
			checkFocus:    false,
			checkReqID:    100,
			expectedMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.SetPendingRequest(tt.setPath, tt.setFocus, tt.setReqID)
			assert.Equal(t, tt.expectedMatch, m.MatchPendingRequest(tt.checkPath, tt.checkFocus, tt.checkReqID))
		})
	}
}

func TestModel_ClearPendingRequest(t *testing.T) {
	m := New()
	m.SetPendingRequest("/path/to/file", true, 42)
	assert.True(t, m.IsPending("/path/to/file", true))

	m.ClearPendingRequest()
	assert.False(t, m.IsPending("/path/to/file", true))
	assert.False(t, m.MatchPendingRequest("/path/to/file", true, 42))
}

func TestModel_IsFresh(t *testing.T) {
	ttl := 5 * time.Second

	t.Run("zero-value lastUpdated is never fresh", func(t *testing.T) {
		m := New()
		assert.False(t, m.IsFresh(ttl))
	})

	t.Run("freshly-set is fresh", func(t *testing.T) {
		m := New()
		m.lastUpdated = time.Now()
		assert.True(t, m.IsFresh(ttl))
	})

	t.Run("set past the TTL is stale", func(t *testing.T) {
		m := New()
		m.lastUpdated = time.Now().Add(-10 * time.Second)
		assert.False(t, m.IsFresh(ttl))
	})
}

func TestModel_SetMetadataLocationAndFocused(t *testing.T) {
	m := New()
	assert.False(t, m.IsFresh(5*time.Second))

	m.SetMetadataLocationAndFocused("/path/to/file", true)

	assert.Equal(t, "/path/to/file", m.GetMetadataLocation())
	assert.Equal(t, true, m.GetMetadataExpectedFocused())
	assert.True(t, m.IsFresh(5*time.Second))
}
