package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildOpenArgs(t *testing.T) {
	tests := []struct {
		name        string
		openCommand string
		filePath    string
		wantProgram string
		wantArgs    []string
	}{
		{
			name:        "program without arguments",
			openCommand: "zathura",
			filePath:    "/home/user/doc.pdf",
			wantProgram: "zathura",
			wantArgs:    []string{"/home/user/doc.pdf"},
		},
		{
			name:        "program with a single argument",
			openCommand: "mpv --force-window=yes",
			filePath:    "/home/user/song.mp3",
			wantProgram: "mpv",
			wantArgs:    []string{"--force-window=yes", "/home/user/song.mp3"},
		},
		{
			name:        "program with several arguments keeps their order",
			openCommand: "feh --start-at --auto-zoom",
			filePath:    "/home/user/pic.png",
			wantProgram: "feh",
			wantArgs:    []string{"--start-at", "--auto-zoom", "/home/user/pic.png"},
		},
		{
			name:        "extra whitespace between tokens is ignored",
			openCommand: "  mpv   --force-window=yes  ",
			filePath:    "/home/user/song.mp3",
			wantProgram: "mpv",
			wantArgs:    []string{"--force-window=yes", "/home/user/song.mp3"},
		},
		{
			name:        "file path is appended after an absolute program path",
			openCommand: "/usr/bin/zathura",
			filePath:    "doc.pdf",
			wantProgram: "/usr/bin/zathura",
			wantArgs:    []string{"doc.pdf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, args, err := buildOpenArgs(tt.openCommand, tt.filePath)
			require.NoError(t, err)
			assert.Equal(t, tt.wantProgram, program)
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func TestBuildOpenArgsEmptyCommand(t *testing.T) {
	for _, openCommand := range []string{"", "   ", "\t\n"} {
		program, args, err := buildOpenArgs(openCommand, "file.txt")
		require.Error(t, err)
		assert.Empty(t, program)
		assert.Nil(t, args)
	}
}
