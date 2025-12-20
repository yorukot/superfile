package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_filePanelUpDown(t *testing.T) {
	testdata := []struct {
		name            string
		panel           FilePanel
		listDown        bool
		mainPanelHeight int
		expectedCursor  int
		expectedRender  int
	}{
		{
			name: "Down movement within renderable range",
			panel: FilePanel{
				Element:     make([]Element, 10),
				Cursor:      0,
				RenderIndex: 0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  1,
			expectedRender:  0,
		},
		{
			name: "Down movement when cursor is at bottom",
			panel: FilePanel{
				Element:     make([]Element, 10),
				Cursor:      6, // 3 - Header lines + 7(0-6 files)
				RenderIndex: 0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  7,
			expectedRender:  1,
		},
		{
			name: "Down movement causing wrap to top",
			panel: FilePanel{
				Element:     make([]Element, 10),
				Cursor:      9, // 3 - Header lines + 7(3-9 files)
				RenderIndex: 3,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
		{
			name: "Up movement within renderable range",
			panel: FilePanel{
				Element:     make([]Element, 10),
				Cursor:      2,
				RenderIndex: 0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  1,
			expectedRender:  0,
		},
		{
			name: "Up movement when cursor is at top",
			panel: FilePanel{
				Element:     make([]Element, 10),
				Cursor:      3,
				RenderIndex: 3,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  2,
			expectedRender:  2,
		},
		{
			name: "Up movement causing wrap to bottom",
			panel: FilePanel{
				Element:     make([]Element, 10),
				Cursor:      0,
				RenderIndex: 0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  9,
			expectedRender:  3,
		},
		{
			name: "Down movement on empty panel",
			panel: FilePanel{
				Element:     make([]Element, 0),
				Cursor:      0,
				RenderIndex: 0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
		{
			name: "Up movement on empty panel",
			panel: FilePanel{
				Element:     make([]Element, 0),
				Cursor:      0,
				RenderIndex: 0,
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
			assert.Equal(t, tt.expectedCursor, tt.panel.Cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.RenderIndex)
		})
	}
}

// TODO : Write tests for File Panel pgUp and pgDown and itemSelectUp/itemSelectDown
