package utils

import (
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func TestResolveAbsPath(t *testing.T) {
	testdata := []struct {
		name        string
		cwd         string
		path        string
		expectedRes string
	}{
		{
			name:        "Path cleaup Test 1",
			cwd:         "/",
			path:        "////",
			expectedRes: "/",
		},
		{
			name:        "Basic test",
			cwd:         "/abc",
			path:        "def",
			expectedRes: "/abc/def",
		},
		{
			name:        "Ignore cwd for abs path",
			cwd:         "/abc",
			path:        "/def",
			expectedRes: "/def",
		},
		{
			name:        "Path cleanup Test 2",
			cwd:         "///abc",
			path:        "./././def",
			expectedRes: "/abc/def",
		},
		{
			name:        "Basic test with ~",
			cwd:         "/",
			path:        "~",
			expectedRes: xdg.Home,
		},
		{
			name:        "~ should not be resolved if not first",
			cwd:         "abc",
			path:        "x/~",
			expectedRes: "abc/x/~",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedRes, ResolveAbsPath(tt.cwd, tt.path))
		})
	}
}
