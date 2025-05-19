package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckFileNameValidity(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{name: "Invalid - single dot",
			input:   ".",
			wantErr: true,
			errMsg:  "file name cannot be '.' or '..'",
		}, {
			name:    "invalid - double dot",
			input:   "..",
			wantErr: true,
			errMsg:  "file name cnnot be '.' or '..'",
		}, {
			name:    "invalid - ends with /.. (platform separator)",
			input:   fmt.Sprintf("testDir%c..", filepath.Separator),
			wantErr: true,
			errMsg:  "file name cannot end with '/.' or '/..'",
		}, {
			name:    "invalid - ends with /. (platform separator)",
			input:   fmt.Sprintf("testDir%c.", filepath.Separator),
			wantErr: true,
			errMsg:  "file name cannot end with '/.' or '/..'",
		}, {
			name:    "valid - normal file name",
			input:   "valid_file.txt",
			wantErr: false,
		},
		{
			name:    "valid - contains dot inside",
			input:   "some.folder.name/file.txt",
			wantErr: false,
		},
		{
			name:    "valid - ends with dot not after separator",
			input:   "somefile.",
			wantErr: false,
		},
		{
			name:    "valid - ends with .. not after separator",
			input:   "somefile..",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkFileNameValidity(tt.input)
			if (err != nil) != tt.wantErr {
				if (err != nil) != tt.wantErr {
					t.Errorf("checkFileNameValidity(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				}
				if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("checkFileNameValidity(%q) error = %q, want to contain %q", tt.input, err.Error(), tt.errMsg)
				}
			}
		})
	}

}
