package processbar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetSortedProcessesCompletionOrder checks the ordering contract documented
// on getSortedProcesses: unfinished processes come first, ordered by completion
// (least progressed first). A directory-only operation has Total == 0; Render()
// displays it at 100%, so it must rank as fully complete here too.
//
// The processes live in a map, whose iteration order is randomized. The order is
// repeated with freshly built maps so that the result must be correct regardless
// of the source iteration order.
func TestGetSortedProcessesCompletionOrder(t *testing.T) {
	tests := []struct {
		name      string
		processes map[string]Process
		expected  []string
	}{
		{
			name: "directory-only process ranks as fully complete",
			processes: map[string]Process{
				"emptyDir": {ID: "emptyDir", State: InOperation, Done: 0, Total: 0},
				"quarter":  {ID: "quarter", State: InOperation, Done: 1, Total: 4},
				"threeQ":   {ID: "threeQ", State: InOperation, Done: 3, Total: 4},
			},
			expected: []string{"quarter", "threeQ", "emptyDir"},
		},
		{
			name: "incomplete process ranks before a completed one",
			processes: map[string]Process{
				"done":    {ID: "done", State: Successful, Done: 2, Total: 2},
				"running": {ID: "running", State: InOperation, Done: 1, Total: 4},
			},
			expected: []string{"running", "done"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for attempt := 0; attempt < 50; attempt++ {
				m := Model{processes: tt.processes}
				got := m.getSortedProcesses()
				require.Len(t, got, len(tt.expected))
				ids := make([]string, len(got))
				for i, p := range got {
					ids[i] = p.ID
				}
				assert.Equal(t, tt.expected, ids,
					"attempt %d: ordering must not depend on map iteration order", attempt)
			}
		})
	}
}
