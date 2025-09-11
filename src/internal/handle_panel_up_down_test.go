package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_filePanelUpDown(t *testing.T) {
	testdata := []struct {
		name            string
		panel           filePanel
		listDown        bool
		mainPanelHeight int
		expectedCursor  int
		expectedRender  int
	}{
		{
			name: "Down movement within renderable range",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  0,
				render:  0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  1,
			expectedRender:  0,
		},
		{
			name: "Down movement when cursor is at bottom",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  6, // 3 - Header lines + 7(0-6 files)
				render:  0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  7,
			expectedRender:  1,
		},
		{
			name: "Down movement causing wrap to top",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  9, // 3 - Header lines + 7(3-9 files)
				render:  3,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
		{
			name: "Up movement within renderable range",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  2,
				render:  0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  1,
			expectedRender:  0,
		},
		{
			name: "Up movement when cursor is at top",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  3,
				render:  3,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  2,
			expectedRender:  2,
		},
		{
			name: "Up movement causing wrap to bottom",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  0,
				render:  0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  9,
			expectedRender:  3,
		},
		{
			name: "Down movement on empty panel",
			panel: filePanel{
				element: make([]element, 0),
				cursor:  0,
				render:  0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
		{
			name: "Up movement on empty panel",
			panel: filePanel{
				element: make([]element, 0),
				cursor:  0,
				render:  0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.panel.ListDown(tt.mainPanelHeight)
			} else {
				tt.panel.ListUp(tt.mainPanelHeight)
			}
			assert.Equal(t, tt.expectedCursor, tt.panel.cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.render)
		})
	}
}

// TODO : Write tests for File Panel pgUp and pgDown and itemSelectUp/itemSelectDown
