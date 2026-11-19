package internal

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golift.io/xtractr"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestExtractCompressFile(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "archive."+format)
			dest := filepath.Join(root, "output")
			files := map[string]string{
				"file.txt":      "first file",
				"nested/中文.txt": "nested file",
			}
			writeExtractionArchive(t, archive, format, files)

			require.NoError(t, extractCompressFile(archive, dest, extractionProcessBar(t)))
			for name, want := range files {
				content, err := os.ReadFile(filepath.Join(dest, filepath.FromSlash(name)))
				require.NoError(t, err)
				require.Equal(t, want, string(content))
			}
		})
	}
}

func TestExtractCompressFileRejectsSiblingPathEscape(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "archive."+format)
			dest := filepath.Join(root, "output")
			outside := filepath.Join(root, "output-sibling", "keep.txt")
			utils.SetupDirectories(t, filepath.Dir(outside))
			utils.SetupFilesWithData(t, []byte("keep me"), outside)
			writeExtractionArchive(t, archive, format, map[string]string{
				"../output-sibling/keep.txt": "overwritten",
			})

			err := extractCompressFile(archive, dest, extractionProcessBar(t))
			require.ErrorIs(t, err, xtractr.ErrInvalidPath)
			content, err := os.ReadFile(outside)
			require.NoError(t, err)
			require.Equal(t, "keep me", string(content))
		})
	}
}

func TestExtractCompressFileRejectsOutputSymlink(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "archive."+format)
			dest := filepath.Join(root, "output")
			outside := filepath.Join(root, "outside", "keep.txt")
			utils.SetupDirectories(t, dest, filepath.Dir(outside))
			utils.SetupFilesWithData(t, []byte("keep me"), outside)
			if err := os.Symlink(filepath.Dir(outside), filepath.Join(dest, "link")); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			writeExtractionArchive(t, archive, format, map[string]string{
				"link/keep.txt": "overwritten",
			})

			err := extractCompressFile(archive, dest, extractionProcessBar(t))
			require.ErrorIs(t, err, xtractr.ErrInvalidPath)
			content, err := os.ReadFile(outside)
			require.NoError(t, err)
			require.Equal(t, "keep me", string(content))
		})
	}
}

func extractionProcessBar(t *testing.T) *processbar.Model {
	t.Helper()
	model := processbar.New()
	model.ListenForChannelUpdates()
	t.Cleanup(model.SendStopListeningMsgBlocking)
	return &model
}

func writeExtractionArchive(t *testing.T, path, format string, files map[string]string) {
	t.Helper()
	var data bytes.Buffer
	switch format {
	case "zip":
		writer := zip.NewWriter(&data)
		for name, content := range files {
			file, err := writer.Create(name)
			require.NoError(t, err)
			_, err = io.WriteString(file, content)
			require.NoError(t, err)
		}
		require.NoError(t, writer.Close())
	case "tar":
		writer := tar.NewWriter(&data)
		for name, content := range files {
			require.NoError(t, writer.WriteHeader(&tar.Header{
				Name: name,
				Mode: utils.ExtractedFileMode,
				Size: int64(len(content)),
			}))
			_, err := io.WriteString(writer, content)
			require.NoError(t, err)
		}
		require.NoError(t, writer.Close())
	default:
		t.Fatalf("unsupported archive format: %s", format)
	}
	require.NoError(t, os.WriteFile(path, data.Bytes(), utils.ExtractedFileMode))
}
