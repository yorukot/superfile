package utils

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveAbsPath(t *testing.T) {
	_, current_file, _, _ := runtime.Caller(0)
	current_file_dir := filepath.Dir(current_file)
	current_file_name := filepath.Base(current_file)
	testdata := []struct{
		name string 
		cwd string
		path string
		expectedRes string
		errorExpected bool
	}{
		{
			name : "Basic Test",
			cwd : "/",
			path : "////",
			expectedRes: "/",
			errorExpected: false,
		},
		{
			name: "non existent file",
			cwd : current_file_dir,
			path : "non_existent_file",
			expectedRes: "",
			errorExpected: true,
		},
		{
			name: "existing file",
			cwd : current_file_dir,
			path : current_file_name,
			expectedRes: current_file,
			errorExpected: false,
		},
		{
			name: "Path cleanup",
			cwd : "///" + current_file_dir,
			path : "./././" + current_file_name,
			expectedRes: current_file,
			errorExpected: false,
		},
		{
			name : "Basic test with ~",
			cwd : "/",
			path : "~",
			expectedRes: xdg.Home,
			errorExpected: false,
		},
		{
			name : "Non abs cwd",
			cwd : "abc",
			path : "~",
			expectedRes: "",
			errorExpected: true,
		},
		
	}
	

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T){
			res, err := ResolveAbsPath(tt.cwd, tt.path);
			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err);
			}
			assert.Equal(t, tt.expectedRes, res)
			
		})
	}
}