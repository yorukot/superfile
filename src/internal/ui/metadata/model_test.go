package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpDown(t *testing.T) {
	defaultMetadata := Metadata{
		data: make([][2]string, 5),
	}
	testdata := []struct {
		name                string
		m                   Model
		listDown            bool // Whether to do listDown or listUp
		expectedRenderIndex int
	}{
		{
			name: "Basic down movement 1",
			m: Model{
				metadata:    defaultMetadata,
				renderIndex: 0,
			},
			listDown:            true,
			expectedRenderIndex: 1,
		},
		{
			name: "Down wraps to top",
			m: Model{
				metadata:    defaultMetadata,
				renderIndex: 4,
			},
			listDown:            true,
			expectedRenderIndex: 0,
		},
		{
			name: "Basic up movement 1",
			m: Model{
				metadata:    defaultMetadata,
				renderIndex: 4,
			},
			listDown:            false,
			expectedRenderIndex: 3,
		},
		{
			name: "Up wraps to top",
			m: Model{
				metadata:    defaultMetadata,
				renderIndex: 0,
			},
			listDown:            false,
			expectedRenderIndex: 4,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.m.ListDown()
			} else {
				tt.m.ListUp()
			}
			assert.Equal(t, tt.expectedRenderIndex, tt.m.renderIndex)
		})
	}
}
