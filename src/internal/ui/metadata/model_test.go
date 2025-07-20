package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpDown(t *testing.T) {
	testdata := []struct {
		name               string
		m                  Model
		listDown           bool // Whether to do listDown or listUp
		expectedRendeIndex int
	}{
		{
			name: "Basic down movement 1",
			m: Model{
				metadata:    make([][2]string, 5),
				renderIndex: 0,
			},
			listDown:           true,
			expectedRendeIndex: 1,
		},
		{
			name: "Down wraps to top",
			m: Model{
				metadata:    make([][2]string, 5),
				renderIndex: 4,
			},
			listDown:           true,
			expectedRendeIndex: 0,
		},
		{
			name: "Basic up movement 1",
			m: Model{
				metadata:    make([][2]string, 5),
				renderIndex: 4,
			},
			listDown:           false,
			expectedRendeIndex: 3,
		},
		{
			name: "Up wraps to top",
			m: Model{
				metadata:    make([][2]string, 5),
				renderIndex: 0,
			},
			listDown:           false,
			expectedRendeIndex: 4,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.m.ListDown()
			} else {
				tt.m.ListUp()
			}
			assert.Equal(t, tt.expectedRendeIndex, tt.m.renderIndex)
		})
	}
}
